// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	"sigs.k8s.io/yaml"
)

func Test_App_FetchDirname_Single(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}
	name := "test-app-fetch-dirname-single"

	appSingleFetchYAML := `---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-app-fetch-dirname-single
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
      paths:
      - foo
      inline:
        paths:
          file.yml: |
            #@ load("@ytt:data", "data")
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              value: #@ data.read("file.csv")
  deploy:
  - kapp: {}
` + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		// if template stage succeeds, assume test pass
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{StdinReader: strings.NewReader(appSingleFetchYAML)})
	})

	out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

	var cr v1alpha1.App
	err := yaml.Unmarshal([]byte(out), &cr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}

	assert.Contains(t, cr.Status.Fetch.Stdout, "path: foo", "Expected destDirName to be used by vendir")
}

func Test_App_FetchDirname_All(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}
	name := "test-app-fetch-dirname-all"

	appAllFetchYAML := fmt.Sprintf(`---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-app-fetch-dirname-multiple
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
  - imgpkgBundle:
      image: k8slt/kctrl-example-pkg:v1.0.0
    path: imgpkg-dest
  - helmChart:
      name: redis
      # Chart version v1, DEPRECATED 
      version: "11.3.4"
      repository:
        url: https://charts.bitnami.com/bitnami
    path: helm-dest
  template:
  - helmTemplate:
      path: helm-dest
  - ytt:
      paths:
      - http-dest
      - git-dest
      - imgpkg-dest/config
  deploy:
  - kapp:
      inspect: {}
      intoNs: %s
`, env.Namespace) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy multiple", func() {
		// if template stage succeeds, assume test pass
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{StdinReader: strings.NewReader(appAllFetchYAML)})
	})

	out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

	var cr v1alpha1.App
	err := yaml.Unmarshal([]byte(out), &cr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}

	testCases := []string{"git", "http", "helm", "imgpkg"}
	for _, c := range testCases {
		assert.Contains(t, cr.Status.Fetch.Stdout, fmt.Sprintf("path: %s-dest", c), "Expected destDirName to be used by vendir")
	}
}
