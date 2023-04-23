// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sources

import (
	"fmt"
	"log"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	buildconfigs "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local/buildconfigs"
)

const (
	VendirSyncDirectory = "upstream"
	vendirFileName      = "vendir.yml"
	includeAllFiles     = "*"
)

const (
	GithubRelease                string = "Github Release"
	Git                          string = "Git Repository"
	HelmRepo                     string = "Helm Chart from Helm Repository"
	ChartFromGit                 string = "Helm Chart from Git Repository"
	LocalDirectory               string = "Local Directory"
	MultipleFetchOptionsSelected string = "MultipleFetchOptionsSelected"
)

type Source struct {
	ui    cmdcore.AuthoringUI
	build buildconfigs.Build
}

func NewSource(ui cmdcore.AuthoringUI, build buildconfigs.Build) Source {
	return Source{ui: ui, build: build}
}

func (f Source) Configure() (string, error) {
	f.ui.PrintHeaderText("Content")
	f.ui.PrintInformationalText("Please provide the location from where your Kubernetes manifests or Helm chart can be fetched. This will be bundled as a part of the package.")

	vendirConfig := NewVendirConfig(vendirFileName)
	err := vendirConfig.Load()
	if err != nil {
		return "", err
	}

	previousFetchOptionSelected := vendirConfig.FetchMode(f.build.HasHelmTemplate())
	if previousFetchOptionSelected == MultipleFetchOptionsSelected {
		f.ui.PrintInformationalText("vendir.yml has multiple sources defined. Running vendir sync without overwriting vendir.yml")
		return "", nil
	}

	options := []string{LocalDirectory, GithubRelease, HelmRepo, Git, ChartFromGit}
	log.Println("=================previousFetchOptionSelected:", previousFetchOptionSelected, f.getPreviousFetchOptionIndex(options, previousFetchOptionSelected))

	choiceOpts := ui.ChoiceOpts{
		Label:   "Enter source",
		Default: f.getPreviousFetchOptionIndex(options, previousFetchOptionSelected),
		Choices: options,
	}
	haveChange := false
	currentFetchOptionIndex := f.getPreviousFetchOptionIndex(options, previousFetchOptionSelected)
	currentFetchOptionSelected := options[currentFetchOptionIndex]
	if previousFetchOptionSelected == "" {
		currentFetchOptionIndex, err = f.ui.AskForChoice(choiceOpts)
		if err != nil {
			return "", err
		}
		log.Println("=================currentFetchOptionIndex:", currentFetchOptionIndex)
		currentFetchOptionSelected = options[currentFetchOptionIndex]
		haveChange = true
	}

	// if previousFetchOptionSelected != "" && currentFetchOptionSelected != previousFetchOptionSelected {
	// 	return "", fmt.Errorf("Transitioning from one fetch option to another is not allowed. Earlier option selected: %s, Current Option selected: %s", previousFetchOptionSelected, currentFetchOptionSelected)
	// }

	// For the local directory options, all files/directories in working directory are used while releasing
	if currentFetchOptionSelected == LocalDirectory {
		return currentFetchOptionSelected, nil
	}

	vendirDirectories := vendirConfig.Directories()
	switch {
	case len(vendirDirectories) > 1:
		return currentFetchOptionSelected, fmt.Errorf("More than 1 directory config found in the vendir file. (hint: Run vendir sync manually)")
	case len(vendirDirectories) == 0:
		err := vendirConfig.InitiatliseDirectories()
		if err != nil {
			return currentFetchOptionSelected, err
		}
	default:
		directory := vendirDirectories[0]
		if len(directory.Contents) > 1 {
			return currentFetchOptionSelected, fmt.Errorf("More than 1 content config found in the vendir file. (hint: Run vendir sync manually)")
		}
	}

	switch currentFetchOptionSelected {
	case GithubRelease:
		return currentFetchOptionSelected, NewGithubReleaseSource(f.ui, vendirConfig).Configure()
	case HelmRepo:
		return currentFetchOptionSelected, NewHelmSource(f.ui, vendirConfig).Configure()
	case Git, ChartFromGit:
		return currentFetchOptionSelected, NewGitSource(f.ui, vendirConfig).Configure()
	}
	return currentFetchOptionSelected, fmt.Errorf("Unexppected: Invalid fetch mode encountered while configuring vendir")
}

func (f Source) getPreviousFetchOptionIndex(manifestOptions []string, previousFetchOption string) int {
	for i, fetchTypeName := range manifestOptions {
		if fetchTypeName == previousFetchOption {
			return i
		}
	}
	return 0
}
