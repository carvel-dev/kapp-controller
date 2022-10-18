// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sources

import (
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	buildconfigs "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local/buildconfigs"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

const (
	VendirSyncDirectory = "upstream"
	vendirFileName      = "vendir.yml"
	includeAllFiles     = "*"
)

const (
	FetchFromGithubRelease       string = "Github Release"
	FetchFromGit                 string = "Git Repository"
	FetchFromHelmRepo            string = "Helm Chart from Helm Repository"
	FetchChartFromGit            string = "Helm Chart from Git Repository"
	FetchFromLocalDirectory      string = "Local Directory"
	MultipleFetchOptionsSelected string = "MultipleFetchOptionsSelected"
)

type FetchConfiguration struct {
	ui    cmdcore.AuthoringUI
	build buildconfigs.Build
}

func NewFetchConfiguration(ui cmdcore.AuthoringUI, build buildconfigs.Build) FetchConfiguration {
	return FetchConfiguration{ui: ui, build: build}
}

func (f FetchConfiguration) Configure() (string, SourceConfiguration, error) {
	f.ui.PrintHeaderText("Content")
	f.ui.PrintInformationalText("Please provide the location from where your Kubernetes manifests or Helm chart can be fetched. This will be bundled as a part of the package.")

	vendirConfig := NewVendirConfig(vendirFileName)
	err := vendirConfig.Load()
	if err != nil {
		return "", nil, err
	}

	isHelmTemplateExistInPreviousOption := f.build.HasHelmTemplate()
	previousFetchOptionSelected := vendirConfig.FetchMode(isHelmTemplateExistInPreviousOption)
	if previousFetchOptionSelected == MultipleFetchOptionsSelected {
		// As this is advanced use case, we dont know how to handle it.
		f.ui.PrintInformationalText("Since vendir is syncing data from multiple resources, we will not reconfigure vendir.yml and run vendir sync.")
		return "", nil, nil
	}

	options := []string{FetchFromLocalDirectory, FetchFromGithubRelease, FetchFromHelmRepo, FetchFromGit, FetchChartFromGit}
	previousFetchOptionIndex := f.getPreviousFetchOptionIndex(options, previousFetchOptionSelected)
	defaultFetchOptionIndex := previousFetchOptionIndex
	choiceOpts := ui.ChoiceOpts{
		Label:   "Enter source",
		Default: defaultFetchOptionIndex,
		Choices: options,
	}
	currentFetchOptionIndex, err := f.ui.AskForChoice(choiceOpts)
	if err != nil {
		return "", nil, err
	}
	currentFetchOptionSelected := options[currentFetchOptionIndex]

	if previousFetchOptionSelected != "" && currentFetchOptionSelected != previousFetchOptionSelected {
		return "", nil, fmt.Errorf("Transitioning from one fetch option to another is not allowed. Earlier option selected: %s, Current Option selected: %s", previousFetchOptionSelected, currentFetchOptionSelected)
	}

	// For the local directory options, all files/directories in working directory are used while releasing
	if currentFetchOptionSelected == FetchFromLocalDirectory {
		return currentFetchOptionSelected, nil, nil
	}

	vendirDirectories := vendirConfig.Directories()
	if len(vendirDirectories) > 1 {
		return currentFetchOptionSelected, nil, fmt.Errorf("More than 1 directory config found in the vendir file. (hint: Run vendir sync manually)")
	}
	if len(vendirDirectories) == 0 {
		err := f.initializeVendirDirectorySection(vendirConfig)
		if err != nil {
			return currentFetchOptionSelected, nil, err
		}
	} else {
		directory := vendirDirectories[0]
		if len(directory.Contents) > 1 {
			return currentFetchOptionSelected, nil, fmt.Errorf("More than 1 content config found in the vendir file. (hint: Run vendir sync manually)")
		}
	}

	switch currentFetchOptionSelected {
	case FetchFromGithubRelease:
		return currentFetchOptionSelected, NewGithubReleaseConfiguration(f.ui, vendirConfig), nil
	case FetchFromHelmRepo:
		return currentFetchOptionSelected, NewHelmConfiguration(f.ui, vendirConfig), nil
	case FetchFromGit, FetchChartFromGit:
		return currentFetchOptionSelected, NewGitConfiguration(f.ui, vendirConfig), nil
	}
	return currentFetchOptionSelected, nil, fmt.Errorf("Unexppected: Invalid fetch mode encountered while configuring vendir")
}

func (f FetchConfiguration) getPreviousFetchOptionIndex(manifestOptions []string, previousFetchOption string) int {
	var previousFetchOptionIndex int
	if previousFetchOption == "" {
		previousFetchOptionIndex = 0
	} else {
		for i, fetchTypeName := range manifestOptions {
			if fetchTypeName == previousFetchOption {
				previousFetchOptionIndex = i
				break
			}
		}
	}
	return previousFetchOptionIndex
}

func (f *FetchConfiguration) initializeVendirDirectorySection(vendirConfig *VendirConfig) error {
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
	vendirConfig.SetDirectories(directories)
	err := vendirConfig.Save()
	if err != nil {
		return err
	}
	return nil
}
