// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"fmt"

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

		fmt.Println("XX Templated output to kapp:\n", tplOutput, "\nXX")

		return kapp.Deploy(appendRebaseRule(tplOutput), a.startFlushingAllStatusUpdates, func(exec.CmdRunResult) {})

	default:
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Unsupported way to deploy"))
	}
}

func appendRebaseRule(tplOutput string) string {
	return tplOutput + `

---
apiVersion: kapp.k14s.io/v1alpha1
kind: Config
rebaseRules:
- ytt:
    overlayContractV1:
      overlay.yml: |
        #@ load("@ytt:data", "data")
        #@ load("@ytt:json", "json")
        #@ load("@ytt:overlay", "overlay")

        #@ def get_rev(annotations):
        #@  if hasattr(annotations, "packaging.carvel.dev/repo-rev"):
        #@    return int(annotations["packaging.carvel.dev/repo-rev"])
        #@  else:
        #@    return -1
        #@  end
        #@ end

        #@ if not hasattr(data.values.existing.metadata.annotations, "packaging.carvel.dev/package-repository-ref"):
        #@   msg = "Error: cannot overwrite package " + data.values.existing.metadata.name + " because it was not created by a package repository."
        #@   print(msg)
        #@   fail(msg)
        #@ end

        #@ if json.encode(data.values.existing.spec) == json.encode(data.values.new.spec):
        #@overlay/match by=overlay.all
        ---
        metadata:
          #@overlay/match missing_ok=True
          annotations:
            #@overlay/match missing_ok=True
            kapp.k14s.io/noop: ""
        #@ elif get_rev(data.values.existing.metadata.annotations) > get_rev(data.values.new.metadata.annotations):
        #@overlay/match by=overlay.all
        ---
        metadata:
          #@overlay/match missing_ok=True
          annotations:
            #@overlay/match missing_ok=True
            kapp.k14s.io/noop: ""
        #@ elif get_rev(data.values.existing.metadata.annotations) < get_rev(data.values.new.metadata.annotations):
        #@ print("replacing existing older rev with newer rev")
        #@ end
  resourceMatchers:
  - apiVersionKindMatcher: {apiVersion: data.packaging.carvel.dev/v1alpha1, kind: Package}
`
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
	kapp.RawOptions = append(kapp.RawOptions, "--dangerous-override-ownership-of-existing-resources=true")
	return a.deployFactory.NewKappPrivileged(kapp, genericOpts, cancelCh)
}
