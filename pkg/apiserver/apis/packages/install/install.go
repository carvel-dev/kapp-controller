// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/packages"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/packages/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func Install(scheme *runtime.Scheme) {
	utilruntime.Must(packages.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}
