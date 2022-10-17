// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sources

import (
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

const (
	LatestVersion = "latest"
)

type GithubReleaseConfiguration struct {
	ui           cmdcore.AuthoringUI
	vendirConfig *VendirConfig
}

func NewGithubReleaseConfiguration(ui cmdcore.AuthoringUI, vendirConfig *VendirConfig) *GithubReleaseConfiguration {
	return &GithubReleaseConfiguration{ui: ui, vendirConfig: vendirConfig}
}

func (g *GithubReleaseConfiguration) Configure() error {
	contents := g.vendirConfig.Contents()
	if contents == nil {
		err := g.initializeContentWithGithubRelease(contents)
		if err != nil {
			return err
		}
	} else if contents[0].GithubRelease == nil {
		err := g.initializeGithubRelease(contents)
		if err != nil {
			return err
		}
	}
	g.ui.PrintHeaderText("Repository details")

	err := g.configureRepoSlug(contents)
	if err != nil {
		return err
	}

	err = g.configureVersion(contents)
	if err != nil {
		return err
	}
	return g.getIncludedPaths(contents)
}

func (g *GithubReleaseConfiguration) initializeGithubRelease(contents []vendirconf.DirectoryContents) error {
	githubReleaseContent := vendirconf.DirectoryContentsGithubRelease{
		DisableAutoChecksumValidation: true,
	}
	contents[0].GithubRelease = &githubReleaseContent
	g.vendirConfig.SetContents(contents)
	return g.vendirConfig.Save()
}

func (g *GithubReleaseConfiguration) initializeContentWithGithubRelease(contents []vendirconf.DirectoryContents) error {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	g.vendirConfig.SetContents(append(contents, vendirconf.DirectoryContents{}))
	return g.initializeGithubRelease(contents)
}

func (g *GithubReleaseConfiguration) configureRepoSlug(contents []vendirconf.DirectoryContents) error {
	defaultSlug := contents[0].GithubRelease.Slug
	g.ui.PrintInformationalText("Slug format is org/repo e.g. vmware-tanzu/simple-app")
	textOpts := ui.TextOpts{
		Label:        "Enter slug for repository",
		Default:      defaultSlug,
		ValidateFunc: nil,
	}
	repoSlug, err := g.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	contents[0].GithubRelease.Slug = strings.TrimSpace(repoSlug)
	g.vendirConfig.SetContents(contents)
	return g.vendirConfig.Save()
}

func (g *GithubReleaseConfiguration) configureVersion(contents []vendirconf.DirectoryContents) error {
	githubReleaseContent := contents[0].GithubRelease
	defaultReleaseTag := g.getDefaultReleaseTag(contents)
	textOpts := ui.TextOpts{
		Label:        "Enter the release tag to be used",
		Default:      defaultReleaseTag,
		ValidateFunc: nil,
	}
	releaseTag, err := g.ui.AskForText(textOpts)
	if err != nil {
		return err
	}
	releaseTag = strings.TrimSpace(releaseTag)

	//TODO Rohit getting the releaseTag even though it is empty bcoz we dont have omitEmpty in the json representation. Might be have to create PR on imgpkg
	if releaseTag == LatestVersion {
		githubReleaseContent.Latest = true
		githubReleaseContent.Tag = ""
	} else {
		githubReleaseContent.Tag = releaseTag
		githubReleaseContent.Latest = false
	}

	contents[0].GithubRelease = githubReleaseContent
	g.vendirConfig.SetContents(contents)
	return g.vendirConfig.Save()
}

func (g *GithubReleaseConfiguration) getDefaultReleaseTag(contents []vendirconf.DirectoryContents) string {
	releaseTag := g.vendirConfig.Contents()[0].GithubRelease.Tag
	if len(releaseTag) > 0 {
		return releaseTag
	}
	return LatestVersion
}

func (g *GithubReleaseConfiguration) getIncludedPaths(contents []vendirconf.DirectoryContents) error {
	g.ui.PrintInformationalText("We need to know which files contain Kubernetes manifests. Multiple files can be included using a comma separator. To include all the files, enter *")
	includedPaths := contents[0].IncludePaths
	defaultIncludedPath := strings.Join(includedPaths, ",")
	if len(includedPaths) == 0 {
		defaultIncludedPath = includeAllFiles
	}
	textOpts := ui.TextOpts{
		Label:        "Enter the paths which contain Kubernetes manifests",
		Default:      defaultIncludedPath,
		ValidateFunc: nil,
	}
	path, err := g.ui.AskForText(textOpts)
	if err != nil {
		return err
	}
	paths := strings.Split(path, ",")
	if path == includeAllFiles {
		paths = nil
	}
	for i := 0; i < len(paths); i++ {
		paths[i] = strings.TrimSpace(paths[i])
	}
	contents[0].IncludePaths = paths
	g.vendirConfig.SetContents(contents)
	return g.vendirConfig.Save()
}
