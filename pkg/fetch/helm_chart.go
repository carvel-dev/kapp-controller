 // Copyright 2020 VMware, Inc.
 // SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/memdir"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type HelmChart struct {
	opts       v1alpha1.AppFetchHelmChart
	nsName     string
	coreClient kubernetes.Interface
}

func NewHelmChart(opts v1alpha1.AppFetchHelmChart, nsName string, coreClient kubernetes.Interface) *HelmChart {
	return &HelmChart{opts, nsName, coreClient}
}

func (t *HelmChart) Retrieve(dstPath string) error {
	if len(t.opts.Name) == 0 {
		return fmt.Errorf("Expected non-empty name")
	}

	chartsDir := memdir.NewTmpDir("fetch-helm-chart")

	err := chartsDir.Create()
	if err != nil {
		return err
	}

	defer chartsDir.Remove()

	err = t.init()
	if err != nil {
		return err
	}

	err = t.fetch(chartsDir.Path())
	if err != nil {
		return err
	}

	chartPath, err := t.findChartDir(chartsDir.Path())
	if err != nil {
		return fmt.Errorf("Finding single helm chart: %s", err)
	}

	err = memdir.NewSubPath("").Extract(chartPath, dstPath)
	if err != nil {
		return err
	}

	return nil
}

func (t *HelmChart) init() error {
	args := []string{"init", "--client-only"}

	var stdoutBs, stderrBs bytes.Buffer

	cmd := exec.Command("helm", args...)
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Init helm: %s (stderr: %s)", err, stderrBs.String())
	}

	return nil
}

func (t *HelmChart) fetch(chartsPath string) error {
	args := []string{"fetch", t.opts.Name, "--untar", "--untardir", chartsPath}

	if len(t.opts.Version) > 0 {
		args = append(args, []string{"--version", t.opts.Version}...)
	}

	if t.opts.Repository != nil {
		if len(t.opts.Repository.URL) == 0 {
			return fmt.Errorf("Expected non-empty repository URL")
		}

		args = append(args, []string{"--repo", t.opts.Repository.URL}...)

		var err error

		args, err = t.addAuthArgs(args)
		if err != nil {
			return fmt.Errorf("Adding helm chart auth info: %s", err)
		}
	}

	var stdoutBs, stderrBs bytes.Buffer

	cmd := exec.Command("helm", args...)
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Fetching helm chart: %s (stderr: %s)", err, stderrBs.String())
	}

	return nil
}

func (t *HelmChart) findChartDir(chartsPath string) (string, error) {
	files, err := ioutil.ReadDir(chartsPath)
	if err != nil {
		return "", err
	}

	var result []os.FileInfo
	for _, file := range files {
		if file.IsDir() {
			result = append(result, file)
		}
	}

	if len(result) != 1 {
		return "", fmt.Errorf("Expected single directory in charts directory")
	}
	return filepath.Join(chartsPath, result[0].Name()), nil
}

func (t *HelmChart) addAuthArgs(args []string) ([]string, error) {
	var authArgs []string

	if t.opts.Repository.SecretRef != nil {
		secret, err := t.coreClient.CoreV1().Secrets(t.nsName).Get(
			t.opts.Repository.SecretRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		for name, val := range secret.Data {
			switch name {
			case corev1.BasicAuthUsernameKey:
				authArgs = append(authArgs, []string{"--username", string(val)}...)
			case corev1.BasicAuthPasswordKey:
				authArgs = append(authArgs, []string{"--password", string(val)}...)
			default:
				return nil, fmt.Errorf("Unknown secret field '%s' in secret '%s'", name, secret.Name)
			}
		}
	}

	return append(args, authArgs...), nil
}
