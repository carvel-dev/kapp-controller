// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"fmt"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

type Factory struct {
	coreClient kubernetes.Interface

	kubeconfigSecrets *KubeconfigSecrets
	serviceAccounts   *ServiceAccounts
}

func NewFactory(coreClient kubernetes.Interface) Factory {
	return Factory{coreClient, NewKubeconfigSecrets(coreClient), NewServiceAccounts(coreClient)}
}

func (f Factory) NewKapp(opts v1alpha1.AppDeployKapp, saName string,
	clusterOpts *v1alpha1.AppCluster, genericOpts GenericOpts, cancelCh chan struct{}) (*Kapp, error) {

	var err error

	switch {
	case len(saName) > 0:
		genericOpts, err = f.serviceAccounts.Find(genericOpts, saName)
		if err != nil {
			return nil, err
		}

	case clusterOpts != nil:
		genericOpts, err = f.kubeconfigSecrets.Find(genericOpts, clusterOpts)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Expected service account or cluster specified")
	}

	return NewKapp(opts, genericOpts, cancelCh), nil
}
