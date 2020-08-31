// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"github.com/k14s/kapp-controller/pkg/exec"
)

type Deploy interface {
	Deploy(tplOutput string, startedApplyingFunc func(),
		changedFunc func(exec.CmdRunResult)) exec.CmdRunResult

	Delete(startedApplyingFunc func(),
		changedFunc func(exec.CmdRunResult)) exec.CmdRunResult

	Inspect() exec.CmdRunResult
}

type GenericOpts struct {
	Name           string
	Namespace      string
	KubeconfigYAML string
}
