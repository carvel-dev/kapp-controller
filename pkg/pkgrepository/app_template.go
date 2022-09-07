// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	datapackagingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	ctltpl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	sigsyaml "sigs.k8s.io/yaml"
)

var errInvalidPackageRepo = exec.NewCmdRunResultWithErr(fmt.Errorf("Invalid package repository content: must contain 'packages/' directory but did not"))

func (a *App) template(dirPath string) exec.CmdRunResult {
	fileInfo, err := os.Lstat(filepath.Join(dirPath, "packages"))
	if err != nil {
		if os.IsNotExist(err) {
			return errInvalidPackageRepo
		}
		return exec.NewCmdRunResultWithErr(err)
	}
	if !fileInfo.IsDir() {
		return errInvalidPackageRepo
	}
	if len(a.app.Spec.Template) != 0 {
		panic("Internal inconsistency: Package repository templates are not configurable")
	}

	appContext := ctltpl.AppContext{Name: a.app.Name, Namespace: a.app.Namespace}

	// We have multiple ytt sections because we want to squash all the user yamls together
	// and then apply our overlays. This way we do multiple distinct ytt passes.

	// First templating pass is to read all the files in
	template1 := kcv1alpha1.AppTemplateYtt{
		IgnoreUnknownComments: true,
		Paths:                 []string{"packages"},
	}
	additionalValues := ctltpl.AdditionalDownwardAPIValues{}

	result, _ := a.templateFactory.NewYtt(template1, appContext, additionalValues).TemplateDir(dirPath)
	if result.Error != nil {
		return result
	}

	// Second templating applies a bunch of overlays and filters,
	// some of which could be migrated to go code
	stream := strings.NewReader(result.Stdout)
	result = a.templateFactory.NewYtt(
		a.yttTemplateCleanRs(), appContext, additionalValues).TemplateStream(stream, dirPath)
	if result.Error != nil {
		return result
	}

	// Intermediate phase deserializes and reserializes all of the resources
	// which guarantees that they only include fields this version of kc knows about.
	resources, err := FilterResources(result.Stdout)
	if err != nil {
		result.Error = err
		return result
	}

	// Third templating inserts a kapp rebase rule to allow noops
	// on identical resources provided by multiple pkgrs
	stream = strings.NewReader(resources)
	result = a.templateFactory.NewYtt(
		a.yttTemplateAddIdenticalRsRebase(), appContext, additionalValues).TemplateStream(stream, dirPath)
	if result.Error != nil {
		return result
	}

	// Optionally use kbld to apply .imgpkg/images.yml if content came from imgpkgBundle
	if a.app.Spec.Fetch[0].ImgpkgBundle != nil {
		stream = strings.NewReader(result.Stdout)
		kbldOpts := kcv1alpha1.AppTemplateKbld{Paths: []string{"-", ".imgpkg/images.yml"}}
		result = a.templateFactory.NewKbld(kbldOpts, appContext).TemplateStream(stream, dirPath)
	}

	return result
}

func (a *App) yttTemplateCleanRs() kcv1alpha1.AppTemplateYtt {
	return kcv1alpha1.AppTemplateYtt{
		Paths: []string{"-"},
		Inline: &kcv1alpha1.AppFetchInline{
			Paths: map[string]string{
				// - Adjust the contents of the repo including adding
				//   annotations and ensuring namespace.
				// - Remove all resources that are not known to this kapp-controller.
				//   It's worth just removing instead of erroring,
				//   since future repo bundles may introduce new kinds.
				"kapp-controller-clean-up.yml": fmt.Sprintf(`
#@ load("@ytt:overlay", "overlay")

#@ pkg = overlay.subset({"apiVersion":"data.packaging.carvel.dev/v1alpha1", "kind": "Package"})
#@ pkgm = overlay.subset({"apiVersion":"data.packaging.carvel.dev/v1alpha1", "kind": "PackageMetadata"})

#@overlay/match by=overlay.not_op(overlay.or_op(pkg, pkgm)),expects="0+"
#@overlay/remove
---

#@overlay/match by=overlay.all,expects="0+"
---
metadata:
  #! Ensure that all resources do not set some random namespace
  #! so that all resource end in the PackageRepository's namespace
  #@overlay/match missing_ok=True
  #@overlay/remove
  namespace:

  #@overlay/match missing_ok=True
  annotations:
    #@overlay/match missing_ok=True
    kapp.k14s.io/disable-original: ""

    #@overlay/match missing_ok=True
    kapp.k14s.io/disable-wait: ""

    #@overlay/match missing_ok=True
    packaging.carvel.dev/package-repository-ref: %s/%s

    #@overlay/match missing_ok=True
    kapp.k14s.io/create-strategy: "fallback-on-update-or-noop"`, a.Namespace(), a.Name()),
			},
		},
	}
}

func (a *App) yttTemplateAddIdenticalRsRebase() kcv1alpha1.AppTemplateYtt {
	var yttRebasePackageRelatedRsByRevision = `
        #@ load("@ytt:data", "data")
        #@ load("@ytt:json", "json")
        #@ load("@ytt:overlay", "overlay")
        #@ load("@ytt:struct", "struct")

        #@ def get_rev(annotations):
        #@   if hasattr(annotations, "packaging.carvel.dev/revision"):
        #@     return [int(x) for x in annotations["packaging.carvel.dev/revision"].split('.')]
        #@   else:
        #@     return [-1]
        #@   end
        #@ end

        #! return 0 iff eq, 1 (or more) iff r1 > r2, -1 (or less) iff r1 < r2
        #@ def cmp_rev(r1, r2):
        #@   size = min(len(r1), len(r2))
        #@   for i in range(size):
        #@     if r1[i] > r2[i]:
        #@        return 1
        #@     elif r1[i] < r2[i]:
        #@        return -1
        #@     end
        #@   end
        #@   return len(r1) - len(r2)
        #@ end

        #@ def filter(kvs, exclude_keys):
        #@   return {k: v for k, v in kvs.items() if not k in exclude_keys}
        #@ end

        #! TODO: next person who adds an annotation to the set of all kapp-controller annotations is gonna have a bad time.
        #@ def annotations_are_identical(existing_anns, new_anns):
        #@   carvel_anns = set(["kapp.k14s.io/disable-original",
        #@                      "kapp.k14s.io/disable-wait",
        #@                      "packaging.carvel.dev/package-repository-ref",
        #@                      "kapp.k14s.io/identity",
        #@                      "kapp.k14s.io/create-strategy"])
        #@   return filter(existing_anns, carvel_anns) == filter(new_anns, carvel_anns)
        #@ end

        #@ def labels_are_identical(existing_labels, new_labels):
        #@   carvel_labels = set(["kapp.k14s.io/app", "kapp.k14s.io/association"])
        #@   return filter(existing_labels, carvel_labels) == filter(new_labels, carvel_labels)
        #@ end

        #@ def specs_are_identical(existing_spec, new_spec):
        #@   if json.encode(new_spec) == json.encode(existing_spec):
        #@     return True, ""
        #@   end
        #@
        #@   # the rest of this method is to help make a better error message
        #@   ex = existing_spec
        #@   nw = new_spec
        #@
        #@   for attr in nw.keys():
        #@     if not attr in ex:
        #@       return False, "adding spec."+attr
        #@     end
        #@   end
        #@
        #@   for attr in ex.keys():
        #@     if not attr in nw:
        #@       return False, "missing spec." + attr
        #@     end
        #@     if str(nw[attr]) != str(ex[attr]):
        #@       return False, "mismatch in spec." + attr
        #@     end
        #@   end
        #@   return False, "mismatch in unknown location"
        #@ end

        #@ def is_identical(existing, new):
        #@   eq, reason = specs_are_identical(struct.decode(existing.spec), struct.decode(new.spec))
        #@   if not eq:
        #@     return False, reason
        #@   end
        #@   if not labels_are_identical(struct.decode(existing.metadata.labels), struct.decode(new.metadata.labels)):
        #@     return False, "mismatch in metadata.labels"
        #@   end
        #@   if not annotations_are_identical(struct.decode(existing.metadata.annotations), struct.decode(new.metadata.annotations)):
        #@     return False, "mismatch in metadata.annotations"
        #@   end
        #@   return True, ""
        #@ end

        #! if the pkgr-ref annotation is missing (and the packages are identical)
        #! then assume ownership - covers upgrade case from old kcs

        #@ pkg_repo_ann = "packaging.carvel.dev/package-repository-ref"
        #@ new_owner = data.values.new.metadata.annotations[pkg_repo_ann]
        #@
        #@ if pkg_repo_ann in data.values.existing.metadata.annotations:
        #@   existing_owner = data.values.existing.metadata.annotations[pkg_repo_ann]
        #@ else:
        #@   existing_owner = new_owner
        #@ end
        #@
        #@ if new_owner != existing_owner:
        #@   identical, reason = is_identical(data.values.existing, data.values.new)
        #@   if identical:

        #@overlay/match by=overlay.all
        ---
        metadata:
          #@overlay/match missing_ok=True
          annotations:
            #@overlay/match missing_ok=True
            kapp.k14s.io/noop: ""

        #@ elif cmp_rev(get_rev(data.values.existing.metadata.annotations), get_rev(data.values.new.metadata.annotations)) > 0:

        #@overlay/match by=overlay.all
        ---
        metadata:
          #@overlay/match missing_ok=True
          annotations:
            #@overlay/match missing_ok=True
            kapp.k14s.io/noop: ""

        #@   elif cmp_rev(get_rev(data.values.existing.metadata.annotations), get_rev(data.values.new.metadata.annotations)) < 0:
        #@     # replacing existing older rev with newer rev
        #@   else:
        #@     fail("Error: Conflicting resources: " + data.values.existing.kind + "/" + data.values.existing.metadata.name + " is already present but not identical (" + reason +")")
        #@   end
        #@ end
`

	return kcv1alpha1.AppTemplateYtt{
		Paths: []string{"-"},
		Inline: &kcv1alpha1.AppFetchInline{
			Paths: map[string]string{
				//   rebase rule to allow multiple repositories to expose identical packages.
				"noop-on-identical-packages.yml": fmt.Sprintf(`
---
apiVersion: kapp.k14s.io/v1alpha1
kind: Config
rebaseRules:
- ytt:
    overlayContractV1:
      overlay.yml: |
%s
  resourceMatchers:
  - apiVersionKindMatcher: {apiVersion: data.packaging.carvel.dev/v1alpha1, kind: Package}
  - apiVersionKindMatcher: {apiVersion: data.packaging.carvel.dev/v1alpha1, kind: PackageMetadata}
`, yttRebasePackageRelatedRsByRevision),
			},
		},
	}
}

// FilterResources takes a multi-doc yaml of the templated
// contents of a PKGR, and filters out unexpected fields on each resource
// by deserializing and re-serializing. This filtering step allows us
// to use newer CRDs (with new fields) in older versions of kc
// without triggering a mistmatch in the rebase "is_identical" checker.
func FilterResources(inputYAML string) (string, error) {
	sch := runtime.NewScheme()

	err := datapackagingv1alpha1.AddToScheme(sch)
	if err != nil {
		return "", err
	}
	deserializer := serializer.NewCodecFactory(sch).UniversalDeserializer()

	docs, err := yamlDocs([]byte(inputYAML))
	if err != nil {
		return "", fmt.Errorf("Parsing stream of resources: %s", err)
	}

	filteredYAMLs := []string{}

	for _, resourceYAML := range docs {
		obj, gvk, err := deserializer.Decode(resourceYAML, nil, nil)
		if err != nil {
			return "", fmt.Errorf("Deserializing resource: %s", err)
		}

		if gvk.Group != datapackagingv1alpha1.SchemeGroupVersion.Group {
			return "", fmt.Errorf("Expected group 'data.packaging.carvel.dev' but was '%s'", gvk.Group)
		}
		if gvk.Version != datapackagingv1alpha1.SchemeGroupVersion.Version {
			return "", fmt.Errorf("Expected version 'v1alpha1' but was '%s'", gvk.Version)
		}
		if gvk.Kind != "Package" && gvk.Kind != "PackageMetadata" {
			return "", fmt.Errorf("Expected kind to be 'Package' or 'PackageMetadata' but was '%s'", gvk.Kind)
		}

		buf, err := sigsyaml.Marshal(obj)
		if err != nil {
			return "", fmt.Errorf("Marshaling resource: %s", err)
		}
		filteredYAMLs = append(filteredYAMLs, string(buf))
	}

	return strings.Join(filteredYAMLs, "\n---\n"), nil
}

func yamlDocs(yamls []byte) ([][]byte, error) {
	var docs [][]byte
	reader := utilyaml.NewYAMLReader(bufio.NewReaderSize(bytes.NewReader(yamls), 4096))

	for {
		docBytes, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		docs = append(docs, docBytes)
	}

	return docs, nil
}
