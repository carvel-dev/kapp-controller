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
