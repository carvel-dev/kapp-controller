// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
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

func (g *GitStep) PreInteract() error { return nil }

func (g *GitStep) Interact() error {
	contents := g.vendirConfig.Directories[0].Contents
	if contents == nil {
		err := g.initializeContentWithGit()
		if err != nil {
			return err
		}
	} else if contents[0].Git == nil {
		err := g.initializeGit()
		if err != nil {
			return err
		}
	}

	err := g.configureGitURL()
	if err != nil {
		return err
	}
	err = g.configureGitRef()
	if err != nil {
		return err
	}
	return g.getIncludedPaths()
}

func (g *GitStep) PostInteract() error { return nil }

func (g *GitStep) initializeContentWithGit() error {
	//TODO Rohit need to check this how it should be done. It is giving path as empty.
	g.vendirConfig.Directories[0].Contents = append(g.vendirConfig.Directories[0].Contents, vendirconf.DirectoryContents{})
	return g.initializeGit()
}

func (g *GitStep) initializeGit() error {
	g.vendirConfig.Directories[0].Contents[0].Git = &vendirconf.DirectoryContentsGit{}
	return SaveVendir(g.vendirConfig)
}

func (g *GitStep) configureGitURL() error {
	g.ui.PrintInformationalText("Both https and ssh URL's are supported, e.g. https://github.com/vmware-tanzu/carvel-kapp-controller")
	gitContent := g.vendirConfig.Directories[0].Contents[0].Git
	defaultURL := gitContent.URL
	textOpts := ui.TextOpts{
		Label:        "Enter Git URL",
		Default:      defaultURL,
		ValidateFunc: nil,
	}
	name, err := g.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	gitContent.URL = strings.TrimSpace(name)
	return SaveVendir(g.vendirConfig)
}

func (g *GitStep) configureGitRef() error {
	g.ui.PrintInformationalText("A git reference can be any branch, tag, commit; origin is the name of the remote.")
	gitContent := g.vendirConfig.Directories[0].Contents[0].Git
	defaultRef := gitContent.Ref
	if defaultRef == "" {
		defaultRef = "origin/main"
	}
	textOpts := ui.TextOpts{
		Label:        "Enter Git Reference",
		Default:      defaultRef,
		ValidateFunc: nil,
	}
	name, err := g.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	gitContent.Ref = strings.TrimSpace(name)
	return SaveVendir(g.vendirConfig)
}

func (g *GitStep) getIncludedPaths() error {
	g.ui.PrintInformationalText(`We need to know which files contain Kubernetes manifests. Multiple files can be included using a comma separator. 
- To include all the files, enter * 
- To include a folder with all the sub-folders and files, enter <FOLDER_NAME>/**/*
- To include all the files inside a folder, enter <FOLDER_NAME>/*`)
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
