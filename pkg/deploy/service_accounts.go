// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/token"
	authenticationv1 "k8s.io/api/authentication/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	caCertPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
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

	tokenMgr := token.NewManager(s.coreClient)

	t, err := tokenMgr.GetServiceAccountToken(nsName, saName, &authenticationv1.TokenRequest{})
	if err != nil {
		return "", fmt.Errorf("Failed to get service account token: %s", err)
	}

	cert, err := os.ReadFile(caCertPath)
	if err != nil {
		return "", fmt.Errorf("Failed to read ca cert from %s: %s", caCertPath, err)
	}

	return s.buildKubeconfig(t.Status.Token, nsName, cert)
}

func (s *ServiceAccounts) buildKubeconfig(token string, nsBytes string, caCert []byte) (string, error) {
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

	caB64Encoded := base64.StdEncoding.EncodeToString(caCert)

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
