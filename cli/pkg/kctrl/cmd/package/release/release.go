// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	appinit "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init"
	cmdapprelease "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/release"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	cmdpkgbuild "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/init"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
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
	lockOutputFolder   = ".imgpkg"
	defaultVersion     = "0.0.0+build.%d"
)

func NewReleaseOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *ReleaseOptions {
	return &ReleaseOptions{ui: cmdcore.NewAuthoringUIImpl(ui), depsFactory: depsFactory, logger: logger}
}

func NewReleaseCmd(o *ReleaseOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release package",
		RunE:  func(cmd *cobra.Command, args []string) error { return o.Run() },
	}

	cmd.Flags().StringVarP(&o.pkgVersion, "version", "v", "", "Version to be released")
	cmd.Flags().StringVar(&o.chdir, "chdir", "", "Working directory with package-build and other config")
	cmd.Flags().StringVar(&o.outputLocation, "copy-to", defaultArtifactDir, "Output location for artifacts")
	cmd.Flags().StringVar(&o.repoOutputLocation, "repo-output", "", "Output location for artifacts in repository bundle format")
	cmd.Flags().BoolVar(&o.debug, "debug", false, "Print verbose debug output")

	return cmd
}

func (o *ReleaseOptions) Run() error {
	if o.pkgVersion == "" {
		o.pkgVersion = fmt.Sprintf("build-%d", time.Now().Unix())
	}

	pkgBuild, err := cmdpkgbuild.NewPackageBuildFromFile("package-build.yml")
	if err != nil {
		return err
	}
	pkg, err := cmdpkgbuild.GetPackage("package-resources.yml")
	if err != nil {
		return err
	}

	o.printPrerequisites()

	err = o.loadExportData(pkgBuild)
	if err != nil {
		return err
	}

	// To be removed and moved to a question in case we have more config/variations around this
	if pkgBuild.Spec.Release == nil || len(pkgBuild.Spec.Release) == 0 {
		pkgBuild.Spec.Release = []appinit.Release{
			{
				Resource: &appinit.ReleaseResource{},
			},
		}
		err = pkgBuild.Save()
		if err != nil {
			return err
		}
	}

	builderOpts := cmdapprelease.AppSpecBuilderOpts{
		BuildTemplate: pkgBuild.GetAppSpec().Template,
		BuildDeploy:   pkgBuild.GetAppSpec().Deploy,
		BuildExport:   *pkgBuild.GetExport(),
		Debug:         o.debug,
	}
	appSpec, err := cmdapprelease.NewAppSpecBuilder(wd, o.depsFactory, o.logger, o.ui, builderOpts).Build()
	if err != nil {
		return err
	}

	for _, release := range pkgBuild.Spec.Release {
		switch {
		case release.Resource != nil:
			err = o.releaseResources(appSpec, *pkgBuild)
			if err != nil {
				return nil
			}
		}
	}

	o.printNextSteps()
	return nil
}

func (o *ReleaseOptions) releaseResources(appSpec kcv1alpha1.AppSpec, pkgBuild cmdpkgbuild.PackageBuild) error {
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

	artifactWriter := NewArtifactWriter(pkgBuild.Name, o.pkgVersion, o.outputLocation, o.ui)
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
				ImgpkgBundle: &appinit.ImgpkgBundle{
					UseKbldImagesLock: true,
				},
			},
		}
	}
	if pkgBuild.Spec.Template.Spec.Export[0].ImgpkgBundle == nil {
		pkgBuild.Spec.Template.Spec.Export[0].ImgpkgBundle = &appinit.ImgpkgBundle{
			UseKbldImagesLock: true,
		}
	}
	defaultImgValue := pkgBuild.Spec.Template.Spec.Export[0].ImgpkgBundle.Image
	o.ui.PrintInformationalText("The bundle created needs to be pushed to an OCI registry. Registry URL format: <REGISTRY_URL/REPOSITORY_NAME> e.g. index.docker.io/k8slt/sample-bundle")
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
	o.ui.PrintHeaderText("Pre-requisites")
	o.ui.PrintInformationalText("1. The host must be authorized to push images to a registry (can be set up by running `docker login`)\n" +
		"2. Package can be released with this command only once `kctrl package init` has been run successfully.\n")
}

func (o *ReleaseOptions) printNextSteps() {
	o.ui.PrintHeaderText("\nNext steps")
	o.ui.PrintInformationalText("1. The artifacts generated by the `--repo-output` flag can be bundled into a repository using the `kctrl package repo release` comand.\n" +
		"2. Package and PackageMetadata YAML generated can be applied to the cluster directly so that it can be installed.")
}
