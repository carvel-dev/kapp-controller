package build

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	build "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/repository/build/build"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"sigs.k8s.io/yaml"
)

type ImgpkgStep struct {
	pkgAuthoringUI  pkgui.IAuthoringUI
	pkgRepoLocation string
	pkgRepoBuild    *build.PackageRepositoryBuild
}

func NewImgPkgStep(ui pkgui.IAuthoringUI, pkgRepoLocation string, pkgRepoBuild *build.PackageRepositoryBuild) *ImgpkgStep {
	imgpkg := ImgpkgStep{
		pkgAuthoringUI:  ui,
		pkgRepoLocation: pkgRepoLocation,
		pkgRepoBuild:    pkgRepoBuild,
	}
	return &imgpkg
}

func (imgpkg ImgpkgStep) PreInteract() error {
	imgpkg.pkgAuthoringUI.PrintInformationalText("To create package repository, we will create an imgpkg bundle first. Imgpkg, a Carvel tool, allows users to package, distribute, and relocate a set of files as one OCI artifact: a bundle. Imgpkg bundles are identified with a unique sha256 digest based on the file contents. Imgpkg uses that digest to ensure that the copied contents are identical to those originally pushed.")
	imgpkg.pkgAuthoringUI.PrintInformationalText("\nA package repository bundle is an imgpkg bundle that holds PackageMetadata and Package CRs. Later on, this bundle can be mentioned in the package repository CR to fetch the package and packageMetadata CRs.")
	imgpkg.pkgAuthoringUI.PrintActionableText("\nCreating the required directory structure for imgpkg bundle\n")
	err := imgpkg.createBundleDir()
	if err != nil {
		return err
	}
	err = imgpkg.createBundlePackagesDir()
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
	util.Execute("rm", []string{"-r", "-f", imgpkg.pkgRepoLocation, "bundle"})
	bundleLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle")
	imgpkg.pkgAuthoringUI.PrintInformationalText("To understand the directory structure of package repository bundle and the purpose of each subdirectory, refer: https://carvel.dev/kapp-controller/docs/latest/packaging-artifact-formats/#package-repository-bundle")
	imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", bundleLocation))
	err := createDirectory(bundleLocation)
	if err != nil {
		return err
	}
	return nil
}

func (imgpkg ImgpkgStep) createBundlePackagesDir() error {
	bundleLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle", "packages")
	imgpkg.pkgAuthoringUI.PrintInformationalText("Packages directory will contain all the Package and PackageMetadata CRs which makes up imgpkg bundle.")
	imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", bundleLocation))
	err := createDirectory(bundleLocation)
	if err != nil {
		return err
	}
	return nil
}

func (imgpkg ImgpkgStep) createBundleDotImgpkgDir() error {
	bundleDotImgPkgLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle", ".imgpkg")
	imgpkg.pkgAuthoringUI.PrintInformationalText(".imgpkg directory will contain the bundle’s lock file. A bundle lock file has the mapping of images(referenced in the package contents such as K8s YAML configurations, etc)to its sha256 digest.")
	//imgpkg.pkgAuthoringUI.PrintInformationalText("It ensures that later on while deployment, we are using the same exact image which we used while creating the bundle as digest are immutable even though tags are.")
	imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", bundleDotImgPkgLocation))
	err := createDirectory(bundleDotImgPkgLocation)
	if err != nil {
		return err
	}
	return nil
}

func (imgpkg *ImgpkgStep) Interact() error {
	imgpkg.pkgAuthoringUI.PrintHeaderText("\nRegistry URL(Step 2/3)")
	imgpkgBundleConf := imgpkg.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle
	if imgpkgBundleConf == nil {
		imgpkg.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle = &v1alpha12.AppFetchImgpkgBundle{}
	}
	defaultRegistryURL := imgpkg.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle.Image
	imgpkg.pkgAuthoringUI.PrintInformationalText("\nNext is to push the imgpkg bundle created above to the OCI registry. Registry URL format: <REGISTRY_URL/REPOSITORY_NAME:TAG> e.g. index.docker.io/k8slt/sample-bundle:v0.1.0")
	textOpts := ui.TextOpts{
		Label:        "Enter the registry url to push the package repository bundle",
		Default:      defaultRegistryURL,
		ValidateFunc: nil,
	}
	registryURL, err := imgpkg.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}

	imgpkg.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle.Image = registryURL
	imgpkg.pkgRepoBuild.WriteToFile(imgpkg.pkgRepoLocation)
	imgpkg.pkgAuthoringUI.PrintHeaderText("\nCreating Package Repository(Step 3/3)")
	return nil
}

func (imgpkg ImgpkgStep) PostInteract() error {

	filesLocation := "packages"
	for _, location := range strings.Split(filesLocation, FilesLocationSeparator) {
		filepath.Walk(location, imgpkg.copyPkgOrPkgMetadataFiles)
	}
	bundleLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle")
	bundledPackagesLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle", "packages")
	imagesFileLocation := filepath.Join(imgpkg.pkgRepoLocation, "bundle", ".imgpkg", "images.yml")

	imgpkg.pkgAuthoringUI.PrintInformationalText("imgpkg bundle configuration is now complete.")
	imgpkg.pkgAuthoringUI.PrintInformationalText("\nLet's use `kbld` to create immutable image reference. Kbld scans all the files in bundle configuration for any references of images and creates a mapping of image tags to a URL with sha256 digest.")
	imgpkg.pkgAuthoringUI.PrintActionableText("Lock image references using Kbld")
	imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("kbld --file %s --imgpkg-lock-output %s", bundleLocation, imagesFileLocation))

	err := runningKbld(bundledPackagesLocation, imagesFileLocation)
	if err != nil {
		return err
	}
	imgpkg.pkgAuthoringUI.PrintInformationalText("\nKbld places the mapping of image tag to its sha digest in images.yml lock file")
	imgpkg.pkgAuthoringUI.PrintActionableText("Printing images.yml")
	imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("cat %s", imagesFileLocation))
	err = imgpkg.printFile(imagesFileLocation)

	err = imgpkg.pushImgpkgBundleToRegistry(bundleLocation)
	if err != nil {
		return err
	}
	imgpkg.pkgAuthoringUI.PrintInformationalText(fmt.Sprintf("We have successfully pushed the package repository imgpkg bundle to the OCI registry.We can use `%s` in our package repository CR fetch section to have access to our package and packageMetadata CRs",
		imgpkg.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle.Image))

	imgpkg.pkgRepoBuild.WriteToFile(imgpkg.pkgRepoLocation)
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
			imgpkg.pkgAuthoringUI.PrintActionableText("Copying file")
			imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("cp %s %s", path, destinationPath))
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
			imgpkg.pkgAuthoringUI.PrintActionableText("Copying file")
			imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("cp %s %s", path, destinationPath))
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

func runningKbld(bundleLocation, imagesFileLocation string) error {
	result := util.Execute("kbld", []string{"--file", bundleLocation, "--imgpkg-lock-output", imagesFileLocation})
	if result.Error != nil {
		return fmt.Errorf("Running kbld.\n %s", result.Stderr)
	}
	return nil
}

func (imgpkg ImgpkgStep) pushImgpkgBundleToRegistry(bundleLoc string) error {
	pushURL := imgpkg.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle.Image
	imgpkg.pkgAuthoringUI.PrintInformationalText("\nNow that our imgpkg bundle is created, we will push the bundle directory by running `imgpkg push`. `Push` command allows users to push the imgpkg bundle from local to registry for consumption.")
	imgpkg.pkgAuthoringUI.PrintActionableText("Running imgpkg push")
	imgpkg.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("imgpkg push --bundle %s --file %s --json", pushURL, bundleLoc))
	result := util.Execute("imgpkg", []string{"push", "--bundle", pushURL, "--file", bundleLoc, "--json"})
	//TODO Rohit it is not showing the actual error
	if result.Error != nil {
		return fmt.Errorf("Imgpkg bundle push failed, check the registry url: %s", pushURL)
	}
	imgpkg.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	return nil
}

func (imgpkg ImgpkgStep) printFile(filePath string) error {
	result := util.Execute("cat", []string{filePath})
	if result.Error != nil {
		return fmt.Errorf("Printing file %s\n %s", filePath, result.Stderr)
	}
	imgpkg.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	return nil
}
