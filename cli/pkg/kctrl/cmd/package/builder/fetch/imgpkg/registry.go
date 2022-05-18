package imgpkg

type RegistryDetails struct {
	RegistryURL string
}

func (createImgpkg CreateImgPkgStep) PopulateRegistryDetails() (RegistryDetails, error) {
	//TODO Rohit have to add example for both docker hub as well as private registry.
	imgpkgPushURL, err := createImgpkg.ui.AskForText("Enter the registry url to push the bundle content")
	if err != nil {
		return RegistryDetails{}, err
	}
	return RegistryDetails{RegistryURL: imgpkgPushURL}, nil

	return RegistryDetails{}, nil
}
