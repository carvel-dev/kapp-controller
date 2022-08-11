// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/internalpackaging/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Install(scheme *runtime.Scheme) {
	v1alpha1.AddToScheme(scheme)
}
