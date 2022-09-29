// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"

	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

const (
	VendirFileName      = "vendir.yml"
	VendirSyncDirectory = "upstream"
	IncludeAllFiles     = "*"
)

type VendirConfigBuilder struct {
	ui          cmdcore.AuthoringUI
	config      *VendirConfig
	fetchOption string
}

func NewVendirConfigBuilder(ui cmdcore.AuthoringUI, config *VendirConfig, fetchOption string) VendirConfigBuilder {
	return VendirConfigBuilder{ui: ui, config: config, fetchOption: fetchOption}
}

func (v VendirConfigBuilder) Configure() error {
	vendirDirectories := v.config.Directories()
	if len(vendirDirectories) > 1 {
		return fmt.Errorf("More than 1 directory config found in the vendir file. (hint: Run vendir sync manually)")

	}
	if len(vendirDirectories) == 0 {
		err := v.initializeVendirDirectorySection()
		if err != nil {
			return err
		}
	} else {
		directory := vendirDirectories[0]
		if len(directory.Contents) > 1 {
			return fmt.Errorf("More than 1 content config found in the vendir file. (hint: Run vendir sync manually)")
		}
	}
	currentFetchOptionSelected := v.fetchOption
	switch currentFetchOptionSelected {
	case FetchFromGithubRelease:
		return NewGithubReleaseConfiguration(v.ui, v.config).Configure()
	case FetchFromHelmRepo:
		return NewHelmConfiguration(v.ui, v.config).Configure()
	case FetchFromGit, FetchChartFromGit:
		return NewGitConfiguration(v.ui, v.config).Configure()
	}
	return fmt.Errorf("Unexppected: Invalid fetch mode encountered while configuring vendir")
}

func (v *VendirConfigBuilder) initializeVendirDirectorySection() error {
	var directory vendirconf.Directory
	directory = vendirconf.Directory{
		Path: VendirSyncDirectory,
		Contents: []vendirconf.DirectoryContents{
			{
				Path: ".",
			},
		},
	}
	directories := []vendirconf.Directory{directory}
	v.config.SetDirectories(directories)
	err := v.config.Save()
	if err != nil {
		return err
	}
	return nil
}
