// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

import "github.com/spf13/cobra"

type YttOverlayFlags struct {
	yttOverlays     bool
	yttOverlayFiles []string
}

func (s *YttOverlayFlags) Set(cmd *cobra.Command) {
	cmd.Flags().StringSliceVar(&s.yttOverlayFiles, "ytt-overlay", nil, "")
	cmd.Flags().BoolVar(&s.yttOverlays, "ytt-overlays", true, "Add or keep ytt overlays")
}
