package e2e

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	workingDir = "kcrl-test"
)

func TestPackageInitAndRelease(t *testing.T) {
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
		t.Fatal(err)
	}

	expectedPackageBuild := `
apiVersion: kctrl.carvel.dev/v1alpha1
kind: PackageBuild
metadata:
  name: testpackage.corp.dev
spec:
  template:
    spec:
      app:
        spec:
          deploy:
          - kapp: {}
          template:
          - helmTemplate:
              path: upstream
          - ytt:
              paths:
              - '-'
          - kbld: {}
      export:
      - includePaths:
        - upstream
`

	expectedPackageResources := `
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: testpackage.corp.dev.0.0.0
spec:
  refName: testpackage.corp.dev
  template:
    spec:
      deploy:
      - kapp: {}
      fetch:
      - git: {}
      template:
      - helmTemplate:
          path: upstream
      - ytt:
          paths:
          - '-'
      - kbld: {}
  valuesSchema:
    openAPIv3: null
  version: 0.0.0
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: testpackage.corp.dev
spec:
  displayName: testpackage
  longDescription: testpackage.corp.dev
  shortDescription: testpackage.corp.dev
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  annotations:
    kctrl.carvel.dev/local-fetch-0: .
  name: testpackage
spec:
  packageRef:
    refName: testpackage.corp.dev
    versionSelection:
      constraints: 0.0.0
  serviceAccountName: testpackage-sa
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0
`

	expectedVendirOutput := `
apiVersion: vendir.k14s.io/v1alpha1
directories:
- contents:
  - helmChart:
      name: enterprise-operator
      repository:
        url: https://mongodb.github.io/helm-charts
      version: 1.16.0
    path: .
  path: upstream
kind: Config
minimumRequiredVersion: ""
`

	expectedPackageMetadata := `
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: testpackage.corp.dev
spec:
  displayName: testpackage
  longDescription: testpackage.corp.dev
  shortDescription: testpackage.corp.dev
`

	expectedPackage := `
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: testpackage.corp.dev.1.0.0
spec:
  refName: testpackage.corp.dev
  template:
    spec:
      deploy:
      - kapp: {}
      fetch:
      - imgpkgBundle:
      template:
      - helmTemplate:
          path: upstream
      - ytt:
          paths:
          - '-'
      - kbld:
          paths:
          - '-'
          - .imgpkg/
  valuesSchema:
    openAPIv3:
      default: null
      nullable: true
  version: 1.0.0
`

	logger.Section("creating a package interactively using pkg init", func() {
		promptOutput := newPromptOutput(t)

		go func() {
			promptOutput.WaitFor("Enter the package reference name")
			promptOutput.Write("testpackage.corp.dev")
			promptOutput.WaitFor("Enter source")
			promptOutput.Write("3")
			promptOutput.WaitFor("Enter helm chart repository URL")
			promptOutput.Write("https://mongodb.github.io/helm-charts")
			promptOutput.WaitFor("Enter helm chart name")
			promptOutput.Write("enterprise-operator")
			promptOutput.WaitFor("Enter helm chart version")
			promptOutput.Write("1.16.0")
		}()

		kappCtrl.RunWithOpts([]string{"pkg", "init", "--tty=true", "--chdir", workingDir},
			RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
				StdoutWriter: promptOutput.BufferedOutputWriter(), Interactive: true})

		keysToBeIgnored := []string{"creationTimestamp:", "releasedAt:"}

		// To handle case till errors on failed syncs are handled better
		_, err := os.Stat(filepath.Join(workingDir, "upstream"))
		require.NoError(t, err)

		// Verify PackageBuild
		out, err := readFile("package-build.yml")
		require.NoErrorf(t, err, "Expected to read package-build.yml")

		expectedPackageBuild = strings.TrimSpace(replaceSpaces(expectedPackageBuild))
		out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
		require.Equal(t, expectedPackageBuild, out, "Expected PackageBuild output to match")

		// Verify package resources
		out, err = readFile("package-resources.yml")
		require.NoErrorf(t, err, "Expected to read package-resources.yml")

		expectedPackageResources = strings.TrimSpace(replaceSpaces(expectedPackageResources))
		out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
		require.Equal(t, expectedPackageResources, out, "Expected package resources output to match")

		// Verify vendir
		out, err = readFile("vendir.yml")
		require.NoErrorf(t, err, "Expected to read vendir.yml")

		expectedVendirOutput = strings.TrimSpace(replaceSpaces(expectedVendirOutput))
		out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
		require.Equal(t, expectedVendirOutput, out, "Expected vendir output to match")
	})

	logger.Section("releasing package using kctrl package release", func() {
		promptOutput := newPromptOutput(t)

		go func() {
			promptOutput.WaitFor("Enter the registry URL")
			promptOutput.Write(env.Image)
		}()

		kappCtrl.RunWithOpts([]string{"pkg", "release", "--version", "1.0.0", "--tty=true", "--chdir", workingDir},
			RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
				StdoutWriter: promptOutput.BufferedOutputWriter(), Interactive: true})

		keysToBeIgnored := []string{"creationTimestamp:", "releasedAt:", "image"}

		// Verify PackageMetadata artifact
		out, err := readFile("./carvel-artifacts/packages/testpackage.corp.dev/metadata.yml")
		require.NoErrorf(t, err, "Expected to read metadata.yml")

		expectedPackageMetadata = strings.TrimSpace(replaceSpaces(expectedPackageMetadata))
		out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
		require.Equal(t, expectedPackageMetadata, out, "Expected PackageMetadata to match")

		// Verify Package artifact
		out, err = readFile("./carvel-artifacts/packages/testpackage.corp.dev/package.yml")
		require.NoErrorf(t, err, "Expected to read package.yml")

		expectedPackage = strings.TrimSpace(replaceSpaces(expectedPackage))
		out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
		require.Equal(t, expectedPackage, out, "Expected Package to match")
	})
}

func readFile(fileName string) (string, error) {
	path := filepath.Join(workingDir, fileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil

}

func replaceSpaces(result string) string {
	// result = strings.Replace(result, " ", "_", -1) // useful for debugging
	result = strings.Replace(result, " \n", " $\n", -1) // explicit endline
	return result
}

// TODO: Make regex more strict. Removes 'images.yaml' from '.imgpkg/images.yml' right now
func clearKeys(keys []string, out string) string {
	for _, key := range keys {
		r := regexp.MustCompile(key + ".*")
		out = r.ReplaceAllString(out, "")

	}
	//removing all empty lines
	r := regexp.MustCompile(`[ ]*[\n\t]*\n`)
	out = r.ReplaceAllString(out, "\n")
	out = strings.ReplaceAll(out, "\n\n", "\n")
	return out
}

type promptOutput struct {
	t            *testing.T
	stringWriter io.Writer
	stringReader io.Reader

	bufferedStdout *bytes.Buffer
}

func newPromptOutput(t *testing.T) promptOutput {
	stringReader, stringWriter, err := os.Pipe()
	require.NoError(t, err)

	bufferedStdout := bytes.Buffer{}

	return promptOutput{t, stringWriter, stringReader, &bufferedStdout}
}

func (p promptOutput) WritePkgRefName() {
	p.stringWriter.Write([]byte("afc.def.ghi\n"))
}

func (p promptOutput) StringWriter() io.Writer         { return p.stringWriter }
func (p promptOutput) StringReader() io.Reader         { return p.stringReader }
func (p promptOutput) BufferedOutputWriter() io.Writer { return p.bufferedStdout }

func (p promptOutput) Write(val string) {
	p.stringWriter.Write([]byte(val + "\n"))
}

func (p promptOutput) WaitFor(text string) {
	attempts := 0
	// Poll buffered output till desired string is found
	for {
		found := strings.Contains(p.bufferedStdout.String(), text)

		if os.Getenv("KCTRL_DEBUG_BUFERED_OUTPUT_TESTS") == "true" {
			fmt.Printf("\n==> \t Buffered Output (Waiting for '%s' | found=%t)\n", text, found)
			fmt.Println(p.bufferedStdout)
			fmt.Printf("--------------------------------------\n\n")
		}

		// Stop waiting if desired string is found
		if found {
			break
		}
		attempts++

		if attempts > 60 {
			fmt.Printf("Timed out waiting for text '%s'", text)
			p.t.Fail()
		}
		// Poll interval
		time.Sleep(1 * time.Second)
	}
}
