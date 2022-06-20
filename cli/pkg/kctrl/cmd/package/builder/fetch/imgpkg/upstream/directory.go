package upstream

import (
	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"strings"
)

type DirectoryStep struct {
	pkgAuthoringUI pkgui.IAuthoringUI
	pkgBuild       *pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewDirectoryStep(ui pkgui.IAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *DirectoryStep {
	return &DirectoryStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
}

func (directoryStep DirectoryStep) PreInteract() error {
	return nil
}

func (directoryStep DirectoryStep) Interact() error {
	contents := directoryStep.pkgBuild.Spec.Vendir.Directories[0].Contents
	if contents == nil {
		directoryStep.initializeContentWithLocalDirectory()
	} else if contents[0].Directory == nil {
		directoryStep.initializeLocalDirectory()
	}
	directoryStep.pkgAuthoringUI.PrintHeaderText("Directory path")

	err := directoryStep.configureFetchDirectoryPath()
	if err != nil {
		return err
	}

	// Here not taking target directory to store fetched content always defaulting it to config/upstream
	directoryStep.pkgAuthoringUI.PrintInformationalText("Fetched content will be always stored in config/upstream directory.\n")

	return nil
}

func (directoryStep DirectoryStep) PostInteract() error {
	return nil
}

func (directoryStep DirectoryStep) configureFetchDirectoryPath() error {
	directoryContent := directoryStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].Directory
	defaultPath := directoryContent.Path
	directoryStep.pkgAuthoringUI.PrintInformationalText("Directory path from where you want to fetch content (will be  relative to ./pkgBuild/bundle)")
	textOpts := ui.TextOpts{
		Label:        "Enter directory path from where to fetch content",
		Default:      defaultPath,
		ValidateFunc: nil,
	}
	path, err := directoryStep.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return err
	}

	directoryContent.Path = strings.TrimSpace(path)
	directoryStep.pkgBuild.WriteToFile()
	return nil
}

func (directoryStep DirectoryStep) initializeLocalDirectory() {
	directoryStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].Directory = &vendirconf.DirectoryContentsDirectory{}
	directoryStep.pkgBuild.WriteToFile()
}

func (directoryStep DirectoryStep) initializeContentWithLocalDirectory() {
	directoryStep.pkgBuild.Spec.Vendir.Directories[0].Contents = append(directoryStep.pkgBuild.Spec.Vendir.Directories[0].Contents, vendirconf.DirectoryContents{})
	directoryStep.initializeLocalDirectory()
}
