// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package local

import (
	"fmt"
	"io"
	"net"
	goexec "os/exec"
	"strings"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

type DetailedCmdRunner struct {
	log        io.Writer
	fullOutput bool
}

var _ exec.CmdRunner = &DetailedCmdRunner{}

func NewDetailedCmdRunner(log io.Writer, fullOutput bool) *DetailedCmdRunner {
	return &DetailedCmdRunner{log, fullOutput}
}

func (r DetailedCmdRunner) Run(cmd *goexec.Cmd) error {
	if r.fullOutput {
		cmd.Stdout = io.MultiWriter(r.log, cmd.Stdout)
		cmd.Stderr = io.MultiWriter(r.log, cmd.Stderr)
	}

	fmt.Fprintf(r.log, "==> Executing %s %v\n", cmd.Path, cmd.Args)
	defer fmt.Fprintf(r.log, "==> Finished executing %s\n\n", cmd.Path)

	return exec.PlainCmdRunner{}.Run(cmd)
}

func (r DetailedCmdRunner) RunWithCancel(cmd *goexec.Cmd, cancelCh chan struct{}) error {
	if r.fullOutput {
		cmd.Stdout = io.MultiWriter(r.log, cmd.Stdout)
		cmd.Stderr = io.MultiWriter(r.log, cmd.Stderr)
	}

	fmt.Fprintf(r.log, "==> Executing %s %v\n", cmd.Path, cmd.Args)
	defer fmt.Fprintf(r.log, "==> Finished executing %s\n\n", cmd.Path)

	addServerDetailsToKubeConfig(cmd)
	return exec.PlainCmdRunner{}.RunWithCancel(cmd, cancelCh)
}

func addServerDetailsToKubeConfig(cmd *goexec.Cmd) {
	envMap := map[string]string{}
	for _, env := range cmd.Env {
		envKeyVal := strings.SplitN(env, "=", 2)
		if len(envKeyVal) > 1 {
			envMap[envKeyVal[0]] = envKeyVal[1]
		} else {
			envMap[envKeyVal[0]] = ""
		}
	}

	var envList []string
	for key, val := range envMap {
		if key == "KAPP_KUBECONFIG_YAML" {
			var envHostPort string
			if len(envMap["KUBERNETES_SERVICE_PORT"]) != 0 {
				envHostPort = net.JoinHostPort(envMap["KUBERNETES_SERVICE_HOST"], envMap["KUBERNETES_SERVICE_PORT"])
			} else {
				envHostPort = envMap["KUBERNETES_SERVICE_HOST"]
			}
			// Is it with https? what are the implications
			val = strings.ReplaceAll(val, "${KAPP_KUBERNETES_SERVICE_HOST_PORT}", envHostPort)
		}
		envList = append(envList, strings.Join([]string{key, val}, "="))

	}
	cmd.Env = envList

}
