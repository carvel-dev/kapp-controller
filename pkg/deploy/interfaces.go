// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"carvel.dev/kapp-controller/pkg/exec"
)

type Deploy interface {
	Deploy(tplOutput string, startedApplyingFunc func(),
		changedFunc func(exec.CmdRunResult)) exec.CmdRunResult

	Delete(startedApplyingFunc func(),
		changedFunc func(exec.CmdRunResult)) exec.CmdRunResult

	Inspect() exec.CmdRunResult
}
