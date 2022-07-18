package fetch

import (
	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"strings"
)

type GitStep struct {
	ui           cmdcore.AuthoringUI
	vendirConfig vendirconf.Config
}

func NewGitStep(ui cmdcore.AuthoringUI, vendirConfig vendirconf.Config) *GitStep {
	return &GitStep{
		ui:           ui,
		vendirConfig: vendirConfig,
	}
}

func (gitStep *GitStep) PreInteract() error {
	return nil
}

func (gitStep *GitStep) Interact() error {
	contents := gitStep.vendirConfig.Directories[0].Contents
	if contents == nil {
		err := gitStep.initializeContentWithGit()
		if err != nil {
			return err
		}
	} else if contents[0].Git == nil {
		err := gitStep.initializeGit()
		if err != nil {
			return err
		}
	}

	err := gitStep.configureGitURL()
	if err != nil {
		return err
	}
	return gitStep.configureGitRef()
	return nil
}

func (gitStep *GitStep) PostInteract() error {

	return nil
}

func (gitStep *GitStep) initializeContentWithGit() error {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	gitStep.vendirConfig.Directories[0].Contents = append(gitStep.vendirConfig.Directories[0].Contents, vendirconf.DirectoryContents{})
	return gitStep.initializeGit()
}

func (gitStep *GitStep) initializeGit() error {
	gitStep.vendirConfig.Directories[0].Contents[0].Git = &vendirconf.DirectoryContentsGit{}
	return SaveVendir(gitStep.vendirConfig)
}

func (gitStep *GitStep) configureGitURL() error {
	gitContent := gitStep.vendirConfig.Directories[0].Contents[0].Git
	defaultURL := gitContent.URL
	textOpts := ui.TextOpts{
		Label:        "Enter Git URL",
		Default:      defaultURL,
		ValidateFunc: nil,
	}
	name, err := gitStep.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	gitContent.URL = strings.TrimSpace(name)
	return SaveVendir(gitStep.vendirConfig)
}

func (gitStep *GitStep) configureGitRef() error {
	gitStep.ui.PrintInformationalText("A git reference can be any branch, tag, commit; origin is the name of the remote (required)")
	gitContent := gitStep.vendirConfig.Directories[0].Contents[0].Git
	defaultRef := gitContent.Ref
	if defaultRef == "" {
		defaultRef = "origin/develop"
	}
	textOpts := ui.TextOpts{
		Label:        "Enter Git Reference",
		Default:      defaultRef,
		ValidateFunc: nil,
	}
	name, err := gitStep.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	gitContent.Ref = strings.TrimSpace(name)
	return SaveVendir(gitStep.vendirConfig)
}
