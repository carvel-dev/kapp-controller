// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

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
	workingDir = "kctrl-test"
)

type E2EAuthoringTestCase struct {
	Name                    string
	InitInteraction         Interaction
	ExpectedPkgBuild        string
	ExpectedPkgResource     string
	ExpectedVendir          string
	ExpectedPackage         string
	ExpectedPackageMetadata string
}

type Interaction struct {
	Prompts []string
	Inputs  []string
}

func (i Interaction) Run(promptOutputObj promptOutput) {
	for ind, prompt := range i.Prompts {
		promptOutputObj.WaitFor(prompt)
		promptOutputObj.Write(i.Inputs[ind])
	}
}

func TestE2EInitAndReleaseCases(t *testing.T) {
	testcases := []E2EAuthoringTestCase{
		{
			Name: "Helm Chart Flow",
			InitInteraction: Interaction{
				Prompts: []string{
					"Enter the package reference name",
					"Enter source",
					"Enter helm chart repository URL",
					"Enter helm chart name",
					"Enter helm chart version",
				},
				Inputs: []string{
					"testpackage.corp.dev",
					"3",
					"https://mongodb.github.io/helm-charts",
					"enterprise-operator",
					"1.16.0",
				},
			},
			ExpectedPkgBuild: `
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
`,
			ExpectedPkgResource: `
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
`,
			ExpectedVendir: `
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
`,
			ExpectedPackageMetadata: `
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: testpackage.corp.dev
spec:
  displayName: testpackage
  longDescription: testpackage.corp.dev
  shortDescription: testpackage.corp.dev
`,
			ExpectedPackage: `
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
`,
		},
		{
			Name: "Github Release Flow",
			InitInteraction: Interaction{
				Prompts: []string{
					"Enter the package reference name",
					"Enter source",
					"Enter slug for repository",
					"Enter the release tag to be used",
					"Enter the paths which contain Kubernetes manifests",
				},
				Inputs: []string{
					"testpackage.corp.dev",
					"2",
					"Dynatrace/dynatrace-operator",
					"v0.6.0",
					"kubernetes.yaml",
				},
			},
			ExpectedPkgBuild: `
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
          - ytt:
              paths:
              - upstream
          - kbld: {}
      export:
      - includePaths:
        - upstream
`,
			ExpectedPkgResource: `
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
      - ytt:
          paths:
          - upstream
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
`,
			ExpectedVendir: `
apiVersion: vendir.k14s.io/v1alpha1
directories:
- contents:
  - githubRelease:
      disableAutoChecksumValidation: true
      slug: Dynatrace/dynatrace-operator
      tag: v0.6.0
    includePaths:
    - kubernetes.yaml
    path: .
  path: upstream
kind: Config
minimumRequiredVersion: ""
`,
			ExpectedPackageMetadata: `
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: testpackage.corp.dev
spec:
  displayName: testpackage
  longDescription: testpackage.corp.dev
  shortDescription: testpackage.corp.dev
`,
			ExpectedPackage: `
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
      - ytt:
          paths:
          - upstream
      - kbld:
          paths:
          - '-'
          - .imgpkg/
  valuesSchema:
    openAPIv3:
      default: null
      nullable: true
  version: 1.0.0
`,
		},
		{
			Name: "Git Repository Flow",
			InitInteraction: Interaction{
				Prompts: []string{
					"Enter the package reference name",
					"Enter source",
					"Enter Git URL",
					"Enter Git Reference",
					"Enter the paths which contain Kubernetes manifests",
				},
				Inputs: []string{
					"testpackage.corp.dev",
					"4",
					"https://github.com/vmware-tanzu/carvel-kapp",
					"origin/develop",
					"examples/simple-app-example/config-1.yml",
				},
			},
			ExpectedPkgBuild: `
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
          - ytt:
              paths:
              - upstream
          - kbld: {}
      export:
      - includePaths:
        - upstream
`,
			ExpectedPkgResource: `
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
      - ytt:
          paths:
          - upstream
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
`,
			ExpectedVendir: `
apiVersion: vendir.k14s.io/v1alpha1
directories:
- contents:
  - git:
      ref: origin/develop
      url: https://github.com/vmware-tanzu/carvel-kapp
    includePaths:
    - examples/simple-app-example/config-1.yml
    path: .
  path: upstream
kind: Config
minimumRequiredVersion: ""
`,
			ExpectedPackageMetadata: `
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: testpackage.corp.dev
spec:
  displayName: testpackage
  longDescription: testpackage.corp.dev
  shortDescription: testpackage.corp.dev
`,
			ExpectedPackage: `
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
      - ytt:
          paths:
          - upstream
      - kbld:
          paths:
          - '-'
          - .imgpkg/
  valuesSchema:
    openAPIv3:
      default: null
      nullable: true
  version: 1.0.0
`,
		},
		{
			Name: "Helm Chart from Git Flow",
			InitInteraction: Interaction{
				Prompts: []string{
					"Enter the package reference name",
					"Enter source",
					"Enter Git URL",
					"Enter Git Reference",
					"Enter the paths which contain Kubernetes manifests",
				},
				Inputs: []string{
					"testpackage.corp.dev",
					"5",
					"https://github.com/rohitagg2020/helm-simple-app",
					"main",
					"simple-app/**/*",
				},
			},
			ExpectedPkgBuild: `
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
              path: upstream/simple-app
          - ytt:
              paths:
              - '-'
          - kbld: {}
      export:
      - includePaths:
        - upstream
`,
			ExpectedPkgResource: `
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
          path: upstream/simple-app
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
`,
			ExpectedVendir: `
apiVersion: vendir.k14s.io/v1alpha1
directories:
- contents:
  - git:
      ref: main
      url: https://github.com/rohitagg2020/helm-simple-app
    includePaths:
    - simple-app/**/*
    path: .
  path: upstream
kind: Config
minimumRequiredVersion: ""
`,
			ExpectedPackageMetadata: `
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: testpackage.corp.dev
spec:
  displayName: testpackage
  longDescription: testpackage.corp.dev
  shortDescription: testpackage.corp.dev
`,
			ExpectedPackage: `
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
          path: upstream/simple-app
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
`,
		},
	}

	env := BuildEnv(t)
	logger := Logger{}
	kappCli := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	for _, testcase := range testcases {
		// verify prompts and input are of same length
		require.EqualValues(t, len(testcase.InitInteraction.Prompts), len(testcase.InitInteraction.Inputs))

		cleanUp := func() {
			os.RemoveAll(workingDir)
		}
		cleanUp()

		err := os.Mkdir(workingDir, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		const pkgDir = "./carvel-artifacts/packages/testpackage.corp.dev/"
		promptOutput := newPromptOutput(t)

		go testcase.InitInteraction.Run(promptOutput)

		logger.Section(fmt.Sprintf("%s: Package init", testcase.Name), func() {
			kappCtrl.RunWithOpts([]string{"pkg", "init", "--tty=true", "--chdir", workingDir},
				RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
					StdoutWriter: promptOutput.BufferedOutputWriter(), Interactive: true})

			// Below key's values will be changed during every run, hence adding these keys to be ignored
			keysToBeIgnored := []string{"creationTimestamp:", "releasedAt:"}

			// Error if upstream folder doesn't exist
			_, err = os.Stat(filepath.Join(workingDir, "upstream"))
			require.NoError(t, err)

			// Verify PackageBuild
			out, err := readFile("package-build.yml")
			require.NoErrorf(t, err, "Expected to read package-build.yml")
			expectedPackageBuild := strings.TrimSpace(replaceSpaces(testcase.ExpectedPkgBuild))
			out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
			require.Equal(t, expectedPackageBuild, out, "Expected PackageBuild output to match")

			// Verify package resources
			out, err = readFile("package-resources.yml")
			require.NoErrorf(t, err, "Expected to read package-resources.yml")
			expectedPackageResources := strings.TrimSpace(replaceSpaces(testcase.ExpectedPkgResource))
			out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
			require.Equal(t, expectedPackageResources, out, "Expected package resources output to match")

			// Verify vendir
			out, err = readFile("vendir.yml")
			require.NoErrorf(t, err, "Expected to read vendir.yml")
			expectedVendirOutput := strings.TrimSpace(replaceSpaces(testcase.ExpectedVendir))
			out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
			require.Equal(t, expectedVendirOutput, out, "Expected vendir output to match")
		})

		logger.Section(fmt.Sprintf("%s: Package release", testcase.Name), func() {
			releaseInteraction := Interaction{
				Prompts: []string{"Enter the registry URL"},
				Inputs:  []string{env.Image},
			}

			go releaseInteraction.Run(promptOutput)

			kappCtrl.RunWithOpts([]string{"pkg", "release", "--version", "1.0.0", "--tty=true", "--chdir", workingDir},
				RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
					StdoutWriter: promptOutput.BufferedOutputWriter(), Interactive: true})

			// Below key's values will be changed during every run, hence adding these keys to be ignored
			keysToBeIgnored := []string{"creationTimestamp:", "releasedAt:", "image"}

			// Verify PackageMetadata artifact
			out, err := readFile(pkgDir + "metadata.yml")
			require.NoErrorf(t, err, "Expected to read metadata.yml")
			expectedPackageMetadata := strings.TrimSpace(replaceSpaces(testcase.ExpectedPackageMetadata))
			out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
			require.Equal(t, expectedPackageMetadata, out, "Expected PackageMetadata to match")

			// Verify Package artifact
			out, err = readFile(pkgDir + "package.yml")
			require.NoErrorf(t, err, "Expected to read package.yml")
			expectedPackage := strings.TrimSpace(replaceSpaces(testcase.ExpectedPackage))
			out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
			require.Equal(t, expectedPackage, out, "Expected Package to match")
		})

		logger.Section(fmt.Sprintf("%s: Testing and installing created Package", testcase.Name), func() {
			cleanUpInstalledPkg := func() {
				switch testcase.Name {
				case "Github Release Flow":
					kubectl.RunWithOpts([]string{"delete", "ns", "dynatrace"}, RunOpts{NoNamespace: true})
				}
				kappCli.RunWithOpts([]string{"delete", "-a", "test-package"},
					RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
				kappCtrl.RunWithOpts([]string{"pkg", "installed", "delete", "-i", "test"},
					RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
			}
			defer cleanUpInstalledPkg()

			switch testcase.Name {
			case "Github Release Flow":
				kubectl.RunWithOpts([]string{"create", "ns", "dynatrace"}, RunOpts{NoNamespace: true})
			}
			kappCli.RunWithOpts([]string{"deploy", "-a", "test-package", "-f", filepath.Join(workingDir, pkgDir), "-c"},
				RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
			kappCtrl.RunWithOpts([]string{"pkg", "available", "list"},
				RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
			kappCtrl.RunWithOpts([]string{"pkg", "install", "-p", "testpackage.corp.dev", "-i", "test", "--version", "1.0.0"},
				RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
		})
	}
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
	t              *testing.T
	stringWriter   io.Writer
	stringReader   io.Reader
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
	for ; attempts < 60; attempts++ {
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

		// Poll interval
		time.Sleep(1 * time.Second)
	}
	if attempts == 60 {
		fmt.Printf("Timed out waiting for text '%s'", text)
		p.t.Fail()
	}
}
