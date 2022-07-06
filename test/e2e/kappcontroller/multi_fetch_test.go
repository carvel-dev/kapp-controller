// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"strings"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func TestMultiFetch(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	yaml1 := `
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-multi-fetch
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
  - inline:
      paths:
        file.yml: |
          #@ load("@ytt:overlay", "overlay")
          #@overlay/match by=overlay.subset({"metadata":{"name":"configmap"}})
          ---
          data:
            #@overlay/match missing_ok=True
            overlay-key: overlay-value
  template:
  - ytt: {}
  deploy:
  - kapp: {}
` + sas.ForNamespaceYAML()

	name := "test-multi-fetch"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1)})
	})

	logger.Section("verify", func() {
		out := kapp.Run([]string{"inspect", "-a", name + ".app", "--raw", "--tty=false", "--filter-kind", "ConfigMap"})

		var cm corev1.ConfigMap

		err := yaml.Unmarshal([]byte(out), &cm)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		if cm.ObjectMeta.Name != "configmap" {
			t.Fatalf(`Expected name to be "configmap" got %#v`, cm.ObjectMeta.Name)
		}

		if cm.Data["key"] != "value" {
			t.Fatalf(`Expected data.key to be "value" got %#v`, cm.Data["key"])
		}

		if cm.Data["overlay-key"] != "overlay-value" {
			t.Fatalf(`Expected data.overlay-key to be "overlay-value" got %#v`, cm.Data["key"])
		}

	})
}
