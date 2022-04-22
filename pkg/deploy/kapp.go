// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"bytes"
	"fmt"
	"os"
	goexec "os/exec"
	"strings"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

const (
	// TODO not a great way to determine whether
	// kapp found changes and started to apply them
	applyOutputMarker = " ---- applying "
)

type Kapp struct {
	opts        v1alpha1.AppDeployKapp
	genericOpts ProcessedGenericOpts
	cancelCh    chan struct{}
	cmdRunner   exec.CmdRunner
}

var _ Deploy = &Kapp{}

// NewKapp takes the kapp yaml from spec.deploy.kapp as arg kapp,
// additional info from the larger app resource (e.g. service account, name, namespace) as genericOpts,
// and a cancel channel that gets passed through to the exec call that runs kapp.
func NewKapp(opts v1alpha1.AppDeployKapp, genericOpts ProcessedGenericOpts,
	cancelCh chan struct{}, cmdRunner exec.CmdRunner) *Kapp {

	return &Kapp{opts, genericOpts, cancelCh, cmdRunner}
}

func (a *Kapp) Deploy(tplOutput string, startedApplyingFunc func(),
	changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {

	args, err := a.addDeployArgs([]string{"deploy", "-f", "-"})
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	args, env := a.addGenericArgs(args)

	cmd := goexec.Command("kapp", args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdin = strings.NewReader(tplOutput)

	resultBuf, doneTrackingOutputCh := a.trackCmdOutput(cmd, startedApplyingFunc, changedFunc)

	err = a.cmdRunner.RunWithCancel(cmd, a.cancelCh)
	close(doneTrackingOutputCh)

	result := resultBuf.Copy()
	result.AttachErrorf("Deploying: %s", err)

	return result
}

func (a *Kapp) Delete(startedApplyingFunc func(), changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {
	args, err := a.addDeleteArgs([]string{"delete"})
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	args, env := a.addGenericArgs(args)

	cmd := goexec.Command("kapp", args...)
	cmd.Env = append(os.Environ(), env...)

	resultBuf, doneTrackingOutputCh := a.trackCmdOutput(cmd, startedApplyingFunc, changedFunc)

	err = a.cmdRunner.RunWithCancel(cmd, a.cancelCh)
	close(doneTrackingOutputCh)

	result := resultBuf.Copy()
	result.AttachErrorf("Deleting: %s", err)

	return result
}

func (a *Kapp) Inspect() exec.CmdRunResult {
	args, err := a.addInspectArgs([]string{
		"inspect",
		// PodMetrics rapidly get/created and removed, hence lets hide them
		// to avoid resource update churn
		// TODO is there a better way to deal with this?
		"--filter", `{"not":{"resource":{"kinds":["PodMetrics"]}}}`,
	})
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	args, env := a.addGenericArgs(args)

	var stdoutBs, stderrBs bytes.Buffer

	cmd := goexec.Command("kapp", args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err = a.cmdRunner.RunWithCancel(cmd, a.cancelCh)

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Inspecting: %s", err)

	return result
}

func (a *Kapp) trackCmdOutput(cmd *goexec.Cmd, startedApplyingFunc func(),
	changedFunc func(exec.CmdRunResult)) (*CmdRunResultBuffer, chan struct{}) {

	liveResult := NewCmdRunResultBuffer()
	doneCh := make(chan struct{})

	cmd.Stdout = WriterFunc(liveResult.WriteStdout)
	cmd.Stderr = WriterFunc(liveResult.WriteStderr)

	// Serialize status updates
	go func() {
		check := time.NewTicker(2 * time.Second)
		defer check.Stop()

		for {
			select {
			case <-check.C:
				resultCopy := liveResult.Copy()

				changedFunc(resultCopy)
				if strings.Contains(resultCopy.Stdout, applyOutputMarker) {
					startedApplyingFunc()
				}

			case <-doneCh:
				return
			}
		}
	}()

	return liveResult, doneCh
}

func (a *Kapp) managedName() string { return a.genericOpts.Name + "-ctrl" }

func (a *Kapp) addDeployArgs(args []string) ([]string, error) {
	if len(a.opts.IntoNs) > 0 {
		args = append(args, []string{"--into-ns", a.opts.IntoNs}...)
	}

	for _, val := range a.opts.MapNs {
		args = append(args, []string{"--map-ns", val}...)
	}

	return a.addRawOpts(args, a.opts.RawOptions, kappAllowedDeployFlagSet)
}

func (a *Kapp) addDeleteArgs(args []string) ([]string, error) {
	if a.opts.Delete != nil {
		return a.addRawOpts(args, a.opts.Delete.RawOptions, kappAllowedDeleteFlagSet)
	}
	return args, nil
}

func (a *Kapp) addInspectArgs(args []string) ([]string, error) {
	if a.opts.Inspect != nil {
		return a.addRawOpts(args, a.opts.Inspect.RawOptions, kappAllowedInspectFlagSet)
	}
	return args, nil
}

func (a *Kapp) addRawOpts(args []string, opts []string, allowedFlagSet exec.FlagSet) ([]string, error) {
	for _, opt := range opts {
		flag, err := exec.NewFlagFromString(opt)
		if err != nil {
			return nil, err
		}
		if allowedFlagSet.Includes(flag.Name) {
			args = append(args, opt)
		} else {
			return nil, fmt.Errorf("Unexpected flag '%s' specified (either forbidden or unknown)", flag.Name)
		}
	}
	return args, nil
}

func (a *Kapp) addGenericArgs(args []string) ([]string, []string) {
	args = append(args, []string{"--app", a.managedName()}...)
	env := []string{}

	if len(a.genericOpts.Namespace) > 0 {
		args = append(args, []string{"--namespace", a.genericOpts.Namespace}...)
	}

	switch {
	case a.genericOpts.Kubeconfig != nil:
		env = append(env, "KAPP_KUBECONFIG_YAML="+a.genericOpts.Kubeconfig.AsYAML())
		args = append(args, "--kubeconfig=/dev/null") // not used due to above env var
	case a.genericOpts.DangerousUsePodServiceAccount:
		// do nothing
	default:
		panic("Internal inconsistency: Unknown kapp service account configuration")
	}

	args = append(args, "--yes")

	return args, env
}
