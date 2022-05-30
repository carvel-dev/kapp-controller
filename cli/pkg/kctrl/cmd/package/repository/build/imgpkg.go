package build

import (
	"fmt"
	"github.com/cppforlife/go-cli-ui/ui"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	build "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/repository/build/build"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/yaml"

	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type ImgpkgStep struct {
	pkgAuthoringUI  pkgui.IPkgAuthoringUI
	pkgRepoLocation string
	pkgRepoBuild    *build.PackageRepositoryBuild
}

func NewImgPkgStep(ui pkgui.IPkgAuthoringUI, pkgRepoLocation string, pkgRepoBuild *build.PackageRepositoryBuild) *ImgpkgStep {
	imgpkg := ImgpkgStep{
		pkgAuthoringUI:  ui,
		pkgRepoLocation: pkgRepoLocation,
		pkgRepoBuild:    pkgRepoBuild,
	}
	return &imgpkg
}

func (imgpkg ImgpkgStep) PreInteract() error {
	imgpkg.pkgAuthoringUI.PrintInformationalText("We have to first create the imgpkg bundle.")
	//TODO Rohit
	err := imgpkg.createBundlePackagesDir()
	if err != nil {
		return err
	}
	err = imgpkg.createBundleDotImgpkgDir()
	if err != nil {
		return err
	}

	return nil
}

func (imgpkg ImgpkgStep) createBundleDir() error {
	bundleLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle")
	imgpkg.pkgAuthoringUI.PrintInformationalText("Bundle directory will act as a parent directory which will contain all the artifacts which makes up our imgpkg bundle.")
	imgpkg.pkgAuthoringUI.PrintActionableText("Creating directory")
	imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", bundleLocation))
	err := createDirectory(bundleLocation)
	if err != nil {
		return err
	}
	return nil
}

func (imgpkg ImgpkgStep) createBundlePackagesDir() error {
	bundleLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle", "packages")
	imgpkg.pkgAuthoringUI.PrintInformationalText("Bundle directory will act as a parent directory which will contain all the artifacts which makes up our imgpkg bundle.")
	imgpkg.pkgAuthoringUI.PrintActionableText("Creating directory")
	imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", bundleLocation))
	err := createDirectory(bundleLocation)
	if err != nil {
		return err
	}
	return nil
}

func (imgpkg ImgpkgStep) createBundleDotImgpkgDir() error {
	bundleDotImgPkgLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle", ".imgpkg")
	imgpkg.pkgAuthoringUI.PrintInformationalText(".imgpkg directory will contain the bundle’s lock file. \n" +
		"A bundle lock file has the mapping of images(referenced in the package contents such as K8s YAML configurations, etc)to its sha256 digest.")
	imgpkg.pkgAuthoringUI.PrintActionableText("Creating directory")
	imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", bundleDotImgPkgLocation))
	err := createDirectory(bundleDotImgPkgLocation)
	if err != nil {
		return err
	}
	return nil
}

func (imgpkg *ImgpkgStep) Interact() error {
	imgpkgBundleConf := imgpkg.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle
	if imgpkgBundleConf == nil {
		imgpkg.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle = &v1alpha12.AppFetchImgpkgBundle{}
	}
	defaultRegistryURL := imgpkg.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle.Image
	textOpts := ui.TextOpts{
		Label:        "Enter the registry url to push the package repository",
		Default:      defaultRegistryURL,
		ValidateFunc: nil,
	}
	registryURL, err := imgpkg.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}

	imgpkg.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle.Image = registryURL
	imgpkg.pkgRepoBuild.WriteToFile(imgpkg.pkgRepoLocation)
	return nil
}

func (imgpkg ImgpkgStep) PostInteract() error {

	filesLocation := imgpkg.pkgRepoBuild.ObjectMeta.Annotations[FilesLocation]

	for _, location := range strings.Split(filesLocation, FilesLocationSeparator) {
		filepath.Walk(location, imgpkg.copyPkgOrPkgMetadataFiles)
	}
	bundleLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle")
	bundledPackagesLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle", "packages")
	imagesFileLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle", ".imgpkg", "images.yml")
	err := runKbld(bundledPackagesLocation, imagesFileLocation)
	if err != nil {
		return err
	}
	err = imgpkg.pushImgpkgBundleToRegistry(bundleLocation)
	if err != nil {
		return err
	}
	return nil
}

func (imgpkg ImgpkgStep) copyPkgOrPkgMetadataFiles(path string, info fs.FileInfo, err error) error {
	if isYamlFile(path, info) {
		bundledPackagesLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle", "packages")
		if isPackageFile(path) {
			pkg, err := getPackage(path)
			if err != nil {
				return err
			}
			subFolderName := pkg.Spec.RefName
			err = createDirectory(filepath.Join(bundledPackagesLocation, subFolderName))
			if err != nil {
				return err
			}
			//TODO what needs to be done if this is empty
			fileName := pkg.Spec.Version
			//TODO what needs to be done if this is empty
			destinationPath := filepath.Join(bundledPackagesLocation, subFolderName, fileName+YMLFileExtension)
			err = copyFile(path, destinationPath)
			if err != nil {
				return err
			}
		} else if isPackageMetadataFile(path) {
			pkgMetadata, err := getPackageMetadata(path)
			if err != nil {
				return err
			}
			subFolderName := pkgMetadata.Name
			err = createDirectory(filepath.Join(bundledPackagesLocation, subFolderName))
			if err != nil {
				return err
			}
			fileName := PackageMetadataFileName
			destinationPath := filepath.Join(bundledPackagesLocation, subFolderName, fileName+YMLFileExtension)
			err = copyFile(path, destinationPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isYamlFile(path string, info fs.FileInfo) bool {
	if !info.IsDir() {
		fileExtension := filepath.Ext(path)
		if fileExtension == YMLFileExtension || fileExtension == YAMLFileExtension {
			return true
		}
	}
	return false
}

func isPackageFile(path string) bool {
	_, err := getPackage(path)
	if err != nil {
		return false
	} else {
		return true
	}
}

func getPackage(path string) (v1alpha1.Package, error) {
	var pkg v1alpha1.Package
	content, err := os.ReadFile(path)
	if err != nil {
		return v1alpha1.Package{}, err
	}
	//TODO should we use unmarshalstrict?
	err = yaml.Unmarshal(content, &pkg)
	if err != nil {
		return v1alpha1.Package{}, err
	}
	if pkg.Kind != "Package" {
		return v1alpha1.Package{}, fmt.Errorf("File %s is not a package yaml file", path)
	}
	return pkg, nil
}

func isPackageMetadataFile(path string) bool {
	_, err := getPackageMetadata(path)
	if err != nil {
		return false
	} else {
		return true
	}
}

func getPackageMetadata(path string) (v1alpha1.PackageMetadata, error) {

	var pkgMetadata v1alpha1.PackageMetadata
	content, err := os.ReadFile(path)
	if err != nil {
		return v1alpha1.PackageMetadata{}, err
	}
	err = yaml.Unmarshal(content, &pkgMetadata)
	if err != nil {
		return v1alpha1.PackageMetadata{}, err
	}
	if pkgMetadata.Kind != "PackageMetadata" {
		return v1alpha1.PackageMetadata{}, fmt.Errorf("File %s is not a packageMetadata yaml file", path)
	}
	return pkgMetadata, nil

}

func createDirectory(dirPath string) error {
	result := util.Execute("mkdir", []string{"-p", dirPath})
	if result.Error != nil {
		return fmt.Errorf("Creating directory: %s", result.Stderr)
	}
	return nil
}

func copyFile(src string, destination string) error {
	result := util.Execute("cp", []string{src, destination})
	if result.Error != nil {
		return fmt.Errorf("Copying file %s: %s", src, result.Stderr)
	}
	return nil
}

func runKbld(bundleLocation, imagesFileLocation string) error {
	result := util.Execute("kbld", []string{"--file", bundleLocation, "--imgpkg-lock-output", imagesFileLocation})
	if result.Error != nil {
		return fmt.Errorf("Running kbld.\n %s", result.Stderr)
	}
	return nil
}

func (imgpkg ImgpkgStep) pushImgpkgBundleToRegistry(bundleLoc string) error {
	pushURL := imgpkg.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle.Image
	imgpkg.pkgAuthoringUI.PrintInformationalText("Running imgpkg to push the bundle directory and indicate what project name and tag to give it.")
	imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("imgpkg push --bundle %s --file %s --json", pushURL, bundleLoc))
	result := util.Execute("imgpkg", []string{"push", "--bundle", pushURL, "--file", bundleLoc, "--json"})
	//TODO Rohit it is not showing the actual error
	if result.Error != nil {
		return fmt.Errorf("Imgpkg bundle push failed, check the registry url: %s", pushURL)
	}
	imgpkg.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	return nil
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
