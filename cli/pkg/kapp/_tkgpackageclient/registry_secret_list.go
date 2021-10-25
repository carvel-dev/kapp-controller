// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	secretgenctrl "github.com/vmware-tanzu/carvel-secretgen-controller/pkg/apis/secretgen2/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

// ListRegistrySecrets lists all registry secrets of type kubernetes.io/dockerconfigjson across the cluster.
func (p *pkgClient) ListRegistrySecrets(o *tkgpackagedatamodel.RegistrySecretOptions) (*corev1.SecretList, error) {
	registrySecretList, err := p.kappClient.ListRegistrySecrets(o.Namespace)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list existing registry secrets in the cluster")
	}

	return registrySecretList, nil
}

// ListSecretExports lists all SecretExports across the cluster.
func (p *pkgClient) ListSecretExports(o *tkgpackagedatamodel.RegistrySecretOptions) (*secretgenctrl.SecretExportList, error) {
	secretExportList, err := p.kappClient.ListSecretExports(o.Namespace)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list existing secret exports in the cluster")
	}

	return secretExportList, nil
}
