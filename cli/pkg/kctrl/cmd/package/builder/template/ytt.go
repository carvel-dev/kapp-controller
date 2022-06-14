package template

import (
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type YttTemplateStep struct {
	pkgAuthoringUI pkgui.IAuthoringUI
	pkgBuild       *pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewYttTemplateStep(ui pkgui.IAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *YttTemplateStep {
	return &YttTemplateStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
}

func (yttTemplateStep YttTemplateStep) PreInteract() error {
	return nil
}

func (yttTemplateStep YttTemplateStep) PostInteract() error {
	return nil
}

func (yttTemplateStep *YttTemplateStep) Interact() error {
	existingPkgTemplates := yttTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template
	if isYttTemplateExist(existingPkgTemplates) {
		return nil
	} else {
		yttTemplateStep.initializeYttTemplate()
		yttTemplateStep.configureYttPath()
	}
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
	yttTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template = append(yttTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template,
		v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{}})
	yttTemplateStep.pkgBuild.WriteToFile()
}

func (yttTemplateStep *YttTemplateStep) configureYttPath() error {
	for _, appTemplate := range yttTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template {
		if appTemplate.Ytt != nil {
			if yttTemplateStep.pkgBuild.Annotations[common.PkgFetchContentAnnotationKey] == common.FetchChartFromHelmRepo {
				appTemplate.Ytt.Paths = []string{"-"}
			} else {
				appTemplate.Ytt.Paths = []string{"config"}
			}
			yttTemplateStep.pkgBuild.WriteToFile()
		}
	}
	return nil
}
