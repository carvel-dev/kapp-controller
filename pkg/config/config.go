// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"fmt"
	"io"
	"net"
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

	caCertsKey            = "caCerts"
	systemCertsFile       = "/etc/pki/tls/certs/ca-bundle.crt"
	backupSystemCertsFile = "/etc/pki/tls/certs/ca-bundle.crt.orig"

	httpProxyKey     = "httpProxy"
	httpsProxyKey    = "httpsProxy"
	httpProxyEnvVar  = "http_proxy"
	httpsProxyEnvVar = "https_proxy"
	noProxyKey       = "noProxy"
	noProxyEnvVar    = "no_proxy"

	skipTLSVerifyKey = "dangerousSkipTLSVerify"

	kubernetesServiceHostEnvVar    = "KUBERNETES_SERVICE_HOST"
	kubernetesServiceHostShorthand = "KAPPCTRL_KUBERNETES_SERVICE_HOST"
)

// Config is populated from the cluster's Secret or ConfigMap and sets behavior of kapp-controller.
// NOTE because config may be populated from a Secret use caution if you're tempted to serialize.
type Config struct {
	caCerts            string
	httpProxy          string
	httpsProxy         string
	noProxy            string
	skipTLSVerify      string
	BackupCaBundlePath string
	SystemCaBundlePath string
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
	err := gc.addTrustedCerts(gc.caCerts)
	if err != nil {
		return fmt.Errorf("Adding trusted certs: %s", err)
	}

	gc.configureProxies()

	return nil
}

// ShouldSkipTLSForAuthority compares a candidate host or host:port against a stored set of allow-listed authorities.
// the allow-list is built from the user-facing flag `dangerousSkipTLSVerify`.
// Note that in some cases the allow-list may contain ports, so the function name could also be ShouldSkipTLSForDomainAndPort
// Note that "authority" is defined in: https://www.rfc-editor.org/rfc/rfc3986#section-3 to mean "host and port"
func (gc *Config) ShouldSkipTLSForAuthority(candidateAuthority string) bool {
	authorities := gc.skipTLSVerify
	if len(authorities) == 0 {
		return false
	}

	host, _, err := net.SplitHostPort(candidateAuthority)
	if err != nil {
		// SplitHostPort considers it to be an error if there's no port at all, but that's a common case for us.
		host = candidateAuthority
	}

	for _, spaceyAuthority := range strings.Split(authorities, ",") {
		// in case user gives domains in form "a, b"
		authority := strings.TrimSpace(spaceyAuthority)
		// check if the host matches the allowed authority, meaning all ports for that host are allowed
		if authority == host {
			return true
		}
		// check the full candidate string in case they both have a port and its the same port
		if authority == candidateAuthority {
			return true
		}
	}

	return false
}

func (gc *Config) addTrustedCerts(certChain string) (err error) {
	backupCertsFilePath := backupSystemCertsFile
	systemCertsFilePath := systemCertsFile

	if gc.BackupCaBundlePath != "" && gc.SystemCaBundlePath != "" {
		backupCertsFilePath = gc.BackupCaBundlePath
		systemCertsFilePath = gc.SystemCaBundlePath
	}

	backupFile, err := os.Open(backupCertsFilePath)
	if err != nil {
		return fmt.Errorf("Opening original certs file: %s", err)
	}
	defer backupFile.Close()

	tmpFile, err := os.CreateTemp(os.TempDir(), "tmp-ca-bundle-")
	if err != nil {
		return fmt.Errorf("Creating tmp certs file: %s", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, backupFile)
	if err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("Copying certs file: %s", err)
	}

	_, err = tmpFile.Write([]byte("\n" + certChain))
	if err != nil {
		_ = tmpFile.Close()
		return err
	}

	if err = tmpFile.Close(); err != nil {
		return err
	}

	err = os.Rename(tmpFile.Name(), systemCertsFilePath)
	if err != nil {
		return fmt.Errorf("Renaming certs file: %s", err)
	}

	return nil
}

func (gc *Config) configureProxies() {
	if httpProxyEnvVar == "" {
		os.Unsetenv(httpProxyEnvVar)
		os.Unsetenv(strings.ToUpper(httpProxyEnvVar))
	} else {
		os.Setenv(httpProxyEnvVar, gc.httpProxy)
		os.Setenv(strings.ToUpper(httpProxyEnvVar), gc.httpProxy)
	}

	if httpsProxyEnvVar == "" {
		os.Unsetenv(httpsProxyEnvVar)
		os.Unsetenv(strings.ToUpper(httpsProxyEnvVar))
	} else {
		os.Setenv(httpsProxyEnvVar, gc.httpsProxy)
		os.Setenv(strings.ToUpper(httpsProxyEnvVar), gc.httpsProxy)
	}

	if noProxyEnvVar == "" {
		os.Unsetenv(noProxyEnvVar)
		os.Unsetenv(strings.ToUpper(noProxyEnvVar))
	} else {
		gc.addKubernetesServiceHostInNoProxy()
		os.Setenv(noProxyEnvVar, gc.noProxy)
		os.Setenv(strings.ToUpper(noProxyEnvVar), gc.noProxy)
	}
}

func (gc *Config) addKubernetesServiceHostInNoProxy() {
	if strings.Contains(gc.noProxy, kubernetesServiceHostShorthand) {
		k8sSvcHost := os.Getenv(kubernetesServiceHostEnvVar)
		gc.noProxy = strings.Replace(gc.noProxy, kubernetesServiceHostShorthand, k8sSvcHost, 1)
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
}
