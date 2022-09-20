// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

const (
	LatestVersion = "latest"
)

type GithubStep struct {
	ui           cmdcore.AuthoringUI
	vendirConfig vendirconf.Config
}

func NewGithubStep(ui cmdcore.AuthoringUI, vendirConfig vendirconf.Config) *GithubStep {
	return &GithubStep{
		ui:           ui,
		vendirConfig: vendirConfig,
	}
}

func (g *GithubStep) PreInteract() error { return nil }

func (g *GithubStep) Interact() error {
	contents := g.vendirConfig.Directories[0].Contents
	if contents == nil {
		err := g.initializeContentWithGithubRelease()
		if err != nil {
			return err
		}
	} else if contents[0].GithubRelease == nil {
		err := g.initializeGithubRelease()
		if err != nil {
			return err
		}
	}
	g.ui.PrintHeaderText("Repository details")

	err := g.configureRepoSlug()
	if err != nil {
		return err
	}

	err = g.configureVersion()
	if err != nil {
		return err
	}
	return g.getIncludedPaths()
}

func (g *GithubStep) configureRepoSlug() error {
	githubReleaseContent := g.vendirConfig.Directories[0].Contents[0].GithubRelease
	defaultSlug := githubReleaseContent.Slug
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

	githubReleaseContent.Slug = strings.TrimSpace(repoSlug)
	return SaveVendir(g.vendirConfig)
}

func (g *GithubStep) configureVersion() error {
	githubReleaseContent := g.vendirConfig.Directories[0].Contents[0].GithubRelease
	defaultReleaseTag := g.getDefaultReleaseTag()
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
	return SaveVendir(g.vendirConfig)
}

func (g *GithubStep) PostInteract() error { return nil }

func (g *GithubStep) initializeGithubRelease() error {
	githubReleaseContent := vendirconf.DirectoryContentsGithubRelease{
		DisableAutoChecksumValidation: true,
	}
	g.vendirConfig.Directories[0].Contents[0].GithubRelease = &githubReleaseContent
	return SaveVendir(g.vendirConfig)
}

func (g *GithubStep) getDefaultReleaseTag() string {
	releaseTag := g.vendirConfig.Directories[0].Contents[0].GithubRelease.Tag
	if len(releaseTag) > 0 {
		return releaseTag
	}
	return LatestVersion
}

func (g *GithubStep) initializeContentWithGithubRelease() error {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	g.vendirConfig.Directories[0].Contents = append(g.vendirConfig.Directories[0].Contents, vendirconf.DirectoryContents{})
	return g.initializeGithubRelease()
}

func (g *GithubStep) getIncludedPaths() error {
	g.ui.PrintInformationalText("We need to know which files contain Kubernetes manifests. Multiple files can be included using a comma separator. To include all the files, enter *")
	includedPaths := g.vendirConfig.Directories[0].Contents[0].IncludePaths
	defaultIncludedPath := strings.Join(includedPaths, ",")
	if len(includedPaths) == 0 {
		defaultIncludedPath = IncludeAllFiles
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
	if path == IncludeAllFiles {
		paths = nil
	}
	for i := 0; i < len(paths); i++ {
		paths[i] = strings.TrimSpace(paths[i])
	}
	g.vendirConfig.Directories[0].Contents[0].IncludePaths = paths
	return SaveVendir(g.vendirConfig)
}
