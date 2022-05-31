package build

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/repository/build/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	createStep.pkgAuthoringUI.PrintInformationalText("Lets start on the package repository creation process.")
	createStep.pkgAuthoringUI.PrintInformationalText("A package repository is a collection of Package and PackageMetadata CRs.")
	createStep.pkgAuthoringUI.PrintInformationalText("Similar to a maven repository or a rpm repository, adding a package repository to a cluster gives users of that cluster the ability to install any of the packages from that repository.")
	createStep.pkgAuthoringUI.PrintInformationalText("We need a directory to act as parent directory. This will be used to store all the information and files required/needed in the package repository creation journey.")
	createStep.pkgAuthoringUI.PrintActionableText("Creating directory")
	createStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("mkdir -p %s", createStep.pkgRepoLocation))
}

func (createStep CreateStep) PreInteract() error {
	createStep.printPreRequisite()
	createStep.printStartBlock()
	err := createStep.createDirectory(createStep.pkgRepoLocation)
	if err != nil {
		return err
	}
	return nil
}

func (createStep CreateStep) printPreRequisite() {
	createStep.pkgAuthoringUI.PrintInformationalText(`Welcome! Before we start on the package creation journey, please ensure the following pre-requites are met:
* The Carvel suite of tools are installed. Do get familiar with the following Carvel tools: imgpkg, kbld, etc.
* You have access to an OCI registry, and you have authenticated locally so that you can push images. e.g. docker login <OCI_REGISTRY_URL>
`)
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

	createStep.pkgAuthoringUI.PrintInformationalText("A package repository name is the name with which it will be referenced while deploying on the cluster.")
	defaultPkgRepoName := createStep.pkgRepoBuild.Spec.PkgRepo.Name
	textOpts := ui.TextOpts{
		Label:        "Enter the package repository name",
		Default:      defaultPkgRepoName,
		ValidateFunc: nil,
	}
	pkgRepoName, err := createStep.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}
	createStep.pkgRepoBuild.Spec.PkgRepo.Name = pkgRepoName
	createStep.pkgRepoBuild.WriteToFile(createStep.pkgRepoLocation)
	err = createStep.configurePackageRepositoryLocation()
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

func (createStep CreateStep) configurePackageRepositoryLocation() error {
	defaultfilesLocation := createStep.pkgRepoBuild.ObjectMeta.Annotations[FilesLocation]
	textOpts := ui.TextOpts{
		Label:        "Enter the directory which contains package and packageMetadata files. Multiple directories can be entered in comma separated format",
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
	createStep.printPackageRepositoryCR()
	createStep.pkgAuthoringUI.PrintInformationalText("This file can be taken and deployed on the Kubernetes cluster to have access to all the packages available to install as part of this repository.")
	createStep.pkgAuthoringUI.PrintInformationalText(fmt.Sprintf("Alternatively, use kctrl to deploy this package repository on the Kubernetes cluster by running `kctrl package repository add -r demo-pkg-repo --url %s`", createStep.pkgRepoBuild.Spec.PkgRepo.Spec.Fetch.ImgpkgBundle.Image))
	return nil
}

func (createStep CreateStep) printPackageRepositoryCR() error {
	createStep.pkgAuthoringUI.PrintInformationalText("Great, we have all the data needed to create the packageRepository.yml")
	pkgRepo := createStep.pkgRepoBuild.Spec.PkgRepo
	pkgRepo.ObjectMeta.CreationTimestamp = v1.NewTime(time.Now())
	pkgRepoData, err := yaml.Marshal(createStep.pkgRepoBuild.Spec.PkgRepo)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	pkgRepoFileLocation := filepath.Join(createStep.pkgRepoLocation, "packageRepository.yml")
	if err != nil {
		return err
	}
	err = writeToFile(pkgRepoFileLocation, pkgRepoData)
	if err != nil {
		createStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("Unable to create package file. %s", err.Error()))
		return err
	}

	createStep.pkgAuthoringUI.PrintActionableText(`Printing packageRepository.yml`)
	createStep.pkgAuthoringUI.PrintCmdExecutionText("cat packageRepository.yml")

	err = createStep.printFile(pkgRepoFileLocation)
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
