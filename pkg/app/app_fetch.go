// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"

	goexec "os/exec"
)

func (a *App) fetch(dstPath string) (string, exec.CmdRunResult) {
	if len(a.app.Spec.Fetch) == 0 {
		return "", exec.NewCmdRunResultWithErr(fmt.Errorf("Expected at least one fetch option"))
	}

	var result exec.CmdRunResult

	vendir := a.fetchFactory.NewVendir(a.app.Namespace)

	// Because vendir doesn't allow placing contents in the vendir root, we
	// place all contents in sub dirs. For backwards compatibility, we must
	// update dstPath to point to dstPath/0 if there is just one fetch step
	for i, fetch := range a.app.Spec.Fetch {
		err := vendir.AddDir(fetch, strconv.Itoa(i))
		if err != nil {
			result.AttachErrorf(fmt.Sprintf("Fetching (%d): ", i)+"%s", err)
			return "", result
		}
	}

	confReader, err := vendir.ConfigReader()
	if err != nil {
		result.AttachErrorf("Fetching: %v", err)
		return "", result
	}

	result = a.runVendir(confReader, dstPath)

	// if only one fetch, update dstPath for backwards compatibility
	if len(a.app.Spec.Fetch) == 1 {
		dstPath = path.Join(dstPath, "0")
	}

	return dstPath, result
}

func (a *App) runVendir(confReader io.Reader, workingDir string) exec.CmdRunResult {
	var stdoutBs, stderrBs bytes.Buffer
	cmd := goexec.Command("vendir", "sync", "-f", "-", "--lock-file", os.DevNull)
	cmd.Dir = workingDir
	cmd.Stdin = confReader
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := cmd.Run()

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Fetching resources: %s", err)

	return result
}
