package deploy

import (
	"bytes"
	"os"
	goexec "os/exec"
	"strings"
	"time"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/exec"
)

const (
	// TODO not a great way to determine whether
	// kapp found changes and started to apply them
	applyOutputMarker = " ---- applying "
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

func (a *Kapp) Deploy(tplOutput string, startedApplyingFunc func(),
	changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {

	args := a.addDeployArgs([]string{"deploy", "-f", "-"})
	args, env := a.addGenericArgs(args)

	cmd := goexec.Command("kapp", args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdin = strings.NewReader(tplOutput)

	resultBuf, doneTrackingOutputCh := a.trackCmdOutput(cmd, startedApplyingFunc, changedFunc)

	err := exec.RunWithCancel(cmd, a.cancelCh)
	close(doneTrackingOutputCh)

	result := resultBuf.Copy()
	result.AttachErrorf("Deploying: %s", err)

	return result
}

func (a *Kapp) Delete(startedApplyingFunc func(), changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {
	args := a.addDeleteArgs([]string{"delete"})
	args, env := a.addGenericArgs(args)

	cmd := goexec.Command("kapp", args...)
	cmd.Env = append(os.Environ(), env...)

	resultBuf, doneTrackingOutputCh := a.trackCmdOutput(cmd, startedApplyingFunc, changedFunc)

	err := exec.RunWithCancel(cmd, a.cancelCh)
	close(doneTrackingOutputCh)

	result := resultBuf.Copy()
	result.AttachErrorf("Deleting: %s", err)

	return result
}

func (a *Kapp) Inspect() exec.CmdRunResult {
	args := a.addInspectArgs([]string{
		"inspect",
		// PodMetrics rapidly get/created and removed, hence lets hide them
		// to avoid resource update churn
		// TODO is there a better way to deal with this?
		"--filter", `{"not":{"resource":{"kinds":["PodMetrics"]}}}`,
	})

	args, env := a.addGenericArgs(args)

	var stdoutBs, stderrBs bytes.Buffer

	cmd := goexec.Command("kapp", args...)
	cmd.Env = append(os.Environ(), env...)
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

func (a *Kapp) addInspectArgs(args []string) []string {
	if a.opts.Inspect != nil {
		for _, opt := range a.opts.Inspect.RawOptions {
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

func (a *Kapp) addGenericArgs(args []string) ([]string, []string) {
	args = append(args, []string{"--app", a.managedName()}...)
	env := []string{}

	if len(a.genericOpts.Namespace) > 0 {
		args = append(args, []string{"--namespace", a.genericOpts.Namespace}...)
	}

	if len(a.genericOpts.KubeconfigYAML) > 0 {
		env = append(env, "KAPP_KUBECONFIG_YAML="+a.genericOpts.KubeconfigYAML)
	}

	args = append(args, "--yes")

	return args, env
}
