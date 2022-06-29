// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	"strings"

	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
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
	yttTemplateStep.appBuild.WriteToFile()
}

func (yttTemplateStep *YttTemplateStep) PostInteract() error {
	existingTemplates := yttTemplateStep.appBuild.Spec.App.Spec.Template
	fetchSource := yttTemplateStep.appBuild.ObjectMeta.Annotations[common.FetchContentAnnotationKey]
	//TODO how to handle if from local directory?
	if fetchSource == common.FetchFromLocalDirectory {
		if !isYttTemplateExist(existingTemplates) {
			yttTemplateStep.initializeYttTemplate()
		}
		return nil
	}

	if fetchSource == common.FetchChartFromHelmRepo || fetchSource == common.FetchChartFromGithub {
		//TODO in case of helm, check if any ytt path exist with '-', otherwise add it.
		return nil
	}

	yttTemplateStep.addUpstreamAsPathIfNotExist()
	return nil
}

func (yttTemplateStep *YttTemplateStep) addUpstreamAsPathIfNotExist() {
	//TODO I am consciously adding upstream folder. Should we do that even in scenario where user doesn't want explicitly. Cant think of the scenario though.
	appTemplates := yttTemplateStep.appBuild.Spec.App.Spec.Template
	for _, appTemplate := range appTemplates {
		if appTemplate.Ytt != nil {
			for _, path := range appTemplate.Ytt.Paths {
				if strings.HasPrefix(path, "upstream") {
					return
				}
			}
		}
	}
	appTemplateWithYtt := v1alpha1.AppTemplate{
		Ytt: &v1alpha1.AppTemplateYtt{
			Paths: []string{"upstream"},
		},
	}
	yttTemplateStep.appBuild.Spec.App.Spec.Template = append([]v1alpha1.AppTemplate{appTemplateWithYtt}, yttTemplateStep.appBuild.Spec.App.Spec.Template...)
	yttTemplateStep.appBuild.WriteToFile()
}
