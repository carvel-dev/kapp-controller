// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
)

type HelmTemplateStep struct {
	ui       cmdcore.AuthoringUI
	appBuild *appbuild.AppBuild
}

func NewHelmTemplateStep(ui cmdcore.AuthoringUI, appBuild *appbuild.AppBuild) *HelmTemplateStep {
	return &HelmTemplateStep{
		ui:       ui,
		appBuild: appBuild,
	}
}

func (helmStep HelmTemplateStep) PreInteract() error {
	helmStep.ui.PrintInformationalText("\nSince we are using a helm chart, we have to add helmTemplate section. In this section, we will add helm Chart path inside package so that helm template can be run.")
	return nil
}

func (helmStep *HelmTemplateStep) PostInteract() error {
	appTemplates := helmStep.appBuild.Spec.App.Spec.Template
	transformAppTemplates := NewTransformAppTemplates(&appTemplates)
	transformAppTemplates.AddUpstreamAsPathToHelmIfNotExist()
	helmStep.appBuild.Spec.App.Spec.Template = transformAppTemplates.GetAppTemplates()
	return helmStep.appBuild.Save()
}

func (helmStep *HelmTemplateStep) Interact() error {

	return nil
}
