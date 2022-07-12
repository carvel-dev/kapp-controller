// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	goexec "os/exec"
	"path/filepath"
	"strings"

	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
)

type ImgpkgRunner struct {
	Image             string
	Version           string
	Paths             []string
	UseKbldImagesLock bool
	ImgLockFilepath   string
	UI                cmdcore.AuthoringUI
}

func (r ImgpkgRunner) Run() (string, error) {
	dir, err := os.MkdirTemp(".", fmt.Sprintf("bundle-%s-*", strings.Replace(r.Image, "/", "-", 1)))
	if err != nil {
		return "", err
	}

	wd, err := os.Getwd()
	for _, path := range r.Paths {
		var stderrBuf bytes.Buffer
		cmd := goexec.Command("cp", "-r", filepath.Join(wd, path), dir)
		cmd.Stderr = &stderrBuf
		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("%s", stderrBuf.String())
		}
	}
	if r.UseKbldImagesLock {
		err = goexec.Command("mkdir", filepath.Join(dir, lockOutputFolder)).Run()
		if err != nil {
			return "", err
		}
		err = goexec.Command("cp", r.ImgLockFilepath, filepath.Join(dir, lockOutputFolder, "images.yml")).Run()
		if err != nil {
			return "", err
		}
	}
	defer os.RemoveAll(dir)

	pushLocation := fmt.Sprintf("%s:%s", r.Image, r.Version)
	var stdoutBuf, stderrBuf bytes.Buffer
	inMemoryStdoutWriter := bufio.NewWriter(&stdoutBuf)
	cmd := goexec.Command("imgpkg", "push", "-b", pushLocation, "-f", dir, "--tty=true")
	// TODO: Stream output
	r.UI.PrintInformationalText("\nAn imgpkg bundle consists of all required YAML configuration bundled into an OCI image" +
		"that can be pushed to an image registry and consumed by the package.\n")
	r.UI.PrintHeaderText("Pushing imgpkg bundle")
	r.UI.PrintCmdExecutionOutput(fmt.Sprintf("$ %s", strings.Join(cmd.Args, " ")))
	cmd.Stdout = io.MultiWriter(inMemoryStdoutWriter)
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		return stdoutBuf.String(), fmt.Errorf("%s", stderrBuf.String())
	}
	err = inMemoryStdoutWriter.Flush()
	if err != nil {
		return "", err
	}
	r.UI.PrintCmdExecutionOutput(stdoutBuf.String())

	return stdoutBuf.String(), nil
}
