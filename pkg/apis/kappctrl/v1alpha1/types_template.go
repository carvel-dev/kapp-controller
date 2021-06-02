// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

// +k8s:openapi-gen=true
type AppTemplate struct {
	Ytt          *AppTemplateYtt          `json:"ytt,omitempty"`
	Kbld         *AppTemplateKbld         `json:"kbld,omitempty"`
	HelmTemplate *AppTemplateHelmTemplate `json:"helmTemplate,omitempty"`
	Kustomize    *AppTemplateKustomize    `json:"kustomize,omitempty"`
	Jsonnet      *AppTemplateJsonnet      `json:"jsonnet,omitempty"`
	Sops         *AppTemplateSops         `json:"sops,omitempty"`
}

// +k8s:openapi-gen=true
type AppTemplateYtt struct {
	IgnoreUnknownComments bool                      `json:"ignoreUnknownComments,omitempty"`
	Strict                bool                      `json:"strict,omitempty"`
	Inline                *AppFetchInline           `json:"inline,omitempty"`
	Paths                 []string                  `json:"paths,omitempty"`
	FileMarks             []string                  `json:"fileMarks,omitempty"`
	ValuesFrom            []AppTemplateValuesSource `json:"valuesFrom,omitempty"`
}

// +k8s:openapi-gen=true
type AppTemplateKbld struct {
	Paths []string `json:"paths,omitempty"`
}

// +k8s:openapi-gen=true
type AppTemplateHelmTemplate struct {
	Name       string                    `json:"name,omitempty"`
	Namespace  string                    `json:"namespace,omitempty"`
	Path       string                    `json:"path,omitempty"`
	ValuesFrom []AppTemplateValuesSource `json:"valuesFrom,omitempty"`
}

// +k8s:openapi-gen=true
type AppTemplateValuesSource struct {
	SecretRef    *AppTemplateValuesSourceRef `json:"secretRef,omitempty"`
	ConfigMapRef *AppTemplateValuesSourceRef `json:"configMapRef,omitempty"`
	Path         string                      `json:"path,omitempty"`
}

// +k8s:openapi-gen=true
type AppTemplateValuesSourceRef struct {
	corev1.LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
}

// TODO implement kustomize
// +k8s:openapi-gen=true
type AppTemplateKustomize struct{}

// TODO implement jsonnet
// +k8s:openapi-gen=true
type AppTemplateJsonnet struct{}

// +k8s:openapi-gen=true
type AppTemplateSops struct {
	PGP   *AppTemplateSopsPGP `json:"pgp,omitempty"`
	Paths []string            `json:"paths,omitempty"`
}

// +k8s:openapi-gen=true
type AppTemplateSopsPGP struct {
	PrivateKeysSecretRef *AppTemplateSopsPGPPrivateKeysSecretRef `json:"privateKeysSecretRef,omitempty"`
}

// +k8s:openapi-gen=true
type AppTemplateSopsPGPPrivateKeysSecretRef struct {
	Name string `json:"name,omitempty"`
}
