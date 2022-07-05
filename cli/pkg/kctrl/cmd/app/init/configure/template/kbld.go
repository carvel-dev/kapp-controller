// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type KbldTemplateStep struct {
	ui       cmdcore.AuthoringUI
	appBuild *appbuild.AppBuild
}

func NewKbldTemplateStep(ui cmdcore.AuthoringUI, appBuild *appbuild.AppBuild) *KbldTemplateStep {
	return &KbldTemplateStep{
		ui:       ui,
		appBuild: appBuild,
	}
}

func (kbldTemplateStep *KbldTemplateStep) PreInteract() error {
	return nil
}

func (kbldTemplateStep *KbldTemplateStep) PostInteract() error {
	return nil
}

func (kbldTemplateStep *KbldTemplateStep) Interact() error {
	existingTemplates := kbldTemplateStep.appBuild.Spec.App.Spec.Template
	if !isKbldTemplateExist(existingTemplates) {
		kbldTemplateStep.initializeKbldTemplate()
	}

	return nil
}

func isKbldTemplateExist(existingTemplates []v1alpha1.AppTemplate) bool {
	for _, appTemplate := range existingTemplates {
		if appTemplate.Kbld != nil {
			return true
		}
	}
	return false
}

func (kbldTemplateStep *KbldTemplateStep) initializeKbldTemplate() {
	kbldTemplateStep.appBuild.Spec.App.Spec.Template = append(kbldTemplateStep.appBuild.Spec.App.Spec.Template,
		v1alpha1.AppTemplate{Kbld: &v1alpha1.AppTemplateKbld{}})
	kbldTemplateStep.appBuild.Save()
}
