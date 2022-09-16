// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func Test_FetchAndDeployImgpkgBundle_Successfully(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	// contents for kappctrl-e2e-bundle
	// available in test/e2e/assets/kappctrl-e2e-bundle
	appYaml := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-bundle-app
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
    - imgpkgBundle:
        image: k8slt/kappctrl-e2e-bundle
  template:
  - kbld:
      paths:
      - .imgpkg/images.yml
      - config.yml
  deploy:
  - kapp:
      inspect: {}
      intoNs: %s
`, env.Namespace) + sas.ForNamespaceYAML()

	name := "test-bundle-app"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(appYaml)})

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

		var cr v1alpha1.App
		err := yaml.Unmarshal([]byte(out), &cr)
		require.NoError(t, err)

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
			require.Contains(t, cr.Status.Fetch.Stdout, "- imgpkgBundle")
			require.Contains(t, cr.Status.Fetch.Stdout, "image: index.docker.io/k8slt/kappctrl-e2e-bundle@sha256:83f86234f68a980490ec66f2d347ad4c8148713073c0993760b8eaaef3eb48d7")
			cr.Status.Fetch.StartedAt = metav1.Time{}
			cr.Status.Fetch.UpdatedAt = metav1.Time{}
			cr.Status.Fetch.Stdout = ""

			// inspect
			require.Contains(t, cr.Status.Inspect.Stdout, "simple-app")
			require.Contains(t, cr.Status.Inspect.Stdout, "Succeeded")
			cr.Status.Inspect.UpdatedAt = metav1.Time{}
			cr.Status.Inspect.Stdout = ""

			// template
			cr.Status.Template.UpdatedAt = metav1.Time{}
			cr.Status.Template.Stderr = ""
		}

		require.Equal(t, expectedStatus, cr.Status)
	})
}

func Test_AppCR_FetchFromCaching(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	pkgiName := "pkgi-to-test-caching"
	registryNamespace := "registry"
	registryName := "test-registry"

	cleanUp := func() {
		// Delete App with kubectl first since kapp
		// deletes ServiceAccount before App
		kubectl.RunWithOpts([]string{"delete", "apps/" + pkgiName}, e2e.RunOpts{AllowError: true})
		kapp.Run([]string{"delete", "-a", pkgiName})
		kapp.Run([]string{"delete", "-a", registryName, "-n", registryNamespace})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy registry without auth", func() {
		kapp.Run([]string{
			"deploy", "-a", registryName,
			"-f", "../assets/registry/no-auth-registry.yml",
			"-f", "../assets/registry/registry-contents.yml",
		})
	})

	var initialFetchUpdate metav1.Time
	logger.Section("deploy PackageInstall", func() {
		pkgiYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.test.carvel.dev.1.0.0
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: "package"
spec:
  refName: pkg.test.carvel.dev
  version: 1.0.0
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: registry-svc.%[3]s.svc.cluster.local:5050/secret-test/test-bundle@sha256:387e2fa04e3c9d7b89f03be4fd3fb661fa3496352ece38a7fc8cba6780c7131d
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
    kapp.k14s.io/change-rule: "upsert after upserting package"
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0
`, env.Namespace, pkgiName, registryNamespace) + sas.ForNamespaceYAML()

		kapp.RunWithOpts([]string{"deploy", "-a", pkgiName, "-f", "-"}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgiYaml),
			OnErrKubectl: []string{"get", "app", pkgiName, "-oyaml"},
		})

		kubectl.Run([]string{"get", "configmap", "e2e-test-map"})
		output := kubectl.Run([]string{"get", "app", pkgiName, "-oyaml"})
		var appcr v1alpha1.App
		require.NoError(t, yaml.Unmarshal([]byte(output), &appcr))
		require.Equal(t, 0, appcr.Status.Fetch.ExitCode)
		initialFetchUpdate = appcr.Status.Fetch.UpdatedAt
	})
	defer kapp.Run([]string{"delete", "-a", pkgiName})

	logger.Section("remove registry", func() {
		kapp.Run([]string{"delete", "-a", registryName, "-n", registryNamespace})
	})

	logger.Section("wait for reconciliation and ensure no error occur while fetching", func() {
		pkgiYaml := fmt.Sprintf(`---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %[2]s
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
    kapp.k14s.io/change-rule: "upsert after upserting package"
spec:
  paused: true
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0
`, env.Namespace, pkgiName, registryNamespace) + sas.ForNamespaceYAML()

		kapp.RunWithOpts([]string{"deploy", "-a", pkgiName, "-f", "-", "-p"}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgiYaml),
			OnErrKubectl: []string{"get", "app", pkgiName, "-oyaml"},
		})

		pkgiYaml = fmt.Sprintf(`
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %[2]s
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
    kapp.k14s.io/change-rule: "upsert after upserting package"
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0
`, env.Namespace, pkgiName, registryNamespace) + sas.ForNamespaceYAML()

		kapp.RunWithOpts([]string{"deploy", "-a", pkgiName, "-f", "-", "-p"}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgiYaml),
			OnErrKubectl: []string{"get", "app", pkgiName, "-oyaml"},
		})
		waitForPKGIReconciliationOrFail(t, kubectl, pkgiName, initialFetchUpdate)

		var appcr v1alpha1.App
		output := kubectl.Run([]string{"get", "app", pkgiName, "-oyaml"})
		require.NoError(t, yaml.Unmarshal([]byte(output), &appcr))
		require.Equal(t, 0, appcr.Status.Fetch.ExitCode)
	})
}

func Test_PackageRepo_FetchFromCaching(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	pkgrName := "pkgr-to-test-caching"
	registryNamespace := "registry"
	registryName := "test-registry"

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", pkgrName})
		kapp.Run([]string{"delete", "-a", registryName, "-n", registryNamespace})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy registry without authentication", func() {
		kapp.Run([]string{
			"deploy", "-a", registryName,
			"-f", "../assets/registry/no-auth-registry.yml",
			"-f", "../assets/registry/registry-contents.yml",
		})
	})

	logger.Section("deploy PackageRepository that auths to registry", func() {
		pkgrYaml := fmt.Sprintf(`---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %[2]s
  namespace: %[1]s
spec:
  syncPeriod: 30s
  fetch:
    image:
      url: registry-svc.%[3]s.svc.cluster.local:5050/secret-test/test-repo@sha256:5deb3a248b2024da70b7a34bf561c91b23d2810bee5aa666d0fba9f10e26c273
`, env.Namespace, pkgrName, registryNamespace) + sas.ForNamespaceYAML()

		kapp.RunWithOpts([]string{"deploy", "-a", pkgrName, "-f", "-"}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgrYaml),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
	})

	logger.Section("remove registry", func() {
		kapp.Run([]string{"delete", "-a", registryName, "-n", registryNamespace})
	})

	logger.Section("pause/unpause package repository and check it does not fail", func() {
		pkgrYaml := fmt.Sprintf(`---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %[2]s
  namespace: %[1]s
spec:
  paused: true
  fetch:
    image:
      url: registry-svc.%[3]s.svc.cluster.local:5050/secret-test/test-repo@sha256:5deb3a248b2024da70b7a34bf561c91b23d2810bee5aa666d0fba9f10e26c273
`, env.Namespace, pkgrName, registryNamespace) + sas.ForNamespaceYAML()

		kapp.RunWithOpts([]string{"deploy", "-a", pkgrName, "-f", "-"}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgrYaml),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})

		pkgrYaml = fmt.Sprintf(`---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %[2]s
  namespace: %[1]s
spec:
  fetch:
    image:
      url: registry-svc.%[3]s.svc.cluster.local:5050/secret-test/test-repo@sha256:5deb3a248b2024da70b7a34bf561c91b23d2810bee5aa666d0fba9f10e26c273
`, env.Namespace, pkgrName, registryNamespace) + sas.ForNamespaceYAML()

		kapp.RunWithOpts([]string{"deploy", "-a", pkgrName, "-f", "-"}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgrYaml),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})

		var pkgr v1alpha12.PackageRepository
		output := kubectl.Run([]string{"get", "packagerepository", pkgrName, "-oyaml"})
		require.NoError(t, yaml.Unmarshal([]byte(output), &pkgr))
		require.Equal(t, int64(3), pkgr.Status.ObservedGeneration)
		require.Equal(t, 0, pkgr.Status.Fetch.ExitCode)
	})
}

func Test_PackageRepo_DoesNotFetchFromCaching(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	pkgrName := "pkgr-to-test-caching"
	registryNamespace := "registry"
	registryName := "test-registry"

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", pkgrName})
		kapp.Run([]string{"delete", "-a", registryName, "-n", registryNamespace})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy registry without authentication", func() {
		kapp.Run([]string{
			"deploy", "-a", registryName,
			"-f", "../assets/registry/no-auth-registry.yml",
			"-f", "../assets/registry/registry-contents.yml",
		})
	})

	logger.Section("deploy PackageRepository that auths to registry", func() {
		pkgrYaml := fmt.Sprintf(`---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %[2]s
  namespace: %[1]s
spec:
  syncPeriod: 30s
  fetch:
    image:
      url: registry-svc.%[3]s.svc.cluster.local:5050/secret-test/test-repo
`, env.Namespace, pkgrName, registryNamespace) + sas.ForNamespaceYAML()

		kapp.RunWithOpts([]string{"deploy", "-a", pkgrName, "-f", "-"}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgrYaml),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})

		kubectl.Run([]string{"get", "package", "pkg.test.carvel.dev.1.0.0"})
	})

	logger.Section("remove registry", func() {
		kapp.Run([]string{"delete", "-a", registryName, "-n", registryNamespace})
	})

	logger.Section("pause/unpause and check that it fails reconciling the package repository", func() {
		pkgrYaml := fmt.Sprintf(`---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %[2]s
  namespace: %[1]s
spec:
  paused: true
  fetch:
    image:
      url: registry-svc.%[3]s.svc.cluster.local:5050/secret-test/test-repo
`, env.Namespace, pkgrName, registryNamespace) + sas.ForNamespaceYAML()

		kapp.RunWithOpts([]string{"deploy", "-a", pkgrName, "-f", "-"}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgrYaml),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		pkgrYaml = fmt.Sprintf(`---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %[2]s
  namespace: %[1]s
spec:
  fetch:
    image:
      url: registry-svc.%[3]s.svc.cluster.local:5050/secret-test/test-repo
`, env.Namespace, pkgrName, registryNamespace) + sas.ForNamespaceYAML()

		_, err := kapp.RunWithOpts([]string{"deploy", "-a", pkgrName, "-f", "-"}, e2e.RunOpts{
			StdinReader: strings.NewReader(pkgrYaml),
			AllowError:  true,
		})

		require.Error(t, err, "Expected fetching error")
	})
}

func waitForPKGIReconciliationOrFail(t *testing.T, kubectl e2e.Kubectl, pkgiName string, initialStartFetch metav1.Time) {
	waitForStatusReconciliationOrFail(t, func(t *testing.T) *v1alpha1.AppStatusFetch {
		output := kubectl.Run([]string{"get", "app", pkgiName, "-oyaml"})
		var appcr v1alpha1.App

		require.NoError(t, yaml.Unmarshal([]byte(output), &appcr))
		return appcr.Status.Fetch
	}, initialStartFetch)
}

func waitForStatusReconciliationOrFail(t *testing.T, getStatus func(t *testing.T) *v1alpha1.AppStatusFetch, initialStartFetch metav1.Time) {
	for i := 0; i < 90; i++ {
		fetchStatus := getStatus(t)
		if fetchStatus.UpdatedAt != initialStartFetch {
			break
		}
		if i == 89 {
			fmt.Printf("%+v\n", fetchStatus)
		}
		require.Less(t, i, 89, "waited too much")
		time.Sleep(1 * time.Second)
	}
}
