// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"k8s.io/client-go/kubernetes"
)

// Factory allows to build various templaters e.g. ytt, cue.
type Factory struct {
	coreClient     kubernetes.Interface
	fetchFactory   fetch.Factory
	kbldAllowBuild bool
	cmdRunner      exec.CmdRunner
}

// NewFactory returns template factory.
func NewFactory(coreClient kubernetes.Interface, fetchFactory fetch.Factory,
	kbldAllowBuild bool, cmdRunner exec.CmdRunner) Factory {

	return Factory{coreClient, fetchFactory, kbldAllowBuild, cmdRunner}
}

// NewYtt returns ytt template.
func (f Factory) NewYtt(opts v1alpha1.AppTemplateYtt, appContext AppContext, additionalValues AdditionalDownwardAPIValues) *Ytt {
	return NewYtt(opts, appContext, f.coreClient, f.fetchFactory, f.cmdRunner, additionalValues)
}

// NewKbld returns kbld template.
func (f Factory) NewKbld(opts v1alpha1.AppTemplateKbld, appContext AppContext) *Kbld {
	return NewKbld(opts, appContext, KbldOpts{AllowBuild: f.kbldAllowBuild}, f.cmdRunner)
}

// NewHelmTemplate returns helm template.
func (f Factory) NewHelmTemplate(
	opts v1alpha1.AppTemplateHelmTemplate, appContext AppContext, additionalValues AdditionalDownwardAPIValues) *HelmTemplate {
	return NewHelmTemplate(opts, appContext, f.coreClient, f.cmdRunner, additionalValues)
}

func (f Factory) NewSops(
	opts v1alpha1.AppTemplateSops, appContext AppContext) *Sops {
	return NewSops(opts, appContext, f.coreClient, f.cmdRunner)
}

// NewCue returns a Cue templater
func (f Factory) NewCue(opts v1alpha1.AppTemplateCue, appContext AppContext, additionalValues AdditionalDownwardAPIValues) Template {
	return newCue(opts, appContext, f.coreClient, f.cmdRunner, additionalValues)
}
