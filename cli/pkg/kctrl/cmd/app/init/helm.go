// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

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

func (h *HelmStep) PreInteract() error { return nil }

func (h *HelmStep) Interact() error {
	contents := h.vendirConfig.Directories[0].Contents
	if contents == nil {
		err := h.initializeContentWithHelmRelease()
		if err != nil {
			return err
		}
	} else if contents[0].HelmChart == nil {
		err := h.initializeHelmRelease()
		if err != nil {
			return err
		}
	}

	err := h.configureHelmChartRepositoryURL()
	if err != nil {
		return err
	}
	err = h.configureHelmChartName()
	if err != nil {
		return err
	}
	return h.configureHelmChartVersion()
}

func (h *HelmStep) initializeHelmRelease() error {
	helmChartContent := vendirconf.DirectoryContentsHelmChart{}
	h.vendirConfig.Directories[0].Contents[0].HelmChart = &helmChartContent
	return SaveVendir(h.vendirConfig)
}

func (h *HelmStep) initializeContentWithHelmRelease() error {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	h.vendirConfig.Directories[0].Contents = append(h.vendirConfig.Directories[0].Contents, vendirconf.DirectoryContents{})
	return h.initializeHelmRelease()
}

func (h *HelmStep) configureHelmChartName() error {
	helmChartContent := h.vendirConfig.Directories[0].Contents[0].HelmChart
	defaultName := helmChartContent.Name
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart name",
		Default:      defaultName,
		ValidateFunc: nil,
	}
	name, err := h.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	helmChartContent.Name = strings.TrimSpace(name)
	return SaveVendir(h.vendirConfig)
}

func (h *HelmStep) configureHelmChartVersion() error {
	helmChartContent := h.vendirConfig.Directories[0].Contents[0].HelmChart
	defaultVersion := helmChartContent.Version
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart version",
		Default:      defaultVersion,
		ValidateFunc: nil,
	}
	version, err := h.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	helmChartContent.Version = strings.TrimSpace(version)
	return SaveVendir(h.vendirConfig)
}

func (h *HelmStep) configureHelmChartRepositoryURL() error {
	helmChartContent := h.vendirConfig.Directories[0].Contents[0].HelmChart
	if helmChartContent.Repository == nil {
		helmChartContent.Repository = &vendirconf.DirectoryContentsHelmChartRepo{}
	}
	defaultUrl := helmChartContent.Repository.URL
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart repository URL",
		Default:      defaultUrl,
		ValidateFunc: nil,
	}
	url, err := h.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	helmChartContent.Repository.URL = strings.TrimSpace(url)
	return SaveVendir(h.vendirConfig)
}

func (h *HelmStep) PostInteract() error { return nil }
