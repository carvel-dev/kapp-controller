// Copyright 2024 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
)

func validateColumns(headers *[]uitable.Header, cols *[]string) error {
	invalidColumns := []string{}
	for _, col := range *cols {
		found := false
		for _, head := range *headers {
			if col == head.Key || col == head.Title {
				found = true
				break
			}
		}
		if !found {
			invalidColumns = append(invalidColumns, col)
		}
	}

	if len(invalidColumns) > 0 {
		return fmt.Errorf("invalid column names: %s", strings.Join(invalidColumns, ","))
	}
	return nil
}

func PrintTable(ui ui.UI, table uitable.Table, columns *[]string) error {
	err := validateColumns(&table.Header, columns)
	if err == nil {
		ui.PrintTable(table)
	}
	return err
}
