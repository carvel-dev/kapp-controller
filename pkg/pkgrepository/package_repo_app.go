// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"fmt"
	"os"
	"time"

	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	kappDebug = os.Getenv("KAPPCTRL_PKGR_KAPP_DEBUG") == "true"
)

func NewPackageRepoApp(pkgRepository *pkgingv1alpha1.PackageRepository) (*kcv1alpha1.App, error) {
	desiredApp := &kcv1alpha1.App{}

	desiredApp.Name = pkgRepository.Name
	desiredApp.Namespace = pkgRepository.Namespace
	desiredApp.DeletionTimestamp = pkgRepository.DeletionTimestamp
	desiredApp.Generation = pkgRepository.Generation

	kappRawOpts := []string{
		"--kube-api-qps=30",
		"--kube-api-burst=40",
		// Default kapp-controller service account allows listing all namespaces
		// but does not allow listing of most resources in such namespaces --
		// instead of spinning wheels trying to list, scope to "current" namespace
		"--dangerous-scope-to-fallback-allowed-namespaces=true",
	}

	if kappDebug {
		kappRawOpts = append(kappRawOpts, "--debug=true")
	}

	kappDeployRawOpts := append([]string{
		"--logs=false",
		"--app-changes-max-to-keep=5",
		// in the case where two packages in different PKGRs are identical,
		// first we have to tell kapp it's ok if the new PKGR takes ownership
		// but then we use the rebase rule (below) to actually not override the ownership
		// but without this flag the rebase rule won't get a chance to run. So in conclusion
		// it's a fake-out and we shouldn't actually take ownership of existing resources.
		"--dangerous-override-ownership-of-existing-resources=true",
		// GKE for some reason does not like high volume of GETs for our API server
		// and ends up taking very long time to respond to SubjectAccessReviews (SAR).
		// We can disable existing check entirely since we use a rebase rule to enforce
		// non-transference of ownership in the case that we have multiple packages of
		// the same name.
		"--existing-non-labeled-resources-check=false",
		// ... could in theory just lower concurrency, but decided to turn it off entirely.
		// (on GKE, 6 was a sweet spot, 10 exhibited hanging behaviour)
		// "--existing-non-labeled-resources-check-concurrency=6",
	}, kappRawOpts...)

	kappDeleteRawOpts := append([]string{}, kappRawOpts...)

	desiredApp.Spec = kcv1alpha1.AppSpec{
		Fetch: []kcv1alpha1.AppFetch{{
			Image:        pkgRepository.Spec.Fetch.Image,
			Inline:       pkgRepository.Spec.Fetch.Inline,
			Git:          pkgRepository.Spec.Fetch.Git,
			HTTP:         pkgRepository.Spec.Fetch.HTTP,
			ImgpkgBundle: pkgRepository.Spec.Fetch.ImgpkgBundle,
		}},
		Template: []kcv1alpha1.AppTemplate{{
			// we have multiple ytt sections because we want to squash all the user yamls together
			// and then apply our overlays. This way we do multiple distinct ytt passes.
			Ytt: &kcv1alpha1.AppTemplateYtt{
				IgnoreUnknownComments: true,
				Paths:                 []string{"packages"},
			},
		}, {
			Ytt: &kcv1alpha1.AppTemplateYtt{
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
    kapp.k14s.io/create-strategy: "fallback-on-update-or-noop"`, pkgRepository.Namespace, pkgRepository.Name),
					},
				},
			},
		}, {
			Ytt: &kcv1alpha1.AppTemplateYtt{
				Paths: []string{"-"},
				Inline: &kcv1alpha1.AppFetchInline{
					Paths: map[string]string{
						//   rebase rule to allow multiple repositories to expose identical packages.
						"noop-on-identical-packages.yml": `
---
apiVersion: kapp.k14s.io/v1alpha1
kind: Config
rebaseRules:
- ytt:
    overlayContractV1:
      overlay.yml: |
        #@ load("@ytt:data", "data")
        #@ load("@ytt:yaml", "yaml")
        #@ load("@ytt:json", "json")
        #@ load("@ytt:overlay", "overlay")
        #@ load("@ytt:struct", "struct")

        #@ def get_rev(annotations):
        #@   if hasattr(annotations, "packaging.carvel.dev/revision"):
        #@     return int(annotations["packaging.carvel.dev/revision"])
        #@   else:
        #@     return -1
        #@   end
        #@ end
        #@ if not hasattr(data.values.existing.metadata.annotations, "packaging.carvel.dev/package-repository-ref"):
        #@   msg = "Error: cannot overwrite package " + data.values.existing.metadata.name + " because it was not created by a package repository."
        #@   print(msg)
        #@   fail(msg)
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
        #@ def specs_are_identical(existing_spec, new_spec, kind):
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
        #@     if nw[attr] != ex[attr]:
        #@       return False, "mismatch in spec." + attr
        #@     end
        #@   end
        #@   return True, ""
        #@ end
        #@
        #@ def is_identical(existing, new):
        #@   eq, reason = specs_are_identical(struct.decode(existing.spec), struct.decode(new.spec), existing.kind)
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

        #@ new_owner = data.values.new.metadata.annotations["packaging.carvel.dev/package-repository-ref"]
        #@ existing_owner = data.values.existing.metadata.annotations["packaging.carvel.dev/package-repository-ref"]
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
        #@ elif get_rev(data.values.existing.metadata.annotations) > get_rev(data.values.new.metadata.annotations):
        #@overlay/match by=overlay.all
        ---
        metadata:
          #@overlay/match missing_ok=True
          annotations:
            #@overlay/match missing_ok=True
            kapp.k14s.io/noop: ""
        #@   elif get_rev(data.values.existing.metadata.annotations) < get_rev(data.values.new.metadata.annotations):
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
`,
					},
				},
			},
		}},
		Deploy: []kcv1alpha1.AppDeploy{{
			Kapp: &kcv1alpha1.AppDeployKapp{
				RawOptions: kappDeployRawOpts,
				Delete: &kcv1alpha1.AppDeployKappDelete{
					RawOptions: kappDeleteRawOpts,
				},
			},
		}},
		Paused: pkgRepository.Spec.Paused,
	}

	if pkgRepository.Spec.SyncPeriod == nil {
		desiredApp.Spec.SyncPeriod = &metav1.Duration{Duration: time.Minute * 10}
	} else {
		desiredApp.Spec.SyncPeriod = pkgRepository.Spec.SyncPeriod
	}

	if desiredApp.Spec.Fetch[0].ImgpkgBundle != nil {
		desiredApp.Spec.Template = append(desiredApp.Spec.Template,
			kcv1alpha1.AppTemplate{Kbld: &kcv1alpha1.AppTemplateKbld{Paths: []string{"-", ".imgpkg/images.yml"}}})
	}

	desiredApp.Status = kcv1alpha1.AppStatus{
		Fetch:                         pkgRepository.Status.Fetch,
		Template:                      pkgRepository.Status.Template,
		Deploy:                        pkgRepository.Status.Deploy,
		GenericStatus:                 pkgRepository.Status.GenericStatus,
		ConsecutiveReconcileSuccesses: pkgRepository.Status.ConsecutiveReconcileSuccesses,
		ConsecutiveReconcileFailures:  pkgRepository.Status.ConsecutiveReconcileFailures,
	}

	return desiredApp, nil
}
