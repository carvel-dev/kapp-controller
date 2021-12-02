// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
)

type ValueTruncated struct {
	V   uitable.Value
	Max int
}

func NewValueTruncated(v uitable.Value, max int) ValueTruncated {
	return ValueTruncated{V: v, Max: max}
}

func (t ValueTruncated) String() string {
	str := t.V.String()
	if len(str) > t.Max {
		return str[:t.Max] + "..."
	}
	return str
}

func (t ValueTruncated) Value() uitable.Value            { return t.V }
func (t ValueTruncated) Compare(other uitable.Value) int { panic("Never called") }
