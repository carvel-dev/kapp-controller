/*
 * Copyright 2020 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package fetch

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Image struct {
	opts       v1alpha1.AppFetchImage
	nsName     string
	coreClient kubernetes.Interface
}

func NewImage(opts v1alpha1.AppFetchImage, nsName string, coreClient kubernetes.Interface) *Image {
	return &Image{opts, nsName, coreClient}
}

func (t *Image) Retrieve(dstPath string) error {
	if len(t.opts.URL) == 0 {
		return fmt.Errorf("Expected non-empty URL")
	}

	args := []string{"pull", "-i", t.opts.URL, "-o", dstPath}

	args, err := t.addAuthArgs(args)
	if err != nil {
		return err
	}

	var stdoutBs, stderrBs bytes.Buffer

	cmd := exec.Command("imgpkg", args...)
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Imgpkg: %s (stderr: %s)", err, stderrBs.String())
	}

	return nil
}

func (t *Image) addAuthArgs(args []string) ([]string, error) {
	var authArgs []string

	if t.opts.SecretRef != nil {
		secret, err := t.coreClient.CoreV1().Secrets(t.nsName).Get(t.opts.SecretRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		for name, val := range secret.Data {
			switch name {
			case "username":
				authArgs = append(authArgs, []string{"--registry-username", string(val)}...)
			case "password":
				authArgs = append(authArgs, []string{"--registry-password", string(val)}...)
			case "token":
				authArgs = append(authArgs, []string{"--registry-token", string(val)}...)
			default:
				return nil, fmt.Errorf("Unknown secret field '%s' in secret '%s'", name, secret.Name)
			}
		}
	}

	if len(authArgs) == 0 {
		authArgs = []string{"--registry-anon"}
	}

	return append(args, authArgs...), nil
}
