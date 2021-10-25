// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/kappclient"
)

type pkgClient struct {
	kappClient kappclient.Client
}

// NewTKGPackageClient instantiates pkgClient
func NewTKGPackageClient(kubeconfigPath string) (TKGPackageClient, error) {
	var err error
	client := &pkgClient{}

	if client.kappClient, err = kappclient.NewKappClient(kubeconfigPath); err != nil {
		return nil, err
	}

	return client, nil
}
