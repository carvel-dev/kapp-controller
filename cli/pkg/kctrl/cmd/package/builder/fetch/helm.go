package fetch

import (
	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type HelmStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	pkgBuild       *pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewHelmStep(ui pkgui.IPkgAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *HelmStep {
	return &HelmStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
}

func (helmStep *HelmStep) PreInteract() error {
	return nil
}

func (helmStep *HelmStep) Interact() error {
	fetchContents := helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch
	if fetchContents == nil {
		helmStep.initializeFetchContentWithHelmRelease()
	} else if fetchContents[0].HelmChart == nil {
		helmStep.initializeHelmChart()
	}

	err := helmStep.configureHelmChartName()
	if err != nil {
		return err
	}
	err = helmStep.configureHelmChartVersion()
	if err != nil {
		return err
	}
	err = helmStep.configureHelmChartRepositoryURL()
	if err != nil {
		return err
	}
	return nil
}

func (helmStep *HelmStep) initializeHelmChart() {
	helmChartContent := v1alpha1.AppFetchHelmChart{}
	helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch[0].HelmChart = &helmChartContent
	helmStep.pkgBuild.WriteToFile(helmStep.pkgLocation)
}

func (helmStep *HelmStep) initializeFetchContentWithHelmRelease() {
	appFetch := v1alpha1.AppFetch{}
	helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch = append(helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch, appFetch)
	helmStep.initializeHelmChart()
}

func (helmStep *HelmStep) configureHelmChartName() error {
	helmChartContent := helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch[0].HelmChart
	defaultName := helmChartContent.Name
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart name",
		Default:      defaultName,
		ValidateFunc: nil,
	}
	name, err := helmStep.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}

	helmChartContent.Name = name
	helmStep.pkgBuild.WriteToFile(helmStep.pkgLocation)
	return nil
}

func (helmStep *HelmStep) configureHelmChartVersion() error {
	helmChartContent := helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch[0].HelmChart
	defaultVersion := helmChartContent.Version
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart version",
		Default:      defaultVersion,
		ValidateFunc: nil,
	}
	version, err := helmStep.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}

	helmChartContent.Version = version
	helmStep.pkgBuild.WriteToFile(helmStep.pkgLocation)
	return nil
}

func (helmStep *HelmStep) configureHelmChartRepositoryURL() error {
	helmChartContent := helmStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch[0].HelmChart
	if helmChartContent.Repository == nil {
		helmChartContent.Repository = &v1alpha1.AppFetchHelmChartRepo{}
	}
	defaultUrl := helmChartContent.Repository.URL
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart repository URL",
		Default:      defaultUrl,
		ValidateFunc: nil,
	}
	url, err := helmStep.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}

	helmChartContent.Repository.URL = url
	helmStep.pkgBuild.WriteToFile(helmStep.pkgLocation)
	return nil
}

func (helmStep *HelmStep) PostInteract() error {
	return nil
}
