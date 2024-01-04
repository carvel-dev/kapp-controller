// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Env struct {
	Namespace         string
	PackagingGlobalNS string
}

func BuildEnv(t *testing.T) Env {
	env := Env{
		Namespace: os.Getenv("KAPPCTRL_E2E_NAMESPACE"),
	}

	if pkgNS := os.Getenv("KAPPCTRL_E2E_PACKAGING_NAMESPACE"); pkgNS != "" {
		env.PackagingGlobalNS = pkgNS
	} else {
		env.PackagingGlobalNS = "kapp-controller-packaging-global"
	}

	env.Validate(t)
	return env
}

func (e Env) Validate(t *testing.T) {
	errStrs := []string{}

	if len(e.Namespace) == 0 {
		errStrs = append(errStrs, "Expected Namespace to be non-empty")
	}

	if len(errStrs) > 0 {
		t.Fatalf("%s", strings.Join(errStrs, "\n"))
	}
}

// GetServerVersion returns the cluster info. 
func (e Env) GetServerVersion() (*version.Info, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("Error getting server version: %v", err)
	}
	return serverVersion, nil
}
