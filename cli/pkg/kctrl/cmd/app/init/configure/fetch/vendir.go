// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"bytes"
	"fmt"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	"os"
	goexec "os/exec"
	"path/filepath"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"sigs.k8s.io/yaml"
)

const (
	VendirFileName      = "vendir.yml"
	VendirSyncDirectory = "upstream"
	IncludeAllFiles     = "*"
)

type VendirStep struct {
	ui       cmdcore.AuthoringUI
	appBuild *appbuild.AppBuild
	config   vendirconf.Config
}

func NewVendirStep(ui cmdcore.AuthoringUI, appBuild *appbuild.AppBuild, config vendirconf.Config) *VendirStep {
	vendirStep := VendirStep{
		ui:       ui,
		appBuild: appBuild,
		config:   config,
	}
	return &vendirStep
}

func (v *VendirStep) PreInteract() error {
	return nil
}

func (v *VendirStep) Interact() error {
	vendirDirectories := v.config.Directories
	if len(vendirDirectories) > 1 {
		//TODO what if we have >1 Directories section in the vendir conf

	}
	if len(vendirDirectories) == 0 {
		vendirDirectories = v.initializeVendirDirectorySection()
	} else {
		directory := vendirDirectories[0]
		if len(directory.Contents) > 1 {
			//TODO what needs to be done in this case
			//As multiple content sections are configured, we dont want to touch them.
			return nil
		}
	}
	currentFetchOptionSelected := v.appBuild.ObjectMeta.Annotations[common.FetchContentAnnotationKey]
	switch currentFetchOptionSelected {
	case common.FetchReleaseArtifactFromGithub:
		githubStep := NewGithubStep(v.ui, v.config)
		err := common.Run(githubStep)
		if err != nil {
			return err
		}
		includedPaths, err := v.getIncludedPaths()
		if err != nil {
			return err
		}
		v.config.Directories[0].Contents[0].IncludePaths = includedPaths

	case common.FetchChartFromHelmRepo:
	}

	SaveVendir(v.config)
	return nil
}

func (v *VendirStep) initializeVendirDirectorySection() []vendirconf.Directory {
	var directory vendirconf.Directory
	directory = vendirconf.Directory{
		Path: VendirSyncDirectory,
		Contents: []vendirconf.DirectoryContents{
			{
				Path: ".",
			},
		},
	}
	directories := []vendirconf.Directory{directory}
	v.config.Directories = directories
	SaveVendir(v.config)
	return directories
}

func (v *VendirStep) getIncludedPaths() ([]string, error) {
	v.ui.PrintInformationalText("We need exact manifest files from the above provided repository which should be included as package content. Multiple files can be included using a comma separator. If you want to include all the files, enter *.")
	includedPaths := v.config.Directories[0].Contents[0].IncludePaths
	defaultIncludedPath := strings.Join(includedPaths, ",")
	if len(includedPaths) == 0 {
		defaultIncludedPath = IncludeAllFiles
	}
	textOpts := ui.TextOpts{
		Label:        "Enter the paths which need to be included as part of this package",
		Default:      defaultIncludedPath,
		ValidateFunc: nil,
	}
	path, err := v.ui.AskForText(textOpts)
	path = strings.TrimSpace(path)
	if err != nil {
		return nil, err
	}
	if path == IncludeAllFiles {
		return []string{}, nil
	}
	paths := strings.Split(path, ",")
	return paths, nil
}

func (v *VendirStep) PostInteract() error {
	err := v.printVendirFile()
	if err != nil {
		return err
	}
	//err = v.runVendirSync()
	if err != nil {
		return err
	}
	return nil
}

func (v *VendirStep) printVendirFile() error {
	vendirFileLocation := filepath.Join(VendirFileName)
	v.ui.PrintActionableText(fmt.Sprintf("Printing %s", vendirFileLocation))
	v.ui.PrintCmdExecutionText(fmt.Sprintf("cat %s", vendirFileLocation))
	err := v.printFile(vendirFileLocation)
	if err != nil {
		return err
	}
	return nil
}

func (v *VendirStep) printFile(filePath string) error {
	var stdoutBs, stderrBs bytes.Buffer

	localCmdRunner := exec.NewPlainCmdRunner()
	cmd := goexec.Command("cat", []string{filePath}...)
	cmd.Stdin = nil
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs
	localCmdRunner.Run(cmd)

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	if result.Error != nil {
		return fmt.Errorf("Printing file: %s", result.Stderr)
	}

	v.ui.PrintCmdExecutionOutput(result.Stdout)

	return nil
}

func (v *VendirStep) runVendirSync() error {
	v.ui.PrintInformationalText("\nNext step is to run `vendir sync` to fetch the data from source to local.")
	v.ui.PrintActionableText("Running vendir sync")
	v.ui.PrintCmdExecutionText("vendir sync -f vendir.yml")
	var stdoutBs, stderrBs bytes.Buffer

	localCmdRunner := exec.NewPlainCmdRunner()
	cmd := goexec.Command("vendir", []string{"sync", "-f", VendirFileName}...)
	cmd.Stdin = nil
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs
	localCmdRunner.Run(cmd)

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	if result.Error != nil {
		return fmt.Errorf("Vendir sync failed. %s", result.Stderr)
	}

	v.ui.PrintInformationalText("\nTo validate that data has been fetched, lets list down the files")
	v.ui.PrintActionableText(fmt.Sprintf("Validating by listing files"))
	v.ui.PrintCmdExecutionText(fmt.Sprintf("ls -lR %s", VendirSyncDirectory))
	err := v.listFiles(VendirSyncDirectory)
	if err != nil {
		return err
	}
	return nil
}

func (v *VendirStep) listFiles(dir string) error {
	var stdoutBs, stderrBs bytes.Buffer

	localCmdRunner := exec.NewPlainCmdRunner()
	cmd := goexec.Command("ls", []string{"-lR", dir}...)
	cmd.Stdin = nil
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs
	localCmdRunner.Run(cmd)

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	if result.Error != nil {
		return fmt.Errorf("Listing files.\n %s", result.Stderr)
	}
	v.ui.PrintCmdExecutionOutput(result.Stdout)
	return nil
}

func ReadVendirConfig() (vendirconf.Config, error) {
	var vendirConfig vendirconf.Config
	exists, err := common.IsFileExists(VendirFileName)
	if err != nil {
		return vendirconf.Config{}, err
	}

	if exists {
		vendirConfig, err = VendirConfigFromExistingFile(VendirFileName)
		if err != nil {
			return vendirconf.Config{}, err
		}
	} else {
		vendirConfig = NewDefaultVendirConfig()
	}
	return vendirConfig, nil
}

func VendirConfigFromExistingFile(filePath string) (vendirconf.Config, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return vendirconf.Config{}, err
	}
	vendirConfig := vendirconf.Config{}
	err = yaml.Unmarshal(content, &vendirConfig)
	if err != nil {
		return vendirconf.Config{}, err
	}
	return vendirConfig, nil
}

func NewDefaultVendirConfig() vendirconf.Config {
	config := vendirconf.Config{
		APIVersion: "vendir.k14s.io/v1alpha1", // TODO: use constant from vendir package
		Kind:       "Config",                  // TODO: use constant from vendir package

	}
	SaveVendir(config)
	return config
}

func SaveVendir(config vendirconf.Config) error {
	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	common.WriteFile(VendirFileName, content)
	return nil
}
