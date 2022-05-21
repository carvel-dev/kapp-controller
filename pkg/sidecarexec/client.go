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
func NewClient(local exec.CmdRunner) Client {
	rpcClient, err := rpc.DialHTTP(serverListenType, serverListenAddr)
	if err != nil {
		// TODO
		panic("Dialing rpc")
	}

	return Client{local, rpcClient}
}

func (r Client) CmdExec() CmdExecClient {
	return CmdExecClient{r.local, r.rpcClient}
}

func (r Client) OSConfig() OSConfigClient {
	return OSConfigClient{r.rpcClient}
}
