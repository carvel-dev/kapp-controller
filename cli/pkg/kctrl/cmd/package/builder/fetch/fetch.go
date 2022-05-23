package fetch

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/fetch/imgpkg"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

const (
	AppFetchInline       string = "Inline"
	AppFetchImage        string = "Image"
	AppFetchHTTP         string = "HTTP"
	AppFetchGit          string = "Git"
	AppFetchHelmChart    string = "HelmChart"
	AppFetchImgpkgBundle string = "Imgpkg(recommended)"
)

type FetchStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	pkgLocation    string
	pkgBuild       *build.PackageBuild
}

func NewFetchStep(ui pkgui.IPkgAuthoringUI, pkgLocation string, pkgBuild *build.PackageBuild) *FetchStep {
	fetchStep := FetchStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
	return &fetchStep
}

func (fetch FetchStep) PreInteract() error {
	fetch.pkgAuthoringUI.PrintInformationalText(`Now, we have to add the configuration which makes up the package for distribution. 
Configuration can be fetched from different types of sources.
Imgpkg is a tool to package, distribute, and relocate Kubernetes configuration and dependent OCI images as one OCI artifact: a bundle.`)
	return nil
}

func (fetch *FetchStep) Interact() error {
	fetchSection := fetch.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch
	var defaultFetchOptionSelected string
	if len(fetchSection) > 1 {
		//As multiple fetch sections are configured, we dont want to touch them.
		return nil
	} else if len(fetchSection) == 0 {
		//Initialize fetch Section
		var appFetchList []v1alpha1.AppFetch
		appFetchList = append(appFetchList, v1alpha1.AppFetch{})
		fetch.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch = appFetchList
	} else {
		defaultFetchOptionSelected = getFetchOptionFromPkgBuild(fetch.pkgBuild)
	}

	var fetchTypeNames = []string{AppFetchImgpkgBundle, AppFetchHelmChart}
	defaultFetchOptionIndex := getDefaultFetchOptionIndex(fetchTypeNames, defaultFetchOptionSelected)
	choiceOpts := ui.ChoiceOpts{
		Label:   "Enter the fetch configuration type",
		Default: defaultFetchOptionIndex,
		Choices: fetchTypeNames,
	}
	fetchOptionSelected, err := fetch.pkgAuthoringUI.AskForChoice(choiceOpts)
	if err != nil {
		return err
	}

	switch fetchTypeNames[fetchOptionSelected] {
	case AppFetchImgpkgBundle:
		imgpkgStep := imgpkg.NewImgPkgStep(fetch.pkgAuthoringUI, fetch.pkgLocation, fetch.pkgBuild)
		err := common.Run(imgpkgStep)
		if err != nil {
			return err
		}
	}
	return nil
}

func getFetchOptionFromPkgBuild(pkgBuild *build.PackageBuild) string {
	appFetch := pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch[0]
	var selectedAppFetch string
	switch {
	case appFetch.Inline != nil:
		selectedAppFetch = AppFetchInline
	case appFetch.Image != nil:
		selectedAppFetch = AppFetchImage
	case appFetch.ImgpkgBundle != nil:
		selectedAppFetch = AppFetchImgpkgBundle
	case appFetch.HTTP != nil:
		selectedAppFetch = AppFetchHTTP
	case appFetch.Git != nil:
		selectedAppFetch = AppFetchGit
	case appFetch.HelmChart != nil:
		selectedAppFetch = AppFetchHelmChart
	default:
		selectedAppFetch = ""
	}
	return selectedAppFetch
}

func getDefaultFetchOptionIndex(fetchTypeNames []string, defaultFetchOptionSelected string) int {
	var defaultFetchOptionIndex int
	if defaultFetchOptionSelected == "" {
		defaultFetchOptionIndex = 0
	} else {
		for i, fetchTypeName := range fetchTypeNames {
			if fetchTypeName == defaultFetchOptionSelected {
				defaultFetchOptionIndex = i
				break
			}
		}
	}
	return defaultFetchOptionIndex
}

func (fetch FetchStep) PostInteract() error {
	return nil
}
