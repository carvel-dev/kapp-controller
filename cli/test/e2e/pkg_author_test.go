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

func TestCreatePackageAuthorInteractively(t *testing.T) {
	keys := []string{"creationTimestamp:", "releasedAt:"}
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	cleanUp := func() {
		//TODO should it return an error
		os.Remove("package-build.yml")
	}
	cleanUp()
	defer cleanUp()

	logger.Section("Creating a package interactively using pkg build create", func() {
		promptOutput := newPromptOutput(t)

		go func() {
			promptOutput.WaitForPkgRefName()
			promptOutput.WritePkgRefName()
			promptOutput.WaitForPkgVersion()
			promptOutput.WritePkgVersion()
			promptOutput.WaitForPkgFetch()
			promptOutput.WritePkgFetch()
			promptOutput.WaitForHelmChartName()
			promptOutput.WriteHelmChartName()
			promptOutput.WaitForHelmChartVersion()
			promptOutput.WriteHelmChartVersion()
			promptOutput.WaitForHelmChartURL()
			promptOutput.WriteHelmChartURL()
			promptOutput.WaitForImgpkgBundleURL()
			promptOutput.WriteImgpkgBundleURL()
		}()

		kappCtrl.RunWithOpts([]string{"pkg", "build", "create"},
			RunOpts{NoNamespace: true, StdinReader: promptOutput.StringReader(),
				StdoutWriter: promptOutput.OutputWriter(), Interactive: true})

		expectedOutput := `
annotations:
  fetch-content-from: Helm Chart from Helm Repository
apiVersion: kctrl.carvel.dev/v1alpha1
kind: PackageBuild
spec:
  imgpkg:
    registryUrl: docker.io/rohitagg2020/mongodb-bundle:1.0.0
  package:
    apiVersion: data.packaging.carvel.dev/v1alpha1
    kind: Package
    metadata:
      name: afc.def.ghi.2.0.0
      namespace: default
    spec:
      licenses:
      - Apache 2.0
      refName: afc.def.ghi
      releaseNotes: |
        Initial release of the simple app package
      template:
        spec:
          deploy:
          - kapp: {}
          fetch:
          - imgpkgBundle:
              image: index.docker.io/rohitagg2020/mongodb-bundle@sha256:21eb6ecef9c7dec256e8722a1da029915e150478101fbb3e0fed58c92fb68e73
          template:
          - helmTemplate:
              path: config/upstream
          - kbld:
              paths:
              - '-'
              - .imgpkg/images.yml
          - ytt:
              paths:
              - '-'
      valuesSchema:
        openAPIv3: null
      version: 2.0.0
  packageMetadata:
    apiVersion: data.packaging.carvel.dev/v1alpha1
    kind: PackageMetadata
    metadata:
      name: afc.def.ghi
      namespace: default
    spec:
      categories:
      - demo
      displayName: afc
      longDescription: Simple app consisting of a k8s deployment and service
      shortDescription: Simple app for demoing
  vendir:
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
      path: config/upstream
    kind: Config
    minimumRequiredVersion: 0.12.0
`
		out, err := readPackageBuild()
		if err != nil {
			fmt.Println(err.Error())
		}
		out = strings.TrimSpace(replaceSpaces(out))
		expectedOutput = strings.TrimSpace(replaceSpaces(expectedOutput))
		out = clearKeys(keys, out)
		require.Equal(t, expectedOutput, out, "output does not match")

	})
}

func replaceSpaces(result string) string {
	// result = strings.Replace(result, " ", "_", -1) // useful for debugging
	result = strings.Replace(result, " \n", " $\n", -1) // explicit endline
	return result
}

func readPackageBuild() (string, error) {
	data, err := os.ReadFile("package-build.yml")
	if err != nil {
		return "", err
	}
	return string(data), nil

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

func (p promptOutput) WaitForPkgRefName() {
	scanner := bufio.NewScanner(p.outputReader)
	for scanner.Scan() {
		// Cannot easily wait for prompt as it's not NL terminated
		if strings.Contains(scanner.Text(), "A package Reference name must be a valid DNS subdomain name") {
			break
		}
	}
	err := scanner.Err()
	require.NoError(p.t, err)
}

func (p promptOutput) WritePkgVersion() {
	p.stringWriter.Write([]byte("2.0.0\n"))
}

func (p promptOutput) WaitForPkgVersion() {
	scanner := bufio.NewScanner(p.outputReader)
	for scanner.Scan() {
		// Cannot easily wait for prompt as it's not NL terminated
		if strings.Contains(scanner.Text(), "It must be valid semver as specified") {
			break
		}
	}
	err := scanner.Err()
	require.NoError(p.t, err)
}

func (p promptOutput) WritePkgFetch() {
	p.stringWriter.Write([]byte("3\n"))

}

func (p promptOutput) WaitForPkgFetch() {
	scanner := bufio.NewScanner(p.outputReader)
	for scanner.Scan() {
		// Cannot easily wait for prompt as it's not NL terminated
		if strings.Contains(scanner.Text(), "we need to fetch the manifest which defines") {
			break
		}
	}
	err := scanner.Err()
	require.NoError(p.t, err)
}

func (p promptOutput) WriteHelmChartName() {
	p.stringWriter.Write([]byte("enterprise-operator\n"))

}

func (p promptOutput) WaitForHelmChartName() {
	scanner := bufio.NewScanner(p.outputReader)
	for scanner.Scan() {
		// Cannot easily wait for prompt as it's not NL terminated
		if strings.Contains(scanner.Text(), "directory will contain the bundle’s lock file") {
			break
		}
	}
	err := scanner.Err()
	require.NoError(p.t, err)
}

func (p promptOutput) WriteHelmChartVersion() {
	p.stringWriter.Write([]byte("1.16.0\n"))

}

func (p promptOutput) WaitForHelmChartVersion() {
	scanner := bufio.NewScanner(p.outputReader)
	for scanner.Scan() {
		// Cannot easily wait for prompt as it's not NL terminated
		if strings.Contains(scanner.Text(), "Enter helm chart name") {
			break
		}
	}
	err := scanner.Err()
	require.NoError(p.t, err)
}
func (p promptOutput) WriteHelmChartURL() {
	p.stringWriter.Write([]byte("https://mongodb.github.io/helm-charts\n"))

}

func (p promptOutput) WaitForHelmChartURL() {
	scanner := bufio.NewScanner(p.outputReader)
	for scanner.Scan() {
		// Cannot easily wait for prompt as it's not NL terminated
		if strings.Contains(scanner.Text(), "Enter helm chart version") {
			break
		}
	}
	err := scanner.Err()
	require.NoError(p.t, err)
}

func (p promptOutput) WaitForImgpkgBundleURL() {
	scanner := bufio.NewScanner(p.outputReader)
	for scanner.Scan() {
		// Cannot easily wait for prompt as it's not NL terminated
		if strings.Contains(scanner.Text(), "imgpkg bundle created above to the OCI registry") {
			break
		}
	}
	err := scanner.Err()
	require.NoError(p.t, err)
}

func (p promptOutput) WriteImgpkgBundleURL() {
	p.stringWriter.Write([]byte("docker.io/rohitagg2020/mongodb-bundle:1.0.0\n"))

}
