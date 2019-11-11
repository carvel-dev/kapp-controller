package deploy

import (
	"bytes"
	"io"
	goexec "os/exec"
	"strings"
	"sync"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/exec"
)

type Kapp struct {
	opts        v1alpha1.AppDeployKapp
	genericOpts GenericOpts
	cancelCh    chan struct{}
}

var _ Deploy = &Kapp{}

func NewKapp(opts v1alpha1.AppDeployKapp, genericOpts GenericOpts, cancelCh chan struct{}) *Kapp {
	return &Kapp{opts, genericOpts, cancelCh}
}

func (a *Kapp) Deploy(tplOutput string, changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {
	args := a.addDeployArgs([]string{"deploy", "-n", a.genericOpts.Namespace, "-a", a.managedName(), "-f", "-", "-y"})

	cmd := goexec.Command("kapp", args...)
	cmd.Stdin = strings.NewReader(tplOutput)
	stdoutBs, stderrBs := a.trackCmdOutput(cmd, changedFunc)

	err := exec.RunWithCancel(cmd, a.cancelCh)

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Deploying: %s", err)

	return result
}

func (a *Kapp) Delete(changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {
	args := a.addDeleteArgs([]string{"delete", "-n", a.genericOpts.Namespace, "-a", a.managedName(), "-y"})

	cmd := goexec.Command("kapp", args...)
	stdoutBs, stderrBs := a.trackCmdOutput(cmd, changedFunc)

	err := exec.RunWithCancel(cmd, a.cancelCh)

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Deleting: %s", err)

	return result
}

func (a *Kapp) Inspect() exec.CmdRunResult {
	var stdoutBs, stderrBs bytes.Buffer

	args := []string{
		"inspect", "-n", a.genericOpts.Namespace, "-a", a.managedName(), "-t",
		// PodMetrics rapidly get/created and removed, hence lets hide them
		// to avoid resource update churn
		// TODO is there a better way to deal with this?
		"--filter", `{"not":{"resource":{"kinds":["PodMetrics"]}}}`,
	}

	cmd := goexec.Command("kapp", args...)
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := exec.RunWithCancel(cmd, a.cancelCh)

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Inspecting: %s", err)

	return result
}

func (a *Kapp) ManagedName() string { return a.managedName() }

func (a *Kapp) trackCmdOutput(cmd *goexec.Cmd, changedFunc func(exec.CmdRunResult)) (*bytes.Buffer, *bytes.Buffer) {
	stdoutBs := &bytes.Buffer{}
	stderrBs := &bytes.Buffer{}

	liveResult := &exec.CmdRunResult{}
	liveResultMux := sync.Mutex{}

	cmd.Stdout = io.MultiWriter(stdoutBs, newBufferingWriter(func(data []byte) {
		liveResultMux.Lock()
		liveResult.Stdout += string(data)
		liveResultCopy := *liveResult
		liveResultMux.Unlock()
		changedFunc(liveResultCopy)
	}))

	cmd.Stderr = io.MultiWriter(stderrBs, newBufferingWriter(func(data []byte) {
		liveResultMux.Lock()
		liveResult.Stderr += string(data)
		liveResultCopy := *liveResult
		liveResultMux.Unlock()
		changedFunc(liveResultCopy)
	}))

	return stdoutBs, stderrBs
}

func (a *Kapp) managedName() string { return a.genericOpts.Name + "-ctrl" }

var (
	kappDisallowedOpts = map[string]bool{
		"--app":                true,
		"--namespace":          true,
		"--file":               true,
		"--kubeconfig":         true,
		"--kubeconfig-context": true,
	}
)

func (a *Kapp) addDeployArgs(args []string) []string {
	if len(a.opts.IntoNs) > 0 {
		args = append(args, []string{"--into-ns", a.opts.IntoNs}...)
	}

	for _, val := range a.opts.MapNs {
		args = append(args, []string{"--map-ns", val}...)
	}

	for _, opt := range a.opts.RawOptions {
		flag, err := exec.NewFlagFromString(opt)
		if err != nil {
			continue
		}

		if _, found := kappDisallowedOpts[flag.Name]; !found {
			args = append(args, opt)
		}
	}

	return args
}

func (a *Kapp) addDeleteArgs(args []string) []string {
	if a.opts.Delete != nil {
		for _, opt := range a.opts.Delete.RawOptions {
			flag, err := exec.NewFlagFromString(opt)
			if err != nil {
				continue
			}

			if _, found := kappDisallowedOpts[flag.Name]; !found {
				args = append(args, opt)
			}
		}
	}
	return args
}
