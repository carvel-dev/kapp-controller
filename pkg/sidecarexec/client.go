// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

// Package sidecarexec provides an implementation of a sidecar container in kapp-controller which
// runs each bundled binary in this separate container.
// This was introduced for security purposes, to reduce the attack vector on kapp-controller container
// by moving the binary exec calls to it's own isolated container.
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
