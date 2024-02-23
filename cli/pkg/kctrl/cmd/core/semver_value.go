// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	versions "carvel.dev/vendir/pkg/vendir/versions"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
)

type ValueSemver struct {
	V string
}

func NewValueSemver(v string) ValueSemver {
	return ValueSemver{V: v}
}

func (t ValueSemver) String() string {
	return t.V
}

func (t ValueSemver) Value() uitable.Value { return t }
func (t ValueSemver) Compare(other uitable.Value) int {
	otherS, _ := versions.NewRelaxedSemver(other.(ValueSemver).V)
	tS, _ := versions.NewRelaxedSemver(t.V)
	switch {
	case tS.Version.EQ(otherS.Version):
		return 0
	case tS.Version.LT(otherS.Version):
		return -1
	default:
		return 1
	}
}
