// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	goexec "os/exec"

	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

const (
	VendirFileName      = "vendir.yml"
	VendirSyncDirectory = "upstream"
	IncludeAllFiles     = "*"
)

type VendirStep struct {
	ui          cmdcore.AuthoringUI
	config      vendirconf.Config
	fetchOption string
}

func NewVendirStep(ui cmdcore.AuthoringUI, config vendirconf.Config, fetchOption string) *VendirStep {
	vendirStep := VendirStep{
		ui:          ui,
		config:      config,
		fetchOption: fetchOption,
	}
	return &vendirStep
}

func (v *VendirStep) PreInteract() error { return nil }

func (v *VendirStep) Interact() error {
	vendirDirectories := v.config.Directories
	if len(vendirDirectories) > 1 {
		return fmt.Errorf("More than 1 directory config found in the vendir file. (hint: Run vendir sync manually)")

	}
	if len(vendirDirectories) == 0 {
		vendirDirectories = v.initializeVendirDirectorySection()
	} else {
		directory := vendirDirectories[0]
		if len(directory.Contents) > 1 {
			return fmt.Errorf("More than 1 content config found in the vendir file. (hint: Run vendir sync manually)")
		}
	}
	currentFetchOptionSelected := v.fetchOption
	switch currentFetchOptionSelected {
	case FetchFromGithubRelease:
		githubStep := NewGithubStep(v.ui, v.config)
		return Run(githubStep)
	case FetchFromHelmRepo:
		helmStep := NewHelmStep(v.ui, v.config)
		return Run(helmStep)
	case FetchFromGit, FetchChartFromGit:
		gitStep := NewGitStep(v.ui, v.config)
		return Run(gitStep)
	}
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

func (v *VendirStep) PostInteract() error {
	v.ui.PrintInformationalText("We will use vendir to fetch the data from the source to the local directory. Vendir allows us to declaratively state what should be in a directory and sync data sources into it. All the information entered above has been persisted into a vendir.yml file.")
	err := v.printVendirFile()
	if err != nil {
		return err
	}
	err = v.runVendirSync()
	if err != nil {
		return err
	}
	return nil
}

func (v *VendirStep) printVendirFile() error {
	vendirFileLocation := VendirFileName
	v.ui.PrintActionableText(fmt.Sprintf("Printing %s \n", vendirFileLocation))
	err := v.printFile(vendirFileLocation)
	if err != nil {
		return err
	}
	return nil
}

func (v *VendirStep) printFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Printing file: %w", err)
	}
	defer func() {
		file.Close()
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		v.ui.PrintCmdExecutionOutput(scanner.Text())
	}
	return nil
}

// TODO vendir sync failure. Reproduce: In case of 429 from github, we dont show errors today.
func (v *VendirStep) runVendirSync() error {
	v.ui.PrintInformationalText("\nNext step is to run `vendir sync` to fetch the data from the source to the local directory. Vendir will sync the data into the upstream folder.")
	v.ui.PrintActionableText("Running vendir sync")
	v.ui.PrintCmdExecutionText("vendir sync -f vendir.yml\n")
	var stdoutBs, stderrBs bytes.Buffer

	localCmdRunner := exec.NewPlainCmdRunner()
	cmd := goexec.Command("vendir", []string{"sync", "-f", VendirFileName}...)
	cmd.Stdin = nil
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs
	err := localCmdRunner.Run(cmd)
	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Fetching resources: %s", err)
	if result.Error != nil {
		return fmt.Errorf("Vendir sync failed. %s", result.Stderr)
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
