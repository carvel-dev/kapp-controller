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
)

type ImgpkgRunner struct {
	Image             string
	Version           string
	Paths             []string
	UseKbldImagesLock bool
	ImgLockFilepath   string
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
		err = goexec.Command("cp", "-r", filepath.Join(wd, ".imgpkg"), dir).Run()
	}
	defer os.RemoveAll(dir)

	pushLocation := fmt.Sprintf("%s:%s", r.Image, r.Version)
	var stdoutBuf, stderrBuf bytes.Buffer
	inMemoryStdoutWriter := bufio.NewWriter(&stdoutBuf)
	cmd := goexec.Command("imgpkg", "push", "-b", pushLocation, "-f", dir, "--tty=true")
	// TODO: Switch to using Authoring UI
	fmt.Printf("Running: %s", strings.Join(cmd.Args, " "))
	cmd.Stdout = io.MultiWriter(os.Stdout, inMemoryStdoutWriter)
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		return stdoutBuf.String(), fmt.Errorf("%s", stderrBuf.String())
	}
	err = inMemoryStdoutWriter.Flush()
	if err != nil {
		return "", err
	}

	return stdoutBuf.String(), nil
}
