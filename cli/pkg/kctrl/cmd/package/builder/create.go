package builder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/template"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"sigs.k8s.io/yaml"
)

type CreateOptions struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	depsFactory    cmdcore.DepsFactory
	logger         logger.Logger

	DefaultValuesFile string

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewCreateOptions(ui ui.UI, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *CreateOptions {
	return &CreateOptions{pkgAuthoringUI: pkgui.NewPackageAuthoringUI(ui), logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewCreateCmd(o *CreateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a package",
		Args:    cobra.ExactArgs(1),
		RunE:    func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Create a package",
				[]string{"package", "build", "create"},
			},
		}.Description("", o.pkgCmdTreeOpts),
		SilenceUsage: true,
	}

	return cmd
}

func (o *CreateOptions) Run(args []string) error {
	//TODO Rohit Should we provide an option to give pkg location?
	pkgLocation, err := GetPkgLocation()
	if err != nil {
		return err
	}
	pkgBuildFilePath := filepath.Join(common.PkgBuildFileName)
	pkgBuild, err := build.GeneratePackageBuild(pkgBuildFilePath)
	if err != nil {
		return err
	}
	createStep := NewCreateStep(o.pkgAuthoringUI, pkgLocation, pkgBuild)
	err = common.Run(createStep)
	if err != nil {
		return err
	}
	return nil
}

type CreateStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	pkgLocation    string
	valuesSchema   v1alpha1.ValuesSchema
	template       template.TemplateStep
	pkgBuild       *build.PackageBuild
}

func NewCreateStep(pkgAuthorUI pkgui.IPkgAuthoringUI, pkgLocation string, pkgBuild build.PackageBuild) *CreateStep {
	return &CreateStep{
		pkgAuthoringUI: pkgAuthorUI,
		pkgLocation:    pkgLocation,
		pkgBuild:       &pkgBuild,
	}
}

func (createStep CreateStep) printStartBlock() {
	createStep.pkgAuthoringUI.PrintHeading("Pre-requisite")
	createStep.pkgAuthoringUI.PrintInformationalText("Welcome! Before we start on the package creation journey, please ensure the following pre-requites are met:\n* The Carvel suite of tools are installed. Do get familiar with the following Carvel tools: ytt, imgpkg, vendir, and kbld.\n* You have access to an OCI registry, and authenticated locally so that images can be pushed. e.g. docker login <REGISTRY URL>\n")
	createStep.pkgAuthoringUI.PrintInformationalText("\nWe need a directory to hold generated Package and PackageMetadata CRs")
	createStep.pkgAuthoringUI.PrintActionableText("Creating directory")
	createStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", createStep.pkgLocation))
}

func (createStep CreateStep) PreInteract() error {
	createStep.printStartBlock()
	err := createStep.createDirectory(createStep.pkgLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createStep CreateStep) createDirectory(dirPath string) error {
	result := util.Execute("mkdir", []string{"-p", dirPath})
	if result.Error != nil {
		return fmt.Errorf("Creating package directory: %s", result.Stderr)
	}
	return nil
}

func (createStep *CreateStep) Interact() error {
	createStep.pkgAuthoringUI.PrintHeading("\nBasic Information(Step 1/3)")
	err := createStep.configureFullyQualifiedName()
	if err != nil {
		return err
	}

	err = createStep.configurePackageVersion()
	if err != nil {
		return err
	}
	/*
		isDefaultPreferenceImmutable := getPackageCreatePreferenceFromPkgBuild(createStep.pkgBuild)
		createStep.pkgAuthoringUI.PrintInformationalText("A package can be created in two ways. Either having immutable reference to the fetch'ed manifest or a direct reference to the fetch'ed manifest. To create an immutable reference, we will use imgpkg and kbld to lock the reference.")
		pkgCreateOptions := []string{common.CreateWithImmutableReference, common.CreateWithDirectReference}
		var defaultPreferenceOptionIndex int
		if isDefaultPreferenceImmutable {
			defaultPreferenceOptionIndex = 0
		} else {
			defaultPreferenceOptionIndex = 1
		}
		choiceOpts := ui.ChoiceOpts{
			Label:   "How to package the manifest",
			Default: defaultPreferenceOptionIndex,
			Choices: pkgCreateOptions,
		}
		preferenceOptionSelectedIndex, err := createStep.pkgAuthoringUI.AskForChoice(choiceOpts)
		if err != nil {
			return err
		}
		if preferenceOptionSelectedIndex == 0 {
			createStep.pkgBuild.Annotations[common.PkgCreatePreferenceAnnotationKey] = "true"
			bundleLocation := filepath.Join(createStep.pkgLocation, "bundle")
			createStep.pkgAuthoringUI.PrintInformationalText("To create an immutable reference, we need to create an imgpkg bundle. Imgpkg, a Carvel tool, allows users to package, distribute, and relocate a set of files as one OCI artifact: a bundle. Imgpkg bundles are identified with a unique sha256 digest based on the file contents. Imgpkg uses that digest to ensure that the copied contents are identical to those originally pushed.")
			createStep.pkgAuthoringUI.PrintInformationalText("\nCleaning up any previous imgpkg bundle directory present.")
			createStep.pkgAuthoringUI.PrintActionableText("Cleaning up any previous bundle directory")
			createStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("rm -r -f %s", bundleLocation))
			util.Execute("rm", []string{"-r", "-f", bundleLocation})
			err := createStep.createBundleDir()
			if err != nil {
				return err
			}

			err = createStep.createBundleDotImgpkgDir()
			if err != nil {
				return err
			}

		} else {
			createStep.pkgBuild.Annotations[common.PkgCreatePreferenceAnnotationKey] = "false"
		}
		createStep.pkgBuild.WriteToFile(createStep.pkgLocation)
	*/
	bundleLocation := filepath.Join(createStep.pkgLocation, "bundle")
	createStep.pkgAuthoringUI.PrintInformationalText("We need to create an imgpkg bundle as part of the package creation process. Imgpkg, a Carvel tool, allows users to package, distribute, and relocate a set of files as one OCI artifact: a bundle. Imgpkg bundles are identified with a unique sha256 digest based on the file contents. Imgpkg uses that digest to ensure that the copied contents are identical to those originally pushed.")
	createStep.pkgAuthoringUI.PrintInformationalText("\nCleaning up any previous imgpkg bundle directory present.")
	createStep.pkgAuthoringUI.PrintActionableText("Cleaning up any previous bundle directory")
	createStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("rm -r -f %s", bundleLocation))
	util.Execute("rm", []string{"-r", "-f", bundleLocation})
	err = createStep.createBundleDir()
	if err != nil {
		return err
	}

	err = createStep.createBundleDotImgpkgDir()
	if err != nil {
		return err
	}
	err = createStep.configureFetchSection()
	if err != nil {
		return err
	}

	err = createStep.configureTemplateSection()
	if err != nil {
		return err
	}

	/*		err = createStep.configureValuesSchema()
			if err != nil {
				return err
			}

	*/
	return nil
}

func getPackageCreatePreferenceFromPkgBuild(pkgBuild *build.PackageBuild) bool {
	if pkgBuild.Annotations[common.PkgCreatePreferenceAnnotationKey] == "false" {
		return false
	}
	return true
}

func (createStep CreateStep) createBundleDir() error {
	bundleLocation := filepath.Join(createStep.pkgLocation, "bundle")
	createStep.pkgAuthoringUI.PrintInformationalText("\nBundle directory will act as a parent directory which will contain all the artifacts which makes up our imgpkg bundle.")
	createStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", bundleLocation))
	err := createStep.createDirectory(bundleLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createStep CreateStep) createBundleDotImgpkgDir() error {
	bundleDotImgPkgLocation := filepath.Join(createStep.pkgLocation, "bundle", ".imgpkg")
	createStep.pkgAuthoringUI.PrintInformationalText("\n.imgpkg directory will contain the bundle’s lock file. A bundle lock file has the mapping of images(referenced in the package contents such as K8s YAML configurations, etc)to its sha256 digest.")
	createStep.pkgAuthoringUI.PrintActionableText("Creating directory")
	createStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", bundleDotImgPkgLocation))
	err := createStep.createDirectory(bundleDotImgPkgLocation)
	if err != nil {
		return err
	}
	return nil
}

//Get Fully Qualified Name of the Package and store it in package-build.yml
func (createStep CreateStep) configureFullyQualifiedName() error {
	createStep.printFQPkgNameBlock()
	defaultFullyQualifiedName := createStep.pkgBuild.Spec.PkgMetadata.Name
	textOpts := ui.TextOpts{
		Label:        "Enter the package reference name",
		Default:      defaultFullyQualifiedName,
		ValidateFunc: validateFQName,
	}
	fqName, err := createStep.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}

	createStep.pkgBuild.Spec.PkgMetadata.Name = fqName
	createStep.pkgBuild.Spec.PkgMetadata.Spec.DisplayName = strings.Split(fqName, ".")[0]
	createStep.pkgBuild.Spec.Pkg.Spec.RefName = fqName
	createStep.pkgBuild.WriteToFile(createStep.pkgLocation)
	return nil
}

func (createStep *CreateStep) printFQPkgNameBlock() {
	createStep.pkgAuthoringUI.PrintInformationalText("A package Reference name must be a valid DNS subdomain name (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names) \n - at least three segments separated by a '.', no trailing '.' e.g. samplepackage.corp.com")
}

//Get Package Version and store it in package-build.yml
func (createStep CreateStep) configurePackageVersion() error {
	createStep.printPkgVersionBlock()
	defaultPkgVersion := createStep.pkgBuild.Spec.Pkg.Spec.Version
	textOpts := ui.TextOpts{
		Label:        "Enter the package version",
		Default:      defaultPkgVersion,
		ValidateFunc: validatePackageSpecVersion,
	}

	pkgVersion, err := createStep.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}
	createStep.pkgBuild.Spec.Pkg.Spec.Version = pkgVersion
	createStep.pkgBuild.Spec.Pkg.Name = createStep.pkgBuild.Spec.Pkg.Spec.RefName + "." + pkgVersion
	createStep.pkgBuild.WriteToFile(createStep.pkgLocation)
	return nil
}

func (createStep *CreateStep) printPkgVersionBlock() {
	createStep.pkgAuthoringUI.PrintInformationalText("A package version is used by PackageInstall to install particular version of the package into the Kubernetes cluster. It must be valid semver as specified by https://semver.org/spec/v2.0.0.html")
}

func (createStep *CreateStep) configureFetchSection() error {
	fetchConfiguration := fetch.NewFetchStep(createStep.pkgAuthoringUI, createStep.pkgLocation, createStep.pkgBuild)
	err := common.Run(fetchConfiguration)
	if err != nil {
		return err
	}
	return nil
}

func (createStep *CreateStep) configureTemplateSection() error {
	templateConfiguration := template.NewTemplateStep(createStep.pkgAuthoringUI, createStep.pkgLocation, createStep.pkgBuild)
	err := common.Run(templateConfiguration)
	if err != nil {
		return err
	}
	return nil
}

func (createStep *CreateStep) configureValuesSchema() error {
	valuesSchema, err := createStep.getValueSchema()
	if err != nil {
		return err
	}
	createStep.valuesSchema = valuesSchema
	return nil
}

func (createStep CreateStep) PostInteract() error {
	createStep.pkgAuthoringUI.PrintInformationalText("Great, we have all the data needed to create the package.yml and package-metadata.yml.")
	err := createStep.printPackageCR(createStep.pkgBuild.GetPackage())
	if err != nil {
		return err
	}
	err = createStep.printPackageMetadataCR(createStep.pkgBuild.GetPackageMetadata())
	if err != nil {
		return err
	}
	createStep.pkgAuthoringUI.PrintInformationalText(fmt.Sprintf("Both the files can be accessed from the following location: %s\n", createStep.pkgLocation))
	createStep.printNextStep()
	return nil
}

func (createStep CreateStep) printPackageCR(pkg v1alpha1.Package) error {
	createStep.pkgAuthoringUI.PrintActionableText("Printing package.yml")
	createStep.pkgAuthoringUI.PrintCmdExecutionText("cat package.yml")

	//TODO: remove this comment. Marshal will make yaml/json

	jsonPackageData, err := json.Marshal(&pkg)
	if err != nil {
		return err
	}
	yaml.JSONToYAML(jsonPackageData)
	packageData, err := yaml.JSONToYAML(jsonPackageData)
	if err != nil {
		return err
	}
	pkgFileLocation := filepath.Join(createStep.pkgLocation, "package.yml")
	if err != nil {
		return err
	}
	err = writeToFile(pkgFileLocation, packageData)
	if err != nil {
		return fmt.Errorf("Unable to create package file. %s", err.Error())
	}

	err = createStep.printFile(pkgFileLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createStep CreateStep) printFile(filePath string) error {
	result := util.Execute("cat", []string{filePath})
	if result.Error != nil {
		return fmt.Errorf("Printing file %s\n %s", filePath, result.Stderr)
	}
	createStep.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	return nil
}

func (createStep CreateStep) printPackageMetadataCR(pkgMetadata v1alpha1.PackageMetadata) error {
	createStep.pkgAuthoringUI.PrintActionableText("Printing packageMetadata.yml")
	createStep.pkgAuthoringUI.PrintCmdExecutionText("cat packageMetadata.yml")
	jsonPackageMetadataData, err := json.Marshal(&pkgMetadata)
	if err != nil {
		return err
	}
	packageMetadataData, err := yaml.JSONToYAML(jsonPackageMetadataData)
	if err != nil {
		return err
	}
	pkgMetadataFileLocation := filepath.Join(createStep.pkgLocation, "package-metadata.yml")
	if err != nil {
		return err
	}
	err = writeToFile(pkgMetadataFileLocation, packageMetadataData)
	if err != nil {
		return fmt.Errorf("Unable to create package metadata file. %s", err.Error())
	}

	err = createStep.printFile(pkgMetadataFileLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createStep CreateStep) printNextStep() {
	createStep.pkgAuthoringUI.PrintInformationalText("\nCreated package can be consumed in following ways:\n1. Package-build.yml is the source of truth for the package authoring. Currently, the package.yml and packageMetadata.yml is created with some default values e.g. Release Notes, Short description, Long description etc. You can change these values to actual values in package-build.yml and rerun `kctrl pkg build create -y`. `-y` will run it in the non-interactive mode and will pick the values described in the package-build.yml\n2. Add the package to package repository by running `kctrl pkg repo build`. Once it has been added to the package repository, test it by running `kctrl pkg install -i <INSTALL_NAME> -p <PACKAGE_NAME> --version <VERSION>\n3. Publish the package on the github repository.\n4. For local testing, add the package to Kubernetes cluster by running `kubectl apply -f package.yml` and then install it by running `kctrl pkg install -i <INSTALL_NAME> -p <PACKAGE_NAME> --version <VERSION> \n")
}

func writeToFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetPkgLocation() (string, error) {
	pwd, _ := os.Getwd()
	//TODO Rohit what should we call the folder name
	pkgLocation, err := filepath.Rel(pwd, filepath.Join(pwd, "pkgBuild"))
	if err != nil {
		return "", err
	}
	return pkgLocation, nil
}
