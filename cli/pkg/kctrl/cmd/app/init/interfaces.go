// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Build interface {
	Save() error
	GetAppSpec() *v1alpha12.AppSpec
	SetAppSpec(*v1alpha12.AppSpec)
	GetObjectMeta() *metav1.ObjectMeta
	SetObjectMeta(*metav1.ObjectMeta)
	SetExport(export *[]Export)
	GetExport() *[]Export
}

type Step interface {
	PreInteract() error
	Interact() error
	PostInteract() error
}

func Run(step Step) error {
	err := step.PreInteract()
	if err != nil {
		return err
	}
	err = step.Interact()
	if err != nil {
		return err
	}
	return step.PostInteract()
}
