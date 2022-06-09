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
	//fetch.pkgAuthoringUI.PrintInformationalText("Now, we have to add the content which makes up the package for distribution. This content, which is essentially the configuration defined by the app operator, can be fetched from different types of sources. Imgpkg is a tool to package, distribute, and relocate Kubernetes configuration and dependent OCI images as one OCI artifact: a bundle.")
	return nil
}

func (fetch *FetchStep) Interact() error {
	fetch.pkgAuthoringUI.PrintHeading("\nPackage Content(Step 2/3)")
	isPreferenceImmutable := fetch.pkgBuild.Annotations[common.PkgCreatePreferenceAnnotationKey]

	defaultManifestOptionSelected := getManifestOptionFromPkgBuild(fetch.pkgBuild)
	fetch.pkgAuthoringUI.PrintInformationalText("In package, we need to fetch the manifest which defines how the application would be deployed in a K8s cluster. This manifest can be in the form of a yaml file used with `kubectl apply ...` or it could be a helm chart used with `helm install ...`. They can be available in any of the following locations. Please select from where to fetch the manifest")
	manifestOptions := []string{common.FetchReleaseArtifactFromGithub, common.FetchManifestFromGithub, common.FetchChartFromHelmRepo, common.FetchChartFromGithub}
	defaultFetchManifestOptionIndex := getDefaultManifestOptionIndex(manifestOptions, defaultManifestOptionSelected)
	choiceOpts := ui.ChoiceOpts{
		Label:   "From where to fetch the manifest",
		Default: defaultFetchManifestOptionIndex,
		Choices: manifestOptions,
	}
	manifestOptionSelectedIndex, err := fetch.pkgAuthoringUI.AskForChoice(choiceOpts)
	if err != nil {
		return err
	}
	//TODO Rohit should we move it up?
	if fetch.pkgBuild.Annotations == nil {
		fetch.pkgBuild.Annotations = make(map[string]string, 0)
	}
	fetch.pkgBuild.Annotations[common.PkgFetchContentAnnotationKey] = manifestOptions[manifestOptionSelectedIndex]
	fetch.pkgBuild.WriteToFile(fetch.pkgLocation)

	if isPreferenceImmutable == "true" {
		imgpkgStep := imgpkg.NewImgPkgStep(fetch.pkgAuthoringUI, fetch.pkgLocation, fetch.pkgBuild)
		err := common.Run(imgpkgStep)
		if err != nil {
			return err
		}
	} else {
		helmStep := NewHelmStep(fetch.pkgAuthoringUI, fetch.pkgLocation, fetch.pkgBuild)
		err := common.Run(helmStep)
		if err != nil {
			return err
		}
	}
	return nil
}

func getManifestOptionFromPkgBuild(pkgBuild *build.PackageBuild) string {
	return pkgBuild.Annotations[common.PkgFetchContentAnnotationKey]
}

func getDefaultManifestOptionIndex(manifestOptions []string, defaultManifestOptionSelected string) int {
	var defaultManifestOptionIndex int
	if defaultManifestOptionSelected == "" {
		defaultManifestOptionIndex = 0
	} else {
		for i, fetchTypeName := range manifestOptions {
			if fetchTypeName == defaultManifestOptionSelected {
				defaultManifestOptionIndex = i
				break
			}
		}
	}
	return defaultManifestOptionIndex
}

func setEarlierUpstreamOptionAsNil(fetchSection []v1alpha1.AppFetch, earlierFetchOption string) {
	switch earlierFetchOption {
	case AppFetchImgpkgBundle:
		fetchSection[0].ImgpkgBundle = nil
	case AppFetchHelmChart:
		fetchSection[0].HelmChart = nil
	case AppFetchGit:
		fetchSection[0].Git = nil
	case AppFetchHTTP:
		fetchSection[0].HTTP = nil
	case AppFetchImage:
		fetchSection[0].Image = nil
	case AppFetchInline:
		fetchSection[0].Inline = nil
	}
	return

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
