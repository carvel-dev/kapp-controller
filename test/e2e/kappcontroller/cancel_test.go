// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestCancel(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	yaml1 := `
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: cancel
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - inline:
      paths:
        file.yml: |
          apiVersion: batch/v1
          kind: Job
          metadata:
            name: sleep
          spec:
            template:
              metadata:
                name: sleep
              spec:
                containers:
                - name: sleep
                  image: busybox
                  command: ['sh', '-c', 'sleep 60']
                restartPolicy: Never
                terminationGracePeriodSeconds: 0
  template:
  - ytt: {}
  deploy:
  - kapp: {}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: filler
data: {}
` + sas.ForNamespaceYAML()

	yaml2 := `
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: cancel
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  canceled: true
  fetch:
  - inline:
      paths:
        file.yml: |
          apiVersion: batch/v1
          kind: Job
          metadata:
            name: sleep
          spec:
            template:
              metadata:
                name: sleep
              spec:
                containers:
                - name: sleep
                  image: busybox
                  command: ['sh', '-c', 'sleep 60']
                restartPolicy: Never
                terminationGracePeriodSeconds: 0
  template:
  - ytt: {}
  deploy:
  - kapp: {}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: filler
data: {}
` + sas.ForNamespaceYAML()

	name := "test-cancel"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("begin deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name, "--wait=false"},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1)})

		waitForDeployToStart(name, kapp, t)
	})

	logger.Section("cancel deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name, "--wait=false"},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml2)})

		waitForReconcileFailed(name, kapp, t)
	})

	app := obtainApp(name, kapp, t)

	logger.Section("initiate delete", func() {
		// kick off deletion but dont wait; do this
		// before un-canceling as otherwise deploy will kick off again
		kapp.RunWithOpts([]string{"delete", "-a", name, "--wait=false", "--filter-kind", "App"}, e2e.RunOpts{})

		// un-cancel now
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name, "--wait=false"},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1)})
	})

	logger.Section("verify canceled status", func() {
		expectedConditions := []v1alpha1.Condition{{
			Type:    v1alpha1.ReconcileFailed,
			Status:  "True",
			Message: "Deploying: Process was canceled",
		}}

		if !reflect.DeepEqual(expectedConditions, app.Status.Conditions) {
			t.Fatalf("Status conditions are not same: %#v vs %#v", expectedConditions, app.Status.Conditions)
		}

		expectedDeploy := &v1alpha1.AppStatusDeploy{
			Finished: true,
			ExitCode: -1,
			Error:    "Deploying: Process was canceled",
		}

		app.Status.Deploy.StartedAt = metav1.Time{}
		app.Status.Deploy.UpdatedAt = metav1.Time{}
		app.Status.Deploy.Stdout = ""

		if !reflect.DeepEqual(expectedDeploy, app.Status.Deploy) {
			t.Fatalf("Status deploy is not same: %#v vs %#v", expectedDeploy, app.Status.Deploy)
		}
	})
}

func obtainApp(name string, kapp e2e.Kapp, t *testing.T) v1alpha1.App {
	out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind", "App"})

	var app v1alpha1.App

	err := yaml.Unmarshal([]byte(out), &app)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}

	return app
}

func waitForReconcileFailed(name string, kapp e2e.Kapp, t *testing.T) {
	count := 0
	for {
		app := obtainApp(name, kapp, t)

		for _, cond := range app.Status.Conditions {
			if cond.Type == v1alpha1.ReconcileFailed {
				return
			}
		}

		if count > 300 {
			t.Fatalf("Failed to reach reconciling state")
		}

		count++
		time.Sleep(1 * time.Second)
	}
}

func waitForDeployToStart(name string, kapp e2e.Kapp, t *testing.T) {
	count := 0
	for {
		app := obtainApp(name, kapp, t)
		switch {
		case app.Status.Deploy != nil && len(app.Status.Deploy.Stdout) > 0:
			return
		case count > 20:
			t.Fatalf("Failed to reach reconciling state")
		default:
			count++
			time.Sleep(1 * time.Second)
		}
	}
}
