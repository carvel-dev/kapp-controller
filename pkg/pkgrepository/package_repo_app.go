// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"time"

	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
				// TODO do not want to interpret packages/**/* as templates
				// but we cannot apply file mark to both .yml and .yaml without
				// one of them possibly being 0 (that results in ytt error currently)
				IgnoreUnknownComments: true,
				Paths:                 []string{"packages"},

				Inline: &kcv1alpha1.AppFetchInline{
					Paths: map[string]string{
						// Remove all resources that are not known to this kapp-controller.
						// It's worth just removing instead of erroring,
						// since future repo bundles may introduce new kinds.
						"kapp-controller-clean-up.yml": `
#@ load("@ytt:overlay", "overlay")

#@ pkg = overlay.subset({"apiVersion":"data.packaging.carvel.dev/v1alpha1", "kind": "Package"})
#@ pkgm = overlay.subset({"apiVersion":"data.packaging.carvel.dev/v1alpha1", "kind": "PackageMetadata"})

#@overlay/match by=overlay.not_op(overlay.or_op(pkg, pkgm)),expects="0+"
#@overlay/remove
---

#@overlay/match by=overlay.or_op(pkg, pkgm),expects="0+"
---
metadata:
  #@overlay/match missing_ok=True
  annotations:
    #@overlay/match missing_ok=True
    kapp.k14s.io/disable-original: ""
    #@overlay/match missing_ok=True
    kapp.k14s.io/disable-wait: ""

#! Ensure that all resources do not set some random namespace
#! so that all resource end in the PackageRepository's namespace
#@overlay/match by=overlay.all,expects="0+"
---
metadata:
  #@overlay/match missing_ok=True
  #@overlay/remove
  namespace:
`,
					},
				},
			},
		}},
		Deploy: []kcv1alpha1.AppDeploy{{
			Kapp: &kcv1alpha1.AppDeployKapp{
				RawOptions: []string{"--wait-timeout=30s", "--kube-api-qps=20", "--kube-api-burst=30"},
				Delete: &kcv1alpha1.AppDeployKappDelete{
					RawOptions: []string{"--wait-timeout=30s", "--kube-api-qps=20", "--kube-api-burst=30"},
				},
			},
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
