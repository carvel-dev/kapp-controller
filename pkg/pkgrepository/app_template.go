// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"fmt"
	"strings"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	ctltpl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
)

func (a *App) template(dirPath string) exec.CmdRunResult {
	if len(a.app.Spec.Template) == 0 {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Expected at least one template option"))
	}

	genericOpts := ctltpl.GenericOpts{Name: a.app.Name, Namespace: a.app.Namespace}

	var result exec.CmdRunResult
	var isStream bool

	for _, tpl := range a.app.Spec.Template {
		var template ctltpl.Template

		switch {
		case tpl.Ytt != nil:
			template = a.templateFactory.NewYtt(*tpl.Ytt, genericOpts)
		case tpl.Kbld != nil:
			template = a.templateFactory.NewKbld(*tpl.Kbld, genericOpts)
		default:
			return exec.NewCmdRunResultWithErr(fmt.Errorf("Unsupported way to template"))
		}

		if isStream {
			result = template.TemplateStream(strings.NewReader(result.Stdout), dirPath)
		} else {
			result, isStream = template.TemplateDir(dirPath)
		}
		if result.Error != nil {
			break
		}
	}

	return result
}
