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
	vendirConfig *VendirConfig
}

func NewHelmStep(ui cmdcore.AuthoringUI, vendirConfig *VendirConfig) *HelmStep {
	return &HelmStep{
		ui:           ui,
		vendirConfig: vendirConfig,
	}
}

func (h *HelmStep) PreInteract() error { return nil }

func (h *HelmStep) Interact() error {
	contents := h.vendirConfig.Contents()
	if contents == nil {
		err := h.initializeContentWithHelmRelease(contents)
		if err != nil {
			return err
		}
	} else if contents[0].HelmChart == nil {
		err := h.initializeHelmRelease(contents)
		if err != nil {
			return err
		}
	}

	err := h.configureHelmChartRepositoryURL(contents)
	if err != nil {
		return err
	}
	err = h.configureHelmChartName(contents)
	if err != nil {
		return err
	}
	return h.configureHelmChartVersion(contents)
}

func (h *HelmStep) initializeHelmRelease(contents []vendirconf.DirectoryContents) error {
	contents[0].HelmChart = &vendirconf.DirectoryContentsHelmChart{}
	h.vendirConfig.SetContents(contents)
	return h.vendirConfig.Save()
}

func (h *HelmStep) initializeContentWithHelmRelease(contents []vendirconf.DirectoryContents) error {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	h.vendirConfig.SetContents(append(h.vendirConfig.Contents(), vendirconf.DirectoryContents{}))
	return h.initializeHelmRelease(contents)
}

func (h *HelmStep) configureHelmChartName(contents []vendirconf.DirectoryContents) error {
	defaultName := contents[0].HelmChart.Name
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart name",
		Default:      defaultName,
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

func (h *HelmStep) configureHelmChartVersion(contents []vendirconf.DirectoryContents) error {
	defaultVersion := contents[0].HelmChart.Version
	textOpts := ui.TextOpts{
		Label:        "Enter helm chart version",
		Default:      defaultVersion,
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

func (h *HelmStep) configureHelmChartRepositoryURL(contents []vendirconf.DirectoryContents) error {
	helmChartContent := contents[0].HelmChart
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
	contents[0].HelmChart = helmChartContent
	h.vendirConfig.SetContents(contents)
	return h.vendirConfig.Save()
}

func (h *HelmStep) PostInteract() error { return nil }
