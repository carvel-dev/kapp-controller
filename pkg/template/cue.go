// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"bytes"
	"fmt"
	"io"
	goexec "os/exec"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	"k8s.io/client-go/kubernetes"
)

type cue struct {
	opts        v1alpha1.AppTemplateCue
	coreClient  kubernetes.Interface
	genericOpts GenericOpts
}

var _ Template = &cue{}

func newCue(opts v1alpha1.AppTemplateCue, genericOpts GenericOpts, coreClient kubernetes.Interface) *cue {
	return &cue{opts: opts, genericOpts: genericOpts, coreClient: coreClient}
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
	args := []string{"export", "--out", "yaml"}

	if len(c.opts.Paths) == 0 {
		paths, err := filepath.Glob(filepath.Join(dirPath, "*.cue"))
		if err != nil {
			return exec.NewCmdRunResultWithErr(fmt.Errorf("Reading files: %w", err))
		}
		args = append(args, paths...)
	} else {
		for _, path := range c.opts.Paths {
			checkedPath, err := memdir.ScopedPath(dirPath, path)
			if err != nil {
				return exec.NewCmdRunResultWithErr(fmt.Errorf("Checking path: %w", err))
			}
			args = append(args, checkedPath)
		}
	}

	vals := Values{c.opts.ValuesFrom, c.genericOpts, c.coreClient}
	paths, valuesCleanUpFunc, err := vals.AsPaths(dirPath)
	if err != nil {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Writing values: %w", err))
	}
	defer valuesCleanUpFunc()

	if c.opts.InputExpression != "" {
		args = append(args, "--path", c.opts.InputExpression)
	}
	args = append(args, paths...)
	if c.opts.OutputExpression != "" {
		args = append(args, "--expression", c.opts.OutputExpression)
	}

	var stdoutBs, stderrBs bytes.Buffer
	cmd := goexec.Command("cue", args...)
	cmd.Stdin = input
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs
	cmd.Dir = dirPath

	err = cmd.Run()

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Templating dir: %s", err)

	return result
}
