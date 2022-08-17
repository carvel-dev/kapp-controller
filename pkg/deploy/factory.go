// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"fmt"

	"github.com/go-logr/logr"
	semver "github.com/k14s/semver/v4"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

// NewKapp configures and returns a deployer of type Kapp
func (f Factory) NewKapp(opts v1alpha1.AppDeployKapp, saName string,
	clusterOpts *v1alpha1.AppCluster, genericOpts GenericOpts, cancelCh chan struct{}) (*Kapp, error) {

	processedGenericOpts, err := f.processOpts(saName, clusterOpts, genericOpts)
	if err != nil {
		return nil, err
	}
	const suffix string = ".app"
	return NewKapp(suffix, opts, processedGenericOpts,
		f.globalKappDeployRawOpts(), cancelCh, f.cmdRunner), nil
}

// processOpts takes generic opts and a ServiceAccount Name, and returns a populated kubeconfig that can connect to a cluster.
// if the saName is empty then you'll connect to a cluster via the clusterOpts inside the genericOpts, otherwise you'll use the specified SA.
func (f Factory) processOpts(saName string, clusterOpts *v1alpha1.AppCluster, genericOpts GenericOpts) (ProcessedGenericOpts, error) {
	var err error
	var processedGenericOpts ProcessedGenericOpts

	switch {
	case len(saName) > 0:
		processedGenericOpts, err = f.serviceAccounts.Find(genericOpts, saName)
		if err != nil {
			return ProcessedGenericOpts{}, err
		}

	case clusterOpts != nil:
		processedGenericOpts, err = f.kubeconfigSecrets.Find(genericOpts, clusterOpts)
		if err != nil {
			return ProcessedGenericOpts{}, err
		}

	default:
		return ProcessedGenericOpts{}, fmt.Errorf("Expected service account or cluster specified")
	}
	return processedGenericOpts, nil
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

// GetClusterVersion returns the kubernetes API version for the cluster which has been supplied to kapp-controller via a kubeconfig
func (f Factory) GetClusterVersion(saName string, specCluster *v1alpha1.AppCluster, genericOpts GenericOpts, log logr.Logger) (semver.Version, error) {
	switch {
	case len(saName) > 0:
		version, err := f.coreClient.Discovery().ServerVersion()
		if err != nil {
			return semver.Version{}, err
		}

		return semver.ParseTolerant(version.GitVersion)
	case specCluster != nil:
		processedGenericOpts, err := f.processOpts(saName, specCluster, GenericOpts{Name: genericOpts.Name, Namespace: genericOpts.Namespace})
		if err != nil {
			return semver.Version{}, err
		}

		config, err := clientcmd.RESTConfigFromKubeConfig([]byte(processedGenericOpts.Kubeconfig.AsYAML()))
		if err != nil {
			return semver.Version{}, err
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return semver.Version{}, err
		}

		version, err := clientset.Discovery().ServerVersion()
		if err != nil {
			return semver.Version{}, err
		}

		return semver.ParseTolerant(version.GitVersion)
	default:
		return semver.Version{}, fmt.Errorf("Expected service account or cluster specified")
	}
}
