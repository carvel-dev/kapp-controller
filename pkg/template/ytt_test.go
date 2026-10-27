// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"carvel.dev/kapp-controller/pkg/exec"
	"carvel.dev/kapp-controller/pkg/fetch"
	"carvel.dev/kapp-controller/pkg/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func Test_Ytt_PathThroughSymlinkOutsideDir_IsRejected(t *testing.T) {
	rootDir := t.TempDir()
	workingDir := filepath.Join(rootDir, "working-dir")
	outsideDir := filepath.Join(rootDir, "outside")
	require.NoError(t, os.MkdirAll(filepath.Join(workingDir, "config"), 0700))
	require.NoError(t, os.Mkdir(outsideDir, 0700))
	require.NoError(t, os.WriteFile(filepath.Join(workingDir, "config", "cm.yml"),
		[]byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: inside\n"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(outsideDir, "cm.yml"),
		[]byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: content-outside-working-dir\n"), 0600))
	require.NoError(t, os.Symlink("config", filepath.Join(workingDir, "link-to-config")))
	require.NoError(t, os.Symlink(outsideDir, filepath.Join(workingDir, "link-to-outside")))

	k8scs := k8sfake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	newYtt := func(paths []string) *template.Ytt {
		return template.NewYtt(v1alpha1.AppTemplateYtt{Paths: paths}, template.AppContext{Name: "app", Namespace: "default"},
			k8scs, fetchFac, exec.NewPlainCmdRunner(), template.AdditionalDownwardAPIValues{})
	}

	t.Run("symlink within working dir is allowed", func(t *testing.T) {
		result, _ := newYtt([]string{"link-to-config/cm.yml"}).TemplateDir(workingDir)
		require.NoError(t, result.Error)
		assert.Contains(t, result.Stdout, "name: inside")
	})

	t.Run("symlink pointing outside working dir is rejected", func(t *testing.T) {
		result, _ := newYtt([]string{"link-to-outside/cm.yml"}).TemplateDir(workingDir)
		assert.ErrorContains(t, result.Error, "resolves outside of the base directory through a symlink")
		assert.NotContains(t, result.Stdout, "content-outside-working-dir")
	})
}
