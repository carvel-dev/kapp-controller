// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package local

import (
	"context"
	"testing"

	"carvel.dev/kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgv1alpha1 "carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1"
	dpv1alpha1 "carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	fakedpkg "carvel.dev/kapp-controller/pkg/apiserver/client/clientset/versioned/fake"
	fakekc "carvel.dev/kapp-controller/pkg/client/clientset/versioned/fake"
	"carvel.dev/kapp-controller/pkg/exec"
	versionsv1alpha1 "carvel.dev/vendir/pkg/vendir/versions/v1alpha1"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes"
	clienttesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const fakeClusterVersion = "v1.27.3"

// fakeClusterClient only answers discovery requests
type fakeClusterClient struct {
	kubernetes.Interface
	discovery discovery.DiscoveryInterface
}

func (c fakeClusterClient) Discovery() discovery.DiscoveryInterface { return c.discovery }

func newFakeClusterClient() fakeClusterClient {
	return fakeClusterClient{discovery: &fakediscovery.FakeDiscovery{
		Fake:               &clienttesting.Fake{Resources: []*metav1.APIResourceList{{GroupVersion: "example.com/v1"}}},
		FakedServerVersion: &version.Info{GitVersion: fakeClusterVersion},
	}}
}

func TestDiscoveryUsesClusterWhenClientIsAvailable(t *testing.T) {
	discoveryClient := (&MinCoreClient{client: newFakeClusterClient()}).Discovery()

	serverVersion, err := discoveryClient.ServerVersion()
	require.NoError(t, err)
	require.Equal(t, fakeClusterVersion, serverVersion.GitVersion)

	groups, err := discoveryClient.ServerGroups()
	require.NoError(t, err)
	require.Equal(t, []string{"example.com/v1"}, metav1.ExtractGroupVersions(groups))
}

func TestDiscoveryIsMockedWhenClientIsNotAvailable(t *testing.T) {
	discoveryClient := (&MinCoreClient{}).Discovery()
	require.IsType(t, &MinDiscoveryClient{}, discoveryClient)

	_, err := discoveryClient.ServerVersion()
	require.NoError(t, err)
}

func TestPackageInstallChecksKubernetesVersionConstraintsAgainstCluster(t *testing.T) {
	tests := []struct {
		name        string
		constraints string
		satisfied   bool
	}{
		{name: "constraint includes cluster version", constraints: ">=1.23.0 <1.28.0", satisfied: true},
		{name: "constraint excludes cluster version", constraints: ">=1.28.0", satisfied: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pkgi := &pkgv1alpha1.PackageInstall{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pkgi", Namespace: "test-ns"},
				Spec: pkgv1alpha1.PackageInstallSpec{
					ServiceAccountName: "test-sa",
					PackageRef: &pkgv1alpha1.PackageRef{
						RefName:          "test-pkg.carvel.dev",
						VersionSelection: &versionsv1alpha1.VersionSelectionSemver{Constraints: "1.0.0"},
					},
				},
			}
			pkg := &dpv1alpha1.Package{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pkg.carvel.dev.1.0.0", Namespace: "test-ns"},
				Spec: dpv1alpha1.PackageSpec{
					RefName:                    "test-pkg.carvel.dev",
					Version:                    "1.0.0",
					KubernetesVersionSelection: &dpv1alpha1.VersionSelection{Constraints: test.constraints},
					Template:                   dpv1alpha1.AppTemplateSpec{Spec: &kcv1alpha1.AppSpec{}},
				},
			}

			kcClient := fakekc.NewSimpleClientset(pkgi)
			coreClient := &MinCoreClient{client: newFakeClusterClient(), localSecrets: &localSecrets{}}
			_, pkgiReconciler := NewReconciler(nil, exec.NewPlainCmdRunner(), logger.NewNoopLogger()).newReconcilers(
				coreClient, kcClient, fakedpkg.NewSimpleClientset(pkg), nil, ReconcileOpts{})

			_, err := pkgiReconciler.Reconcile(context.TODO(), reconcile.Request{
				NamespacedName: types.NamespacedName{Name: pkgi.Name, Namespace: pkgi.Namespace},
			})

			_, appErr := kcClient.KappctrlV1alpha1().Apps(pkgi.Namespace).Get(context.TODO(), pkgi.Name, metav1.GetOptions{})
			if test.satisfied {
				require.NoError(t, err)
				require.NoError(t, appErr)
			} else {
				require.ErrorContains(t, err, "kubernetes-version-check")
				require.True(t, errors.IsNotFound(appErr))
			}
		})
	}
}
