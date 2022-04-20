package upstream

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
)

type GithubStep struct {
	RepoSlug                      string `json:"slug"`
	ReleaseTag                    string `json:"tag"`
	Ui                            ui.UI  `json:"-"`
	DisableAutoChecksumValidation bool   `json:"disableAutoChecksumValidation"`
}

func NewGithubStep(ui ui.UI) *GithubStep {
	return &GithubStep{
		Ui: ui,
	}
}

func (g *GithubStep) PreInteract() error {
	return nil
}

func (g *GithubStep) PostInteract() error {
	return nil
}

func (g *GithubStep) Interact() error {
	repoSlug, err := g.Ui.AskForText("Enter slug for repository(org/repo)")
	if err != nil {
		return err
	}
	g.RepoSlug = repoSlug
	releaseTag, err := g.getVersion()
	if err != nil {
		return err
	}
	g.ReleaseTag = releaseTag
	g.DisableAutoChecksumValidation = true
	return nil
}

func (g GithubStep) getVersion() (string, error) {
	var useLatestVersion bool
	for {
		input, err := g.Ui.AskForText("Do you want to use the latest released version(y/n)")
		if err != nil {
			return "", err
		}
		var isValidInput bool
		useLatestVersion, isValidInput = common.ValidateInputYesOrNo(input)
		if isValidInput {
			break
		} else {
			input, _ = g.Ui.AskForText("Invalid input. (must be 'y','n','Y','N')")
		}
	}

	//if useLatestVersion {
	if useLatestVersion {

	} else {
		g.Ui.PrintBlock([]byte("# Ok. Then we have to mention the specific release tag which makes up the package configuration"))
		releaseTag, err := g.Ui.AskForText("Enter the release tag")
		if err != nil {
			return "", err
		}
		return releaseTag, nil
	}
	return "", nil
}

func (g *GithubStep) Run() error {
	g.PreInteract()
	g.Interact()
	g.PostInteract()
	return nil
}
