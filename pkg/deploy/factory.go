package deploy

import (
	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

type Factory struct {
	coreClient        kubernetes.Interface
	kubeconfigSecrets *KubeconfigSecrets
}

func NewFactory(coreClient kubernetes.Interface) Factory {
	return Factory{coreClient, NewKubeconfigSecrets(coreClient)}
}

func (f Factory) NewKapp(opts v1alpha1.AppDeployKapp, clusterOpts *v1alpha1.AppCluster,
	genericOpts GenericOpts, cancelCh chan struct{}) (*Kapp, error) {

	genericOpts, err := f.kubeconfigSecrets.Find(genericOpts, clusterOpts)
	if err != nil {
		return nil, err
	}

	return NewKapp(opts, genericOpts, cancelCh), nil
}
