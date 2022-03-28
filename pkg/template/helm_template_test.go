// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
)

// Create a temporary testing directory containing a Chart.yaml file
func newTestChartPath(specVersion string) (*memdir.TmpDir, error) {
	chartSpecContent := fmt.Sprintf(`
apiVersion: %s
appVersion: 5.6.2
name: myApp
`, specVersion)

	tmpDir := memdir.NewTmpDir("chartPath")
	err := tmpDir.Create()
	if err != nil {
		return nil, err
	}

	chartPath := tmpDir.Path()

	err = ioutil.WriteFile(filepath.Join(chartPath, "Chart.yaml"), []byte(chartSpecContent), 0600)
	if err != nil {
		return nil, err
	}

	return tmpDir, nil
}
