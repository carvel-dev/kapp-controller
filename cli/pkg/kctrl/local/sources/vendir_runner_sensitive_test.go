// Copyright 2026 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sources

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cmdcore "carvel.dev/kapp-controller/cli/pkg/kctrl/cmd/core"
)

type vendirCaptureUI struct {
	cmdcore.AuthoringUI
	output bytes.Buffer
}

func (u *vendirCaptureUI) PrintInformationalText(s string)  { u.output.WriteString(s) }
func (u *vendirCaptureUI) PrintActionableText(s string)     { u.output.WriteString(s) }
func (u *vendirCaptureUI) PrintCmdExecutionText(s string)   { u.output.WriteString(s) }
func (u *vendirCaptureUI) PrintCmdExecutionOutput(s string) { u.output.WriteString(s) }
func TestVendirRunnerDoesNotPrintConfiguration(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.WriteFile(vendirFileName, []byte("https://user:password-canary@example.com?token=query-canary"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "vendir"), []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	capture := &vendirCaptureUI{}
	if err := NewVendirRunner(capture).Sync("git"); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(capture.output.String(), "canary") {
		t.Fatal("configuration contents exposed")
	}
}
