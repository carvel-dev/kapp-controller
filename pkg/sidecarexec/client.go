// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"net/rpc"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

// Client provides access to sidecarexec API.
type Client struct {
	local     exec.CmdRunner
	rpcClient *rpc.Client
}

// NewClient returns a new Client.
func NewClient(local exec.CmdRunner) (Client, error) {
	rpcClient, err := rpc.DialHTTP(serverListenType, serverListenAddr)
	if err != nil {
		return Client{}, err
	}
	return Client{local, rpcClient}, nil
}

// CmdExec returns command execution implementation.
func (r Client) CmdExec() CmdExecClient {
	return CmdExecClient{r.local, r.rpcClient}
}

// OSConfig returns runtime environment configuration implementation.
func (r Client) OSConfig() OSConfigClient {
	return OSConfigClient{r.rpcClient}
}
