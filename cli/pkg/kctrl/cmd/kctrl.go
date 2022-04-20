// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"io"

	"github.com/cppforlife/cobrautil"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	cmdpkg "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package"
	pkgavail "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/available"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder"
	pkginst "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/installed"
	pkgrepo "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/repository"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/version"
)

type KctrlOptions struct {
	ui            *ui.ConfUI
	logger        *logger.UILogger
	configFactory cmdcore.ConfigFactory
	depsFactory   cmdcore.DepsFactory

	UIFlags         UIFlags
	LoggerFlags     LoggerFlags
	KubeAPIFlags    cmdcore.KubeAPIFlags
	KubeconfigFlags cmdcore.KubeconfigFlags
}

func NewKctrlOptions(ui *ui.ConfUI, configFactory cmdcore.ConfigFactory,
	depsFactory cmdcore.DepsFactory) *KctrlOptions {

	return &KctrlOptions{ui: ui, logger: logger.NewUILogger(ui),
		configFactory: configFactory, depsFactory: depsFactory}
}

func NewDefaultKctrlCmd(ui *ui.ConfUI) *cobra.Command {
	configFactory := cmdcore.NewConfigFactoryImpl()
	depsFactory := cmdcore.NewDepsFactoryImpl(configFactory, ui)
	options := NewKctrlOptions(ui, configFactory, depsFactory)
	flagsFactory := cmdcore.NewFlagsFactory(configFactory, depsFactory)
	return NewKctrlCmd(options, flagsFactory)
}

func NewKctrlCmd(o *KctrlOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kctrl",
		Short: "kctrl helps to manage packages and repositories on your Kubernetes cluster",

		RunE: cobrautil.ShowHelp,

		// Affects children as well
		SilenceErrors: true,
		SilenceUsage:  true,

		// Disable docs header
		DisableAutoGenTag: true,
		Version:           version.Version,

		// TODO bash completion
	}

	cmd.SetOutput(uiBlockWriter{o.ui}) // setting output for cmd.Help()

	cmd.SetUsageTemplate(cobrautil.HelpSectionsUsageTemplate([]cobrautil.HelpSection{
		cmdcore.PackageHelpGroup,
		cmdcore.RestOfCommandsHelpGroup,
	}))

	pkgOpts := cmdcore.PackageCommandTreeOpts{BinaryName: "kctrl", PositionalArgs: false, Color: true, JSON: true}

	SetGlobalFlags(o, cmd, flagsFactory, pkgOpts)

	ConfigurePathResolvers(o, cmd, flagsFactory)

	cmd.AddCommand(NewVersionCmd(NewVersionOptions(o.ui, o.depsFactory), flagsFactory))

	pkgCmd := cmdpkg.NewCmd()
	AddPackageCommands(o, pkgCmd, flagsFactory, pkgOpts)

	cmd.AddCommand(pkgCmd)

	appCmd := app.NewCmd()
	appCmd.AddCommand(app.NewGetCmd(app.NewGetOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	appCmd.AddCommand(app.NewListCmd(app.NewListOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	appCmd.AddCommand(app.NewStatusCmd(app.NewStatusOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	appCmd.AddCommand(app.NewPauseCmd(app.NewPauseOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	appCmd.AddCommand(app.NewKickCmd(app.NewKickOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	appCmd.AddCommand(app.NewDeleteCmd(app.NewDeleteOptions(o.ui, o.depsFactory, o.logger), flagsFactory))

	cmd.AddCommand(appCmd)

	ConfigureGlobalFlags(o, cmd, flagsFactory, pkgOpts.PositionalArgs)

	cmd.AddCommand(NewCmdCompletion())

	return cmd
}

func SetGlobalFlags(o *KctrlOptions, cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory, opts cmdcore.PackageCommandTreeOpts) {
	o.UIFlags.Set(cmd, flagsFactory, opts)
	o.LoggerFlags.Set(cmd, flagsFactory)
	o.KubeAPIFlags.Set(cmd, flagsFactory)
	o.KubeconfigFlags.Set(cmd, flagsFactory, opts)
}

func ConfigurePathResolvers(o *KctrlOptions, cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	o.configFactory.ConfigurePathResolver(o.KubeconfigFlags.Path.Value)
	o.configFactory.ConfigureContextResolver(o.KubeconfigFlags.Context.Value)
	o.configFactory.ConfigureYAMLResolver(o.KubeconfigFlags.YAML.Value)
}

func ConfigureGlobalFlags(o *KctrlOptions, cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory, positionalNameArg bool) {
	finishDebugLog := func(cmd *cobra.Command) {
		origRunE := cmd.RunE
		if origRunE != nil {
			cmd.RunE = func(cmd2 *cobra.Command, args []string) error {
				defer o.logger.DebugFunc("CommandRun").Finish()
				return origRunE(cmd2, args)
			}
		}
	}

	configureGlobal := cobrautil.WrapRunEForCmd(func(*cobra.Command, []string) error {
		o.UIFlags.ConfigureUI(o.ui)
		o.LoggerFlags.Configure(o.logger)
		o.KubeAPIFlags.Configure(o.configFactory)
		return nil
	})

	// Last one runs first
	// TODO: Add validation for number of arguments when positionalNameArg is true
	if positionalNameArg {
		cobrautil.VisitCommands(cmd, finishDebugLog, cobrautil.ReconfigureCmdWithSubcmd,
			configureGlobal, cobrautil.WrapRunEForCmd(cobrautil.ResolveFlagsForCmd))
	} else {
		cobrautil.VisitCommands(cmd, finishDebugLog, cobrautil.ReconfigureCmdWithSubcmd,
			cobrautil.ReconfigureLeafCmds(cobrautil.DisallowExtraArgs), configureGlobal, cobrautil.WrapRunEForCmd(cobrautil.ResolveFlagsForCmd))
	}
}

func AddPackageCommands(o *KctrlOptions, cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory, opts cmdcore.PackageCommandTreeOpts) {
	pkgrepoCmd := pkgrepo.NewCmd()
	pkgrepoCmd.AddCommand(pkgrepo.NewListCmd(pkgrepo.NewListOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgrepoCmd.AddCommand(pkgrepo.NewGetCmd(pkgrepo.NewGetOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgrepoCmd.AddCommand(pkgrepo.NewDeleteCmd(pkgrepo.NewDeleteOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgrepoCmd.AddCommand(pkgrepo.NewAddCmd(pkgrepo.NewAddOrUpdateOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgrepoCmd.AddCommand(pkgrepo.NewUpdateCmd(pkgrepo.NewAddOrUpdateOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))

	pkgiCmd := pkginst.NewCmd()
	pkgiCmd.AddCommand(pkginst.NewListCmd(pkginst.NewListOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgiCmd.AddCommand(pkginst.NewGetCmd(pkginst.NewGetOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgiCmd.AddCommand(pkginst.NewCreateCmd(pkginst.NewCreateOrUpdateOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgiCmd.AddCommand(pkginst.NewUpdateCmd(pkginst.NewCreateOrUpdateOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgiCmd.AddCommand(pkginst.NewDeleteCmd(pkginst.NewDeleteOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgiCmd.AddCommand(pkginst.NewPauseCmd(pkginst.NewPauseOrKickOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgiCmd.AddCommand(pkginst.NewKickCmd(pkginst.NewPauseOrKickOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgiCmd.AddCommand(pkginst.NewStatusCmd(pkginst.NewStatusOptions(o.ui, o.depsFactory, o.logger), flagsFactory))

	pkgaCmd := pkgavail.NewCmd()
	pkgaCmd.AddCommand(pkgavail.NewListCmd(pkgavail.NewListOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	pkgaCmd.AddCommand(pkgavail.NewGetCmd(pkgavail.NewGetOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))

	pkgbCmd := pkgbuilder.NewCmd()
	pkgbCmd.AddCommand(pkgbuilder.NewCreateCmd(pkgbuilder.NewCreateOptions(o.ui, o.logger, opts)))

	cmd.AddCommand(pkgrepoCmd)
	cmd.AddCommand(pkgiCmd)
	cmd.AddCommand(pkgaCmd)
	cmd.AddCommand(pkginst.NewInstallCmd(pkginst.NewCreateOrUpdateOptions(o.ui, o.depsFactory, o.logger, opts), flagsFactory))
	cmd.AddCommand(pkgbCmd)
}

func AttachGlobalFlags(o *KctrlOptions, cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory, opts cmdcore.PackageCommandTreeOpts) {
	SetGlobalFlags(o, cmd, flagsFactory, opts)
	ConfigurePathResolvers(o, cmd, flagsFactory)
	ConfigureGlobalFlags(o, cmd, flagsFactory, opts.PositionalArgs)
}

func AttachKctrlPackageCommandTree(cmd *cobra.Command, confUI *ui.ConfUI, opts cmdcore.PackageCommandTreeOpts) {
	configFactory := cmdcore.NewConfigFactoryImpl()
	depsFactory := cmdcore.NewDepsFactoryImpl(configFactory, confUI)
	options := NewKctrlOptions(confUI, configFactory, depsFactory)
	flagsFactory := cmdcore.NewFlagsFactory(configFactory, depsFactory)

	AddPackageCommands(options, cmd, flagsFactory, opts)
	AttachGlobalFlags(options, cmd, flagsFactory, opts)
}

type uiBlockWriter struct {
	ui ui.UI
}

var _ io.Writer = uiBlockWriter{}

func (w uiBlockWriter) Write(p []byte) (n int, err error) {
	w.ui.PrintBlock(p)
	return len(p), nil
}
