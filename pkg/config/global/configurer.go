package global

import (
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	configMapName = "kapp-controller-config"

	caCertsKey      = "caCerts"
	systemCertsFile = "/etc/ssl/certs/ca-certificates.crt"

	httpProxyKey     = "httpProxy"
	httpsProxyKey    = "httpsProxy"
	httpProxyEnvVar  = "http_proxy"
	httpsProxyEnvVar = "https_proxy"
	noProxyKey       = "noProxy"
	noProxyEnvVar    = "no_proxy"
)

type GlobalConfigurer struct {
	namespace string
	client    kubernetes.Interface
}

func NewGlobalConfigurer() (*GlobalConfigurer, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConf := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	namespace, _, err := kubeConf.Namespace()
	if err != nil {
		return nil, fmt.Errorf("Getting namespace: %s", err)
	}

	restConfig := config.GetConfigOrDie()
	coreClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("Buidling client: %s", err)
	}

	return &GlobalConfigurer{
		namespace: namespace,
		client:    coreClient,
	}, nil
}

func (gc *GlobalConfigurer) Configure() error {
	configMap, err := gc.configMap()
	if err != nil {
		return fmt.Errorf("Fetching config map: %s", err)
	}

	if configMap == nil {
		return nil
	}

	err = gc.addTrustedCerts(configMap.Data[caCertsKey])
	if err != nil {
		return fmt.Errorf("Adding trusted certs: %s", err)
	}

	gc.configureProxies(configMap.Data[httpProxyKey], configMap.Data[httpsProxyKey], configMap.Data[noProxyKey])

	return nil
}

func (gc *GlobalConfigurer) addTrustedCerts(certChain string) (err error) {
	if certChain == "" {
		return nil
	}

	var file *os.File
	file, err = os.OpenFile(systemCertsFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("Opening certs file: %s", err)
	}

	_, err = file.Write([]byte("\n" + certChain))
	if err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}

func (gc *GlobalConfigurer) configureProxies(httpProxy, httpsProxy, noProxy string) {
	if httpProxy != "" {
		os.Setenv(httpProxyEnvVar, httpProxy)
		os.Setenv(strings.ToUpper(httpProxyEnvVar), httpProxy)
	}

	if httpsProxy != "" {
		os.Setenv(httpsProxyEnvVar, httpsProxy)
		os.Setenv(strings.ToUpper(httpsProxyEnvVar), httpsProxy)
	}

	if noProxy != "" {
		os.Setenv(noProxyEnvVar, noProxy)
		os.Setenv(strings.ToUpper(noProxyEnvVar), noProxy)
	}
}

func (gc *GlobalConfigurer) configMap() (*v1.ConfigMap, error) {
	configMap, err := gc.client.CoreV1().ConfigMaps(gc.namespace).Get(configMapName, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		return nil, nil
	}

	return configMap, err
}
