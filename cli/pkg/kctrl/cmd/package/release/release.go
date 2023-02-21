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
	cmdapprelease "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/release"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	cmdpkg "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local"
	buildconfigs "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local/buildconfigs"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcdatav1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
)

type ReleaseOptions struct {
	ui          cmdcore.AuthoringUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	pkgVersion            string
	chdir                 string
	outputLocation        string
	repoOutputLocation    string
	debug                 bool
	generateOpenAPISchema bool
	tag                   string
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
	cmd.Flags().StringVarP(&o.tag, "tag", "t", "", "Tag pushed with imgpkg bundle (default build-<TIMESTAMP>)")
	cmd.Flags().BoolVar(&o.generateOpenAPISchema, "openapi-schema", true, "Generates openapi schema for ytt and helm templated files and adds it to generated package")

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

	pkgBuild, err := buildconfigs.NewPackageBuildFromFile(buildconfigs.PkgBuildFileName)
	if err != nil {
		return err
	}

	pkgConfigs, err := local.NewConfigFromFiles([]string{cmdpkg.PkgResourcesFileName})
	if err != nil {
		return err
	}
	if len(pkgConfigs.PkgMetadatas) != 1 || len(pkgConfigs.PkgMetadatas) != 1 {
		return fmt.Errorf("Reading package-resource.yml: file malformed. (hint: delete the file and run `package init` again)")
	}

	o.printPrerequisites()

	err = o.loadExportData(pkgBuild)
	if err != nil {
		return err
	}

	// To be removed and moved to a question in case we have more config/variations around this
	if pkgBuild.Spec.Release == nil || len(pkgBuild.Spec.Release) == 0 {
		pkgBuild.Spec.Release = []buildconfigs.Release{{Resource: &buildconfigs.ReleaseResource{}}}
		err = pkgBuild.Save()
		if err != nil {
			return err
		}
	}

	buildAppSpec := pkgBuild.GetAppSpec()
	if buildAppSpec == nil {
		return fmt.Errorf("Releasing package: 'package init' was not run successfully. (hint: re-run the 'init' command)")
	}
	builderOpts := cmdapprelease.AppSpecBuilderOpts{
		BuildTemplate: buildAppSpec.Template,
		BuildDeploy:   buildAppSpec.Deploy,
		BuildExport:   pkgBuild.GetExport(),
		Debug:         o.debug,
		BundleTag:     o.tag,
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

func (o *ReleaseOptions) releaseResources(appSpec kcv1alpha1.AppSpec, pkgBuild buildconfigs.PackageBuild,
	packageTemplate *kcdatav1alpha1.Package, metadataTemplate *kcdatav1alpha1.PackageMetadata) error {
	var valuesSchema *kcdatav1alpha1.ValuesSchema
	var err error
	if o.generateOpenAPISchema {
		valuesSchema, err = generateValuesSchema(pkgBuild)
		if err != nil {
			return err
		}
	}

	artifactWriter := NewArtifactWriter(pkgBuild.Name, o.pkgVersion, packageTemplate, metadataTemplate, o.outputLocation, o.ui)
	err = artifactWriter.Write(&appSpec, valuesSchema)
	if err != nil {
		return err
	}

	if o.repoOutputLocation != "" {
		err = artifactWriter.WriteRepoOutput(&appSpec, valuesSchema, o.repoOutputLocation)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateValuesSchema(pkgBuild buildconfigs.PackageBuild) (*kcdatav1alpha1.ValuesSchema, error) {
	if pkgBuild.Spec.Template.Spec.App.Spec.Template != nil {
		// As of today, PackageInstall values file is applicable only for the first templating step.
		// https://github.com/carvel-dev/kapp-controller/blob/develop/pkg/packageinstall/app.go#L103
		templateStage := pkgBuild.Spec.Template.Spec.App.Spec.Template[0]
		switch {
		case templateStage.HelmTemplate != nil:
			return NewHelmValuesSchemaGen(templateStage.HelmTemplate.Path).Schema()
		case templateStage.Ytt != nil:
			return NewValuesSchemaGen(templateStage.Ytt.Paths).Schema()
		}
	}
	return nil, nil
}

func (o *ReleaseOptions) loadExportData(pkgBuild *buildconfigs.PackageBuild) error {
	if len(pkgBuild.Spec.Template.Spec.Export) == 0 {
		pkgBuild.Spec.Template.Spec.Export = []buildconfigs.Export{
			{
				ImgpkgBundle: &buildconfigs.ImgpkgBundle{UseKbldImagesLock: true},
			},
		}
	}
	if pkgBuild.Spec.Template.Spec.Export[0].ImgpkgBundle == nil {
		pkgBuild.Spec.Template.Spec.Export[0].ImgpkgBundle = &buildconfigs.ImgpkgBundle{UseKbldImagesLock: true}
	}
	defaultImgValue := pkgBuild.Spec.Template.Spec.Export[0].ImgpkgBundle.Image
	o.ui.PrintInformationalText("The bundle created needs to be pushed to an OCI registry. (format: <REGISTRY_URL/REPOSITORY_NAME>) e.g. index.docker.io/k8slt/sample-bundle")
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
	o.ui.PrintInformationalText("1. Host is authorized to push images to a registry (can be set up by running `docker login`)\n" +
		"2. `package init` ran successfully.\n\n")
}

func (o *ReleaseOptions) printNextSteps() {
	o.ui.PrintHeaderText("\nNext steps")
	o.ui.PrintInformationalText("1. The artifacts generated by the `--repo-output` flag can be bundled into a PackageRepository by using the `package repository release` command.\n" +
		"2. Generated Package and PackageMetadata manifests can be applied to the cluster directly.\n")
}
