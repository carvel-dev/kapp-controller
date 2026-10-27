// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig_test

import (
	"testing"

	"carvel.dev/kapp-controller/pkg/kubeconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientcmd "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/yaml"
)

func Test_NewKubeconfigRestricted_RemovesLocalAuthConfig(t *testing.T) {
	input := `
apiVersion: v1
kind: Config
current-context: ctx
clusters:
- name: cluster
  cluster:
    server: https://example.com
    certificate-authority: /etc/ca.crt
    certificate-authority-data: Y2EtZGF0YQ==
contexts:
- name: ctx
  context:
    cluster: cluster
    user: user
    namespace: ns
users:
- name: user
  user:
    token: token-value
    tokenFile: /var/run/token
    client-certificate: /etc/client.crt
    client-key: /etc/client.key
    client-certificate-data: Y2VydC1kYXRh
    client-key-data: a2V5LWRhdGE=
    exec:
      apiVersion: client.authentication.k8s.io/v1
      command: /bin/sh
      args: ["-c", "echo"]
    auth-provider:
      name: gcp
      config:
        cmd-path: /bin/sh
        cmd-args: -c echo
`
	restricted, err := kubeconfig.NewKubeconfigRestricted(input)
	require.NoError(t, err)

	var result clientcmd.Config
	require.NoError(t, yaml.Unmarshal([]byte(restricted.AsYAML()), &result))

	require.Len(t, result.AuthInfos, 1)
	authInfo := result.AuthInfos[0].AuthInfo
	assert.Nil(t, authInfo.AuthProvider)
	assert.Nil(t, authInfo.Exec)
	assert.Empty(t, authInfo.TokenFile)
	assert.Empty(t, authInfo.ClientCertificate)
	assert.Empty(t, authInfo.ClientKey)
	assert.Equal(t, "token-value", authInfo.Token)
	assert.Equal(t, []byte("cert-data"), authInfo.ClientCertificateData)
	assert.Equal(t, []byte("key-data"), authInfo.ClientKeyData)

	require.Len(t, result.Clusters, 1)
	assert.Empty(t, result.Clusters[0].Cluster.CertificateAuthority)
	assert.Equal(t, []byte("ca-data"), result.Clusters[0].Cluster.CertificateAuthorityData)
	assert.Equal(t, "https://example.com", result.Clusters[0].Cluster.Server)

	assert.NotContains(t, restricted.AsYAML(), "/bin/sh")
}
