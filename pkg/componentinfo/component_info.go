// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

// Package componentinfo provides access to version and configuration information about components of the system.
package componentinfo

import (
	"context"
	"fmt"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"carvel.dev/kapp-controller/pkg/kubeconfig"
	"github.com/k14s/semver/v4"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ComponentInfo provides information about components of system
type ComponentInfo struct {
	coreClient            kubernetes.Interface
	clusterAccess         *kubeconfig.Kubeconfig
	kappControllerVersion string
}

// NewComponentInfo returns a ComponentInfo
func NewComponentInfo(coreClient kubernetes.Interface, clusterAccess *kubeconfig.Kubeconfig, kappControllerVersion string) *ComponentInfo {
	return &ComponentInfo{coreClient: coreClient, clusterAccess: clusterAccess, kappControllerVersion: kappControllerVersion}
}

// KappControllerVersion returns the running KC version
func (ci *ComponentInfo) KappControllerVersion() (semver.Version, error) {
	v, err := semver.ParseTolerant(ci.kappControllerVersion)
	if err != nil {
		return semver.Version{}, err
	}
	return v, nil
}

// KubernetesVersion returns the running K8s version depending on AppSpec
// If AppSpec points to external cluster, we use that k8s version instead
func (ci *ComponentInfo) KubernetesVersion(serviceAccountName string, specCluster *v1alpha1.AppCluster, objMeta *metav1.ObjectMeta) (semver.Version, error) {
	switch {
	case len(serviceAccountName) > 0:
		v, err := ci.coreClient.Discovery().ServerVersion()
		if err != nil {
			return semver.Version{}, err
		}
		return ci.parseAndScrubVersion(v.String())

	case specCluster != nil:
		accessInfo, err := ci.clusterAccess.ClusterAccess(serviceAccountName, specCluster, kubeconfig.AccessLocation{Name: objMeta.Name, Namespace: objMeta.Namespace}, "")
		if err != nil {
			return semver.Version{}, err
		}
		config, err := clientcmd.RESTConfigFromKubeConfig([]byte(accessInfo.Kubeconfig.AsYAML()))
		if err != nil {
			return semver.Version{}, err
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return semver.Version{}, err
		}

		v, err := clientset.Discovery().ServerVersion()
		if err != nil {
			return semver.Version{}, err
		}
		return ci.parseAndScrubVersion(v.String())

	default:
		return semver.Version{}, fmt.Errorf("Expected service account or cluster specified")
	}
}

// KubernetesAPIs returns the available kubernetes Group/Version resources
func (ci *ComponentInfo) KubernetesAPIs() ([]string, error) {
	groups, err := ci.coreClient.Discovery().ServerGroups()
	if err != nil {
		return []string{}, err
	}

	return metav1.ExtractGroupVersions(groups), nil
}

// NamespaceStatus returns the status of the App namespace
func (ci *ComponentInfo) NamespaceStatus(name string) (v1.NamespaceStatus, error) {
	namespace, err := ci.coreClient.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return v1.NamespaceStatus{}, err
	}
	return namespace.Status, nil
}

// parseAndScrubVersion parses version string and removes Pre and Build from the version
func (*ComponentInfo) parseAndScrubVersion(version string) (semver.Version, error) {
	retv, err := semver.ParseTolerant(version)
	if err != nil {
		return retv, err
	}
	retv.Pre = semver.PRVersion{}
	retv.Build = semver.BuildMeta{}
	return retv, nil
}
