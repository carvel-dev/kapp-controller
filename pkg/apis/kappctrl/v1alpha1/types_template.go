// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type AppTemplate struct {
	Ytt          *AppTemplateYtt          `json:"ytt,omitempty"`
	Kbld         *AppTemplateKbld         `json:"kbld,omitempty"`
	HelmTemplate *AppTemplateHelmTemplate `json:"helmTemplate,omitempty"`
	Kustomize    *AppTemplateKustomize    `json:"kustomize,omitempty"`
	Jsonnet      *AppTemplateJsonnet      `json:"jsonnet,omitempty"`
	Sops         *AppTemplateSops         `json:"sops,omitempty"`
}

type AppTemplateYtt struct {
	IgnoreUnknownComments bool            `json:"ignoreUnknownComments,omitempty"`
	Strict                bool            `json:"strict,omitempty"`
	Inline                *AppFetchInline `json:"inline,omitempty"`
	Paths                 []string        `json:"paths,omitempty"`
	FileMarks             []string        `json:"fileMarks,omitempty"`
}

type AppTemplateKbld struct {
	Paths []string `json:"paths,omitempty"`
}

type AppTemplateHelmTemplate struct {
	Name       string                                `json:"name,omitempty"`
	Namespace  string                                `json:"namespace,omitempty"`
	Path       string                                `json:"path,omitempty"`
	ValuesFrom []AppTemplateHelmTemplateValuesSource `json:"valuesFrom,omitempty"`
}

type AppTemplateHelmTemplateValuesSource struct {
	SecretRef    *AppTemplateHelmTemplateValuesSourceRef `json:"secretRef,omitempty"`
	ConfigMapRef *AppTemplateHelmTemplateValuesSourceRef `json:"configMapRef,omitempty"`
	Path         string                                  `json:"path,omitempty"`
}

type AppTemplateHelmTemplateValuesSourceRef struct {
	corev1.LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
}

// TODO implement kustomize
type AppTemplateKustomize struct{}

// TODO implement jsonnet
type AppTemplateJsonnet struct{}

type AppTemplateSops struct {
	PGP   *AppTemplateSopsPGP `json:"pgp,omitempty"`
	Paths []string            `json:"paths,omitempty"`
}

type AppTemplateSopsPGP struct {
	PrivateKeysSecretRef *AppTemplateSopsPGPPrivateKeysSecretRef `json:"privateKeysSecretRef,omitempty"`
}

type AppTemplateSopsPGPPrivateKeysSecretRef struct {
	Name string `json:"name,omitempty"`
}
