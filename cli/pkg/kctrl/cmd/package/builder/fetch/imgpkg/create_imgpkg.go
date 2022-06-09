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
	//TODO ROhit
	return nil
}

func (createImgPkgStep CreateImgPkgStep) createDirectory(dirPath string) error {
	result := util.Execute("mkdir", []string{"-p", dirPath})
	if result.Error != nil {
		return fmt.Errorf("Creating package directory.\n %s", result.Stderr)
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) Interact() error {
	upstreamStep := upstream.NewUpstreamStep(createImgPkgStep.pkgAuthoringUI, createImgPkgStep.pkgLocation, createImgPkgStep.pkgBuild)
	err := common.Run(upstreamStep)
	if err != nil {
		return err
	}
	//createImgPkgStep.pkgAuthoringUI.PrintInformationalText("To push the image onto registry, ensure that `docker login` is done onto registry. If not done, open a separate tab and run `docker login` and enter the valid credentials to login successfully.")
	registryDetails, err := createImgPkgStep.GetRegistryDetails()
	if err != nil {
		return err
	}
	createImgPkgStep.populateRegistryDetailsInPkgBuild(registryDetails)
	createImgPkgStep.pkgAuthoringUI.PrintHeading("Creating package(Step 3/3)")
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
	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("imgpkg bundle configuration is now complete.")
	err := createImgPkgStep.runKbld(bundleLocation, imagesFileLocation)
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
	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("\nLet's use `kbld` to create immutable image reference. Kbld scans all the files in bundle configuration for any references of images and creates a mapping of image tags to a URL with sha256 digest.")
	if createImgPkgStep.isHelmContent() {
		tempLocation := filepath.Join(createImgPkgStep.pkgLocation, "temp")
		chartLocation, err := getPathFromVendirConf(createImgPkgStep.pkgBuild)
		if err != nil {
			return err
		}
		helmChartLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle", chartLocation)
		createImgPkgStep.pkgAuthoringUI.PrintInformationalText("Kbld needs valid yml files. Most of the Helm Charts are templated, which means kbld can't parse them as it is.First, run `helm template` on the chart to create a valid yml files. And then we will run kbld on them. Output of helm template will be stored in the temp directory")
		createImgPkgStep.pkgAuthoringUI.PrintActionableText("Running helm template to create valid yml files")
		createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("helm3 template %s --output-dir %s", helmChartLocation, tempLocation))
		result := util.Execute("helm3", []string{"template", helmChartLocation, "--output-dir", tempLocation})
		if result.Error != nil {
			return fmt.Errorf("Running helm template.\n %s", result.Stderr)
		}
		bundleLocation = tempLocation
	}
	createImgPkgStep.pkgAuthoringUI.PrintActionableText("Lock image references using Kbld")
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("kbld --file %s --imgpkg-lock-output %s", bundleLocation, imagesFileLocation))
	result := util.Execute("kbld", []string{"--file", bundleLocation, "--imgpkg-lock-output", imagesFileLocation})
	if result.Error != nil {
		return fmt.Errorf("Running kbld.\n %s", result.Stderr)
	}

	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("\nKbld places the mapping of image tag to its sha digest in images.yml lock file")
	createImgPkgStep.pkgAuthoringUI.PrintActionableText("Printing images.yml")
	createImgPkgStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("Running cat %s", imagesFileLocation))
	err := createImgPkgStep.printFile(imagesFileLocation)
	if err != nil {
		return err
	}

	return nil
}

func getPathFromVendirConf(pkgBuild *pkgbuild.PackageBuild) (string, error) {
	//This means that helmChart has been mentioned directly in the fetch section of Package.
	if pkgBuild.Spec.Vendir == nil {
		return "", nil
	}

	directories := pkgBuild.Spec.Vendir.Directories
	if directories == nil {
		return "", fmt.Errorf("No helm chart reference in Vendir")
	}

	var path string
	for _, directory := range directories {
		directoryPath := directory.Path
		for _, content := range directories[0].Contents {
			if content.HelmChart != nil {
				path = directoryPath + "/" + content.Path
				break
			}
		}
	}
	return path, nil
}

func (createImgPkgStep CreateImgPkgStep) isHelmContent() bool {
	if createImgPkgStep.pkgBuild.Annotations[common.PkgFetchContentAnnotationKey] == common.FetchChartFromHelmRepo {
		return true
	}
	return false
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
	createImgPkgStep.pkgAuthoringUI.PrintInformationalText("\nNow that our imgpkg bundle is created, we will push the bundle directory by running `imgpkg push`. `Push` command allows users to push the imgpkg bundle from local to registry for consumption.")
	createImgPkgStep.pkgAuthoringUI.PrintActionableText("Running imgpkg push")
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
