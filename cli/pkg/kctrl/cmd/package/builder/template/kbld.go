package template

import (
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type KbldTemplateStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	pkgBuild       *pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewKbldTemplateStep(ui pkgui.IPkgAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *KbldTemplateStep {
	return &KbldTemplateStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
}

func (kbldTemplateStep KbldTemplateStep) PreInteract() error {
	kbldTemplateStep.pkgAuthoringUI.PrintInformationalText(`We need to provide the values to ytt.`)
	return nil
}

func (kbldTemplateStep KbldTemplateStep) PostInteract() error {
	return nil
}

func (kbldTemplateStep *KbldTemplateStep) Interact() error {
	existingPkgTemplates := kbldTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template
	//TODO Rohit needs to be removed
	if existingPkgTemplates == nil {
		kbldTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template = append(kbldTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template, v1alpha1.AppTemplate{})
	}
	if isKbldTemplateExist(existingPkgTemplates) {
		if isKbldTemplateExistOnlyOnce(existingPkgTemplates) {
			kbldTemplateStep.configureKbldPaths()
		} else {
			//If there are > 1 helmTemplate section, then we dont want to touch them as they had been intentionally added.
			//TODO Rohit should we throw an error here?
			return nil
		}
	} else {
		kbldTemplateStep.initializeKbldTemplate()
		kbldTemplateStep.configureKbldPaths()
	}

	kbldTemplateStep.pkgAuthoringUI.PrintInformationalText("Adding path to the ytt template section")
	return nil
}

func (kbldTemplateStep *KbldTemplateStep) askForPath() bool {
	return false
}

func isKbldTemplateExist(existingTemplates []v1alpha1.AppTemplate) bool {
	for _, appTemplate := range existingTemplates {
		if appTemplate.Ytt != nil {
			return true
		}
	}
	return false
}

func isKbldTemplateExistOnlyOnce(existingTemplates []v1alpha1.AppTemplate) bool {
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

func (kbldTemplateStep *KbldTemplateStep) initializeKbldTemplate() {
	kbldTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template = append(kbldTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template,
		v1alpha1.AppTemplate{Kbld: &v1alpha1.AppTemplateKbld{}})
	kbldTemplateStep.pkgBuild.WriteToFile(kbldTemplateStep.pkgLocation)
}

func (kbldTemplateStep *KbldTemplateStep) configureKbldPaths() error {
	for _, appTemplate := range kbldTemplateStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template {
		if appTemplate.Kbld != nil {
			defaultPaths := appTemplate.Kbld.Paths
			//If paths are already configured, do not change them
			if len(defaultPaths) == 0 {
				defaultPaths = append(defaultPaths, "-")
			}
			//TODO Rohit check if this works
			appTemplate.Kbld.Paths = defaultPaths
			kbldTemplateStep.pkgBuild.WriteToFile(kbldTemplateStep.pkgLocation)
		}
	}
	return nil
}
