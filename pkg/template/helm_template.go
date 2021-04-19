// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	goexec "os/exec"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	for _, source := range t.opts.ValuesFrom {
		var paths []string
		var err error

		valuesDir := memdir.NewTmpDir("helm-template-values")

		err = valuesDir.Create()
		if err != nil {
			result := exec.CmdRunResult{}
			result.AttachErrorf("Templating: %s", err)
			return result, true
		}

		defer valuesDir.Remove()

		switch {
		case source.SecretRef != nil:
			paths, err = t.writeFromSecret(valuesDir.Path(), *source.SecretRef)

		case source.ConfigMapRef != nil:
			paths, err = t.writeFromConfigMap(valuesDir.Path(), *source.ConfigMapRef)

		case len(source.Path) > 0:
			checkedPath, err := memdir.ScopedPath(dirPath, source.Path)
			if err == nil {
				paths = append(paths, checkedPath)
			}

		default:
			err = fmt.Errorf("Expected either secretRef, configMapRef or path as a source")
		}
		if err != nil {
			result := exec.CmdRunResult{}
			result.AttachErrorf("Writing paths: %s", err)
			return result, true
		}

		for _, path := range paths {
			args = append(args, []string{"--values", path}...)
		}
	}

	var stdoutBs, stderrBs bytes.Buffer

	cmd := goexec.Command(helmCmdCtx.BinaryName, args...)
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

func (t *HelmTemplate) writeFromSecret(dstPath string, secretRef v1alpha1.AppTemplateHelmTemplateValuesSourceRef) ([]string, error) {
	secret, err := t.coreClient.CoreV1().Secrets(t.genericOpts.Namespace).Get(secretRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var result []string

	for name, val := range secret.Data {
		path, err := t.writeFile(dstPath, name, val)
		if err != nil {
			return nil, err
		}
		result = append(result, path)
	}

	return result, nil
}

func (t *HelmTemplate) writeFromConfigMap(dstPath string, configMapRef v1alpha1.AppTemplateHelmTemplateValuesSourceRef) ([]string, error) {
	configMap, err := t.coreClient.CoreV1().ConfigMaps(t.genericOpts.Namespace).Get(configMapRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var result []string

	for name, val := range configMap.Data {
		path, err := t.writeFile(dstPath, name, []byte(val))
		if err != nil {
			return nil, err
		}
		result = append(result, path)
	}

	return result, nil
}

func (t *HelmTemplate) writeFile(dstPath, subPath string, content []byte) (string, error) {
	newPath, err := memdir.ScopedPath(dstPath, subPath)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(newPath, content, 0600)
	if err != nil {
		return "", fmt.Errorf("Writing file '%s': %s", newPath, err)
	}

	return newPath, nil
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
