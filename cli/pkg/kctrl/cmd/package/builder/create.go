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

const (
	PkgBuildFileName = "package-build.yml"
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
	pkgLocation := GetPkgLocation()
	pkgBuildFilePath := filepath.Join(pkgLocation, PkgBuildFileName)
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
	createStep.pkgAuthoringUI.PrintInformationalText("\nLets start on the package creation process.")
	createStep.pkgAuthoringUI.PrintActionableText("\nCreating directory")
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
		createStep.pkgAuthoringUI.PrintErrorText(fmt.Sprintf("Error creating package directory.Error is: %s",
			result.Stderr))
		return result.Error
	}
	return nil
}

func (createStep *CreateStep) Interact() error {

	err := createStep.configureFullyQualifiedName()
	if err != nil {
		return err
	}

	err = createStep.configurePackageVersion()
	if err != nil {
		return err
	}

	err = createStep.configureFetchSection()
	if err != nil {
		return err
	}
	/*
		err = createStep.configureTemplateSection()
		if err != nil {
			return err
		}

			err = createStep.configureValuesSchema()
			if err != nil {
				return err
			}

	*/

	return nil
}

//Get Fully Qualified Name of the Package and store it in package-build.yml
//Get Fully Qualified Name of the Package and store it in package-build.yml
func (createStep CreateStep) configureFullyQualifiedName() error {
	createStep.printFQPkgNameBlock()
	defaultFullyQualifiedName := createStep.pkgBuild.Spec.PkgMetadata.Name
	textOpts := ui.TextOpts{
		Label:        "Enter the fully qualified package name",
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
	createStep.pkgAuthoringUI.PrintInformationalText(`
A package name must be a fully qualified name. 
It must consist of at least three segments separated by a '.'
Fully Qualified Name cannot have a trailing '.' e.g. samplepackage.corp.com`)
}

//Get Package Version and store it in package-build.yml
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
	createStep.pkgAuthoringUI.PrintInformationalText(`A package can have multiple versions. 
These versions are used by PackageInstall to install specific version of the package into the Kubernetes cluster.`)
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
	templateConfiguration := template.NewTemplateStep(createStep.pkgAuthoringUI)
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
	createStep.pkgAuthoringUI.PrintInformationalText("Great, we have all the data needed to build the package.yml and package-metadata.yml.")
	createStep.printPackageCR(createStep.pkgBuild.GetPackage())
	createStep.printPackageMetadataCR(createStep.pkgBuild.GetPackageMetadata())
	createStep.pkgAuthoringUI.PrintInformationalText(fmt.Sprintf("Both the files can be accessed from the following location: %s\n", createStep.pkgLocation))
	return nil
}

func (createStep CreateStep) printPackageCR(pkg v1alpha1.Package) error {
	createStep.pkgAuthoringUI.PrintInformationalText("This is how the package.yml will look like")
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
		createStep.pkgAuthoringUI.PrintErrorText(fmt.Sprintf("Unable to create package file. %s", err.Error()))
		return err
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
		createStep.pkgAuthoringUI.PrintErrorText(fmt.Sprintf("Error printing file %s.Error is: %s", filePath, result.ErrorStr()))
		return result.Error
	}
	createStep.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	return nil
}

func (createStep CreateStep) printPackageMetadataCR(pkgMetadata v1alpha1.PackageMetadata) error {
	createStep.pkgAuthoringUI.PrintInformationalText("\nThis is how the packageMetadata.yml will look like")
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
		createStep.pkgAuthoringUI.PrintErrorText(fmt.Sprintf("Unable to create package metadata file. %s", err.Error()))
		return err
	}

	err = createStep.printFile(pkgMetadataFileLocation)
	if err != nil {
		return err
	}
	return nil
}

func writeToFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetPkgLocation() string {
	pwd, _ := os.Getwd()
	//TODO Rohit what should we call the folder name
	return filepath.Join(pwd, "pkgBuild")
}
