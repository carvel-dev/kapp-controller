// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

// Package template provides a factory pattern design of instantiating a new "templater" in kapp-controller
// some examples include cue, helm, and ytt
package template

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/olm/convert"
)

type olmRegistry struct {
	opts             v1alpha1.AppTemplateOLMRegistry
	coreClient       kubernetes.Interface
	appContext       AppContext
	cmdRunner        exec.CmdRunner
	additionalValues AdditionalDownwardAPIValues
}

var _ Template = &olmRegistry{}

func newOLMRegistry(opts v1alpha1.AppTemplateOLMRegistry, appContext AppContext,
	coreClient kubernetes.Interface, cmdRunner exec.CmdRunner) *olmRegistry {

	return &olmRegistry{opts: opts, appContext: appContext,
		coreClient: coreClient, cmdRunner: cmdRunner}
}

// TemplateDir works on directory returning templating result,
// and boolean indicating whether subsequent operations
// should operate on result, or continue operating on the directory
func (c *olmRegistry) TemplateDir(dirPath string) (exec.CmdRunResult, bool) {
	return c.template(dirPath, nil), true
}

func namespacedName(namespace, name string) string {
	if namespace == "" {
		return name
	}
	return fmt.Sprintf("%s/%s", namespace, name)
}

// TemplateStream works on a stream returning templating result.
// dirPath is provided for context from which to reference additional inputs.
func (c *olmRegistry) TemplateStream(_ io.Reader, _ string) exec.CmdRunResult {
	return exec.NewCmdRunResultWithErr(fmt.Errorf("Templating stream is not supported")) // TODO: Implement
}

func (c *olmRegistry) template(dirPath string, _ io.Reader) exec.CmdRunResult {
	bundleRoot := filepath.Join(dirPath, c.opts.BundleRoot)

	rv1Bundle, err := convert.LoadRegistryV1(os.DirFS(bundleRoot))
	if err != nil {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Loading registry+v1 bundle: %w", err))
	}

	plainManifestBundle, err := convert.Convert(*rv1Bundle, c.appContext.Namespace, c.opts.TargetNamespaces)
	if err != nil {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Converting registry+v1 bundle to plain manifests: %w", err))
	}

	buf := bytes.NewBuffer(nil)
	for i, obj := range plainManifestBundle.Objects {
		objData, err := yaml.Marshal(obj)
		if err != nil {
			return exec.NewCmdRunResultWithErr(fmt.Errorf("Encoding plain manifest object [%d] %v, %v: %w", i, obj.GetObjectKind().GroupVersionKind(), namespacedName(obj.GetNamespace(), obj.GetName()), err))
		}
		buf.WriteString("---\n")
		buf.Write(objData)
	}

	return exec.CmdRunResult{
		Stdout:   buf.String(),
		Finished: true,
	}
}
