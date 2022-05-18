package fetch

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/fetch/imgpkg"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type FetchStep struct {
	ui          ui.UI
	pkgLocation string
	pkgBuild    *build.PackageBuild
}

func NewFetchStep(ui ui.UI, pkgLocation string, pkgBuild *build.PackageBuild) *FetchStep {
	fetchStep := FetchStep{
		ui:          ui,
		pkgLocation: pkgLocation,
		pkgBuild:    pkgBuild,
	}
	return &fetchStep
}

func (fetch FetchStep) PreInteract() error {
	str := `Now, we have to add the configuration which makes up the package for distribution. 
Configuration can be fetched from different types of sources.
Imgpkg is a tool to package, distribute, and relocate Kubernetes configuration and dependent OCI images as one OCI artifact: a bundle.`
	fetch.ui.BeginLinef(str)
	return nil
}

func (fetch *FetchStep) Interact() error {
	fetchSection := fetch.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch
	if len(fetchSection) == 0 {
		fetch.initializeFetchSection()
	}
	var fetchOptionSelected int
	var fetchTypeNames = []string{"Imgpkg(recommended)", "HelmChart", "Inline"}
	//TODO while reading we have to dissect fetch section and see what was the configuration used and make that as default.
	fetchOptionSelected, err := fetch.ui.AskForChoice("Enter the fetch configuration type", fetchTypeNames)
	if err != nil {
		return err
	}
	switch fetchTypeNames[fetchOptionSelected] {
	case "Imgpkg(recommended)":
		imgpkgStep := imgpkg.NewImgPkgStep(fetch.ui, fetch.pkgLocation, fetch.pkgBuild)
		err := common.Run(imgpkgStep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fetch FetchStep) PostInteract() error {
	return nil
}

func (fetch FetchStep) initializeFetchSection() {
	var appFetchList []v1alpha1.AppFetch
	appFetchList = append(appFetchList, v1alpha1.AppFetch{})
	fetch.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch = appFetchList
}
