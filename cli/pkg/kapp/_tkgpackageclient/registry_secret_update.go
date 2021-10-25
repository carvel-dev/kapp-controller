// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"

	secretgenctrl "github.com/vmware-tanzu/carvel-secretgen-controller/pkg/apis/secretgen2/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

var (
	secret       = &corev1.Secret{}
	secretExport = &secretgenctrl.SecretExport{}
)

// UpdateRegistrySecret updates a registry Secret in the cluster
func (p *pkgClient) UpdateRegistrySecret(o *tkgpackagedatamodel.RegistrySecretOptions) error {
	var (
		registry string
		username string
		password string
		dataMap  map[string]interface{}
	)

	if err := p.kappClient.GetClient().Get(context.Background(), crtclient.ObjectKey{Name: o.SecretName, Namespace: o.Namespace}, secret); err != nil {
		if apierrors.IsNotFound(err) {
			return errors.New(fmt.Sprintf("secret '%s' does not exist in namespace '%s'", o.SecretName, o.Namespace))
		}
		return err
	}

	secretToUpdate := secret.DeepCopy()
	if err := json.Unmarshal(secretToUpdate.Data[corev1.DockerConfigJsonKey], &dataMap); err != nil {
		return err
	}

	auths, ok := dataMap["auths"]
	if !ok {
		return errors.New(fmt.Sprintf("no 'auths' entry exists in secret '%s'", o.SecretName))
	}

	entries := auths.(map[string]interface{})
	if len(entries) != 1 {
		return errors.New(fmt.Sprintf("updating secret '%s' is not allowed as multiple registry entries exists", o.SecretName))
	}

	for reg, v := range entries {
		registry = reg
		credentials := v.(map[string]interface{})
		currentUsername, ok := credentials["username"]
		if !ok {
			return errors.New(fmt.Sprintf("no 'username' entry exists in secret '%s'", o.SecretName))
		}
		username = currentUsername.(string)
		currentPassword, ok := credentials["password"]
		if !ok {
			return errors.New(fmt.Sprintf("no 'password' entry exists in secret '%s'", o.SecretName))
		}
		password = currentPassword.(string)
	}

	if o.Username != "" {
		username = o.Username
	}

	if o.Password != "" {
		password = o.Password
	}

	dockerCfg := DockerConfigJSON{Auths: map[string]dockerConfigEntry{registry: {Username: username, Password: password}}}
	dockerCfgContent, err := json.Marshal(dockerCfg)
	if err != nil {
		return err
	}
	secretToUpdate.Data[corev1.DockerConfigJsonKey] = dockerCfgContent

	if err := p.kappClient.GetClient().Update(context.Background(), secretToUpdate); err != nil {
		return errors.Wrap(err, "failed to update Secret resource")
	}

	if err := p.UpdateSecretExport(o); err != nil {
		return err
	}

	return nil
}

// UpdateSecretExport updates the SecretExport resource in the cluster
func (p *pkgClient) UpdateSecretExport(o *tkgpackagedatamodel.RegistrySecretOptions) error {
	if o.Export.ExportToAllNamespaces == nil {
		return nil
	}

	if *o.Export.ExportToAllNamespaces {
		err := p.kappClient.GetClient().Get(context.Background(), crtclient.ObjectKey{Name: o.SecretName, Namespace: o.Namespace}, secretExport)
		if err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
			secretExport = p.newSecretExport(o.SecretName, o.Namespace)
			if err := p.kappClient.GetClient().Create(context.Background(), secretExport); err != nil {
				return errors.Wrap(err, "failed to create SecretExport resource")
			}
			return nil
		}
		secretExportToUpdate := secretExport.DeepCopy()
		secretExportToUpdate.Spec = secretgenctrl.SecretExportSpec{ToNamespaces: []string{"*"}}
		if err := p.kappClient.GetClient().Update(context.Background(), secretExportToUpdate); err != nil {
			return errors.Wrap(err, "failed to update SecretExport resource")
		}
	} else { // un-export already exported secrets
		secretExport = &secretgenctrl.SecretExport{
			TypeMeta:   metav1.TypeMeta{Kind: tkgpackagedatamodel.KindSecretExport, APIVersion: secretgenctrl.SchemeGroupVersion.String()},
			ObjectMeta: metav1.ObjectMeta{Name: o.SecretName, Namespace: o.Namespace},
		}
		if err := p.kappClient.GetClient().Delete(context.Background(), secretExport); err != nil {
			if !apierrors.IsNotFound(err) {
				return errors.Wrap(err, "failed to delete SecretExport resource")
			}
		}
	}

	return nil
}
