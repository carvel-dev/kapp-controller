// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	kcConfigName = "kapp-controller-config"

	caCertsKey      = "caCerts"
	systemCertsFile = "/etc/pki/tls/certs/ca-bundle.crt"

	httpProxyKey     = "httpProxy"
	httpsProxyKey    = "httpsProxy"
	httpProxyEnvVar  = "http_proxy"
	httpsProxyEnvVar = "https_proxy"
	noProxyKey       = "noProxy"
	noProxyEnvVar    = "no_proxy"

	skipTLSVerifyKey = "dangerousSkipTLSVerify"
)

// Config is populated from the cluster's Secret or ConfigMap and sets behavior of kapp-controller.
// NOTE because config may be populated from a Secret use caution if you're tempted to serialize.
type Config struct {
	caCerts       string
	httpProxy     string
	httpsProxy    string
	noProxy       string
	skipTLSVerify string
	populated     bool
}

// findExternalConfig will populate exactly one of its return values and the others will be nil.
// we prefer to populate secret, fall back to configMap, and return unrecoverable errors if they occur.
func findExternalConfig(namespace string, client kubernetes.Interface) (*v1.Secret, *v1.ConfigMap, error) {
	secret, err := client.CoreV1().Secrets(namespace).Get(context.Background(), kcConfigName, metav1.GetOptions{})
	// NOTE: to avoid nested ifs we are checking err == nil,  instead of != nil.
	if err == nil { // happy path return
		return secret, nil, nil
	}
	if !errors.IsNotFound(err) { // other error than NotFound so we're not gonna look for configMap
		return nil, nil, err
	}

	configMap, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), kcConfigName, metav1.GetOptions{})
	if err == nil { // second happiest path return
		return nil, configMap, nil
	}
	if !errors.IsNotFound(err) { // other error than NotFound
		return nil, nil, err
	}

	// nothing found, no errors, return triple-nil
	return nil, nil, nil
}

// GetConfig populates the Config struct from k8s resources.
// GetConfig prefers a secret named kcConfigName but if that does not exist falls back to a configMap of same name.
func GetConfig(client kubernetes.Interface) (*Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConf := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	namespace, _, err := kubeConf.Namespace()
	if err != nil {
		return nil, fmt.Errorf("Getting namespace: %s", err)
	}

	secret, configMap, err := findExternalConfig(namespace, client)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	if secret != nil {
		config.addSecretDataToConfig(secret)
	} else if configMap != nil {
		config.addConfigMapDataToConfig(configMap)
	}

	return config, nil
}

func (gc *Config) Apply() error {
	if !gc.populated {
		return nil
	}

	err := gc.addTrustedCerts(gc.caCerts)
	if err != nil {
		return fmt.Errorf("Adding trusted certs: %s", err)
	}

	gc.configureProxies(gc.httpProxy, gc.httpsProxy, gc.noProxy)

	return nil
}

func (gc *Config) ShouldSkipTLSForDomain(candidateDomain string) bool {
	if !gc.populated {
		return false
	}

	domains := gc.skipTLSVerify
	if len(domains) == 0 {
		return false
	}

	for _, domain := range strings.Split(domains, ",") {
		// in case user gives domains in form "a, b"
		if strings.TrimSpace(domain) == candidateDomain {
			return true
		}
	}

	return false
}

func (gc *Config) addTrustedCerts(certChain string) (err error) {
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

func (gc *Config) configureProxies(httpProxy, httpsProxy, noProxy string) {
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

func (gc *Config) addSecretDataToConfig(secret *v1.Secret) {
	extractedValues := map[string]string{}
	for key, value := range secret.Data {
		extractedValues[key] = string(value)
	}

	gc.addDataToConfig(extractedValues)
}

func (gc *Config) addConfigMapDataToConfig(configMap *v1.ConfigMap) {
	gc.addDataToConfig(configMap.Data)
}

func (gc *Config) addDataToConfig(data map[string]string) {
	gc.caCerts = data[caCertsKey]
	gc.httpProxy = data[httpProxyKey]
	gc.httpsProxy = data[httpsProxyKey]
	gc.noProxy = data[noProxyKey]
	gc.skipTLSVerify = data[skipTLSVerifyKey]
	gc.populated = true
}
