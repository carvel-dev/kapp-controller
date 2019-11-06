package template

import (
	"github.com/k14s/kapp-controller/pkg/exec"
	"io"
)

type Template interface {
	TemplateDir(dirPath string) exec.CmdRunResult
	TemplateStream(io.Reader) exec.CmdRunResult
}

type GenericOpts struct {
	Name      string
	Namespace string
}
