// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

// Kubeconfig provides access to the kubernetes cluster
// It initializes the service-account token cache
type Kubeconfig struct {
	kubeconfigSecrets *Secrets
	serviceAccounts   *ServiceAccounts
	log               logr.Logger
}

// AccessLocation contains the name/namespace of the resource which provides cluster access
// for example, a service account has a name and namespace
type AccessLocation struct {
	Name      string
	Namespace string
}

// AccessInfo provides a kubernetes kubeconfig for use to access the cluster
type AccessInfo struct {
	Name      string
	Namespace string

	Kubeconfig                    *Restricted
	DangerousUsePodServiceAccount bool
}

// NewKubeconfig creates a Kubeconfig with a new serviceaccount cache and kubeconfigsecrets
func NewKubeconfig(coreClient kubernetes.Interface, log logr.Logger) *Kubeconfig {
	return &Kubeconfig{
		kubeconfigSecrets: NewKubeconfigSecrets(coreClient),
		serviceAccounts:   NewServiceAccounts(coreClient, log),
		log:               log,
	}
}

// ClusterAccess takes cluster info and a ServiceAccount Name, and returns a populated kubeconfig that can connect to a cluster.
// if the saName is empty then you'll connect to a cluster via the clusterOpts inside the genericOpts, otherwise you'll use the specified SA.
func (k Kubeconfig) ClusterAccess(saName string, clusterOpts *v1alpha1.AppCluster, accessLocation AccessLocation) (AccessInfo, error) {
	var err error
	var clusterAccessInfo AccessInfo

	switch {
	case len(saName) > 0:
		clusterAccessInfo, err = k.serviceAccounts.Find(accessLocation, saName)
		if err != nil {
			return AccessInfo{}, err
		}

	case clusterOpts != nil:
		clusterAccessInfo, err = k.kubeconfigSecrets.Find(accessLocation, clusterOpts)
		if err != nil {
			return AccessInfo{}, err
		}

	default:
		return AccessInfo{}, fmt.Errorf("Expected service account or cluster specified")
	}
	return clusterAccessInfo, nil
}
