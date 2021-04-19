// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template_test

import (
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
)

func TestHelmTemplateCmdLookup(t *testing.T) {
	tests := []struct {
		chartSpecVersion string                // i.e v1 or v2
		binaryName       string                // returned helm binary name
		args             func(string) []string // returned helm template arguments
	}{
		{
			"v1",
			"helmv2",
			func(path string) []string {
				return []string{"template", path, "--name", "testRelease", "--namespace", "testNs"}
			},
		},
		{
			"v2",
			"helm",
			func(path string) []string {
				return []string{"template", "testRelease", path, "--namespace", "testNs", "--include-crds"}
			},
		},
		{
			"v3", // Non existent today but will fallback to newer version of the binary and format
			"helm",
			func(path string) []string {
				return []string{"template", "testRelease", path, "--namespace", "testNs", "--include-crds"}
			},
		},
	}

	for _, tc := range tests {
		tmpDir, err := newTestChartPath(tc.chartSpecVersion)
		if err != nil {
			t.Fatal(err)
		}
		defer tmpDir.Remove()

		cmd, err := template.NewHelmTemplateCmdArgs("testRelease", tmpDir.Path(), "testNs")
		if err != nil {
			t.Fatal(err)
		}

		if got, want := cmd.BinaryName, tc.binaryName; got != want {
			t.Errorf("got=%q, want=%q", got, want)
		}

		if got, want := cmd.Args, tc.args(tmpDir.Path()); !reflect.DeepEqual(got, want) {
			t.Errorf("got=%q, want=%q", got, want)
		}
	}

}

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

	err = ioutil.WriteFile(path.Join(chartPath, "Chart.yaml"), []byte(chartSpecContent), 0600)
	if err != nil {
		return nil, err
	}

	return tmpDir, nil
}
