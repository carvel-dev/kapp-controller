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
	err = createImgPkgStep.createBundleConfigDir()
	if err != nil {
		return err
	}
	err = createImgPkgStep.createBundleDotImgpkgDir()
	if err != nil {
		return err
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) createDirectory(dirPath string) error {
	result := util.Execute("mkdir", []string{"-p", dirPath})
	if result.Error != nil {
		createImgPkgStep.pkgAuthoringUI.PrintErrorText(fmt.Sprintf("Error creating package directory.Error is: %s", result.ErrorStr()))
		return result.Error
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

func (createImgPkgStep CreateImgPkgStep) createBundleConfigDir() error {
	bundleConfigLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle", "config")
	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("Config directory will contain package contents such as Kubernetes YAML configuration, ytt templates, Helm templates, etc.")
	createImgPkgStep.pkgAuthoringUI.PrintActionableText("Creating directory")
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", bundleConfigLocation))
	err := createImgPkgStep.createDirectory(bundleConfigLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) createBundleDotImgpkgDir() error {
	bundleDotImgPkgLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle", ".imgpkg")
	createImgPkgStep.pkgAuthoringUI.PrintInformationalText(".imgpkg directory will contain the bundle’s lock file. \n" +
		"A bundle lock file has the mapping of images(referenced in the package contents such as K8s YAML configurations, etc)to its sha256 digest.")
	createImgPkgStep.pkgAuthoringUI.PrintActionableText("Creating directory")
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", bundleDotImgPkgLocation))
	err := createImgPkgStep.createDirectory(bundleDotImgPkgLocation)
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
	err := createImgPkgStep.runkbld(bundleLocation, imagesFileLocation)
	if err != nil {
		createImgPkgStep.pkgAuthoringUI.PrintErrorText(err.Error())
		return err
	}

	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("\nLets see how the images.yml file looks like:")
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

func (createImgPkgStep CreateImgPkgStep) runkbld(bundleLocation, imagesFileLocation string) error {
	result := util.Execute("kbld", []string{"--file", bundleLocation, "--imgpkg-lock-output", imagesFileLocation})
	if result.Error != nil {
		createImgPkgStep.pkgAuthoringUI.PrintErrorText(fmt.Sprintf("Error running kbld.Error is: %s", result.ErrorStr()))
		return result.Error
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) printFile(filePath string) error {
	result := util.Execute("cat", []string{filePath})
	if result.Error != nil {
		createImgPkgStep.pkgAuthoringUI.PrintErrorText(fmt.Sprintf("Error printing file %s.Error is: %s", filePath, result.ErrorStr()))
		return result.Error
	}
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	return nil
}

func (createImgPkgStep CreateImgPkgStep) pushImgpkgBundleToRegistry(bundleLoc string) (string, error) {
	pushURL := createImgPkgStep.pkgBuild.Spec.Imgpkg.RegistryURL
	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("\nRunning imgpkg to push the bundle directory and indicate what project name and tag to give it.")
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("imgpkg push --bundle %s --file %s --json", pushURL, bundleLoc))

	//TODO Rohit it is not showing the actual error
	result := util.Execute("imgpkg", []string{"push", "--bundle", pushURL, "--file", bundleLoc, "--json"})
	if result.Error != nil {
		createImgPkgStep.pkgAuthoringUI.PrintErrorText(fmt.Sprintf("Imgpkg bundle push failed. Reason: %s", result.Stderr))
		return "", result.Error
	}
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	bundleURL := getBundleURL(result.Stdout)
	return bundleURL, nil
}

type ImgpkgPushOutput struct {
	Lines  []string    `json:"Lines"`
	Tables interface{} `json:"Tables"`
	Blocks interface{} `json:"Blocks"`
}

func getBundleURL(output string) string {
	var imgPkgPushOutput ImgpkgPushOutput
	json.Unmarshal([]byte(output), &imgPkgPushOutput)
	for _, val := range imgPkgPushOutput.Lines {
		if strings.HasPrefix(val, "Pushed") {
			bundleURL := strings.Split(val, " ")[1]
			return strings.Trim(bundleURL, "'")
		}
	}
	return ""

}
