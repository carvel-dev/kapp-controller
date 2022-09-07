// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"
	"time"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	"sigs.k8s.io/yaml"
)

func TestConfig_HTTPProxy(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, "kapp-controller", logger}

	configName := "test-config-http-proxy-config"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", configName})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("inspect controller logs for empty proxy vars at startup", func() {
		// app name must match the app name being deployed in hack/deploy-test.sh
		out := kubectl.Run([]string{"logs", "deployment/kapp-controller", "-c", "kapp-controller-sidecarexec"})
		assert.Contains(t, out, "Clearing http_proxy", "should be empty value, so set")
		assert.Contains(t, out, "Clearing https_proxy", "should be empty value, so clear")
		assert.Contains(t, out, "Clearing no_proxy", "should be empty value, so set")
	})

	logger.Section("change proxy configuration", func() {
		config := `
apiVersion: v1
kind: Secret
metadata:
  name: kapp-controller-config
  namespace: kapp-controller
stringData:
  httpProxy: proxy-svc.proxy-server.svc.cluster.local:80
  httpsProxy: proxy-svc.proxy-server.svc.cluster.local:80
  noProxy: docker.io,KAPPCTRL_KUBERNETES_SERVICE_HOST
`
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", configName}, e2e.RunOpts{StdinReader: strings.NewReader(config)})

		// Since config propagation is async, just wait a little bit
		time.Sleep(2 * time.Second)
	})

	logger.Section("inspect controller logs for propagation of proxy env vars", func() {
		// app name must match the app name being deployed in hack/deploy-test.sh
		out := kubectl.Run([]string{"logs", "deployment/kapp-controller", "-c", "kapp-controller-sidecarexec"})
		assert.Contains(t, out, "Setting http_proxy", "should be non-empty value, so set")
		assert.Contains(t, out, "Setting https_proxy", "should be non-empty value, so clear")
		assert.Contains(t, out, "Setting no_proxy", "should be non-empty value, so set")
	})
}

func TestConfig_KappDeployRawOptions(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	configName := "test-config-kapp-deploy-raw-opts-config"
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
		kapp.Run([]string{"delete", "-a", configName})
	}
	cleanUpApp()
	defer cleanUpApp()

	logger.Section("deploy controller config to set global label", func() {
		config := `
apiVersion: v1
kind: Secret
metadata:
  name: kapp-controller-config
  namespace: kapp-controller
stringData:
  kappDeployRawOptions: "[\"--diff-changes=true\", \"--labels=kc-test=kc-test-val\"]"
`
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", configName}, e2e.RunOpts{StdinReader: strings.NewReader(config)})

		// Since config propagation is async, just wait a little bit
		time.Sleep(2 * time.Second)
	})

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
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	name := "test-https"
	pkgrName := "test-https-pkgr"
	httpsServerName := "test-https-server"
	configName := "test-config-trust-ca-config"

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
      # use https to exercise CA certificate validation
      # When updating address, certs and keys must be regenerated
      # for server and added to e2e/assets/https-server
      url: https://https-svc.https-server.svc.cluster.local:443/deployment.yml
  template:
  - ytt: {}
  deploy:
  - kapp:
      inspect: {}
      intoNs: %s
`, env.Namespace) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
		kapp.Run([]string{"delete", "-a", pkgrName})
		kapp.Run([]string{"delete", "-a", configName})
		kapp.Run([]string{"delete", "-a", httpsServerName})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy controller config to trust CA cert", func() {
		config := `
apiVersion: v1
kind: Secret
metadata:
  name: kapp-controller-config
  namespace: kapp-controller
stringData:
  # Must match the second cert in the cert chain in test/e2e/assets/self-signed-https-server.yml
  caCerts: |
    -----BEGIN CERTIFICATE-----
    MIIEaTCCAtGgAwIBAgIQMnHSoj2so3Ye4U0CepDOfTANBgkqhkiG9w0BAQsFADA9
    MQwwCgYDVQQGEwNVU0ExFjAUBgNVBAoTDUNsb3VkIEZvdW5kcnkxFTATBgNVBAMT
    DGdlbmVyYXRlZC1jYTAgFw0yMjAzMjIxNjA4NDNaGA8yMTIyMDIyNjE2MDg0M1ow
    PTEMMAoGA1UEBhMDVVNBMRYwFAYDVQQKEw1DbG91ZCBGb3VuZHJ5MRUwEwYDVQQD
    EwxnZW5lcmF0ZWQtY2EwggGiMA0GCSqGSIb3DQEBAQUAA4IBjwAwggGKAoIBgQDm
    1mAC3HRlZd7ZTlPPB2K5AxHl8luSGmRm4UnYXxxCaoKNJAfP9Fr/f7NOXSss/R02
    F9JKH9UIAOaxSvyGnQegbbpRkRwvgPt76TSMrvwq/Qvr+beocJXeIbgNXY18/SLe
    jDyMJezDhWcOYolXOWD6+pNzJ5QjenidO82LVKOtp9umHRMqZbaBhW0AbN9WwV1e
    YM+iU/l9Ql7H+meDAioGP/NSduHtyD6dtgfFGVxwKEoU0HmVwCMsgcU5DVbexk01
    SDFOHNv1adfKIB0NQNZZNuT45QV3En2jON79EP7QQQ3kcX65BRv+AWsP0TNoa8SI
    Tma097oFnoats7JpcGptcgCafaZq1suGs2Lcc004cCOvcquw6ow3hXw0YCKZHDNO
    TGPdylU8T3FTrB9gJMBrwCs7OqjCL83m6vr68vICswNch6jaVaTkiRheTfjUyShP
    GmUsCvv/yT5sBt6kjzlCTtGlSKDOYxEqoMbvsV34Cb1qUUjoalYKfsn3Fo6ttVMC
    AwEAAaNjMGEwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0O
    BBYEFOxC7vVMU3/CO7q9Ylp1zVL9AW2cMB8GA1UdIwQYMBaAFOxC7vVMU3/CO7q9
    Ylp1zVL9AW2cMA0GCSqGSIb3DQEBCwUAA4IBgQDICCkhON4+AIxHrbtK1rfuF7vK
    Ck1yL0k482H+FH1bkMXCTGTtBPsk9yvG0mGzSi6f52euh2m+ZKWp5MXRcPJT5OUC
    59oXZhLiHBeQRQ5cJRXxz7OgsGORwjWIrjU1mHq6xAIwl59v0QennCDHUFzu6nPw
    6dgy+excnZ4KJmH70D3/QRxfj2nuxe5KyobTyOQIawRl1TTgLSRMchiDp23TbIWe
    ZLiyb2CoWdRfQfEanwYbAavyhYNQJCWLwDExBYEV5Ep6hr1g5E8jHN6f+/0a5nkK
    GES8ooNXEsm9QTuA2Cnvf8a9jYoRAHrMoL0KlaP+0HikjFoySafl5UFdm/iEWVRV
    fmDRVhlZZ0bHX/0jR1woV/Nlz3dRysMH4M7/FKsuPFYg9xOfqa0PwBFNK0Os1jM7
    WM+DlzZxGMBd7QKW7xCdEuUmKxB8gQw0LvStYM/38MB5KMDtFo/uTIkr1HsEpSNG
    lYEKi+1KNYrJFl+DIUQVWoC+fi0Doiqor2D2Zkk=
    -----END CERTIFICATE-----
`
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", configName},
			e2e.RunOpts{StdinReader: strings.NewReader(config)})

		// Since config propagation is async, just wait a little bit
		time.Sleep(2 * time.Second)
	})

	logger.Section("deploy https server with self signed certs", func() {
		kapp.Run([]string{"deploy", "-f", "../assets/https-server/", "-a", httpsServerName})
	})

	logger.Section("deploy app that fetches content from http server", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(yaml1),
			OnErrKubectl: []string{"get", "app/test-https", "-oyaml"},
		})

		out := kubectl.Run([]string{"get", "cm", "http-server-returned-cm", "-oyaml"})
		assert.Contains(t, out, "content: http-server-returned-content\n")
	})

	logger.Section("deploy package repository that fetches content from http server", func() {
		pkgrConfig := `---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: test-https-pkgr
spec:
  fetch:
    http:
      # use https to exercise CA certificate validation
      # When updating address, certs and keys must be regenerated
      # for server and added to e2e/assets/https-server
      url: https://https-svc.https-server.svc.cluster.local:443/packages.tar
`

		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgrName}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgrConfig),
			OnErrKubectl: []string{"get", "pkgr/test-https-pkgr", "-oyaml"},
		})

		out := kubectl.Run([]string{"get", "pkg", "package-behind-ca-cert.carvel.dev.1.0.0", "-oyaml"})
		assert.Contains(t, out, "name: package-behind-ca-cert.carvel.dev.1.0.0\n")
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
