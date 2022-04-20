package builder

import (
	"encoding/json"
	"fmt"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/util"
	"os"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/template"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

type CreateOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	ValuesSchema      bool
	DefaultValuesFile string

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewCreateOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *CreateOptions {
	return &CreateOptions{ui: ui, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewCreateCmd(o *CreateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"g"},
		Short:   "Create a package",
		Args:    cobra.ExactArgs(1),
		RunE:    func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Create a package",
				[]string{"package", "builder", "create", "-p", "pkg-a"},
			},
		}.Description("-p", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations:  map[string]string{"table": ""},
	}

	o.NamespaceFlags.Set(cmd, flagsFactory, o.pkgCmdTreeOpts)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "package", "p", "", "Set package name (required)")
	} else {
		//TODO Rohit - Revisit this
		cmd.Use = "get PACKAGE_NAME or PACKAGE_NAME/VERSION"
		cmd.Args = cobra.ExactArgs(1)
	}
	return cmd
}

func (create CreateStep) getStartBlock() []byte {
	str := fmt.Sprintf(`
# Lets start on the package creation for %s
# First we need a directory to store all configurations
# Creating directory %s
# mkdir -p %s
`, create.PkgName, create.PkgLocation, create.PkgLocation)
	return []byte(str)
}

type CreateStep struct {
	Ui                 ui.UI
	PkgLocation        string
	PkgVersionLocation string
	PkgVersion         string
	PkgName            string
	FqName             string
	ValuesSchema       v1alpha1.ValuesSchema
	Fetch              fetch.FetchStep
	Template           template.TemplateStep
	Maintainers        []v1alpha1.Maintainer
}

func NewCreateStep(ui ui.UI, pkgName string, pkgLocation string) *CreateStep {
	return &CreateStep{
		Ui:          ui,
		PkgName:     pkgName,
		PkgLocation: pkgLocation,
	}
}

func (create CreateStep) PreInteract() error {
	create.Ui.PrintBlock(create.getStartBlock())
	output, err := util.Execute("mkdir", []string{"-p", create.PkgLocation})
	if err != nil {
		create.Ui.ErrorLinef("Error creating package directory.Error is: %s", err.Error())
		return err
	}
	create.Ui.PrintBlock([]byte(output))
	return nil
}

func (create *CreateStep) getPkgVersionBlock() []byte {
	str := `
# A package can have multiple versions. These versions are used by PackageInstall to install specific version of the package into the Kubernetes cluster. 
`
	return []byte(str)
}

func (create *CreateStep) getFQPkgNameBlock() []byte {
	str := `# A package name must be a fully qualified name. 
# It must consist of at least three segments separated by a '.'
# Fully Qualified Name cannot have a trailing '.' e.g. Valid input: samplePackage.corp.com`
	return []byte(str)
}

func (create *CreateStep) Interact() error {
	//Get Package Version
	create.Ui.PrintBlock(create.getPkgVersionBlock())
	pkgVersion, err := create.Ui.AskForText("Enter the package version")
	if err != nil {
		return err
	}
	create.PkgVersion = pkgVersion
	pkgVersionLocation := fmt.Sprintf("%s/%s", create.PkgLocation, create.PkgVersion)
	str := fmt.Sprintf(`# All the files related to this version will be stored in the following location: %s
# Let's create this directory.
# mkdir -p %s`, pkgVersionLocation, pkgVersionLocation)
	create.Ui.PrintLinef(str)
	output, err := util.Execute("mkdir", []string{"-p", pkgVersionLocation})
	if err != nil {
		return err
	}
	create.Ui.PrintBlock([]byte(output))
	create.PkgVersionLocation = pkgVersionLocation
	//Get Fully Qualified Name of the Package
	create.Ui.PrintBlock(create.getFQPkgNameBlock())
	defaultFQName := fmt.Sprintf("%s.corp.com", create.PkgName)
	fqName, err := create.Ui.AskForText(fmt.Sprintf("Enter the fully qualified package name(default: %s)", defaultFQName))
	//TODO Rohit should we perform the validation
	if err != nil {
		return err
	}
	if fqName == "" {
		fqName = defaultFQName
	}
	create.FqName = fqName

	fetchConfiguration := fetch.NewFetchStep(create.Ui, create.PkgName, create.PkgLocation, create.PkgVersionLocation)
	fetchConfiguration.Run()
	create.Fetch = *fetchConfiguration

	templateConfiguration := template.NewTemplateStep(create.Ui)
	templateConfiguration.Run()
	create.Template = *templateConfiguration

	valuesSchema, err := create.getValueSchema()
	if err != nil {
		return err
	}
	create.ValuesSchema = valuesSchema

	maintainerNames, err := create.Ui.AskForText("Enter the Maintainer's Name. Multiple names can be provided by comma separated values")
	if err != nil {
		return err
	}
	var maintainers []v1alpha1.Maintainer
	for _, maintainerName := range strings.Split(maintainerNames, ",") {
		maintainers = append(maintainers, v1alpha1.Maintainer{Name: maintainerName})
	}
	create.Maintainers = maintainers
	return nil
}

func (create CreateStep) getValueSchema() (v1alpha1.ValuesSchema, error) {
	valuesSchema := v1alpha1.ValuesSchema{}
	var isValueSchemaSpecified bool
	var isValidInput bool
	input, err := create.Ui.AskForText("Do you want to specify the values Schema(y/n)")
	if err != nil {
		return valuesSchema, err
	}
	for {
		isValueSchemaSpecified, isValidInput = common.ValidateInputYesOrNo(input)
		if !isValidInput {
			input, err = create.Ui.AskForText("Invalid input. (must be 'y','n','Y','N')")
			if err != nil {
				return valuesSchema, err
			}
			continue
		}
		if isValueSchemaSpecified {
			valuesSchemaFileLocation, err := create.Ui.AskForText("Enter the values schema file location")
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

func (create CreateStep) PostInteract() error {
	err := create.createAndPrintPackageCR()
	if err != nil {
		return err
	}
	err = create.createAndPrintPackageMetadataCR()
	if err != nil {
		return err
	}
	return nil
}

func (create CreateStep) createAndPrintPackageCR() error {
	str := `# Great, we have all the data needed to builder the package.yml and packageMetadata.yml. 
# This is how the package.yml will look like
# cat package.yaml
`
	create.Ui.PrintBlock([]byte(str))
	pkg := create.populatePkgFromCreate()

	//TODO: remove this comment. Marshal will make yaml/json
	jsonPackageData, err := json.Marshal(&pkg)
	yaml.JSONToYAML(jsonPackageData)
	packageData, err := yaml.JSONToYAML(jsonPackageData)
	pkgFileLocation := create.PkgVersionLocation + "/package.yaml"
	if err != nil {
		return err
	}
	writeToFile(pkgFileLocation, packageData)

	output, err := util.Execute("cat", []string{pkgFileLocation})
	if err != nil {
		return err
	}
	create.Ui.PrintBlock([]byte(output))
	return nil
}

func (create CreateStep) createAndPrintPackageMetadataCR() error {
	str := `
# This is how the packageMetadata.yml will look like
# cat packageMetadata.yaml
`
	create.Ui.PrintBlock([]byte(str))
	pkgMetadata := create.populatePkgMetadataFromCreate()
	jsonPackageMetadataData, err := json.Marshal(&pkgMetadata)
	packageMetadataData, err := yaml.JSONToYAML(jsonPackageMetadataData)
	pkgMetadataFileLocation := create.PkgVersionLocation + "/packageMetadata.yaml"
	if err != nil {
		return err
	}
	writeToFile(pkgMetadataFileLocation, packageMetadataData)

	output, err := util.Execute("cat", []string{pkgMetadataFileLocation})
	if err != nil {
		return err
	}

	create.Ui.PrintBlock([]byte(output))
	str = fmt.Sprintf(`# Both the files can be accessed from the following location: %s
`, create.PkgVersionLocation)
	create.Ui.PrintBlock([]byte(str))
	return nil
}

func writeToFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (create CreateStep) populatePkgMetadataFromCreate() v1alpha1.PackageMetadata {
	packageMetadataContent := v1alpha1.PackageMetadata{
		TypeMeta:   v1.TypeMeta{Kind: "PackageMetadata", APIVersion: "data.packaging.carvel.dev/v1alpha1"},
		ObjectMeta: v1.ObjectMeta{Name: create.FqName},
		Spec: v1alpha1.PackageMetadataSpec{
			DisplayName:      create.PkgName,
			LongDescription:  "A long description",
			ShortDescription: "A short description",
			ProviderName:     "VMWare",
			Maintainers:      create.Maintainers,
		},
	}
	return packageMetadataContent
}

func (create CreateStep) populatePkgFromCreate() v1alpha1.Package {
	packageContent := v1alpha1.Package{
		TypeMeta:   v1.TypeMeta{Kind: "Package", APIVersion: "data.packaging.carvel.dev/v1alpha1"},
		ObjectMeta: v1.ObjectMeta{Namespace: "default", Name: create.FqName + "." + create.PkgVersion},
		Spec: v1alpha1.PackageSpec{
			RefName:                         create.FqName,
			Version:                         create.PkgVersion,
			Licenses:                        []string{"Apache 2.0", "MIT"},
			ReleasedAt:                      v1.Time{time.Now()},
			CapactiyRequirementsDescription: "",
			ReleaseNotes:                    "",
			Template: v1alpha1.AppTemplateSpec{Spec: &v1alpha12.AppSpec{
				ServiceAccountName: "",
				Cluster:            nil,
				Fetch:              create.Fetch.AppFetch,
				Template:           create.Template.AppTemplate,
				Deploy: []v1alpha12.AppDeploy{
					v1alpha12.AppDeploy{Kapp: &v1alpha12.AppDeployKapp{}},
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

func (o *CreateOptions) Run(args []string) error {
	//TODO Rohit Should we provide an option to give pkg location?
	pkgLocation := GetPkgLocation(o.Name)
	createPkg := NewCreateStep(o.ui, o.Name, pkgLocation)
	createPkg.PreInteract()
	createPkg.Interact()
	createPkg.PostInteract()
	return nil
}

func GetPkgLocation(pkgName string) string {
	var pkgLocation string
	pkgLocation, _ = os.UserHomeDir()
	return pkgLocation + "/.kctrl/" + pkgName
}
