// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ghodss/yaml"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_PackageInstalled_FromPackageInstall_Successfully(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}
	name := "instl-pkg-test"

	installPkgYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.test.carvel.dev
  namespace: %[1]s
spec:
  # This is the name we want to reference in resources such as PackageInstall.
  displayName: "Test Package in repo"
  shortDescription: "Package used for testing"
  longDescription: "A longer, more detailed description of what the package contains and what it is for"
  providerName: Carvel
  maintainers:
  - name: carvel
  categories:
  - testing
  supportDescription: "Description of support provided for the package"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageVersion
metadata:
  name: pkg.test.carvel.dev.1.0.0
  namespace: %[1]s
spec:
  packageName: pkg.test.carvel.dev
  version: 1.0.0
  licenses:
  - Apache 2.0
  capactiyRequirementsDescription: "cpu: 1,RAM: 2, Disk: 3"
  releaseNotes: |
    - Introduce simple-app package
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
      - kapp: {}
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %[2]s
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageVersionRef:
    packageName: pkg.test.carvel.dev
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
`, env.Namespace, name) + sas.ForNamespaceYAML()

	cleanUp := func() {
		// Delete App with kubectl first since kapp
		// deletes ServiceAccount before App
		kubectl.RunWithOpts([]string{"delete", "apps/" + name}, RunOpts{AllowError: true})
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	// Create Repo, PackageInstall, and App from YAML
	kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, RunOpts{StdinReader: strings.NewReader(installPkgYaml)})

	kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "pkgi/" + name, "--timeout", "1m"})
	kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "apps/" + name, "--timeout", "1m"})
	out := kubectl.Run([]string{"get", fmt.Sprintf("apps/%s", name), "-o", "yaml"})

	var cr v1alpha1.App
	err := yaml.Unmarshal([]byte(out), &cr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}

	expectedStatus := v1alpha1.AppStatus{
		GenericStatus: v1alpha1.GenericStatus{
			Conditions: []v1alpha1.AppCondition{{
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
		cr.Status.Deploy.StartedAt = metav1.Time{}
		cr.Status.Deploy.UpdatedAt = metav1.Time{}
		cr.Status.Deploy.Stdout = ""

		// fetch
		cr.Status.Fetch.StartedAt = metav1.Time{}
		cr.Status.Fetch.UpdatedAt = metav1.Time{}
		cr.Status.Fetch.Stdout = ""

		if !strings.Contains(cr.Status.Inspect.Stdout, "simple-app") && !strings.Contains(cr.Status.Inspect.Stdout, "Succeeded") {
			t.Fatalf("Expected to find simple-app resources created but got:\n%s", cr.Status.Inspect.Stdout)
		}
		cr.Status.Inspect.UpdatedAt = metav1.Time{}
		cr.Status.Inspect.Stdout = ""

		// template
		cr.Status.Template.UpdatedAt = metav1.Time{}
		cr.Status.Template.Stderr = ""
	}

	if !reflect.DeepEqual(expectedStatus, cr.Status) {
		t.Fatalf("\nStatus is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedStatus, cr.Status)
	}
}

func Test_PackageInstallStatus_DisplaysUsefulErrorMessage_ForDeploymentFailure(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}
	name := "instl-pkg-test-fail"

	installPkgYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.fail.carvel.dev
spec:
  displayName: "Test Package in repo"
  shortDescription: "Package used for testing"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageVersion
metadata:
  name: pkg.fail.carvel.dev.1.0.0
spec:
  packageName: pkg.fail.carvel.dev
  version: 1.0.0
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: k8slt/kctrl-example-pkg:v1.0.0
      template:
      - ytt:
          paths:
          - "config.yml"
          - "values.yml"
      - kbld:
          paths:
          - "-"
          - ".imgpkg/images.yml"
      deploy:
      - kapp:
          # this is done intentionally for testing
          intoNs: does-not-exist
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %s
  namespace: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageVersionRef:
    packageName: pkg.fail.carvel.dev
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
`, name, env.Namespace) + sas.ForNamespaceYAML()

	cleanUp := func() {
		// Delete App with kubectl first since kapp
		// deletes ServiceAccount before App
		kubectl.RunWithOpts([]string{"delete", "apps/" + name}, RunOpts{AllowError: true})
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	// Create Repo, PackageInstall, and App from YAML
	kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, RunOpts{StdinReader: strings.NewReader(installPkgYaml)})

	// wait for status to update for PackageInstall
	var cr pkgingv1alpha1.PackageInstall
	retry(t, 30*time.Second, func() error {
		out := kubectl.Run([]string{"get", fmt.Sprintf("pkgi/%s", name), "-o", "yaml"})
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal: %s", err)
		}

		if !strings.Contains(cr.Status.UsefulErrorMessage, "kapp: Error") {
			return fmt.Errorf("\nExpected useful error message to contain deploy error\nGot:\n%s", cr.Status.UsefulErrorMessage)
		}

		if !strings.Contains(cr.Status.FriendlyDescription, "Error (see .status.usefulErrorMessage for details)") {
			return fmt.Errorf("\nExpected friendly description to contain error\nGot:\n%s", cr.Status.FriendlyDescription)
		}

		return err
	})
}

func Test_PackageInstalled_FromPackageInstall_DeletionFailureBlocks(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}
	name := "instl-pkg-failure-block-test"

	// contents of this bundle (k8slt/kappctrl-e2e-repo)
	// under examples/packaging-demo
	installPkgYaml := fmt.Sprintf(`---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: basic.test.carvel.dev
  # cluster scoped
spec:
  fetch:
    imgpkgBundle:
      image: index.docker.io/k8slt/kc-e2e-test-repo@sha256:0ae0f32ef92d2362339b47055a6ea2042bc114a7dd36cf339bf05df4d1cc1b9b
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %s
  namespace: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageVersionRef:
    packageName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0
`, name, env.Namespace) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("Install package", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
			RunOpts{StdinReader: strings.NewReader(installPkgYaml)})

		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "pkgi/" + name, "--timeout", "1m"})
	})

	logger.Section("Delete service account so that PackageInstall deletion would fail", func() {
		kubectl.Run([]string{"delete", "serviceaccount", "kappctrl-e2e-ns-sa"})
	})

	logger.Section("Check that deletion of PackageInstall results in failure conditions", func() {
		// No waiting for deletion since it's blocked
		kubectl.Run([]string{"delete", "pkgi", name, "--wait=false"})
		kubectl.Run([]string{"wait", "--for=condition=DeleteFailed", "pkgi", name, "--timeout", "1m"})
	})

	logger.Section("Bring back service account and see that kubectl delete succeeds", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-", "--filter-kind", "ServiceAccount"},
			RunOpts{StdinReader: strings.NewReader(installPkgYaml)})

		kubectl.Run([]string{"delete", "pkgi", name, "--wait=true"})
	})
}
