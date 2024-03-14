// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func Test_AppPause(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	name := "app-pause"

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	// create originally as not paused otherwise
	// App will not create original resources.
	appYaml := appYAML(name, env.Namespace, "original", false)
	_, err := kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	if err != nil {
		t.Fatalf("Expected initial app dfeploy to succeed, it did not: %v", err)
	}

	// update app
	appYaml = appYAML(name, env.Namespace, "change", true)
	_, err = kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	if err != nil {
		t.Fatalf("Expected initial app dfeploy to succeed, it did not: %v", err)
	}

	var cr v1alpha1.App
	out, err := kubectl.RunWithOpts([]string{"get", "apps/" + name, "-o", "yaml"}, e2e.RunOpts{AllowError: true})
	if err != nil {
		t.Fatalf("failed to get App %s: %v", name, err)
	}

	err = yaml.Unmarshal([]byte(out), &cr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}

	if cr.Status.FriendlyDescription != "Canceled/paused" {
		t.Fatalf("expected App to have status Canceled/paused\nGot: %s", cr.Status.FriendlyDescription)
	}

	out, err = kubectl.RunWithOpts([]string{"get", "configmap/configmap", "-o", "yaml"}, e2e.RunOpts{AllowError: true})
	if err != nil {
		t.Fatalf("failed to get configmap/configmap")
	}

	var cm corev1.ConfigMap
	err = yaml.Unmarshal([]byte(out), &cm)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}

	if cm.Data["key"] != "original" {
		t.Fatalf("configmap message was updated despite App being paused\nGot: %s", cm.Data["key"])
	}
}

func appYAML(name, namespace, configMapValue string, paused bool) string {
	sas := e2e.ServiceAccounts{namespace}
	return fmt.Sprintf(`---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  namespace: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  paused: %t
  fetch:
  - inline:
      paths:
        file.yml: |
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: configmap
          data:
            key: %s
  template:
  - ytt: {}
  deploy:
  - kapp: {}`, name, namespace, paused, configMapValue) + sas.ForNamespaceYAML()
}
