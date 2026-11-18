// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"os"
	"path/filepath"
	"testing"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"carvel.dev/kapp-controller/pkg/deploy"
	"carvel.dev/kapp-controller/pkg/exec"
	"carvel.dev/kapp-controller/pkg/fetch"
	"carvel.dev/kapp-controller/pkg/kubeconfig"
	"carvel.dev/kapp-controller/pkg/metrics"
	"carvel.dev/kapp-controller/pkg/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func Test_Fetch_PathOutsideWorkingDir_IsRejected(t *testing.T) {
	log := logf.Log.WithName("kc")

	rootDir := t.TempDir()
	workingDir := filepath.Join(rootDir, "working-dir")
	outsideDir := filepath.Join(rootDir, "outside")
	outsideFile := filepath.Join(outsideDir, "keep.txt")
	require.NoError(t, os.Mkdir(workingDir, 0700))
	require.NoError(t, os.Mkdir(outsideDir, 0700))
	require.NoError(t, os.WriteFile(outsideFile, []byte("original"), 0600))

	appSpec := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "simple-app",
			Namespace: "default",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{{
				Inline: &v1alpha1.AppFetchInline{Paths: map[string]string{"file.yml": "foo: bar"}},
				Path:   "../outside",
			}},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)

	app := NewApp(appSpec, Hooks{}, fetchFac, tmpFac, deployFac, log, Opts{}, metrics.NewMetrics(), FakeComponentInfo{})

	_, result := app.fetch(workingDir)
	assert.ErrorContains(t, result.Error, "invalid fetch path '../outside'")

	contents, err := os.ReadFile(outsideFile)
	require.NoError(t, err, "file outside the working directory should not be removed")
	assert.Equal(t, "original", string(contents))
}
