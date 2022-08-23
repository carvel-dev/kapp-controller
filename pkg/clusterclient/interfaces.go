// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package clusterclient

// GenericOpts contains the basic information for discovery of a kubernetes workload
type GenericOpts struct {
	Name      string
	Namespace string
}

// ProcessedGenericOpts provides a kubeconfig for access to the cluster, and the ability to toggle usage of a privileged
// pod service account
type ProcessedGenericOpts struct {
	Name      string
	Namespace string

	Kubeconfig                    *KubeconfigRestricted
	DangerousUsePodServiceAccount bool
}
