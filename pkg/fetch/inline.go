// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Inline struct helps you fetch files or strings from a storage location such as a configmap or secret.
// For an example if you had a controller in your cluster that generated certs and stored them in a configmap or secret,
// you could use this inline fetcher to retrieve those certs to pass them on to a templating step.
type Inline struct {
	opts       v1alpha1.AppFetchInline
	nsName     string
	coreClient kubernetes.Interface
}

func NewInline(opts v1alpha1.AppFetchInline, nsName string, coreClient kubernetes.Interface) *Inline {
	return &Inline{opts, nsName, coreClient}
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
	secret, err := t.coreClient.CoreV1().Secrets(t.nsName).Get(context.Background(), secretRef.Name, metav1.GetOptions{})
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
	configMap, err := t.coreClient.CoreV1().ConfigMaps(t.nsName).Get(context.Background(), configMapRef.Name, metav1.GetOptions{})
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
