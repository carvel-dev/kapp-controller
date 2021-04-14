// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var (
	trailingSpace = regexp.MustCompile("\\s+\n")
)

type CmdRunResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
	Finished bool
}

func NewCmdRunResultWithErr(err error) CmdRunResult {
	var result CmdRunResult
	result.AttachErrorf("%s", err)
	return result
}

func (r *CmdRunResult) AttachErrorf(msg string, err error) {
	r.Finished = true
	if err != nil {
		r.ExitCode = -1
		if exitError, ok := err.(*exec.ExitError); ok {
			r.ExitCode = exitError.ExitCode()
		}

		if err.Error() == "exit status 1" {
			r.Error = fmt.Errorf(msg, "Error (see .status.usefulErrorMessage for details)")
		} else {
			r.Error = fmt.Errorf(msg, err)
		}
	}
}

func (r CmdRunResult) ErrorStr() string {
	if r.Error != nil {
		return r.Error.Error()
	}
	return ""
}

func (r CmdRunResult) WithFriendlyYAMLStrings() CmdRunResult {
	// YAML can format muliline strings nicely
	// if they do not have trailing spaces right before newlines
	return CmdRunResult{
		Stdout:   trailingSpace.ReplaceAllString(strings.TrimSpace(r.Stdout), "\n"),
		Stderr:   trailingSpace.ReplaceAllString(strings.TrimSpace(r.Stderr), "\n"),
		ExitCode: r.ExitCode,
		Error:    r.Error,
		Finished: r.Finished,
	}
}
