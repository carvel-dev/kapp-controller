package upstream

import (
	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"strings"
)

const (
	latestVersion = "latest"
)

type GithubStep struct {
	pkgAuthoringUI pkgui.IAuthoringUI
	pkgBuild       *pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewGithubStep(ui pkgui.IAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *GithubStep {
	return &GithubStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
}

func (githubStep *GithubStep) PreInteract() error {
	return nil
}

func (githubStep *GithubStep) Interact() error {
	contents := githubStep.pkgBuild.Spec.Vendir.Directories[0].Contents
	if contents == nil {
		githubStep.initializeContentWithGithubRelease()
	} else if contents[0].GithubRelease == nil {
		githubStep.initializeGithubRelease()
	}
	githubStep.pkgAuthoringUI.PrintHeaderText("Repository details")

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
	githubReleaseContent := githubStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].GithubRelease
	defaultSlug := githubReleaseContent.Slug
	githubStep.pkgAuthoringUI.PrintInformationalText("Slug format is org/repo e.g. vmware-tanzu/simple-app")
	textOpts := ui.TextOpts{
		Label:        "Enter slug for repository",
		Default:      defaultSlug,
		ValidateFunc: nil,
	}
	repoSlug, err := githubStep.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}

	githubReleaseContent.Slug = strings.TrimSpace(repoSlug)
	githubStep.pkgBuild.WriteToFile()
	return nil
}

func (githubStep GithubStep) configureVersion() error {
	githubReleaseContent := githubStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].GithubRelease
	defaultReleaseTag := githubStep.getDefaultReleaseTag()
	textOpts := ui.TextOpts{
		Label:        "Enter the release tag to be used. Leave empty to use the latest version",
		Default:      defaultReleaseTag,
		ValidateFunc: nil,
	}
	releaseTag, err := githubStep.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}
	releaseTag = strings.TrimSpace(releaseTag)

	//TODO Rohit getting the releaseTag even though it is empty bcoz we dont have omitEmpty in the json representation. Might be have to create PR on imgpkg
	if len(releaseTag) == 0 {
		githubReleaseContent.Latest = true
		githubReleaseContent.Tag = ""
	} else {
		githubReleaseContent.Tag = releaseTag
		githubReleaseContent.Latest = false
	}
	githubStep.pkgBuild.WriteToFile()
	return nil
}

func (githubStep *GithubStep) PostInteract() error {
	return nil
}

func (githubStep *GithubStep) initializeGithubRelease() {
	githubReleaseContent := vendirconf.DirectoryContentsGithubRelease{
		DisableAutoChecksumValidation: true,
	}
	githubStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].GithubRelease = &githubReleaseContent
	githubStep.pkgBuild.WriteToFile()
}

func (githubStep *GithubStep) getDefaultReleaseTag() string {
	releaseTag := githubStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].GithubRelease.Tag
	if len(releaseTag) > 0 {
		return releaseTag
	}
	return ""
}

func (githubStep *GithubStep) initializeContentWithGithubRelease() {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	githubStep.pkgBuild.Spec.Vendir.Directories[0].Contents = append(githubStep.pkgBuild.Spec.Vendir.Directories[0].Contents, vendirconf.DirectoryContents{})
	githubStep.initializeGithubRelease()
}
