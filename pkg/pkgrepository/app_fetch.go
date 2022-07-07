// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"bytes"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

func (a *App) fetch(dstPath string) (string, exec.CmdRunResult) {
	if len(a.app.Spec.Fetch) == 0 {
		return "", exec.NewCmdRunResultWithErr(fmt.Errorf("Expected at least one fetch option"))
	}

	var result exec.CmdRunResult

	vendir := a.fetchFactory.NewVendir(a.app.Namespace)

	// Because vendir doesn't allow placing contents in the vendir root, we
	// place all contents in sub dirs. For backwards compatibility, we must
	// update dstPath to point to dstPath/0 if there is just one fetch step
	for i, fetch := range a.app.Spec.Fetch {
		err := vendir.AddDir(fetch, strconv.Itoa(i))
		if err != nil {
			result.AttachErrorf(fmt.Sprintf("Fetching (%d): ", i)+"%s", err)
			return "", result
		}
	}

	conf, err := vendir.ConfigBytes()
	if err != nil {
		result.AttachErrorf("Fetching: %v", err)
		return "", result
	}

	result = vendir.Run(conf, dstPath)

	// retry if error occurs before reporting failure.
	// This is mainly done to support private registry
	// authentication for images/bundles since placeholder
	// secrets may not be populated in time.
	if result.Error != nil && a.HasImageOrImgpkgBundle() {
		// Only retrying once resulted in flaky behavior
		// for private auth so use 3 iterations.
		for i := 0; i < 3; i++ {
			// Sleep for 2 seconds to allow secretgen-controller
			// to update placeholder secret(s).
			time.Sleep(2 * time.Second)
			newConf, err := vendir.ConfigBytes()
			if err != nil {
				result.AttachErrorf("Fetching: %v", err)
				return "", result
			}
			if bytes.Equal(conf, newConf) {
				// no secrets/configmaps have changed, no point in retrying
				continue
			}
			result = vendir.Run(newConf, dstPath)
			if result.Error == nil {
				break
			}
		}
		if result.Error != nil {
			return "", result
		}
	}

	// if only one fetch, update dstPath for backwards compatibility
	if len(a.app.Spec.Fetch) == 1 {
		dstPath = path.Join(dstPath, "0")
	}

	return dstPath, result
}
