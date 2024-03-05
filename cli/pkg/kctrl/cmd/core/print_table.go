// Copyright 2024 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
)

func PrintTable(ui ui.UI, table uitable.Table) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	ui.PrintTable(table)
	return nil
}
