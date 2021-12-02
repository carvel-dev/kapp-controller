// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"time"

	"github.com/spf13/cobra"
)

type WaitFlags struct {
	Enabled       bool
	CheckInterval time.Duration
	Timeout       time.Duration
}

func (f *WaitFlags) Set(cmd *cobra.Command, flagsFactory FlagsFactory) {
	cmd.Flags().BoolVar(&f.Enabled, "wait", true, "Wait for reconciliation to complete")
	cmd.Flags().DurationVar(&f.CheckInterval, "wait-check-interval", 1*time.Second, "Amount of time to sleep between checks while waiting")
	cmd.Flags().DurationVar(&f.Timeout, "wait-timeout", 5*time.Minute, "Maximum amount of time to wait in wait phase")
}
