package imgpkg

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"strings"
)

type ImgpkgStep struct {
	pkgAuthoringUI pkgui.IAuthoringUI
	pkgLocation    string
	pkgBuild       *build.PackageBuild
}

func NewImgPkgStep(ui pkgui.IAuthoringUI, pkgLocation string, pkgBuild *build.PackageBuild) *ImgpkgStep {
	imgpkg := ImgpkgStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
	return &imgpkg
}

func (imgpkg ImgpkgStep) PreInteract() error {
	return nil
}

func (imgpkg *ImgpkgStep) Interact() error {
	var isImgpkgBundleCreated bool
	isImgpkgBundleCreated = false
	existingImgPkgBundleConf := imgpkg.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch[0].ImgpkgBundle
	if existingImgPkgBundleConf == nil {
		imgpkg.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch[0].ImgpkgBundle = &v1alpha1.AppFetchImgpkgBundle{}
	}

	if isImgpkgBundleCreated {
		textOpts := ui.TextOpts{
			Label:        "Enter the imgpkg bundle url",
			Default:      "",
			ValidateFunc: nil,
		}
		image, err := imgpkg.pkgAuthoringUI.AskForText(textOpts)
		image = strings.TrimSpace(image)
		imgpkg.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch[0].ImgpkgBundle.Image = image
		imgpkg.pkgBuild.WriteToFile()
		if err != nil {
			return err
		}
	} else {
		imgpkg.pkgAuthoringUI.PrintInformationalText("In package, we need to fetch the manifest which defines how the application would be deployed in a K8s cluster. This manifest can be in the form of a yaml file used with `kubectl apply ...` or it could be a helm chart used with `helm install ...`. They can be available in any of the following locations. Please select from where to fetch the manifest")
		options := []string{common.FetchReleaseArtifactFromGithub, common.FetchManifestFromGithub, common.FetchChartFromHelmRepo, common.FetchChartFromGithub, common.FetchFromLocalDirectory}
		defaultFetchOptionSelected := getPreviousSelectedFetchOption(imgpkg.pkgBuild)
		defaultFetchOptionIndex := getDefaultFetchOptionIndex(options, defaultFetchOptionSelected)
		choiceOpts := ui.ChoiceOpts{
			Label:   "From where to fetch the manifest",
			Default: defaultFetchOptionIndex,
			Choices: options,
		}
		manifestOptionSelectedIndex, err := imgpkg.pkgAuthoringUI.AskForChoice(choiceOpts)
		if err != nil {
			return err
		}
		if imgpkg.pkgBuild.Annotations == nil {
			imgpkg.pkgBuild.Annotations = make(map[string]string, 0)
		}
		imgpkg.pkgBuild.Annotations[common.PkgFetchContentAnnotationKey] = options[manifestOptionSelectedIndex]
		imgpkg.pkgBuild.WriteToFile()

		createImgPkgStep := NewCreateImgPkgStep(imgpkg.pkgAuthoringUI, imgpkg.pkgLocation, imgpkg.pkgBuild)
		err = common.Run(createImgPkgStep)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPreviousSelectedFetchOption(pkgBuild *build.PackageBuild) string {
	return pkgBuild.Annotations[common.PkgFetchContentAnnotationKey]
}

func getDefaultFetchOptionIndex(manifestOptions []string, defaultManifestOptionSelected string) int {
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

func (imgpkg ImgpkgStep) PostInteract() error {
	return nil
}
