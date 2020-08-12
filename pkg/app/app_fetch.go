/*
 * Copyright 2020 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package app

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/exec"
	"github.com/k14s/kapp-controller/pkg/memdir"
)

func (a *App) fetch(dstPath string) exec.CmdRunResult {
	if len(a.app.Spec.Fetch) == 0 {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Expected at least one fetch option"))
	}

	var result exec.CmdRunResult

	if len(a.app.Spec.Fetch) == 1 {
		fetch := a.app.Spec.Fetch[0]

		err := a.fetchOne(fetch, dstPath)
		if err != nil {
			result.AttachErrorf("Fetching (0): %s", err)
		}
	} else {
		for i, fetch := range a.app.Spec.Fetch {
			dstSubPath := path.Join(dstPath, strconv.Itoa(i))

			err := os.Mkdir(dstSubPath, os.FileMode(0700))
			if err != nil {
				result.AttachErrorf(fmt.Sprintf("Fetching (%d): ", i)+"%s", err)
				break
			}

			err = a.fetchOne(fetch, dstSubPath)
			if err != nil {
				result.AttachErrorf(fmt.Sprintf("Fetching (%d): ", i)+"%s", err)
				break
			}
		}
	}

	return result
}

func (a *App) fetchOne(fetch v1alpha1.AppFetch, dstPath string) error {
	tmpDstDir := memdir.NewTmpDir("fetch-one")

	err := tmpDstDir.Create()
	if err != nil {
		return err
	}

	defer tmpDstDir.Remove()

	var subPath string

	switch {
	case fetch.Inline != nil:
		err = a.fetchFactory.NewInline(*fetch.Inline, a.app.Namespace).Retrieve(tmpDstDir.Path())
		if err != nil {
			return err
		}

	case fetch.Image != nil:
		subPath = fetch.Image.SubPath
		err = a.fetchFactory.NewImage(*fetch.Image, a.app.Namespace).Retrieve(tmpDstDir.Path())
		if err != nil {
			return fmt.Errorf("Fetching registry image: %s", err)
		}

	case fetch.HTTP != nil:
		subPath = fetch.HTTP.SubPath
		err = a.fetchFactory.NewHTTP(*fetch.HTTP, a.app.Namespace).Retrieve(tmpDstDir.Path())
		if err != nil {
			return fmt.Errorf("Fetching HTTP asset: %s", err)
		}

	case fetch.Git != nil:
		subPath = fetch.Git.SubPath
		err = a.fetchFactory.NewGit(*fetch.Git, a.app.Namespace).Retrieve(tmpDstDir.Path())
		if err != nil {
			return fmt.Errorf("Fetching git repo: %s", err)
		}

	case fetch.HelmChart != nil:
		err = a.fetchFactory.NewHelmChart(*fetch.HelmChart, a.app.Namespace).Retrieve(tmpDstDir.Path())
		if err != nil {
			return fmt.Errorf("Fetching helm chart: %s", err)
		}

	default:
		return fmt.Errorf("Unsupported way to fetch templates")
	}

	return memdir.NewSubPath(subPath).Extract(tmpDstDir.Path(), dstPath)
}
