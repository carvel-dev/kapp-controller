// Copyright 2026 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"context"
	"fmt"
	"testing"

	kcv1alpha1 "carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1"
	"carvel.dev/kapp-controller/pkg/client/clientset/versioned/fake"
	"carvel.dev/kapp-controller/pkg/deploy"
	"carvel.dev/kapp-controller/pkg/exec"
	"carvel.dev/kapp-controller/pkg/fetch"
	"carvel.dev/kapp-controller/pkg/kubeconfig"
	"carvel.dev/kapp-controller/pkg/metrics"
	"carvel.dev/kapp-controller/pkg/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func TestUnblockDeletionRetriesOnConflict(t *testing.T) {
	pkgr := pkgrWithFinalizers()
	kappcs := fake.NewSimpleClientset(&pkgr)

	updateAttempts := 0
	kappcs.PrependReactor("update", "packagerepositories", func(_ k8stesting.Action) (bool, runtime.Object, error) {
		updateAttempts++
		if updateAttempts <= 2 {
			return true, nil, errors.NewConflict(schema.GroupResource{Group: "packaging.carvel.dev", Resource: "packagerepositories"},
				pkgr.Name, fmt.Errorf("the object has been modified; please apply your changes to the latest version and try again"))
		}
		// Fall through to the default reactor to persist the update
		return false, nil, nil
	})

	err := newTestCRDApp(t, &pkgr, kappcs).unblockDeletion()
	require.NoError(t, err)
	assert.Equal(t, 3, updateAttempts)

	updatedPkgr, err := kappcs.PackagingV1alpha1().PackageRepositories(pkgr.Namespace).Get(context.Background(), pkgr.Name, metav1.GetOptions{})
	require.NoError(t, err)
	assert.Empty(t, updatedPkgr.Finalizers)
}

func TestUnblockDeletionDoesNotRetryOnNonConflictError(t *testing.T) {
	pkgr := pkgrWithFinalizers()
	kappcs := fake.NewSimpleClientset(&pkgr)

	updateAttempts := 0
	kappcs.PrependReactor("update", "packagerepositories", func(_ k8stesting.Action) (bool, runtime.Object, error) {
		updateAttempts++
		return true, nil, errors.NewInternalError(fmt.Errorf("etcd unavailable"))
	})

	err := newTestCRDApp(t, &pkgr, kappcs).unblockDeletion()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Updating PackageRepository: Internal error occurred: etcd unavailable")
	assert.Equal(t, 1, updateAttempts)
}

func pkgrWithFinalizers() pkgingv1alpha1.PackageRepository {
	now := metav1.Now()
	return pkgingv1alpha1.PackageRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "repo",
			Namespace:         "default",
			DeletionTimestamp: &now,
			Finalizers:        []string{deleteFinalizerName, deletePrevFinalizerName},
		},
		Spec: pkgingv1alpha1.PackageRepositorySpec{
			Fetch: &pkgingv1alpha1.PackageRepositoryFetch{
				Inline: &kcv1alpha1.AppFetchInline{},
			},
		},
	}
}

func newTestCRDApp(t *testing.T, pkgr *pkgingv1alpha1.PackageRepository, kappcs *fake.Clientset) *CRDApp {
	log := logf.Log.WithName("kc")

	app, err := NewPackageRepoApp(pkgr)
	require.NoError(t, err)

	k8scs := k8sfake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)

	return NewCRDApp(app, pkgr, log, kappcs, metrics.NewMetrics(), fetchFac, tmpFac, deployFac)
}
