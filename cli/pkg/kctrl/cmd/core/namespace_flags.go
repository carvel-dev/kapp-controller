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

var sharedNamespaces = []string{
	"default",
	"kube-public",
}

type NamespaceFlags struct {
	Name                    string
	AllowedSharedNamespaces []string
}

func (s *NamespaceFlags) Set(cmd *cobra.Command, flagsFactory FlagsFactory) {
	name := flagsFactory.NewNamespaceNameFlag(&s.Name, "KCTRL_NAMESPACE")
	cmd.Flags().VarP(name, "namespace", "n", "Specified namespace ($KCTRL_NAMESPACE or default from kubeconfig)")
}

func (s *NamespaceFlags) SetWithPackageCommandTreeOpts(cmd *cobra.Command, flagsFactory FlagsFactory, opts PackageCommandTreeOpts) {
	namespaceEnvVariableKey := fmt.Sprintf("%s_NAMESPACE", strings.ToUpper(opts.BinaryName))
	name := flagsFactory.NewNamespaceNameFlag(&s.Name, namespaceEnvVariableKey)
	cmd.Flags().VarP(name, "namespace", "n", fmt.Sprintf("Specified namespace ($%s or default from kubeconfig)", namespaceEnvVariableKey))
	cmd.Flags().StringArrayVar(&s.AllowedSharedNamespaces, "dangerous-allow-use-of-shared-namespace", []string{}, "Allow use of shared namespaces (optional, comma separated strings)")
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

func (s *NamespaceFlags) CheckForDisallowedSharedNamespaces() error {
	for _, ns := range sharedNamespaces {
		if s.Name == ns {
			for _, allowedNs := range s.AllowedSharedNamespaces {
				if ns == allowedNs {
					return nil
				}
			}
			return fmt.Errorf(`Creating sensitive resources in a shared namespace (%s)
			(hint: Specify a namespace using the '-n' flag or use kubeconfig to change default namespace 'kubectl config set-context --current --namespace=private-namespace'.
			Or use '--dangerous-allow-allow-use-of-shared-namespace=%s' to allow use of shared namespace)`, s.Name, s.Name)
		}
	}
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
