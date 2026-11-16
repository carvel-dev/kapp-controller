// Copyright 2026 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package dev

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	kcv1alpha1 "carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	fakekc "carvel.dev/kapp-controller/pkg/client/clientset/versioned/fake"
	"github.com/cppforlife/go-cli-ui/ui"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type appCaptureUI struct {
	ui.UI
	output bytes.Buffer
}

func (u *appCaptureUI) PrintBlock(b []byte)                      { u.output.Write(b) }
func (u *appCaptureUI) PrintLinef(f string, args ...interface{}) { fmt.Fprintf(&u.output, f, args...) }
func TestDevResourceDisplayDoesNotPrintInlineFiles(t *testing.T) {
	capture := &appCaptureUI{UI: ui.NewNoopUI()}
	app := &kcv1alpha1.App{ObjectMeta: metav1.ObjectMeta{Name: "sample", Namespace: "default"}, Spec: kcv1alpha1.AppSpec{Fetch: []kcv1alpha1.AppFetch{{Inline: &kcv1alpha1.AppFetchInline{Paths: map[string]string{"credentials.txt": "password-canary"}}}}}}
	client := fakekc.NewSimpleClientset(app)
	options := &DevOptions{ui: capture}
	if err := options.printRs(app.ObjectMeta, client); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(capture.output.String(), "canary") || !strings.Contains(capture.output.String(), "app/sample") {
		t.Fatal("unsafe resource display")
	}
}
