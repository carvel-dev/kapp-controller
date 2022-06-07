// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	datapackagingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	ctltpl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	otherkyaml "sigs.k8s.io/yaml" // TODO: dude seriously there's so many yamls what would you call this one
)

var rebaseRule string = `
        #@ load("@ytt:data", "data")
        #@ load("@ytt:yaml", "yaml")
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

        #! return 0 iff eq, 1 (or more) iff r1 is gt r2, -1(or less) iff r1 < r2
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

        #@ def filter(d, s):
        #@   return {x: v for x, v in d.items() if not x in s}
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
        #@
        #@ def specs_are_identical(existing_spec, new_spec):
        #@   if str(new_spec) == str(existing_spec):
        #@     return True, ""
        #@   end
        #@
        #!   the rest of this method is to help make a better error message
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
        #@
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
        #@
        #@   return True, ""
        #@ end

        #! if the pkgr-ref annotation is missing (and the packages are identical) then assume ownership - covers upgrade case from old kcs
        #@ new_owner = data.values.new.metadata.annotations["packaging.carvel.dev/package-repository-ref"]
        #@ if 'packaging.carvel.dev/package-repository-ref' in data.values.existing.metadata.annotations:
        #@   existing_owner = data.values.existing.metadata.annotations["packaging.carvel.dev/package-repository-ref"]
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
        #@     print("replacing existing older rev with newer rev")
        #@   else:
        #@     msg = "Error: Conflicting Resources: " + data.values.existing.kind + "/" + data.values.existing.metadata.name + " is already present but not identical (" + reason +")"
        #@     print(msg)
        #@     fail(msg)
        #@   end
        #@ end
  resourceMatchers:
  - apiVersionKindMatcher: {apiVersion: data.packaging.carvel.dev/v1alpha1, kind: Package}
  - apiVersionKindMatcher: {apiVersion: data.packaging.carvel.dev/v1alpha1, kind: PackageMetadata}
`

func (a *App) template(dirPath string) exec.CmdRunResult {
	// we have multiple ytt sections because we want to squash all the user yamls together
	// and then apply our overlays. This way we do multiple distinct ytt passes.
	template1 := kcv1alpha1.AppTemplateYtt{
		IgnoreUnknownComments: true,
		Paths:                 []string{"packages"},
	}
	template2 := kcv1alpha1.AppTemplateYtt{
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
	template3 := kcv1alpha1.AppTemplateYtt{
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
`, rebaseRule),
			},
		},
	}

	genericOpts := ctltpl.GenericOpts{Name: a.app.Name, Namespace: a.app.Namespace}

	// first templating pass is to read all the files in
	var result exec.CmdRunResult
	var template ctltpl.Template
	template = a.templateFactory.NewYtt(template1, genericOpts)
	result, _ = template.TemplateDir(dirPath)
	if result.Error != nil {
		return result
	}

	// second templating applies a bunch of overlays and filters, some of which could be migrated to go code
	template = a.templateFactory.NewYtt(template2, genericOpts)
	result = template.TemplateStream(strings.NewReader(result.Stdout), dirPath)
	if result.Error != nil {
		return result
	}

	// intermediate phase deserializes and reserializes all of the resources
	// which guarantees that they only include fields this version of kc knows about.
	resources, err := FilterResources(result.Stdout)
	if err != nil {
		result.Error = err
		return result
	}

	// third templating inserts a kapp rebase rule to allow noops on identical resources provided by multiple pkgrs
	template = a.templateFactory.NewYtt(template3, genericOpts)
	result = template.TemplateStream(strings.NewReader(resources), dirPath)

	return result
}

// FilterResources takes a multi-doc yaml of the templated
// contents of a PKGR, and filters out unexpected fields on each resource by deserializing and re-serializing
// This filtering step allows us to use newer CRDs (with new fields) in older versions of kc
// without triggering a mistmatch in the rebase "is_identical" checker.
func FilterResources(yamlss string) (string, error) {
	yamls := []byte(yamlss)
	sch := runtime.NewScheme()
	err := scheme.AddToScheme(sch)
	if err != nil {
		return "", err
	}
	err = datapackagingv1alpha1.AddToScheme(sch)
	if err != nil {
		return "", err
	}
	decoder := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	filteredYamls := []string{}
	docs, err := yamlDocs(yamls)
	if err != nil {
		return "", err
	}

	for _, resourceYAML := range docs {
		obj, gKV, err := decoder(resourceYAML, nil, nil)
		if err != nil {
			return "", err
		}
		kind := gKV.Kind
		if gKV.Group != "data.packaging.carvel.dev" {
			return "", fmt.Errorf("Expected group 'data.packaging.carvel.dev' but was: %s", gKV.Group)
		}

		if gKV.Version != "v1alpha1" {
			return "", fmt.Errorf("Expected version 'v1alpha1' but was: %s", gKV.Version)
		}

		switch kind {
		case "Package":
			p := obj.(*datapackagingv1alpha1.Package)
			buf, err := otherkyaml.Marshal(p)
			if err != nil {
				return "", err
			}
			filteredYamls = append(filteredYamls, string(buf))
		case "PackageMetadata":
			p := obj.(*datapackagingv1alpha1.PackageMetadata)
			buf, err := otherkyaml.Marshal(p)
			if err != nil {
				return "", err
			}
			filteredYamls = append(filteredYamls, string(buf))
		default:
			return "", fmt.Errorf("PKGR contained unexpected kind: %s", kind)
		}

	}
	return strings.Join(filteredYamls, "\n---\n"), nil
}

func yamlDocs(yamls []byte) ([][]byte, error) {
	var docs [][]byte

	fileBytes := yamls
	reader := kyaml.NewYAMLReader(bufio.NewReaderSize(bytes.NewReader(fileBytes), 4096))

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
