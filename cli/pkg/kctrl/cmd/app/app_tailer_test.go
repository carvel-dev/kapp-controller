// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0
package app

import (
	"testing"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/stretchr/testify/require"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	internalv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/typed/internalpackaging/v1alpha1"
	kappctrlv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/typed/kappctrl/v1alpha1"
	packagingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/typed/packaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	discovery "k8s.io/client-go/discovery"
)

func TestAppSuccessTail(t *testing.T) {
	fakeVersionedInterface := &FakeVersionedInterface{t}

	successStatus := kcv1alpha1.AppStatus{
		Fetch: &kcv1alpha1.AppStatusFetch{
			UpdatedAt: metav1.Now(),
			StartedAt: metav1.Now(),
			ExitCode:  0,
			Stdout:    "vendir success",
		},
		Template: &kcv1alpha1.AppStatusTemplate{
			UpdatedAt: metav1.Now(),
			ExitCode:  0,
		},
		Deploy: &kcv1alpha1.AppStatusDeploy{
			UpdatedAt: metav1.Now(),
			StartedAt: metav1.Now(),
			ExitCode:  0,
			Stdout:    "kapp success",
		},
	}
	appTailer := NewAppTailer("default", "test-app", ui.NewNoopUI(), fakeVersionedInterface, AppTailerOpts{})
	appTailer.stopperChan = make(chan struct{})
	err := appTailer.printTillCurrent(successStatus)

	require.NoError(t, err)
}

func TestAppFetchFail(t *testing.T) {
	fakeVersionedInterface := &FakeVersionedInterface{t}

	successStatus := kcv1alpha1.AppStatus{
		Fetch: &kcv1alpha1.AppStatusFetch{
			UpdatedAt: metav1.Now(),
			StartedAt: metav1.Now(),
			ExitCode:  1,
			Stderr:    "vendir fail",
		},
	}
	appTailer := NewAppTailer("default", "test-app", ui.NewNoopUI(), fakeVersionedInterface, AppTailerOpts{})
	appTailer.stopperChan = make(chan struct{})
	err := appTailer.printTillCurrent(successStatus)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Fetch failed")
}

func TestAppTemplateFail(t *testing.T) {
	fakeVersionedInterface := &FakeVersionedInterface{t}

	successStatus := kcv1alpha1.AppStatus{
		Fetch: &kcv1alpha1.AppStatusFetch{
			UpdatedAt: metav1.Now(),
			StartedAt: metav1.Now(),
			ExitCode:  0,
			Stdout:    "vendir success",
		},
		Template: &kcv1alpha1.AppStatusTemplate{
			UpdatedAt: metav1.Now(),
			ExitCode:  1,
			Stderr:    "ytt fail",
		},
	}
	appTailer := NewAppTailer("default", "test-app", ui.NewNoopUI(), fakeVersionedInterface, AppTailerOpts{})
	appTailer.stopperChan = make(chan struct{})
	err := appTailer.printTillCurrent(successStatus)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Template failed")
}

func TestAppDeployFail(t *testing.T) {
	fakeVersionedInterface := &FakeVersionedInterface{t}

	successStatus := kcv1alpha1.AppStatus{
		Fetch: &kcv1alpha1.AppStatusFetch{
			UpdatedAt: metav1.Now(),
			StartedAt: metav1.Now(),
			ExitCode:  0,
			Stdout:    "vendir success",
		},
		Template: &kcv1alpha1.AppStatusTemplate{
			UpdatedAt: metav1.Now(),
			ExitCode:  0,
		},
		Deploy: &kcv1alpha1.AppStatusDeploy{
			UpdatedAt: metav1.Now(),
			StartedAt: metav1.Now(),
			Finished:  true,
			ExitCode:  1,
			Stderr:    "kapp fail",
		},
	}
	appTailer := NewAppTailer("default", "test-app", ui.NewNoopUI(), fakeVersionedInterface, AppTailerOpts{})
	appTailer.stopperChan = make(chan struct{})
	err := appTailer.printTillCurrent(successStatus)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Deploy failed")
}

type FakeVersionedInterface struct {
	t *testing.T
}

func (c *FakeVersionedInterface) Discovery() discovery.DiscoveryInterface { return nil }
func (c *FakeVersionedInterface) InternalV1alpha1() internalv1alpha1.InternalV1alpha1Interface {
	return nil
}
func (c *FakeVersionedInterface) KappctrlV1alpha1() kappctrlv1alpha1.KappctrlV1alpha1Interface {
	return nil
}
func (c *FakeVersionedInterface) PackagingV1alpha1() packagingv1alpha1.PackagingV1alpha1Interface {
	return nil
}
