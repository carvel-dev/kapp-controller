// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	kcconfig "github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func Test_NewConfig_ReturnsSecret_WhenBothConfigMapAndSecretExist(t *testing.T) {
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

	config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(configMap, secret))
	assert.NoError(t, err)

	assert.Equal(t, kcconfig.ProxyOpts{
		HTTPProxy: "proxy-svc.proxy-server.svc.cluster.local:80",
	}, config.ProxyOpts())
}

func Test_NewConfig_PackageInstallDefaultSyncPeriod(t *testing.T) {
	t.Run("with empty config value, returns 10m", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{},
		}
		config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, 10*time.Minute, config.PackageInstallDefaultSyncPeriod())
	})

	t.Run("with value, returns 80s", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"packageInstallDefaultSyncPeriod": []byte("80s"),
			},
		}
		config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, 80*time.Second, config.PackageInstallDefaultSyncPeriod())
	})
}

func Test_NewConfig_AppDefaultSyncPeriod(t *testing.T) {
	t.Run("with empty config value, returns 30s", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{},
		}
		config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, 30*time.Second, config.AppDefaultSyncPeriod())
	})

	t.Run("with value", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"appDefaultSyncPeriod": []byte("1m20s"),
			},
		}
		config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, 80*time.Second, config.AppDefaultSyncPeriod())
	})

	t.Run("with too small of a value, returns 30s", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"appDefaultSyncPeriod": []byte("1s"),
			},
		}
		config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, 30*time.Second, config.AppDefaultSyncPeriod())
	})
}

func Test_NewConfig_AppMinimumSyncPeriod(t *testing.T) {
	t.Run("with empty config value, returns 30s", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{},
		}
		config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, 30*time.Second, config.AppMinimumSyncPeriod())
	})

	t.Run("with value", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"appMinimumSyncPeriod": []byte("1m20s"),
			},
		}
		config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, 80*time.Second, config.AppMinimumSyncPeriod())
	})

	t.Run("with too small of a value, returns 30s", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"appMinimumSyncPeriod": []byte("1s"),
			},
		}
		config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, 30*time.Second, config.AppMinimumSyncPeriod())
	})
}

func Test_NewConfig_KappDeployRawOptions(t *testing.T) {
	defaultRawOptions := []string{
		"--app-changes-max-to-keep=5", "--apply-timeout=5m", "--diff-anchored=true",
	}
	t.Run("with empty config value, returns just default", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{},
		}
		config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, defaultRawOptions, config.KappDeployRawOptions())
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
		config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, defaultRawOptions, config.KappDeployRawOptions())
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
		config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
		assert.NoError(t, err)
		assert.Equal(t, appendNewSlice(defaultRawOptions, "--key=val"), config.KappDeployRawOptions())
	})

	t.Run("clears previously set value when secret is gone", func(t *testing.T) {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kapp-controller-config",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"kappDeployRawOptions": []byte("[\"--key=val\"]"),
			},
		}
		client := k8sfake.NewSimpleClientset(secret)

		config, err := kcconfig.NewConfig(client)
		assert.NoError(t, err)
		assert.Equal(t, appendNewSlice(defaultRawOptions, "--key=val"), config.KappDeployRawOptions())

		err = client.CoreV1().Secrets("default").Delete(
			context.Background(), "kapp-controller-config", metav1.DeleteOptions{})
		assert.NoError(t, err)

		assert.NoError(t, config.Reload())
		assert.Equal(t, defaultRawOptions, config.KappDeployRawOptions())
	})
}

func appendNewSlice(act []string, items ...string) []string {
	newslice := make([]string, 0)
	newslice = append(newslice, act...)
	return append(newslice, items...)
}

func Test_NewConfig_ReturnsConfigMap_WhenOnlyConfigMapExists(t *testing.T) {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"httpProxy": "proxy-svc.proxy-server.svc.cluster.local:80",
		},
	}

	config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(configMap))
	assert.NoError(t, err)

	assert.Equal(t, kcconfig.ProxyOpts{
		HTTPProxy: "proxy-svc.proxy-server.svc.cluster.local:80",
	}, config.ProxyOpts())
}

func Test_NewConfig_ReturnsSecret_WhenOnlySecretExists(t *testing.T) {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"httpProxy": []byte("proxy-svc.proxy-server.svc.cluster.local:80"),
		},
	}

	config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
	assert.NoError(t, err)

	assert.Equal(t, kcconfig.ProxyOpts{
		HTTPProxy: "proxy-svc.proxy-server.svc.cluster.local:80",
	}, config.ProxyOpts())
}

func Test_NewConfig(t *testing.T) {
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

	config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
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

	config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(secret))
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

	config, err := kcconfig.NewConfig(k8sfake.NewSimpleClientset(configMap))
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
