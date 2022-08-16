// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"os"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdappinit "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
)

type ReleaseOptions struct {
	ui          cmdcore.AuthoringUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	pkgVersion     string
	chdir          string
	outputLocation string
	debug          bool
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
		Short: "Release package",
		RunE:  func(cmd *cobra.Command, args []string) error { return o.Run() },
	}

	cmd.Flags().StringVarP(&o.pkgVersion, "version", "v", "", "Version to be released")
	cmd.Flags().StringVar(&o.chdir, "chdir", "", "Working directory with package-build and other config")
	cmd.Flags().StringVar(&o.outputLocation, "copy-to", defaultArtifactDir, "Output location for artifacts")
	cmd.Flags().BoolVar(&o.debug, "debug", false, "Print verbose debug output")

	return cmd
}

func (o *ReleaseOptions) Run() error {
	if o.chdir != "" {
		err := os.Chdir(o.chdir)
		if err != nil {
			return err
		}
	}

	o.printPrerequisites()

	appBuild, err := cmdappinit.NewAppBuildFromFile(cmdappinit.FileName)
	if err != nil {
		return err
	}

	builderOpts := AppSpecBuilderOpts{
		BuildTemplate: appBuild.GetAppSpec().Template,
		BuildDeploy:   appBuild.GetAppSpec().Deploy,
		BuildExport:   *appBuild.GetExport(),
		Debug:         o.debug,
	}
	_, err = NewAppSpecBuilder(o.depsFactory, o.logger, o.ui, builderOpts).Build()
	if err != nil {
		return err
	}

	// TODO: Write app resources

	return nil
}

func (o *ReleaseOptions) printPrerequisites() {
	o.ui.PrintHeaderText("Pre-requisites")
	o.ui.PrintInformationalText("1. The host must be authorized to push images to a registry (can be set up by running `docker login`)\n" +
		"2. an app can be released with this command only once `kctrl app init` has been run successfully.\n")
}
