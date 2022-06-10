// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
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
		Template: []kcv1alpha1.AppTemplate{}, // Template step hardcoded into app_template.go
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
