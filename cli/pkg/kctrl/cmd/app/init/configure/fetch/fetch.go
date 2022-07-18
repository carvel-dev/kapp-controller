// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"fmt"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/interfaces/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/interfaces/step"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

type FetchStep struct {
	ui                        cmdcore.AuthoringUI
	build                     build.Build
	isAppCommandRunExplicitly bool
}

func NewFetchStep(ui cmdcore.AuthoringUI, build build.Build, isAppCommandRunExplicitly bool) *FetchStep {
	fetchStep := FetchStep{
		ui:                        ui,
		build:                     build,
		isAppCommandRunExplicitly: isAppCommandRunExplicitly,
	}
	return &fetchStep
}

func (fetchStep *FetchStep) PreInteract() error {
	return nil
}

func (fetchStep *FetchStep) Interact() error {
	fetchStep.ui.PrintHeaderText("Content (Step 2/3)")
	fetchStep.ui.PrintInformationalText("We need to fetch the manifest which defines how the package would be deployed in a K8s cluster. This manifest can be in the form of a yaml file used with `kubectl apply ...` or it could be a helm chart used with `helm install ...`. They can be available in any of the following locations. Please select one among them")

	var vendirConfig vendirconf.Config
	vendirConfig, err := ReadVendirConfig()
	if err != nil {
		return err
	}
	isHelmTemplateExistInPreviousOption := fetchStep.helmTemplateExistInAppBuild()
	previousFetchOptionSelected := GetFetchOptionFromVendir(vendirConfig, isHelmTemplateExistInPreviousOption)

	options := []string{FetchFromLocalDirectory, FetchFromGithubRelease, FetchChartFromHelmRepo, FetchManifestFromGit, FetchChartFromGit}
	previousFetchOptionIndex := getPreviousFetchOptionIndex(options, previousFetchOptionSelected)
	defaultFetchOptionIndex := previousFetchOptionIndex
	choiceOpts := ui.ChoiceOpts{
		Label:   "Enter configuration source",
		Default: defaultFetchOptionIndex,
		Choices: options,
	}
	currentFetchOptionIndex, err := fetchStep.ui.AskForChoice(choiceOpts)
	if err != nil {
		return err
	}
	currentFetchOptionSelected := options[currentFetchOptionIndex]

	if previousFetchOptionSelected != "" && currentFetchOptionSelected != previousFetchOptionSelected {
		return fmt.Errorf("Transitioning from one fetch option to another is not allowed. Earlier option selected: %s, Current Option selected: %s", previousFetchOptionSelected, currentFetchOptionSelected)
	}

	buildObjectMeta := fetchStep.build.GetObjectMeta()
	if buildObjectMeta.Annotations == nil {
		buildObjectMeta.Annotations = make(map[string]string)
	}
	buildObjectMeta.Annotations[FetchContentAnnotationKey] = currentFetchOptionSelected
	fetchStep.build.SetObjectMeta(buildObjectMeta)
	// For a local directory, we will be including everything.
	if currentFetchOptionSelected == FetchFromLocalDirectory {
		return nil
	}

	vendirStep := NewVendirStep(fetchStep.ui, vendirConfig, currentFetchOptionSelected)
	return step.Run(vendirStep)
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

func (fetchStep *FetchStep) helmTemplateExistInAppBuild() bool {
	appSpec := fetchStep.build.GetAppSpec()
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

func (fetchStep *FetchStep) PostInteract() error {
	return nil
}
