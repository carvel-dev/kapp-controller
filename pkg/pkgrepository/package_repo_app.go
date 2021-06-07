// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func NewPackageRepoApp(pkgRepository *pkgingv1alpha1.PackageRepository) (*kcv1alpha1.App, error) {
	desiredApp := &kcv1alpha1.App{}

	desiredApp.Name = pkgRepository.Name
	desiredApp.Namespace = pkgRepository.Namespace
	desiredApp.DeletionTimestamp = pkgRepository.DeletionTimestamp
	desiredApp.Generation = pkgRepository.Generation

	desiredApp.Spec = kcv1alpha1.AppSpec{
		Fetch: []kcv1alpha1.AppFetch{{
			Image:        pkgRepository.Spec.Fetch.Image,
			Git:          pkgRepository.Spec.Fetch.Git,
			HTTP:         pkgRepository.Spec.Fetch.HTTP,
			ImgpkgBundle: pkgRepository.Spec.Fetch.ImgpkgBundle,
		}},
		Template: []kcv1alpha1.AppTemplate{{
			Ytt: &kcv1alpha1.AppTemplateYtt{
				IgnoreUnknownComments: true,
				Paths:                 []string{"packages"},
			},
		}},
		Deploy: []kcv1alpha1.AppDeploy{{
			Kapp: &kcv1alpha1.AppDeployKapp{},
		}},
		Paused: pkgRepository.Spec.Paused,
	}

	if pkgRepository.Spec.SyncPeriod == nil {
		desiredApp.Spec.SyncPeriod = &metav1.Duration{Duration: time.Minute * 5}
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
