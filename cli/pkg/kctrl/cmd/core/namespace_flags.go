// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/cppforlife/cobrautil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type NamespaceFlags struct {
	Name string
}

func (s *NamespaceFlags) Set(cmd *cobra.Command, flagsFactory FlagsFactory) {
	name := flagsFactory.NewNamespaceNameFlag(&s.Name, "KCTRL_NAMESPACE")
	cmd.Flags().VarP(name, "namespace", "n", "Specified namespace ($KCTRL_NAMESPACE or default from kubeconfig)")
}

func (s *NamespaceFlags) SetWithPackageCommandTreeOpts(cmd *cobra.Command, flagsFactory FlagsFactory, opts PackageCommandTreeOpts) {
	namespaceEnvVariableKey := fmt.Sprintf("%s_NAMESPACE", strings.ToUpper(opts.BinaryName))
	name := flagsFactory.NewNamespaceNameFlag(&s.Name, namespaceEnvVariableKey)
	cmd.Flags().VarP(name, "namespace", "n", fmt.Sprintf("Specified namespace ($%s or default from kubeconfig)", namespaceEnvVariableKey))
}

type NamespaceNameFlag struct {
	value          *string
	configFactory  ConfigFactory
	envVariableKey string
}

var _ pflag.Value = &NamespaceNameFlag{}
var _ cobrautil.ResolvableFlag = &NamespaceNameFlag{}

func NewNamespaceNameFlag(value *string, configFactory ConfigFactory, envVariableKey string) *NamespaceNameFlag {
	return &NamespaceNameFlag{value, configFactory, envVariableKey}
}

func (s *NamespaceNameFlag) Set(val string) error {
	*s.value = val
	return nil
}

func (s *NamespaceNameFlag) Type() string   { return "string" }
func (s *NamespaceNameFlag) String() string { return "" } // default for usage

func (s *NamespaceNameFlag) Resolve() error {
	value, err := s.resolveValue()
	if err != nil {
		return err
	}

	*s.value = value

	return nil
}

func (s *NamespaceNameFlag) resolveValue() (string, error) {
	if s.value != nil && len(*s.value) > 0 {
		return *s.value, nil
	}

	envVal := os.Getenv(s.envVariableKey)
	if len(envVal) > 0 {
		return envVal, nil
	}

	configVal, err := s.configFactory.DefaultNamespace()
	if err != nil {
		return configVal, nil
	}

	if len(configVal) > 0 {
		return configVal, nil
	}

	return "", fmt.Errorf("Expected to non-empty namespace name")
}
