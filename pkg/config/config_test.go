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

	config, err := kcconfig.GetConfig(k8sfake.NewSimpleClientset(configMap, secret))
	assert.NoError(t, err)

	assert.Equal(t, kcconfig.ProxyOpts{
		HTTPProxy: "proxy-svc.proxy-server.svc.cluster.local:80",
	}, config.ProxyOpts())
}

func Test_GetConfig_KappDeployRawOptions(t *testing.T) {
	t.Run("with empty config value, returns just default", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{},
		}
		config, err := kcconfig.GetConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, []string{"--app-changes-max-to-keep=5"}, config.KappDeployRawOptions())
	})

	t.Run("with empty config value, returns just default", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"kappDeployRawOptions": []byte(""),
			},
		}
		config, err := kcconfig.GetConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, []string{"--app-changes-max-to-keep=5"}, config.KappDeployRawOptions())
	})

	t.Run("with populated config value, returns default and user set", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"kappDeployRawOptions": []byte("[\"--key=val\"]"),
			},
		}
		config, err := kcconfig.GetConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, []string{"--app-changes-max-to-keep=5", "--key=val"}, config.KappDeployRawOptions())
	})
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

	config, err := kcconfig.GetConfig(k8sfake.NewSimpleClientset(configMap))
	assert.NoError(t, err)

	assert.Equal(t, kcconfig.ProxyOpts{
		HTTPProxy: "proxy-svc.proxy-server.svc.cluster.local:80",
	}, config.ProxyOpts())
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

	config, err := kcconfig.GetConfig(k8sfake.NewSimpleClientset(secret))
	assert.NoError(t, err)

	assert.Equal(t, kcconfig.ProxyOpts{
		HTTPProxy: "proxy-svc.proxy-server.svc.cluster.local:80",
	}, config.ProxyOpts())
}

func Test_GetConfig(t *testing.T) {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"httpProxy":  []byte("http-proxy"),
			"httpsProxy": []byte("https-proxy"),
			"noProxy":    []byte("no-proxy"),
			"caCerts":    []byte("ca-certs"),
		},
	}

	config, err := kcconfig.GetConfig(k8sfake.NewSimpleClientset(secret))
	assert.NoError(t, err)

	assert.Equal(t, kcconfig.ProxyOpts{
		HTTPProxy:  "http-proxy",
		HTTPSProxy: "https-proxy",
		NoProxy:    "no-proxy",
	}, config.ProxyOpts())

	assert.Equal(t, "ca-certs", config.CACerts())
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

	os.Setenv("KUBERNETES_SERVICE_HOST", "10.96.0.1")
	defer os.Unsetenv("KUBERNETES_SERVICE_HOST")

	config, err := kcconfig.GetConfig(k8sfake.NewSimpleClientset(secret))
	assert.NoError(t, err)

	assert.Equal(t, kcconfig.ProxyOpts{NoProxy: "10.96.0.1"}, config.ProxyOpts())
}

func Test_ShouldSkipTLSForAuthority(t *testing.T) {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"dangerousSkipTLSVerify": "always.trustworthy.com, selectively.trusted.net:123456, [1fff:0:a88:85a3::ac1f]:8001, 1aaa:0:a88:85a3::ac1f",
		},
	}

	config, err := kcconfig.GetConfig(k8sfake.NewSimpleClientset(configMap))
	assert.NoError(t, err)

	assert.False(t, config.ShouldSkipTLSForAuthority("some.random.org"))
	assert.True(t, config.ShouldSkipTLSForAuthority("always.trustworthy.com"))
	assert.True(t, config.ShouldSkipTLSForAuthority("always.trustworthy.com:12345"))
	assert.False(t, config.ShouldSkipTLSForAuthority("selectively.trusted.net"))
	assert.False(t, config.ShouldSkipTLSForAuthority("selectively.trusted.net:8888"))
	assert.True(t, config.ShouldSkipTLSForAuthority("selectively.trusted.net:123456"))
	assert.True(t, config.ShouldSkipTLSForAuthority("[1fff:0:a88:85a3::ac1f]:8001"))
	assert.False(t, config.ShouldSkipTLSForAuthority("[1fff:0:a88:85a3::ac1f]:8888"))
	assert.False(t, config.ShouldSkipTLSForAuthority("1fff:0:a88:85a3::ac1f"))
	assert.True(t, config.ShouldSkipTLSForAuthority("1aaa:0:a88:85a3::ac1f"))
	assert.True(t, config.ShouldSkipTLSForAuthority("[1aaa:0:a88:85a3::ac1f]:888"))
}
