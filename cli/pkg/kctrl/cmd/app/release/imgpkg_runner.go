// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"bytes"
	"fmt"
	"os"
	goexec "os/exec"
	"path/filepath"
	"strings"

	dircopy "github.com/otiai10/copy"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
)

type ImgpkgRunner struct {
	BundlePath        string
	Paths             []string
	UseKbldImagesLock bool
	ImgLockFilepath   string
	UI                cmdcore.AuthoringUI
}

func (r ImgpkgRunner) Run() (string, error) {
	tmpDir, err := os.MkdirTemp(".", fmt.Sprintf("bundle-%s-*", strings.Replace(r.BundlePath, "/", "-", -1)))
	if err != nil {
		return "", err
	}

	// Copy required files into temporary bundle location
	for _, path := range r.Paths {
		err = dircopy.Copy(path, tmpDir)
		if err != nil {
			return "", err
		}
	}
	if r.UseKbldImagesLock {
		err = os.Mkdir(filepath.Join(tmpDir, LockOutputFolder), os.ModePerm)
		if err != nil {
			return "", err
		}
		err = dircopy.Copy(r.ImgLockFilepath, filepath.Join(tmpDir, LockOutputFolder, LockOutputFile))
		if err != nil {
			return "", err
		}
	}
	defer os.RemoveAll(tmpDir)

	// TODO: Stream output
	r.UI.PrintInformationalText("\nAn imgpkg bundle consists of all required YAML configuration bundled into an OCI image " +
		"that can be pushed to an image registry and consumed by the package.\n")
	r.UI.PrintHeaderText("Pushing imgpkg bundle")

	cmd := goexec.Command("imgpkg", "push", "-b", r.BundlePath, "-f", tmpDir, "--tty=true")
	r.UI.PrintCmdExecutionOutput(fmt.Sprintf("$ %s", strings.Join(cmd.Args, " ")))

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		return stdoutBuf.String(), fmt.Errorf("%s", stderrBuf.String())
	}
	r.UI.PrintCmdExecutionOutput(stdoutBuf.String())

	bundleURL, err := r.imgpkgBundleURLFromStdout(stdoutBuf.String())
	if err != nil {
		return "", err
	}

	return bundleURL, nil
}

func (r *ImgpkgRunner) imgpkgBundleURLFromStdout(imgpkgStdout string) (string, error) {
	lines := strings.Split(imgpkgStdout, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Pushed") {
			line = strings.TrimPrefix(line, "Pushed")
			line = strings.Replace(line, "'", "", -1)
			line = strings.Replace(line, " ", "", -1)
			return line, nil
		}
	}
	return "", fmt.Errorf("Could not get imgpkg bundle location")
}
