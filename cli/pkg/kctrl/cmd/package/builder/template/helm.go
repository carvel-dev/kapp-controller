package template

import (
	"fmt"
	"path/filepath"

	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type HelmTemplateStep struct {
	pkgAuthoringUI pkgui.IAuthoringUI
	pkgBuild       *pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewHelmTemplateStep(ui pkgui.IAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *HelmTemplateStep {
	return &HelmTemplateStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
}

func (helmStep HelmTemplateStep) PreInteract() error {
	helmStep.pkgAuthoringUI.PrintInformationalText("\nSince we are using a helm chart, we have to add helmTemplate section. In this section, we will add helm Chart path inside package so that helm template can be run.")
	return nil
}

func (helmStep HelmTemplateStep) PostInteract() error {
	return nil
}

func (helmStep *HelmTemplateStep) Interact() error {
	existingPkgTemplates := helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template
	if isHelmTemplateExist(existingPkgTemplates) {
		if isHelmTemplateExistOnlyOnce(existingPkgTemplates) {
			helmStep.configureHelmChartPath()
		} else {
			//If there are > 1 helmTemplate section, then we dont want to touch them as they had been intentionally added.
			//TODO Rohit should we throw an error here?
			return nil
		}
	} else {
		helmStep.initializeHelmTemplate()
		helmStep.configureHelmChartPath()
	}

	return nil
}

func isHelmTemplateExist(existingTemplates []v1alpha1.AppTemplate) bool {
	for _, appTemplate := range existingTemplates {
		if appTemplate.HelmTemplate != nil {
			return true
		}
	}
	return false
}

func isHelmTemplateExistOnlyOnce(existingTemplates []v1alpha1.AppTemplate) bool {
	var count int
	for _, appTemplate := range existingTemplates {
		if appTemplate.HelmTemplate != nil {
			if count == 0 {
				count++
			} else {
				return false
			}
		}
	}
	return true
}

func (helmStep HelmTemplateStep) initializeHelmTemplate() {
	helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template = append([]v1alpha1.AppTemplate{
		v1alpha1.AppTemplate{HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{}}}, helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template...)
	helmStep.pkgBuild.WriteToFile()
}

func (helmStep HelmTemplateStep) configureHelmChartPath() error {
	var chartPath string
	for _, appTemplate := range helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template {
		if appTemplate.HelmTemplate != nil {

			path, err := helmStep.getPath()
			if err != nil {
				return err
			}
			chartPath = path

			appTemplate.HelmTemplate.Path = chartPath
			helmStep.pkgBuild.WriteToFile()
		}
	}
	return nil
}

func (helmStep HelmTemplateStep) getPath() (string, error) {
	if helmStep.pkgBuild.Spec.Vendir == nil {
		return "", fmt.Errorf("No Vendir configuration exist to get the path for helm template.")
	}

	directories := helmStep.pkgBuild.Spec.Vendir.Directories
	if directories == nil {
		return "", fmt.Errorf("No helm chart reference in Vendir")
	}
	var path string
	for _, directory := range directories {
		directoryPath := directory.Path
		for _, content := range directories[0].Contents {
			if content.HelmChart != nil {
				path = filepath.Join(directoryPath, content.Path)
				break
			}
		}
	}
	helmStep.pkgAuthoringUI.PrintInformationalText(fmt.Sprintf("As package is being created with immutable reference, helm chart is located in the %s directory. Add same in the helmTemplate path section", path))
	return path, nil
}
