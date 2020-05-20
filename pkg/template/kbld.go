package template

import (
	"bytes"
	"io"
	goexec "os/exec"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/exec"
)

type Kbld struct {
	opts        v1alpha1.AppTemplateKbld
	genericOpts GenericOpts
}

var _ Template = &Kbld{}

func NewKbld(opts v1alpha1.AppTemplateKbld, genericOpts GenericOpts) *Kbld {
	return &Kbld{opts, genericOpts}
}

func (t *Kbld) TemplateDir(dirPath string) exec.CmdRunResult {
	return t.template(dirPath, nil)
}

func (t *Kbld) TemplateStream(input io.Reader) exec.CmdRunResult {
	return t.template(stdinPath, input)
}

func (t *Kbld) template(dirPath string, input io.Reader) exec.CmdRunResult {
	args := t.addArgs([]string{"-f", dirPath})

	var stdoutBs, stderrBs bytes.Buffer

	cmd := goexec.Command("kbld", args...)
	cmd.Stdin = input
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := cmd.Run()

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Templating dir: %s", err)

	return result
}

func (t *Kbld) addArgs(args []string) []string {
	return args
}
