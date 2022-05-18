// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"bytes"
	"fmt"
	"io"
	goexec "os/exec"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	"k8s.io/client-go/kubernetes"
)

type Ytt struct {
	opts         v1alpha1.AppTemplateYtt
	genericOpts  GenericOpts
	coreClient   kubernetes.Interface
	fetchFactory fetch.Factory
	cmdRunner    exec.CmdRunner
}

var _ Template = &Ytt{}

func NewYtt(opts v1alpha1.AppTemplateYtt, genericOpts GenericOpts,
	coreClient kubernetes.Interface, fetchFactory fetch.Factory, cmdRunner exec.CmdRunner) *Ytt {

	return &Ytt{opts, genericOpts, coreClient, fetchFactory, cmdRunner}
}

func (t *Ytt) TemplateDir(dirPath string) (exec.CmdRunResult, bool) {
	return t.template(dirPath, nil), true
}

func (t *Ytt) TemplateStream(input io.Reader, dirPath string) exec.CmdRunResult {
	return t.template(dirPath, input)
}

func (t *Ytt) template(dirPath string, input io.Reader) exec.CmdRunResult {
	args := t.addArgs([]string{})

	args, err := t.addPaths(dirPath, input, args)
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

	args = t.addFileMarks(args)

	{ // Add values files
		vals := Values{t.opts.ValuesFrom, t.genericOpts, t.coreClient}

		paths, valuesCleanUpFunc, err := vals.AsPaths(dirPath)
		if err != nil {
			return exec.NewCmdRunResultWithErr(err)
		}

		defer valuesCleanUpFunc()

		for _, path := range paths {
			args = append(args, []string{"--data-values-file", path}...)
		}
	}

	var stdoutBs, stderrBs bytes.Buffer

	cmd := goexec.Command("ytt", args...)
	cmd.Stdin = input
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err = t.cmdRunner.Run(cmd)

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

func (t *Ytt) addPaths(dirPath string, input io.Reader, args []string) ([]string, error) {
	// TODO we currently do not allow file path remapping -f new-path=actual-path syntax
	// If explicit paths provided, expect user specify stdin explicitly
	switch {
	case len(t.opts.Paths) > 0:
		for _, path := range t.opts.Paths {
			if path == stdinPath {
				if input == nil {
					return nil, fmt.Errorf("Expected stdin to be available when using it as path, but was not")
				}
				args = append(args, "-f", path)
			} else {
				checkedPath, err := memdir.ScopedPath(dirPath, path)
				if err != nil {
					return nil, fmt.Errorf("Checking path: %s", err)
				}
				args = append(args, "-f", checkedPath)
			}
		}
		return args, nil

	case input != nil:
		return append(args, "-f", "-"), nil

	default:
		return append(args, "-f", dirPath), nil
	}
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

func (t *Ytt) addFileMarks(args []string) []string {
	for _, fileMark := range t.opts.FileMarks {
		args = append(args, []string{"--file-mark", fileMark}...)
	}

	return args
}
