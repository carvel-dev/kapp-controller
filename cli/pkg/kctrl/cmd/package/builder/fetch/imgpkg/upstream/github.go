package upstream

import (
	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

const (
	latestVersion = "latest"
)

type GithubStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	pkgBuild       *pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewGithubStep(ui pkgui.IPkgAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *GithubStep {
	return &GithubStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
}

func (g *GithubStep) PreInteract() error {
	return nil
}

func (g *GithubStep) Interact() error {
	contents := g.pkgBuild.Spec.Vendir.Directories[0].Contents
	if contents == nil {
		g.initializeContentWithGithubRelease()
	} else if contents[0].GithubRelease == nil {
		g.initializeGithubRelease()
	}

	err := g.configureRepoSlug()
	if err != nil {
		return err
	}

	err = g.configureVersion()
	if err != nil {
		return err
	}

	return nil
}

func (g *GithubStep) configureRepoSlug() error {
	githubReleaseContent := g.pkgBuild.Spec.Vendir.Directories[0].Contents[0].GithubRelease
	defaultSlug := githubReleaseContent.Slug
	textOpts := ui.TextOpts{
		Label:        "Enter slug for repository(org/repo)",
		Default:      defaultSlug,
		ValidateFunc: nil,
	}
	repoSlug, err := g.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}

	githubReleaseContent.Slug = repoSlug
	g.pkgBuild.WriteToFile(g.pkgLocation)
	return nil
}

func (g GithubStep) configureVersion() error {
	githubReleaseContent := g.pkgBuild.Spec.Vendir.Directories[0].Contents[0].GithubRelease
	defaultReleaseTag := g.getDefaultReleaseTag()
	textOpts := ui.TextOpts{
		Label:        "Enter the release tag to be used. Leave empty to use the latest version",
		Default:      defaultReleaseTag,
		ValidateFunc: nil,
	}

	releaseTag, err := g.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}
	//TODO Rohit getting the releaseTag even though it is empty bcoz we dont have omitEmpty in the json representation. Might be have to create PR on imgpkg
	if len(releaseTag) == 0 {
		githubReleaseContent.Latest = true
		githubReleaseContent.Tag = ""
	} else {
		githubReleaseContent.Tag = releaseTag
		githubReleaseContent.Latest = false
	}

	g.pkgBuild.WriteToFile(g.pkgLocation)
	return nil
}

func (g *GithubStep) PostInteract() error {
	return nil
}

func (g *GithubStep) initializeGithubRelease() {
	githubReleaseContent := vendirconf.DirectoryContentsGithubRelease{
		DisableAutoChecksumValidation: true,
	}
	g.pkgBuild.Spec.Vendir.Directories[0].Contents[0].GithubRelease = &githubReleaseContent
	g.pkgBuild.WriteToFile(g.pkgLocation)
}

func (g *GithubStep) getDefaultReleaseTag() string {
	releaseTag := g.pkgBuild.Spec.Vendir.Directories[0].Contents[0].GithubRelease.Tag
	if len(releaseTag) > 0 {
		return releaseTag
	}
	return ""
}

func (g *GithubStep) initializeContentWithGithubRelease() {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	g.pkgBuild.Spec.Vendir.Directories[0].Contents = append(g.pkgBuild.Spec.Vendir.Directories[0].Contents, vendirconf.DirectoryContents{})
	g.initializeGithubRelease()
}
