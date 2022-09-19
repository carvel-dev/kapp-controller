// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/cppforlife/cobrautil"
	uierrs "github.com/cppforlife/go-cli-ui/errors"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd"

	// Import to initialize client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	err := nonExitingMain()
	if err != nil {
		os.Exit(1)
	}
}

// nonExitingMain does not use os.Exit to make sure Go runs defers
func nonExitingMain() error {
	rand.Seed(time.Now().UTC().UnixNano())

	// TODO logs
	// TODO log flags used

	confUI := ui.NewConfUI(ui.NewNoopLogger())
	defer confUI.Flush()

	command := cmd.NewDefaultKctrlCmd(confUI)

	err := command.Execute()
	if err != nil {
		confUI.ErrorLinef("kctrl: Error: %v", uierrs.NewMultiLineError(err))
		return err
	}

	if !cobrautil.IsCobraInternalCommand(os.Args) {
		confUI.PrintLinef("Succeeded")
	}
	return nil
}
