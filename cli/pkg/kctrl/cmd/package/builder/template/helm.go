package template

import (
	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

const (
	YttFilesLocation int = iota
	Inline
)

type HelmTemplateStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	pkgBuild       pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewHelmTemplateStep(ui pkgui.IPkgAuthoringUI) *HelmTemplateStep {
	return &HelmTemplateStep{
		pkgAuthoringUI: ui,
	}
}

func (helmStep HelmTemplateStep) PreInteract() error {
	helmStep.pkgAuthoringUI.PrintInformationalText(`We need to add path to the helm Chart so that helm template can be run.`)
	return nil
}

func (helmStep HelmTemplateStep) PostInteract() error {
	return nil
}

func (helmStep *HelmTemplateStep) Interact() error {
	existingPkgTemplates := helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template
	if existingPkgTemplates == nil {
		helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template = append(helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template, v1alpha1.AppTemplate{})
	}
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

	helmStep.pkgAuthoringUI.PrintInformationalText(`Adding path to the helm template section`)
	return nil
}

func (helmStep HelmTemplateStep) askForPath() bool {
	return true
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

func (y *HelmTemplateStep) Run() error {
	y.PreInteract()
	y.Interact()
	y.PostInteract()
	return nil
}

func (helmStep HelmTemplateStep) initializeHelmTemplate() {
	helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template = append(helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template,
		v1alpha1.AppTemplate{HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{}})
	helmStep.pkgBuild.WriteToFile(helmStep.pkgLocation)
}

func (helmStep HelmTemplateStep) configureHelmChartPath() error {
	var path string
	for _, appTemplate := range helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Template {
		if appTemplate.HelmTemplate != nil {
			defaultPath := appTemplate.HelmTemplate.Path
			if helmStep.askForPath() {
				input, err := helmStep.pkgAuthoringUI.AskForText(ui.TextOpts{
					Label:        "Enter the path where helm chart is located",
					Default:      defaultPath,
					ValidateFunc: nil, //I think we can add some validation.
				})
				if err != nil {
					return err
				}
				path = input
			} else {
				path = getPathFromVendir(helmStep.pkgBuild)
			}
			//TODO Rohit check if this works
			appTemplate.HelmTemplate.Path = path
			helmStep.pkgBuild.WriteToFile(helmStep.pkgLocation)
		}
	}
	return nil
}

func getPathFromVendir(pkgBuild pkgbuilder.PackageBuild) string {
	return ""
}
