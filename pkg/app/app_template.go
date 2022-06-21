// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"strings"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	ctltpl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
)

func (a *App) template(dirPath string) exec.CmdRunResult {
	if len(a.app.Spec.Template) == 0 {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Expected at least one template option"))
	}

	appContext := ctltpl.AppContext{
		Name:      a.app.Name,
		Namespace: a.app.Namespace,
		Metadata:  asPartialObjectMetadata(a.app),
	}

	var result exec.CmdRunResult
	var isStream bool

	for _, tpl := range a.app.Spec.Template {
		var template ctltpl.Template

		switch {
		case tpl.Ytt != nil:
			template = a.templateFactory.NewYtt(*tpl.Ytt, appContext)
		case tpl.Kbld != nil:
			template = a.templateFactory.NewKbld(*tpl.Kbld, appContext)
		case tpl.HelmTemplate != nil:
			template = a.templateFactory.NewHelmTemplate(*tpl.HelmTemplate, appContext)
		case tpl.Sops != nil:
			template = a.templateFactory.NewSops(*tpl.Sops, appContext)
		case tpl.Cue != nil:
			template = a.templateFactory.NewCue(*tpl.Cue, appContext)
		default:
			result.AttachErrorf("%s", fmt.Errorf("Unsupported way to template"))
			return result
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

func asPartialObjectMetadata(m v1alpha1.App) ctltpl.PartialObjectMetadata {
	return ctltpl.PartialObjectMetadata{
		ObjectMeta: ctltpl.ObjectMeta{
			Name:        m.GetName(),
			Namespace:   m.GetNamespace(),
			UID:         m.GetUID(),
			Labels:      m.GetLabels(),
			Annotations: m.GetAnnotations(),
		},
	}
}
