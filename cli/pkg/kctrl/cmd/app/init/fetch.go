// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
)

type FetchStep struct {
	ui                        cmdcore.AuthoringUI
	build                     Build
	isAppCommandRunExplicitly bool
}

func NewFetchStep(ui cmdcore.AuthoringUI, build Build, isAppCommandRunExplicitly bool) *FetchStep {
	return &FetchStep{
		ui:                        ui,
		build:                     build,
		isAppCommandRunExplicitly: isAppCommandRunExplicitly,
	}
}

func (f *FetchStep) PreInteract() error { return nil }

func (f *FetchStep) Interact() error {
	f.ui.PrintHeaderText("Content (Step 2/3)")
	f.ui.PrintInformationalText("We need to fetch the manifest which defines how the package would be deployed in a K8s cluster. This manifest can be in the form of a yaml file used with `kubectl apply ...` or it could be a helm chart used with `helm install ...`. They can be available in any of the following locations. Please select one among them")

	vendirConfig, err := ReadVendirConfig()
	if err != nil {
		return err
	}
	isHelmTemplateExistInPreviousOption := f.helmTemplateExistInAppBuild()
	previousFetchOptionSelected := GetFetchOptionFromVendir(vendirConfig, isHelmTemplateExistInPreviousOption)

	options := []string{FetchFromLocalDirectory, FetchFromGithubRelease, FetchFromHelmRepo, FetchFromGit, FetchChartFromGit}
	previousFetchOptionIndex := getPreviousFetchOptionIndex(options, previousFetchOptionSelected)
	defaultFetchOptionIndex := previousFetchOptionIndex
	choiceOpts := ui.ChoiceOpts{
		Label:   "Enter configuration source",
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
	// For a local directory, we will be including everything.
	if currentFetchOptionSelected == FetchFromLocalDirectory {
		return nil
	}

	vendirStep := NewVendirStep(f.ui, vendirConfig, currentFetchOptionSelected)
	return Run(vendirStep)
}

func getPreviousFetchOptionIndex(manifestOptions []string, previousFetchOption string) int {
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

func (f *FetchStep) helmTemplateExistInAppBuild() bool {
	appSpec := f.build.GetAppSpec()
	if appSpec == nil || appSpec.Template == nil {
		return false
	}
	appTemplates := appSpec.Template
	for _, appTemplate := range appTemplates {
		if appTemplate.HelmTemplate != nil {
			return true
		}
	}
	return false
}

func (f *FetchStep) PostInteract() error { return nil }
