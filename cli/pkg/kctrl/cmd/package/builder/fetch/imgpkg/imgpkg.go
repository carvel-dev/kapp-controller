package imgpkg

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type ImgpkgStep struct {
	ui           ui.UI
	pkgLocation  string
	ImgpkgBundle v1alpha1.AppFetchImgpkgBundle
}

func NewImgPkgStep(ui ui.UI, pkgLocation string) *ImgpkgStep {
	imgpkg := ImgpkgStep{
		ui:          ui,
		pkgLocation: pkgLocation,
	}
	return &imgpkg
}

func (imgpkg ImgpkgStep) PreInteract() error {
	return nil
}

func (imgpkg *ImgpkgStep) Interact() error {
	var isImgpkgCreated bool
	input, _ := imgpkg.ui.AskForText("Is the imgpkg bundle already created(y/n)")

	for {
		var isValidInput bool
		isImgpkgCreated, isValidInput = common.ValidateInputYesOrNo(input)
		if isValidInput {
			break
		} else {
			input, _ = imgpkg.ui.AskForText("Invalid input (must be 'y','n','Y','N')")
		}
	}

	if isImgpkgCreated {

	} else {
		createImgPkgStep := NewCreateImgPkgStep(imgpkg.ui, imgpkg.pkgLocation)
		err := createImgPkgStep.Run()
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

func (imgpkg *ImgpkgStep) Run() error {
	err := imgpkg.PreInteract()
	if err != nil {
		return err
	}
	err = imgpkg.Interact()
	if err != nil {
		return err
	}
	err = imgpkg.PostInteract()
	if err != nil {
		return err
	}
	return nil
}
