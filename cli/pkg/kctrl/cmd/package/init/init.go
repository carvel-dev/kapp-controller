// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	appInit "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/configure/fetch"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	cmdlocal "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	v1alpha13 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	PkgBuildFileName     = "package-build.yml"
	PkgResourcesFileName = "package-resources.yml"

	PkgAPIVersion         = "data.packaging.carvel.dev/v1alpha1"
	PkgMetadataAPIVersion = "data.packaging.carvel.dev/v1alpha1"
	PkgInstallAPIVersion  = "packaging.carvel.dev/v1alpha1"
	PkgKind               = "Package"
	PkgMetadataKind       = "PackageMetadata"
	PkgInstallKind        = "PackageInstall"

	defaultPkgRefName = "samplepackage.corp.com"
	defaultPkgVersion = "0.0.0"

	YAMLSeparator = "---"
)

type InitOptions struct {
	ui          cmdcore.AuthoringUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger
}

func NewInitOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *InitOptions {
	return &InitOptions{ui: cmdcore.NewAuthoringUIImpl(ui), depsFactory: depsFactory, logger: logger}
}

func NewInitCmd(o *InitOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Package",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	return cmd
}

func (o *InitOptions) Run() error {
	o.ui.PrintHeaderText("\nPre-requisite")
	o.ui.PrintInformationalText("Welcome! Before we start on the package creation journey, please ensure the following pre-requites are met:\n* The Carvel suite of tools are installed. Do get familiar with the following Carvel tools: ytt, imgpkg, vendir, and kbld.\n* You have access to an OCI registry, and authenticated locally so that images can be pushed. e.g. docker login <REGISTRY URL>\n")

	pkgBuild, err := GetPackageBuild(PkgBuildFileName)
	if err != nil {
		return err
	}
	var configs cmdlocal.Configs
	resourcesFileExists, err := common.IsFileExists(PkgResourcesFileName)
	if err != nil {
		return err
	}
	if resourcesFileExists {
		configs, err = cmdlocal.NewConfigFromFiles([]string{PkgResourcesFileName})
		if err != nil {
			return err
		}
	}

	pkg, err := getPackage(configs)
	if err != nil {
		return err
	}
	pkgMetadata, err := getPackageMetadata(configs)
	if err != nil {
		return err
	}
	// TODO we get an error if package-resources.yml file exist but there is no packageInstall in it.
	// Probably, needs to make changes to local Package and adopt them in dev deploy.
	pkgInstall, err := getPackageInstall(configs)
	if err != nil {
		return err
	}

	createStep := NewCreateStep(o.ui, pkgBuild, pkg, pkgMetadata, pkgInstall, o.logger, o.depsFactory)
	createStep.pkg = pkg
	createStep.pkgMetadata = pkgMetadata
	err = common.Run(createStep)
	if err != nil {
		return err
	}
	return nil
}

func getPackage(configs cmdlocal.Configs) (*v1alpha1.Package, error) {
	var pkg v1alpha1.Package
	if len(configs.Pkgs) > 1 {
		return nil, fmt.Errorf("More than 1 Package found")
	}
	if configs.Pkgs != nil {
		pkg = configs.Pkgs[0]
	} else {
		pkg = v1alpha1.Package{
			TypeMeta: metav1.TypeMeta{
				Kind:       PkgKind,
				APIVersion: PkgAPIVersion,
			},
		}
	}
	return &pkg, nil
}

func getPackageMetadata(configs cmdlocal.Configs) (*v1alpha1.PackageMetadata, error) {
	var pkgMetadata v1alpha1.PackageMetadata
	if len(configs.PkgMetadatas) > 1 {
		return nil, fmt.Errorf("More than 1 PackageMetadata found")
	}
	if configs.PkgMetadatas != nil {
		pkgMetadata = configs.PkgMetadatas[0]
	} else {
		pkgMetadata = v1alpha1.PackageMetadata{
			TypeMeta: metav1.TypeMeta{
				Kind:       PkgMetadataKind,
				APIVersion: PkgMetadataAPIVersion,
			},
		}
	}
	return &pkgMetadata, nil
}

func getPackageInstall(configs cmdlocal.Configs) (*v1alpha12.PackageInstall, error) {
	var pkgInstall v1alpha12.PackageInstall
	if len(configs.PkgInstalls) > 1 {
		return nil, fmt.Errorf("More than 1 PackageInstall found")
	}
	if configs.PkgMetadatas != nil {
		pkgInstall = configs.PkgInstalls[0]
	} else {
		pkgInstall = v1alpha12.PackageInstall{
			TypeMeta: metav1.TypeMeta{
				Kind:       PkgInstallKind,
				APIVersion: PkgInstallAPIVersion,
			},
		}
	}
	return &pkgInstall, nil
}

type CreateStep struct {
	ui          cmdcore.AuthoringUI
	pkgBuild    *PackageBuild
	pkg         *v1alpha1.Package
	pkgMetadata *v1alpha1.PackageMetadata
	pkgInstall  *v1alpha12.PackageInstall
	logger      logger.Logger
	depsFactory cmdcore.DepsFactory
}

func NewCreateStep(ui cmdcore.AuthoringUI, pkgBuild *PackageBuild, pkg *v1alpha1.Package, pkgMetadata *v1alpha1.PackageMetadata, pkgInstall *v1alpha12.PackageInstall, logger logger.Logger, depsFactory cmdcore.DepsFactory) *CreateStep {
	return &CreateStep{
		ui:          ui,
		pkgBuild:    pkgBuild,
		pkg:         pkg,
		pkgMetadata: pkgMetadata,
		pkgInstall:  pkgInstall,
		logger:      logger,
		depsFactory: depsFactory,
	}
}

func (createStep *CreateStep) PreInteract() error {
	return nil
}

func (createStep *CreateStep) Interact() error {
	err := createStep.configurePackageReferenceName()
	if err != nil {
		return err
	}

	appBuild := createStep.generateAppBuildFromPackageBuild()
	appCreateStep := appInit.NewCreateStep(createStep.ui, appBuild, createStep.logger, createStep.depsFactory, false)
	err = common.Run(appCreateStep)
	if err != nil {
		return err
	}
	createStep.pkgBuild.Spec.Template = *appCreateStep.GetAppBuild()
	createStep.pkgBuild.Save()
	return nil
}

func (createStep CreateStep) configurePackageReferenceName() error {
	createStep.printPackageReferenceNameBlock()
	defaultPackageRefName := createStep.getDefaultPackageRefName()
	textOpts := ui.TextOpts{
		Label:        "Enter the package reference name",
		Default:      defaultPackageRefName,
		ValidateFunc: validateRefName,
	}
	pkgRefName, err := createStep.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	createStep.pkgMetadata.Name = pkgRefName
	createStep.pkgMetadata.Spec.DisplayName = strings.Split(pkgRefName, ".")[0]

	shortDesc := createStep.pkgMetadata.Spec.ShortDescription
	if len(shortDesc) == 0 {
		createStep.pkgMetadata.Spec.ShortDescription = pkgRefName
	}

	longDesc := createStep.pkgMetadata.Spec.LongDescription
	if len(longDesc) == 0 {
		createStep.pkgMetadata.Spec.LongDescription = pkgRefName
	}

	createStep.pkg.Spec.RefName = pkgRefName
	err = createStep.Save()
	if err != nil {
		return err
	}
	return nil
}

func (createStep *CreateStep) printPackageReferenceNameBlock() {
	createStep.ui.PrintInformationalText("\nA package Reference name must be a valid DNS subdomain name (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names) \n - at least three segments separated by a '.', no trailing '.' e.g. samplepackage.corp.com")
}

func (createStep *CreateStep) PostInteract() error {
	currentTime := metav1.NewTime(time.Now())
	createStep.pkgMetadata.ObjectMeta.CreationTimestamp = currentTime
	createStep.pkg.ObjectMeta.CreationTimestamp = currentTime
	createStep.pkg.Spec.ReleasedAt = currentTime
	createStep.pkgBuild.CreationTimestamp = currentTime
	createStep.pkgBuild.ObjectMeta.Name = createStep.pkg.Spec.RefName
	createStep.pkgInstall.CreationTimestamp = currentTime

	err := createStep.updatePackageInstall()
	if err != nil {
		return err
	}

	err = createStep.updatePackage()
	if err != nil {
		return err
	}
	err = createStep.Save()
	if err != nil {
		return err
	}
	//createStep.ui.PrintInformationalText(fmt.Sprintf("Both the files can be accessed from the following location: %s\n", createStep.pkgLocation))
	createStep.printInformation()
	createStep.printNextStep()
	return nil
}

func (createStep CreateStep) updatePackageInstall() error {
	existingPkgInstall := createStep.pkgInstall
	if existingPkgInstall.ObjectMeta.Annotations == nil {
		existingPkgInstall.ObjectMeta.Annotations = make(map[string]string)
		existingPkgInstall.ObjectMeta.Annotations[fetch.LocalFetchAnnotationKey] = "."
	}

	if len(existingPkgInstall.Name) == 0 {
		existingPkgInstall.Name = createStep.pkgMetadata.Spec.DisplayName
	}

	if len(existingPkgInstall.Spec.ServiceAccountName) == 0 {
		existingPkgInstall.Spec.ServiceAccountName = createStep.pkgMetadata.Spec.DisplayName + "-sa"
	}
	if existingPkgInstall.Spec.PackageRef == nil {
		existingPkgInstall.Spec.PackageRef = &v1alpha12.PackageRef{
			RefName: createStep.pkg.Spec.RefName,
		}
	}
	// TODO Check whether we should add version constraint as well.
	if len(existingPkgInstall.Spec.PackageRef.RefName) == 0 {
		existingPkgInstall.Spec.PackageRef.RefName = createStep.pkg.Spec.RefName
	}
	return nil
}

func (createStep CreateStep) updatePackage() error {
	existingPkg := createStep.pkg

	if len(existingPkg.Spec.Version) == 0 {
		existingPkg.Spec.Version = defaultPkgVersion
	}
	existingPkg.Name = existingPkg.Spec.RefName + "." + existingPkg.Spec.Version

	if existingPkg.Spec.Template.Spec == nil {
		existingPkg.Spec.Template.Spec = &v1alpha13.AppSpec{}
	}
	if len(existingPkg.Spec.Template.Spec.Template) == 0 {
		existingPkg.Spec.Template.Spec.Template = createStep.pkgBuild.Spec.Template.Spec.App.Spec.Template
	}

	if len(existingPkg.Spec.Template.Spec.Deploy) == 0 {
		existingPkg.Spec.Template.Spec.Deploy = createStep.pkgBuild.Spec.Template.Spec.App.Spec.Deploy
	}

	return nil
}

func (createStep CreateStep) printPackageCR(pkg v1alpha1.Package) error {

	return nil
}

func (createStep CreateStep) printFile(filePath string) error {
	return nil
}

func (createStep CreateStep) printPackageMetadataCR(pkgMetadata v1alpha1.PackageMetadata) error {

	return nil
}

func (createStep CreateStep) printNextStep() {
	createStep.ui.PrintInformationalText("\n**Next steps**")
	createStep.ui.PrintInformationalText("\nCreated package can be consumed in following ways:\n1. Add the package to package repository by running `kctrl pkg repo build`. Once it has been added to the package repository, test it by running `kctrl pkg install -i <INSTALL_NAME> -p <PACKAGE_NAME> --version <VERSION>`\n2. Publish the package on the github repository.\n")
}

func (createStep CreateStep) printInformation() {
	createStep.ui.PrintInformationalText("\n**Information**\npackage-build.yml is generated as part of this flow. This file can be used for further updating and adding complex scenarios while using the `kctrl pkg build create` command. Please read the link'ed documentation for more explanation.")
}

func (createStep CreateStep) generateAppBuildFromPackageBuild() *appbuild.AppBuild {
	appBuild := createStep.pkgBuild.Spec.Template
	return &appBuild
}

func (createStep CreateStep) getDefaultPackageRefName() string {
	if len(createStep.pkgMetadata.Name) != 0 {
		return createStep.pkgMetadata.Name
	}
	return defaultPkgRefName
}

// Save method will save all the resources i.e. PackageBuild, PackageInstall, Package and PackageMetadata
func (createStep CreateStep) Save() error {
	// Save PackageBuild
	content, err := yaml.Marshal(createStep.pkgBuild)
	if err != nil {
		return err
	}
	err = WriteFile(PkgBuildFileName, content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}

	// Save Package
	content, err = yaml.Marshal(createStep.pkg)
	if err != nil {
		return err
	}
	err = WriteFile(PkgResourcesFileName, content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}

	// Add YAML Separator
	err = WriteFile(PkgResourcesFileName, []byte(YAMLSeparator+"\n"), os.O_APPEND|os.O_WRONLY)
	if err != nil {
		return err
	}

	// Save PackageMetadata
	content, err = yaml.Marshal(createStep.pkgMetadata)
	if err != nil {
		return err
	}
	err = WriteFile(PkgResourcesFileName, content, os.O_APPEND|os.O_WRONLY)
	if err != nil {
		return err
	}

	// Add YAML Separator
	err = WriteFile(PkgResourcesFileName, []byte(YAMLSeparator+"\n"), os.O_APPEND|os.O_WRONLY)
	if err != nil {
		return err
	}

	// Save PackageInstall
	content, err = yaml.Marshal(createStep.pkgInstall)
	if err != nil {
		return err
	}
	err = WriteFile(PkgResourcesFileName, content, os.O_APPEND|os.O_WRONLY)
	if err != nil {
		return err
	}

	return nil
}

// Write binary content to file
func WriteFile(filePath string, data []byte, flag int) error {
	file, err := os.OpenFile(filePath, flag, 0777)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(data)
	if err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("Unable to close the file %s\n %s", filePath, err.Error())
	}
	return nil
}
