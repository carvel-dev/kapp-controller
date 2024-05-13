// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging"
	"carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func Install(scheme *runtime.Scheme) {
	utilruntime.Must(datapackaging.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}
