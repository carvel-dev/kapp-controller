// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

var (
	serverListenType = "unix"
	serverListenAddr = os.Getenv("KAPPCTRL_SIDECAREXEC_SOCK")
)

// Server accepts RPCs to execute commands or configure runtime environment.
type Server struct {
	cmdExec *CmdExec
	log     logr.Logger
}

// ServerOpts accepts Server's configuration.
type ServerOpts struct {
	AllowedCmdNames []string
}

// NewServer returns a new Server.
func NewServer(local exec.CmdRunner, opts ServerOpts, log logr.Logger) *Server {
	allowedCmdNames := map[string]struct{}{}
	for _, cmd := range opts.AllowedCmdNames {
		allowedCmdNames[cmd] = struct{}{}
	}
	return &Server{&CmdExec{local, allowedCmdNames}, log}
}

// Serve starts an RPC server.
func (r *Server) Serve() error {
	// See which methods satisfy criteria: https://pkg.go.dev/net/rpc#pkg-overview
	// e.g.   func (t *T) MethodName(argType T1, replyType *T2) error

	err := rpc.Register(r.cmdExec)
	if err != nil {
		return fmt.Errorf("Registering CmdExec RPC methods: %s", err)
	}

	err = rpc.Register(NewOSConfig(r.log))
	if err != nil {
		return fmt.Errorf("Registering OSConfig RPC methods: %s", err)
	}

	rpc.HandleHTTP()

	// Socket file may not be cleaned up upon unexpected process termination
	// and without removal will result in "already in use" error messages.
	if serverListenType == "unix" {
		err = os.RemoveAll(serverListenAddr)
		if err != nil {
			return fmt.Errorf("Removing (unbinding) all listen addr: %s", err)
		}
	}

	listener, err := net.Listen(serverListenType, serverListenAddr)
	if err != nil {
		return fmt.Errorf("Listening RPC: %s", err)
	}
	return http.Serve(listener, nil)
}
