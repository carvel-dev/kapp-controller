// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	appinit "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init"
	cmdapprelease "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/release"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	cmdpkgbuild "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/init"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcdatav1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
)

type ReleaseOptions struct {
	ui          cmdcore.AuthoringUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	pkgVersion         string
	chdir              string
	outputLocation     string
	repoOutputLocation string
	debug              bool
}

const (
	defaultArtifactDir = "carvel-artifacts"
	defaultVersion     = "0.0.0+build.%d"
)

func NewReleaseOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *ReleaseOptions {
	return &ReleaseOptions{ui: cmdcore.NewAuthoringUIImpl(ui), depsFactory: depsFactory, logger: logger}
}

func NewReleaseCmd(o *ReleaseOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release package (experimental)",
		RunE:  func(cmd *cobra.Command, args []string) error { return o.Run() },
		Annotations: map[string]string{
			cmdcore.PackageAuthoringCommandsHelpGroup.Key: cmdcore.PackageAuthoringCommandsHelpGroup.Value,
		},
	}

	cmd.Flags().StringVarP(&o.pkgVersion, "version", "v", "", "Version to be released")
	cmd.Flags().StringVar(&o.chdir, "chdir", "", "Location of the working directory")
	cmd.Flags().StringVar(&o.outputLocation, "copy-to", defaultArtifactDir, "Output location for artifacts")
	cmd.Flags().StringVar(&o.repoOutputLocation, "repo-output", "", "Output location for artifacts in repository bundle format")
	cmd.Flags().BoolVar(&o.debug, "debug", false, "Print verbose debug output")

	return cmd
}

func (o *ReleaseOptions) Run() error {
	if o.pkgVersion == "" {
		o.pkgVersion = fmt.Sprintf(defaultVersion, time.Now().Unix())
	}

	if o.chdir != "" {
		err := os.Chdir(o.chdir)
		if err != nil {
			return err
		}
	}

	pkgBuild, err := cmdpkgbuild.NewPackageBuildFromFile("package-build.yml")
	if err != nil {
		return err
	}

	pkgConfigs, err := local.NewConfigFromFiles([]string{"package-resources.yml"})
	if err != nil {
		return err
	}
	if len(pkgConfigs.PkgMetadatas) != 1 || len(pkgConfigs.PkgMetadatas) != 1 {
		return fmt.Errorf("Reading package-resource.yml: file malformed. (hint: delete the file and run `kctrl package init` again)")
	}

	o.printPrerequisites()

	err = o.loadExportData(pkgBuild)
	if err != nil {
		return err
	}

	// To be removed and moved to a question in case we have more config/variations around this
	if pkgBuild.Spec.Release == nil || len(pkgBuild.Spec.Release) == 0 {
		pkgBuild.Spec.Release = []appinit.Release{{Resource: &appinit.ReleaseResource{}}}
		err = pkgBuild.Save()
		if err != nil {
			return err
		}
	}

	buildAppSpec := pkgBuild.GetAppSpec()
	if buildAppSpec == nil {
		return fmt.Errorf("Releasing package: `kctrl pkg init` was not run successfully. (hint: re-run the `init` command)")
	}
	builderOpts := cmdapprelease.AppSpecBuilderOpts{
		BuildTemplate: buildAppSpec.Template,
		BuildDeploy:   buildAppSpec.Deploy,
		BuildExport:   *pkgBuild.GetExport(),
		Debug:         o.debug,
	}
	appSpec, err := cmdapprelease.NewAppSpecBuilder(o.depsFactory, o.logger, o.ui, builderOpts).Build()
	if err != nil {
		return err
	}

	for _, release := range pkgBuild.Spec.Release {
		switch {
		case release.Resource != nil:
			err = o.releaseResources(appSpec, *pkgBuild, &pkgConfigs.Pkgs[0], &pkgConfigs.PkgMetadatas[0])
			if err != nil {
				return nil
			}
		}
	}

	o.printNextSteps()
	return nil
}

func (o *ReleaseOptions) releaseResources(appSpec kcv1alpha1.AppSpec, pkgBuild cmdpkgbuild.PackageBuild,
	packageTemplate *kcdatav1alpha1.Package, metadataTemplate *kcdatav1alpha1.PackageMetadata) error {
	var yttPaths []string
	for _, templateStage := range pkgBuild.Spec.Template.Spec.App.Spec.Template {
		if templateStage.Ytt != nil {
			yttPaths = append(yttPaths, templateStage.Ytt.Paths...)
		}
	}
	valuesSchema, err := NewValuesSchemaGen(yttPaths).Schema()
	if err != nil {
		return err
	}

	artifactWriter := NewArtifactWriter(pkgBuild.Name, o.pkgVersion, packageTemplate, metadataTemplate, o.outputLocation, o.ui)
	err = artifactWriter.Write(&appSpec, *valuesSchema)
	if err != nil {
		return err
	}

	if o.repoOutputLocation != "" {
		err = artifactWriter.WriteRepoOutput(&appSpec, *valuesSchema, o.repoOutputLocation)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ReleaseOptions) loadExportData(pkgBuild *cmdpkgbuild.PackageBuild) error {
	if len(pkgBuild.Spec.Template.Spec.Export) == 0 {
		pkgBuild.Spec.Template.Spec.Export = []appinit.Export{
			{
				ImgpkgBundle: &appinit.ImgpkgBundle{UseKbldImagesLock: true},
			},
		}
	}
	if pkgBuild.Spec.Template.Spec.Export[0].ImgpkgBundle == nil {
		pkgBuild.Spec.Template.Spec.Export[0].ImgpkgBundle = &appinit.ImgpkgBundle{UseKbldImagesLock: true}
	}
	defaultImgValue := pkgBuild.Spec.Template.Spec.Export[0].ImgpkgBundle.Image
	o.ui.PrintInformationalText("The bundle created needs to be pushed to an OCI registry." +
		" Registry URL format: <REGISTRY_URL/REPOSITORY_NAME> e.g. index.docker.io/k8slt/sample-bundle")
	textOpts := ui.TextOpts{
		Label:        "Enter the registry URL",
		Default:      defaultImgValue,
		ValidateFunc: nil,
	}
	imgValue, err := o.ui.AskForText(textOpts)
	if err != nil {
		return err
	}
	pkgBuild.Spec.Template.Spec.Export[0].ImgpkgBundle.Image = strings.TrimSpace(imgValue)
	return pkgBuild.Save()
}

func (o *ReleaseOptions) printPrerequisites() {
	o.ui.PrintHeaderText("\nPrerequisites")
	o.ui.PrintInformationalText("1. The host must be authorized to push images to a registry (can be set up by running `docker login`)\n" +
		"2. Package can be released only after `kctrl package init` has been run successfully.\n\n")
}

func (o *ReleaseOptions) printNextSteps() {
	o.ui.PrintHeaderText("\nNext steps")
	o.ui.PrintInformationalText("1. The artifacts generated by the `--repo-output` flag can be bundled into a PackageRepository by using the `kctrl package repo release` command.\n" +
		"2. Generated Package and PackageMetadata manifests can be applied to the cluster directly.")
}
