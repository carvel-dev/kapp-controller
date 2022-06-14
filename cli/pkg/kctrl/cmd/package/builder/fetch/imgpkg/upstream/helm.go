package upstream

import (
	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"strings"
)

type HelmStep struct {
	pkgAuthoringUI pkgui.IAuthoringUI
	pkgBuild       *pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewHelmStep(ui pkgui.IAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *HelmStep {
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
	contents := helmStep.pkgBuild.Spec.Vendir.Directories[0].Contents
	if contents == nil {
		helmStep.initializeContentWithHelmRelease()
	} else if contents[0].HelmChart == nil {
		helmStep.initializeHelmRelease()
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

func (helmStep *HelmStep) initializeHelmRelease() {
	helmChartContent := vendirconf.DirectoryContentsHelmChart{
		HelmVersion: "3",
	}
	helmStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].HelmChart = &helmChartContent
	helmStep.pkgBuild.WriteToFile()
}

func (helmStep *HelmStep) initializeContentWithHelmRelease() {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	helmStep.pkgBuild.Spec.Vendir.Directories[0].Contents = append(helmStep.pkgBuild.Spec.Vendir.Directories[0].Contents, vendirconf.DirectoryContents{})
	helmStep.initializeHelmRelease()
}

func (helmStep *HelmStep) configureHelmChartName() error {
	helmChartContent := helmStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].HelmChart
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

	helmChartContent.Name = strings.TrimSpace(name)
	helmStep.pkgBuild.WriteToFile()
	return nil
}

func (helmStep *HelmStep) configureHelmChartVersion() error {
	helmChartContent := helmStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].HelmChart
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

	helmChartContent.Version = strings.TrimSpace(version)
	helmStep.pkgBuild.WriteToFile()
	return nil
}

func (helmStep *HelmStep) configureHelmChartRepositoryURL() error {
	helmChartContent := helmStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].HelmChart
	if helmChartContent.Repository == nil {
		helmChartContent.Repository = &vendirconf.DirectoryContentsHelmChartRepo{}
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

	helmChartContent.Repository.URL = strings.TrimSpace(url)
	helmStep.pkgBuild.WriteToFile()
	return nil
}

func (helmStep *HelmStep) PostInteract() error {
	return nil
}
