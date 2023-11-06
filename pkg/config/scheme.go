// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	internalpkging "carvel.dev/kapp-controller/pkg/apis/internalpackaging/install"
	installkctrl "carvel.dev/kapp-controller/pkg/apis/kappctrl/install"
	pkging "carvel.dev/kapp-controller/pkg/apis/packaging/install"
	datapackaging "carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging/install"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	Scheme = scheme.Scheme
)

func init() {
	pkging.Install(Scheme)
	internalpkging.Install(Scheme)
	installkctrl.Install(Scheme)
	datapackaging.Install(Scheme)
}
