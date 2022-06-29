// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"k8s.io/client-go/kubernetes"
)

// Factory allows to build various deployers.
// kapp deployer is standard, kapp privileged deployer
// should only be used for PackageRepository reconciliation.
type Factory struct {
	coreClient kubernetes.Interface
	kappConfig KappConfiguration

	kubeconfigSecrets *KubeconfigSecrets
	serviceAccounts   *ServiceAccounts

	cmdRunner exec.CmdRunner
}

// KappConfiguration provides a way to inject shared kapp settings.
type KappConfiguration interface {
	KappDeployRawOptions() []string
}

// NewFactory returns deploy factory.
func NewFactory(coreClient kubernetes.Interface,
	kappConfig KappConfiguration, cmdRunner exec.CmdRunner, log logr.Logger) Factory {

	return Factory{coreClient, kappConfig,
		NewKubeconfigSecrets(coreClient), NewServiceAccounts(coreClient, log), cmdRunner}
}

func (f Factory) NewKapp(opts v1alpha1.AppDeployKapp, saName string,
	clusterOpts *v1alpha1.AppCluster, genericOpts GenericOpts, cancelCh chan struct{}) (*Kapp, error) {

	var err error
	var processedGenericOpts ProcessedGenericOpts
	const suffix string = ".app"

	switch {
	case len(saName) > 0:
		processedGenericOpts, err = f.serviceAccounts.Find(genericOpts, saName)
		if err != nil {
			return nil, err
		}

	case clusterOpts != nil:
		processedGenericOpts, err = f.kubeconfigSecrets.Find(genericOpts, clusterOpts)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Expected service account or cluster specified")
	}

	return NewKapp(suffix, opts, processedGenericOpts,
		f.globalKappDeployRawOpts(), cancelCh, f.cmdRunner), nil
}

// NewKappPrivileged is used for package repositories where users aren't required to provide
// a service account, so it will install resources using its own privileges.
func (f Factory) NewKappPrivilegedForPackageRepository(opts v1alpha1.AppDeployKapp,
	genericOpts GenericOpts, cancelCh chan struct{}) (*Kapp, error) {

	const suffix string = ".pkgr"

	pgo := ProcessedGenericOpts{
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
