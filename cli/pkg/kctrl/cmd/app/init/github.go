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
	vendirConfig *VendirConfig
}

func NewGithubStep(ui cmdcore.AuthoringUI, vendirConfig *VendirConfig) *GithubStep {
	return &GithubStep{
		ui:           ui,
		vendirConfig: vendirConfig,
	}
}

func (g *GithubStep) PreInteract() error { return nil }

func (g *GithubStep) Interact() error {
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

func (g *GithubStep) configureRepoSlug(contents []vendirconf.DirectoryContents) error {
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

func (g *GithubStep) configureVersion(contents []vendirconf.DirectoryContents) error {
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

func (g *GithubStep) PostInteract() error { return nil }

func (g *GithubStep) initializeGithubRelease(contents []vendirconf.DirectoryContents) error {
	githubReleaseContent := vendirconf.DirectoryContentsGithubRelease{
		DisableAutoChecksumValidation: true,
	}
	contents[0].GithubRelease = &githubReleaseContent
	g.vendirConfig.SetContents(contents)
	return g.vendirConfig.Save()
}

func (g *GithubStep) getDefaultReleaseTag(contents []vendirconf.DirectoryContents) string {
	releaseTag := g.vendirConfig.Contents()[0].GithubRelease.Tag
	if len(releaseTag) > 0 {
		return releaseTag
	}
	return LatestVersion
}

func (g *GithubStep) initializeContentWithGithubRelease(contents []vendirconf.DirectoryContents) error {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	g.vendirConfig.SetContents(append(contents, vendirconf.DirectoryContents{}))
	return g.initializeGithubRelease(contents)
}

func (g *GithubStep) getIncludedPaths(contents []vendirconf.DirectoryContents) error {
	g.ui.PrintInformationalText("We need to know which files contain Kubernetes manifests. Multiple files can be included using a comma separator. To include all the files, enter *")
	includedPaths := contents[0].IncludePaths
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
	contents[0].IncludePaths = paths
	g.vendirConfig.SetContents(contents)
	return g.vendirConfig.Save()
}
