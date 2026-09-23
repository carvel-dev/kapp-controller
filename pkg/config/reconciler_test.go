// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"context"
	"errors"
	"testing"

	kcconfig "carvel.dev/kapp-controller/pkg/config"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type fakeOSConfig struct {
	applyCACertsErr error
	applyProxyErr   error

	applyCACertsCalled bool
	applyProxyCalled   bool
}

func (f *fakeOSConfig) ApplyCACerts(string) error {
	f.applyCACertsCalled = true
	return f.applyCACertsErr
}

func (f *fakeOSConfig) ApplyProxy(kcconfig.ProxyOpts) error {
	f.applyProxyCalled = true
	return f.applyProxyErr
}

func newTestConfig(t *testing.T) *kcconfig.Config {
	t.Helper()

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"caCerts": []byte("ca-certs"),
		},
	}

	config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
	require.NoError(t, err)
	return config
}

// Test_Reconciler_Reconcile_ReturnsErr_WhenApplyCACertsFails ensures that a
// transient failure applying CA certs (e.g. the sidecarexec RPC connection
// not being ready yet) is surfaced as an error so that controller-runtime
// requeues the request with backoff, instead of being silently swallowed.
func Test_Reconciler_Reconcile_ReturnsErr_WhenApplyCACertsFails(t *testing.T) {
	osConfig := &fakeOSConfig{applyCACertsErr: errors.New("dial unix sidecarexec.sock: connect: no such file or directory")}

	reconciler := kcconfig.NewReconciler(
		k8sfake.NewSimpleClientset(), newTestConfig(t), osConfig, logr.Discard())

	_, err := reconciler.Reconcile(context.Background(), reconcile.Request{})
	assert.Error(t, err)
	assert.True(t, osConfig.applyCACertsCalled)
	assert.False(t, osConfig.applyProxyCalled, "should not attempt to apply proxy opts once CA certs failed to apply")
}

// Test_Reconciler_Reconcile_ReturnsErr_WhenApplyProxyFails mirrors the CA
// certs case above for proxy configuration.
func Test_Reconciler_Reconcile_ReturnsErr_WhenApplyProxyFails(t *testing.T) {
	osConfig := &fakeOSConfig{applyProxyErr: errors.New("dial unix sidecarexec.sock: connect: no such file or directory")}

	reconciler := kcconfig.NewReconciler(
		k8sfake.NewSimpleClientset(), newTestConfig(t), osConfig, logr.Discard())

	_, err := reconciler.Reconcile(context.Background(), reconcile.Request{})
	assert.Error(t, err)
	assert.True(t, osConfig.applyCACertsCalled)
	assert.True(t, osConfig.applyProxyCalled)
}

// Test_Reconciler_Reconcile_Succeeds_WhenOSConfigApplied ensures the happy
// path continues to succeed with no requeue.
func Test_Reconciler_Reconcile_Succeeds_WhenOSConfigApplied(t *testing.T) {
	osConfig := &fakeOSConfig{}

	reconciler := kcconfig.NewReconciler(
		k8sfake.NewSimpleClientset(), newTestConfig(t), osConfig, logr.Discard())

	result, err := reconciler.Reconcile(context.Background(), reconcile.Request{})
	assert.NoError(t, err)
	assert.Equal(t, reconcile.Result{}, result)
	assert.True(t, osConfig.applyCACertsCalled)
	assert.True(t, osConfig.applyProxyCalled)
}
