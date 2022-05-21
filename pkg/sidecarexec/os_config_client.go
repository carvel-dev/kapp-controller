// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"fmt"
	"net/rpc"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
)

type OSConfigClient struct {
	rpcClient *rpc.Client
}

func (r OSConfigClient) ApplyCACerts(chain string) error {
	err := r.rpcClient.Call("OSConfig.ApplyCACerts", chain, nil)
	if err != nil {
		return fmt.Errorf("Internal run comm: %s", err)
	}
	return nil
}

func (r OSConfigClient) ApplyProxy(in config.ProxyOpts) error {
	err := r.rpcClient.Call("OSConfig.ApplyProxy", in, nil)
	if err != nil {
		return fmt.Errorf("Internal run comm: %s", err)
	}
	return nil
}
