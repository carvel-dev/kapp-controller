// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

const (
	FetchContentAnnotationKey = "fetch-content-from"
	LocalFetchAnnotationKey   = "kctrl.carvel.dev/local-fetch-0"
)

const (
	FetchFromGithubRelease  string = "Github Release"
	FetchManifestFromGit    string = "Git Repository"
	FetchChartFromHelmRepo  string = "Helm Chart from Helm Repository"
	FetchChartFromGit       string = "Helm Chart from Git Repository"
	FetchFromLocalDirectory string = "Local Directory"
)
