// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
)

func Test_Managed_Name_App_Migration(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}
	oldName := "simple-app-ctrl"
	name := "simple-app"

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
		kapp.Run([]string{"delete", "-a", oldName})
	}
	cleanUp()
	defer cleanUp()

	config := `---
kind: ConfigMap 
apiVersion: v1 
metadata:
  name: simple-configmap-e2e-test-migration`

	appYaml := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
    - inline:
        paths:
          file.yml: |
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              key: value
  template:
    - ytt: {}
  deploy:
    - kapp:
        inspect: {}`, name) + sas.ForNamespaceYAML()

	logger.Section("deploy a simple configmap with the managed name", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", oldName},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(config)})
	})

	logger.Section("deploy an AppCR with managed name", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	_, err := kapp.RunWithOpts([]string{"inspect", "-a", name + ".app", "--raw", "--tty=false", "--filter-kind=App"}, e2e.RunOpts{AllowError: true})
	assert.NoError(t, err, fmt.Sprintf("expected %s to exist but does not", name))

	_, err = kapp.RunWithOpts([]string{"inspect", "-a", oldName, "--raw", "--tty=false", "--filter-kind=App"}, e2e.RunOpts{AllowError: true})
	assert.Error(t, err, fmt.Sprintf("expected %s not to exist", oldName))
	assert.ErrorContainsf(t, err, "does not exist", "expected 'does not exist' error", oldName)
}

func Test_Managed_Name_Package_Repository_Migration(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}
	oldName := "simple-pkgr-ctrl"
	name := "simple-pkgr"

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
		kapp.Run([]string{"delete", "-a", oldName})
	}
	cleanUp()
	defer cleanUp()

	config := `---
kind: ConfigMap 
apiVersion: v1 
metadata:
  name: simple-configmap-e2e-test-migration`

	pkgrYaml := fmt.Sprintf(`
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name:  %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  fetch:
    inline:
      paths:
        packages/pkg.test.carvel.dev/pkg.test.carvel.dev.0.0.1.yml: |
          ---
          apiVersion: data.packaging.carvel.dev/v1alpha1
          kind: Package
          metadata:
            name: pkg.test.carvel.dev.0.0.1
          spec:
            refName: pkg.test.carvel.dev
            version: 0.0.1
            template:
              spec: {}`, name) + sas.ForNamespaceYAML()

	logger.Section("deploy a simple configmap with the managed name", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", oldName}, e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(config)})
	})

	logger.Section("deploy a PackageRepository with managed name", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(pkgrYaml)})
	})

	_, err := kapp.RunWithOpts([]string{"inspect", "-a", name + ".pkgr", "--raw", "--tty=false", "--filter-kind=PackageRepository"}, e2e.RunOpts{AllowError: true})
	assert.NoError(t, err, fmt.Sprintf("expected %s to exist but does not", name))

	_, err = kapp.RunWithOpts([]string{"inspect", "-a", oldName, "--raw", "--tty=false", "--filter-kind=PackageRepository"}, e2e.RunOpts{AllowError: true})
	assert.Error(t, err, fmt.Sprintf("expected %s not to exist", oldName))
	assert.ErrorContainsf(t, err, "does not exist", "expected 'does not exist' error", oldName)
}
