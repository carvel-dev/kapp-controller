package builder

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	pkgAuthoringUI    pkgui.IAuthoringUI
	depsFactory       cmdcore.DepsFactory
	logger            logger.Logger
	DefaultValuesFile string
	pkgVersion        string
	pkgCmdTreeOpts    cmdcore.PackageCommandTreeOpts
}

func NewCreateOptions(ui ui.UI, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *CreateOptions {
	return &CreateOptions{pkgAuthoringUI: pkgui.NewAuthoringUI(ui), logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
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

	cmd.PersistentFlags().StringVarP(&o.pkgVersion, "version", "v", "", "Version of a package (in semver format)")

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
	createStep := NewCreateStep(o.pkgAuthoringUI, pkgLocation, o.pkgVersion, pkgBuild)
	err = common.Run(createStep)
	if err != nil {
		return err
	}
	return nil
}

type CreateStep struct {
	pkgAuthoringUI pkgui.IAuthoringUI
	pkgLocation    string
	pkgVersion     string
	valuesSchema   v1alpha1.ValuesSchema
	template       template.TemplateStep
	pkgBuild       *build.PackageBuild
}

func NewCreateStep(pkgAuthorUI pkgui.IAuthoringUI, pkgLocation, pkgVersion string, pkgBuild build.PackageBuild) *CreateStep {
	return &CreateStep{
		pkgAuthoringUI: pkgAuthorUI,
		pkgLocation:    pkgLocation,
		pkgBuild:       &pkgBuild,
		pkgVersion:     pkgVersion,
	}
}

func (createStep CreateStep) printStartBlock() {
	createStep.pkgAuthoringUI.PrintHeaderText("\nPre-requisite")
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
	createStep.pkgAuthoringUI.PrintHeaderText("\nBasic Information(Step 1/3)")
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

	shortDesc := createStep.pkgBuild.Spec.PkgMetadata.Spec.ShortDescription
	if len(shortDesc) == 0 {
		createStep.pkgBuild.Spec.PkgMetadata.Spec.ShortDescription = fqName
	}

	longDesc := createStep.pkgBuild.Spec.PkgMetadata.Spec.LongDescription
	if len(longDesc) == 0 {
		createStep.pkgBuild.Spec.PkgMetadata.Spec.LongDescription = fqName
	}

	createStep.pkgBuild.Spec.Pkg.Spec.RefName = fqName
	createStep.pkgBuild.WriteToFile()
	return nil
}

func (createStep *CreateStep) printFQPkgNameBlock() {
	createStep.pkgAuthoringUI.PrintInformationalText("A package Reference name must be a valid DNS subdomain name (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names) \n - at least three segments separated by a '.', no trailing '.' e.g. samplepackage.corp.com")
}

//Get Package Version and store it in package-build.yml
func (createStep CreateStep) configurePackageVersion() error {
	var (
		pkgVersion string
		err        error
	)

	if createStep.pkgVersion != "" {
		pkgVersion = createStep.pkgVersion
	} else {
		createStep.printPkgVersionBlock()
		defaultPkgVersion := createStep.pkgBuild.Spec.Pkg.Spec.Version
		textOpts := ui.TextOpts{
			Label:        "Enter the package version",
			Default:      defaultPkgVersion,
			ValidateFunc: validatePackageSpecVersion,
		}

		pkgVersion, err = createStep.pkgAuthoringUI.AskForText(textOpts)
		if err != nil {
			return err
		}
	}

	createStep.pkgBuild.Spec.Pkg.Spec.Version = pkgVersion
	createStep.pkgBuild.Spec.Pkg.Name = createStep.pkgBuild.Spec.Pkg.Spec.RefName + "." + pkgVersion
	createStep.pkgBuild.WriteToFile()
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
	createStep.updateTimeStamp()

	createStep.pkgAuthoringUI.PrintInformationalText("\nGreat, we have all the data needed to create the package.yml and package-metadata.yml.")
	err := createStep.printPackageCR(createStep.pkgBuild.GetPackage())
	if err != nil {
		return err
	}
	err = createStep.printPackageMetadataCR(createStep.pkgBuild.GetPackageMetadata())
	if err != nil {
		return err
	}
	createStep.pkgAuthoringUI.PrintInformationalText(fmt.Sprintf("Both the files can be accessed from the following location: %s\n", createStep.pkgLocation))
	createStep.printInformation()
	createStep.printNextStep()
	return nil
}

func (createStep CreateStep) updateTimeStamp() error {
	timestamp := v1.NewTime(time.Now().UTC()).Rfc3339Copy()
	createStep.pkgBuild.Spec.Pkg.ObjectMeta.CreationTimestamp = timestamp
	createStep.pkgBuild.Spec.Pkg.Spec.ReleasedAt = timestamp
	createStep.pkgBuild.Spec.PkgMetadata.ObjectMeta.CreationTimestamp = timestamp
	err := createStep.pkgBuild.WriteToFile()
	if err != nil {
		return err
	}
	return nil
}

func (createStep CreateStep) printPackageCR(pkg v1alpha1.Package) error {
	createStep.pkgAuthoringUI.PrintActionableText("Printing package.yml")
	createStep.pkgAuthoringUI.PrintCmdExecutionText("cat package.yml")
	packageData, err := yaml.Marshal(&pkg)
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

	packageMetadataData, err := yaml.Marshal(&pkgMetadata)
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
	createStep.pkgAuthoringUI.PrintInformationalText("\n**Next steps**")
	createStep.pkgAuthoringUI.PrintInformationalText("\nCreated package can be consumed in following ways:\n1. Add the package to package repository by running `kctrl pkg repo build`. Once it has been added to the package repository, test it by running `kctrl pkg install -i <INSTALL_NAME> -p <PACKAGE_NAME> --version <VERSION>`\n2. Publish the package on the github repository.\n")
}

func (createStep CreateStep) printInformation() {
	createStep.pkgAuthoringUI.PrintInformationalText("\n**Information**\npackage-build.yml is generated as part of this flow. This file can be used for further updating and adding complex scenarios while using the `kctrl pkg build create` command. Please read the link'ed documentation for more explanation.")
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
	pkgLocation, err := filepath.Rel(pwd, filepath.Join(pwd, "pkgbuild"))
	if err != nil {
		return "", err
	}
	return pkgLocation, nil
}
