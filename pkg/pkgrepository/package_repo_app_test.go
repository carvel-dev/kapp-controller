// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"testing"
	"time"

	pkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewPackageRepoApp(t *testing.T) {
	expectedSyncPeriod := metav1.Duration{15 * time.Minute}

	inputApp := pkgingv1alpha1.PackageRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ThisIsJustATest",
			Namespace: "Testing12",
		},
		Spec: pkgingv1alpha1.PackageRepositorySpec{
			SyncPeriod: &expectedSyncPeriod,
			Fetch:      &pkgingv1alpha1.PackageRepositoryFetch{},
		},
	}

	outputApp, err := NewPackageRepoApp(&inputApp)
	assert.NoError(t, err)
	assert.Equal(t, "ThisIsJustATest", outputApp.Name)
	assert.Equal(t, "Testing12", outputApp.Namespace)
	require.NotNil(t, outputApp.Spec.SyncPeriod)
	assert.Equal(t, expectedSyncPeriod, *outputApp.Spec.SyncPeriod)
}
