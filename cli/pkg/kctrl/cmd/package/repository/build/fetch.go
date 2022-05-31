package build

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/repository/build/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
)

const (
	Git          string = "Git"
	ImgpkgBundle string = "Imgpkg(recommended)"
)

type FetchStep struct {
	pkgAuthoringUI  pkgui.IPkgAuthoringUI
	pkgRepoLocation string
	pkgRepoBuild    *build.PackageRepositoryBuild
}

func NewFetchStep(ui pkgui.IPkgAuthoringUI, pkgLocation string, pkgRepoBuild *build.PackageRepositoryBuild) *FetchStep {
	fetchStep := FetchStep{
		pkgAuthoringUI:  ui,
		pkgRepoLocation: pkgLocation,
		pkgRepoBuild:    pkgRepoBuild,
	}
	return &fetchStep
}

func (fetch FetchStep) PreInteract() error {
	return nil
}

func (fetch *FetchStep) Interact() error {
	fetchSection := fetch.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch
	if fetchSection == nil {
		//Initialize fetch Section
		pkgRepo := &v1alpha1.PackageRepositoryFetch{}
		fetch.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch = pkgRepo
	}

	fetchOptionSelected := ImgpkgBundle

	switch fetchOptionSelected {
	case ImgpkgBundle:
		imgpkgStep := NewImgPkgStep(fetch.pkgAuthoringUI, fetch.pkgRepoLocation, fetch.pkgRepoBuild)
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
