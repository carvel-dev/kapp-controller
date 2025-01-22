// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
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
          - .imgpkg/images.yml
  valuesSchema:
    openAPIv3:
      properties:
        agent:
          properties:
            name:
              default: mongodb-agent
              type: string
            version:
              default: 11.12.0.7388-1
              type: string
          type: object
        database:
          description: Database
          properties:
            name:
              default: mongodb-enterprise-database
              type: string
            version:
              default: 2.0.2
              type: string
          type: object
        initAppDb:
          description: Application Database
          properties:
            name:
              default: mongodb-enterprise-init-appdb
              type: string
            version:
              default: 1.0.9
              type: string
          type: object
        initDatabase:
          properties:
            name:
              default: mongodb-enterprise-init-database
              type: string
            version:
              default: 1.0.9
              type: string
          type: object
        initOpsManager:
          properties:
            name:
              default: mongodb-enterprise-init-ops-manager
              type: string
            version:
              default: 1.0.7
              type: string
          type: object
        managedSecurityContext:
          default: false
          description: Set this to true if your cluster is managing SecurityContext
            for you. If running OpenShift (Cloud, Minishift, etc.), set this to true.
          type: boolean
        mongodb:
          properties:
            name:
              default: mongodb-enterprise-appdb-database
              type: string
            repo:
              default: quay.io/mongodb
              type: string
          type: object
        multiCluster:
          properties:
            clusters:
              default: []
              items: {}
              type: array
            kubeConfigSecretName:
              default: mongodb-enterprise-operator-multi-cluster-kubeconfig
              type: string
          type: object
        operator:
          properties:
            affinity:
              default: {}
              type: object
            createOperatorServiceAccount:
              default: true
              description: Create operator-service account
              type: boolean
            deployment_name:
              default: mongodb-enterprise-operator
              description: Name of the deployment of the operator pod
              type: string
            env:
              default: prod
              description: Execution environment for the operator, dev or prod. Use
                dev for more verbose logging
              type: string
            name:
              default: mongodb-enterprise-operator
              description: Name that will be assigned to most of internal Kubernetes
                objects like Deployment, ServiceAccount, Role etc.
              type: string
            nodeSelector:
              default: {}
              type: object
            operator_image_name:
              default: mongodb-enterprise-operator
              description: Name of the operator image
              type: string
            tolerations:
              default: []
              items: {}
              type: array
            vaultSecretBackend:
              properties:
                enabled:
                  default: false
                  description: set to true if you want the operator to store secrets
                    in Vault
                  type: boolean
                tlsSecretRef:
                  default: ""
                  type: string
              type: object
            version:
              default: 1.16.0
              description: Version of mongodb-enterprise-operator
              type: string
            watchedResources:
              default: []
              description: The Custom Resources that will be watched by the Operator.
                Needs to be changed if only some of the CRDs are installed
              items:
                default: mongodb
                type: string
              type: array
          type: object
        opsManager:
          description: Ops Manager
          properties:
            name:
              default: mongodb-enterprise-ops-manager
              type: string
          type: object
        registry:
          description: Registry
          properties:
            agent:
              default: quay.io/mongodb
              type: string
            appDb:
              default: quay.io/mongodb
              type: string
            database:
              default: quay.io/mongodb
              type: string
            imagePullSecrets:
              oneOf:
              - default: null
                nullable: true
                type: integer
              - default: null
                nullable: true
                type: number
              - default: null
                nullable: true
                type: boolean
              - default: null
                nullable: true
                type: string
              - default: null
                nullable: true
                type: object
              - default: null
                items: {}
                nullable: true
                type: array
            initAppDb:
              default: quay.io/mongodb
              type: string
            initDatabase:
              default: quay.io/mongodb
              type: string
            initOpsManager:
              default: quay.io/mongodb
              type: string
            operator:
              default: quay.io/mongodb
              description: Specify if images are pulled from private registry
              type: string
            opsManager:
              default: quay.io/mongodb
              type: string
            pullPolicy:
              default: Always
              description: 'TODO: specify for each image and move there?'
              type: string
          type: object
        subresourceEnabled:
          default: true
          description: Set this to false to disable subresource utilization It might
            be required on some versions of Openshift
          type: boolean
      type: object
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
          - .imgpkg/images.yml
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
          - .imgpkg/images.yml
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
          - .imgpkg/images.yml
  valuesSchema:
    openAPIv3:
      properties:
        env:
          properties:
            hello_msg:
              default: hello
              description: response provided by the server
              type: string
          type: object
        fullnameOverride:
          default: ""
          type: string
        imageProp:
          properties:
            pullPolicy:
              default: IfNotPresent
              type: string
            repository:
              default: docker.io/dkalinin/k8s-simple-app
              type: string
            tag:
              default: latest
              description: Overrides the image tag whose default is the chart appVersion.
              type: string
          type: object
        imagePullSecrets:
          default: []
          items: {}
          type: array
        nameOverride:
          default: ""
          type: string
        podAnnotations:
          default: {}
          type: object
        replicaCount:
          default: 1
          type: integer
        service:
          properties:
            port:
              default: 80
              type: integer
            type:
              default: ClusterIP
              type: string
          type: object
        serviceAccount:
          properties:
            annotations:
              default: {}
              description: Annotations to add to the service account
              type: object
            create:
              default: true
              description: Specifies whether a service account should be created
              type: boolean
            name:
              default: ""
              description: The name of the service account to use. If not set and
                create is true, a name is generated using the fullname template
              type: string
          type: object
      type: object
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
		defer cleanUp()

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
			verifyBoilerplate(t, testcase)
		})

		// Ensure that defaults are retained by running in non-interactive mode
		logger.Section(fmt.Sprintf("%s: Package init non-interactive", testcase.Name), func() {
			kappCtrl.RunWithOpts([]string{"pkg", "init", "--tty=true", "--chdir", workingDir}, RunOpts{NoNamespace: true})
			verifyBoilerplate(t, testcase)
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
			keysToBeIgnored := []string{"creationTimestamp:", "releasedAt:", "image:"}

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

func verifyBoilerplate(t *testing.T, testcase E2EAuthoringTestCase) {
	// Below key's values will be changed during every run, hence adding these keys to be ignored
	keysToBeIgnored := []string{"creationTimestamp:", "releasedAt:"}

	// Error if upstream folder doesn't exist
	_, err := os.Stat(filepath.Join(workingDir, "upstream"))
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

func TestE2EInitAndReleaseCaseDisableOpenAPISchemaGeneration(t *testing.T) {
	input := E2EAuthoringTestCase{
		Name: "Disable OpenAPI Schema generation while releasing the package",
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
          - .imgpkg/images.yml
  valuesSchema:
    openAPIv3: null
  version: 1.0.0
`,
	}
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}
	kappCli := Kapp{t, env.Namespace, env.KappBinaryPath, logger}

	cleanUp := func() {
		os.RemoveAll(workingDir)
	}
	cleanUp()
	defer cleanUp()

	err := os.Mkdir(workingDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	const pkgDir = "./carvel-artifacts/packages/testpackage.corp.dev/"
	promptOutput := newPromptOutput(t)
	go input.InitInteraction.Run(promptOutput)

	logger.Section(fmt.Sprintf("%s: Package init", input.Name), func() {
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
		expectedPackageBuild := strings.TrimSpace(replaceSpaces(input.ExpectedPkgBuild))
		out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
		require.Equal(t, expectedPackageBuild, out, "Expected PackageBuild output to match")

		// Verify package resources
		out, err = readFile("package-resources.yml")
		require.NoErrorf(t, err, "Expected to read package-resources.yml")
		expectedPackageResources := strings.TrimSpace(replaceSpaces(input.ExpectedPkgResource))
		out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
		require.Equal(t, expectedPackageResources, out, "Expected package resources output to match")

		// Verify vendir
		out, err = readFile("vendir.yml")
		require.NoErrorf(t, err, "Expected to read vendir.yml")
		expectedVendirOutput := strings.TrimSpace(replaceSpaces(input.ExpectedVendir))
		out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
		require.Equal(t, expectedVendirOutput, out, "Expected vendir output to match")
	})

	logger.Section(fmt.Sprintf("%s: Package release", input.Name), func() {
		releaseInteraction := Interaction{
			Prompts: []string{"Enter the registry URL"},
			Inputs:  []string{env.Image},
		}

		go releaseInteraction.Run(promptOutput)

		kappCtrl.RunWithOpts([]string{"pkg", "release", "--version", "1.0.0", "--tty=true", "--chdir", workingDir, "--openapi-schema=false"},
			RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
				StdoutWriter: promptOutput.BufferedOutputWriter(), Interactive: true})

		// Below key's values will be changed during every run, hence adding these keys to be ignored
		keysToBeIgnored := []string{"creationTimestamp:", "releasedAt:", "image:"}

		// Verify PackageMetadata artifact
		out, err := readFile(pkgDir + "metadata.yml")
		require.NoErrorf(t, err, "Expected to read metadata.yml")
		expectedPackageMetadata := strings.TrimSpace(replaceSpaces(input.ExpectedPackageMetadata))
		out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
		require.Equal(t, expectedPackageMetadata, out, "Expected PackageMetadata to match")

		// Verify Package artifact
		out, err = readFile(pkgDir + "package.yml")
		require.NoErrorf(t, err, "Expected to read package.yml")
		expectedPackage := strings.TrimSpace(replaceSpaces(input.ExpectedPackage))
		out = clearKeys(keysToBeIgnored, strings.TrimSpace(replaceSpaces(out)))
		require.Equal(t, expectedPackage, out, "Expected Package to match")
	})

	logger.Section(fmt.Sprintf("%s: Testing and installing created Package", input.Name), func() {
		cleanUpInstalledPkg := func() {
			kappCli.RunWithOpts([]string{"delete", "-a", "test-package"},
				RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
			kappCtrl.RunWithOpts([]string{"pkg", "installed", "delete", "-i", "test"},
				RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
		}
		defer cleanUpInstalledPkg()

		kappCli.RunWithOpts([]string{"deploy", "-a", "test-package", "-f", filepath.Join(workingDir, pkgDir), "-c"},
			RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
		kappCtrl.RunWithOpts([]string{"pkg", "available", "list"},
			RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
		kappCtrl.RunWithOpts([]string{"pkg", "install", "-p", "testpackage.corp.dev", "-i", "test", "--version", "1.0.0"},
			RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
	})
}

func TestPackageInitAndReleaseWithTag(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kctrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	const pkgDir = "./carvel-artifacts/packages/testpackage.corp.dev/"
	promptOutput := newPromptOutput(t)

	cleanUp := func() {
		os.RemoveAll(workingDir)
		kapp.RunWithOpts([]string{"delete", "-a", "test-package"},
			RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
		kctrl.RunWithOpts([]string{"pkg", "installed", "delete", "-i", "test"},
			RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
	}
	cleanUp()
	defer cleanUp()

	err := os.Mkdir(workingDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	logger.Section("Package init", func() {
		interaction := Interaction{
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
		}

		go interaction.Run(promptOutput)
		kctrl.RunWithOpts([]string{"pkg", "init", "--tty=true", "--chdir", workingDir},
			RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
				StdoutWriter: promptOutput.BufferedOutputWriter(), Interactive: true})

		// Error if upstream folder doesn't exist
		_, err = os.Stat(filepath.Join(workingDir, "upstream"))
		require.NoError(t, err)
	})

	logger.Section("Package release", func() {
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
      - ytt:
          paths:
          - upstream
      - kbld:
          paths:
          - '-'
          - .imgpkg/images.yml
  valuesSchema:
    openAPIv3:
      default: null
      nullable: true
  version: 1.0.0
`
		releaseInteraction := Interaction{
			Prompts: []string{"Enter the registry URL"},
			Inputs:  []string{env.Image},
		}
		tag := "1.0.0"

		go releaseInteraction.Run(promptOutput)
		kctrl.RunWithOpts([]string{"pkg", "release", "--version", "1.0.0", "--tty=true", "--chdir", workingDir, "--tag", tag},
			RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
				StdoutWriter: promptOutput.BufferedOutputWriter(), Interactive: true})

		out, err := readFile(pkgDir + "package.yml")
		require.NoErrorf(t, err, "Expected to read package.yml")

		// Below key's values will be changed during every run, hence removing them
		out = clearKeys([]string{"creationTimestamp:", "releasedAt:", "image:"}, strings.TrimSpace(replaceSpaces(out)))
		expectedPackage = strings.TrimSpace(replaceSpaces(expectedPackage))
		require.Equal(t, expectedPackage, out, "Expected Package to match")

		cmd := exec.Command("imgpkg", []string{"pull", "-b", fmt.Sprintf("%s:%s", env.Image, tag), "-o", filepath.Join(workingDir, "tmp")}...)
		err = cmd.Run()
		require.NoErrorf(t, err, "Expected imgpkg pull to succeed")
	})

	logger.Section("Testing and installing created Package", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", "test-package", "-f", filepath.Join(workingDir, pkgDir), "-c"},
			RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
		kctrl.RunWithOpts([]string{"pkg", "available", "list"},
			RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
		kctrl.RunWithOpts([]string{"pkg", "install", "-p", "testpackage.corp.dev", "-i", "test", "--version", "1.0.0"},
			RunOpts{StdinReader: promptOutput.StringReader(), StdoutWriter: promptOutput.BufferedOutputWriter()})
	})
}

func TestPackageInitAndReleaseWithInvalidTag(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kctrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}
	promptOutput := newPromptOutput(t)

	os.RemoveAll(workingDir)
	defer os.RemoveAll(workingDir)

	err := os.Mkdir(workingDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	logger.Section("Package init", func() {
		interaction := Interaction{
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
		}

		go interaction.Run(promptOutput)
		kctrl.RunWithOpts([]string{"pkg", "init", "--tty=true", "--chdir", workingDir},
			RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
				StdoutWriter: promptOutput.BufferedOutputWriter(), Interactive: true})

		// Error if upstream folder doesn't exist
		_, err = os.Stat(filepath.Join(workingDir, "upstream"))
		require.NoError(t, err)
	})

	logger.Section("Package release", func() {
		releaseInteraction := Interaction{
			Prompts: []string{"Enter the registry URL"},
			Inputs:  []string{env.Image},
		}
		tag := "1.0.0+"

		go releaseInteraction.Run(promptOutput)
		_, err := kctrl.RunWithOpts([]string{"pkg", "release", "--version", "1.0.0", "--tty=true", "--chdir", workingDir, "--tag", tag},
			RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
				StdoutWriter: promptOutput.BufferedOutputWriter(), Interactive: true, AllowError: true})

		require.Error(t, err)
	})
}
