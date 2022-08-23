// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package clusterclient

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

// ClusterClient provides access to the kubernetes cluster
// It initializes the service-account token cache
type ClusterClient struct {
	coreClient        kubernetes.Interface
	kubeconfigSecrets *KubeconfigSecrets
	serviceAccounts   *ServiceAccounts
	log               logr.Logger
}

// CoreClient provides the underlying client for kubernetes
func (c *ClusterClient) CoreClient() kubernetes.Interface {
	return c.coreClient
}

// NewClusterClient creates a ClusterClient with a new serviceaccount cache and kubeconfigsecrets
func NewClusterClient(coreClient kubernetes.Interface, log logr.Logger) *ClusterClient {
	return &ClusterClient{
		coreClient:        coreClient,
		kubeconfigSecrets: NewKubeconfigSecrets(coreClient),
		serviceAccounts:   NewServiceAccounts(coreClient, log),
		log:               log,
	}
}

// ProcessOpts takes generic opts and a ServiceAccount Name, and returns a populated kubeconfig that can connect to a cluster.
// if the saName is empty then you'll connect to a cluster via the clusterOpts inside the genericOpts, otherwise you'll use the specified SA.
func (c ClusterClient) ProcessOpts(saName string, clusterOpts *v1alpha1.AppCluster, genericOpts GenericOpts) (ProcessedGenericOpts, error) {
	var err error
	var processedGenericOpts ProcessedGenericOpts

	switch {
	case len(saName) > 0:
		processedGenericOpts, err = c.serviceAccounts.Find(genericOpts, saName)
		if err != nil {
			return ProcessedGenericOpts{}, err
		}

	case clusterOpts != nil:
		processedGenericOpts, err = c.kubeconfigSecrets.Find(genericOpts, clusterOpts)
		if err != nil {
			return ProcessedGenericOpts{}, err
		}

	default:
		return ProcessedGenericOpts{}, fmt.Errorf("Expected service account or cluster specified")
	}
	return processedGenericOpts, nil
}
