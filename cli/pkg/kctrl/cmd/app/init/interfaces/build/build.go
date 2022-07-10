// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Build interface {
	Save() error
	GetAppSpec() *v1alpha12.AppSpec
	SetAppSpec(*v1alpha12.AppSpec)
	GetObjectMeta() *metav1.ObjectMeta
	SetObjectMeta(*metav1.ObjectMeta)
	SetExport(export *[]appbuild.Export)
	GetExport() *[]appbuild.Export
}
