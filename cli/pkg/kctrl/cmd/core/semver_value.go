// Copyright 20210 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions"
)

type ValueSemver struct {
	V uitable.Value
}

func NewValueSemver(v uitable.Value) ValueSemver {
	return ValueSemver{V: v}
}

func (t ValueSemver) String() string {
	return t.V.String()
}

func (t ValueSemver) Value() uitable.Value { return t }
func (t ValueSemver) Compare(other uitable.Value) int {
	otherS, _ := versions.NewRelaxedSemver(other.String())
	tS, _ := versions.NewRelaxedSemver(t.String())
	switch {
	case tS.Version.EQ(otherS.Version):
		return 0
	case tS.Version.LT(otherS.Version):
		return -1
	default:
		return 1
	}
}
