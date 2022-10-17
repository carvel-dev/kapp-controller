package e2e

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

const (
	packagesDir       = "packages"
	pkgRepoOutputFile = "package-repository.yml"
	pkgrName          = "testpackagerepo.corp.dev"
)

func TestPackageRepositoryReleaseInteractively(t *testing.T) {
	cleanUp := func() {
		os.RemoveAll(workingDir)
	}
	cleanUp()
	defer cleanUp()

	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}
	kappCli := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	promptOutput := newPromptOutput(t)
	setup(t)

	logger.Section("Creating a package repository interactively using pkg repo release", func() {

		go func() {
			promptOutput.WaitFor("Enter the package repository name")
			promptOutput.Write(pkgrName)
			promptOutput.WaitFor("Enter the registry url")
			promptOutput.Write(env.Image)
		}()

		kappCtrl.RunWithOpts([]string{"pkg", "repo", "release", "--tty=true", "--chdir", workingDir, "--version", "1.0.0"},
			RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
				StdoutWriter: promptOutput.BufferedOutputWriter(), Interactive: true})

		keysToBeIgnored := []string{"creationTimestamp:", "image"}
		verifyPackageRepoBuild(t, keysToBeIgnored)
		verifyPackageRepository(t, keysToBeIgnored)
	})

	logger.Section(fmt.Sprintf("Installing package repository"), func() {
		pkgrKappAppName := "test-package-repo-app"
		cleanUpInstalledPkgRepo := func() {
			kappCli.RunWithOpts([]string{"delete", "-a", pkgrKappAppName},
				RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
		}
		defer cleanUpInstalledPkgRepo()

		kappCli.RunWithOpts([]string{"deploy", "-a", pkgrKappAppName, "-f", filepath.Join(workingDir, pkgRepoOutputFile), "-c"},
			RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
		out, _ := kappCtrl.RunWithOpts([]string{"package", "repository", "get", "-r", pkgrName, "--json"}, RunOpts{})

		output := uitest.JSONUIFromBytes(t, []byte(out))
		for _, m := range output.Tables[0].Rows {
			val, ok := m["source"]
			if ok {
				require.Contains(t, val, fmt.Sprintf("(imgpkg) (1.0.0) %s", env.Image))
			}
		}
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
  annotations:
    tag: 1.0.0
  name: testpackagerepo.corp.dev
spec:
  fetch:
    imgpkgBundle:
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0
`
	out, err := readFile(pkgRepoOutputFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	out = strings.TrimSpace(replaceSpaces(out))
	packageRepoBuildExpectedOutput = strings.TrimSpace(replaceSpaces(packageRepoBuildExpectedOutput))
	out = clearKeys(keysToBeIgnored, out)
	require.Equal(t, packageRepoBuildExpectedOutput, out, "output does not match")
}

func setup(t *testing.T) {
	pkg_metadata_yaml := `---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: test-pkg.carvel.dev
spec:
  displayName: "Carvel Test Package"
  shortDescription: "Carvel package for testing installation"
`
	pkg_yaml := `---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: test-pkg.carvel.dev.1.0.0
spec:
  refName: test-pkg.carvel.dev
  version: 1.0.0
  valuesSchema:
    openAPIv3:
      properties:
        app_port:
          default: 80
          description: App port
          type: integer
        app_name:
          description: App Name
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
      - kapp: {}
`
	err := os.MkdirAll(filepath.Join(workingDir, packagesDir), fs.ModePerm)
	if err != nil {
		t.Errorf("Unable to create packages directory: %s", err.Error())
	}
	err = os.WriteFile(filepath.Join(workingDir, packagesDir, "package.yml"), []byte(pkg_yaml), fs.ModePerm)
	if err != nil {
		t.Errorf("Unable to create package.yml file: %s", err.Error())
	}
	err = os.WriteFile(filepath.Join(workingDir, packagesDir, "metadata.yml"), []byte(pkg_metadata_yaml), fs.ModePerm)
	if err != nil {
		t.Errorf("Unable to create metadata.yml file: %s", err.Error())
	}
}
