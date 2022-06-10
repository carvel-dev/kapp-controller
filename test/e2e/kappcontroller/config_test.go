// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"testing"
	"strings"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"sigs.k8s.io/yaml"
)

func TestConfig_HTTPProxy(t *testing.T) {
	assert := assert.New(t)
	logger := e2e.Logger{}
	kubectl := e2e.Kubectl{t, "kapp-controller", logger}

	// Proxy configured in config-test/secret-config.yml
	logger.Section("inspect controller logs for propagation of proxy env vars", func() {
		// app name must match the app name being deployed in hack/deploy-test.sh
		out := kubectl.Run([]string{"logs", "deployment/kapp-controller", "-c", "kapp-controller-sidecarexec"})
		assert.Contains(out, "Setting http_proxy", "should be non-empty value, so set")
		assert.Contains(out, "Clearing https_proxy", "should be empty value, so clear")
		assert.Contains(out, "Setting no_proxy", "should be non-empty value, so set")
	})
}

func TestConfig_KappDeployRawOptions(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	name := "global-kapp-deploy-raw-opts"
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
          file.yml: |
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              key: value
  template:
    - ytt: {}
  deploy:
    - kapp:
        rawOptions: ["--labels=local-lbl=local-lbl-val"]`, name) + sas.ForNamespaceYAML()

	cleanUpApp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUpApp()
	defer cleanUpApp()

	// Global label is configured in config-test/secret-config.yml
	logger.Section("deploy and check that kc-test label is set", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})

		// Exactly one app should have global and local label combination
		out := kapp.Run([]string{"ls", "-A", "--filter-labels", "kc-test=kc-test-val,local-lbl=local-lbl-val", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))
		assert.Equal(t, len(resp.Tables[0].Rows), 1)
	})
}

func TestConfig_TrustCACerts(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	// When updating, certs and keys must be regenerated for server and added to server.go and config-test/config-map.yml
	serverNamespace := "https-server"

	yaml1 := fmt.Sprintf(`---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-https
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - http:
      url: https://https-svc.%s.svc.cluster.local:443/deployment.yml
  template:
  - ytt: {}
  deploy:
  - kapp:
      inspect: {}
      intoNs: %s
`, serverNamespace, env.Namespace) + sas.ForNamespaceYAML()

	name := "test-https"
	httpsServerName := "test-https-server"

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
		kapp.Run([]string{"delete", "-a", httpsServerName, "-n", serverNamespace})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy https server with self signed certs", func() {
		kapp.Run([]string{"deploy", "-f", "../assets/https-server/server.yml", "-f", "../assets/https-server/certs-for-custom-ca.yml", "-a", httpsServerName})
	})

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1)})

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

		var cr v1alpha1.App

		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		expectedStatus := v1alpha1.AppStatus{
			GenericStatus: v1alpha1.GenericStatus{
				Conditions: []v1alpha1.Condition{{
					Type:   v1alpha1.ReconcileSucceeded,
					Status: corev1.ConditionTrue,
				}},
				ObservedGeneration:  1,
				FriendlyDescription: "Reconcile succeeded",
			},
			Deploy: &v1alpha1.AppStatusDeploy{
				ExitCode: 0,
				Finished: true,
			},
			Fetch: &v1alpha1.AppStatusFetch{
				ExitCode: 0,
			},
			Inspect: &v1alpha1.AppStatusInspect{
				ExitCode: 0,
			},
			Template: &v1alpha1.AppStatusTemplate{
				ExitCode: 0,
			},
			ConsecutiveReconcileSuccesses: 1,
		}

		{
			// deploy
			if !strings.Contains(cr.Status.Deploy.Stdout, "Wait to:") {
				t.Fatalf("Expected non-empty deploy output: '%s'", cr.Status.Deploy.Stdout)
			}
			cr.Status.Deploy.StartedAt = metav1.Time{}
			cr.Status.Deploy.UpdatedAt = metav1.Time{}
			cr.Status.Deploy.Stdout = ""

			// fetch
			if !strings.Contains(cr.Status.Fetch.Stdout, "kind: LockConfig") {
				t.Fatalf("Expected non-empty fetch output: '%s'", cr.Status.Fetch.Stdout)
			}
			cr.Status.Fetch.StartedAt = metav1.Time{}
			cr.Status.Fetch.UpdatedAt = metav1.Time{}
			cr.Status.Fetch.Stdout = ""

			// inspect
			if !strings.Contains(cr.Status.Inspect.Stdout, "Resources in app 'test-https-ctrl'") {
				t.Fatalf("Expected non-empty inspect output: '%s'", cr.Status.Inspect.Stdout)
			}
			cr.Status.Inspect.UpdatedAt = metav1.Time{}
			cr.Status.Inspect.Stdout = ""

			// template
			cr.Status.Template.UpdatedAt = metav1.Time{}
			cr.Status.Template.Stderr = ""
		}

		if !reflect.DeepEqual(expectedStatus, cr.Status) {
			t.Fatalf("Status is not same: %#v vs %#v", expectedStatus, cr.Status)
		}
	})

}

func TestConfig_SkipTLSVerify(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	// If this changes, the skip-tls-verify domain must be updated to match
	name := "test-skip-tls"
	registryName := "test-registry"
	configName := "test-config-skip-tls-verify-config"

	yaml1 := fmt.Sprintf(`---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - imgpkgBundle:
      image: registry-svc.registry.svc.cluster.local:443/my-repo/image
  template:
  - ytt: {}
  deploy:
  - kapp:
      inspect: {}
`, name) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
		kapp.Run([]string{"delete", "-a", registryName})
		kapp.Run([]string{"delete", "-a", configName})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy controller config to skip tls for specific domain", func() {
		config := `
apiVersion: v1
kind: Secret
metadata:
  name: kapp-controller-config
  namespace: kapp-controller
stringData:
  dangerousSkipTLSVerify: registry-svc.registry.svc.cluster.local
`
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", configName}, e2e.RunOpts{StdinReader: strings.NewReader(config)})

		// Since config propagation is async, just wait a little bit
		time.Sleep(2 * time.Second)
	})

	logger.Section("deploy registry with self signed certs", func() {
		kapp.Run([]string{"deploy", "-f", "../assets/registry/registry.yml", "-f", "../assets/registry/certs-for-skip-tls.yml", "-a", registryName})
	})

	logger.Section("deploy app that fetches contents from registry", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{AllowError: true, StdinReader: strings.NewReader(yaml1)})
		assert.Error(t, err, "Expected fetching error")

		var cr v1alpha1.App

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})
		assert.NoError(t, yaml.Unmarshal([]byte(out), &cr))

		// To avoid the complexity associated with preloading our deployed registry
		// with an image, and because imgpkg will fall back to http to pull the image
		// if https fails, this explicitly expects an error and asserts the error message
		// is due to a manifest unknown error and not TLS verification.
		assert.NotNil(t, cr.Status.Fetch)
		assert.Equal(t, cr.Status.Fetch.ExitCode, 1)
		assert.NotContains(t, cr.Status.Fetch.Stderr, "x509: certificate signed by unknown authority")
		assert.Contains(t, cr.Status.Fetch.Stderr, "MANIFEST_UNKNOWN: manifest unknown;")
	})
}
