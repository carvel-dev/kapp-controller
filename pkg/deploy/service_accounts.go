// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"context"
	"encoding/base64"
	"fmt"

	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceAccounts struct {
	coreClient kubernetes.Interface
}

func NewServiceAccounts(coreClient kubernetes.Interface) *ServiceAccounts {
	return &ServiceAccounts{coreClient}
}

func (s *ServiceAccounts) Find(genericOpts GenericOpts, saName string) (ProcessedGenericOpts, error) {
	kubeconfigYAML, err := s.fetchServiceAccount(genericOpts.Namespace, saName)
	if err != nil {
		return ProcessedGenericOpts{}, err
	}

	kubeconfigRestricted, err := NewKubeconfigRestricted(kubeconfigYAML)
	if err != nil {
		return ProcessedGenericOpts{}, err
	}

	pgoForSA := ProcessedGenericOpts{
		Name:       genericOpts.Name,
		Namespace:  "", // Assume kubeconfig contains preferred namespace from SA
		Kubeconfig: kubeconfigRestricted,
	}

	return pgoForSA, nil
}

func (s *ServiceAccounts) fetchServiceAccount(nsName string, saName string) (string, error) {
	if len(nsName) == 0 {
		return "", fmt.Errorf("Internal inconsistency: Expected namespace name to not be empty")
	}
	if len(saName) == 0 {
		return "", fmt.Errorf("Internal inconsistency: Expected service account name to not be empty")
	}

	// sa, err := s.coreClient.CoreV1().ServiceAccounts(nsName).Get(context.Background(), saName, metav1.GetOptions{})
	// if err != nil {
	// 	return "", fmt.Errorf("Getting service account: %s", err)
	// }

	t, err := s.coreClient.CoreV1().ServiceAccounts(nsName).CreateToken(context.Background(), saName, &authenticationv1.TokenRequest{}, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create token: %s", err)
	}

	c, err := s.coreClient.CoreV1().ConfigMaps(nsName).Get(context.Background(), "kube-root-ca.crt", metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get 'kube-root-ca.crt' configmap: %s", err)
	}

	return s.buildKubeconfig(t.Status.Token, nsName, c.Data["ca.crt"])

	// for _, secretRef := range sa.Secrets {
	// 	secret, err := s.coreClient.CoreV1().Secrets(nsName).Get(context.Background(), secretRef.Name, metav1.GetOptions{})
	// 	if err != nil {
	// 		return "", fmt.Errorf("Getting service account secret: %s", err)
	// 	}

	// 	if secret.Type != corev1.SecretTypeServiceAccountToken {
	// 		continue
	// 	}

	// 	return s.buildKubeconfig(secret, []byte(t.Status.Token))
	// }

	// return "", fmt.Errorf("Expected to find one service account token secret, but found none")
}

func (s *ServiceAccounts) buildKubeconfig(token string, nsBytes string, caCert string) (string, error) {
	// caBytes, found := secret.Data[corev1.ServiceAccountRootCAKey]
	// if !found {
	// 	return "", fmt.Errorf("Expected to find service account token ca")
	// }

	// tokenBytes, found := secret.Data[corev1.ServiceAccountTokenKey]
	// if !found {
	// 	return "", fmt.Errorf("Expected to find service account token value")
	// }

	// nsBytes, found := secret.Data[corev1.ServiceAccountNamespaceKey]
	// if !found {
	// 	return "", fmt.Errorf("Expected to find service account token namespace")
	// }

	const kubeconfigYAMLTpl = `
apiVersion: v1
kind: Config
clusters:
- name: dst-cluster
  cluster:
    certificate-authority-data: "%s"
    server: https://${KAPP_KUBERNETES_SERVICE_HOST_PORT}
users:
- name: dst-user
  user:
    token: "%s"
contexts:
- name: dst-ctx
  context:
    cluster: dst-cluster
    namespace: "%s"
    user: dst-user
current-context: dst-ctx
`

	caB64Encoded := base64.StdEncoding.EncodeToString([]byte(caCert))

	return fmt.Sprintf(kubeconfigYAMLTpl, caB64Encoded, []byte(token), []byte(nsBytes)), nil
}

/*

Example SA + secret:

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: app1-sa
  namespace: app1
secrets:
- name: app1-sa-token-grr7z
---
apiVersion: v1
kind: Secret
metadata:
  name: app1-sa-token-grr7z
  namespace: app1
  annotations:
    kubernetes.io/service-account.name: app1-sa
    kubernetes.io/service-account.uid: 26675b19-769a-4145-a386-7ca2b3ab3435
type: kubernetes.io/service-account-token
data:
  ca.crt: LS0tLS...
  namespace: a2FwcC1jb250cm9sbGVy
  token: ZXlKaGJ...

*/
