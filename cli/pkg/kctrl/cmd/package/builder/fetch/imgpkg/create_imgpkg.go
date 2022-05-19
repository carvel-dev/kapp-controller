package imgpkg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/fetch/imgpkg/upstream"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	"k8s.io/apimachinery/pkg/util/json"
)

type CreateImgPkgStep struct {
	ui          ui.UI
	pkgLocation string
	pkgBuild    *pkgbuilder.PackageBuild
}

func NewCreateImgPkgStep(ui ui.UI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *CreateImgPkgStep {
	return &CreateImgPkgStep{
		ui:          ui,
		pkgLocation: pkgLocation,
		pkgBuild:    pkgBuild,
	}
}

func (createImgPkgStep CreateImgPkgStep) PreInteract() error {
	createImgPkgStep.ui.BeginLinef("We have to first create the imgpkg bundle.")
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
		createImgPkgStep.ui.ErrorLinef("Error creating package directory.Error is: %s", result.ErrorStr())
		return result.Error
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) createBundleDir() error {
	bundleLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle")
	str := fmt.Sprintf(` 
Bundle directory will act as a parent directory which will contain all the artifacts which makes up our imgpkg bundle.
Creating directory %s.
	$ mkdir -p %s
`, bundleLocation, bundleLocation)
	createImgPkgStep.ui.BeginLinef(str)
	err := createImgPkgStep.createDirectory(bundleLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) createBundleConfigDir() error {
	bundleConfigLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle", "config")
	str := fmt.Sprintf(`
Config directory will contain package contents such as Kubernetes YAML configuration, ytt templates, Helm templates, etc.
Creating directory %s. 
	$ mkdir -p %s
`, bundleConfigLocation, bundleConfigLocation)
	createImgPkgStep.ui.BeginLinef(str)
	err := createImgPkgStep.createDirectory(bundleConfigLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) createBundleDotImgpkgDir() error {
	bundleDotImgPkgLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle", ".imgpkg")
	str := fmt.Sprintf(`
.imgpkg directory will contain the bundle’s lock file. A bundle lock file has the mapping of images(referenced in the package contents such as K8s YAML configurations, etc)to its sha256 digest.
Creating directory %s. 
	$ mkdir -p %s
`, bundleDotImgPkgLocation, bundleDotImgPkgLocation)
	createImgPkgStep.ui.BeginLinef(str)
	err := createImgPkgStep.createDirectory(bundleDotImgPkgLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) Interact() error {
	upstreamStep := upstream.NewUpstreamStep(createImgPkgStep.ui, createImgPkgStep.pkgLocation, createImgPkgStep.pkgBuild)
	err := common.Run(upstreamStep)
	if err != nil {
		return err
	}

	createImgPkgStep.ui.BeginLinef("To push the image onto registry, ensure that `docker login` is done onto registry. If not done, open a separate tab and run `docker login` and enter the valid credentials to login successfully.")
	registryDetails, err := createImgPkgStep.PopulateRegistryDetails()
	if err != nil {
		return err
	}
	createImgPkgStep.populateRegistryDetailsInPkgBuild(registryDetails)
	return nil
}

func (createImgPkgStep *CreateImgPkgStep) PostInteract() error {
	imagesFileLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle", ".imgpkg", "images.yml")
	bundleLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle")
	str := fmt.Sprintf(`
imgpkg bundle configuration is now complete. Let's use kbld to lock it down.
kbld allows to build the imgpkg bundle with immutable image references.
kbld scans a package configuration for any references to images and creates a mapping of image tags to a URL with a sha256 digest. 
This mapping will then be placed into an images.yml lock file in your bundle/.imgpkg directory. Running kbld now.
	$ kbld --file %s --imgpkg-lock-output %s`, bundleLocation, imagesFileLocation)
	createImgPkgStep.ui.BeginLinef(str)

	err := createImgPkgStep.runkbld(bundleLocation, imagesFileLocation)
	if err != nil {
		createImgPkgStep.ui.BeginLinef(err.Error())
		return err
	}

	str = fmt.Sprintf(`
Lets see how the images.yml file looks like:
Running cat %s
`, imagesFileLocation)
	createImgPkgStep.ui.BeginLinef(str)
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
		createImgPkgStep.ui.ErrorLinef("Error running kbld.Error is: %s", result.ErrorStr())
		return result.Error
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) printFile(filePath string) error {
	result := util.Execute("cat", []string{filePath})
	if result.Error != nil {
		createImgPkgStep.ui.ErrorLinef("Error printing file %s.Error is: %s", filePath, result.ErrorStr())
		return result.Error
	}
	createImgPkgStep.ui.PrintBlock([]byte(result.Stdout))
	return nil
}

func (createImgPkgStep CreateImgPkgStep) pushImgpkgBundleToRegistry(bundleLoc string) (string, error) {

	pushURL := createImgPkgStep.pkgBuild.Spec.Imgpkg.RegistryURL
	str := fmt.Sprintf(`Running imgpkg to push the bundle directory and indicate what project name and tag to give it.
 	$ imgpkg push --bundle %s --file %s --json
`, pushURL, bundleLoc)
	createImgPkgStep.ui.BeginLinef(str)

	//TODO Rohit it is not showing the actual error
	result := util.Execute("imgpkg", []string{"push", "--bundle", pushURL, "--file", bundleLoc, "--registry-username", "--json"})
	if result.Error != nil {
		createImgPkgStep.ui.ErrorLinef("Imgpkg bundle push failed. Reason: %s", result.Stderr)
		return "", result.Error
	}
	createImgPkgStep.ui.BeginLinef(result.Stdout)
	bundleURL := getBundleURL(result.Stdout)
	return bundleURL, nil
}

func (createImgPkgStep CreateImgPkgStep) populateRegistryDetailsInPkgBuild(registryDetails RegistryDetails) {
	imgpkgConf := pkgbuilder.Imgpkg{
		RegistryURL: registryDetails.RegistryURL,
	}
	createImgPkgStep.pkgBuild.Spec.Imgpkg = &imgpkgConf
	createImgPkgStep.pkgBuild.WriteToFile(createImgPkgStep.pkgLocation)
	return
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
