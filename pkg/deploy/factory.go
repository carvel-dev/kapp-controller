// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"fmt"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/clusterclient"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

// Factory allows to build various deployers.
// kapp deployer is standard, kapp privileged deployer
// should only be used for PackageRepository reconciliation.
type Factory struct {
	kappConfig    KappConfiguration
	clusterClient *clusterclient.ClusterClient

	cmdRunner exec.CmdRunner
}

// KappConfiguration provides a way to inject shared kapp settings.
type KappConfiguration interface {
	KappDeployRawOptions() []string
}

// NewFactory returns deploy factory.
func NewFactory(clusterClient *clusterclient.ClusterClient, kappConfig KappConfiguration, cmdRunner exec.CmdRunner) Factory {
	return Factory{clusterClient: clusterClient, kappConfig: kappConfig, cmdRunner: cmdRunner}
}

// NewKapp configures and returns a deployer of type Kapp
func (f Factory) NewKapp(opts v1alpha1.AppDeployKapp, saName string,
	clusterOpts *v1alpha1.AppCluster, genericOpts clusterclient.GenericOpts, cancelCh chan struct{}) (*Kapp, error) {

	processedGenericOpts, err := f.clusterClient.ProcessOpts(saName, clusterOpts, genericOpts)
	fmt.Println(processedGenericOpts.Kubeconfig)
	if err != nil {
		return nil, err
	}
	const suffix string = ".app"
	return NewKapp(suffix, opts, processedGenericOpts,
		f.globalKappDeployRawOpts(), cancelCh, f.cmdRunner), nil
}

// NewKappPrivileged is used for package repositories where users aren't required to provide
// a service account, so it will install resources using its own privileges.
func (f Factory) NewKappPrivilegedForPackageRepository(opts v1alpha1.AppDeployKapp,
	genericOpts clusterclient.GenericOpts, cancelCh chan struct{}) (*Kapp, error) {

	const suffix string = ".pkgr"

	pgo := clusterclient.ProcessedGenericOpts{
		Name:      genericOpts.Name,
		Namespace: genericOpts.Namespace,
		// Just use the default service account. Mainly
		// used for PackageRepos now so users do not need
		// to specify serviceaccount via PackageRepo CR.
		Kubeconfig:                    nil,
		DangerousUsePodServiceAccount: true,
	}

	return NewKapp(suffix, opts, pgo, f.globalKappDeployRawOpts(), cancelCh, f.cmdRunner), nil
}

func (f Factory) globalKappDeployRawOpts() []string {
	if f.kappConfig != nil {
		return f.kappConfig.KappDeployRawOptions()
	}
	return nil
}
