// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"github.com/ghodss/yaml"
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

func Test_YttTemplate_ValuesFrom(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}

	name := "ytt-values-from"
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
          cm.yml: |
            #@ load("@ytt:data", "data")
            #@ load("@ytt:yaml", "yaml")
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: cm-result
            data:
              values: #@ yaml.encode(data.values)
          vals.yml: |
            from_path: true
  template:
    - ytt:
        paths:
        - cm.yml
        valuesFrom:
        - secretRef:
            name: secret-values
        - configMapRef:
            name: cm-values
        - path: vals.yml
  deploy:
    - kapp: {}
---
apiVersion: v1
kind: Secret
metadata:
  name: secret-values
stringData:
  vals.yml: |
    from_secret: true
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-values
data:
  vals.yml: |
    from_cm: true
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
		out := kubectl.Run([]string{"get", "configmap", "cm-result", "-o", "yaml"})

		var cm corev1.ConfigMap

		err := yaml.Unmarshal([]byte(out), &cm)
		if err != nil {
			t.Fatalf("Unmarshaling result config map: %s", err)	
		}

		expectedOut := `from_secret: true
from_cm: true
from_path: true
`

		if cm.Data["values"] != expectedOut {
			t.Fatalf("Values '%s' does not match expected value", cm.Data["values"])
		}
	})
}
