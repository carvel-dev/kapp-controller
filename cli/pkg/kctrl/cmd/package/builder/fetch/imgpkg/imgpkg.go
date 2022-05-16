package imgpkg

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type ImgpkgStep struct {
	ui           ui.UI
	pkgLocation  string
	ImgpkgBundle v1alpha1.AppFetchImgpkgBundle
	pkgBuild     *build.PackageBuild
}

func NewImgPkgStep(ui ui.UI, pkgLocation string, pkgBuild *build.PackageBuild) *ImgpkgStep {
	imgpkg := ImgpkgStep{
		ui:          ui,
		pkgLocation: pkgLocation,
		pkgBuild:    pkgBuild,
	}
	return &imgpkg
}

func (imgpkg ImgpkgStep) PreInteract() error {
	return nil
}

func (imgpkg *ImgpkgStep) Interact() error {
	var isImgpkgBundleCreated bool
	isImgpkgBundleCreated = false

	if isImgpkgBundleCreated {
		//TODO Rohit should we add some information here.
		//imgpkg.ui.BeginLinef("")
		image, err := imgpkg.ui.AskForText("Enter the imgpkg bundle url")
		if err != nil {
			return err
		}
		imgpkg.ImgpkgBundle.Image = image
	} else {
		createImgPkgStep := NewCreateImgPkgStep(imgpkg.ui, imgpkg.pkgLocation, imgpkg.pkgBuild)
		err := common.Run(createImgPkgStep)
		if err != nil {
			return err
		}
		imgpkg.ImgpkgBundle.Image = createImgPkgStep.image
	}
	return nil
}

func (imgpkg ImgpkgStep) PostInteract() error {
	return nil
}
