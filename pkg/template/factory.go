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

func (f Factory) NewYtt(opts v1alpha1.AppTemplateYtt, genericOpts GenericOpts) *Ytt {
	return NewYtt(opts, genericOpts, f.coreClient, f.fetchFactory, f.cmdRunner)
}

func (f Factory) NewKbld(opts v1alpha1.AppTemplateKbld, genericOpts GenericOpts) *Kbld {
	return NewKbld(opts, genericOpts, KbldOpts{AllowBuild: f.kbldAllowBuild}, f.cmdRunner)
}

func (f Factory) NewHelmTemplate(
	opts v1alpha1.AppTemplateHelmTemplate, genericOpts GenericOpts) *HelmTemplate {
	return NewHelmTemplate(opts, genericOpts, f.coreClient, f.cmdRunner)
}

func (f Factory) NewSops(
	opts v1alpha1.AppTemplateSops, genericOpts GenericOpts) *Sops {
	return NewSops(opts, genericOpts, f.coreClient, f.cmdRunner)
}

// NewCue returns a Cue templater
func (f Factory) NewCue(opts v1alpha1.AppTemplateCue, genericOpts GenericOpts) Template {
	return newCue(opts, genericOpts, f.coreClient, f.cmdRunner)
}
