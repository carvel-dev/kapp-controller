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
	VendirGitConf           string = "Git"
	VendirHgConf            string = "Hg"
	VendirHTTPConf          string = "HTTP"
	VendirImageConf         string = "Image"
	VendirImgpkgBundleConf  string = "Imgpkg"
	VendirGithubReleaseConf string = "Github Release(recommended)"
	VendirHelmChartConf     string = "HelmChart"
	VendirDirectoryConf     string = "Directory"
	VendirManualConf        string = "Manual"
	VendirInlineConf        string = "Inline"
)

type UpstreamStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	PkgLocation    string
	pkgBuild       *pkgbuilder.PackageBuild
}

func NewUpstreamStep(ui pkgui.IPkgAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *UpstreamStep {
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
	var earlierUpstreamOptionSelected string
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
		earlierUpstreamOptionSelected = getUpstreamOptionFromPkgBuild(upstreamStep.pkgBuild)
	}

	manifestOptionSelected := getManifestOptionFromPkgBuild(upstreamStep.pkgBuild)
	if manifestOptionSelected != earlierUpstreamOptionSelected {
		setEarlierUpstreamOptionAsNil(vendirDirectories, earlierUpstreamOptionSelected)
	}
	switch manifestOptionSelected {
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
		upstreamStep.pkgBuild.WriteToFile(upstreamStep.PkgLocation)
	case common.FetchChartFromHelmRepo:
		helmStep := NewHelmStep(upstreamStep.pkgAuthoringUI, upstreamStep.PkgLocation, upstreamStep.pkgBuild)
		err := common.Run(helmStep)
		if err != nil {
			return err
		}
	}

	return nil
}

func getManifestOptionFromPkgBuild(pkgBuild *pkgbuilder.PackageBuild) string {
	return pkgBuild.Annotations[common.PkgFetchContentAnnotationKey]
}

func (upstreamStep *UpstreamStep) OldInteract() error {
	var defaultUpstreamOptionSelected string

	vendirDirectories := upstreamStep.pkgBuild.Spec.Vendir.Directories
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
		defaultUpstreamOptionSelected = getUpstreamOptionFromPkgBuild(upstreamStep.pkgBuild)
	}

	var upstreamTypeNames = []string{VendirGithubReleaseConf, VendirHelmChartConf}

	defaultUpstreamOptionIndex := getDefaultUpstreamOptionIndex(upstreamTypeNames, defaultUpstreamOptionSelected)
	choiceOpts := ui.ChoiceOpts{
		Label:   "Enter the upstream type",
		Default: defaultUpstreamOptionIndex,
		Choices: upstreamTypeNames,
	}
	upstreamTypeSelectedIndex, err := upstreamStep.pkgAuthoringUI.AskForChoice(choiceOpts)
	if err != nil {
		//TODO Rohit error handling
	}
	if defaultUpstreamOptionIndex != upstreamTypeSelectedIndex {
		setEarlierUpstreamOptionAsNil(vendirDirectories, defaultUpstreamOptionSelected)
	}
	switch upstreamTypeNames[upstreamTypeSelectedIndex] {
	case VendirGithubReleaseConf:
		githubStep := NewGithubStep(upstreamStep.pkgAuthoringUI, upstreamStep.PkgLocation, upstreamStep.pkgBuild)
		err := common.Run(githubStep)
		if err != nil {
			return err
		}
	case VendirHelmChartConf:
		helmStep := NewHelmStep(upstreamStep.pkgAuthoringUI, upstreamStep.PkgLocation, upstreamStep.pkgBuild)
		err := common.Run(helmStep)
		if err != nil {
			return err
		}
	}
	if VendirHelmChartConf != upstreamTypeNames[upstreamTypeSelectedIndex] {
		includedPaths, err := upstreamStep.getIncludedPaths()
		if err != nil {
			return err
		}
		upstreamStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].IncludePaths = includedPaths
		upstreamStep.pkgBuild.WriteToFile(upstreamStep.PkgLocation)
	}
	return nil
}

func setEarlierUpstreamOptionAsNil(vendirDirectories []vendirconf.Directory, earlierUpstreamOption string) {
	if vendirDirectories[0].Contents == nil {
		return
	}
	switch earlierUpstreamOption {
	case VendirGitConf:
		vendirDirectories[0].Contents[0].Git = nil
	case VendirHgConf:
		vendirDirectories[0].Contents[0].Hg = nil
	case VendirHTTPConf:
		vendirDirectories[0].Contents[0].HTTP = nil
	case VendirImageConf:
		vendirDirectories[0].Contents[0].Image = nil
	case VendirImgpkgBundleConf:
		vendirDirectories[0].Contents[0].ImgpkgBundle = nil
	case VendirGithubReleaseConf:
		vendirDirectories[0].Contents[0].GithubRelease = nil
	case VendirHelmChartConf:
		vendirDirectories[0].Contents[0].HelmChart = nil
	case VendirDirectoryConf:
		vendirDirectories[0].Contents[0].Directory = nil
	case VendirManualConf:
		vendirDirectories[0].Contents[0].Manual = nil
	case VendirInlineConf:
		vendirDirectories[0].Contents[0].Inline = nil
	}
	return
}

func getDefaultUpstreamOptionIndex(upstreamTypeNames []string, defaultUpstreamOptionSelected string) int {
	var defaultUpstreamOptionIndex int
	if defaultUpstreamOptionSelected == "" {
		defaultUpstreamOptionIndex = 0
	} else {
		for i, upstreamTypeName := range upstreamTypeNames {
			if upstreamTypeName == defaultUpstreamOptionSelected {
				defaultUpstreamOptionIndex = i
				break
			}
		}
	}
	return defaultUpstreamOptionIndex
}

func getUpstreamOptionFromPkgBuild(pkgBuild *pkgbuilder.PackageBuild) string {
	dirContents := pkgBuild.Spec.Vendir.Directories[0].Contents
	if dirContents == nil {
		return ""
	}
	content := pkgBuild.Spec.Vendir.Directories[0].Contents[0]
	var selectedUpstreamOption string
	switch {
	case content.Git != nil:
		selectedUpstreamOption = VendirGitConf
	case content.Hg != nil:
		selectedUpstreamOption = VendirHgConf
	case content.HTTP != nil:
		selectedUpstreamOption = VendirHTTPConf
	case content.Image != nil:
		selectedUpstreamOption = VendirImageConf
	case content.ImgpkgBundle != nil:
		selectedUpstreamOption = VendirImgpkgBundleConf
	case content.GithubRelease != nil:
		selectedUpstreamOption = common.FetchReleaseArtifactFromGithub
	case content.HelmChart != nil:
		selectedUpstreamOption = common.FetchChartFromHelmRepo
	case content.Directory != nil:
		selectedUpstreamOption = VendirDirectoryConf
	case content.Manual != nil:
		selectedUpstreamOption = VendirManualConf
	case content.Inline != nil:
		selectedUpstreamOption = VendirInlineConf
	default:
		selectedUpstreamOption = ""
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
	upstreamStep.pkgBuild.WriteToFile(upstreamStep.PkgLocation)
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
	/*err = upstreamStep.printVendirLockFile()
	if err != nil {
		return err
	}*/
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
		//TODO Rohit how are you sure that this is the error.
		fmt.Println("File already exist")
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

func (upstreamStep *UpstreamStep) printVendirLockFile() error {
	vendirLockFileLocation := filepath.Join(upstreamStep.PkgLocation, "bundle", "vendir.lock.yml")
	upstreamStep.pkgAuthoringUI.PrintInformationalText("Vendir sync creates one more file i.e. bundle/vendir.lock.yml. This lock file contains mapping of the resolved release tag to the specific GitHub release and the target path where data is synchronized.")
	upstreamStep.pkgAuthoringUI.PrintActionableText("Printing Vendir.lock.yml")
	upstreamStep.pkgAuthoringUI.PrintCmdExecutionText(fmt.Sprintf("cat %s", vendirLockFileLocation))
	err := upstreamStep.printFile(vendirLockFileLocation)
	if err != nil {
		return err
	}
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
		defaultIncludedPath = "*"
	}
	textOpts := ui.TextOpts{
		Label:        "Enter the paths which need to be included as part of this package",
		Default:      defaultIncludedPath,
		ValidateFunc: nil,
	}
	path, err := upstreamStep.pkgAuthoringUI.AskForText(textOpts)
	if err != nil {
		return nil, err
	}
	if path == "*" {
		return []string{}, nil
	}
	paths := strings.Split(path, ",")
	return paths, nil
}
