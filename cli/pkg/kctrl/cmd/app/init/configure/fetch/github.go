// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

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

func (githubStep *GithubStep) PreInteract() error {
	return nil
}

func (githubStep *GithubStep) Interact() error {
	contents := githubStep.vendirConfig.Directories[0].Contents
	if contents == nil {
		githubStep.initializeContentWithGithubRelease()
	} else if contents[0].GithubRelease == nil {
		githubStep.initializeGithubRelease()
	}
	githubStep.ui.PrintHeaderText("Repository details")

	err := githubStep.configureRepoSlug()
	if err != nil {
		return err
	}

	err = githubStep.configureVersion()
	if err != nil {
		return err
	}

	return nil
}

func (githubStep *GithubStep) configureRepoSlug() error {
	githubReleaseContent := githubStep.vendirConfig.Directories[0].Contents[0].GithubRelease
	defaultSlug := githubReleaseContent.Slug
	githubStep.ui.PrintInformationalText("Slug format is org/repo e.g. vmware-tanzu/simple-app")
	textOpts := ui.TextOpts{
		Label:        "Enter slug for repository",
		Default:      defaultSlug,
		ValidateFunc: nil,
	}
	repoSlug, err := githubStep.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	githubReleaseContent.Slug = strings.TrimSpace(repoSlug)
	SaveVendir(githubStep.vendirConfig)
	return nil
}

func (githubStep *GithubStep) configureVersion() error {
	githubReleaseContent := githubStep.vendirConfig.Directories[0].Contents[0].GithubRelease
	defaultReleaseTag := githubStep.getDefaultReleaseTag()
	textOpts := ui.TextOpts{
		Label:        "Enter the release tag to be used",
		Default:      defaultReleaseTag,
		ValidateFunc: nil,
	}
	releaseTag, err := githubStep.ui.AskForText(textOpts)
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
	SaveVendir(githubStep.vendirConfig)
	return nil
}

func (githubStep *GithubStep) PostInteract() error {
	return nil
}

func (githubStep *GithubStep) initializeGithubRelease() {
	githubReleaseContent := vendirconf.DirectoryContentsGithubRelease{
		DisableAutoChecksumValidation: true,
	}
	githubStep.vendirConfig.Directories[0].Contents[0].GithubRelease = &githubReleaseContent
	SaveVendir(githubStep.vendirConfig)
}

func (githubStep *GithubStep) getDefaultReleaseTag() string {
	releaseTag := githubStep.vendirConfig.Directories[0].Contents[0].GithubRelease.Tag
	if len(releaseTag) > 0 {
		return releaseTag
	}
	return LatestVersion
}

func (githubStep *GithubStep) initializeContentWithGithubRelease() {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	githubStep.vendirConfig.Directories[0].Contents = append(githubStep.vendirConfig.Directories[0].Contents, vendirconf.DirectoryContents{})
	githubStep.initializeGithubRelease()
}
