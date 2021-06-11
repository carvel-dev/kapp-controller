// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeconfigSecrets struct {
	coreClient kubernetes.Interface
}

func NewKubeconfigSecrets(coreClient kubernetes.Interface) *KubeconfigSecrets {
	return &KubeconfigSecrets{coreClient}
}

func (s *KubeconfigSecrets) Find(genericOpts GenericOpts,
	clusterOpts *v1alpha1.AppCluster) (ProcessedGenericOpts, error) {

	if clusterOpts == nil {
		return ProcessedGenericOpts{}, fmt.Errorf("Internal inconsistency: Expected cluster to not be nil")
	}

	if clusterOpts.KubeconfigSecretRef == nil {
		return ProcessedGenericOpts{}, fmt.Errorf("Expected kubeconfig secret reference to be specified")
	}

	kubeconfigYAML, err := s.fetchKubeconfigYAML(genericOpts.Namespace, clusterOpts.KubeconfigSecretRef)
	if err != nil {
		return ProcessedGenericOpts{}, err
	}

	kubeconfigRestricted, err := NewKubeconfigRestricted(kubeconfigYAML)
	if err != nil {
		return ProcessedGenericOpts{}, err
	}

	pgoForCluster := ProcessedGenericOpts{
		Name: genericOpts.Name,
		// Override destination namespace; if it's empty
		// assume kubeconfig contains preferred namespace
		Namespace:  clusterOpts.Namespace,
		Kubeconfig: kubeconfigRestricted,
	}

	return pgoForCluster, nil
}

func (s *KubeconfigSecrets) fetchKubeconfigYAML(nsName string,
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
