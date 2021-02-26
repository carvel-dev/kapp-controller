// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"strings"
	"testing"
)

func Test_AppReconcileOccurs_WhenSecretUpdated(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}

	name := "simple-app-with-secret"
	// syncPeriod set to 20 minutes so that test
	// won't pass because of reconcile from time sync.
	appYaml := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  syncPeriod: 20m
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - git:
      url: https://github.com/k14s/k8s-simple-app-example
      ref: origin/develop
      subPath: config-step-2-template
  template:
  - ytt:
      inline:
        pathsFrom:
        - secretRef:
            name: simple-app-values
  deploy:
  - kapp: {}
---
apiVersion: v1
kind: Secret
metadata:
  name: simple-app-values
stringData:
  values2.yml: |
    #@data/values
    ---
    hello_msg: original
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: exec-curl
spec:
  selector:
    matchLabels:
      app: exec-curl
  template:
    metadata:
      labels:
        app: exec-curl
    spec:
      containers:
        - image: k8s.gcr.io/echoserver:1.4
          imagePullPolicy: IfNotPresent
          name: echoserver`, name) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			RunOpts{IntoNs: true, StdinReader: strings.NewReader(appYaml)})

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

		var cr v1alpha1.App
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		// check for successful deployment
		if cr.Status.Deploy == nil || !cr.Status.Deploy.Finished || cr.Status.Deploy.ExitCode != 0 {
			t.Fatalf("Expected simple-app deployment to succeed but got:\n%s", cr.Status.Deploy.Stdout)
		}
	})

	logger.Section("update secret", func() {
		updatedSecret := `
---
apiVersion: v1
kind: Secret
metadata:
  name: simple-app-values
stringData:
  values2.yml: |
    #@data/values
    ---
    hello_msg: updated`

		// Update secret
		kubectl.RunWithOpts([]string{"apply", "-f", "-"}, RunOpts{StdinReader: strings.NewReader(updatedSecret)})

		// Make sure App reconciles from secret update/reconcile succeeds before asserting on response from App
		kubectl.Run([]string{"wait", "--for=condition=Reconciling", "apps/" + name, "--timeout", "1m"})
		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "apps/" + name, "--timeout", "1m"})
	})

	logger.Section("check App uses new secret", func() {
		// get ClusterIP of simple-app to use with kubectl exec
		svcIP := kubectl.Run([]string{"get", "svc/simple-app", "-o", "jsonpath={.spec.clusterIP}"})
		// kubectl exec into debug pod and run curl against simple-app
		// to get response with secret
		appResponse, _ := kubectl.RunWithOpts([]string{"exec", "deployment/exec-curl", "-n", env.Namespace, "--", "curl", svcIP, "-s"}, RunOpts{NoNamespace: true})

		if !strings.Contains(appResponse, "updated") {
			t.Fatalf("\nSecret message was not updated to Hello updated!\nGot:%s", appResponse)
		}
	})

	logger.Section("check App status", func() {
		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

		var cr v1alpha1.App
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		if cr.Status.ConsecutiveReconcileSuccesses != 3 {
			t.Fatalf("Expected only two App reconciles but got %d", cr.Status.ConsecutiveReconcileSuccesses)
		}
	})
}
