// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"
	"os"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	appinit "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init"
	interfaces "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	cmdlocal "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"sigs.k8s.io/yaml"
)

const (
	pkgBuildFileName     = "package-build.yml"
	pkgResourcesFileName = "package-resources.yml"
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

	pkgBuild, err := o.newOrExistingPackageBuild()
	if err != nil {
		return err
	}

	pkg, pkgMetadata, pkgInstall, err := o.newOrExistingPackageResources()

	o.ui.PrintHeaderText("Prerequisites")
	o.ui.PrintInformationalText(`Welcome! Before we start, please ensure the following pre-requisites are met:
- The Carvel suite of tools are installed and you are familiar with them: ytt, imgpkg, vendir, and kbld.
`)

	o.ui.PrintHeaderText("\nBasic Information (Step 1/3)")

	pkgRefName, err := o.readPackageRefName(pkgMetadata.Name)
	if err != nil {
		return err
	}

	// TODO: @praveenrewar Can this be don any better?
	pkgMetadata.Name = pkgRefName
	pkg.Spec.RefName = pkgRefName
	pkgMetadata.Spec.DisplayName = strings.Split(pkgRefName, ".")[0]

	shortDesc := pkgMetadata.Spec.ShortDescription
	if len(shortDesc) == 0 {
		pkgMetadata.Spec.ShortDescription = pkgRefName
	}

	longDesc := pkgMetadata.Spec.LongDescription
	if len(longDesc) == 0 {
		pkgMetadata.Spec.LongDescription = pkgRefName
	}

	err = o.SavePackageResources(pkg, pkgMetadata, pkgInstall)
	if err != nil {
		return err
	}
	pkgBuild.Save()
	if err != nil {
		return err
	}

	// TODO: @praveenrewar Remove the step part and use only relevant code from Fetch
	appCreateStep := appinit.NewCreateStep(o.ui, pkgBuild, o.logger, o.depsFactory, false)
	err = interfaces.Run(appCreateStep)
	if err != nil {
		return err
	}

	pkgBuild.Save()
	if err != nil {
		return err
	}

	// TODO: @praveenrewar Remove the header text as we are only saving data to files and doesn't require user input
	// hence shouldn't be called as Step 3/3
	o.ui.PrintHeaderText("\nOutput (Step 3/3)")

	pkgBuild.SetObjectMeta(&metav1.ObjectMeta{
		Name: pkg.Spec.RefName,
	})

	// TODO: @praveeenrewar Refactor the updates
	o.updatePackageInstall(pkgInstall, pkg.Spec.RefName, pkgMetadata.Spec.DisplayName)
	o.updatePackage(pkg, pkgBuild)

	err = o.SavePackageResources(pkg, pkgMetadata, pkgInstall)
	if err != nil {
		return err
	}
	pkgBuild.Save()
	if err != nil {
		return err
	}

	o.ui.PrintInformationalText("Successfully updated package-build.yml\n")
	o.ui.PrintInformationalText("Successfully updated package-resources.yml\n")
	o.ui.PrintHeaderText("\n**Next steps**")
	o.ui.PrintInformationalText(`Created files can be consumed in following ways:
1. Optionally, use 'kctrl dev deploy' to iterate on the package and deploy locally.
2. Use 'kctrl pkg release' to release the package.
3. Use 'kctrl pkg release --repo-output repo/' to release and add package to package repository.
`)
	return nil
}

func (o *InitOptions) newOrExistingPackageBuild() (*PackageBuild, error) {
	content, err := os.ReadFile(pkgBuildFileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return &PackageBuild{}, err
		}
		return &PackageBuild{TypeMeta: metav1.TypeMeta{
			Kind:       "PackageBuild",
			APIVersion: "kctrl.carvel.dev/v1alpha1",
		}}, nil
	}

	var packageBuild PackageBuild
	err = yaml.Unmarshal(content, &packageBuild)
	return &packageBuild, err
}

func (o *InitOptions) newOrExistingPackageResources() (*v1alpha1.Package,
	*v1alpha1.PackageMetadata, *pkgv1alpha1.PackageInstall, error) {
	var configs cmdlocal.Configs
	configs, err := cmdlocal.NewConfigFromFiles([]string{pkgResourcesFileName})
	if err != nil {
		if os.IsNotExist(err) {
			return &v1alpha1.Package{}, &v1alpha1.PackageMetadata{}, &pkgv1alpha1.PackageInstall{}, err
		}
	}

	pkg := &v1alpha1.Package{TypeMeta: metav1.TypeMeta{
		Kind:       "Package",
		APIVersion: "data.packaging.carvel.dev/v1alpha1",
	}}

	if configs.Pkgs != nil {
		if len(configs.Pkgs) > 1 {
			return &v1alpha1.Package{}, &v1alpha1.PackageMetadata{}, &pkgv1alpha1.PackageInstall{}, fmt.Errorf("More than 1 Package found")
		}
		pkg = &configs.Pkgs[0]
	}

	pkgMetadata := &v1alpha1.PackageMetadata{TypeMeta: metav1.TypeMeta{
		Kind:       "PackageMetadata",
		APIVersion: "data.packaging.carvel.dev/v1alpha1",
	}}

	if configs.PkgMetadatas != nil {
		if len(configs.PkgMetadatas) > 1 {
			return &v1alpha1.Package{}, &v1alpha1.PackageMetadata{}, &pkgv1alpha1.PackageInstall{}, fmt.Errorf("More than 1 PackageMetadata found")
		}
		pkgMetadata = &configs.PkgMetadatas[0]
	}

	pkgInstall := &pkgv1alpha1.PackageInstall{TypeMeta: metav1.TypeMeta{
		Kind:       "PackageInstall",
		APIVersion: "packaging.carvel.dev/v1alpha1",
	}}

	// TODO we get an error if package-resources.yml file exist but there is no packageInstall in it.
	// Probably, needs to make changes to local Package and adopt them in dev deploy.
	if configs.PkgInstalls != nil {
		if len(configs.PkgInstalls) > 1 {
			return &v1alpha1.Package{}, &v1alpha1.PackageMetadata{}, &pkgv1alpha1.PackageInstall{}, fmt.Errorf("More than 1 PackageInstall found")
		}
		pkgInstall = &configs.PkgInstalls[0]
	}

	return pkg, pkgMetadata, pkgInstall, nil
}

func (o *InitOptions) readPackageRefName(packageMetadataName string) (string, error) {
	o.ui.PrintInformationalText(`A package reference name must be a valid DNS subdomain name (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
 - at least three segments separated by a '.', no trailing '.' e.g. samplepackage.corp.com`)

	defaultPkgRefName := "samplepackage.corp.com"
	if len(packageMetadataName) > 0 {
		defaultPkgRefName = packageMetadataName
	}

	return o.ui.AskForText(ui.TextOpts{
		Label:        "Enter the package reference name",
		Default:      defaultPkgRefName,
		ValidateFunc: o.validateRefName,
	})
}

func (o *InitOptions) updatePackageInstall(pkgInstall *pkgv1alpha1.PackageInstall, refName, displayName string) {
	if pkgInstall.ObjectMeta.Annotations == nil {
		pkgInstall.ObjectMeta.Annotations = make(map[string]string)
		pkgInstall.ObjectMeta.Annotations[appinit.LocalFetchAnnotationKey] = "."
	}

	if len(pkgInstall.Name) == 0 {
		pkgInstall.Name = displayName
	}

	if len(pkgInstall.Spec.ServiceAccountName) == 0 {
		pkgInstall.Spec.ServiceAccountName = displayName + "-sa"
	}
	if pkgInstall.Spec.PackageRef == nil {
		pkgInstall.Spec.PackageRef = &pkgv1alpha1.PackageRef{
			RefName: refName,
		}
	}
	// TODO Check whether we should add version constraint as well.
	if len(pkgInstall.Spec.PackageRef.RefName) == 0 {
		pkgInstall.Spec.PackageRef.RefName = refName
	}
}

func (o *InitOptions) updatePackage(pkg *v1alpha1.Package, pkgBuild *PackageBuild) {
	if len(pkg.Spec.Version) == 0 {
		pkg.Spec.Version = "0.0.0"
	}
	pkg.Name = pkg.Spec.RefName + "." + pkg.Spec.Version

	if pkg.Spec.Template.Spec == nil {
		pkg.Spec.Template.Spec = &kcv1alpha1.AppSpec{}
	}
	if len(pkg.Spec.Template.Spec.Template) == 0 {
		pkg.Spec.Template.Spec.Template = pkgBuild.GetAppSpec().Template
	}

	if len(pkg.Spec.Template.Spec.Deploy) == 0 {
		pkg.Spec.Template.Spec.Deploy = pkgBuild.GetAppSpec().Deploy
	}
}

func (o *InitOptions) SavePackageResources(pkg *v1alpha1.Package,
	pkgMetadata *v1alpha1.PackageMetadata, pkgInstall *pkgv1alpha1.PackageInstall) error {
	pkgYAML, err := yaml.Marshal(pkg)
	if err != nil {
		return err
	}

	pkgMetadataYAML, err := yaml.Marshal(pkgMetadata)
	if err != nil {
		return err
	}

	pkgInstallYAML, err := yaml.Marshal(pkgInstall)
	if err != nil {
		return err
	}

	packageResources := fmt.Sprintf(`%s
---
%s
---
%s`, string(pkgYAML), string(pkgMetadataYAML), string(pkgInstallYAML))

	return os.WriteFile(pkgResourcesFileName, []byte(packageResources), os.ModePerm)
}

// TODO should we use the same validation used in kapp controller. But that accepts other parameter. ValidatePackageMetadataName in validations.go file
func (o *InitOptions) validateRefName(name string) (bool, string, error) {
	if len(name) == 0 {
		return false, "Fully qualified name of a package cannot be empty", nil
	}
	if errs := validation.IsDNS1123Subdomain(name); len(errs) > 0 {
		return false, strings.Join(errs, ","), nil
	}
	if len(strings.Split(name, ".")) < 3 {
		return false, fmt.Sprintf("Invalid name: %s should be a fully qualified name with at least three segments separated by dots", name), nil
	}
	return true, "", nil
}
