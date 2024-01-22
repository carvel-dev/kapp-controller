// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

var (
	kappAllowedSharedOpts = []string{
		// Globals
		"--column",
		"--debug",
		"--json",
		"--tty",

		"--dangerous-ignore-failing-api-services",
		"--dangerous-scope-to-fallback-allowed-namespaces",

		// Filtering
		"--filter",
		"--filter-age",
		"--filter-kind",
		"--filter-kind-name",
		"--filter-kind-ns",
		"--filter-kind-ns-name",
		"--filter-name",
		"--filter-ns",

		"--kube-api-qps",
		"--kube-api-burst",
	}

	kappAllowedChangeOpts = []string{
		// Diffing
		"--diff-changes",
		"--diff-against-last-applied",
		"--diff-context",
		"--diff-line-numbers",
		"--diff-mask",
		"--diff-run",
		"--diff-summary",
		"--diff-anchored",

		// Applying
		"--apply-check-interval",
		"--apply-concurrency",
		"--apply-default-update-strategy",
		"--apply-ignored",
		"--apply-timeout",
		"--exit-early-on-apply-error",

		// Waiting
		"--wait",
		"--wait-check-interval",
		"--wait-concurrency",
		"--wait-ignored",
		"--wait-timeout",
		"--exit-early-on-wait-error",
	}
)

var (
	kappAllowedDeployFlagSet = exec.NewFlagSet(kappAllowedSharedOpts, kappAllowedChangeOpts, []string{
		"--dangerous-allow-empty-list-of-resources",

		"--existing-non-labeled-resources-check",
		"--existing-non-labeled-resources-check-concurrency",
		"--dangerous-override-ownership-of-existing-resources",

		"--into-ns",
		"--map-ns",

		"--logs",
		"--logs-all",

		"--app-changes-max-to-keep",

		"--labels",
		"--patch",
	})

	kappAllowedInspectFlagSet = exec.NewFlagSet(kappAllowedSharedOpts, []string{
		"--raw",
		"--status",
		"--tree",
	})

	kappAllowedDeleteFlagSet = exec.NewFlagSet(kappAllowedSharedOpts, kappAllowedChangeOpts)
)
