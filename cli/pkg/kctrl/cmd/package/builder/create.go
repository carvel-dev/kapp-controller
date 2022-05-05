package builder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/template"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kappctrlapis "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
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
		Aliases: []string{"g"},
		Short:   "Create a package",
		Args:    cobra.ExactArgs(1),
		RunE:    func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Create a package",
				[]string{"package", "builder", "create"},
			},
		}.Description("", o.pkgCmdTreeOpts),
		SilenceUsage: true,
	}

	return cmd
}

func (o *CreateOptions) Run(args []string) error {
	//TODO Rohit Should we provide an option to give pkg location?
	pkgLocation := GetPkgLocation()
	createStep := NewCreateStep(o.ui, pkgLocation)
	err := createStep.Run()
	if err != nil {
		return err
	}
	return nil
}

type CreateStep struct {
	ui           ui.UI
	pkgLocation  string
	pkgVersion   string
	fqName       string
	valuesSchema v1alpha1.ValuesSchema
	fetch        fetch.FetchStep
	template     template.TemplateStep
	maintainers  []v1alpha1.Maintainer
}

func NewCreateStep(ui ui.UI, pkgLocation string) *CreateStep {
	return &CreateStep{
		ui:          ui,
		pkgLocation: pkgLocation,
	}
}

func (createStep *CreateStep) Run() error {
	err := createStep.PreInteract()
	if err != nil {
		return err
	}
	err = createStep.Interact()
	if err != nil {
		return err
	}
	err = createStep.PostInteract()
	if err != nil {
		return err
	}
	return nil
}

func (createStep CreateStep) getStartBlock() string {
	str := fmt.Sprintf(`
Lets start on the package creation process.
Creating directory %s
	$ mkdir -p %s
`, createStep.pkgLocation, createStep.pkgLocation)
	return str
}

func (create CreateStep) PreInteract() error {
	create.ui.BeginLinef(create.getStartBlock())
	output, err := util.Execute("mkdir", []string{"-p", create.pkgLocation})
	if err != nil {
		create.ui.ErrorLinef("Error creating package directory.Error is: %s", err.Error())
		return err
	}
	create.ui.BeginLinef(output)
	return nil
}

func (create *CreateStep) getPkgVersionBlock() string {
	str := `
A package can have multiple versions. These versions are used by PackageInstall to install specific version of the package into the Kubernetes cluster.`
	return str
}

func (create *CreateStep) getFQPkgNameBlock() string {
	str := `A package name must be a fully qualified name. 
It must consist of at least three segments separated by a '.'
Fully Qualified Name cannot have a trailing '.' e.g. samplePackage.corp.com`
	return str
}

func (create *CreateStep) Interact() error {
	//Get Fully Qualified Name of the Package
	create.ui.BeginLinef(create.getFQPkgNameBlock())
	fqName, err := create.ui.AskForText("Enter the fully qualified package name")
	//TODO Rohit validateFQName format
	if err != nil {
		return err
	}
	create.fqName = fqName

	//Get Package Version
	create.ui.BeginLinef(create.getPkgVersionBlock())
	pkgVersion, err := create.ui.AskForText("Enter the package version")
	//TODO Rohit any validation needed here
	if err != nil {
		return err
	}
	create.pkgVersion = pkgVersion

	err = create.configureFetchSection()
	if err != nil {
		return err
	}
	/*err = create.configureTemplateSection()
	if err != nil {
		return err
	}
	err = create.configureValuesSchema()
	if err != nil {
		return err
	}

	*/

	err = create.configureMaintainers()
	if err != nil {
		return err
	}
	return nil
}

func (create *CreateStep) configureMaintainers() error {
	maintainerNames, err := create.ui.AskForText("Enter the Maintainer's Name. Multiple names can be provided by comma separated values")
	if err != nil {
		return err
	}
	var maintainers []v1alpha1.Maintainer
	for _, maintainerName := range strings.Split(maintainerNames, ",") {
		maintainers = append(maintainers, v1alpha1.Maintainer{Name: maintainerName})
	}
	create.maintainers = maintainers
	return nil
}

func (create *CreateStep) configureFetchSection() error {
	fetchConfiguration := fetch.NewFetchStep(create.ui, create.pkgLocation)
	err := fetchConfiguration.Run()
	if err != nil {
		return err
	}
	create.fetch = *fetchConfiguration
	return nil
}

func (create *CreateStep) configureTemplateSection() error {
	templateConfiguration := template.NewTemplateStep(create.ui)
	err := templateConfiguration.Run()
	if err != nil {
		return err
	}
	create.template = *templateConfiguration
	return nil
}

func (create *CreateStep) configureValuesSchema() error {
	valuesSchema, err := create.getValueSchema()
	if err != nil {
		return err
	}
	create.valuesSchema = valuesSchema
	return nil
}

func (create CreateStep) getValueSchema() (v1alpha1.ValuesSchema, error) {
	valuesSchema := v1alpha1.ValuesSchema{}
	var isValueSchemaSpecified bool
	var isValidInput bool
	input, err := create.ui.AskForText("Do you want to specify the values Schema(y/n)")
	if err != nil {
		return valuesSchema, err
	}
	for {
		isValueSchemaSpecified, isValidInput = common.ValidateInputYesOrNo(input)
		if !isValidInput {
			input, err = create.ui.AskForText("Invalid input. (must be 'y','n','Y','N')")
			if err != nil {
				return valuesSchema, err
			}
			continue
		}
		if isValueSchemaSpecified {
			valuesSchemaFileLocation, err := create.ui.AskForText("Enter the values schema file location")
			if err != nil {
				return valuesSchema, err
			}
			valuesSchemaData, err := readDataFromFile(valuesSchemaFileLocation)
			if err != nil {
				return valuesSchema, err
			}
			valuesSchema = v1alpha1.ValuesSchema{
				OpenAPIv3: runtime.RawExtension{
					Raw: valuesSchemaData,
				},
			}
		} else {
			break
		}
	}
	return valuesSchema, nil

}
func readDataFromFile(fileLocation string) ([]byte, error) {
	//TODO should we read it in a buffer
	data, err := os.ReadFile(fileLocation)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (createStep CreateStep) PostInteract() error {
	pkgBuilder, err := createStep.createPackageBuilder()
	if err != nil {
		return err
	}
	err = createStep.printPackageCR(pkgBuilder.GetPackage())
	if err != nil {
		return err
	}
	err = createStep.printPackageMetadataCR(pkgBuilder.GetPackageMetadata())
	if err != nil {
		return err
	}
	str := fmt.Sprintf(`
Both the files can be accessed from the following location: %s
`, createStep.pkgLocation)
	createStep.ui.PrintBlock([]byte(str))

	return nil
}

func (create CreateStep) createPackageBuilder() (PackageBuild, error) {
	pkgBuilder := create.populatePkgBuilder()
	jsonPackageData, err := json.Marshal(&pkgBuilder)
	if err != nil {
		return PackageBuild{}, err
	}
	yaml.JSONToYAML(jsonPackageData)
	packageData, err := yaml.JSONToYAML(jsonPackageData)
	if err != nil {
		return PackageBuild{}, err
	}
	pkgFileLocation := filepath.Join(create.pkgLocation, "package-build.yml")
	if err != nil {
		return PackageBuild{}, err
	}
	err = writeToFile(pkgFileLocation, packageData)
	if err != nil {
		create.ui.ErrorLinef("Unable to create package file. %s", err.Error())
		return PackageBuild{}, err
	}
	return pkgBuilder, nil
}

func (create CreateStep) printPackageCR(pkg v1alpha1.Package) error {
	str := `Great, we have all the data needed to builder the package.yml and package-metadata.yml. 
This is how the package.yml will look like
	$ cat package.yml
`
	create.ui.PrintBlock([]byte(str))

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
	pkgFileLocation := filepath.Join(create.pkgLocation, "package.yml")
	if err != nil {
		return err
	}
	err = writeToFile(pkgFileLocation, packageData)
	if err != nil {
		create.ui.ErrorLinef("Unable to create package file. %s", err.Error())
		return err
	}

	output, err := util.Execute("cat", []string{pkgFileLocation})
	if err != nil {
		return err
	}
	create.ui.PrintBlock([]byte(output))
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

	output, err := util.Execute("cat", []string{pkgMetadataFileLocation})
	if err != nil {
		return err
	}

	createStep.ui.PrintBlock([]byte(output))
	return nil
}

func writeToFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (create CreateStep) populatePkgMetadataSection() v1alpha1.PackageMetadata {
	packageMetadataContent := v1alpha1.PackageMetadata{
		TypeMeta:   v1.TypeMeta{Kind: "PackageMetadata", APIVersion: "data.packaging.carvel.dev/v1alpha1"},
		ObjectMeta: v1.ObjectMeta{Name: create.fqName},
		Spec: v1alpha1.PackageMetadataSpec{
			DisplayName:      strings.Split(create.fqName, ".")[0],
			LongDescription:  "A long description",
			ShortDescription: "A short description",
			ProviderName:     "",
			Maintainers:      create.maintainers,
		},
	}
	return packageMetadataContent
}

func (create CreateStep) populatePkgSection() v1alpha1.Package {
	packageContent := v1alpha1.Package{
		TypeMeta:   v1.TypeMeta{Kind: "Package", APIVersion: "data.packaging.carvel.dev/v1alpha1"},
		ObjectMeta: v1.ObjectMeta{Namespace: "default", Name: create.fqName + "." + create.pkgVersion},
		Spec: v1alpha1.PackageSpec{
			RefName:                         create.fqName,
			Version:                         create.pkgVersion,
			Licenses:                        []string{"Apache 2.0", "MIT"},
			ReleasedAt:                      v1.Time{time.Now()},
			CapactiyRequirementsDescription: "",
			ReleaseNotes:                    "",
			Template: v1alpha1.AppTemplateSpec{Spec: &kappctrlapis.AppSpec{
				ServiceAccountName: "",
				Cluster:            nil,
				Fetch:              create.fetch.AppFetch,
				Template: []kappctrlapis.AppTemplate{
					kappctrlapis.AppTemplate{Ytt: &kappctrlapis.AppTemplateYtt{}},
				},
				Deploy: []kappctrlapis.AppDeploy{
					kappctrlapis.AppDeploy{Kapp: &kappctrlapis.AppDeployKapp{}},
				},
				Paused:     false,
				Canceled:   false,
				SyncPeriod: nil,
				NoopDelete: false,
			}},
			ValuesSchema: v1alpha1.ValuesSchema{
				OpenAPIv3: runtime.RawExtension{},
			},
		},
	}
	return packageContent
}

func (createStep CreateStep) populatePkgBuilder() PackageBuild {
	pkgBuilder := PackageBuild{
		//TODO Rohit what should we call it. I think PackageBuild
		TypeMeta: v1.TypeMeta{Kind: "PackageBuilder", APIVersion: "kctrl.carvel.dev/v1alpha1"},
		Spec: Spec{
			Pkg:         createStep.populatePkgSection(),
			PkgMetadata: createStep.populatePkgMetadataSection(),
		},
	}
	return pkgBuilder
}

func GetPkgLocation() string {
	pwd, _ := os.Getwd()
	//TODO Rohit what should we call the folder name
	return filepath.Join(pwd, "pkgBuilder")
}
