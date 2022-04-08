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

	err = stubTrustedCerts(t, config)
	assert.NoError(t, err)

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

	err = stubTrustedCerts(t, config)
	assert.NoError(t, err)

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

	err = stubTrustedCerts(t, config)
	assert.NoError(t, err)

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

	err = stubTrustedCerts(t, config)
	assert.NoError(t, err)

	assert.Nil(t, config.Apply(), "unexpected error after running config.Apply()", err)

	expected := "10.96.0.1"
	noProxyActual := os.Getenv("no_proxy")

	assert.Equal(t, expected, noProxyActual)
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
	k8scs := k8sfake.NewSimpleClientset(configMap)
	config, err := kcconfig.GetConfig(k8scs)
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

func Test_TrustedCertsCreateConfig(t *testing.T) {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"caCerts": "cert-42",
		},
	}
	k8scs := k8sfake.NewSimpleClientset(configMap)
	config, err := kcconfig.GetConfig(k8scs)
	assert.NoError(t, err)

	backup, certs, close, err := createCertTempFiles(t)
	assert.NoError(t, err)
	defer close()

	config.BackupCaBundlePath = backup.Name()
	config.SystemCaBundlePath = certs.Name()

	assert.NoError(t, config.Apply(), "unexpected error after running config.Apply()")

	contents, err := os.ReadFile(config.SystemCaBundlePath)
	assert.NoError(t, err)

	assert.Contains(t, string(contents), "cert-42")
}

func Test_TrustedCertsUpdateConfig(t *testing.T) {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"caCerts": "cert-42",
		},
	}
	k8scs := k8sfake.NewSimpleClientset(configMap)
	config, err := kcconfig.GetConfig(k8scs)
	assert.NoError(t, err)

	backup, certs, close, err := createCertTempFiles(t)
	assert.NoError(t, err)
	defer close()

	config.BackupCaBundlePath = backup.Name()
	config.SystemCaBundlePath = certs.Name()

	assert.NoError(t, config.Apply(), "unexpected error after running config.Apply()")

	contents, err := os.ReadFile(config.SystemCaBundlePath)
	assert.NoError(t, err)
	assert.Contains(t, string(contents), "cert-42")

	// update config
	configMap.Data["caCerts"] = "cert-43"

	k8scs = k8sfake.NewSimpleClientset(configMap)
	config, err = kcconfig.GetConfig(k8scs)
	assert.NoError(t, err)

	config.BackupCaBundlePath = backup.Name()
	config.SystemCaBundlePath = certs.Name()

	assert.NoError(t, config.Apply(), "unexpected error after running config.Apply()")

	contents, err = os.ReadFile(config.SystemCaBundlePath)
	assert.NoError(t, err)

	assert.Contains(t, string(contents), "cert-43")
}

func Test_TrustedCertsDeleteConfig(t *testing.T) {
	backup, certs, close, err := createCertTempFiles(t)
	assert.NoError(t, err)
	defer close()

	backup.Write([]byte("this-is-the-old-content"))

	// no config found
	k8scs := k8sfake.NewSimpleClientset()
	config, err := kcconfig.GetConfig(k8scs)
	assert.NoError(t, err)

	config.BackupCaBundlePath = backup.Name()
	config.SystemCaBundlePath = certs.Name()

	assert.NoError(t, config.Apply(), "unexpected error after running config.Apply()")

	contents, err := os.ReadFile(config.SystemCaBundlePath)
	assert.NoError(t, err)

	// restored to the backup without any additional data
	assert.Contains(t, string(contents), "this-is-the-old-content")
}

func stubTrustedCerts(t *testing.T, gc *kcconfig.Config) error {
	backup, certs, close, err := createCertTempFiles(t)
	if err != nil {
		return err
	}
	defer close()

	gc.BackupCaBundlePath = backup.Name()
	gc.SystemCaBundlePath = certs.Name()

	return nil
}

func createCertTempFiles(t *testing.T) (backup *os.File, certs *os.File, close func(), err error) {
	backup, err = os.CreateTemp("", "backup.crt")
	if err != nil {
		return nil, nil, nil, err
	}

	certs, err = os.CreateTemp("", "certs.crt")
	if err != nil {
		return nil, nil, nil, err
	}

	return backup, certs, func() {
		backup.Close()
		certs.Close()
	}, nil
}
