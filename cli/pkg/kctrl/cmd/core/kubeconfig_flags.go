// Copyright 2024 The Carvel Authors.
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

type KubeconfigFlags struct {
	Path    *KubeconfigPathFlag
	Context *KubeconfigContextFlag
	YAML    *KubeconfigYAMLFlag
}

func (f *KubeconfigFlags) Set(cmd *cobra.Command, flagsFactory FlagsFactory, opts PackageCommandTreeOpts) {
	kubeconfigEnvVariableKey := fmt.Sprintf("%s_KUBECONFIG", strings.ToUpper(opts.BinaryName))
	f.Path = NewKubeconfigPathFlag(kubeconfigEnvVariableKey)
	cmd.PersistentFlags().Var(f.Path, "kubeconfig", fmt.Sprintf("Path to the kubeconfig file ($%s)", kubeconfigEnvVariableKey))

	kubeconfigContextEnvVariableKey := fmt.Sprintf("%s_KUBECONFIG_CONTEXT", strings.ToUpper(opts.BinaryName))
	f.Context = NewKubeconfigContextFlag(kubeconfigContextEnvVariableKey)
	cmd.PersistentFlags().Var(f.Context, "kubeconfig-context", fmt.Sprintf("Kubeconfig context override ($%s)", kubeconfigContextEnvVariableKey))

	kubeconfigYamlEnvVariableKey := fmt.Sprintf("%s_KUBECONFIG_YAML", strings.ToUpper(opts.BinaryName))
	f.YAML = NewKubeconfigYAMLFlag(kubeconfigYamlEnvVariableKey)
	cmd.PersistentFlags().Var(f.YAML, "kubeconfig-yaml", fmt.Sprintf("Kubeconfig contents as YAML ($%s)", kubeconfigYamlEnvVariableKey))
}

type KubeconfigPathFlag struct {
	value          string
	envVariableKey string
}

var _ pflag.Value = &KubeconfigPathFlag{}
var _ cobrautil.ResolvableFlag = &KubeconfigPathFlag{}

func NewKubeconfigPathFlag(envVariableKey string) *KubeconfigPathFlag {
	return &KubeconfigPathFlag{envVariableKey: envVariableKey}
}

func (s *KubeconfigPathFlag) Set(val string) error {
	s.value = val
	return nil
}

func (s *KubeconfigPathFlag) Type() string   { return "string" }
func (s *KubeconfigPathFlag) String() string { return "" } // default for usage

func (s *KubeconfigPathFlag) Value() (string, error) {
	err := s.Resolve()
	if err != nil {
		return "", err
	}

	return s.value, nil
}

func (s *KubeconfigPathFlag) Resolve() error {
	if len(s.value) > 0 {
		return nil
	}

	s.value = s.resolveValue()

	return nil
}

func (s *KubeconfigPathFlag) resolveValue() string {
	path := os.Getenv(s.envVariableKey)
	if len(path) > 0 {
		return path
	}

	return ""
}

type KubeconfigContextFlag struct {
	value          string
	envVariableKey string
}

var _ pflag.Value = &KubeconfigContextFlag{}
var _ cobrautil.ResolvableFlag = &KubeconfigPathFlag{}

func NewKubeconfigContextFlag(envVariableKey string) *KubeconfigContextFlag {
	return &KubeconfigContextFlag{envVariableKey: envVariableKey}
}

func (s *KubeconfigContextFlag) Set(val string) error {
	s.value = val
	return nil
}

func (s *KubeconfigContextFlag) Type() string   { return "string" }
func (s *KubeconfigContextFlag) String() string { return "" } // default for usage

func (s *KubeconfigContextFlag) Value() (string, error) {
	err := s.Resolve()
	if err != nil {
		return "", err
	}

	return s.value, nil
}

func (s *KubeconfigContextFlag) Resolve() error {
	if len(s.value) > 0 {
		return nil
	}

	s.value = os.Getenv(s.envVariableKey)

	return nil
}

type KubeconfigYAMLFlag struct {
	value          string
	envVariableKey string
}

var _ pflag.Value = &KubeconfigYAMLFlag{}
var _ cobrautil.ResolvableFlag = &KubeconfigPathFlag{}

func NewKubeconfigYAMLFlag(envVariableKey string) *KubeconfigYAMLFlag {
	return &KubeconfigYAMLFlag{envVariableKey: envVariableKey}
}

func (s *KubeconfigYAMLFlag) Set(val string) error {
	s.value = val
	return nil
}

func (s *KubeconfigYAMLFlag) Type() string   { return "string" }
func (s *KubeconfigYAMLFlag) String() string { return "" } // default for usage

func (s *KubeconfigYAMLFlag) Value() (string, error) {
	err := s.Resolve()
	if err != nil {
		return "", err
	}

	return s.value, nil
}

func (s *KubeconfigYAMLFlag) Resolve() error {
	if len(s.value) > 0 {
		return nil
	}

	s.value = os.Getenv(s.envVariableKey)

	return nil
}
