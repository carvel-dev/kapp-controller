// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	installinstpkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/install"
	installkctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/install"
	installpkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/package/install"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	Scheme = scheme.Scheme
)

func init() {
	installpkg.Install(Scheme)
	installkctrl.Install(Scheme)
	installinstpkg.Install(Scheme)
}
