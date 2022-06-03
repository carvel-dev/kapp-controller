package imgpkg

import (
	"fmt"
	"path/filepath"
	"strings"

	pkgbuild "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/fetch/imgpkg/upstream"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	"k8s.io/apimachinery/pkg/util/json"
)

type CreateImgPkgStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	pkgLocation    string
	pkgBuild       *pkgbuild.PackageBuild
}

func NewCreateImgPkgStep(ui pkgui.IPkgAuthoringUI, pkgLocation string, pkgBuild *pkgbuild.PackageBuild) *CreateImgPkgStep {
	return &CreateImgPkgStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
}

func (createImgPkgStep CreateImgPkgStep) PreInteract() error {
	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("We have to first create the imgpkg bundle.")
	//TODO ROhit
	bundleLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle")
	util.Execute("rm", []string{"-r", "-f", bundleLocation})
	err := createImgPkgStep.createBundleDir()
	if err != nil {
		return err
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) createDirectory(dirPath string) error {
	result := util.Execute("mkdir", []string{"-p", dirPath})
	if result.Error != nil {
		return fmt.Errorf("Creating package directory.\n %s", result.Stderr)
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) createBundleDir() error {
	bundleLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle")
	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("Bundle directory will act as a parent directory which will contain all the artifacts which makes up our imgpkg bundle.")
	createImgPkgStep.pkgAuthoringUI.PrintActionableText("Creating directory")
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", bundleLocation))
	err := createImgPkgStep.createDirectory(bundleLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) Interact() error {
	upstreamStep := upstream.NewUpstreamStep(createImgPkgStep.pkgAuthoringUI, createImgPkgStep.pkgLocation, createImgPkgStep.pkgBuild)
	err := common.Run(upstreamStep)
	if err != nil {
		return err
	}
	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("To push the image onto registry, ensure that `docker login` is done onto registry. \nIf not done, open a separate tab and run `docker login` and enter the valid credentials to login successfully.")
	registryDetails, err := createImgPkgStep.GetRegistryDetails()
	if err != nil {
		return err
	}
	createImgPkgStep.populateRegistryDetailsInPkgBuild(registryDetails)
	return nil
}

func (createImgPkgStep CreateImgPkgStep) populateRegistryDetailsInPkgBuild(registryDetails RegistryDetails) {
	imgpkgConf := pkgbuild.Imgpkg{
		RegistryURL: registryDetails.RegistryURL,
	}
	createImgPkgStep.pkgBuild.Spec.Imgpkg = &imgpkgConf
	createImgPkgStep.pkgBuild.WriteToFile(createImgPkgStep.pkgLocation)
	return
}

func (createImgPkgStep *CreateImgPkgStep) PostInteract() error {
	imagesFileLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle", ".imgpkg", "images.yml")
	bundleLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle")
	createImgPkgStep.pkgAuthoringUI.PrintInformationalText(`imgpkg bundle configuration is now complete. Let's use kbld to lock it down.
kbld allows to build the imgpkg bundle with immutable image references.
kbld scans a package configuration for any references to images and creates a mapping of image tags to a URL with a sha256 digest. 
This mapping will then be placed into an images.yml lock file in your bundle/.imgpkg directory. Running kbld now.`)
	createImgPkgStep.pkgAuthoringUI.PrintActionableText("Running kbld")
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("kbld --file %s --imgpkg-lock-output %s", bundleLocation, imagesFileLocation))
	err := createImgPkgStep.runKbld(bundleLocation, imagesFileLocation)
	if err != nil {
		return err
	}

	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("Lets see how the images.yml file looks like:")
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("Running cat %s", imagesFileLocation))
	err = createImgPkgStep.printFile(imagesFileLocation)
	if err != nil {
		return err
	}
	bundleURL, err := createImgPkgStep.pushImgpkgBundleToRegistry(bundleLocation)
	if err != nil {
		return err
	}
	createImgPkgStep.pkgBuild.Spec.Pkg.Spec.Template.Spec.Fetch[0].ImgpkgBundle.Image = bundleURL
	createImgPkgStep.pkgBuild.WriteToFile(createImgPkgStep.pkgLocation)
	return nil
}

func (createImgPkgStep CreateImgPkgStep) runKbld(bundleLocation, imagesFileLocation string) error {
	result := util.Execute("kbld", []string{"--file", bundleLocation, "--imgpkg-lock-output", imagesFileLocation})
	if result.Error != nil {
		return fmt.Errorf("Running kbld.\n %s", result.Stderr)
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) printFile(filePath string) error {
	result := util.Execute("cat", []string{filePath})
	if result.Error != nil {
		return fmt.Errorf("Printing file %s\n %s", filePath, result.Stderr)
	}
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	return nil
}

func (createImgPkgStep CreateImgPkgStep) pushImgpkgBundleToRegistry(bundleLoc string) (string, error) {
	pushURL := createImgPkgStep.pkgBuild.Spec.Imgpkg.RegistryURL
	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("Running imgpkg to push the bundle directory and indicate what project name and tag to give it.")
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("imgpkg push --bundle %s --file %s --json", pushURL, bundleLoc))

	//TODO Rohit it is not showing the actual error
	result := util.Execute("imgpkg", []string{"push", "--bundle", pushURL, "--file", bundleLoc, "--json"})
	if result.Error != nil {
		return "", fmt.Errorf("Imgpkg bundle push failed, check the registry url: %s", pushURL)
	}
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	bundleURL, err := getBundleURL(result.Stdout)
	if err != nil {
		return "", err
	}
	return bundleURL, nil
}

type ImgpkgPushOutput struct {
	Lines  []string    `json:"Lines"`
	Tables interface{} `json:"Tables"`
	Blocks interface{} `json:"Blocks"`
}

func getBundleURL(output string) (string, error) {
	var imgPkgPushOutput ImgpkgPushOutput
	err := json.Unmarshal([]byte(output), &imgPkgPushOutput)
	if err != nil {
		return "", err
	}
	for _, val := range imgPkgPushOutput.Lines {
		if strings.HasPrefix(val, "Pushed") {
			bundleURL := strings.Split(val, " ")[1]
			return strings.Trim(bundleURL, "'"), nil
		}
	}
	return "", fmt.Errorf("No Bundle URL generated after doing imgpkg push")
}
