// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config

type ProxyOpts struct {
	HTTPProxy  string
	HTTPsProxy string
	NoProxy    string
}

type OSConfig interface {
	ApplyCACerts(string) error
	ApplyProxy(ProxyOpts) error
}
