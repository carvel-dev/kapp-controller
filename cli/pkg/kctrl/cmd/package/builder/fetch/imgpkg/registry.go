package imgpkg

import "github.com/cppforlife/go-cli-ui/ui"

type RegistryDetails struct {
	RegistryURL string
}

func (createImgpkg CreateImgPkgStep) GetRegistryDetails() (RegistryDetails, error) {
	//TODO Rohit have to add example for both docker hub as well as private registry.
	defaultRegistryURL := createImgpkg.pkgBuild.Spec.Imgpkg.RegistryURL
	textOpts := ui.TextOpts{
		Label:        "Enter the registry url to push the bundle content",
		Default:      defaultRegistryURL,
		ValidateFunc: nil,
	}
	imgpkgPushURL, err := createImgpkg.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return RegistryDetails{}, err
	}
	return RegistryDetails{RegistryURL: imgpkgPushURL}, nil

	return RegistryDetails{}, nil
}
