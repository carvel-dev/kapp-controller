package deploy

import (
	"github.com/k14s/kapp-controller/pkg/exec"
)

type Deploy interface {
	Deploy(tplOutput string, changedFunc func(exec.CmdRunResult)) exec.CmdRunResult
	Delete(changedFunc func(exec.CmdRunResult)) exec.CmdRunResult
	Inspect() exec.CmdRunResult
	ManagedName() string
}

type GenericOpts struct {
	Name      string
	Namespace string
}
