// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	goexec "os/exec"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/exec"
	"github.com/k14s/kapp-controller/pkg/memdir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type HelmTemplate struct {
	opts        v1alpha1.AppTemplateHelmTemplate
	genericOpts GenericOpts
	coreClient  kubernetes.Interface
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

	args := []string{
		"template", chartPath,
		"--namespace", t.genericOpts.Namespace,
		"--name", t.genericOpts.Name,
	}

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

	cmd := goexec.Command("helm", args...)
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := cmd.Run()

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
