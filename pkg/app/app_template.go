// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"strings"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"carvel.dev/kapp-controller/pkg/exec"
	ctltpl "carvel.dev/kapp-controller/pkg/template"
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

	additionalValues := a.buildDownwardAPIAdditionalValues()

	for _, tpl := range a.app.Spec.Template {
		var template ctltpl.Template

		switch {
		case tpl.Ytt != nil:
			template = a.templateFactory.NewYtt(*tpl.Ytt, appContext, additionalValues)
		case tpl.Kbld != nil:
			template = a.templateFactory.NewKbld(*tpl.Kbld, appContext)
		case tpl.HelmTemplate != nil:
			template = a.templateFactory.NewHelmTemplate(*tpl.HelmTemplate, appContext, additionalValues)
		case tpl.Sops != nil:
			template = a.templateFactory.NewSops(*tpl.Sops, appContext)
		case tpl.Cue != nil:
			template = a.templateFactory.NewCue(*tpl.Cue, appContext, additionalValues)
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

func (a *App) buildDownwardAPIAdditionalValues() ctltpl.AdditionalDownwardAPIValues {
	return ctltpl.AdditionalDownwardAPIValues{
		KubernetesVersion: func() (string, error) {
			if a.memoizedKubernetesVersion == "" {
				v, err := a.compInfo.KubernetesVersion(a.app.Spec.ServiceAccountName, a.app.Spec.Cluster, &a.app.ObjectMeta)
				if err != nil {
					return "", fmt.Errorf("Unable to get kubernetes version: %s", err)
				}
				a.memoizedKubernetesVersion = v.String()
			}
			return a.memoizedKubernetesVersion, nil
		},
		KappControllerVersion: func() (string, error) {
			v, err := a.compInfo.KappControllerVersion()
			if err != nil {
				return "", fmt.Errorf("Unable to get kapp-controller version: %s", err)
			}
			return v.String(), nil
		},
		KubernetesAPIs: func() ([]string, error) {
			if len(a.memoizedKubernetesAPIs) == 0 {
				v, err := a.compInfo.KubernetesAPIs()
				if err != nil {
					return []string{}, fmt.Errorf("Unable to list all server apigroups/version: %s", err)
				}
				a.memoizedKubernetesAPIs = v
			}
			return a.memoizedKubernetesAPIs, nil
		},
	}
}
