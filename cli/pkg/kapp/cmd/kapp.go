// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"io"

	"github.com/cppforlife/cobrautil"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/cmd/core"
	cmdpkg "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/cmd/package"
	pkgavail "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/cmd/package/available"
	pkginst "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/cmd/package/install"
	pkgrepo "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/cmd/package/repository"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/logger"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/version"
)

type KappOptions struct {
	ui            *ui.ConfUI
	logger        *logger.UILogger
	configFactory cmdcore.ConfigFactory
	depsFactory   cmdcore.DepsFactory

	UIFlags         UIFlags
	LoggerFlags     LoggerFlags
	KubeAPIFlags    cmdcore.KubeAPIFlags
	KubeconfigFlags cmdcore.KubeconfigFlags
}

func NewKappOptions(ui *ui.ConfUI, configFactory cmdcore.ConfigFactory,
	depsFactory cmdcore.DepsFactory) *KappOptions {

	return &KappOptions{ui: ui, logger: logger.NewUILogger(ui),
		configFactory: configFactory, depsFactory: depsFactory}
}

func NewDefaultKappCmd(ui *ui.ConfUI) *cobra.Command {
	configFactory := cmdcore.NewConfigFactoryImpl()
	depsFactory := cmdcore.NewDepsFactoryImpl(configFactory, ui)
	options := NewKappOptions(ui, configFactory, depsFactory)
	flagsFactory := cmdcore.NewFlagsFactory(configFactory, depsFactory)
	return NewKappCmd(options, flagsFactory)
}

func NewKappCmd(o *KappOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kapp",
		Short: "kapp helps to manage applications on your Kubernetes cluster",

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
		cmdcore.AppHelpGroup,
		cmdcore.AppSupportHelpGroup,
		cmdcore.MiscHelpGroup,
		cmdcore.RestOfCommandsHelpGroup,
	}))

	o.UIFlags.Set(cmd, flagsFactory)
	o.LoggerFlags.Set(cmd, flagsFactory)
	o.KubeAPIFlags.Set(cmd, flagsFactory)
	o.KubeconfigFlags.Set(cmd, flagsFactory)

	o.configFactory.ConfigurePathResolver(o.KubeconfigFlags.Path.Value)
	o.configFactory.ConfigureContextResolver(o.KubeconfigFlags.Context.Value)
	o.configFactory.ConfigureYAMLResolver(o.KubeconfigFlags.YAML.Value)

	cmd.AddCommand(NewVersionCmd(NewVersionOptions(o.ui), flagsFactory))

	pkgrepoCmd := pkgrepo.NewCmd()
	//pkgrepoCmd.AddCommand(pkgrepo.NewCmd())
	// appCmd.AddCommand(cmdtools.NewDiffCmd(cmdtools.NewDiffOptions(o.ui, o.depsFactory), flagsFactory))
	pkgrepoCmd.AddCommand(pkgrepo.NewListCmd(pkgrepo.NewListOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	pkgrepoCmd.AddCommand(pkgrepo.NewGetCmd(pkgrepo.NewGetOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	pkgrepoCmd.AddCommand(pkgrepo.NewDeleteCmd(pkgrepo.NewDeleteOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	pkgrepoCmd.AddCommand(pkgrepo.NewAddCmd(pkgrepo.NewAddOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	pkgrepoCmd.AddCommand(pkgrepo.NewUpdateCmd(pkgrepo.NewUpdateOptions(o.ui, o.depsFactory, o.logger), flagsFactory))

	pkgiCmd := pkginst.NewCmd()
	pkgiCmd.AddCommand(pkginst.NewListCmd(pkginst.NewListOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	pkgiCmd.AddCommand(pkginst.NewGetCmd(pkginst.NewGetOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	pkgiCmd.AddCommand(pkginst.NewCreateCmd(pkginst.NewCreateOrUpdateOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	pkgiCmd.AddCommand(pkginst.NewUpdateCmd(pkginst.NewCreateOrUpdateOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	pkgiCmd.AddCommand(pkginst.NewDeleteCmd(pkginst.NewDeleteOptions(o.ui, o.depsFactory, o.logger), flagsFactory))

	pkgaCmd := pkgavail.NewCmd()
	pkgaCmd.AddCommand(pkgavail.NewListCmd(pkgavail.NewListOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	pkgaCmd.AddCommand(pkgavail.NewGetCmd(pkgavail.NewGetOptions(o.ui, o.depsFactory, o.logger), flagsFactory))

	pkgCmd := cmdpkg.NewCmd()
	pkgCmd.AddCommand(pkgrepoCmd)
	pkgCmd.AddCommand(pkgiCmd)
	pkgCmd.AddCommand(pkgaCmd)
	pkgCmd.AddCommand(pkginst.NewInstallCmd(pkginst.NewCreateOrUpdateOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	// appCmd.AddCommand(cmdtools.NewDiffCmd(cmdtools.NewDiffOptions(o.ui, o.depsFactory), flagsFactory))
	// appCmd.AddCommand(cmdtools.NewListLabelsCmd(cmdtools.NewListLabelsOptions(o.ui, o.depsFactory, o.logger), flagsFactory))
	cmd.AddCommand(pkgCmd)

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
	cobrautil.VisitCommands(cmd, finishDebugLog, cobrautil.ReconfigureCmdWithSubcmd,
		cobrautil.ReconfigureLeafCmds(cobrautil.DisallowExtraArgs), configureGlobal, cobrautil.WrapRunEForCmd(cobrautil.ResolveFlagsForCmd))

	// Completion command have to be added after the VisitCommands
	// This due to the ReconfigureLeafCmds that we do not want to have enforced for the completion
	// This configurations forces all nodes to do not accept extra args, but the completion requires 1 extra arg
	cmd.AddCommand(NewCmdCompletion())
	return cmd
}

type uiBlockWriter struct {
	ui ui.UI
}

var _ io.Writer = uiBlockWriter{}

func (w uiBlockWriter) Write(p []byte) (n int, err error) {
	w.ui.PrintBlock(p)
	return len(p), nil
}
