package fetch

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/fetch/imgpkg"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type FetchStep struct {
	pkgAuthoringUI pkgui.IAuthoringUI
	pkgLocation    string
	pkgBuild       *build.PackageBuild
}

func NewFetchStep(pkgAuthoringUI pkgui.IAuthoringUI, pkgLocation string, pkgBuild *build.PackageBuild) *FetchStep {
	fetchStep := FetchStep{
		pkgAuthoringUI: pkgAuthoringUI,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
	return &fetchStep
}

func (fetch FetchStep) PreInteract() error {
	return nil
}

func (fetch *FetchStep) Interact() error {
	if fetch.pkgBuild.Spec.Pkg.Spec.Template.Spec == nil {
		fetch.pkgBuild.Spec.Pkg.Spec.Template.Spec = &v1alpha1.AppSpec{}
	}
	pkgFetchSection := fetch.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch
	fetch.pkgAuthoringUI.PrintHeaderText("\nPackage Content(Step 2/3)")

	if len(pkgFetchSection) == 0 {
		fetch.initializeFetchWithImgpkgBundle()
	} else if len(pkgFetchSection) > 1 {
		// As multiple fetch section are configured, we dont want to touch them.
		return nil
	} else if len(fetch.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch) == 1 && pkgFetchSection[0].ImgpkgBundle == nil {
		// Return as we support only creation of imgpkg bundle
		return nil
	}

	//Package Authoring will always create the imgpkg bundle and put that in the fetch section.
	imgpkgStep := imgpkg.NewImgPkgStep(fetch.pkgAuthoringUI, fetch.pkgLocation, fetch.pkgBuild)
	err := common.Run(imgpkgStep)
	if err != nil {
		return err
	}
	return nil
}

func (fetch FetchStep) PostInteract() error {
	return nil
}

func (fetch FetchStep) initializeFetchWithImgpkgBundle() {
	fetch.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch = append(fetch.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch, v1alpha1.AppFetch{ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{}})
}
