// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/clusterclient"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
)

// Factory allows to build various templaters e.g. ytt, cue.
type Factory struct {
	clusterClient  *clusterclient.ClusterClient
	fetchFactory   fetch.Factory
	kbldAllowBuild bool
	cmdRunner      exec.CmdRunner
}

// NewFactory returns template factory.
func NewFactory(clusterClient *clusterclient.ClusterClient, fetchFactory fetch.Factory, kbldAllowBuild bool, cmdRunner exec.CmdRunner) Factory {
	return Factory{clusterClient: clusterClient, fetchFactory: fetchFactory, kbldAllowBuild: kbldAllowBuild, cmdRunner: cmdRunner}
}

// NewYtt returns ytt template.
func (f Factory) NewYtt(opts v1alpha1.AppTemplateYtt, appContext AppContext) *Ytt {
	valuesFactory := ValuesFactory{
		appContext:   appContext,
		coreClient:   f.clusterClient.CoreClient(),
		fetchFactory: f.fetchFactory,
	}
	return NewYtt(opts, appContext, f.clusterClient, f.fetchFactory, f.cmdRunner, valuesFactory)
}

// NewKbld returns kbld template.
func (f Factory) NewKbld(opts v1alpha1.AppTemplateKbld, appContext AppContext) *Kbld {
	return NewKbld(opts, appContext, KbldOpts{AllowBuild: f.kbldAllowBuild}, f.cmdRunner)
}

// NewHelmTemplate returns helm template.
func (f Factory) NewHelmTemplate(opts v1alpha1.AppTemplateHelmTemplate, appContext AppContext) *HelmTemplate {
	valuesFactory := ValuesFactory{
		appContext:   appContext,
		coreClient:   f.clusterClient.CoreClient(),
		fetchFactory: f.fetchFactory,
	}
	return NewHelmTemplate(opts, appContext, f.cmdRunner, valuesFactory)
}

func (f Factory) NewSops(
	opts v1alpha1.AppTemplateSops, appContext AppContext) *Sops {
	return NewSops(opts, appContext, f.clusterClient.CoreClient(), f.cmdRunner)
}

// NewCue returns a Cue templater
func (f Factory) NewCue(opts v1alpha1.AppTemplateCue, appContext AppContext) Template {
	valuesFactory := ValuesFactory{
		appContext:   appContext,
		coreClient:   f.clusterClient.CoreClient(),
		fetchFactory: f.fetchFactory,
	}
	return newCue(opts, appContext, f.cmdRunner, valuesFactory, f.fetchFactory)
}
