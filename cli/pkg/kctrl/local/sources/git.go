// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sources

import (
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

type SourceConfiguration interface {
	Configure() error
}

type GitSource struct {
	ui           cmdcore.AuthoringUI
	vendirConfig *VendirConfig
}

var _ SourceConfiguration = &GitSource{}

func NewGitSource(ui cmdcore.AuthoringUI, vendirConfig *VendirConfig) *GitSource {
	return &GitSource{ui: ui, vendirConfig: vendirConfig}
}

func (g *GitSource) Configure() error {
	contents := g.vendirConfig.Contents()
	if contents == nil {
		err := g.initializeContentWithGit(contents)
		if err != nil {
			return err
		}
	} else if contents[0].Git == nil {
		err := g.initializeGit(contents)
		if err != nil {
			return err
		}
	}

	err := g.configureGitURL(contents)
	if err != nil {
		return err
	}
	err = g.configureGitRef(contents)
	if err != nil {
		return err
	}
	return g.getIncludedPaths(contents)
}

func (g *GitSource) initializeContentWithGit(contents []vendirconf.DirectoryContents) error {
	g.vendirConfig.SetContents(append(contents, vendirconf.DirectoryContents{}))
	return g.initializeGit(contents)
}

func (g *GitSource) initializeGit(contents []vendirconf.DirectoryContents) error {
	contents[0].Git = &vendirconf.DirectoryContentsGit{}
	g.vendirConfig.SetContents(contents)
	return g.vendirConfig.Save()
}

func (g *GitSource) configureGitURL(contents []vendirconf.DirectoryContents) error {
	g.ui.PrintInformationalText("Both https and ssh URL's are supported, e.g. https://github.com/vmware-tanzu/carvel-kapp-controller")
	defaultURL := contents[0].Git.URL
	textOpts := ui.TextOpts{
		Label:        "Enter Git URL",
		Default:      defaultURL,
		ValidateFunc: nil,
	}
	name, err := g.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	contents[0].Git.URL = strings.TrimSpace(name)
	g.vendirConfig.SetContents(contents)
	return g.vendirConfig.Save()
}

func (g *GitSource) configureGitRef(contents []vendirconf.DirectoryContents) error {
	g.ui.PrintInformationalText("A git reference can be any branch, tag, commit; origin is the name of the remote.")
	defaultRef := contents[0].Git.Ref
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

	contents[0].Git.Ref = strings.TrimSpace(name)
	g.vendirConfig.SetContents(contents)
	return g.vendirConfig.Save()
}

func (g *GitSource) getIncludedPaths(contents []vendirconf.DirectoryContents) error {
	g.ui.PrintInformationalText(`We need to know which files contain Kubernetes manifests. Multiple files can be included using a comma separator. 
- To include all the files, enter * 
- To include a folder with all the sub-folders and files, enter <FOLDER_NAME>/**/*
- To include all the files inside a folder, enter <FOLDER_NAME>/*`)
	includedPaths := contents[0].IncludePaths
	defaultIncludedPath := strings.Join(includedPaths, ",")
	if len(includedPaths) == 0 {
		defaultIncludedPath = includeAllFiles
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
	if path == includeAllFiles {
		paths = nil
	}
	for i := 0; i < len(paths); i++ {
		paths[i] = strings.TrimSpace(paths[i])
	}
	contents[0].IncludePaths = paths
	g.vendirConfig.SetContents(contents)
	return g.vendirConfig.Save()
}
