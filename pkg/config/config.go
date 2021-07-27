// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"os"
	"strings"

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

// TODO: Use struct for keys
type Config struct {
	configMap *v1.ConfigMap
}

func GetConfig(client kubernetes.Interface) (*Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConf := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	namespace, _, err := kubeConf.Namespace()
	if err != nil {
		return nil, fmt.Errorf("Getting namespace: %s", err)
	}

	// TODO: Try to simplify nesting
	configMap := &v1.ConfigMap{}
	secret, err := client.CoreV1().Secrets(namespace).Get(context.Background(), kcConfigName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			configMap, err = client.CoreV1().ConfigMaps(namespace).Get(context.Background(), kcConfigName, metav1.GetOptions{})
			if errors.IsNotFound(err) {
				return &Config{}, nil
			}
			return &Config{}, nil
		}
		return nil, err
	}

	if secret != nil {
		err := addSecretDataToConfigMap(configMap, secret)
		if err != nil {
			return nil, err
		}
	}

	return &Config{configMap}, nil
}

func (gc *Config) Apply() error {
	if gc.configMap == nil {
		return nil
	}

	configMap := gc.configMap
	err := gc.addTrustedCerts(gc.configMap.Data[caCertsKey])
	if err != nil {
		return fmt.Errorf("Adding trusted certs: %s", err)
	}

	gc.configureProxies(configMap.Data[httpProxyKey], configMap.Data[httpsProxyKey], configMap.Data[noProxyKey])

	return nil
}

func (gc *Config) ShouldSkipTLSForDomain(candidateDomain string) bool {
	if gc.configMap == nil {
		return false
	}

	domains, exists := gc.configMap.Data[skipTLSVerifyKey]
	if !exists {
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

// TODO: No base64 decoding needed
func addSecretDataToConfigMap(configMap *v1.ConfigMap, secret *v1.Secret) error {
	keys := []string{caCertsKey, httpProxyKey, httpsProxyKey, noProxyKey, skipTLSVerifyKey}
	configMap.Data = map[string]string{}
	for _, key := range keys {
		if value, valueExists := secret.Data[key]; valueExists {
			decodedValue, err := base64.StdEncoding.DecodeString(string(value))
			if err != nil {
				return err
			}
			configMap.Data[key] = string(decodedValue)
		}
	}

	return nil
}
