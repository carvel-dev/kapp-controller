// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	goexec "os/exec"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	"k8s.io/client-go/kubernetes"
	kyaml "sigs.k8s.io/yaml"
)

type HelmTemplate struct {
	opts        v1alpha1.AppTemplateHelmTemplate
	genericOpts GenericOpts
	coreClient  kubernetes.Interface
}

// HelmTemplateCmdArgs represents the binary and arguments used during templating
type HelmTemplateCmdArgs struct {
	BinaryName string
	Args       []string
}

var _ Template = &HelmTemplate{}

func NewHelmTemplate(opts v1alpha1.AppTemplateHelmTemplate,
	genericOpts GenericOpts, coreClient kubernetes.Interface) *HelmTemplate {
	return &HelmTemplate{opts, genericOpts, coreClient}
}

func (t *HelmTemplate) TemplateDir(dirPath string) (exec.CmdRunResult, bool) {
	chartPath := dirPath

	if len(t.opts.Path) > 0 {
		checkedPath, err := memdir.ScopedPath(dirPath, t.opts.Path)
		if err != nil {
			return exec.NewCmdRunResultWithErr(err), true
		}
		chartPath = checkedPath
	}

	name := t.genericOpts.Name
	if len(t.opts.Name) > 0 {
		name = t.opts.Name
	}

	namespace := t.genericOpts.Namespace
	if len(t.opts.Namespace) > 0 {
		namespace = t.opts.Namespace
	}

	// Return Helm binary name and arguments based on the version defined on the Chart.yaml file.
	// NOTE: This will be removed once we remove retro-compatibility with Helm V2 binary
	helmCmdCtx, err := NewHelmTemplateCmdArgs(name, chartPath, namespace)
	if err != nil {
		return exec.NewCmdRunResultWithErr(err), true
	}

	// Actual helm template arguments
	args := helmCmdCtx.Args

	{ // Add values files
		vals := Values{t.opts.ValuesFrom, t.genericOpts, t.coreClient}

		paths, valuesCleanUpFunc, err := vals.AsPaths(dirPath)
		if err != nil {
			return exec.NewCmdRunResultWithErr(err), true
		}

		defer valuesCleanUpFunc()

		for _, path := range paths {
			args = append(args, []string{"--values", path}...)
		}
	}

	var stdoutBs, stderrBs bytes.Buffer

	cmd := goexec.Command(helmCmdCtx.BinaryName, args...)
	// "Reset" kubernetes vars just in case, even though helm template should not reach out to cluster
	cmd.Env = append(os.Environ(), "KUBERNETES_SERVICE_HOST=not-real", "KUBERNETES_SERVICE_PORT=not-real")
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err = cmd.Run()

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Templating helm chart: %s", err)

	return result, true
}

func (t *HelmTemplate) TemplateStream(_ io.Reader, _ string) exec.CmdRunResult {
	return exec.NewCmdRunResultWithErr(fmt.Errorf("Templating data is not supported"))
}

// auxiliary struct used for Chart.yaml unmarshalling
type chartSpec struct {
	APIVersion string
}

// DEPRECATED. This method will not be required once support for Helm 2 is dropped
// Returns the Helm Binary Name and the arguments required to be passed to the "helm template" subcommand
// The returned values depend on the ApiVersion property inside the Chart.yaml file.
// apiVersion==v1 will fallback to old Helm 2 binary and command format.
func NewHelmTemplateCmdArgs(releaseName, chartPath, namespace string) (*HelmTemplateCmdArgs, error) {
	const (
		helmBinaryName = "helm"

		helm2BinaryName       = "helmv2" // DEPRECATED
		helm2ChartSpecVersion = "v1"     // DEPRECATED
	)

	// Load [chartPath]/Chart.yaml and inspect apiVersion value.
	bs, err := ioutil.ReadFile(filepath.Join(chartPath, "Chart.yaml"))
	if err != nil {
		return nil, fmt.Errorf("Reading Chart.yaml: %w", err)
	}

	var chartSpec chartSpec

	err = kyaml.Unmarshal(bs, &chartSpec)
	if err != nil {
		return nil, fmt.Errorf("Unmarshaling Chart.yaml: %w", err)
	}

	// By default, use Helm 3+ format except for chart.apiSpec=v1
	res := &HelmTemplateCmdArgs{
		BinaryName: helmBinaryName,
		Args:       []string{"template", releaseName, chartPath, "--namespace", namespace, "--include-crds"},
	}

	// DEPRECATED: Helm V2 will be removed in a future release
	if chartSpec.APIVersion == helm2ChartSpecVersion {
		res.BinaryName = helm2BinaryName
		res.Args = []string{"template", chartPath, "--name", releaseName, "--namespace", namespace}
	}

	return res, nil
}
