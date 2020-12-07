// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	kyaml "sigs.k8s.io/yaml"
)

type Inline struct {
	opts       v1alpha1.AppFetchInline
	nsName     string
	coreClient kubernetes.Interface
}

func NewInline(opts v1alpha1.AppFetchInline, nsName string, coreClient kubernetes.Interface) *Inline {
	return &Inline{opts, nsName, coreClient}
}

func (t *Inline) VendirRes(dirPath string) (vendirconf.Directory, [][]byte, error) {
	dir := NewVendir().InlineDirConf(t.opts, dirPath)

	resources, err := t.resources()
	if err != nil {
		return vendirconf.Directory{}, nil, fmt.Errorf("Fecthing resources: %v", err)
	}

	return dir, resources, nil

}

func (t *Inline) resources() ([][]byte, error) {
	var resources [][]byte

	for _, source := range t.opts.PathsFrom {
		switch {
		case source.SecretRef != nil:
			bytes, err := t.secretBytes(*source.SecretRef)
			if err != nil {
				return nil, err
			}

			resources = append(resources, bytes)

		case source.ConfigMapRef != nil:
			bytes, err := t.configMapBytes(*source.SecretRef)
			if err != nil {
				return nil, err
			}

			resources = append(resources, bytes)
		}
	}

	return resources, nil
}

func (t *Inline) secretBytes(secretRef v1alpha1.AppFetchInlineSourceRef) ([]byte, error) {
	secret, err := t.coreClient.CoreV1().Secrets(t.nsName).Get(secretRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	secret.TypeMeta.Kind = "Secret"
	secret.TypeMeta.APIVersion = "v1"

	return kyaml.Marshal(secret)
}

func (t *Inline) configMapBytes(configMapRef v1alpha1.AppFetchInlineSourceRef) ([]byte, error) {
	configMap, err := t.coreClient.CoreV1().ConfigMaps(t.nsName).Get(configMapRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	configMap.TypeMeta.Kind = "ConfigMap"
	configMap.TypeMeta.APIVersion = "v1"

	return kyaml.Marshal(configMap)
}

func (t *Inline) Retrieve(dstPath string) error {
	for path, content := range t.opts.Paths {
		err := t.writeFile(dstPath, path, []byte(content))
		if err != nil {
			return err
		}
	}

	for _, source := range t.opts.PathsFrom {
		switch {
		case source.SecretRef != nil:
			err := t.writeFromSecret(dstPath, *source.SecretRef)
			if err != nil {
				return err
			}

		case source.ConfigMapRef != nil:
			err := t.writeFromConfigMap(dstPath, *source.ConfigMapRef)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("Expected either secretRef or configMapRef as a source")
		}
	}

	return nil
}

func (t *Inline) writeFromSecret(dstPath string, secretRef v1alpha1.AppFetchInlineSourceRef) error {
	secret, err := t.coreClient.CoreV1().Secrets(t.nsName).Get(secretRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for name, val := range secret.Data {
		err := t.writeFile(dstPath, filepath.Join(secretRef.DirectoryPath, name), val)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Inline) writeFromConfigMap(dstPath string, configMapRef v1alpha1.AppFetchInlineSourceRef) error {
	configMap, err := t.coreClient.CoreV1().ConfigMaps(t.nsName).Get(configMapRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for name, val := range configMap.Data {
		err := t.writeFile(dstPath, filepath.Join(configMapRef.DirectoryPath, name), []byte(val))
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Inline) writeFile(dstPath, subPath string, content []byte) error {
	newPath, err := memdir.ScopedPath(dstPath, subPath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(newPath, content, 0600)
	if err != nil {
		return fmt.Errorf("Writing file '%s': %s", newPath, err)
	}

	return nil
}
