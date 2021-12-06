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

type WaitFlagsOpts struct {
	AllowDisableWait bool
	DefaultInterval  time.Duration
	DefaultTimeout   time.Duration
}

func (f *WaitFlags) Set(cmd *cobra.Command, flagsFactory FlagsFactory, opts *WaitFlagsOpts) {
	if opts.AllowDisableWait {
		cmd.Flags().BoolVar(&f.Enabled, "wait", true, "Wait for reconciliation to complete")
	}
	f.Enabled = true
	cmd.Flags().DurationVar(&f.CheckInterval, "wait-check-interval", opts.DefaultInterval, "Amount of time to sleep between checks while waiting")
	cmd.Flags().DurationVar(&f.Timeout, "wait-timeout", opts.DefaultTimeout, "Maximum amount of time to wait in wait phase")
}
