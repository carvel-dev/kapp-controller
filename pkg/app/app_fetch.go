// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strconv"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"

	goexec "os/exec"

	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

func (a *App) fetch(dstPath string) (string, exec.CmdRunResult) {
	if len(a.app.Spec.Fetch) == 0 {
		return "", exec.NewCmdRunResultWithErr(fmt.Errorf("Expected at least one fetch option"))
	}

	var (
		result      exec.CmdRunResult
		resourceYml [][]byte
	)

	vendirConf := vendirconf.Config{
		APIVersion: "vendir.k14s.io/v1alpha1", // TODO: use constant from vendir package
		Kind:       "Config",                  // TODO: use constant from vendir package
	}

	// Because vendir doesn't allow placing contents in the vendir root, we
	// place all contents in sub dirs. For backwards compatibility, we must
	// update dstPath to point to dstPath/0 if there is just one fetch step
	for i, fetch := range a.app.Spec.Fetch {
		dir, resources, err := a.vendirResFor(fetch, strconv.Itoa(i))
		if err != nil {
			result.AttachErrorf(fmt.Sprintf("Fetching (%d): ", i)+"%s", err)
			return "", result
		}

		vendirConf.Directories = append(vendirConf.Directories, dir)
		resourceYml = append(resourceYml, resources...)
	}

	confReader, err := a.confReader(vendirConf, resourceYml)
	if err != nil {
		result.AttachErrorf("%v", err)
		return "", result
	}

	result = a.runVendir(confReader, dstPath)

	// if only one fetch, update dstPath for backwards compatibility
	if len(a.app.Spec.Fetch) == 1 {
		dstPath = path.Join(dstPath, "0")
	}

	return dstPath, result
}

func (a *App) vendirResFor(fetch v1alpha1.AppFetch, destPath string) (vendirconf.Directory, [][]byte, error) {
	switch {
	case fetch.Inline != nil:
		return a.fetchFactory.NewInline(*fetch.Inline, a.app.Namespace).VendirRes(destPath)
	case fetch.Image != nil:
		return a.fetchFactory.NewImage(*fetch.Image, a.app.Namespace).VendirRes(destPath)
	case fetch.HTTP != nil:
		return a.fetchFactory.NewHTTP(*fetch.HTTP, a.app.Namespace).VendirRes(destPath)
	case fetch.Git != nil:
		return a.fetchFactory.NewGit(*fetch.Git, a.app.Namespace).VendirRes(destPath)
	case fetch.HelmChart != nil:
		return a.fetchFactory.NewHelmChart(*fetch.HelmChart, a.app.Namespace).VendirRes(destPath)
	}

	return vendirconf.Directory{}, nil, fmt.Errorf("Unsupported way to fetch templates")
}

func (a *App) confReader(conf vendirconf.Config, resourceYml [][]byte) (io.Reader, error) {
	vendirConfBytes, err := conf.AsBytes()
	if err != nil {
		return nil, err
	}

	finalConfig := bytes.Join(append(resourceYml, vendirConfBytes), []byte("---\n"))

	return bytes.NewReader(finalConfig), nil
}

func (a *App) runVendir(confReader io.Reader, workingDir string) exec.CmdRunResult {
	var stdoutBs, stderrBs bytes.Buffer
	cmd := goexec.Command("vendir", "sync", "-f", "-", "--lock-file", "/dev/null")
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
