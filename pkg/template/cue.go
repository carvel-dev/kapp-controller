// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"bytes"
	"fmt"
	"io"
	goexec "os/exec"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

type cue struct {
	opts        v1alpha1.AppTemplateCue
	genericOpts GenericOpts
}

var _ Template = &cue{}

func newCue(opts v1alpha1.AppTemplateCue, genericOpts GenericOpts) *cue {
	return &cue{opts: opts, genericOpts: genericOpts}
}

// TemplateDir works on directory returning templating result,
// and boolean indicating whether subsequent operations
// should operate on result, or continue operating on the directory
func (c *cue) TemplateDir(dirPath string) (exec.CmdRunResult, bool) {
	return c.template(dirPath, nil), true
}

// TemplateStream works on a stream returning templating result.
// dirPath is provided for context from which to reference additonal inputs.
func (c *cue) TemplateStream(stream io.Reader, dirPath string) exec.CmdRunResult {
	return exec.NewCmdRunResultWithErr(fmt.Errorf("Templating stream is not supported")) // TODO: Implement
}

func (c *cue) template(dirPath string, input io.Reader) exec.CmdRunResult {
	var stdoutBs, stderrBs bytes.Buffer
	args := []string{"export", "--out", "yaml"}
	args = append(args, c.opts.Paths...)

	cmd := goexec.Command("cue", args...)
	cmd.Stdin = input
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs
	cmd.Dir = dirPath

	err := cmd.Run()

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Templating dir: %s", err)

	return result
}
