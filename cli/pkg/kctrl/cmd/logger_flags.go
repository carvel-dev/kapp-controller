// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
)

type LoggerFlags struct {
	Debug bool
}

func (f *LoggerFlags) Set(cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	cmd.PersistentFlags().BoolVar(&f.Debug, "debug", false, "Include debug output")
}

func (f *LoggerFlags) Configure(logger *logger.UILogger) {
	logger.SetDebug(f.Debug)
}
