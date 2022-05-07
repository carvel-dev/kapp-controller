package fetch

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/fetch/imgpkg"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

const (
	Imgpkg int = iota
	HelmChart
	Inline
)

var fetchTypeNames = []string{"Imgpkg(recommended)", "HelmChart", "Inline"}

type FetchStep struct {
	ui          ui.UI
	pkgLocation string
	AppFetch    []v1alpha1.AppFetch
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
	var appFetchList []v1alpha1.AppFetch
	var fetchOptionSelected int
	fetchOptionSelected, err := fetch.ui.AskForChoice("Enter the fetch configuration type", fetchTypeNames)
	if err != nil {
		return err
	}
	//TODO Rohit This is error prone. How can we make it better. Option: Switch over the list items e.g. fetchTypeNames[fetchOptionSelected]
	switch fetchOptionSelected {
	case Imgpkg:
		imgpkgStep := imgpkg.NewImgPkgStep(fetch.ui, fetch.pkgLocation, fetch.pkgBuild)
		err := imgpkgStep.Run()
		if err != nil {
			return err
		}
		appFetchList = append(appFetchList, v1alpha1.AppFetch{
			ImgpkgBundle: &imgpkgStep.ImgpkgBundle})
	}

	fetch.AppFetch = appFetchList
	return nil
}

func (fetch FetchStep) PostInteract() error {
	return nil
}

func (fetch *FetchStep) Run() error {
	err := fetch.PreInteract()
	if err != nil {
		return err
	}
	err = fetch.Interact()
	if err != nil {
		return err
	}
	err = fetch.PostInteract()
	if err != nil {
		return err
	}
	return nil
}
