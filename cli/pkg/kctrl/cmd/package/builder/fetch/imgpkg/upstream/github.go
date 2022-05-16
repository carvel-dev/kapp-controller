package upstream

import (
	"github.com/cppforlife/go-cli-ui/ui"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

const (
	latestVersion = "latest"
)

type GithubStep struct {
	ui            ui.UI
	GithubRelease *vendirconf.DirectoryContentsGithubRelease `json:"githubRelease,omitempty"`
}

func NewGithubStep(ui ui.UI) *GithubStep {
	return &GithubStep{
		ui: ui,
	}
}

func (g *GithubStep) PreInteract() error {
	return nil
}

func (g *GithubStep) PostInteract() error {
	return nil
}

func (g *GithubStep) Interact() error {
	repoSlug, err := g.ui.AskForText("Enter slug for repository(org/repo)")
	if err != nil {
		return err
	}
	var releaseTag string
	var latest bool

	releaseVersion, err := g.getVersion()
	if releaseVersion == latestVersion {
		latest = true
	} else {
		releaseTag = releaseVersion
	}
	if err != nil {
		return err
	}
	//TODO Rohit getting the releaseTag even though it is empty bcoz we dont have omitEmpty in the json representation. Might be have to create PR on imgpkg
	directoryContentsGithubRelease := vendirconf.DirectoryContentsGithubRelease{
		Slug:                          repoSlug,
		Tag:                           releaseTag,
		Latest:                        latest,
		DisableAutoChecksumValidation: true,
	}
	g.GithubRelease = &directoryContentsGithubRelease
	return nil
}

func (g GithubStep) getVersion() (string, error) {
	//TODO Rohit check when you press ctrl-C, does it generate an error
	releaseTag, err := g.ui.AskForText("Enter the release tag to be used. Leave empty to use the latest version")
	if err != nil {
		return "", err
	}
	if len(releaseTag) == 0 {
		return latestVersion, nil
	}
	return releaseTag, nil
}
