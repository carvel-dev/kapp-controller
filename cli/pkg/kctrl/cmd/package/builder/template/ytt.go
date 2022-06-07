package template

import (
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"

	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
)

type YttTemplateStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	pkgBuild       *pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewYttTemplateStep(ui pkgui.IPkgAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *YttTemplateStep {
	return &YttTemplateStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
}

func (yttTemplateStep YttTemplateStep) PreInteract() error {
	yttTemplateStep.pkgAuthoringUI.PrintInformationalText(`We need to provide the values to ytt.`)
	return nil
}

func (yttTemplateStep YttTemplateStep) PostInteract() error {
	return nil
}

func (yttTemplateStep *YttTemplateStep) Interact() error {
	existingPkgTemplates := yttTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template
	//TODO Rohit needs to be removed
	if existingPkgTemplates == nil {
		yttTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template = append(yttTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template, v1alpha1.AppTemplate{})
	}
	if isYttTemplateExist(existingPkgTemplates) {
		if isYttTemplateExistOnlyOnce(existingPkgTemplates) {
			yttTemplateStep.configureYttPath()
		} else {
			//If there are > 1 helmTemplate section, then we dont want to touch them as they had been intentionally added.
			//TODO Rohit should we throw an error here?
			return nil
		}
	} else {
		yttTemplateStep.initializeYttTemplate()
		yttTemplateStep.configureYttPath()
	}

	yttTemplateStep.pkgAuthoringUI.PrintInformationalText("Adding path to the ytt template section")
	return nil
}

func (yttTemplateStep *YttTemplateStep) askForPath() bool {
	return false
}

func isYttTemplateExist(existingTemplates []v1alpha1.AppTemplate) bool {
	for _, appTemplate := range existingTemplates {
		if appTemplate.Ytt != nil {
			return true
		}
	}
	return false
}

func isYttTemplateExistOnlyOnce(existingTemplates []v1alpha1.AppTemplate) bool {
	var count int
	for _, appTemplate := range existingTemplates {
		if appTemplate.Ytt != nil {
			if count == 0 {
				count++
			} else {
				return false
			}
		}
	}
	return true
}

func (yttTemplateStep *YttTemplateStep) initializeYttTemplate() {
	yttTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template = append(yttTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template,
		v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{}})
	yttTemplateStep.pkgBuild.WriteToFile(yttTemplateStep.pkgLocation)
}

func (yttTemplateStep *YttTemplateStep) configureYttPath() error {
	for _, appTemplate := range yttTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template {
		if appTemplate.Ytt != nil {
			defaultPaths := appTemplate.Ytt.Paths
			if yttTemplateStep.askForPath() {

			} else {
				if len(defaultPaths) == 0 {
					defaultPaths = append(defaultPaths, yttTemplateStep.getPathFromFetchConf())
				}
			}
			//TODO Rohit check if this works
			appTemplate.Ytt.Paths = defaultPaths
			yttTemplateStep.pkgBuild.WriteToFile(yttTemplateStep.pkgLocation)
		}
	}
	return nil
}

func (yttTemplateStep *YttTemplateStep) getPathFromFetchConf() string {
	//This means that helmChart has been mentioned directly in the fetch section of Package.
	if yttTemplateStep.pkgBuild.Spec.Vendir == nil {
		return "-"
	}

	directories := yttTemplateStep.pkgBuild.Spec.Vendir.Directories
	if directories == nil {
		return "-"
	}

	var path string
	for _, directory := range directories {
		directoryPath := directory.Path
		for _, content := range directories[0].Contents {
			if content.Directory != nil {
				path = directoryPath + "/" + content.Path
				break
			}
		}
	}
	return path
}
