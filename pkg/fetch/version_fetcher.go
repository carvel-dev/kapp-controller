// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"fmt"

	"github.com/k14s/semver/v4"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/clusterclient"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// VersionFetcher supplies a generic interface for retrieving versions for components of the system
type VersionFetcher interface {
	GetKappControllerVersion() string
	GetKubernetesVersion() semver.Version
}

// VersionFetch allows retrieving versions for components of the system
type VersionFetch struct {
	clusterClient     *clusterclient.ClusterClient
	controllerVersion string
}

// NewVersionFetcher creates a VersionFetch object based on the clusterClient (for kubernetes access, and the static controller version)
func NewVersionFetcher(clusterclient *clusterclient.ClusterClient, controllerVersion string) *VersionFetch {
	return &VersionFetch{
		clusterClient:     clusterclient,
		controllerVersion: controllerVersion,
	}
}

// GetKappControllerVersion returns the live kapp-controller version
func (vf VersionFetch) GetKappControllerVersion() string {
	return vf.controllerVersion
}

// GetKubernetesVersion returns the live kubernetes version (local or external cluster based on the spec)
func (vf VersionFetch) GetKubernetesVersion(saName string, specCluster *v1alpha1.AppCluster, genericOpts clusterclient.GenericOpts) (semver.Version, error) {
	switch {
	case len(saName) > 0:
		version, err := vf.clusterClient.CoreClient().Discovery().ServerVersion()
		if err != nil {
			return semver.Version{}, err
		}

		return semver.ParseTolerant(version.GitVersion)
	case specCluster != nil:
		processedGenericOpts, err := vf.clusterClient.ProcessOpts(saName, specCluster, clusterclient.GenericOpts{Name: genericOpts.Name, Namespace: genericOpts.Namespace})
		if err != nil {
			return semver.Version{}, err
		}

		config, err := clientcmd.RESTConfigFromKubeConfig([]byte(processedGenericOpts.Kubeconfig.AsYAML()))
		if err != nil {
			return semver.Version{}, err
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return semver.Version{}, err
		}

		version, err := clientset.Discovery().ServerVersion()
		if err != nil {
			return semver.Version{}, err
		}

		return semver.ParseTolerant(version.GitVersion)
	default:
		return semver.Version{}, fmt.Errorf("Expected service account or clusterclient specified")
	}
}
