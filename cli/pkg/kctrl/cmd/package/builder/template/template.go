package template

import (
	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
)

const (
	Ytt int = iota
	HelmTemplate
)

type TemplateStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	pkgBuild       *pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewTemplateStep(ui pkgui.IPkgAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *TemplateStep {
	templateStep := TemplateStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
	return &templateStep
}

func (t TemplateStep) PreInteract() error {
	t.pkgAuthoringUI.PrintInformationalText("Next step is to define which templating tool to be used to render the package template. A package template can be rendered by various different tools.")
	return nil
}

func (templateStep TemplateStep) PostInteract() error {
	return nil
}

func (templateStep *TemplateStep) Interact() error {

	if templateStep.isHelmTemplateRequired() {
		helmTemplateStep := NewHelmTemplateStep(templateStep.pkgAuthoringUI, templateStep.pkgLocation, templateStep.pkgBuild)
		err := common.Run(helmTemplateStep)
		if err != nil {
			return err
		}
	}

	err := templateStep.configureYtt()
	if err != nil {
		return err
	}

	err = templateStep.configureKbld()
	if err != nil {
		return err
	}

	templateType, err := templateStep.pkgAuthoringUI.AskForChoice(ui.ChoiceOpts{
		Label:   "Enter the templating tool to be used",
		Default: 0,
		Choices: []string{"ytt(recommended)", "helmTemplate"},
	})
	if err != nil {
		return err
	}
	//var appTemplateList []v1alpha1.AppTemplate
	switch templateType {
	case Ytt:
		/*yttTemplateStep := NewYttTemplateStep(templateStep.pkgAuthoringUI)
		yttTemplateStep.Run()
		yttAppTemplate := v1alpha1.AppTemplate{
			Ytt: &yttTemplateStep.appTemplateYtt,
		}
		appTemplateList = append(appTemplateList, yttAppTemplate)
		templateStep.AppTemplate = appTemplateList*/
	case HelmTemplate:
	}
	return nil
}

func (templateStep TemplateStep) isHelmTemplateRequired() bool {
	return true
}

func (templateStep TemplateStep) configureKbld() error {
	return nil
}

func (templateStep TemplateStep) configureYtt() error {
	//yttTemplateStep := NewYttTemplateStep(templateStep.pkgAuthoringUI)
	return nil
}
