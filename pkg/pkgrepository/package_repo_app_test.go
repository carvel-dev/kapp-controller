package pkgrepository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
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
