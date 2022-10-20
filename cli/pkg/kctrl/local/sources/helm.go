// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sources

import (
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

type HelmSource struct {
	ui           cmdcore.AuthoringUI
	vendirConfig *VendirConfig
}

var _ SourceConfiguration = &HelmSource{}

func NewHelmSource(ui cmdcore.AuthoringUI, vendirConfig *VendirConfig) *HelmSource {
	return &HelmSource{ui: ui, vendirConfig: vendirConfig}
}

func (h *HelmSource) Configure() error {
	contents := h.vendirConfig.Contents()
	if contents == nil {
		err := h.initializeContentWithHelmRelease()
		if err != nil {
			return err
		}
	} else if contents[0].HelmChart == nil {
		err := h.initializeHelmRelease(contents)
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

func (h *HelmSource) initializeHelmRelease(contents []vendirconf.DirectoryContents) error {
	contents[0].HelmChart = &vendirconf.DirectoryContentsHelmChart{}
	h.vendirConfig.SetContents(contents)
	return h.vendirConfig.Save()
}

func (h *HelmSource) initializeContentWithHelmRelease() error {
	contents := h.vendirConfig.Contents()
	if contents == nil {
		contents = append(h.vendirConfig.Contents(), vendirconf.DirectoryContents{})
	}
	if contents[0].HelmChart == nil {
		contents[0].HelmChart = &vendirconf.DirectoryContentsHelmChart{}
	}
	h.vendirConfig.SetContents(contents)
	return h.vendirConfig.Save()
}

func (h *HelmSource) configureHelmChartName() error {
	contents := h.vendirConfig.Contents()
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart name",
		Default:      contents[0].HelmChart.Name,
		ValidateFunc: nil,
	}
	name, err := h.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	contents[0].HelmChart.Name = strings.TrimSpace(name)
	h.vendirConfig.SetContents(contents)
	return h.vendirConfig.Save()
}

func (h *HelmSource) configureHelmChartVersion() error {
	contents := h.vendirConfig.Contents()
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart version",
		Default:      contents[0].HelmChart.Version,
		ValidateFunc: nil,
	}
	version, err := h.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	contents[0].HelmChart.Version = strings.TrimSpace(version)
	h.vendirConfig.SetContents(contents)
	return h.vendirConfig.Save()
}

func (h *HelmSource) configureHelmChartRepositoryURL() error {
	contents := h.vendirConfig.Contents()
	helmChartContent := contents[0].HelmChart
	if helmChartContent.Repository == nil {
		helmChartContent.Repository = &vendirconf.DirectoryContentsHelmChartRepo{}
	}
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart repository URL",
		Default:      helmChartContent.Repository.URL,
		ValidateFunc: nil,
	}
	url, err := h.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	helmChartContent.Repository.URL = strings.TrimSpace(url)
	contents[0].HelmChart = helmChartContent
	h.vendirConfig.SetContents(contents)
	return h.vendirConfig.Save()
}
