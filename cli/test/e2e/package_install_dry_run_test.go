package e2e

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageInstallDryRun(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	logger.Section("package install dry-run", func() {
		expectedOutput := `---
apiVersion: v1
kind: ServiceAccount
metadata:
  annotations:
    packaging.carvel.dev/package: test-kctrl-test
    tkg.tanzu.vmware.com/tanzu-package: test-kctrl-test
  creationTimestamp: null
  name: test-kctrl-test-sa
  namespace: kctrl-test
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    packaging.carvel.dev/package: test-kctrl-test
    tkg.tanzu.vmware.com/tanzu-package: test-kctrl-test
  creationTimestamp: null
  name: test-kctrl-test-cluster-role
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    packaging.carvel.dev/package: test-kctrl-test
    tkg.tanzu.vmware.com/tanzu-package: test-kctrl-test
  creationTimestamp: null
  name: test-kctrl-test-cluster-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: test-kctrl-test-cluster-role
subjects:
- kind: ServiceAccount
  name: test-kctrl-test-sa
  namespace: kctrl-test
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  annotations:
    packaging.carvel.dev/package-ClusterRole: test-kctrl-test-cluster-role
    packaging.carvel.dev/package-ClusterRoleBinding: test-kctrl-test-cluster-rolebinding
    packaging.carvel.dev/package-ServiceAccount: test-kctrl-test-sa
    tkg.tanzu.vmware.com/tanzu-package-ClusterRole: test-kctrl-test-cluster-role
    tkg.tanzu.vmware.com/tanzu-package-ClusterRoleBinding: test-kctrl-test-cluster-rolebinding
    tkg.tanzu.vmware.com/tanzu-package-ServiceAccount: test-kctrl-test-sa
  creationTimestamp: null
  name: test
  namespace: kctrl-test
spec:
  packageRef:
    refName: test.carvel.dev
    versionSelection:
      constraints: 1.0.0
      prereleases: {}
  serviceAccountName: test-kctrl-test-sa
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0`

		output := kappCtrl.Run([]string{"package", "install", "-i", "test", "-p", "test.carvel.dev", "--version", "1.0.0", "--dry-run"})
		require.Contains(t, output, expectedOutput)
	})

	logger.Section("package install dry-run with service account", func() {
		expectedOutput := `---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  creationTimestamp: null
  name: test
  namespace: kctrl-test
spec:
  packageRef:
    refName: test.carvel.dev
    versionSelection:
      constraints: 1.0.0
      prereleases: {}
  serviceAccountName: test-svc-acc
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0`

		output := kappCtrl.Run([]string{"package", "install", "-i", "test", "-p", "test.carvel.dev", "--version", "1.0.0",
			"--service-account-name", "test-svc-acc", "--dry-run"})
		require.Contains(t, output, expectedOutput)
	})

	logger.Section("package install dry-run with service account and values file", func() {
		expectedOutput := `---
apiVersion: v1
kind: Secret
metadata:
  annotations:
    packaging.carvel.dev/package: test-kctrl-test
    tkg.tanzu.vmware.com/tanzu-package: test-kctrl-test
  creationTimestamp: null
  name: test-kctrl-test-values
  namespace: kctrl-test
stringData:
  values.yaml: |
    foo: bar
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  annotations:
    packaging.carvel.dev/package-Secret: test-kctrl-test-values
    tkg.tanzu.vmware.com/tanzu-package-Secret: test-kctrl-test-values
  creationTimestamp: null
  name: test
  namespace: kctrl-test
spec:
  packageRef:
    refName: test.carvel.dev
    versionSelection:
      constraints: 1.0.0
      prereleases: {}
  serviceAccountName: test-svc-acc
  values:
  - secretRef:
      name: test-kctrl-test-values
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0`

		output, _ := kappCtrl.RunWithOpts([]string{"package", "install", "-i", "test", "-p", "test.carvel.dev", "--version", "1.0.0", "--values-file", "-",
			"--service-account-name", "test-svc-acc", "--dry-run"},
			RunOpts{StdinReader: strings.NewReader("foo: bar\n")})
		require.Contains(t, output, expectedOutput)
	})
}
