// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"fmt"
	"strings"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	ctldep "github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

func (a *App) deploy(tplOutput string) exec.CmdRunResult {
	err := a.blockDeletion()
	if err != nil {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Blocking for deploy: %s", err))
	}

	if len(a.app.Spec.Deploy) != 1 {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Expected exactly one deploy option"))
	}

	dep := a.app.Spec.Deploy[0]

	switch {
	case dep.Kapp != nil:
		kapp, err := a.newKapp(*dep.Kapp, make(chan struct{}))
		if err != nil {
			return exec.NewCmdRunResultWithErr(fmt.Errorf("Preparing kapp: %s", err))
		}

		appName := a.app.Name
		if !strings.HasSuffix(appName, ".app") {
			appName += ".pkgr"
		}

		return kapp.Deploy(tplOutput, appName, a.startFlushingAllStatusUpdates, func(exec.CmdRunResult) {})
	default:
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Unsupported way to deploy"))
	}
}

func (a *App) delete() exec.CmdRunResult {
	if len(a.app.Spec.Deploy) != 1 {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Expected exactly one deploy option"))
	}

	var result exec.CmdRunResult

	if !a.app.Spec.NoopDelete {
		dep := a.app.Spec.Deploy[0]

		switch {
		case dep.Kapp != nil:
			kapp, err := a.newKapp(*dep.Kapp, make(chan struct{}))
			if err != nil {
				return exec.NewCmdRunResultWithErr(fmt.Errorf("Preparing kapp: %s", err))
			}

			result = kapp.Delete(a.startFlushingAllStatusUpdates, func(exec.CmdRunResult) {})

		default:
			result.AttachErrorf("%s", fmt.Errorf("Unsupported way to delete"))
		}
	}

	if result.Error == nil {
		err := a.unblockDeletion()
		if err != nil {
			return exec.NewCmdRunResultWithErr(fmt.Errorf("Unblocking for deploy: %s", err))
		}
	}

	return result
}

func (a *App) newKapp(kapp v1alpha1.AppDeployKapp, cancelCh chan struct{}) (*ctldep.Kapp, error) {
	genericOpts := ctldep.GenericOpts{Name: a.app.Name, Namespace: a.app.Namespace}
	return a.deployFactory.NewKappPrivileged(kapp, genericOpts, cancelCh)
}
