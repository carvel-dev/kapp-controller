// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	corev1 "k8s.io/api/core/v1"
)

func Test_YttTemplate_UsesFileMarks(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

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
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("check ConfigMap exists", func() {
		kubectl.Run([]string{"get", "configmap", name})
	})
}

func Test_YttTemplate_ValuesFrom(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

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
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
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

func Test_SecretsAndConfigMapsWithCustomPathsCanReconcile(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	name := "inline-pathsfrom-directorypath"
	appYaml := `---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: simple-app2
  namespace: kappctrl-test
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - inline:
      paths:
        file-from-fetch-paths: ""
        template/dir/blah.txt: ""
      pathsFrom:
        - secretRef:
            name: from-secret
            directoryPath: fetch/dis/dir
        # both secrets have same key, so without directoryPath one would overwrite the other
        - secretRef:
            name: another-from-secret
            directoryPath: fetch/dat/dir
  template:
  - ytt:
      inline:
        pathsFrom:
        - secretRef:
            name: from-secret
            directoryPath: template/dir
        paths:
          config.yml: |
            #@ load("@ytt:data", "data")
            #@ load("@ytt:yaml", "yaml")

            #@ print(data.list())

            kind: ConfigMap
            apiVersion: v1
            metadata:
              name: foo
            data:
              data: #@ yaml.encode(data.list())
  deploy:
  - kapp: {}
---
kind: Secret
apiVersion: v1
metadata:
  name: from-secret
  namespace: kappctrl-test
stringData:
  file-from-secret: "thisIsData"
---
kind: Secret
apiVersion: v1
metadata:
  name: another-from-secret
  namespace: kappctrl-test
stringData:
  file-from-secret: "SameKeyDifferentValues"
` + sas.ForNamespaceYAML()
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("check configmap is populated", func() {
		out := kubectl.Run([]string{"get", "configmap/foo", "-o", "yaml"})
		expectedContentItems := []string{
			"file-from-fetch-paths",
			"config.yml",
			"fetch/dis/dir/file-from-secret",
			"fetch/dat/dir/file-from-secret",
			"template/dir/blah.txt",
			"template/dir/file-from-secret",
		}

		for _, item := range expectedContentItems {
			if !strings.Contains(out, item) {
				t.Fatal("failed", fmt.Errorf("configmap %s missing item: %s", out, item))
			}
		}
	})
}
