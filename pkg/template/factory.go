package template

import (
	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/fetch"
	"k8s.io/client-go/kubernetes"
)

type Factory struct {
	coreClient   kubernetes.Interface
	fetchFactory fetch.Factory
}

func NewFactory(coreClient kubernetes.Interface, fetchFactory fetch.Factory) Factory {
	return Factory{coreClient, fetchFactory}
}

func (f Factory) NewYtt(opts v1alpha1.AppTemplateYtt, genericOpts GenericOpts) *Ytt {
	return NewYtt(opts, genericOpts, f.fetchFactory)
}

func (f Factory) NewKbld(opts v1alpha1.AppTemplateKbld, genericOpts GenericOpts) *Kbld {
	return NewKbld(opts, genericOpts)
}

func (f Factory) NewSops(opts v1alpha1.AppTemplateSops, genericOpts GenericOpts) *Sops {
	return NewSops(opts, genericOpts)
}

func (f Factory) NewHelmTemplate(
	opts v1alpha1.AppTemplateHelmTemplate, genericOpts GenericOpts) *HelmTemplate {
	return NewHelmTemplate(opts, genericOpts, f.coreClient)
}
