package imgpkg

import (
	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuild "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"strings"
)

type RegistryDetails struct {
	RegistryURL string
}

func (createImgpkgStep CreateImgPkgStep) GetRegistryDetails() (RegistryDetails, error) {
	createImgpkgStep.pkgAuthoringUI.PrintInformationalText("\nNext is to push the imgpkg bundle created above to the OCI registry. Registry URL format: <REGISTRY_URL/REPOSITORY_NAME:TAG> e.g. index.docker.io/k8slt/sample-bundle:v0.1.0")
	defaultRegistryURL := createImgpkgStep.pkgBuild.Spec.Imgpkg.RegistryURL
	textOpts := ui.TextOpts{
		Label:        "Enter the registry url to push the bundle content",
		Default:      defaultRegistryURL,
		ValidateFunc: nil,
	}
	imgpkgPushURL, err := createImgpkgStep.pkgAuthoringUI.AskForText(textOpts)
	imgpkgPushURL = strings.TrimSpace(imgpkgPushURL)
	if err != nil {
		return RegistryDetails{}, err
	}
	imgpkgConf := pkgbuild.Imgpkg{
		RegistryURL: imgpkgPushURL,
	}
	createImgpkgStep.pkgBuild.Spec.Imgpkg = &imgpkgConf
	createImgpkgStep.pkgBuild.WriteToFile()

	return RegistryDetails{RegistryURL: imgpkgPushURL}, nil

}
