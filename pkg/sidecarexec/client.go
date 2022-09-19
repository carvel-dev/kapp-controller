// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"net/rpc"
	"sync"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

// Client provides access to sidecarexec API.
type Client struct {
	local     exec.CmdRunner
	rpcClient rpcClient
}

// NewClient returns a new Client.
func NewClient(local exec.CmdRunner) (Client, error) {
	return Client{local, &reconnectingRPCClient{}}, nil
}

// CmdExec returns command execution implementation.
func (r Client) CmdExec() CmdExecClient {
	return CmdExecClient{r.local, r.rpcClient}
}

// OSConfig returns runtime environment configuration implementation.
func (r Client) OSConfig() OSConfigClient {
	return OSConfigClient{r.rpcClient}
}

type rpcClient interface {
	Call(serviceMethod string, args any, reply any) error
}

type reconnectingRPCClient struct {
	clientLock sync.Mutex
	client     *rpc.Client
}

func (c *reconnectingRPCClient) Call(serviceMethod string, args any, reply any) error {
	client, err := c.connect(nil)
	if err != nil {
		return err
	}

	err = client.Call(serviceMethod, args, reply)
	if err == rpc.ErrShutdown {
		refreshedClient, err := c.connect(client)
		if err != nil {
			return err
		}

		err = refreshedClient.Call(serviceMethod, args, reply)
	}
	return err
}

// connect a client (which is nil or disconnected) in a thread-safe manner.
// This is intended for clients which encountered a connection error,
// so that we only try to reconnect if we haven't already done so in a different thread.
func (c *reconnectingRPCClient) connect(disconnectedClient *rpc.Client) (*rpc.Client, error) {
	c.clientLock.Lock()
	defer c.clientLock.Unlock()

	if c.client != nil {
		if disconnectedClient == nil || disconnectedClient != c.client {
			return c.client, nil
		}
		_ = c.client.Close()
		c.client = nil
	}

	client, err := rpc.DialHTTP(serverListenType, serverListenAddr)
	if err != nil {
		return nil, err
	}
	c.client = client
	return c.client, nil
}
