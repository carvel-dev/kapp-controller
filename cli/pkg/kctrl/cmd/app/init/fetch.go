// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
)

const (
	FetchContentAnnotationKey = "fetch-content-from"
	LocalFetchAnnotationKey   = "kctrl.carvel.dev/local-fetch-0"
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
	build Build
}

func NewFetchConfiguration(ui cmdcore.AuthoringUI, build Build) FetchConfiguration {
	return FetchConfiguration{ui: ui, build: build}
}

func (f FetchConfiguration) Configure() error {
	f.ui.PrintHeaderText("Content")
	f.ui.PrintInformationalText("Please provide the location from where your Kubernetes manifests or Helm chart can be fetched. This will be bundled as a part of the package.")

	vendirConfig := NewVendirConfig(VendirFileName)
	err := vendirConfig.Load()
	if err != nil {
		return err
	}

	isHelmTemplateExistInPreviousOption := f.build.HasHelmTemplate()
	previousFetchOptionSelected := vendirConfig.FetchMode(isHelmTemplateExistInPreviousOption)
	if previousFetchOptionSelected == MultipleFetchOptionsSelected {
		// As this is advanced use case, we dont know how to handle it.
		f.ui.PrintInformationalText("Since vendir is syncing data from multiple resources, we will not reconfigure vendir.yml and run vendir sync.")
		return nil
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
		return err
	}
	currentFetchOptionSelected := options[currentFetchOptionIndex]

	if previousFetchOptionSelected != "" && currentFetchOptionSelected != previousFetchOptionSelected {
		return fmt.Errorf("Transitioning from one fetch option to another is not allowed. Earlier option selected: %s, Current Option selected: %s", previousFetchOptionSelected, currentFetchOptionSelected)
	}

	buildObjectMeta := f.build.GetObjectMeta()
	if buildObjectMeta.Annotations == nil {
		buildObjectMeta.Annotations = make(map[string]string)
	}
	buildObjectMeta.Annotations[FetchContentAnnotationKey] = currentFetchOptionSelected
	f.build.SetObjectMeta(buildObjectMeta)
	// For the local directory options, all files/directories in working directory are used while releasing
	if currentFetchOptionSelected == FetchFromLocalDirectory {
		return nil
	}

	return NewVendirConfigBuilder(f.ui, vendirConfig, currentFetchOptionSelected).Configure()
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
