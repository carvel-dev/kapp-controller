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
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"sigs.k8s.io/yaml"
)

const (
	PkgBuildFileName = "package_build.yml"
)

type CreateOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	DefaultValuesFile string

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewCreateOptions(ui ui.UI, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *CreateOptions {
	return &CreateOptions{ui: ui, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
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
	createStep := NewCreateStep(o.ui, pkgLocation, pkgBuild)
	err = common.Run(createStep)
	if err != nil {
		return err
	}
	return nil
}

type CreateStep struct {
	ui           ui.UI
	pkgLocation  string
	valuesSchema v1alpha1.ValuesSchema
	template     template.TemplateStep
	pkgBuild     *build.PackageBuild
}

func NewCreateStep(ui ui.UI, pkgLocation string, pkgBuild build.PackageBuild) *CreateStep {
	return &CreateStep{
		ui:          ui,
		pkgLocation: pkgLocation,
		pkgBuild:    &pkgBuild,
	}
}

func (createStep CreateStep) getStartBlock() string {
	str := fmt.Sprintf(`
Lets start on the package creation process.
Creating directory %s
	$ mkdir -p %s
`, createStep.pkgLocation, createStep.pkgLocation)
	return str
}

func (createStep CreateStep) PreInteract() error {
	createStep.ui.BeginLinef(createStep.getStartBlock())
	err := createStep.createDirectory(createStep.pkgLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createStep CreateStep) createDirectory(dirPath string) error {
	result := util.Execute("mkdir", []string{"-p", dirPath})
	if result.Error != nil {
		createStep.ui.ErrorLinef("Error creating package directory.Error is: %s", result.Stderr)
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
func (createStep CreateStep) configureFullyQualifiedName() error {
	createStep.ui.BeginLinef(createStep.getFQPkgNameBlock())
	var fqName string
	var err error
	for {
		fqName, err = createStep.ui.AskForText("Enter the fully qualified package name")
		if err != nil {
			return err
		}
		err = validateFQName(fqName)
		if err == nil {
			break
		}
		createStep.ui.ErrorLinef("Invalid Package Name. %s", err.Error())
	}

	createStep.pkgBuild.Spec.PkgMetadata.Name = fqName
	createStep.pkgBuild.Spec.PkgMetadata.Spec.DisplayName = strings.Split(fqName, ".")[0]
	createStep.pkgBuild.Spec.Pkg.Spec.RefName = fqName
	createStep.pkgBuild.WriteToFile(createStep.pkgLocation)
	return nil
}

func (createStep *CreateStep) getFQPkgNameBlock() string {
	str := `
A package name must be a fully qualified name. 
It must consist of at least three segments separated by a '.'
Fully Qualified Name cannot have a trailing '.' e.g. samplepackage.corp.com`
	return str
}

//Get Package Version and store it in package-build.yml
func (createStep CreateStep) configurePackageVersion() error {
	createStep.ui.BeginLinef(createStep.getPkgVersionBlock())
	var pkgVersion string
	var err error
	for {
		pkgVersion, err = createStep.ui.AskForText("Enter the package version")
		if err != nil {
			return err
		}
		err = validatePackageSpecVersion(pkgVersion)
		if err == nil {
			break
		}
		createStep.ui.ErrorLinef("Invalid package version. %s", err.Error())
	}

	createStep.pkgBuild.Spec.Pkg.Spec.Version = pkgVersion
	createStep.pkgBuild.Spec.Pkg.Name = createStep.pkgBuild.Spec.Pkg.Spec.RefName + "." + pkgVersion
	createStep.pkgBuild.WriteToFile(createStep.pkgLocation)
	return nil
}

func (createStep *CreateStep) getPkgVersionBlock() string {
	str := `A package can have multiple versions. These versions are used by PackageInstall to install specific version of the package into the Kubernetes cluster.`
	return str
}

func (createStep *CreateStep) configureFetchSection() error {
	fetchConfiguration := fetch.NewFetchStep(createStep.ui, createStep.pkgLocation, createStep.pkgBuild)
	err := common.Run(fetchConfiguration)
	if err != nil {
		return err
	}
	return nil
}

func (createStep *CreateStep) configureTemplateSection() error {
	templateConfiguration := template.NewTemplateStep(createStep.ui)
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
	str := `Great, we have all the data needed to builder the package.yml and package-metadata.yml.`
	createStep.ui.BeginLinef(str)
	createStep.printPackageCR(createStep.pkgBuild.GetPackage())
	createStep.printPackageMetadataCR(createStep.pkgBuild.GetPackageMetadata())
	str = fmt.Sprintf(`
Both the files can be accessed from the following location: %s
`, createStep.pkgLocation)
	createStep.ui.PrintBlock([]byte(str))
	return nil
}

func (createStep CreateStep) printPackageCR(pkg v1alpha1.Package) error {
	str := `This is how the package.yml will look like
	$ cat package.yml
`
	createStep.ui.PrintBlock([]byte(str))

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
		createStep.ui.ErrorLinef("Unable to create package file. %s", err.Error())
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
		createStep.ui.ErrorLinef("Error printing file %s.Error is: %s", filePath, result.ErrorStr())
		return result.Error
	}
	createStep.ui.PrintBlock([]byte(result.Stdout))
	return nil
}

func (createStep CreateStep) printPackageMetadataCR(pkgMetadata v1alpha1.PackageMetadata) error {
	str := `
This is how the packageMetadata.yml will look like
	$ cat packageMetadata.yml
`
	createStep.ui.PrintBlock([]byte(str))
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
		createStep.ui.ErrorLinef("Unable to create package metadata file. %s", err.Error())
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
	return filepath.Join(pwd, "pkgBuilder")
}
