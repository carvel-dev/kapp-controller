// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installedpkg

import (
	instpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/package/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/scheme"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func NewApp(existingApp *v1alpha1.App, installedPkg *instpkgv1alpha1.InstalledPackage, pkg pkgv1alpha1.Package) (*v1alpha1.App, error) {
	desiredApp := existingApp.DeepCopy()

	desiredApp.Name = installedPkg.Name
	desiredApp.Namespace = installedPkg.Namespace
	desiredApp.Spec = pkg.Spec.Template.Spec
	desiredApp.Spec.ServiceAccountName = installedPkg.Spec.ServiceAccountName

	err := controllerutil.SetControllerReference(installedPkg, desiredApp, scheme.Scheme)
	if err != nil {
		return &v1alpha1.App{}, err
	}

	for i, templateStep := range desiredApp.Spec.Template {
		if templateStep.HelmTemplate != nil {
			for _, value := range installedPkg.Spec.Values {
				templateStep.HelmTemplate.ValuesFrom = append(templateStep.HelmTemplate.ValuesFrom, kcv1alpha1.AppTemplateHelmTemplateValuesSource{
					SecretRef: &kcv1alpha1.AppTemplateHelmTemplateValuesSourceRef{
						LocalObjectReference: corev1.LocalObjectReference{Name: value.SecretRef.Name},
					},
				})
			}
			desiredApp.Spec.Template[i] = templateStep
			break
		}

		if templateStep.Ytt != nil {
			if templateStep.Ytt.Inline == nil {
				templateStep.Ytt.Inline = &kcv1alpha1.AppFetchInline{}
			}

			for _, value := range installedPkg.Spec.Values {
				templateStep.Ytt.Inline.PathsFrom = append(templateStep.Ytt.Inline.PathsFrom, kcv1alpha1.AppFetchInlineSource{
					SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
						LocalObjectReference: corev1.LocalObjectReference{Name: value.SecretRef.Name},
					},
				})
			}
			desiredApp.Spec.Template[i] = templateStep
			break
		}
	}

	return desiredApp, nil
}
