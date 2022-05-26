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
		"--wait-timeout=30s",
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
		// GKE for some reason does not like high volume of GETs for our API server
		// and ends up taking very long time to respond to SubjectAccessReviews (SAR).
		// We can disable existing check entirely since we do not allow to "take ownership"
		// of other Package/PackageMetadata and we do not _need_ to adopt Packages that
		// are not owned by kapp already.
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
`, pkgRepository.Namespace, pkgRepository.Name),
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