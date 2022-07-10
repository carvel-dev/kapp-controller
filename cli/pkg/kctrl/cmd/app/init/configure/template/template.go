// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/configure/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/interfaces/build"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
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
	defaultYttPath := UpstreamFolderName

	// Add helmTemplate
	if templateStep.isHelmTemplateRequired() {
		defaultYttPath = StdIn
		appTemplateWithHelm := v1alpha1.AppTemplate{
			HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{
				Path: UpstreamFolderName,
			}}
		appTemplates = append(appTemplates, appTemplateWithHelm)
	}

	// Add yttTemplate
	appTemplateWithYtt := v1alpha1.AppTemplate{
		Ytt: &v1alpha1.AppTemplateYtt{
			Paths: []string{defaultYttPath},
		}}
	appTemplates = append(appTemplates, appTemplateWithYtt)

	// Add kbldTemplate
	appTemplates = append(appTemplates, v1alpha1.AppTemplate{Kbld: &v1alpha1.AppTemplateKbld{}})

	appSpec.Template = appTemplates
	templateStep.build.SetAppSpec(appSpec)
	return templateStep.build.Save()
}

func (templateStep TemplateStep) isHelmTemplateRequired() bool {
	if templateStep.build.GetObjectMeta().Annotations[fetch.FetchContentAnnotationKey] == fetch.FetchChartFromHelmRepo {
		return true
	}
	return false
}
