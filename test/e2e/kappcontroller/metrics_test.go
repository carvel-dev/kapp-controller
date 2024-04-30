// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"carvel.dev/kapp-controller/test/e2e"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMetrics(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	pkgRepoYAML := fmt.Sprintf(`
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: minimal-repo.tanzu.carvel.dev
  namespace: %s
  annotations:
    kapp.k14s.io/disable-original: ""
spec:
  fetch:
    inline:
      paths:

        packages/pkg.test.carvel.dev/pkg.test.carvel.dev.0.0.0.yml: |
          ---
          apiVersion: data.packaging.carvel.dev/v1alpha1
          kind: Package
          metadata:
            name: pkg.test.carvel.dev.0.0.0
          spec:
            refName: pkg.test.carvel.dev
            version: 0.0.0
            template:
              spec: {}
`, env.Namespace)

	installPkgYAML := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg.test.carvel.dev
  namespace: %[1]s
spec:
  # This is the name we want to reference in resources such as PackageInstall.
  displayName: "Test PackageMetadata in repo"
  shortDescription: "PackageMetadata used for testing"
  longDescription: "A longer, more detailed description of what the package contains and what it is for"
  providerName: Carvel
  maintainers:
  - name: carvel
  categories:
  - testing
  supportDescription: "Description of support provided for the package"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.test.carvel.dev.1.0.0
  namespace: %[1]s
spec:
  refName: pkg.test.carvel.dev
  version: 1.0.0
  licenses:
  - Apache 2.0
  capactiyRequirementsDescription: "cpu: 1,RAM: 2, Disk: 3"
  releaseNotes: |
    - Introduce simple-app package
  releasedAt: 2021-05-05T18:57:06Z
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: k8slt/kctrl-example-pkg:v1.0.0
      template:
      - ytt: {}
      - kbld:
          paths:
          - "-"
          - ".imgpkg/images.yml"
      deploy:
      - kapp:
          inspect: {}
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: instl-pkg-test
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0
  values:
  - secretRef:
      name: pkg-demo-values
---
apiVersion: v1
kind: Secret
metadata:
  name: pkg-demo-values
stringData:
  values.yml: |
    hello_msg: "hi"
`, env.Namespace) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", "simple-app.app"})
		kapp.Run([]string{"delete", "-a", "simple-app"})
		kapp.Run([]string{"delete", "-a", "default-ns-rbac"})
		kapp.Run([]string{"delete", "-a", "instl-pkg-test"})
		kapp.Run([]string{"delete", "-a", "minimal-repo.tanzu.carvel.dev"})
	}
	cleanUp()
	defer cleanUp()

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	// port-forwarding goroutine
	go func() {
		defer wg.Done()
		portForward(ctx)
	}()

	// Allow some time for port-forwarding
	time.Sleep(2 * time.Second)

	kapp.Run([]string{"deploy", "-a", "default-ns-rbac", "-f",
		"https://raw.githubusercontent.com/carvel-dev/kapp-controller/develop/examples/rbac/default-ns.yml"})

	kapp.Run([]string{"deploy", "-a", "simple-app", "-f",
		"https://raw.githubusercontent.com/k14s/kapp-controller/develop/examples/simple-app-git/1.yml"})

	kapp.RunWithOpts([]string{"deploy", "-a", "minimal-repo.tanzu.carvel.dev", "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(pkgRepoYAML)})

	kapp.RunWithOpts([]string{"deploy", "-a", "instl-pkg-test", "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(installPkgYAML)})

	t.Logf("Hitting URL")

	resp, err := http.Get("http://localhost:8080/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check if the response is successful
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)

	response := string(bodyBytes)

	bodyContains := assert.Contains(t, response, "kappctrl_reconcile_deploy_time_seconds") &&
		assert.Contains(t, response, "kappctrl_reconcile_fetch_time_seconds") &&
		assert.Contains(t, response, "kappctrl_reconcile_template_time_seconds") &&
		assert.Contains(t, response, "kappctrl_reconcile_time_seconds")

	assert.True(t, bodyContains)

	// Stop port-forwarding by canceling the context
	cancel()

	// Wait for the port-forwarding goroutine to complete
	wg.Wait()
}

func portForward(ctx context.Context) {
	cmd := exec.CommandContext(ctx, "kubectl", "port-forward", "svc/packaging-api", "8080", "-n", "kapp-controller")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.Canceled {
			fmt.Println("Port-forwarding stopped by context cancellation.")
		} else {
			fmt.Printf("Error running kubectl port-forward: %v\n", err)
		}
	}
}
