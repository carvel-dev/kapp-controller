// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

type HelmStep struct {
	ui           cmdcore.AuthoringUI
	vendirConfig vendirconf.Config
}

func NewHelmStep(ui cmdcore.AuthoringUI, vendirConfig vendirconf.Config) *HelmStep {
	return &HelmStep{
		ui:           ui,
		vendirConfig: vendirConfig,
	}
}

func (helmStep *HelmStep) PreInteract() error {
	return nil
}

func (helmStep *HelmStep) Interact() error {
	contents := helmStep.vendirConfig.Directories[0].Contents
	if contents == nil {
		err := helmStep.initializeContentWithHelmRelease()
		if err != nil {
			return err
		}
	} else if contents[0].HelmChart == nil {
		err := helmStep.initializeHelmRelease()
		if err != nil {
			return err
		}
	}

	err := helmStep.configureHelmChartRepositoryURL()
	if err != nil {
		return err
	}
	err = helmStep.configureHelmChartName()
	if err != nil {
		return err
	}
	return helmStep.configureHelmChartVersion()
}

func (helmStep *HelmStep) initializeHelmRelease() error {
	helmChartContent := vendirconf.DirectoryContentsHelmChart{
		// TODO if the binary name is helm2, then it wont work.
		HelmVersion: "3",
	}
	helmStep.vendirConfig.Directories[0].Contents[0].HelmChart = &helmChartContent
	return SaveVendir(helmStep.vendirConfig)
}

func (helmStep *HelmStep) initializeContentWithHelmRelease() error {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	helmStep.vendirConfig.Directories[0].Contents = append(helmStep.vendirConfig.Directories[0].Contents, vendirconf.DirectoryContents{})
	return helmStep.initializeHelmRelease()
}

func (helmStep *HelmStep) configureHelmChartName() error {
	helmChartContent := helmStep.vendirConfig.Directories[0].Contents[0].HelmChart
	defaultName := helmChartContent.Name
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart name",
		Default:      defaultName,
		ValidateFunc: nil,
	}
	name, err := helmStep.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	helmChartContent.Name = strings.TrimSpace(name)
	return SaveVendir(helmStep.vendirConfig)
}

func (helmStep *HelmStep) configureHelmChartVersion() error {
	helmChartContent := helmStep.vendirConfig.Directories[0].Contents[0].HelmChart
	defaultVersion := helmChartContent.Version
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart version",
		Default:      defaultVersion,
		ValidateFunc: nil,
	}
	version, err := helmStep.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	helmChartContent.Version = strings.TrimSpace(version)
	return SaveVendir(helmStep.vendirConfig)
}

func (helmStep *HelmStep) configureHelmChartRepositoryURL() error {
	helmChartContent := helmStep.vendirConfig.Directories[0].Contents[0].HelmChart
	if helmChartContent.Repository == nil {
		helmChartContent.Repository = &vendirconf.DirectoryContentsHelmChartRepo{}
	}
	defaultUrl := helmChartContent.Repository.URL
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart repository URL",
		Default:      defaultUrl,
		ValidateFunc: nil,
	}
	url, err := helmStep.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	helmChartContent.Repository.URL = strings.TrimSpace(url)
	return SaveVendir(helmStep.vendirConfig)
}

func (helmStep *HelmStep) PostInteract() error {
	return nil
}
