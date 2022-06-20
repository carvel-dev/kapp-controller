package upstream

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"sigs.k8s.io/yaml"
)

const (
	VendirGithubReleaseConf string = "Github Release"
	VendirHelmChartConf     string = "HelmChart"
	VendirDirectoryConf     string = "Directory"
	IncludeAllFiles         string = "*"
)

type UpstreamStep struct {
	pkgAuthoringUI pkgui.IAuthoringUI
	PkgLocation    string
	pkgBuild       *pkgbuilder.PackageBuild
}

func NewUpstreamStep(ui pkgui.IAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *UpstreamStep {
	return &UpstreamStep{
		pkgAuthoringUI: ui,
		PkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
}

func (upstreamStep *UpstreamStep) PreInteract() error {
	return nil
}

func (upstreamStep *UpstreamStep) Interact() error {
	vendirDirectories := upstreamStep.pkgBuild.Spec.Vendir.Directories
	currentUpstreamOptionSelected := getCurrentSelectedUpstreamOption(upstreamStep.pkgBuild)
	upstreamOptionToVendirContentMap := map[string]string{
		common.FetchReleaseArtifactFromGithub: VendirGithubReleaseConf,
		common.FetchChartFromHelmRepo:         VendirHelmChartConf,
		common.FetchFromLocalDirectory:        VendirDirectoryConf,
	}

	if len(vendirDirectories) > 1 {
		//As multiple upstream directories are configured, we dont want to touch them.
		return nil
	}
	if len(vendirDirectories) == 0 {
		vendirDirectories = upstreamStep.initializeVendirDirectoryConf()
	} else {
		directory := vendirDirectories[0]
		if len(directory.Contents) > 1 {
			//As multiple content sections are configured, we dont want to touch them.
			return nil
		}
		var previousUpstreamOptionSelected string
		previousUpstreamOptionSelected = getPreviousSelectedUpstreamOption(upstreamStep.pkgBuild)
		if currentUpstreamOptionSelected != previousUpstreamOptionSelected {
			setEarlierUpstreamOptionAsNil(vendirDirectories, upstreamOptionToVendirContentMap[previousUpstreamOptionSelected])
			resetIncludedPathAsNil(vendirDirectories)
		}
	}

	switch currentUpstreamOptionSelected {
	case common.FetchReleaseArtifactFromGithub:
		githubStep := NewGithubStep(upstreamStep.pkgAuthoringUI, upstreamStep.PkgLocation, upstreamStep.pkgBuild)
		err := common.Run(githubStep)
		if err != nil {
			return err
		}
		includedPaths, err := upstreamStep.getIncludedPaths()
		if err != nil {
			return err
		}
		upstreamStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].IncludePaths = includedPaths
		upstreamStep.pkgBuild.WriteToFile()
	case common.FetchChartFromHelmRepo:
		helmStep := NewHelmStep(upstreamStep.pkgAuthoringUI, upstreamStep.PkgLocation, upstreamStep.pkgBuild)
		err := common.Run(helmStep)
		if err != nil {
			return err
		}
	case common.FetchFromLocalDirectory:
		directoryStep := NewDirectoryStep(upstreamStep.pkgAuthoringUI, upstreamStep.PkgLocation, upstreamStep.pkgBuild)
		err := common.Run(directoryStep)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCurrentSelectedUpstreamOption(pkgBuild *pkgbuilder.PackageBuild) string {
	return pkgBuild.Annotations[common.PkgFetchContentAnnotationKey]
}

func setEarlierUpstreamOptionAsNil(vendirDirectories []vendirconf.Directory, earlierUpstreamOption string) {
	if vendirDirectories[0].Contents == nil {
		return
	}
	switch earlierUpstreamOption {
	case VendirGithubReleaseConf:
		vendirDirectories[0].Contents[0].GithubRelease = nil
	case VendirHelmChartConf:
		vendirDirectories[0].Contents[0].HelmChart = nil
	case VendirDirectoryConf:
		vendirDirectories[0].Contents[0].Directory = nil
	}
	return
}

func resetIncludedPathAsNil(vendirDirectories []vendirconf.Directory) {
	if vendirDirectories[0].Contents == nil {
		return
	}
	vendirDirectories[0].Contents[0].IncludePaths = nil

}

func getPreviousSelectedUpstreamOption(pkgBuild *pkgbuilder.PackageBuild) string {
	dirContents := pkgBuild.Spec.Vendir.Directories[0].Contents
	if dirContents == nil {
		return ""
	}
	content := pkgBuild.Spec.Vendir.Directories[0].Contents[0]
	var selectedUpstreamOption string
	switch {
	case content.GithubRelease != nil:
		selectedUpstreamOption = common.FetchReleaseArtifactFromGithub
	case content.HelmChart != nil:
		selectedUpstreamOption = common.FetchChartFromHelmRepo
	case content.Directory != nil:
		selectedUpstreamOption = common.FetchFromLocalDirectory
	}
	return selectedUpstreamOption
}

func (upstreamStep *UpstreamStep) initializeVendirDirectoryConf() []vendirconf.Directory {
	var directory vendirconf.Directory
	directory = vendirconf.Directory{
		Path: "config/upstream",
		Contents: []vendirconf.DirectoryContents{
			{
				Path: ".",
			},
		},
	}
	directories := []vendirconf.Directory{directory}
	upstreamStep.pkgBuild.Spec.Vendir.Directories = directories
	upstreamStep.pkgBuild.WriteToFile()
	return directories
}

func (upstreamStep *UpstreamStep) PostInteract() error {
	err := upstreamStep.createVendirFile()
	if err != nil {
		return err
	}
	err = upstreamStep.printVendirFile()
	if err != nil {
		return err
	}
	err = upstreamStep.runVendirSync()
	if err != nil {
		return err
	}
	return nil
}

func (upstreamStep *UpstreamStep) createVendirFile() error {
	vendirFileLocation := filepath.Join(upstreamStep.PkgLocation, "bundle", "vendir.yml")
	upstreamStep.pkgAuthoringUI.PrintInformationalText("We have all the information needed to fetch the data from source to local. We will use vendir for this purpose. Vendir allows to declaratively state what should be in a directory and sync any number of data sources into it. Lets create vendir.yml file. Vendir.yml is used to fetch the data.This is how `vendir.yml` file looks.")
	data, err := yaml.Marshal(&upstreamStep.pkgBuild.Spec.Vendir)
	if err != nil {
		return fmt.Errorf("Unable to create vendir.yml\n %s", err.Error())
	}
	f, err := os.Create(vendirFileLocation)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (upstreamStep *UpstreamStep) printVendirFile() error {
	upstreamStep.pkgAuthoringUI.PrintActionableText("Printing vendir.yml")
	vendirFileLocation := filepath.Join(upstreamStep.PkgLocation, "bundle", "vendir.yml")
	upstreamStep.pkgAuthoringUI.PrintCmdExecutionText("cat vendir.yml")
	err := upstreamStep.printFile(vendirFileLocation)
	if err != nil {
		return err
	}
	return nil
}

func (upstreamStep *UpstreamStep) printFile(filePath string) error {
	result := util.Execute("cat", []string{filePath})
	if result.Error != nil {
		return fmt.Errorf("Printing file %s\n %s", filePath, result.Stderr)
	}
	upstreamStep.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	return nil
}

func (upstreamStep *UpstreamStep) runVendirSync() error {
	bundleLocation := filepath.Join(upstreamStep.PkgLocation, "bundle")
	upstreamStep.pkgAuthoringUI.PrintInformationalText("\nNext step is to run `vendir sync` to fetch the data from source to local.")
	upstreamStep.pkgAuthoringUI.PrintActionableText("Running vendir sync")
	upstreamStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("vendir sync --chdir %s -f vendir.yml", bundleLocation))
	result := util.Execute("vendir", []string{"sync", "--chdir", bundleLocation, "-f", "vendir.yml"})
	if result.Error != nil {
		return fmt.Errorf("while running vendir sync. %s", result.Stderr)
	}
	configLocation := filepath.Join(upstreamStep.PkgLocation, "bundle", "config")
	upstreamStep.pkgAuthoringUI.PrintInformationalText("\nTo validate that data has been fetched, lets list down the files")
	upstreamStep.pkgAuthoringUI.PrintActionableText(fmt.Sprintf("Validating by listing files"))
	upstreamStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("ls -lR %s", configLocation))
	err := upstreamStep.listFiles(configLocation)
	if err != nil {
		return err
	}
	return nil
}

func (upstreamStep UpstreamStep) listFiles(dir string) error {
	result := util.Execute("ls", []string{"-lR", dir})
	if result.Error != nil {
		return fmt.Errorf("Listing files.\n %s", result.Stderr)
	}
	upstreamStep.pkgAuthoringUI.PrintCmdExecutionOutput(result.Stdout)
	return nil
}

func (upstreamStep UpstreamStep) getIncludedPaths() ([]string, error) {
	upstreamStep.pkgAuthoringUI.PrintInformationalText("We need exact manifest files from the above provided repository which should be included as package content. Multiple files can be included using a comma separator. If you want to include all the files, enter *.")
	includedPaths := upstreamStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].IncludePaths
	defaultIncludedPath := strings.Join(includedPaths, ",")
	if len(includedPaths) == 0 {
		defaultIncludedPath = IncludeAllFiles
	}
	textOpts := ui.TextOpts{
		Label:        "Enter the paths which need to be included as part of this package",
		Default:      defaultIncludedPath,
		ValidateFunc: nil,
	}
	path, err := upstreamStep.pkgAuthoringUI.AskForText(textOpts)
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
