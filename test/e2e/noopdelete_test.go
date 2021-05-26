// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"
)

func Test_NoopDelete_DeletesAfterServiceAccountDeleted(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}
	name := "instl-pkg-noop-delete"
	cfgMapName := "configmap"

	installPkgYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: noopdelete.carvel.dev
spec:
  longDescription: noopdelete-test
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageVersion
metadata:
  name: noopdelete.carvel.dev.1.0.0
spec:
  packageName: noopdelete.carvel.dev
  version: 1.0.0
  template:
    spec:
      fetch:
      - inline:
          paths:
            file.yml: |
              apiVersion: v1
              kind: ConfigMap
              metadata:
                name: %s
              data:
                key: value
      template:
      - ytt: {}
      deploy:
      - kapp: {}
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: InstalledPackage
metadata:
 name: %s
 namespace: %s
 annotations:
   kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/installedpackages
spec:
 serviceAccountName: kappctrl-e2e-ns-sa
 noopDelete: true
 packageVersionRef:
   packageName: noopdelete.carvel.dev
   versionSelection:
     constraint: 1.0.0
`, cfgMapName, name, env.Namespace) + sas.ForNamespaceYAML()

	cleanUpIpkg := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUpConfigMap := func() {
		kubectl.Run([]string{"delete", "configmap", cfgMapName})
	}

	cleanUpIpkg()
	defer cleanUpIpkg()
	defer cleanUpConfigMap()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, RunOpts{StdinReader: strings.NewReader(installPkgYaml)})
	})

	kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "apps/" + name, "--timeout", "1m"})
	logger.Section("delete Service Account and InstalledPackage", func() {
		kubectl.Run([]string{"delete", "serviceaccount", "kappctrl-e2e-ns-sa"})
		cleanUpIpkg()
	})

	logger.Section("check ConfigMap still exists after delete", func() {
		kubectl.Run([]string{"get", "configmap/" + cfgMapName})
	})
}
