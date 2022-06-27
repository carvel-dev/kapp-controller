// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	appInit "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init"
	appBuild "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/init/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	PkgBuildFileName     = "package-build.yml"
	PkgResourcesFileName = "package-resources.yml"
)

type InitOptions struct {
	ui          cmdcore.IAuthoringUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	KbldBuild  bool
	pkgVersion string
}

func NewInitOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *InitOptions {
	return &InitOptions{ui: cmdcore.NewAuthoringUI(ui), depsFactory: depsFactory, logger: logger}
}

func NewInitCmd(o *InitOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Package",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	cmd.PersistentFlags().StringVarP(&o.pkgVersion, "version", "v", "", "Version of a package (in semver format)")
	return cmd
}

func (o *InitOptions) Run() error {
	o.ui.PrintHeaderText("\nPre-requisite")
	o.ui.PrintInformationalText("Welcome! Before we start on the package creation journey, please ensure the following pre-requites are met:\n* The Carvel suite of tools are installed. Do get familiar with the following Carvel tools: ytt, imgpkg, vendir, and kbld.\n* You have access to an OCI registry, and authenticated locally so that images can be pushed. e.g. docker login <REGISTRY URL>\n")

	pkgBuildFilePath := filepath.Join(PkgBuildFileName)
	pkgBuild, err := build.GetPackageBuild(pkgBuildFilePath)
	if err != nil {
		return err
	}
	pkgMetadata, err := build.GetPackageMetadata(PkgResourcesFileName)
	if err != nil {
		return err
	}
	pkg, err := build.GetPackage(PkgResourcesFileName)
	if err != nil {
		return err
	}
	pkgInstall, err := build.GetPackageInstall(PkgResourcesFileName)
	if err != nil {
		return err
	}
	createStep := NewCreateStep(o.ui, pkgBuild, pkg, pkgMetadata, pkgInstall, o.logger, o.pkgVersion, o.depsFactory)
	createStep.pkg = pkg
	createStep.pkgMetadata = pkgMetadata
	err = common.Run(createStep)
	if err != nil {
		return err
	}
	return nil
}

type CreateStep struct {
	ui          cmdcore.IAuthoringUI
	pkgBuild    *build.PackageBuild
	pkg         *v1alpha1.Package
	pkgMetadata *v1alpha1.PackageMetadata
	pkgInstall  *v1alpha12.PackageInstall
	pkgVersion  string
	logger      logger.Logger
	depsFactory cmdcore.DepsFactory
}

func NewCreateStep(ui cmdcore.IAuthoringUI, pkgBuild *build.PackageBuild, pkg *v1alpha1.Package, pkgMetadata *v1alpha1.PackageMetadata, pkgInstall *v1alpha12.PackageInstall, logger logger.Logger, pkgVersion string, depsFactory cmdcore.DepsFactory) *CreateStep {
	return &CreateStep{
		ui:          ui,
		pkgBuild:    pkgBuild,
		pkg:         pkg,
		pkgMetadata: pkgMetadata,
		pkgInstall:  pkgInstall,
		pkgVersion:  pkgVersion,
		logger:      logger,
		depsFactory: depsFactory,
	}
}

func (createStep CreateStep) printStartBlock() {
	createStep.ui.PrintHeaderText("\nPre-requisite")
	createStep.ui.PrintInformationalText("Welcome! Before we start on the package creation journey, please ensure the following pre-requites are met:\n* The Carvel suite of tools are installed. Do get familiar with the following Carvel tools: ytt, imgpkg, vendir, and kbld.\n* You have access to an OCI registry, and authenticated locally so that images can be pushed. e.g. docker login <REGISTRY URL>\n")
	createStep.ui.PrintInformationalText("\nWe need a directory to hold generated Package and PackageMetadata CRs")
}

func (createStep *CreateStep) PreInteract() error {
	createStep.printStartBlock()
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
	appBuild := createStep.generateAppBuildFromPackageBuild()
	appCreateStep := appInit.NewCreateStep(createStep.ui, appBuild, createStep.logger, createStep.depsFactory, false)
	err = common.Run(appCreateStep)
	if err != nil {
		return err
	}
	createStep.populatePkgBuildFromAppBuild(appCreateStep.GetAppBuild())
	createStep.pkgBuild.WriteToFile()
	return nil
}

//Get Fully Qualified Name of the Package and store it in package-build.yml
func (createStep CreateStep) configureFullyQualifiedName() error {
	createStep.printFQPkgNameBlock()
	defaultFullyQualifiedName := createStep.pkgMetadata.Name
	textOpts := ui.TextOpts{
		Label:        "Enter the package reference name",
		Default:      defaultFullyQualifiedName,
		ValidateFunc: validateFQName,
	}
	fqName, err := createStep.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	createStep.pkgMetadata.Name = fqName
	createStep.pkgMetadata.Spec.DisplayName = strings.Split(fqName, ".")[0]

	shortDesc := createStep.pkgMetadata.Spec.ShortDescription
	if len(shortDesc) == 0 {
		createStep.pkgMetadata.Spec.ShortDescription = fqName
	}

	longDesc := createStep.pkgMetadata.Spec.LongDescription
	if len(longDesc) == 0 {
		createStep.pkgMetadata.Spec.LongDescription = fqName
	}

	createStep.pkg.Spec.RefName = fqName
	//TODO createStep saves all the resources.
	createStep.pkgBuild.WriteToFile()
	return nil
}

func (createStep *CreateStep) printFQPkgNameBlock() {
	createStep.ui.PrintInformationalText("A package Reference name must be a valid DNS subdomain name (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names) \n - at least three segments separated by a '.', no trailing '.' e.g. samplepackage.corp.com")
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
		defaultPkgVersion := createStep.pkg.Spec.Version
		textOpts := ui.TextOpts{
			Label:        "Enter the package version",
			Default:      defaultPkgVersion,
			ValidateFunc: validatePackageSpecVersion,
		}

		pkgVersion, err = createStep.ui.AskForText(textOpts)
		if err != nil {
			return err
		}
	}

	createStep.pkg.Spec.Version = pkgVersion
	createStep.pkg.Name = createStep.pkg.Spec.RefName + "." + pkgVersion
	createStep.pkgBuild.WriteToFile()
	return nil
}

func (createStep *CreateStep) printPkgVersionBlock() {
	createStep.ui.PrintInformationalText("A package version is used by PackageInstall to install particular version of the package into the Kubernetes cluster. It must be valid semver as specified by https://semver.org/spec/v2.0.0.html")
}

func (createStep *CreateStep) PostInteract() error {
	var data []byte
	pkgInstall, err := createStep.updatePackageInstall()
	if err != nil {
		return err
	}

	pkgInstallContent, err := yaml.Marshal(pkgInstall)
	if err != nil {
		return err
	}
	data = append(data, pkgInstallContent...)

	pkgContent, err := yaml.Marshal(createStep.pkg)
	if err != nil {
		return err
	}
	data = append(data, []byte("---\n")...)
	data = append(data, pkgContent...)
	err = common.WriteFile(PkgResourcesFileName, data)
	//createStep.ui.PrintInformationalText(fmt.Sprintf("Both the files can be accessed from the following location: %s\n", createStep.pkgLocation))
	createStep.printInformation()
	createStep.printNextStep()
	return nil
}

//TODO should we do like this or the way we are generating AppFromAppBuild.
func (createStep CreateStep) updatePackageInstall() (v1alpha12.PackageInstall, error) {
	existingPkgInstall := createStep.pkgInstall
	if existingPkgInstall.ObjectMeta.Annotations == nil {
		existingPkgInstall.ObjectMeta.Annotations = make(map[string]string)
		existingPkgInstall.ObjectMeta.Annotations[common.LocalFetchAnnotationKey] = "."
	}

	if len(existingPkgInstall.Name) == 0 {
		existingPkgInstall.Name = createStep.pkg.Name
	}

	if len(existingPkgInstall.Spec.ServiceAccountName) == 0 {
		existingPkgInstall.Spec.ServiceAccountName = "sa-" + createStep.pkg.Name
	}
	if existingPkgInstall.Spec.PackageRef == nil {
		existingPkgInstall.Spec.PackageRef = &v1alpha12.PackageRef{
			RefName: createStep.pkg.Spec.RefName,
		}
	}
	//TODO Check whether we should add version constraint as well.
	if len(existingPkgInstall.Spec.PackageRef.RefName) == 0 {
		existingPkgInstall.Spec.PackageRef.RefName = createStep.pkg.Spec.RefName
	}
	return *existingPkgInstall, nil
}

func (createStep CreateStep) printPackageCR(pkg v1alpha1.Package) error {
	createStep.ui.PrintActionableText("Printing package.yml")
	createStep.ui.PrintCmdExecutionText("cat package.yml")
	pkg.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now())
	pkg.Spec.ReleasedAt = metav1.NewTime(time.Now())
	packageData, err := yaml.Marshal(&pkg)
	if err != nil {
		return err
	}
	pkgFileLocation := filepath.Join("package.yml")
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
	return nil
}

func (createStep CreateStep) printPackageMetadataCR(pkgMetadata v1alpha1.PackageMetadata) error {
	createStep.ui.PrintActionableText("Printing packageMetadata.yml")
	createStep.ui.PrintCmdExecutionText("cat packageMetadata.yml")
	pkgMetadata.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now())
	packageMetadataData, err := yaml.Marshal(&pkgMetadata)
	if err != nil {
		return err
	}
	pkgMetadataFileLocation := filepath.Join("package-metadata.yml")
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
	createStep.ui.PrintInformationalText("\n**Next steps**")
	createStep.ui.PrintInformationalText("\nCreated package can be consumed in following ways:\n1. Add the package to package repository by running `kctrl pkg repo build`. Once it has been added to the package repository, test it by running `kctrl pkg install -i <INSTALL_NAME> -p <PACKAGE_NAME> --version <VERSION>`\n2. Publish the package on the github repository.\n")
}

func (createStep CreateStep) printInformation() {
	createStep.ui.PrintInformationalText("\n**Information**\npackage-build.yml is generated as part of this flow. This file can be used for further updating and adding complex scenarios while using the `kctrl pkg build create` command. Please read the link'ed documentation for more explanation.")
}

func (createStep CreateStep) generateAppBuildFromPackageBuild() *appBuild.AppBuild {
	appBuild := createStep.pkgBuild.Spec.Template
	return &appBuild
}

func (createStep CreateStep) populatePkgBuildFromAppBuild(appBuild *appBuild.AppBuild) error {
	createStep.pkgBuild.Spec.Template = *appBuild
	return nil
}

func writeToFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
