package upstream

import (
	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

const (
	latestVersion = "latest"
)

type GithubStep struct {
	ui          ui.UI
	pkgBuild    *pkgbuilder.PackageBuild
	pkgLocation string
}

func NewGithubStep(ui ui.UI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *GithubStep {
	return &GithubStep{
		ui:          ui,
		pkgLocation: pkgLocation,
		pkgBuild:    pkgBuild,
	}
}

func (g *GithubStep) PreInteract() error {
	return nil
}

func (g *GithubStep) Interact() error {
	contents := g.pkgBuild.Spec.Vendir.Directories[0].Contents
	if contents == nil || contents[0].GithubRelease == nil {
		g.initializeContentWithGithubReleaseGithubRelease()
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
	_ = githubReleaseContent.Slug
	repoSlug, err := g.ui.AskForText("Enter slug for repository(org/repo)")
	if err != nil {
		return err
	}

	githubReleaseContent.Slug = repoSlug
	g.pkgBuild.WriteToFile(g.pkgLocation)
	return nil
}

func (g GithubStep) configureVersion() error {
	githubReleaseContent := g.pkgBuild.Spec.Vendir.Directories[0].Contents[0].GithubRelease
	_ = g.getDefaultReleaseTag()

	releaseTag, err := g.ui.AskForText("Enter the release tag to be used. Leave empty to use the latest version")
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

func (g *GithubStep) initializeContentWithGithubReleaseGithubRelease() {
	g.pkgBuild.Spec.Vendir.Directories[0].Contents = append(g.pkgBuild.Spec.Vendir.Directories[0].Contents, vendirconf.DirectoryContents{})
	g.initializeGithubRelease()
}
