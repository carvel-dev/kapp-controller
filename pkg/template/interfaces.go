// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"io"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

const (
	stdinPath = "-"
)

type Template interface {
	// TemplateDir works on directory returning templating result,
	// and boolean indicating whether subsequent operations
	// should operate on result, or continue operating on the directory
	TemplateDir(dirPath string) (exec.CmdRunResult, bool)
	// TemplateStream works on a stream returning templating result.
	// dirPath is provided for context from which to reference additonal inputs.
	TemplateStream(stream io.Reader, dirPath string) exec.CmdRunResult
}

type GenericOpts struct {
	Name      string
	Namespace string
}
