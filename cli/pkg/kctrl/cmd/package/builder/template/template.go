package template

import (
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

func (templateStep TemplateStep) PreInteract() error {
	//templateStep.pkgAuthoringUI.PrintInformationalText("\nNext step is to add the template section of the package.")
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

	return nil
}

func (templateStep TemplateStep) isHelmTemplateRequired() bool {
	if templateStep.pkgBuild.Annotations[common.PkgFetchContentAnnotationKey] == common.FetchChartFromHelmRepo {
		return true
	}
	return false
}

func (templateStep TemplateStep) configureKbld() error {
	kbldTemplateStep := NewKbldTemplateStep(templateStep.pkgAuthoringUI, templateStep.pkgLocation, templateStep.pkgBuild)
	return common.Run(kbldTemplateStep)
}

func (templateStep TemplateStep) configureYtt() error {
	yttTemplateStep := NewYttTemplateStep(templateStep.pkgAuthoringUI, templateStep.pkgLocation, templateStep.pkgBuild)
	return common.Run(yttTemplateStep)
}
