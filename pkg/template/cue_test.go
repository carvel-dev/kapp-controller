// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"os"
	goexec "os/exec"
	"path/filepath"
	"testing"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"carvel.dev/kapp-controller/pkg/exec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

type cueRecordingCmdRunner struct {
	cmds []*goexec.Cmd
}

var _ exec.CmdRunner = &cueRecordingCmdRunner{}

func (r *cueRecordingCmdRunner) Run(cmd *goexec.Cmd) error {
	r.cmds = append(r.cmds, cmd)
	return nil
}

func (r *cueRecordingCmdRunner) RunWithCancel(cmd *goexec.Cmd, _ chan struct{}) error {
	return r.Run(cmd)
}

func Test_Cue_Paths(t *testing.T) {
	testCases := []struct {
		desc         string
		paths        []string
		expectedArgs []string
		expectedErr  string
	}{
		{desc: "no paths", expectedArgs: []string{"cue", "export", "--out", "yaml", "."}},
		{desc: "relative paths", paths: []string{"file.cue", "./dir", "a/b.cue"},
			expectedArgs: []string{"cue", "export", "--out", "yaml", "file.cue", "./dir", "a/b.cue"}},
		{desc: "path with dash in the middle", paths: []string{"my-file.cue"},
			expectedArgs: []string{"cue", "export", "--out", "yaml", "my-file.cue"}},
		{desc: "path with leading double dash", paths: []string{"file.cue", "--outfile=/tmp/file"},
			expectedErr: "invalid cue path '--outfile=/tmp/file': must not start with '-'"},
		{desc: "path with leading single dash", paths: []string{"-f"},
			expectedErr: "invalid cue path '-f': must not start with '-'"},
		{desc: "path outside of directory", paths: []string{"../file.cue"},
			expectedErr: "Checking path: Invalid path: ../file.cue"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cmdRunner := &cueRecordingCmdRunner{}
			c := newCue(v1alpha1.AppTemplateCue{Paths: tc.paths}, AppContext{},
				k8sfake.NewSimpleClientset(), cmdRunner, AdditionalDownwardAPIValues{})

			result, _ := c.TemplateDir(t.TempDir())

			if tc.expectedErr != "" {
				assert.ErrorContains(t, result.Error, tc.expectedErr)
				assert.Empty(t, cmdRunner.cmds, "cue should not be invoked")
				return
			}
			require.NoError(t, result.Error)
			require.Len(t, cmdRunner.cmds, 1)
			assert.Equal(t, tc.expectedArgs, cmdRunner.cmds[0].Args[:len(tc.expectedArgs)])
		})
	}
}

func Test_Cue_PathIsNotParsedAsCueFlag(t *testing.T) {
	rootDir := t.TempDir()
	workingDir := filepath.Join(rootDir, "working-dir")
	outsideFile := filepath.Join(rootDir, "outside.txt")
	require.NoError(t, os.Mkdir(workingDir, 0700))
	require.NoError(t, os.WriteFile(filepath.Join(workingDir, "cm.cue"),
		[]byte("apiVersion: \"v1\"\nkind: \"ConfigMap\"\nmetadata: name: \"cm\"\n"), 0600))
	require.NoError(t, os.WriteFile(outsideFile, []byte("original"), 0600))

	opts := v1alpha1.AppTemplateCue{Paths: []string{"cm.cue", "--outfile=" + outsideFile, "--force"}}
	c := newCue(opts, AppContext{}, k8sfake.NewSimpleClientset(), exec.NewPlainCmdRunner(), AdditionalDownwardAPIValues{})

	result, _ := c.TemplateDir(workingDir)

	assert.ErrorContains(t, result.Error, "must not start with '-'")
	contents, err := os.ReadFile(outsideFile)
	require.NoError(t, err)
	assert.Equal(t, "original", string(contents), "file outside the working directory should not be modified")
}
