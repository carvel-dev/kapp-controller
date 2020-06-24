package deploy

import (
	"github.com/k14s/kapp-controller/pkg/exec"
)

var (
	kappAllowedSharedOpts = []string{
		// Globals
		"--column",
		"--debug",
		"--json",
		"--tty",

		"--dangerous-ignore-failing-api-services",

		// Filtering
		"--filter",
		"--filter-age",
		"--filter-kind",
		"--filter-kind-name",
		"--filter-kind-ns",
		"--filter-kind-ns-name",
		"--filter-name",
		"--filter-ns",
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

		// Applying
		"--apply-check-interval",
		"--apply-concurrency",
		"--apply-default-update-strategy",
		"--apply-ignored",
		"--apply-timeout",

		// Waiting
		"--wait",
		"--wait-check-interval",
		"--wait-concurrency",
		"--wait-ignored",
		"--wait-timeout",
	}
)

var (
	kappAllowedDeployFlagSet = exec.NewFlagSet(kappAllowedSharedOpts, kappAllowedChangeOpts, []string{
		"--dangerous-allow-empty-list-of-resources",
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
