// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"bytes"
	"fmt"
	"io"
	"os"
	goexec "os/exec"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	"k8s.io/client-go/kubernetes"
)

type HelmTemplate struct {
	opts       v1alpha1.AppTemplateHelmTemplate
	appContext AppContext
	coreClient kubernetes.Interface
	cmdRunner  exec.CmdRunner
}

// HelmTemplateCmdArgs represents the binary and arguments used during templating
type HelmTemplateCmdArgs struct {
	BinaryName string
	Args       []string
}

var _ Template = &HelmTemplate{}

func NewHelmTemplate(opts v1alpha1.AppTemplateHelmTemplate,
	appContext AppContext, coreClient kubernetes.Interface,
	cmdRunner exec.CmdRunner) *HelmTemplate {

	return &HelmTemplate{opts, appContext, coreClient, cmdRunner}
}

func (t *HelmTemplate) TemplateDir(dirPath string) (exec.CmdRunResult, bool) {
	return t.template(dirPath, nil), true
}

// TemplateStream works on a stream returning templating result.
// dirPath is provided for context from which to reference additional inputs.
func (t *HelmTemplate) TemplateStream(stream io.Reader, dirPath string) exec.CmdRunResult {
	return t.template(dirPath, stream)
}

func (t *HelmTemplate) template(dirPath string, input io.Reader) exec.CmdRunResult {
	chartPath := dirPath

	if len(t.opts.Path) > 0 {
		checkedPath, err := memdir.ScopedPath(dirPath, t.opts.Path)
		if err != nil {
			return exec.NewCmdRunResultWithErr(err)
		}
		chartPath = checkedPath
	}

	name := t.appContext.Name
	if len(t.opts.Name) > 0 {
		name = t.opts.Name
	}

	namespace := t.appContext.Namespace
	if len(t.opts.Namespace) > 0 {
		namespace = t.opts.Namespace
	}

	args := []string{"template", name, chartPath, "--namespace", namespace, "--include-crds"}

	{ // Add values files
		vals := Values{t.opts.ValuesFrom, t.appContext, t.coreClient}

		paths, valuesCleanUpFunc, err := vals.AsPaths(dirPath)
		if err != nil {
			return exec.NewCmdRunResultWithErr(err)
		}

		defer valuesCleanUpFunc()

		for _, path := range paths {
			if path == stdinPath && input == nil {
				return exec.NewCmdRunResultWithErr(
					fmt.Errorf("Expected stdin to be available when using it as path, but was not"))
			}
			args = append(args, []string{"--values", path}...)
		}
	}

	var stdoutBs, stderrBs bytes.Buffer

	cmd := goexec.Command("helm", args...)
	// "Reset" kubernetes vars just in case, even though helm template should not reach out to cluster
	cmd.Env = append(os.Environ(), "KUBERNETES_SERVICE_HOST=not-real", "KUBERNETES_SERVICE_PORT=not-real")
	cmd.Stdin = input
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := t.cmdRunner.Run(cmd)

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Templating helm chart: %s", err)

	return result
}
