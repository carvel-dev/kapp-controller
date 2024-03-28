// Copyright 2024 The Carvel Authors.
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

type GithubReleaseSource struct {
	ui           cmdcore.AuthoringUI
	vendirConfig *VendirConfig
}

var _ SourceConfiguration = &GithubReleaseSource{}

func NewGithubReleaseSource(ui cmdcore.AuthoringUI, vendirConfig *VendirConfig) *GithubReleaseSource {
	return &GithubReleaseSource{ui: ui, vendirConfig: vendirConfig}
}

func (g *GithubReleaseSource) Configure() error {
	err := g.initializeContentWithGithubRelease()
	if err != nil {
		return err
	}

	g.ui.PrintHeaderText("Repository details")

	err = g.configureRepoSlug()
	if err != nil {
		return err
	}

	err = g.configureVersion()
	if err != nil {
		return err
	}
	return g.getIncludedPaths()
}

func (g *GithubReleaseSource) initializeContentWithGithubRelease() error {
	contents := g.vendirConfig.Contents()
	if contents == nil {
		contents = append(contents, vendirconf.DirectoryContents{})
	}
	if contents[0].GithubRelease == nil {
		contents[0].GithubRelease = &vendirconf.DirectoryContentsGithubRelease{DisableAutoChecksumValidation: true}
	}
	g.vendirConfig.SetContents(contents)
	return g.vendirConfig.Save()
}

func (g *GithubReleaseSource) configureRepoSlug() error {
	contents := g.vendirConfig.Contents()
	g.ui.PrintInformationalText("Slug format is org/repo e.g. vmware-tanzu/simple-app")
	textOpts := ui.TextOpts{
		Label:        "Enter slug for repository",
		Default:      contents[0].GithubRelease.Slug,
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

func (g *GithubReleaseSource) configureVersion() error {
	contents := g.vendirConfig.Contents()
	githubReleaseContent := contents[0].GithubRelease
	textOpts := ui.TextOpts{
		Label:        "Enter the release tag to be used",
		Default:      g.getDefaultReleaseTag(contents),
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

func (g *GithubReleaseSource) getDefaultReleaseTag(contents []vendirconf.DirectoryContents) string {
	releaseTag := g.vendirConfig.Contents()[0].GithubRelease.Tag
	if len(releaseTag) > 0 {
		return releaseTag
	}
	return LatestVersion
}

func (g *GithubReleaseSource) getIncludedPaths() error {
	contents := g.vendirConfig.Contents()
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
