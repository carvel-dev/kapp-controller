// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	instpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/scheme"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	appNs = "kapp-controller"
)

func NewApp(existingApp *kcv1alpha1.App, pkgRepository *instpkgv1alpha1.PackageRepository) (*kcv1alpha1.App, error) {
	desiredApp := existingApp.DeepCopy()

	desiredApp.Name = pkgRepository.Name
	desiredApp.Namespace = appNs

	err := controllerutil.SetControllerReference(pkgRepository, desiredApp, scheme.Scheme)
	if err != nil {
		return &kcv1alpha1.App{}, err
	}

	desiredApp.Spec = kcv1alpha1.AppSpec{
		// TODO since we are assuming that we are inside kapp-controller NS, use its SA
		ServiceAccountName: "kapp-controller-sa",
		Fetch: []kcv1alpha1.AppFetch{{
			Image:        pkgRepository.Spec.Fetch.Image,
			Git:          pkgRepository.Spec.Fetch.Git,
			HTTP:         pkgRepository.Spec.Fetch.HTTP,
			ImgpkgBundle: pkgRepository.Spec.Fetch.Bundle,
		}},
		Template: []kcv1alpha1.AppTemplate{{
			Ytt: &kcv1alpha1.AppTemplateYtt{
				IgnoreUnknownComments: true,
			},
		}},
		Deploy: []kcv1alpha1.AppDeploy{{
			Kapp: &kcv1alpha1.AppDeployKapp{},
		}},
	}

	if desiredApp.Spec.Fetch[0].ImgpkgBundle != nil {
		desiredApp.Spec.Template = append(desiredApp.Spec.Template,
			kcv1alpha1.AppTemplate{Kbld: &kcv1alpha1.AppTemplateKbld{Paths: []string{"-", ".imgpkg/images.yml"}}})
	}

	return desiredApp, nil
}
