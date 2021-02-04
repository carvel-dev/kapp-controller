// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config

const (
	SecretK8sCorev1BasicAuthUsernameKey = "username"
	SecretK8sCorev1BasicAuthPasswordKey = "password"

	SecretK8sCoreV1SSHAuthPrivateKey = "ssh-privatekey"
	SecretSSHAuthKnownHosts          = "ssh-knownhosts" // not part of k8s

	SecretToken = "token"
)

// There structs have minimal used set of fields from their K8s representations.

type GenericMetadata struct {
	Name string
}

type Secret struct {
	APIVersion string
	Kind       string

	Metadata GenericMetadata
	Data     map[string][]byte
}

// nolint:golint
type ConfigMap struct {
	APIVersion string
	Kind       string

	Metadata GenericMetadata
	Data     map[string]string
}
