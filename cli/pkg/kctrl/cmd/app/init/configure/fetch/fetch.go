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
	fetchStep.ui.PrintHeaderText("\nApp Content(Step 2/3)")

	fetchStep.ui.PrintInformationalText("We need to fetch the manifest which defines how the application would be deployed in a K8s cluster. This manifest can be in the form of a yaml file used with `kubectl apply ...` or it could be a helm chart used with `helm install ...`. They can be available in any of the following locations. Please select from where to fetch the manifest")
	options := []string{common.FetchReleaseArtifactFromGithub, common.FetchManifestFromGithub, common.FetchChartFromHelmRepo, common.FetchChartFromGithub, common.FetchFromLocalDirectory}
	previousFetchOption := getPreviousFetchOption(fetchStep.appBuild)
	previousFetchOptionIndex := getPreviousFetchOptionIndex(options, previousFetchOption)
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
	fetchStep.appBuild.ObjectMeta.Annotations[common.FetchContentAnnotationKey] = currentFetchOptionSelected
	if fetchStep.isAppCommandRunExplicitly {
		fetchStep.appBuild.WriteToFile()
	}
	//For a local directory, we will be including everything.
	if currentFetchOptionSelected == common.FetchFromLocalDirectory {
		fetchStep.ui.PrintInformationalText("For local directory, we are going to include everything as part of `init` command.")
		return nil
	}
	var vendirConfig vendirconf.Config
	if defaultFetchOptionIndex != currentFetchOptionIndex {
		//resetVendirConf
		vendirConfig = NewDefaultVendirConfig()
	} else {
		vendirConfig, err = ReadVendirConfig()
		if err != nil {
			return err
		}
	}

	vendirStep := NewVendirStep(fetchStep.ui, fetchStep.appBuild, vendirConfig)
	err = common.Run(vendirStep)
	if err != nil {
		return err
	}
	return nil
}

func getPreviousFetchOption(appBuild *appbuild.AppBuild) string {
	return appBuild.ObjectMeta.Annotations[common.FetchContentAnnotationKey]
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

func (fetch *FetchStep) PostInteract() error {
	return nil
}
