// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1"
	"carvel.dev/kapp-controller/test/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func Test_PackageInstalled_FromPackageInstall_Successfully(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}
	name := "instl-pkg-test"

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
  name: %[2]s
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
`, env.Namespace, name) + sas.ForNamespaceYAML()

	cleanUp := func() {
		// Delete App with kubectl first since kapp
		// deletes ServiceAccount before App
		kubectl.RunWithOpts([]string{"delete", "apps/" + name}, e2e.RunOpts{AllowError: true})
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	// Create Repo, PackageInstall, and App from YAML
	kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(installPkgYAML)})

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
		cr.Status.Deploy.StartedAt = metav1.Time{}
		cr.Status.Deploy.UpdatedAt = metav1.Time{}
		cr.Status.Deploy.Stdout = ""
		cr.Status.Deploy.KappDeployStatus = nil

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

	assert.Equal(t, expectedStatus, cr.Status)
}

// helpers for the below tests allows us to "template" in and out version constraints
// if kcVerC is "" then we omit that field, same for k8sVerC
// annotations will be added to the PKGI, and are expected to be key-names only (vals will be "")
func generatePkgYAML(kcVerC string, k8sVerC string, pkgVer string) string {
	pkgYAML := fmt.Sprintf(`
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.test.carvel.dev.%[1]s
spec:
`, pkgVer)
	// too bad there's no tools for templating yaml.
	if kcVerC != "" {
		pkgYAML += fmt.Sprintf(`  kappControllerVersionSelection:
    constraints: "%s"
`, kcVerC)
	}
	if k8sVerC != "" {
		pkgYAML += fmt.Sprintf(`  kubernetesVersionSelection:
    constraints: "%s"
`, k8sVerC)
	}
	pkgYAML += fmt.Sprintf(`  refName: pkg.test.carvel.dev
  version: %[1]s
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
`, pkgVer)
	return pkgYAML
}

func generatePkgiYAML(pkgiName, kcVerC, k8sVerC string, annotations []string) string {
	installPkgYAML := generatePkgYAML(kcVerC, k8sVerC, "1.0.0") + fmt.Sprintf(`
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %[1]s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
`, pkgiName)
	if len(annotations) > 0 {
		for _, ann := range annotations {
			installPkgYAML += fmt.Sprintf(`    "%s": ""`, ann)
		}
	}
	installPkgYAML += `
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: ">0.0.0"
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
`
	return installPkgYAML
}

func Test_PackageInstalled_FromPackageInstall_VersionConstraints(t *testing.T) {
	const (
		kubernetesVersionOverrideAnnotation     = "packaging.carvel.dev/ignore-kubernetes-version-selection"
		kappControllerVersionOverrideAnnotation = "packaging.carvel.dev/ignore-kapp-controller-version-selection"
	)

	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}
	name := "instl-pkg-test-version-constraints"

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("PackageInstall fails due to k8s version constraint", func() {
		installPkgYAML := generatePkgiYAML(name, "", "<1.0.0", []string{}) + sas.ForNamespaceYAML()
		out, err := kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
			e2e.RunOpts{
				StdinReader: strings.NewReader(installPkgYAML),
				AllowError:  true})
		assert.Contains(t, out, "after-kubernetes-version-check=0")
		require.Error(t, err)
		// Without calling cleanup after each stage, somehow it seems like the controller doesn't
		// get the updated annotation.
		cleanUp()
	})

	logger.Section("PackageInstall succeeds with a kubernetes version override annotation", func() {
		installPkgYAML := generatePkgiYAML(name, "", "<1.0.0",
			[]string{kubernetesVersionOverrideAnnotation}) + sas.ForNamespaceYAML()
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(installPkgYAML)})
		cleanUp()
	})

	logger.Section("PackageInstall succeeds with a kubernetes version within a range", func() {
		installPkgYAML := generatePkgiYAML(name, "", ">0.0.0 <10.0.0", []string{}) + sas.ForNamespaceYAML()
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(installPkgYAML)})
		cleanUp()
	})

	logger.Section("PackageInstall fails due to kapp-controller version constraint", func() {
		installPkgYAML := generatePkgiYAML(name, ">=3.0.0", "", []string{}) + sas.ForNamespaceYAML()
		out, err := kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
			e2e.RunOpts{
				StdinReader: strings.NewReader(installPkgYAML),
				AllowError:  true})
		assert.Contains(t, out, "after-kapp-controller-version-check=0")
		require.Error(t, err)
		cleanUp()
	})

	logger.Section("PackageInstall succeeds with a kapp-controller version override annotation", func() {
		installPkgYAML := generatePkgiYAML(name, ">=3.0.0", "",
			[]string{kappControllerVersionOverrideAnnotation}) + sas.ForNamespaceYAML()
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(installPkgYAML)})
		cleanUp()
	})

	logger.Section("PackageInstall succeeds with a kapp-controller version within a range", func() {
		installPkgYAML := generatePkgiYAML(name, ">0.0.0 <10.0.0", "", []string{}) + sas.ForNamespaceYAML()
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(installPkgYAML)})
	})
}

func Test_PackageInstallStatus_DisplaysUsefulErrorMessage_ForDeploymentFailure(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}
	name := "instl-pkg-test-fail"

	installPkgYAML := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg.fail.carvel.dev
  annotations:
    kapp.k14s.io/change-group: "package"
spec:
  displayName: "Test PackageMetadata in repo"
  shortDescription: "PackageMetadata used for testing"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.fail.carvel.dev.1.0.0
  annotations:
    kapp.k14s.io/change-group: "package"
spec:
  refName: pkg.fail.carvel.dev
  version: 1.0.0
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: k8slt/kctrl-example-pkg:v1.0.0
      template:
      - ytt:
          paths:
          - config/
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
    kapp.k14s.io/change-rule: "upsert after upserting package"
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.fail.carvel.dev
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
		kubectl.RunWithOpts([]string{"delete", "apps/" + name}, e2e.RunOpts{AllowError: true})
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	// Create Repo, PackageInstall, and App from YAML
	kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
		e2e.RunOpts{StdinReader: strings.NewReader(installPkgYAML), AllowError: true})

	var cr pkgingv1alpha1.PackageInstall
	out := kubectl.Run([]string{"get", fmt.Sprintf("pkgi/%s", name), "-o", "yaml"})
	err := yaml.Unmarshal([]byte(out), &cr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}

	if !strings.Contains(cr.Status.UsefulErrorMessage, "kapp: Error") {
		t.Fatalf("\nExpected useful error message to contain deploy error\nGot:\n%s", cr.Status.UsefulErrorMessage)
	}

	if !strings.Contains(cr.Status.FriendlyDescription, "Error (see .status.usefulErrorMessage for details)") {
		t.Fatalf("\nExpected friendly description to contain error\nGot:\n%s", cr.Status.FriendlyDescription)
	}
}

func Test_PackageInstalled_FromPackageInstall_DeletionFailureBlocks(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}
	name := "instl-pkg-failure-block-test"

	// contents of this bundle (k8slt/kappctrl-e2e-repo)
	// under examples/packaging-demo
	installPkgYAML := fmt.Sprintf(`---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: basic.test.carvel.dev
  annotations:
    kapp.k14s.io/change-group: "packagerepo"
spec:
  fetch:
    imgpkgBundle:
      image: index.docker.io/k8slt/kc-e2e-test-repo@sha256:ddd93b67b97c1460580ca1afd04326d16900dc716c4357cade85b83deab76f1c
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %s
  namespace: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
    kapp.k14s.io/change-rule: "upsert after upserting packagerepo"
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0
`, name, env.Namespace) + sas.ForNamespaceYAML()

	cleanUp := func() {
		// Need to recreate ServiceAccount in event test fails
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-", "--filter-kind", "ServiceAccount"},
			e2e.RunOpts{StdinReader: strings.NewReader(installPkgYAML)})
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("Create PackageRepository and PackageInstall", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(installPkgYAML)})
	})

	logger.Section("Delete service account so that PackageInstall deletion would fail", func() {
		kubectl.Run([]string{"delete", "serviceaccount", "kappctrl-e2e-ns-sa"})
	})

	logger.Section("Check that deletion of PackageInstall results in failure conditions", func() {
		// No waiting for deletion since it's blocked
		kubectl.Run([]string{"delete", "pkgi", name, "--wait=false"})
		for i := 1; i < 33; i += i {
			out := kubectl.Run([]string{"get", "pkgi", name})
			// we expected for the delete to fail, so once we see this we're done
			if strings.Contains(out, "Delete failed: Error (see .status.usefulErrorMessage for details)") {
				break
			}
			// succeeded is the state the resource was previously in, so it's ok if it's in this state transiently, but no other state is expected/acceptable
			if !strings.Contains(out, "Reconcile succeeded") {
				fmt.Println("got unexpected Description: \n", out)
				fmt.Println(kubectl.Run([]string{"get", "pkgi", name, "-oyaml"}))
				break
			}
			time.Sleep(time.Duration(i) * time.Second)
		}
	})
	// this kubectl Run functions as a test assertion.
	kubectl.Run([]string{"wait", "--for=condition=DeleteFailed", "pkgi", name, "--timeout", "1s"})

	logger.Section("Bring back service account and see that kubectl delete succeeds", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-", "--filter-kind", "ServiceAccount"},
			e2e.RunOpts{StdinReader: strings.NewReader(installPkgYAML)})
		kubectl.Run([]string{"delete", "pkgi", name})
	})
}

func Test_PackageInstall_UsesExistingAppWithSameName(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}
	name := "pkg-instl-uses-app"

	appYAML := fmt.Sprintf(`
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
  template:
  - ytt: {}
  deploy:
    - kapp: {}
`, name) + sas.ForNamespaceYAML()

	pkginstallYAML := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg.test.carvel.dev
  namespace: %[1]s
spec:
  displayName: "Test PackageMetadata"
  shortDescription: "PackageMetadata used for testing"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.test.carvel.dev.1.0.0
  namespace: %[1]s
spec:
  refName: pkg.test.carvel.dev
  version: 1.0.0
  template:
    spec:
      fetch:
      - inline:
          paths:
            file.yml: |
              apiVersion: v1
              kind: ConfigMap
              metadata:
                name: configmap
      template:
      - ytt: 
          paths:
          - file.yml
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
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0
`, env.Namespace, name) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("Create App CR", func() {
		kubectl.RunWithOpts([]string{"apply", "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(appYAML)})
		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "app/" + name, "--timeout", "1m"})
	})

	logger.Section("Create PackageInstall with same name as App CR", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(pkginstallYAML)})
	})

	logger.Section("Assert that App spec is different from PackageInstall take over", func() {
		out := kubectl.Run([]string{"get", fmt.Sprintf("apps/%s", name), "-o", "yaml"})
		var cr v1alpha1.App
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("failed to unmarshal: %s", err)
		}

		tmpl := cr.Spec.Template[0]
		if tmpl.Ytt != nil && len(tmpl.Ytt.Paths) != 0 {
			if tmpl.Ytt.Paths[0] != "file.yml" {
				t.Fatalf("\nExpected App spec.template.ytt.paths to contain file.yml\nGot: %s", tmpl.Ytt.Paths[0])
			}
		} else {
			t.Fatalf("\nExpected App spec.template.ytt.paths to contain file.yml\nGot: %s", tmpl)
		}
	})

	logger.Section("Delete PackageInstall and expect App with same name to be deleted", func() {
		cleanUp()
		// Since the App was created first without the PackageInstall, the result of the
		// PackageInstall using the existing App should be that it will be deleted when
		// PackageInstall gets deleted.
		out, err := kubectl.RunWithOpts([]string{"get", fmt.Sprintf("apps/%s", name)}, e2e.RunOpts{AllowError: true})
		if err == nil {
			t.Fatalf("Expected no App to be found after PackageInstall is deleted\nGot: %s", out)
		}

		if !strings.Contains(err.Error(), "NotFound") {
			t.Fatalf("Expected error from kubectl get app to show App not found.\nGot: %s", err.Error())
		}
	})
}

func Test_PackageInstall_UpgradesToNewVersion_Successfully(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	name := "instl-pkg-upgrade-test"

	cleanUp := func() {
		// Delete App with kubectl first since kapp
		// deletes ServiceAccount before App
		kubectl.RunWithOpts([]string{"delete", "apps/" + name}, e2e.RunOpts{AllowError: true})
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("Create PackageInstall using version Package version 1.0.0", func() {
		pkgInstallYAML := packageInstallVersionInYAML(name, env.Namespace, "1.0.0")
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(pkgInstallYAML)})
	})

	logger.Section("Check PackageInstall with version 1.0.0 success", func() {
		out := kubectl.Run([]string{"get", fmt.Sprintf("pkgi/%s", name), "-o", "yaml"})

		var cr pkgingv1alpha1.PackageInstall
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		expectedStatus := pkgingv1alpha1.PackageInstallStatus{
			GenericStatus: v1alpha1.GenericStatus{
				Conditions: []v1alpha1.Condition{{
					Type:   v1alpha1.ReconcileSucceeded,
					Status: corev1.ConditionTrue,
				}},
				ObservedGeneration:  1,
				FriendlyDescription: "Reconcile succeeded",
			},
			Version:              "1.0.0",
			LastAttemptedVersion: "1.0.0",
		}

		if !reflect.DeepEqual(expectedStatus, cr.Status) {
			t.Fatalf("\nStatus is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedStatus, cr.Status)
		}
	})

	logger.Section("Create PackageInstall using version Package version 2.0.0", func() {
		pkgInstallYAML := packageInstallVersionInYAML(name, env.Namespace, "2.0.0")
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(pkgInstallYAML)})
	})

	logger.Section("Check PackageInstall with version 2.0.0 success", func() {
		outPkgi := kubectl.Run([]string{"get", fmt.Sprintf("pkgi/%s", name), "-o", "yaml"})

		var crPkgi pkgingv1alpha1.PackageInstall
		err := yaml.Unmarshal([]byte(outPkgi), &crPkgi)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		expectedStatus := pkgingv1alpha1.PackageInstallStatus{
			GenericStatus: v1alpha1.GenericStatus{
				Conditions: []v1alpha1.Condition{{
					Type:   v1alpha1.ReconcileSucceeded,
					Status: corev1.ConditionTrue,
				}},
				ObservedGeneration:  2,
				FriendlyDescription: "Reconcile succeeded",
			},
			Version:              "2.0.0",
			LastAttemptedVersion: "2.0.0",
		}

		if !reflect.DeepEqual(expectedStatus, crPkgi.Status) {
			t.Fatalf("\nStatus is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedStatus, crPkgi.Status)
		}

		outCm := kubectl.Run([]string{"get", "configmap/configmap-version", "-o", "yaml"})
		var cm corev1.ConfigMap
		err = yaml.Unmarshal([]byte(outCm), &cm)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		if cm.Data["version"] != "2.0.0" {
			t.Fatalf("Expected Package ConfigMap data to use version 2.0.0.\nGot:\n%s", cm.Data["version"])
		}
	})
}

func packageInstallVersionInYAML(name, namespace, version string) string {
	sas := e2e.ServiceAccounts{namespace}

	return fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
 name: pkg.carvel.dev
spec:
 displayName: "Test PackageMetadata"
 shortDescription: "PackageMetadata used for testing"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
 name: pkg.carvel.dev.1.0.0
spec:
 refName: pkg.carvel.dev
 version: 1.0.0
 template:
   spec:
     fetch:
     - inline:
         paths:
           file.yml: |
             apiVersion: v1
             kind: ConfigMap
             metadata:
               name: configmap-version
             data:
               version: 1.0.0
     template:
     - ytt: {}
     deploy:
     - kapp: {}
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
 name: pkg.carvel.dev.2.0.0
spec:
 refName: pkg.carvel.dev
 version: 2.0.0
 template:
   spec:
     fetch:
     - inline:
         paths:
           file.yml: |
             apiVersion: v1
             kind: ConfigMap
             metadata:
               name: configmap-version
             data:
               version: 2.0.0
     template:
     - ytt: {}
     deploy:
     - kapp: {}
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
 name: %[1]s
 annotations:
   kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
spec:
 serviceAccountName: kappctrl-e2e-ns-sa
 packageRef:
   refName: pkg.carvel.dev
   versionSelection:
     constraints: %[2]s
`, name, version) + sas.ForNamespaceYAML()
}

func Test_PackageInstall_DefaultSyncPeriod(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{T: t, Namespace: env.Namespace, L: logger}
	sas := e2e.ServiceAccounts{Namespace: env.Namespace}
	kubectl := e2e.Kubectl{T: t, Namespace: env.Namespace, L: logger}

	name := "install-pkg-default-sync-period-test"

	packageInstallYaml := fmt.Sprintf(`
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: basic.test.carvel.dev
  annotations:
    kapp.k14s.io/change-group: "packagerepo"
spec:
  fetch:
    imgpkgBundle:
      image: index.docker.io/k8slt/kc-e2e-test-repo@sha256:ddd93b67b97c1460580ca1afd04326d16900dc716c4357cade85b83deab76f1c
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %s
  namespace: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
    kapp.k14s.io/change-rule: "upsert after upserting packagerepo"
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0`, name, env.Namespace) + sas.ForNamespaceYAML()

	cleanUpApp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUpApp()
	defer cleanUpApp()

	logger.Section("deploy PackageInstall and App created has default syncPeriod set", func() {

		// Create Repo and PackageInstall from YAML
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(packageInstallYaml)})

		// syncPeriod of the App created via PackageInstall should match to default
		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "pkgi/" + name, "--timeout", "1m"})
		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "apps/" + name, "--timeout", "1m"})
		out := kubectl.Run([]string{"get", fmt.Sprintf("apps/%s", name), "-o", "yaml"})

		var cr v1alpha1.App
		err := yaml.Unmarshal([]byte(out), &cr)

		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		expectedSyncPeriod := time.Duration(10) * time.Minute
		actualSyncPeriod := cr.Spec.SyncPeriod.Duration

		assert.Equal(t, actualSyncPeriod, expectedSyncPeriod)
	})
}

func Test_PackageInstall_CustomSyncPeriod(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{T: t, Namespace: env.Namespace, L: logger}
	sas := e2e.ServiceAccounts{Namespace: env.Namespace}
	kubectl := e2e.Kubectl{T: t, Namespace: env.Namespace, L: logger}

	name := "install-pkg-custom-sync-period-test"

	packageInstallYaml := fmt.Sprintf(`
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: basic.test.carvel.dev
  annotations:
    kapp.k14s.io/change-group: "packagerepo"
spec:
  fetch:
    imgpkgBundle:
      image: index.docker.io/k8slt/kc-e2e-test-repo@sha256:ddd93b67b97c1460580ca1afd04326d16900dc716c4357cade85b83deab76f1c
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %s
  namespace: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
    kapp.k14s.io/change-rule: "upsert after upserting packagerepo"
spec:
  syncPeriod: 30s
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0`, name, env.Namespace) + sas.ForNamespaceYAML()

	cleanUpApp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUpApp()
	defer cleanUpApp()

	logger.Section("deploy PackageInstall and App created has default syncPeriod set", func() {

		// Create Repo and PackageInstall from YAML
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(packageInstallYaml)})

		// syncPeriod of the App created via PackageInstall should match to configured value part of PackageInstall.
		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "pkgi/" + name, "--timeout", "1m"})
		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "apps/" + name, "--timeout", "1m"})
		out := kubectl.Run([]string{"get", fmt.Sprintf("apps/%s", name), "-o", "yaml"})

		var cr v1alpha1.App
		err := yaml.Unmarshal([]byte(out), &cr)

		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		expectedSyncPeriod := time.Duration(30) * time.Second
		actualSyncPeriod := cr.Spec.SyncPeriod.Duration

		assert.Equal(t, actualSyncPeriod, expectedSyncPeriod)
	})
}
