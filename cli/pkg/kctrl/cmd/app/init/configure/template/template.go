// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/configure/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/interfaces/build"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"strings"
)

const (
	UpstreamFolderName = "upstream"
	StdIn              = "-"
)

type TemplateStep struct {
	ui    cmdcore.AuthoringUI
	build build.Build
}

func NewTemplateStep(ui cmdcore.AuthoringUI, build build.Build) *TemplateStep {
	templateStep := TemplateStep{
		ui:    ui,
		build: build,
	}
	return &templateStep
}

func (templateStep *TemplateStep) PreInteract() error {
	return nil
}

func (templateStep *TemplateStep) Interact() error {
	return nil
}

func (templateStep *TemplateStep) PostInteract() error {
	appSpec := templateStep.build.GetAppSpec()
	if appSpec == nil {
		appSpec = &v1alpha1.AppSpec{}
	}
	existingTemplates := appSpec.Template

	/* In case of pkg init rerun, we assume this will be already populated and hence we return from here.
	We dont want to reset the modification user have done. */
	if len(existingTemplates) > 0 {
		return nil
	}

	var appTemplates []v1alpha1.AppTemplate
	if existingTemplates == nil {
		appTemplates = []v1alpha1.AppTemplate{}
	}
	var defaultYttPaths []string

	fetchSource := templateStep.build.GetObjectMeta().Annotations[fetch.FetchContentAnnotationKey]
	// Add helmTemplate
	switch fetchSource {
	case fetch.FetchChartFromHelmRepo:
		defaultYttPaths = []string{StdIn}
		appTemplateWithHelm := v1alpha1.AppTemplate{
			HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{
				Path: UpstreamFolderName,
			}}
		appTemplates = append(appTemplates, appTemplateWithHelm)
	case fetch.FetchFromLocalDirectory:
		templateStep.ui.PrintInformationalText("We need to include files/directory which should be part of this package. Multiple values can be included using a comma separator.")
		var defaultIncludedPath string
		textOpts := ui.TextOpts{
			Label:        "Enter the paths which need to be included as part of this package",
			Default:      defaultIncludedPath,
			ValidateFunc: nil,
		}
		includePaths, err := templateStep.ui.AskForText(textOpts)
		if err != nil {
			return err
		}
		defaultYttPaths = strings.Split(includePaths, ",")
	default:
		defaultYttPaths = []string{}
	}

	// Add yttTemplate
	appTemplateWithYtt := v1alpha1.AppTemplate{
		Ytt: &v1alpha1.AppTemplateYtt{
			Paths: defaultYttPaths,
		},
	}

	appTemplates = append(appTemplates, appTemplateWithYtt)

	// Add kbldTemplate
	appTemplates = append(appTemplates, v1alpha1.AppTemplate{Kbld: &v1alpha1.AppTemplateKbld{}})

	appSpec.Template = appTemplates
	templateStep.build.SetAppSpec(appSpec)
	return templateStep.build.Save()
}
