// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	appinit "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init"
	interfaces "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init"
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
	chdir       string
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

	cmd.Flags().StringVar(&o.chdir, "chdir", "", "Working directory with package-build and other config")
	return cmd
}

func (o *InitOptions) Run() error {

	if o.chdir != "" {
		err := os.Chdir(o.chdir)
		if err != nil {
			return err
		}
	}

	o.ui.PrintHeaderText("\nPrerequisites")
	o.ui.PrintInformationalText("Welcome! Before we start, please ensure the following pre-requites are met:\n- The Carvel suite of tools are installed. Do get familiar with the following Carvel tools: ytt, imgpkg, vendir, and kbld.\n")

	pkgBuild, err := GetPackageBuild(PkgBuildFileName)
	if err != nil {
		return err
	}
	var configs cmdlocal.Configs
	resourcesFileExists, err := interfaces.IsFileExists(PkgResourcesFileName)
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
	err = interfaces.Run(createStep)
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
	build       interfaces.Build
	pkg         *v1alpha1.Package
	pkgMetadata *v1alpha1.PackageMetadata
	pkgInstall  *v1alpha12.PackageInstall
	logger      logger.Logger
	depsFactory cmdcore.DepsFactory
}

func NewCreateStep(ui cmdcore.AuthoringUI, pkgBuild interfaces.Build, pkg *v1alpha1.Package, pkgMetadata *v1alpha1.PackageMetadata, pkgInstall *v1alpha12.PackageInstall, logger logger.Logger, depsFactory cmdcore.DepsFactory) *CreateStep {
	return &CreateStep{
		ui:          ui,
		build:       pkgBuild,
		pkg:         pkg,
		pkgMetadata: pkgMetadata,
		pkgInstall:  pkgInstall,
		logger:      logger,
		depsFactory: depsFactory,
	}
}

func (c *CreateStep) PreInteract() error { return nil }

func (c *CreateStep) Interact() error {
	c.ui.PrintHeaderText("\nBasic Information (Step 1/3)")
	err := c.configurePackageReferenceName()
	if err != nil {
		return err
	}

	appCreateStep := appinit.NewCreateStep(c.ui, c.build, c.logger, c.depsFactory, false)
	err = interfaces.Run(appCreateStep)
	if err != nil {
		return err
	}
	c.build.Save()
	return nil
}

func (c CreateStep) configurePackageReferenceName() error {
	c.printPackageReferenceNameBlock()
	defaultPackageRefName := c.getDefaultPackageRefName()
	textOpts := ui.TextOpts{
		Label:        "Enter the package reference name",
		Default:      defaultPackageRefName,
		ValidateFunc: validateRefName,
	}
	pkgRefName, err := c.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	c.pkgMetadata.Name = pkgRefName
	c.pkgMetadata.Spec.DisplayName = strings.Split(pkgRefName, ".")[0]

	shortDesc := c.pkgMetadata.Spec.ShortDescription
	if len(shortDesc) == 0 {
		c.pkgMetadata.Spec.ShortDescription = pkgRefName
	}

	longDesc := c.pkgMetadata.Spec.LongDescription
	if len(longDesc) == 0 {
		c.pkgMetadata.Spec.LongDescription = pkgRefName
	}

	c.pkg.Spec.RefName = pkgRefName
	err = c.Save()
	if err != nil {
		return err
	}
	return nil
}

func (c *CreateStep) printPackageReferenceNameBlock() {
	c.ui.PrintInformationalText("A package reference name must be a valid DNS subdomain name (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names) \n - at least three segments separated by a '.', no trailing '.' e.g. samplepackage.corp.com")
}

func (c *CreateStep) PostInteract() error {
	c.ui.PrintHeaderText("\nOutput (Step 3/3)")
	buildObjectMeta := &metav1.ObjectMeta{
		Name: c.pkg.Spec.RefName,
	}
	c.build.SetObjectMeta(buildObjectMeta)
	c.updatePackageInstall()
	c.updatePackage()

	err := c.Save()
	if err != nil {
		return err
	}
	c.ui.PrintInformationalText("Successfully updated package-build.yml\n")
	c.ui.PrintInformationalText("Successfully updated package-resources.yml\n")
	c.printNextStep()
	return nil
}

func (c CreateStep) updatePackageInstall() {

	existingPkgInstall := c.pkgInstall
	if existingPkgInstall.ObjectMeta.Annotations == nil {
		existingPkgInstall.ObjectMeta.Annotations = make(map[string]string)
		existingPkgInstall.ObjectMeta.Annotations[appinit.LocalFetchAnnotationKey] = "."
	}

	if len(existingPkgInstall.Name) == 0 {
		existingPkgInstall.Name = c.pkgMetadata.Spec.DisplayName
	}

	if len(existingPkgInstall.Spec.ServiceAccountName) == 0 {
		existingPkgInstall.Spec.ServiceAccountName = c.pkgMetadata.Spec.DisplayName + "-sa"
	}
	if existingPkgInstall.Spec.PackageRef == nil {
		existingPkgInstall.Spec.PackageRef = &v1alpha12.PackageRef{
			RefName: c.pkg.Spec.RefName,
		}
	}
	// TODO Check whether we should add version constraint as well.
	if len(existingPkgInstall.Spec.PackageRef.RefName) == 0 {
		existingPkgInstall.Spec.PackageRef.RefName = c.pkg.Spec.RefName
	}
}

func (c CreateStep) updatePackage() {
	existingPkg := c.pkg

	if len(existingPkg.Spec.Version) == 0 {
		existingPkg.Spec.Version = defaultPkgVersion
	}
	existingPkg.Name = existingPkg.Spec.RefName + "." + existingPkg.Spec.Version

	if existingPkg.Spec.Template.Spec == nil {
		existingPkg.Spec.Template.Spec = &v1alpha13.AppSpec{}
	}
	if len(existingPkg.Spec.Template.Spec.Template) == 0 {
		existingPkg.Spec.Template.Spec.Template = c.build.GetAppSpec().Template
	}

	if len(existingPkg.Spec.Template.Spec.Deploy) == 0 {
		existingPkg.Spec.Template.Spec.Deploy = c.build.GetAppSpec().Deploy
	}
}

func (c CreateStep) printNextStep() {
	c.ui.PrintHeaderText("\n**Next steps**")
	c.ui.PrintInformationalText("Created files can be consumed in following ways:\n1. Optionally, use 'kctrl dev deploy' to iterate on the package and deploy locally.\n2. Use 'kctrl pkg release' to release the package.\n3. Use 'kctrl pkg release --repo-output repo/' to release and add package to package repository.\n")
}

func (c CreateStep) getDefaultPackageRefName() string {
	if len(c.pkgMetadata.Name) != 0 {
		return c.pkgMetadata.Name
	}
	return defaultPkgRefName
}

// Save method will save all the resources i.e. PackageBuild, PackageInstall, Package and PackageMetadata
func (c CreateStep) Save() error {
	// Save PackageBuild
	err := c.build.Save()
	if err != nil {
		return err
	}

	// Save Package
	content, err := yaml.Marshal(c.pkg)
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
	content, err = yaml.Marshal(c.pkgMetadata)
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
	content, err = yaml.Marshal(c.pkgInstall)
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
