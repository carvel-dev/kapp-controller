package imgpkg

const (
	DockerHub int = iota
	PrivateRegistry
)

type RegistryAuthDetails struct {
	RegistryURL string
	Username    string
	Password    string
}

func (createImgpkg CreateImgPkgStep) PopulateRegistryAuthDetails() (RegistryAuthDetails, error) {
	registry, err := createImgpkg.Ui.AskForChoice("Where do you want to push the bundle", []string{"DockerHub", "Private Registry"})
	if err != nil {
		return RegistryAuthDetails{}, err
	}
	switch registry {
	case DockerHub:
	case PrivateRegistry:
		registryURL, err := createImgpkg.Ui.AskForText("Enter the registry URL")
		if err != nil {
			return RegistryAuthDetails{}, err
		}
		/*
			username, err := createImgpkg.Ui.AskForText("Registry UserName")
			if err != nil {
				return RegistryAuthDetails{}, err
			}
			password, err := createImgpkg.Ui.AskForPassword("Registry Password")
			if err != nil {
				return RegistryAuthDetails{}, err
			}
		*/

		return RegistryAuthDetails{RegistryURL: registryURL, Username: "", Password: ""}, nil
	}
	return RegistryAuthDetails{}, nil
}
