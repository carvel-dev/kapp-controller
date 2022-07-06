// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/configure/fetch"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type YttTemplateStep struct {
	ui       cmdcore.AuthoringUI
	appBuild *appbuild.AppBuild
}

func NewYttTemplateStep(ui cmdcore.AuthoringUI, appBuild *appbuild.AppBuild) *YttTemplateStep {
	return &YttTemplateStep{
		ui:       ui,
		appBuild: appBuild,
	}
}

func (yttTemplateStep *YttTemplateStep) PreInteract() error {
	return nil
}

func (yttTemplateStep *YttTemplateStep) Interact() error {
	return nil
}

func (yttTemplateStep *YttTemplateStep) PostInteract() error {
	existingTemplates := yttTemplateStep.appBuild.Spec.App.Spec.Template
	fetchSource := yttTemplateStep.appBuild.ObjectMeta.Annotations[fetch.FetchContentAnnotationKey]
	//TODO how to handle if from local directory?
	if fetchSource == fetch.FetchFromLocalDirectory {
		if !isYttTemplateExist(existingTemplates) {
			yttTemplateStep.initializeYttTemplate()
		}
		return nil
	}

	if fetchSource == fetch.FetchChartFromHelmRepo || fetchSource == fetch.FetchChartFromGithub {
		yttTemplateStep.addStdInAsPathToYttTemplateIfNotExist()
		return nil
	}

	yttTemplateStep.addUpstreamAsPathToYttTemplateIfNotExist()
	return nil
}

func (yttTemplateStep *YttTemplateStep) addStdInAsPathToYttTemplateIfNotExist() {
	appTemplates := yttTemplateStep.appBuild.Spec.App.Spec.Template
	transformAppTemplates := NewTransformAppTemplates(&appTemplates)
	transformAppTemplates.AddStdInAsPathToYttIfNotExist()
	yttTemplateStep.appBuild.Spec.App.Spec.Template = transformAppTemplates.GetAppTemplates()
	yttTemplateStep.appBuild.Save()
	return
}

func isYttTemplateExist(existingTemplates []v1alpha1.AppTemplate) bool {
	for _, appTemplate := range existingTemplates {
		if appTemplate.Ytt != nil {
			return true
		}
	}
	return false
}

func (yttTemplateStep *YttTemplateStep) initializeYttTemplate() {
	yttTemplateStep.appBuild.Spec.App.Spec.Template = append(yttTemplateStep.appBuild.Spec.App.Spec.Template,
		v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{}})
	yttTemplateStep.appBuild.Save()
}

func (yttTemplateStep *YttTemplateStep) addUpstreamAsPathToYttTemplateIfNotExist() {
	appTemplates := yttTemplateStep.appBuild.Spec.App.Spec.Template
	//TODO I am consciously adding upstream folder. Should we do that even in scenario where user doesn't want explicitly. Cant think of the scenario though.
	transformAppTemplates := NewTransformAppTemplates(&appTemplates)
	transformAppTemplates.AddUpstreamAsPathToYttIfNotExist()
	yttTemplateStep.appBuild.Spec.App.Spec.Template = transformAppTemplates.GetAppTemplates()
	yttTemplateStep.appBuild.Save()
}
