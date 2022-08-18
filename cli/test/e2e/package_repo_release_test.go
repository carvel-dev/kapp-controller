package e2e

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	packagesDir = "packages"
)

func TestPackageRepositoryReleaseInteractively(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	cleanUp := func() {
		os.RemoveAll(workingDir)
	}
	cleanUp()
	defer cleanUp()

	err := os.Mkdir(workingDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Unable to create working directory: %s", err.Error())
	}

	logger.Section("Creating a package repository interactively using pkg repo release", func() {
		promptOutput := newPromptOutput(t)

		go func() {
			promptOutput.WaitFor("Enter the package repository name")
			promptOutput.Write("testpackagerepo.corp.dev")
			promptOutput.WaitFor("Enter the registry url")
			promptOutput.Write(env.Image)
		}()

		err = os.Mkdir(filepath.Join(workingDir, packagesDir), fs.ModePerm)
		if err != nil {
			t.Errorf("Unable to create packages directory: %s", err.Error())
		}

		kappCtrl.RunWithOpts([]string{"pkg", "repo", "release", "--tty=true", "--chdir", workingDir},
			RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
				StdoutWriter: promptOutput.BufferedOutputWriter(), Interactive: true})

		keysToBeIgnored := []string{"creationTimestamp:", "image"}
		verifyPackageRepoBuild(t, keysToBeIgnored)
		verifyPackageRepository(t, keysToBeIgnored)
	})
}

func verifyPackageRepoBuild(t *testing.T, keysToBeIgnored []string) {
	packageRepoBuildExpectedOutput := `
apiVersion: kctrl.carvel.dev/v1alpha1
kind: PackageRepositoryBuild
metadata:
  name: testpackagerepo.corp.dev
spec:
  export:
    imgpkgBundle:
`
	out, err := readFile("pkgrepo-build.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	out = strings.TrimSpace(replaceSpaces(out))
	packageRepoBuildExpectedOutput = strings.TrimSpace(replaceSpaces(packageRepoBuildExpectedOutput))
	out = strings.TrimSpace(clearKeys(keysToBeIgnored, out))
	require.Equal(t, packageRepoBuildExpectedOutput, out, "output does not match")

}

func verifyPackageRepository(t *testing.T, keysToBeIgnored []string) {
	packageRepoBuildExpectedOutput := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: testpackagerepo.corp.dev
spec:
  fetch:
    imgpkgBundle:
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0
`
	out, err := readFile("package-repository.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	out = strings.TrimSpace(replaceSpaces(out))
	packageRepoBuildExpectedOutput = strings.TrimSpace(replaceSpaces(packageRepoBuildExpectedOutput))
	out = clearKeys(keysToBeIgnored, out)
	require.Equal(t, packageRepoBuildExpectedOutput, out, "output does not match")
}
