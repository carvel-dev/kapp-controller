// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installedpkg

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/scheme"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func NewApp(existingApp *v1alpha1.App, installedPkg *pkgingv1alpha1.InstalledPackage, pkgVersion datapkgingv1alpha1.PackageVersion) (*v1alpha1.App, error) {
	desiredApp := existingApp.DeepCopy()

	desiredApp.Name = installedPkg.Name
	desiredApp.Namespace = installedPkg.Namespace
	desiredApp.Spec = *pkgVersion.Spec.Template.Spec
	desiredApp.Spec.ServiceAccountName = installedPkg.Spec.ServiceAccountName
	desiredApp.Spec.SyncPeriod = installedPkg.Spec.SyncPeriod
	desiredApp.Spec.NoopDelete = installedPkg.Spec.NoopDelete
	desiredApp.Spec.Paused = installedPkg.Spec.Paused
	desiredApp.Spec.Canceled = installedPkg.Spec.Canceled
	desiredApp.Spec.Cluster = installedPkg.Spec.Cluster

	err := controllerutil.SetControllerReference(installedPkg, desiredApp, scheme.Scheme)
	if err != nil {
		return &v1alpha1.App{}, err
	}

	for i, templateStep := range desiredApp.Spec.Template {
		if templateStep.HelmTemplate != nil {
			for _, value := range installedPkg.Spec.Values {
				templateStep.HelmTemplate.ValuesFrom = append(templateStep.HelmTemplate.ValuesFrom, kcv1alpha1.AppTemplateValuesSource{
					SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
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
				templateStep.Ytt.ValuesFrom = append(templateStep.Ytt.ValuesFrom, kcv1alpha1.AppTemplateValuesSource{
					SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
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
