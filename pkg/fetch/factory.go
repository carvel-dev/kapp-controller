// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/clusterclient"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

type SkipTLSConfig interface {
	ShouldSkipTLSForAuthority(authority string) bool
}

// Factory allows to build various fetchers. At this point
// most of the fetching is performed via vendir.
type Factory struct {
	clusterClient     *clusterclient.ClusterClient
	vendirOpts        VendirOpts
	cmdRunner         exec.CmdRunner
	controllerVersion string
}

// NewFactory returns a Factory.
func NewFactory(clusterClient *clusterclient.ClusterClient, vendirOpts VendirOpts, cmdRunner exec.CmdRunner, controllerVersion string) Factory {
	return Factory{clusterClient: clusterClient, vendirOpts: vendirOpts, cmdRunner: cmdRunner, controllerVersion: controllerVersion}
}

func (f Factory) NewInline(opts v1alpha1.AppFetchInline, nsName string) *Inline {
	return NewInline(opts, nsName, f.clusterClient.CoreClient())
}

// TODO: pass v1alpha1.Vendir opts here once api is exapnded
func (f Factory) NewVendir(nsName string) *Vendir {
	return NewVendir(nsName, f.clusterClient.CoreClient(), f.vendirOpts, f.cmdRunner)
}

// NewVersionFetcher returns a fetcher which can fetch system component versions
func (f Factory) NewVersionFetcher() *VersionFetch {
	return NewVersionFetcher(f.clusterClient, f.controllerVersion)
}
