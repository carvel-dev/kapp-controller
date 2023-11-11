// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	goexec "os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/kubeconfig"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	"sigs.k8s.io/yaml"
)

const (
	// TODO not a great way to determine whether
	// kapp found changes and started to apply them
	applyOutputMarker = " ---- applying "
)

type Kapp struct {
	appSuffix           string
	opts                v1alpha1.AppDeployKapp
	clusterAccess       kubeconfig.AccessInfo
	globalDeployRawOpts []string
	cancelCh            chan struct{}
	cmdRunner           exec.CmdRunner
	appMeta             *Meta
	appNamespace        string
}

var _ Deploy = &Kapp{}

// NewKapp takes the kapp yaml from spec.deploy.kapp as arg kapp,
// additional info from the larger app resource (e.g. service account, name, namespace) as genericOpts,
// and a cancel channel that gets passed through to the exec call that runs kapp.
func NewKapp(appSuffix string, opts v1alpha1.AppDeployKapp, clusterAccess kubeconfig.AccessInfo,
	globalDeployRawOpts []string, cancelCh chan struct{}, cmdRunner exec.CmdRunner, appNamespace string) *Kapp {

	return &Kapp{appSuffix, opts, clusterAccess, globalDeployRawOpts, cancelCh, cmdRunner, nil, appNamespace}
}

// Deploy takes the output from templating, and the app name,
// it shells out, running kapp deploy ...
func (a *Kapp) Deploy(tplOutput string, startedApplyingFunc func(),
	changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {

	tmpMetadataDir := memdir.NewTmpDir("app_metadata")
	defer tmpMetadataDir.Remove()

	err := tmpMetadataDir.Create()
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	metadataFile := filepath.Join(tmpMetadataDir.Path(), "app-metadata.yml")

	args, err := a.addDeployArgs([]string{"deploy", "--app-metadata-file-output", metadataFile, "--prev-app", a.oldManagedName(), "-f", "-"})
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	args, env := a.addGenericArgs(args, a.clusterAccess.Name+a.appSuffix)

	cmd := goexec.Command("kapp", args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdin = strings.NewReader(tplOutput)

	resultBuf, doneTrackingOutputCh := a.trackCmdOutput(cmd, startedApplyingFunc, changedFunc)

	err = a.cmdRunner.RunWithCancel(cmd, a.cancelCh)
	close(doneTrackingOutputCh)

	result := resultBuf.Copy()
	result.AttachErrorf("Deploying: %s", err)

	a.appMeta = nil
	a.trySaveAppMeta(metadataFile)

	return result
}

// Delete takes the app name, it shells out, running kapp delete ...
func (a *Kapp) Delete(startedApplyingFunc func(), changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {
	args, err := a.addDeleteArgs([]string{"delete", "--prev-app", a.oldManagedName()})
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	args, env := a.addGenericArgs(args, a.clusterAccess.Name+a.appSuffix)

	cmd := goexec.Command("kapp", args...)
	cmd.Env = append(os.Environ(), env...)

	resultBuf, doneTrackingOutputCh := a.trackCmdOutput(cmd, startedApplyingFunc, changedFunc)

	err = a.cmdRunner.RunWithCancel(cmd, a.cancelCh)
	close(doneTrackingOutputCh)

	result := resultBuf.Copy()
	result.AttachErrorf("Deleting: %s", err)

	return result
}

// Inspect takes the app name, it shells out, running kapp inspect ...
func (a *Kapp) Inspect() exec.CmdRunResult {
	args, err := a.addInspectArgs([]string{
		"inspect",
		// PodMetrics rapidly get/created and removed, hence lets hide them
		// to avoid resource update churn
		// TODO is there a better way to deal with this?
		"--filter", `{"not":{"resource":{"kinds":["PodMetrics"]}}}`,
		"--tty",
	})
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	args, env := a.addGenericArgs(args, a.clusterAccess.Name+a.appSuffix)

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

// Meta contains app meta allowing for an AppCR to surface the associated namespaces and GKs
type Meta struct {
	LabelKey   string `yaml:"labelKey"`
	LabelValue string `yaml:"labelValue"`
	LastChange struct {
		Namespaces []string `yaml:"namespaces"`
	} `yaml:"lastChange"`
	UsedGKs []struct {
		Group string `yaml:"Group"`
		Kind  string `yaml:"Kind"`
	} `yaml:"usedGKs"`
}

// InternalAppMeta exposes the internal configmap kapp maintains on every app deploy
func (a *Kapp) InternalAppMeta() (*Meta, error) {
	if a.appMeta == nil {
		return nil, errors.New("Unable to retrieve kapp internal config map")
	}

	return a.appMeta, nil
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

// This is the old naming schema for KC owned kapp apps.
// The new convention is x.app for AppCRs / PKGIs and x.pkgr for PackageRepositories.
func (a *Kapp) oldManagedName() string { return a.clusterAccess.Name + "-ctrl" }

func (a *Kapp) addDeployArgs(args []string) ([]string, error) {
	if len(a.opts.IntoNs) > 0 {
		args = append(args, []string{"--into-ns", a.opts.IntoNs}...)
	}

	for _, val := range a.opts.MapNs {
		args = append(args, []string{"--map-ns", val}...)
	}

	// Global raw options are applied first to be able to override them within an App
	args, err := a.addRawOpts(args, a.globalDeployRawOpts, kappAllowedDeployFlagSet)
	if err != nil {
		return nil, err
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

func (a *Kapp) addGenericArgs(args []string, appName string) ([]string, []string) {
	args = append(args, []string{"--app", appName}...)
	env := []string{}

	if len(a.clusterAccess.Namespace) > 0 {
		args = append(args, []string{"--namespace", a.clusterAccess.Namespace}...)
	}

	if len(a.clusterAccess.DeployNamespace) > 0 {
		args = append(args, []string{"--app-namespace", a.clusterAccess.DeployNamespace}...)
	}

	switch {
	case a.clusterAccess.Kubeconfig != nil:
		env = append(env, "KAPP_KUBECONFIG_YAML="+a.clusterAccess.Kubeconfig.AsYAML())
		args = append(args, "--kubeconfig=/dev/null") // not used due to above env var
	case a.clusterAccess.DangerousUsePodServiceAccount:
		// do nothing
	default:
		panic("Internal inconsistency: Unknown kapp service account configuration")
	}

	args = append(args, "--yes")

	return args, env
}

// trySaveAppMeta if unable to save the kapp configmap metadata, then continue and do not fail the deploy.
func (a *Kapp) trySaveAppMeta(metadataFileName string) {
	metadataFile, err := os.ReadFile(metadataFileName)
	if err != nil {
		return
	}

	appMetadata := Meta{}
	err = yaml.Unmarshal(metadataFile, &appMetadata)
	if err != nil {
		return
	}

	a.appMeta = &appMetadata
}
