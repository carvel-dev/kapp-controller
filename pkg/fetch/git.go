// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"fmt"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kyaml "sigs.k8s.io/yaml"
)

type Git struct {
	opts       v1alpha1.AppFetchGit
	nsName     string
	coreClient kubernetes.Interface
}

func NewGit(opts v1alpha1.AppFetchGit, nsName string, coreClient kubernetes.Interface) *Git {
	return &Git{opts, nsName, coreClient}
}

func (t *Git) VendirRes(dirPath string) (vendirconf.Directory, [][]byte, error) {
	dir := NewVendir().GitDirConf(t.opts, dirPath)

	resources, err := t.resources()
	if err != nil {
		return vendirconf.Directory{}, nil, fmt.Errorf("Fecthing resources: %v", err)
	}

	return dir, resources, nil
}

func (t *Git) resources() ([][]byte, error) {
	if t.opts.SecretRef == nil {
		return nil, nil
	}

	secret, err := t.coreClient.CoreV1().Secrets(t.nsName).Get(t.opts.SecretRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// typed clients drop GVK or resource (https://github.com/kubernetes/kubernetes/issues/80609)
	secret.TypeMeta.Kind = "Secret"
	secret.TypeMeta.APIVersion = "v1"

	sBytes, err := kyaml.Marshal(secret)
	if err != nil {
		return nil, err
	}

	return [][]byte{sBytes}, nil
}
