package fetch

import (
	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

type Factory struct {
	coreClient kubernetes.Interface
}

func NewFactory(coreClient kubernetes.Interface) Factory {
	return Factory{coreClient}
}

func (f Factory) NewHelmChart(opts v1alpha1.AppFetchHelmChart, nsName string) *HelmChart {
	return NewHelmChart(opts, nsName, f.coreClient)
}

func (f Factory) NewHTTP(opts v1alpha1.AppFetchHTTP, nsName string) *HTTP {
	return NewHTTP(opts, nsName, f.coreClient)
}

func (f Factory) NewGit(opts v1alpha1.AppFetchGit, nsName string) *Git {
	return NewGit(opts, nsName, f.coreClient)
}

func (f Factory) NewImage(opts v1alpha1.AppFetchImage, nsName string) *Image {
	return NewImage(opts, nsName, f.coreClient)
}

func (f Factory) NewInline(opts v1alpha1.AppFetchInline, nsName string) *Inline {
	return NewInline(opts, nsName, f.coreClient)
}
