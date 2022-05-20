// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	goexec "os/exec"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

var (
	serverListenType = "unix"
	serverListenAddr = os.Getenv("KAPPCTRL_SIDECAREXEC_SOCK")
)

type CmdInput struct {
	Command string
	Args    []string
	Stdin   []byte
	Env     []string
	Dir     string
}

type CmdOutput struct {
	Stdout []byte
	Stderr []byte
	Error  string
}

type Server struct {
	serverMethods *ServerMethods
}

type ServerOpts struct {
	AllowedCmdNames []string
}

// NewServer returns a new Server.
func NewServer(local exec.CmdRunner, opts ServerOpts) *Server {
	allowedCmdNames := map[string]struct{}{}
	for _, cmd := range opts.AllowedCmdNames {
		allowedCmdNames[cmd] = struct{}{}
	}
	return &Server{&ServerMethods{local, allowedCmdNames}}
}

func (r *Server) Serve() error {
	err := rpc.Register(r.serverMethods)
	if err != nil {
		return fmt.Errorf("Registering RPC methods: %s", err)
	}

	rpc.HandleHTTP()

	listener, err := net.Listen(serverListenType, serverListenAddr)
	if err != nil {
		return err
	}
	go http.Serve(listener, nil)
	select {} // TODO
}

type ServerMethods struct {
	local           exec.CmdRunner
	allowedCmdNames map[string]struct{}
}

func (r ServerMethods) Run(input CmdInput, output *CmdOutput) error {
	if _, found := r.allowedCmdNames[input.Command]; !found {
		return fmt.Errorf("Command '%s' is not allowed", input.Command)
	}

	cmd := goexec.Command(input.Command, input.Args...)

	if len(input.Stdin) > 0 {
		cmd.Stdin = bytes.NewBuffer(input.Stdin)
	}
	if len(input.Env) > 0 {
		cmd.Env = input.Env
	}
	if len(input.Dir) > 0 {
		cmd.Dir = input.Dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := r.local.Run(cmd)
	if err != nil {
		output.Error = err.Error()
	}

	output.Stdout = stdout.Bytes()
	output.Stderr = stderr.Bytes()
	return nil
}
