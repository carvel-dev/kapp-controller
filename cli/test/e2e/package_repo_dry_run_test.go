package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageRepoDryRun(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	logger.Section("dry-run package repo add", func() {
		expectedOutput := `apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  creationTimestamp: null
  name: test-repo
  namespace: kctrl-test
spec:
  fetch:
    imgpkgBundle:
      image: registry.carvel.dev/project/repo:1.0.0
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0`

		output := kappCtrl.Run([]string{"package", "repo", "add", "-r", "test-repo", "--url", "registry.carvel.dev/project/repo:1.0.0", "--dry-run"})
		require.Contains(t, output, expectedOutput)
	})
}
