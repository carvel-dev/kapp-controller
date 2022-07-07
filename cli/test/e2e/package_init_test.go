package e2e

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestPackageInitInteractively(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	cleanUp := func() {
		os.Remove("package-build.yml")
		os.Remove("package-resources.yml")
		os.Remove("vendir.yml")
		os.Remove("app-build.yml")
		os.RemoveAll("upstream")
	}
	cleanUp()
	defer cleanUp()

	logger.Section("Creating a package interactively using pkg init", func() {
		promptOutput := newPromptOutput(t)

		go func() {
			promptOutput.WaitFor("A package Reference name must be a valid DNS subdomain name")
			promptOutput.Write("testpackage.corp.dev")
			promptOutput.WaitFor("need to fetch the manifest which defines")
			promptOutput.Write("3")
			promptOutput.WaitFor("From where to fetch the")
			promptOutput.Write("enterprise-operator")
			promptOutput.WaitFor("Enter helm chart name")
			promptOutput.Write("1.16.0")
			promptOutput.WaitFor("Enter helm chart version")
			promptOutput.Write("https://mongodb.github.io/helm-charts")
		}()

		kappCtrl.RunWithOpts([]string{"pkg", "init", "--tty=true"},
			RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
				StdoutWriter: promptOutput.OutputWriter(), Interactive: true})

		keysToBeIgnored := []string{"creationTimestamp:", "releasedAt:"}
		verifyPackageBuild(t, keysToBeIgnored)
		verifyPackageResources(t, keysToBeIgnored)
		verifyVendir(t, keysToBeIgnored)

	})
}

func verifyPackageBuild(t *testing.T, keysToBeIgnored []string) {
	packageBuildExpectedOutput := `
apiVersion: kctrl.carvel.dev/v1alpha1
kind: PackageBuild
metadata:
  name: testpackage.corp.dev
spec:
  template:
    apiVersion: kctrl.carvel.dev/v1alpha1
    kind: AppBuild
    metadata:
      annotations:
        fetch-content-from: Helm Chart from Helm Repository
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
	out, err := readFile("package-build.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	out = strings.TrimSpace(replaceSpaces(out))
	packageBuildExpectedOutput = strings.TrimSpace(replaceSpaces(packageBuildExpectedOutput))
	out = clearKeys(keysToBeIgnored, out)
	require.Equal(t, packageBuildExpectedOutput, out, "output does not match")

}

func verifyPackageResources(t *testing.T, keysToBeIgnored []string) {
	packagesResourcesExpectedOutput := `
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
  serviceAccountName: testpackage-sa
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0
`
	out, err := readFile("package-resources.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	out = strings.TrimSpace(replaceSpaces(out))
	packagesResourcesExpectedOutput = strings.TrimSpace(replaceSpaces(packagesResourcesExpectedOutput))
	out = clearKeys(keysToBeIgnored, out)
	require.Equal(t, packagesResourcesExpectedOutput, out, "output does not match")
}

func verifyVendir(t *testing.T, keysToBeIgnored []string) {
	vendirExpectedOutput := `
apiVersion: vendir.k14s.io/v1alpha1
directories:
- contents:
  - helmChart:
      helmVersion: "3"
      name: enterprise-operator
      repository:
        url: https://mongodb.github.io/helm-charts
      version: 1.16.0
    path: .
  path: upstream
kind: Config
minimumRequiredVersion: ""
`
	out, err := readFile("vendir.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	out = strings.TrimSpace(replaceSpaces(out))
	vendirExpectedOutput = strings.TrimSpace(replaceSpaces(vendirExpectedOutput))
	out = clearKeys(keysToBeIgnored, out)
	require.Equal(t, vendirExpectedOutput, out, "output does not match")
}

func readFile(fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
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
	outputWriter io.Writer
	outputReader io.Reader
}

func newPromptOutput(t *testing.T) promptOutput {
	stringReader, stringWriter, err := os.Pipe()
	require.NoError(t, err)

	outputReader, outputWriter, err := os.Pipe()
	require.NoError(t, err)

	return promptOutput{t, stringWriter, stringReader, outputWriter, outputReader}
}

func (p promptOutput) WritePkgRefName() {
	p.stringWriter.Write([]byte("afc.def.ghi\n"))
}

func (p promptOutput) OutputWriter() io.Writer { return p.outputWriter }
func (p promptOutput) OutputReader() io.Reader { return p.outputReader }
func (p promptOutput) StringWriter() io.Writer { return p.stringWriter }
func (p promptOutput) StringReader() io.Reader { return p.stringReader }

func (p promptOutput) Write(val string) {
	p.stringWriter.Write([]byte(val + "\n"))
}

func (p promptOutput) WaitFor(text string) {
	scanner := bufio.NewScanner(p.outputReader)
	for scanner.Scan() {
		// Cannot easily wait for prompt as it's not NL terminated
		if strings.Contains(scanner.Text(), text) {
			break
		}
	}
	err := scanner.Err()
	require.NoError(p.t, err)
}
