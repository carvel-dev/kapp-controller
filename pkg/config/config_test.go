// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	kcconfig "github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func Test_GetConfig_ReturnsSecret_WhenBothConfigMapAndSecretExist(t *testing.T) {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"httpProxy": "wrong-proxy",
		},
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"httpProxy": []byte("proxy-svc.proxy-server.svc.cluster.local:80"),
		},
	}

	defer os.Unsetenv("http_proxy")

	k8scs := k8sfake.NewSimpleClientset(configMap, secret)

	config, err := kcconfig.GetConfig(k8scs)
	assert.Nil(t, err, "unexpected error after running config.GetConfig()", err)

	assert.Nil(t, config.Apply(), "unexpected error after running config.Apply()", err)

	expected := "proxy-svc.proxy-server.svc.cluster.local:80"
	httpProxyActual := os.Getenv("http_proxy")

	assert.Equal(t, expected, httpProxyActual)
}

func Test_GetConfig_ReturnsConfigMap_WhenOnlyConfigMapExists(t *testing.T) {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"httpProxy": "proxy-svc.proxy-server.svc.cluster.local:80",
		},
	}

	defer os.Unsetenv("http_proxy")

	k8scs := k8sfake.NewSimpleClientset(configMap)

	config, err := kcconfig.GetConfig(k8scs)
	assert.Nil(t, err, "unexpected error after running config.GetConfig()", err)

	assert.Nil(t, config.Apply(), "unexpected error after running config.Apply()", err)

	expected := "proxy-svc.proxy-server.svc.cluster.local:80"
	httpProxyActual := os.Getenv("http_proxy")

	assert.Equal(t, expected, httpProxyActual)
}

func Test_GetConfig_ReturnsSecret_WhenOnlySecretExists(t *testing.T) {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"httpProxy": []byte("proxy-svc.proxy-server.svc.cluster.local:80"),
		},
	}

	defer os.Unsetenv("http_proxy")

	k8scs := k8sfake.NewSimpleClientset(secret)

	config, err := kcconfig.GetConfig(k8scs)
	assert.Nil(t, err, "unexpected error after running config.GetConfig()", err)

	assert.Nil(t, config.Apply(), "unexpected error after running config.Apply()", err)

	expected := "proxy-svc.proxy-server.svc.cluster.local:80"
	httpProxyActual := os.Getenv("http_proxy")

	assert.Equal(t, expected, httpProxyActual)
}

func Test_KubernetesServiceHost_IsSet(t *testing.T) {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"noProxy": []byte("KAPPCTRL_KUBERNETES_SERVICE_HOST"),
		},
	}

	defer os.Unsetenv("no_proxy")
	defer os.Unsetenv("KUBERNETES_SERVICE_HOST")

	os.Setenv("KUBERNETES_SERVICE_HOST", "10.96.0.1")

	k8scs := k8sfake.NewSimpleClientset(secret)

	config, err := kcconfig.GetConfig(k8scs)
	assert.Nil(t, err, "unexpected error after running config.GetConfig()", err)

	assert.Nil(t, config.Apply(), "unexpected error after running config.Apply()", err)

	expected := "10.96.0.1"
	noProxyActual := os.Getenv("no_proxy")

	assert.Equal(t, expected, noProxyActual)
}

func Test_ShouldSkipTLSForDomain(t *testing.T) {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"dangerousSkipTLSVerify": "always.trustworthy.com, selectively.trusted.net:123456",
		},
	}
	k8scs := k8sfake.NewSimpleClientset(configMap)
	config, err := kcconfig.GetConfig(k8scs)
	assert.NoError(t, err)

	assert.False(t, config.ShouldSkipTLSForDomain("some.random.org"))
	assert.True(t, config.ShouldSkipTLSForDomain("always.trustworthy.com"))
	assert.False(t, config.ShouldSkipTLSForDomain("selectively.trusted.net"))
	assert.True(t, config.ShouldSkipTLSForDomain("selectively.trusted.net:123456"))
}
