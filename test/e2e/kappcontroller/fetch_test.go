// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
)

func Test_App_FetchPath(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	appSingleFetchPathYAML := `---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-app-fetch-path-single
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - inline:
      paths:
        file.csv: |
          Year,Make,Model,Description,Price
          1997,Ford,E350,"ac, abs, moon",3000.00
    path: foo
  template:
  - ytt:
      inline:
        paths:
          file.yml: |
            #@ load("@ytt:assert", "assert")
            #@ load("@ytt:data", "data")
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              key: #@ assert.equals(data.list(), ["file.yml", "foo/file.csv"])
  deploy:
  - kapp:
      rawOptions: ["--dangerous-allow-empty-list-of-resources=true"]
` + sas.ForNamespaceYAML()

	appSingleFetchNoPathYAML := `---
  apiVersion: kappctrl.k14s.io/v1alpha1
  kind: App
  metadata:
    name: test-app-fetch-path-single
    annotations:
      kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
  spec:
    serviceAccountName: kappctrl-e2e-ns-sa
    fetch:
    - inline:
        paths:
          file.csv: |
            Year,Make,Model,Description,Price
            1997,Ford,E350,"ac, abs, moon",3000.00
    template:
    - ytt:
        inline:
          paths:
            file.yml: |
              #@ load("@ytt:assert", "assert")
              #@ load("@ytt:data", "data")
              apiVersion: v1
              kind: ConfigMap
              metadata:
                name: configmap
              data:
                key: #@ assert.equals(data.list(), ["file.csv", "file.yml"])
    deploy:
    - kapp: {}
` + sas.ForNamespaceYAML()

	appAllFetchMultiplePathYAML := fmt.Sprintf(`---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-app-fetch-path-multiple
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - git:
      url: https://github.com/k14s/kapp
      ref: origin/develop
      subPath: examples/gitops/guestbook
    path: git-dest
  - http:
      url: https://raw.githubusercontent.com/k14s/kapp/db2cc63e12e988235eb8815af8edc0ca0cfaa79c/examples/gitops/nginx/config.yml
      sha256: 35d7a2798393d93548922bcb5bb19dafc59535a8ce1c0afa8130e8137ba43d9f
    path: http-dest
  template:
  - ytt:
      inline:
        paths:
          file.yml: |
            #@ load("@ytt:assert", "assert")
            #@ load("@ytt:data", "data")
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              key: #@ assert.equals(data.list(), ["file.yml", "git-dest/README.md", "git-dest/all-in-one.yml", "http-dest/config.yml"])
  deploy:
  - kapp:
      intoNs: %s
`, env.Namespace) + sas.ForNamespaceYAML()

	appAllFetchMultipleNoPathYAML := fmt.Sprintf(`---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-app-fetch-path-multiple
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - git:
      url: https://github.com/k14s/kapp
      ref: origin/develop
      subPath: examples/gitops/guestbook
    path: git-dest
  - http:
      url: https://raw.githubusercontent.com/k14s/kapp/db2cc63e12e988235eb8815af8edc0ca0cfaa79c/examples/gitops/nginx/config.yml
      sha256: 35d7a2798393d93548922bcb5bb19dafc59535a8ce1c0afa8130e8137ba43d9f
  template:
  - ytt:
      inline:
        paths:
          file.yml: |
            #@ load("@ytt:assert", "assert")
            #@ load("@ytt:data", "data")
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              key: #@ assert.equals(data.list(), ["file.yml", "1/config.yml", "git-dest/README.md", "git-dest/all-in-one.yml"])
  deploy:
  - kapp:
      intoNs: %s
`, env.Namespace) + sas.ForNamespaceYAML()

	tests := []struct {
		desc           string
		name           string
		deploymentYAML string
	}{
		{
			"Single Fetch - path given",
			"test-single-fetch-path",
			appSingleFetchPathYAML,
		},
		{
			"Single Fetch - no path",
			"test-single-fetch-no-path",
			appSingleFetchNoPathYAML,
		},
		{
			"Multiple Fetch - path",
			"test-multiple-fetch-path",
			appAllFetchMultiplePathYAML,
		},
		{
			"Multiple Fetch - no path",
			"test-multiple-fetch-no-path",
			appAllFetchMultipleNoPathYAML,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			cleanUp := func() {
				kapp.Run([]string{"delete", "-a", tc.name})
			}
			cleanUp()
			defer cleanUp()

			logger.Section("deploy", func() {
				// if template stage succeeds, assume test pass
				kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", tc.name},
					e2e.RunOpts{StdinReader: strings.NewReader(tc.deploymentYAML)})
			})
		})
	}
}
