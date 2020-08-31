// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"bytes"
	"fmt"
	"io"
	goexec "os/exec"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/exec"
	"github.com/k14s/kapp-controller/pkg/fetch"
	"github.com/k14s/kapp-controller/pkg/memdir"
)

type Ytt struct {
	opts         v1alpha1.AppTemplateYtt
	genericOpts  GenericOpts
	fetchFactory fetch.Factory
}

var _ Template = &Ytt{}

func NewYtt(opts v1alpha1.AppTemplateYtt,
	genericOpts GenericOpts, fetchFactory fetch.Factory) *Ytt {
	return &Ytt{opts, genericOpts, fetchFactory}
}

func (t *Ytt) TemplateDir(dirPath string) exec.CmdRunResult {
	return t.template(dirPath, nil)
}

func (t *Ytt) TemplateStream(input io.Reader) exec.CmdRunResult {
	return t.template(stdinPath, input)
}

func (t *Ytt) template(dirPath string, input io.Reader) exec.CmdRunResult {
	args := t.addArgs([]string{})

	args, err := t.addPaths(dirPath, args)
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	args, inlineDir, err := t.addInlinePaths(args)
	if inlineDir != nil {
		defer inlineDir.Remove()
	}
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	var stdoutBs, stderrBs bytes.Buffer

	cmd := goexec.Command("ytt", args...)
	cmd.Stdin = input
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err = cmd.Run()

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Templating dir: %s", err)

	return result
}

func (t *Ytt) addArgs(args []string) []string {
	if t.opts.IgnoreUnknownComments {
		args = append(args, "--ignore-unknown-comments")
	}
	return args
}

func (t *Ytt) addPaths(dirPath string, args []string) ([]string, error) {
	if len(t.opts.Paths) > 0 {
		// Disallow path selection when consuming a stream
		if dirPath == stdinPath {
			return nil, fmt.Errorf("Paths cannot be used with template stream")
		}

		for _, path := range t.opts.Paths {
			checkedPath, err := memdir.ScopedPath(dirPath, path)
			if err != nil {
				return nil, fmt.Errorf("Checking path: %s", err)
			}
			args = append(args, []string{"-f", checkedPath}...)
		}

		return args, nil
	}

	return append(args, []string{"-f", dirPath}...), nil
}

func (t *Ytt) addInlinePaths(args []string) ([]string, *memdir.TmpDir, error) {
	if t.opts.Inline == nil {
		return args, nil, nil
	}

	inlineDir := memdir.NewTmpDir("template-ytt-inline")

	err := inlineDir.Create()
	if err != nil {
		return nil, nil, err
	}

	inline := t.fetchFactory.NewInline(*t.opts.Inline, t.genericOpts.Namespace)

	err = inline.Retrieve(inlineDir.Path())
	if err != nil {
		return nil, inlineDir, err
	}

	args = append(args, []string{"-f", inlineDir.Path()}...)

	return args, inlineDir, nil
}
