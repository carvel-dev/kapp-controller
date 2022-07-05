// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/configure/fetch"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
)

type TemplateStep struct {
	ui       cmdcore.AuthoringUI
	appBuild *appbuild.AppBuild
}

func NewTemplateStep(ui cmdcore.AuthoringUI, appBuild *appbuild.AppBuild) *TemplateStep {
	templateStep := TemplateStep{
		ui:       ui,
		appBuild: appBuild,
	}
	return &templateStep
}

func (templateStep *TemplateStep) PreInteract() error {
	return nil
}

func (templateStep *TemplateStep) Interact() error {
	if templateStep.appBuild.Spec.App == nil {
		templateStep.appBuild.Spec.App = &v1alpha12.AppTemplateSpec{}
	}
	existingTemplates := templateStep.appBuild.Spec.App.Spec.Template
	if existingTemplates == nil {
		appTemplate := []v1alpha1.AppTemplate{}
		templateStep.appBuild.Spec.App.Spec.Template = appTemplate
	}
	if templateStep.isHelmTemplateRequired() {
		//TODO Add code here to add helm Template
	}
	err := templateStep.configureYtt()
	if err != nil {
		return err
	}

	err = templateStep.configureKbld()
	if err != nil {
		return err
	}

	return nil
}

func (templateStep TemplateStep) isHelmTemplateRequired() bool {
	if templateStep.appBuild.ObjectMeta.Annotations[fetch.FetchContentAnnotationKey] == fetch.FetchChartFromHelmRepo {
		return true
	}
	return false
}

func (templateStep TemplateStep) configureYtt() error {
	yttTemplateStep := NewYttTemplateStep(templateStep.ui, templateStep.appBuild)
	return common.Run(yttTemplateStep)
}

func (templateStep TemplateStep) configureKbld() error {
	kbldTemplateStep := NewKbldTemplateStep(templateStep.ui, templateStep.appBuild)
	return common.Run(kbldTemplateStep)
}

func (templateStep *TemplateStep) PostInteract() error {
	return nil
}
