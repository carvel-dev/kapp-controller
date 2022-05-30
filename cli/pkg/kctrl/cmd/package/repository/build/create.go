package build

import (
	"encoding/json"
	"fmt"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/repository/build/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

const (
	PackageRepositoryBuildFileName = "pkgrepo-build.yml"
	FilesLocation                  = "filesLocation"
	FilesLocationSeparator         = ","
	PackageMetadataFileName        = "metadata"
	PushToImgpkgBundle             = "ImgpkgBundle"
	YAMLFileExtension              = ".yaml"
	YMLFileExtension               = ".yml"
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
		Short:   "Create a package repository",
		Args:    cobra.ExactArgs(1),
		RunE:    func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Create a package",
				[]string{"package", "repository", "build", "create"},
			},
		}.Description("", o.pkgCmdTreeOpts),
		SilenceUsage: true,
	}

	return cmd
}

func (o *CreateOptions) Run(args []string) error {
	//TODO Rohit Should we provide an option to give pkg location?
	pkgRepoLocation := GetPkgRepoLocation()
	pkgRepoBuildFilePath := filepath.Join(pkgRepoLocation, PackageRepositoryBuildFileName)
	pkgRepoBuild, err := build.GeneratePackageRepositoryBuild(pkgRepoBuildFilePath)
	if err != nil {
		return err
	}
	createStep := NewCreateStep(o.pkgAuthoringUI, pkgRepoLocation, &pkgRepoBuild)
	err = common.Run(createStep)
	if err != nil {
		return err
	}
	return nil
}

type CreateStep struct {
	pkgAuthoringUI  pkgui.IPkgAuthoringUI
	pkgRepoLocation string
	pkgRepoBuild    *build.PackageRepositoryBuild
}

func NewCreateStep(pkgAuthorUI pkgui.IPkgAuthoringUI, pkgRepoLocation string, pkgRepoBuild *build.PackageRepositoryBuild) *CreateStep {
	return &CreateStep{
		pkgAuthoringUI:  pkgAuthorUI,
		pkgRepoLocation: pkgRepoLocation,
		pkgRepoBuild:    pkgRepoBuild,
	}
}

func (createStep CreateStep) printStartBlock() {
	createStep.pkgAuthoringUI.PrintInformationalText("\nLets start on the package creation process.")
	createStep.pkgAuthoringUI.PrintActionableText("\nCreating directory")
	createStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", createStep.pkgRepoLocation))
}

func (createStep CreateStep) PreInteract() error {
	createStep.printStartBlock()
	err := createStep.createDirectory(createStep.pkgRepoLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createStep CreateStep) createDirectory(dirPath string) error {
	result := util.Execute("mkdir", []string{"-p", dirPath})
	if result.Error != nil {
		createStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("Error creating package directory.Error is: %s",
			result.Stderr))
		return result.Error
	}
	return nil
}

func (createStep *CreateStep) Interact() error {

	err := createStep.configurePackageLocation()
	if err != nil {
		return err
	}

	fetchStep := NewFetchStep(createStep.pkgAuthoringUI, createStep.pkgRepoLocation, createStep.pkgRepoBuild)
	err = common.Run(fetchStep)
	if err != nil {
		return err
	}

	return nil
}

func (createStep CreateStep) configurePackageLocation() error {
	defaultfilesLocation := createStep.pkgRepoBuild.ObjectMeta.Annotations[FilesLocation]
	textOpts := ui.TextOpts{
		Label:        "Enter the location of package and packageMetadata files",
		Default:      defaultfilesLocation,
		ValidateFunc: validatePathExists,
	}
	filesLocation, err := createStep.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}

	if createStep.pkgRepoBuild.ObjectMeta.Annotations == nil {
		createStep.pkgRepoBuild.ObjectMeta.Annotations = map[string]string{}
	}
	createStep.pkgRepoBuild.ObjectMeta.Annotations[FilesLocation] = filesLocation
	createStep.pkgRepoBuild.WriteToFile(createStep.pkgRepoLocation)
	return nil
}

func (createStep CreateStep) PostInteract() error {
	return nil
}

func (createStep CreateStep) printPackageRepositoryCR(pkg v1alpha1.Package) error {
	createStep.pkgAuthoringUI.PrintInformationalText(`This is how the package.yml will look like`)
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
	pkgFileLocation := filepath.Join(createStep.pkgRepoLocation, "package.yml")
	if err != nil {
		return err
	}
	err = writeToFile(pkgFileLocation, packageData)
	if err != nil {
		createStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("Unable to create package file. %s", err.Error()))
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
		return fmt.Errorf("Error printing file %s.Error is: %s", filePath, result.ErrorStr())
	}
	createStep.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	return nil
}

func writeToFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetPkgRepoLocation() string {
	pwd, _ := os.Getwd()
	//TODO Rohit what should we call the folder name
	return filepath.Join(pwd, "repoBuild")
}
