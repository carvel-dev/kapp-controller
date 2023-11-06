// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"carvel.dev/kapp-controller/pkg/exec"
	"carvel.dev/kapp-controller/pkg/kubeconfig"
	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
)

// Factory allows to build various deployers.
// kapp deployer is standard, kapp privileged deployer
// should only be used for PackageRepository reconciliation.
type Factory struct {
	coreClient kubernetes.Interface
	kappConfig KappConfiguration

	kubeconfig *kubeconfig.Kubeconfig
	cmdRunner  exec.CmdRunner
}

// KappConfiguration provides a way to inject shared kapp settings.
type KappConfiguration interface {
	KappDeployRawOptions() []string
}

// NewFactory returns deploy factory.
func NewFactory(coreClient kubernetes.Interface, kubeconfig *kubeconfig.Kubeconfig,
	kappConfig KappConfiguration, cmdRunner exec.CmdRunner, _ logr.Logger) Factory {

	return Factory{coreClient, kappConfig, kubeconfig, cmdRunner}
}

// NewKapp configures and returns a deployer of type Kapp
func (f Factory) NewKapp(opts v1alpha1.AppDeployKapp, saName string,
	clusterOpts *v1alpha1.AppCluster, cancelCh chan struct{}, location kubeconfig.AccessLocation,
	defaultNamespace string, appNamespace string) (*Kapp, error) {

	clusterAccess, err := f.kubeconfig.ClusterAccess(saName, clusterOpts, location, defaultNamespace)
	if err != nil {
		return nil, err
	}

	const suffix string = ".app"
	return NewKapp(suffix, opts, clusterAccess,
		f.globalKappDeployRawOpts(), cancelCh, f.cmdRunner, appNamespace), nil
}

// NewKappPrivileged is used for package repositories where users aren't required to provide
// a service account, so it will install resources using its own privileges.
func (f Factory) NewKappPrivilegedForPackageRepository(opts v1alpha1.AppDeployKapp, clusterAccess kubeconfig.AccessInfo, cancelCh chan struct{}, appNamespace string) (*Kapp, error) {

	const suffix string = ".pkgr"

	kconfAccess := kubeconfig.AccessInfo{
		Name:                          clusterAccess.Name,
		Namespace:                     clusterAccess.Namespace,
		Kubeconfig:                    nil,
		DangerousUsePodServiceAccount: true,
	}

	return NewKapp(suffix, opts, kconfAccess, f.globalKappDeployRawOpts(), cancelCh, f.cmdRunner, appNamespace), nil
}

func (f Factory) globalKappDeployRawOpts() []string {
	if f.kappConfig != nil {
		return f.kappConfig.KappDeployRawOptions()
	}
	return nil
}
