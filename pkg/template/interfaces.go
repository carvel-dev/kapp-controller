// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"io"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	stdinPath = "-"
)

type Template interface {
	// TemplateDir works on directory returning templating result,
	// and boolean indicating whether subsequent operations
	// should operate on result, or continue operating on the directory
	TemplateDir(dirPath string) (exec.CmdRunResult, bool)
	// TemplateStream works on a stream returning templating result.
	// dirPath is provided for context from which to reference additional inputs.
	TemplateStream(stream io.Reader, dirPath string) exec.CmdRunResult
}

// GenericOpts carries App information used across API boundaries.
// Primarily used in a context when templating with values
type GenericOpts struct {
	Name      string
	Namespace string
	Metadata  *PartialObjectMetadata
}

// PartialObjectMetadata represents an v1alpha1.App with a subset of Metadata fields exposed.
// Used to control which metadata fields an operator can query (using jsonpath) to provide as a Value when templating
type PartialObjectMetadata struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

// ObjectMeta represents a subset of v1.ObjectMetadata fields
type ObjectMeta struct {
	Name        string            `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Namespace   string            `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
	UID         types.UID         `json:"uid,omitempty" protobuf:"bytes,3,opt,name=uid,casttype=k8s.io/kubernetes/pkg/types.UID"`
	Labels      map[string]string `json:"labels,omitempty" protobuf:"bytes,4,rep,name=labels"`
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,5,rep,name=annotations"`
}
