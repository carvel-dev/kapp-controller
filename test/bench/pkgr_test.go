// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package bench

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// we'll assert against these thresholds below; maps represent "for x packages, what's the max allowed seconds to deploy or delete"?
// thresholds set by observing runs locally, rounding, multiplying by 1.25, and rounding:
var (
	deploySecondsForPackageCount = map[int]float64{50: 5, 100: 7, 500: 22.5}
	deleteSecondsForPackageCount = map[int]float64{50: 7, 100: 8, 500: 25}
)

func Benchmark_pkgr_with_500_packages(b *testing.B) {
	runWithPkgsAndVersions(b, 100, 5)
}

func Benchmark_pkgr_with_100_packages(b *testing.B) {
	runWithPkgsAndVersions(b, 50, 2)
}

func Benchmark_pkgr_with_50_packages(b *testing.B) {
	runWithPkgsAndVersions(b, 50, 1)
}

func runWithPkgsAndVersions(b *testing.B, numPackages int, numVersionsPerPackage int) {
	pkgrFileName := writePkgr(b, numPackages, numVersionsPerPackage)
	defer os.Remove(pkgrFileName)

	cleanup := func() {
		cmd := exec.Command("kapp", "delete", "-a", appName(pkgrFileName), "-y")
		cmd.Run()
	}
	cleanup()
	defer cleanup()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		deployAndDeletePkgr(b, pkgrFileName, numPackages*numVersionsPerPackage)
	}
}

// given a file name like "repo-500.yaml" returns "repo-500"
func appName(fileName string) string {
	return fileName[:len(fileName)-5]
}

func deployAndDeletePkgr(b *testing.B, pkgrFileName string, totalPackages int) {
	t1 := time.Now()
	cmd := exec.Command("kapp", "deploy", "-f", pkgrFileName, "-a", appName(pkgrFileName), "-y")
	output, err := cmd.Output()
	require.NoError(b, err, string(output))
	t2 := time.Now()

	cmd = exec.Command("kapp", "delete", "-a", appName(pkgrFileName), "-y")
	output, err = cmd.Output()
	require.NoError(b, err, string(output))
	t3 := time.Now()

	deployTime := t2.Sub(t1).Seconds()
	deleteTime := t3.Sub(t2).Seconds()

	assert.Less(b, deployTime, deploySecondsForPackageCount[totalPackages], "Seconds deploying were too slow for a pkgr with ", totalPackages, " packages.")
	assert.Less(b, deleteTime, deleteSecondsForPackageCount[totalPackages], "Seconds deleting were too slow for a pkgr with ", totalPackages, " packages.")
	b.ReportMetric(deployTime, "DeploySeconds")
	b.ReportMetric(deleteTime, "DeleteSeconds")
}

func writePkgr(b *testing.B, numPackages int, numVersions int) string {
	totalPackages := numPackages * numVersions

	preamble := fmt.Sprintf(`
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: minimal-repo-%d.tanzu.carvel.dev
  namespace: kapp-controller-packaging-global
  annotations:
    kapp.k14s.io/disable-original: ""
spec:
  fetch:
    inline:
      paths:
`, totalPackages)

	pkgStr := `
        packages/pkg.test.carvel.dev/pkg%[1]d.test.carvel.dev.0.%[2]d.0.yml: |
          ---
          apiVersion: data.packaging.carvel.dev/v1alpha1
          kind: Package
          metadata:
            name: pkg%[1]d.test.carvel.dev.0.%[2]d.0
          spec:
            refName: pkg%[1]d.test.carvel.dev
            version: 0.%[2]d.0
            template:
              spec: {}
`
	fname := fmt.Sprintf("pkgr-%d.yaml", totalPackages)
	f, err := os.Create(fname)
	require.NoError(b, err)
	defer f.Close()

	f.WriteString(preamble)
	for i := 0; i < numPackages; i++ {
		for j := 0; j < numVersions; j++ {
			_, err := f.WriteString(fmt.Sprintf(pkgStr, i, j))
			require.NoError(b, err)
		}
	}
	return fname
}
