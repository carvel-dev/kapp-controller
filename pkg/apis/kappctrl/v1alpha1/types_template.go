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
}

type AppTemplateKbld struct{}

type AppTemplateSopsArgs struct {
	KMSKeys    []string `json:"kmsKeys,omitempty"`
	GCPKms     []string `json:"gcpKms,omitempty"`
	AzureKV    []string `json:"azureKV,omitempty"`
	PGP        []string `json:"pgp,omitempty"`
	AWSProifle string   `json:"awsProfile,omitempty"`
	IgnoreMac  bool     `json:"ignoreMac,omitempty"`
}

type AppTemplateSops struct {
	Match      string               `json:"match,omitempty"`
	MergeFiles bool                 `json:"mergeFiles,omitempty"`
	Args       *AppTemplateSopsArgs `json:"args,omitempty"`
}

type AppTemplateHelmTemplate struct {
	ValuesFrom []AppTemplateHelmTemplateValuesSource `json:"valuesFrom,omitempty"`
}

type AppTemplateHelmTemplateValuesSource struct {
	SecretRef    *AppTemplateHelmTemplateValuesSourceRef `json:"secretRef,omitempty"`
	ConfigMapRef *AppTemplateHelmTemplateValuesSourceRef `json:"configMapRef,omitempty"`
}

type AppTemplateHelmTemplateValuesSourceRef struct {
	corev1.LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
}

// TODO implement kustomize
type AppTemplateKustomize struct{}

// TODO implement jsonnet
type AppTemplateJsonnet struct{}
