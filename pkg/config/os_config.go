// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package config

// ProxyOpts contains proxy configuration for the system.
type ProxyOpts struct {
	HTTPProxy  string
	HTTPSProxy string
	NoProxy    string
}

// OSConfig configures runtime environment with necessary
// CA certificates and proxy configuration.
type OSConfig interface {
	ApplyCACerts(string) error
	ApplyProxy(ProxyOpts) error
}
