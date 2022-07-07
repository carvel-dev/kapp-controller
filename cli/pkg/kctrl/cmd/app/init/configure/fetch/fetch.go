// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

type FetchStep struct {
	ui                        cmdcore.AuthoringUI
	appBuild                  *appbuild.AppBuild
	isAppCommandRunExplicitly bool
	hasFetchOptionChanged     bool
}

func NewFetchStep(ui cmdcore.AuthoringUI, appBuild *appbuild.AppBuild, isAppCommandRunExplicitly bool) *FetchStep {
	fetchStep := FetchStep{
		ui:                        ui,
		appBuild:                  appBuild,
		isAppCommandRunExplicitly: isAppCommandRunExplicitly,
	}
	return &fetchStep
}

func (fetchStep *FetchStep) PreInteract() error {
	return nil
}

func (fetchStep *FetchStep) Interact() error {
	if fetchStep.isAppCommandRunExplicitly {
		fetchStep.ui.PrintHeaderText("\nApp Content(Step 2/3)")
		fetchStep.ui.PrintInformationalText("We need to fetch the manifest which defines how the application would be deployed in a K8s cluster. This manifest can be in the form of a yaml file used with `kubectl apply ...` or it could be a helm chart used with `helm install ...`. They can be available in any of the following locations. Please select from where to fetch the manifest")
	} else {
		fetchStep.ui.PrintHeaderText("Package Content(Step 2/3)")
		fetchStep.ui.PrintInformationalText("We need to fetch the manifest which defines how the package would be deployed in a K8s cluster. This manifest can be in the form of a yaml file used with `kubectl apply ...` or it could be a helm chart used with `helm install ...`. They can be available in any of the following locations. Please select from where to fetch the manifest")
	}

	var vendirConfig vendirconf.Config
	vendirConfig, err := ReadVendirConfig()
	if err != nil {
		return err
	}
	isHelmTemplateExistInPreviousOption := fetchStep.helmTemplateExistInAppBuild()
	previousFetchOptionSelected := getPreviousFetchOptionFromVendir(vendirConfig, isHelmTemplateExistInPreviousOption)

	options := []string{FetchReleaseArtifactFromGithub, FetchManifestFromGithub, FetchChartFromHelmRepo, FetchChartFromGithub, FetchFromLocalDirectory}
	previousFetchOptionIndex := getPreviousFetchOptionIndex(options, previousFetchOptionSelected)
	defaultFetchOptionIndex := previousFetchOptionIndex
	choiceOpts := ui.ChoiceOpts{
		Label:   "From where to fetch the manifest",
		Default: defaultFetchOptionIndex,
		Choices: options,
	}
	currentFetchOptionIndex, err := fetchStep.ui.AskForChoice(choiceOpts)
	if err != nil {
		return err
	}
	if fetchStep.appBuild.ObjectMeta.Annotations == nil {
		fetchStep.appBuild.ObjectMeta.Annotations = make(map[string]string)
	}
	currentFetchOptionSelected := options[currentFetchOptionIndex]
	fetchStep.appBuild.ObjectMeta.Annotations[FetchContentAnnotationKey] = currentFetchOptionSelected
	if fetchStep.isAppCommandRunExplicitly {
		fetchStep.appBuild.Save()
	}
	//For a local directory, we will be including everything.
	if currentFetchOptionSelected == FetchFromLocalDirectory {
		fetchStep.ui.PrintInformationalText("For local directory, we are going to include everything as part of `init` command.")
		return nil
	}

	// TODO handle a scenario where previousFetchOptionSelected is Helm from Chart(or anything similar) to currentFetchOptionSelected is Local Dir.
	// Need to remove vendir.yml in this case.
	// One more edge case: user move from github release to helm. In that case, we should remove the whole template section I think.
	if currentFetchOptionSelected != previousFetchOptionSelected {
		fetchStep.hasFetchOptionChanged = true
		vendirConfig = NewDefaultVendirConfig()
	}

	vendirStep := NewVendirStep(fetchStep.ui, vendirConfig, currentFetchOptionSelected)
	return common.Run(vendirStep)
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

func (fetchStep *FetchStep) PostInteract() error {
	return nil
}

func (fetchStep *FetchStep) helmTemplateExistInAppBuild() bool {
	if fetchStep.appBuild.Spec.App == nil || fetchStep.appBuild.Spec.App.Spec == nil || fetchStep.appBuild.Spec.App.Spec.Template == nil {
		return false
	}
	appTemplates := fetchStep.appBuild.Spec.App.Spec.Template
	for _, appTemplate := range appTemplates {
		if appTemplate.HelmTemplate != nil {
			return true
		}
	}
	return false
}

func (fetchStep *FetchStep) HasFetchOptionChanged() bool {
	return fetchStep.hasFetchOptionChanged
}
