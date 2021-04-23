// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

// +k8s:openapi-gen=true
type AppDeploy struct {
	Kapp *AppDeployKapp `json:"kapp,omitempty"`
}

// +k8s:openapi-gen=true
type AppDeployKapp struct {
	IntoNs     string   `json:"intoNs,omitempty"`
	MapNs      []string `json:"mapNs,omitempty"`
	RawOptions []string `json:"rawOptions,omitempty"`

	Inspect *AppDeployKappInspect `json:"inspect,omitempty"`
	Delete  *AppDeployKappDelete  `json:"delete,omitempty"`
}

// +k8s:openapi-gen=true
type AppDeployKappInspect struct {
	RawOptions []string `json:"rawOptions,omitempty"`
}

// +k8s:openapi-gen=true
type AppDeployKappDelete struct {
	RawOptions []string `json:"rawOptions,omitempty"`
}
