// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"
)

func Test_YttTemplate_UsesFileMarks(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}

	name := "configmap-with-non-yml-ext-file"
	// This App's ConfigMap is in a non yaml file
	// named file. In order for ytt to render the YAML,
	// the file needs a file-mark to denote it is a
	// YAML file.
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
          file: |
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              key: value
  template:
    - ytt:
        paths:
          - file
        fileMarks:
          - file:type=yaml-plain
  deploy:
    - kapp: {}
`, name) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("check ConfigMap exists", func() {
		kubectl.Run([]string{"get", "configmap", name})
	})
}
