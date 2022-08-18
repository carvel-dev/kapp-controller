// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package clusterstuff

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Manager struct {
	caches     map[string]interface{}
	log        logr.Logger
	coreClient kubernetes.Interface
}

// any function that takes no arguments but when called will get the cluster version for you (lazy)
type GetsVersion func() (*version.Info, error)

func GetClusterVersionLater (saName string, specCluster *v1alpha1.AppCluster, objMeta *metav1.ObjectMeta, log logr.Logger, coreClient kubernetes.Interface) GetsVersion {
	return func() (*version.Info, error){
		return GetClusterVersion(saName, specCluster, objMeta, log, coreClient)
	}
}

// GetClusterVersion returns the kubernetes API version for the cluster which has been supplied to kapp-controller via a kubeconfig
func (m *Manager) GetClusterVersion(saName string, specCluster *v1alpha1.AppCluster, objMeta *metav1.ObjectMeta) (*version.Info, error) {
	return GetClusterVersion(saName, specCluster, objMeta, m.log, m.coreClient)
}

// GetClusterVersion returns the kubernetes API version for the cluster which has been supplied to kapp-controller via a kubeconfig
func GetClusterVersion(saName string, specCluster *v1alpha1.AppCluster, objMeta *metav1.ObjectMeta, log logr.Logger, coreClient kubernetes.Interface) (*version.Info, error) {
	switch {
	case len(saName) > 0:
		return coreClient.Discovery().ServerVersion()
	case specCluster != nil:
		processedGenericOpts, err := deploy.ProcessOpts(saName, specCluster, deploy.GenericOpts{Name: objMeta.Name, Namespace: objMeta.Namespace}, coreClient, log)
		if err != nil {
			return nil, err
		}

		config, err := clientcmd.RESTConfigFromKubeConfig([]byte(processedGenericOpts.Kubeconfig.AsYAML()))
		if err != nil {
			return nil, err
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		return clientset.Discovery().ServerVersion()
	default:
		return nil, fmt.Errorf("Expected service account or cluster specified")
	}
}

// ProcessOpts takes generic opts and a ServiceAccount Name, and returns a populated kubeconfig that can connect to a cluster.
// if the saName is empty then you'll connect to a cluster via the clusterOpts inside the genericOpts, otherwise you'll use the specified SA.
func ProcessOpts(saName string, clusterOpts *v1alpha1.AppCluster, genericOpts deploy.GenericOpts, coreClient kubernetes.Interface, log logr.Logger) (deploy.ProcessedGenericOpts, error) {
	var err error
	var processedGenericOpts deploy.ProcessedGenericOpts

	switch {
	case len(saName) > 0:
		processedGenericOpts, err = NewServiceAccounts(coreClient, log).Find(genericOpts, saName)
		if err != nil {
			return ProcessedGenericOpts{}, err
		}

	case clusterOpts != nil:
		processedGenericOpts, err = NewKubeconfigSecrets(coreClient).Find(genericOpts, clusterOpts)
		if err != nil {
			return ProcessedGenericOpts{}, err
		}

	default:
		return ProcessedGenericOpts{}, fmt.Errorf("Expected service account or cluster specified")
	}
	return processedGenericOpts, nil
}
