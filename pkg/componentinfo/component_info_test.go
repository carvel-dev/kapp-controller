// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package componentinfo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/componentinfo"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/kubeconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_prereleases_are_removed(t *testing.T) {
	fakek8s := fake.NewSimpleClientset()

	// mock the kubernetes server version
	fakeDiscovery, _ := fakek8s.Discovery().(*fakediscovery.FakeDiscovery)
	fakeDiscovery.FakedServerVersion = &version.Info{
		GitVersion: "v0.20.0-gke.100",
	}

	ci := componentinfo.NewComponentInfo(fakek8s, &kubeconfig.Kubeconfig{}, "0.40.0")

	version, err := ci.KubernetesVersion("saname", &v1alpha1.AppCluster{}, &metav1.ObjectMeta{})
	assert.NoError(t, err)
	assert.Equal(t, "0.20.0", version.String())
}
