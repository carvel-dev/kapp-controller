// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Secrets gets cluster access based on a secret
type Secrets struct {
	coreClient kubernetes.Interface
}

// NewKubeconfigSecrets returns a Secrets
func NewKubeconfigSecrets(coreClient kubernetes.Interface) *Secrets {
	return &Secrets{coreClient}
}

// Find takes the location of the credentials secret and returns information to access the cluster
func (s *Secrets) Find(accessLocation AccessLocation,
	clusterOpts *v1alpha1.AppCluster) (AccessInfo, error) {

	if clusterOpts == nil {
		return AccessInfo{}, fmt.Errorf("Internal inconsistency: Expected cluster to not be nil")
	}

	if clusterOpts.KubeconfigSecretRef == nil {
		return AccessInfo{}, fmt.Errorf("Expected kubeconfig secret reference to be specified")
	}

	kubeconfigYAML, err := s.fetchKubeconfigYAML(accessLocation.Namespace, clusterOpts.KubeconfigSecretRef)
	if err != nil {
		return AccessInfo{}, err
	}

	kubeconfigRestricted, err := NewKubeconfigRestricted(kubeconfigYAML)
	if err != nil {
		return AccessInfo{}, err
	}

	pgoForCluster := AccessInfo{
		Name: accessLocation.Name,
		// Override destination namespace; if it's empty
		// assume kubeconfig contains preferred namespace
		Namespace: clusterOpts.Namespace,
		// Use provided namespace as app namespace
		DeployNamespace: clusterOpts.Namespace,
		Kubeconfig:      kubeconfigRestricted,
	}

	return pgoForCluster, nil
}

func (s *Secrets) fetchKubeconfigYAML(nsName string,
	secretRef *v1alpha1.AppClusterKubeconfigSecretRef) (string, error) {

	if len(nsName) == 0 {
		return "", fmt.Errorf("Internal inconsistency: Expected namespace name to not be empty")
	}
	if len(secretRef.Name) == 0 {
		return "", fmt.Errorf("Internal inconsistency: Expected service name to not be empty")
	}

	secret, err := s.coreClient.CoreV1().Secrets(nsName).Get(
		context.Background(), secretRef.Name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("Getting kubeconfig secret: %s", err)
	}

	key := secretRef.Key
	if len(key) == 0 {
		key = "value"
	}

	val, found := secret.Data[key]
	if !found {
		var otherKeys []string
		for otherKey := range secret.Data {
			otherKeys = append(otherKeys, otherKey)
		}

		return "", fmt.Errorf("Expected to find key '%s' in secret (keys: %s)",
			key, strings.Join(otherKeys, ", "))
	}

	return string(val), nil
}
