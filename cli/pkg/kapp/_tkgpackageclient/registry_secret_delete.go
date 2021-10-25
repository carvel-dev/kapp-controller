// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	secretgenctrl "github.com/vmware-tanzu/carvel-secretgen-controller/pkg/apis/secretgen2/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

// DeleteRegistrySecret deletes a registry secret from the cluster
func (p *pkgClient) DeleteRegistrySecret(o *tkgpackagedatamodel.RegistrySecretOptions) (bool, error) {
	secretExport = &secretgenctrl.SecretExport{
		TypeMeta:   metav1.TypeMeta{Kind: tkgpackagedatamodel.KindSecretExport, APIVersion: secretgenctrl.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{Name: o.SecretName, Namespace: o.Namespace},
	}
	if err := p.kappClient.GetClient().Delete(context.Background(), secretExport); err != nil {
		if !apierrors.IsNotFound(err) {
			return true, errors.Wrap(err, "failed to delete SecretExport resource")
		}
	}

	secret = &corev1.Secret{
		TypeMeta:   metav1.TypeMeta{Kind: tkgpackagedatamodel.KindSecret, APIVersion: corev1.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{Name: o.SecretName, Namespace: o.Namespace},
	}
	if err := p.kappClient.GetClient().Delete(context.Background(), secret); err != nil {
		if !apierrors.IsNotFound(err) {
			return true, errors.Wrap(err, "failed to delete Secret resource")
		}
		return false, nil
	}

	return true, nil
}
