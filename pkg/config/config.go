// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Config is populated from the cluster's Secret or ConfigMap and sets behavior of kapp-controller.
// NOTE because config may be populated from a Secret use caution if you're tempted to serialize.
type Config struct {
	client    kubernetes.Interface
	namespace string

	dataLock sync.RWMutex
	data     configData
}

// configData keeps all configuration data in one struct
// so that it's easy to swap all configuration atomically.
type configData struct {
	caCerts              string
	proxyOpts            ProxyOpts
	kappDeployRawOptions []string
	skipTLSVerify        string

	appDefaultSyncPeriod            time.Duration
	appMinimumSyncPeriod            time.Duration
	packageInstallDefaultSyncPeriod time.Duration
}

const (
	kcConfigName = "kapp-controller-config"
)

// NewConfig populates the Config struct from k8s resources.
// NewConfig prefers a secret named kcConfigName but
// if that does not exist falls back to a configMap of same name.
func NewConfig(client kubernetes.Interface) (*Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConf := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})

	namespace, _, err := kubeConf.Namespace()
	if err != nil {
		return nil, fmt.Errorf("Getting namespace: %s", err)
	}

	config := &Config{client: client, namespace: namespace}

	return config, config.Reload()
}

// Reload reloads configuration (proxy, CA certs, etc.) from ConfigMap/Secret.
// All configuration is cleared if ConfigMap or Secret are not found.
func (gc *Config) Reload() error {
	secret, configMap, err := gc.findExternalConfig()
	if err != nil {
		return err
	}
	switch {
	case secret != nil:
		return gc.addSecretDataToConfig(secret)
	case configMap != nil:
		return gc.addDataToConfig(configMap.Data)
	default:
		// Clear out previously set data
		return gc.addDataToConfig(map[string]string{})
	}
}

// findExternalConfig will populate exactly one of its return values and the others will be nil.
// we prefer to populate secret, fall back to configMap, and return unrecoverable errors if they occur.
func (gc *Config) findExternalConfig() (*v1.Secret, *v1.ConfigMap, error) {
	secret, err := gc.client.CoreV1().Secrets(
		gc.namespace).Get(context.Background(), kcConfigName, metav1.GetOptions{})
	// NOTE: to avoid nested ifs we are checking err == nil,  instead of != nil.
	if err == nil { // happy path return
		return secret, nil, nil
	}
	if !errors.IsNotFound(err) { // other error than NotFound so we're not gonna look for configMap
		return nil, nil, err
	}

	configMap, err := gc.client.CoreV1().ConfigMaps(
		gc.namespace).Get(context.Background(), kcConfigName, metav1.GetOptions{})
	if err == nil { // second happiest path return
		return nil, configMap, nil
	}
	if !errors.IsNotFound(err) { // other error than NotFound
		return nil, nil, err
	}

	// nothing found, no errors, return triple-nil
	return nil, nil, nil
}

// CACerts returns configured CA certificates in PEM format.
func (gc *Config) CACerts() string {
	gc.dataLock.RLock()
	defer gc.dataLock.RUnlock()

	return gc.data.caCerts
}

// ProxyOpts returns configured proxy configuration.
func (gc *Config) ProxyOpts() ProxyOpts {
	gc.dataLock.RLock()
	defer gc.dataLock.RUnlock()

	return gc.data.proxyOpts
}

// ShouldSkipTLSForAuthority compares a candidate host or host:port against a stored set of allow-listed authorities.
// the allow-list is built from the user-facing flag `dangerousSkipTLSVerify`.
// Note that in some cases the allow-list may contain ports, so the function name could also be ShouldSkipTLSForDomainAndPort
// Note that "authority" is defined in: https://www.rfc-editor.org/rfc/rfc3986#section-3 to mean "host and port"
func (gc *Config) ShouldSkipTLSForAuthority(candidateAuthority string) bool {
	gc.dataLock.RLock()
	defer gc.dataLock.RUnlock()

	authorities := gc.data.skipTLSVerify
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

// KappDeployRawOptions returns user configured kapp raw options
func (gc *Config) KappDeployRawOptions() []string {
	gc.dataLock.RLock()
	defer gc.dataLock.RUnlock()

	kappOptions := make([]string, 0)

	// Configure kapp to keep only 5 app changes as it seems that
	// larger number of ConfigMaps negative affects other controllers on the cluster.
	// Eventually kapp can be smart enough to keep minimal number of app changes.
	// Set default first so that it can be overridden by user provided options.
	// return append([]string{"--app-changes-max-to-keep=5"}, gc.data.kappDeployRawOptions...)
	kappOptions = append(kappOptions, "--app-changes-max-to-keep=5")
	kappOptions = append(kappOptions, "--apply-timeout=5m")
	kappOptions = append(kappOptions, "--diff-anchored=true")
	kappOptions = append(kappOptions, gc.data.kappDeployRawOptions...)
	return kappOptions
}

// AppDefaultSyncPeriod returns duration that is used by Apps
// that do not explicitly specify sync period.
func (gc *Config) AppDefaultSyncPeriod() time.Duration {
	const lowestDefault = 30 * time.Second
	if gc.data.appDefaultSyncPeriod > lowestDefault {
		return gc.data.appDefaultSyncPeriod
	}
	return lowestDefault
}

// AppMinimumSyncPeriod returns duration that is used as a lowest
// sync period App would use for reconciliation. This value
// takes precedence over any sync period that is lower.
func (gc *Config) AppMinimumSyncPeriod() time.Duration {
	const min = 30 * time.Second
	if gc.data.appMinimumSyncPeriod > min {
		return gc.data.appMinimumSyncPeriod
	}
	return min
}

// PackageInstallDefaultSyncPeriod returns duration that is used by Apps
// that do not explicitly specify sync period.
func (gc *Config) PackageInstallDefaultSyncPeriod() time.Duration {
	const defaultSyncPeriod = 10 * time.Minute
	const minDefaultSyncPeriod = 30 * time.Second
	if gc.data.packageInstallDefaultSyncPeriod != 0 {
		if gc.data.packageInstallDefaultSyncPeriod > minDefaultSyncPeriod {
			return gc.data.packageInstallDefaultSyncPeriod
		}
		return minDefaultSyncPeriod
	}
	return defaultSyncPeriod
}

func (gc *Config) addSecretDataToConfig(secret *v1.Secret) error {
	extractedValues := map[string]string{}
	for key, value := range secret.Data {
		extractedValues[key] = string(value)
	}
	return gc.addDataToConfig(extractedValues)
}

func (gc *Config) addDataToConfig(rawData map[string]string) error {
	data := configData{
		caCerts: rawData["caCerts"],
		proxyOpts: ProxyOpts{
			HTTPProxy:  rawData["httpProxy"],
			HTTPSProxy: rawData["httpsProxy"],
			NoProxy:    gc.replaceServiceHostPlaceholder(rawData["noProxy"]),
		},
		skipTLSVerify: rawData["dangerousSkipTLSVerify"],
	}

	if val := rawData["appDefaultSyncPeriod"]; len(val) > 0 {
		dur, err := time.ParseDuration(val)
		if err != nil {
			return fmt.Errorf("Unmarshaling appDefaultSyncPeriod as duration: %s", err)
		}
		data.appDefaultSyncPeriod = dur
	}

	if val := rawData["packageInstallDefaultSyncPeriod"]; len(val) > 0 {
		dur, err := time.ParseDuration(val)
		if err != nil {
			return fmt.Errorf("Unmarshaling packageInstallDefaultSyncPeriod as duration: %s", err)
		}
		data.packageInstallDefaultSyncPeriod = dur
	}

	if val := rawData["appMinimumSyncPeriod"]; len(val) > 0 {
		dur, err := time.ParseDuration(val)
		if err != nil {
			return fmt.Errorf("Unmarshaling appMinimumSyncPeriod as duration: %s", err)
		}
		data.appMinimumSyncPeriod = dur
	}

	if val := rawData["kappDeployRawOptions"]; len(val) > 0 {
		var opts []string
		err := json.Unmarshal([]byte(val), &opts)
		if err != nil {
			return fmt.Errorf("Unmarshaling kappDeployRawOptions as JSON: %s", err)
		}
		// Allowed flags will be verified before kapp is invoked within Kapp class.
		// (See pkg/deploy/kapp_restrict.go).
		data.kappDeployRawOptions = opts
	}

	gc.dataLock.Lock()
	defer gc.dataLock.Unlock()
	gc.data = data

	return nil
}

func (Config) replaceServiceHostPlaceholder(val string) string {
	const (
		kubernetesServiceHostEnvVar    = "KUBERNETES_SERVICE_HOST"
		kubernetesServiceHostShorthand = "KAPPCTRL_KUBERNETES_SERVICE_HOST"
	)
	if strings.Contains(val, kubernetesServiceHostShorthand) {
		k8sSvcHost := os.Getenv(kubernetesServiceHostEnvVar)
		val = strings.Replace(val, kubernetesServiceHostShorthand, k8sSvcHost, 1)
	}
	return val
}
