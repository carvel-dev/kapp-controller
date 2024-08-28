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
		tagExpectedOutput := `apiVersion: packaging.carvel.dev/v1alpha1
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
		semverExpectedOutput := `apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  creationTimestamp: null
  name: test-repo
  namespace: kctrl-test
spec:
  fetch:
    imgpkgBundle:
      image: registry.carvel.dev/project/repo
      tagSelection:
        semver:
          constraints: 1.0.0
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0
`
		tagSemverExpectedOutput := `apiVersion: packaging.carvel.dev/v1alpha1
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
  observedGeneration: 0
`

		tagOutput := kappCtrl.Run([]string{"package", "repo", "add", "-r", "test-repo", "--url",
			"registry.carvel.dev/project/repo:1.0.0", "--semver-tag-constraints", "1.0.0", "--dry-run"})
		semverOutput := kappCtrl.Run([]string{"package", "repo", "add", "-r", "test-repo", "--url",
			"registry.carvel.dev/project/repo", "--semver-tag-constraints", "1.0.0", "--dry-run"})
		tagSemverOutput := kappCtrl.Run([]string{"package", "repo", "add", "-r", "test-repo", "--url",
			"registry.carvel.dev/project/repo:1.0.0", "--semver-tag-constraints", "1.0.0", "--dry-run"})
		require.Contains(t, tagOutput, tagExpectedOutput)
		require.Contains(t, semverOutput, semverExpectedOutput)
		require.Contains(t, tagSemverOutput, tagSemverExpectedOutput)
	})
}

func TestPackageRepoSecretRefDryRun(t *testing.T) {
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
      secretRef:
        name: regcred
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0`

		output := kappCtrl.Run([]string{"package", "repo", "add", "-r", "test-repo", "--url",
			"registry.carvel.dev/project/repo:1.0.0", "--semver-tag-constraints", "1.0.0", "--secret-ref", "regcred", "--dry-run"})
		require.Contains(t, output, expectedOutput)
	})
}
