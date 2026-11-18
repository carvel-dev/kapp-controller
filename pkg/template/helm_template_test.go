// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package template_test

import (
	"fmt"
	"os"
	goexec "os/exec"
	"path/filepath"
	"testing"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"carvel.dev/kapp-controller/pkg/exec"
	"carvel.dev/kapp-controller/pkg/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

type recordingCmdRunner struct {
	cmds []*goexec.Cmd
}

var _ exec.CmdRunner = &recordingCmdRunner{}

func (r *recordingCmdRunner) Run(cmd *goexec.Cmd) error {
	r.cmds = append(r.cmds, cmd)
	return nil
}

func (r *recordingCmdRunner) RunWithCancel(cmd *goexec.Cmd, _ chan struct{}) error {
	return r.Run(cmd)
}

func Test_HelmTemplate_NameAndNamespace(t *testing.T) {
	appContext := template.AppContext{Name: "app-name", Namespace: "app-ns"}

	testCases := []struct {
		desc         string
		opts         v1alpha1.AppTemplateHelmTemplate
		expectedName string
		expectedNs   string
		expectedErr  string
	}{
		{desc: "defaults to app name and namespace", expectedName: "app-name", expectedNs: "app-ns"},
		{desc: "custom name and namespace", opts: v1alpha1.AppTemplateHelmTemplate{Name: "my-release", Namespace: "my-ns"},
			expectedName: "my-release", expectedNs: "my-ns"},
		{desc: "name with dash in the middle", opts: v1alpha1.AppTemplateHelmTemplate{Name: "my-release-1"},
			expectedName: "my-release-1", expectedNs: "app-ns"},
		{desc: "name with leading double dash", opts: v1alpha1.AppTemplateHelmTemplate{Name: "--post-renderer=/bin/sh"},
			expectedErr: "invalid helm template name '--post-renderer=/bin/sh': must not start with '-'"},
		{desc: "name with leading single dash", opts: v1alpha1.AppTemplateHelmTemplate{Name: "-g"},
			expectedErr: "invalid helm template name '-g': must not start with '-'"},
		{desc: "namespace with leading dash", opts: v1alpha1.AppTemplateHelmTemplate{Namespace: "--debug"},
			expectedErr: "invalid helm template namespace '--debug': must not start with '-'"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cmdRunner := &recordingCmdRunner{}
			helmTemplate := template.NewHelmTemplate(tc.opts, appContext, k8sfake.NewSimpleClientset(),
				cmdRunner, template.AdditionalDownwardAPIValues{})

			result, _ := helmTemplate.TemplateDir(t.TempDir())

			if tc.expectedErr != "" {
				assert.ErrorContains(t, result.Error, tc.expectedErr)
				assert.Empty(t, cmdRunner.cmds, "helm should not be invoked")
				return
			}
			require.NoError(t, result.Error)
			require.Len(t, cmdRunner.cmds, 1)
			args := cmdRunner.cmds[0].Args
			assert.Equal(t, []string{"helm", "template", tc.expectedName}, args[:3])
			assert.Equal(t, []string{"--namespace", tc.expectedNs}, args[4:6])
		})
	}
}

func Test_HelmTemplate_NameIsNotParsedAsHelmFlag(t *testing.T) {
	rootDir := t.TempDir()
	chartDir := filepath.Join(rootDir, "chart")
	outsideFile := filepath.Join(rootDir, "outside.txt")
	require.NoError(t, os.MkdirAll(filepath.Join(chartDir, "templates"), 0700))
	require.NoError(t, os.WriteFile(filepath.Join(chartDir, "Chart.yaml"),
		[]byte("apiVersion: v2\nname: test-chart\nversion: 0.1.0\n"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(chartDir, "templates", "cm.yaml"), []byte(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm
data:
  value: {{ .Values.value | default "" | quote }}
`), 0600))
	require.NoError(t, os.WriteFile(outsideFile, []byte("content-outside-chart-dir"), 0600))

	opts := v1alpha1.AppTemplateHelmTemplate{Name: fmt.Sprintf("--set-file=value=%s", outsideFile)}
	helmTemplate := template.NewHelmTemplate(opts, template.AppContext{Name: "app", Namespace: "default"},
		k8sfake.NewSimpleClientset(), exec.NewPlainCmdRunner(), template.AdditionalDownwardAPIValues{})

	result, _ := helmTemplate.TemplateDir(chartDir)

	assert.ErrorContains(t, result.Error, "must not start with '-'")
	assert.NotContains(t, result.Stdout, "content-outside-chart-dir")
}
