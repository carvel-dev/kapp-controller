// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"k8s.io/client-go/kubernetes"
)

type SkipTLSConfig interface {
	ShouldSkipTLSForAuthority(authority string) bool
}

// Factory allows to build various fetchers. At this point
// most of the fetching is performed via vendir.
type Factory struct {
	coreClient kubernetes.Interface
	vendirOpts VendirOpts
	cmdRunner  exec.CmdRunner
}

// NewFactory returns a Factory.
func NewFactory(coreClient kubernetes.Interface, vendirOpts VendirOpts, cmdRunner exec.CmdRunner) Factory {
	return Factory{coreClient, vendirOpts, cmdRunner}
}

func (f Factory) NewInline(opts v1alpha1.AppFetchInline, nsName string) *Inline {
	return NewInline(opts, nsName, f.coreClient)
}

// TODO: pass v1alpha1.Vendir opts here once api is exapnded
func (f Factory) NewVendir(nsName string) *Vendir {
	return NewVendir(nsName, f.coreClient, f.vendirOpts, f.cmdRunner)
}
