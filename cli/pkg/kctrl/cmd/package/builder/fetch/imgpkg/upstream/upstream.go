package upstream

import (
	"fmt"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	"os"
	"path/filepath"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"sigs.k8s.io/yaml"
)

const (
	GithubRelease int = iota
	HelmChart
	Image
)

type Content struct {
	Path          string     `json:"path"`
	GithubRelease GithubStep `json:"githubRelease,omitempty"`
	IncludePaths  []string   `json:"includePaths"`
}

type Directory struct {
	Path     string    `json:"path"`
	Contents []Content `json:"contents"`
}
type UpstreamStep struct {
	ApiVersion             string      `json:"apiVersion"`
	Kind                   string      `json:"kind"`
	MinimumRequiredVersion string      `json:"minimumRequiredVersion"`
	Directories            []Directory `json:"directories"`
	ui                     ui.UI       `json:"-"`
	PkgLocation            string      `json:"-"`
}

func NewUpstreamStep(ui ui.UI, pkgLocation string) *UpstreamStep {
	return &UpstreamStep{
		ui:          ui,
		PkgLocation: pkgLocation,
	}
}

func (upstreamStep *UpstreamStep) PreInteract() error {
	str := `
In Carvel, An upstream source is the location from where we want to sync the software configuration.
Different types of upstream available are`
	upstreamStep.ui.BeginLinef(str)
	return nil
}

func (upstreamStep *UpstreamStep) PostInteract() error {
	upstreamStep.populateUpstreamMetadata()
	err := upstreamStep.createVendirFile()
	if err != nil {
		return err
	}
	err = upstreamStep.printVendirFile()
	if err != nil {
		return err
	}
	err = upstreamStep.syncDataFromUpstream()
	if err != nil {
		return err
	}
	err = upstreamStep.printVendirLockFile()
	if err != nil {
		return err
	}
	return nil
}

func (upstreamStep *UpstreamStep) createVendirFile() error {
	str := `We have all the information needed to sync the upstream.
Lets build vendir.yml file with above inputs.
`
	upstreamStep.ui.BeginLinef(str)
	data, err := yaml.Marshal(&upstreamStep)
	vendirFileLocation := filepath.Join(upstreamStep.PkgLocation, "bundle", "vendir.yml")
	if err != nil {
		fmt.Errorf("Unable to build vendir.yml")
		return err
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
	vendirFileLocation := filepath.Join(upstreamStep.PkgLocation, "bundle", "vendir.yml")
	str := `	$ cat vendir.yml`
	upstreamStep.ui.BeginLinef(str)
	fmt.Println()
	resp, err := util.Execute("cat", []string{vendirFileLocation})
	if err != nil {
		fmt.Println("Unable to read vendir.yaml file")
		return err
	}
	upstreamStep.ui.PrintBlock([]byte(resp))
	return nil
}

func (upstreamStep *UpstreamStep) populateUpstreamMetadata() {
	upstreamStep.ApiVersion = "vendir.k14s.io/v1alpha1"
	upstreamStep.Kind = "Config"
	upstreamStep.MinimumRequiredVersion = "0.12.0"
}

func (upstreamStep *UpstreamStep) printVendirLockFile() error {
	vendirLockFileLocation := filepath.Join(upstreamStep.PkgLocation, "bundle", "vendir.lock.yml")
	str := fmt.Sprintf(`
After running vendir sync, there is one more file created i.e. bundle/vendir.lock.yml
This lock file resolves the release tag to the specific GitHub release and declares that the config is the synchronization target path.
Lets see its content
	$ cat %s
---
`, vendirLockFileLocation)
	upstreamStep.ui.BeginLinef(str)
	output, err := util.Execute("cat", []string{vendirLockFileLocation})
	if err != nil {
		return err
	}
	upstreamStep.ui.PrintBlock([]byte(output))
	return nil
}

func (upstreamStep *UpstreamStep) syncDataFromUpstream() error {
	bundleLocation := filepath.Join(upstreamStep.PkgLocation, "bundle")
	str := fmt.Sprintf(`
Next step is to run vendir to sync the data from upstream.
	$ vendir sync --chdir %s
`, bundleLocation)
	upstreamStep.ui.BeginLinef(str)
	_, err := util.Execute("vendir", []string{"sync", "--chdir", bundleLocation})
	if err != nil {
		fmt.Printf("Error while running vendir sync. Error is: %s", err.Error())
		return err
	}
	configLocation := filepath.Join(upstreamStep.PkgLocation, "bundle", "config")
	str = fmt.Sprintf(`To ensure that data has been synced, lets do
	$ ls -l %s
`, configLocation)
	upstreamStep.ui.BeginLinef(str)
	output, err := util.Execute("ls", []string{"-l", configLocation})
	if err != nil {
		return err
	}
	upstreamStep.ui.BeginLinef(output)
	return nil
}

func (upstreamStep *UpstreamStep) Interact() error {
	upstreamTypeSelected, err := upstreamStep.ui.AskForChoice("Enter the upstream type", []string{"Github Release", "HelmChart", "Image"})
	if err != nil {
		//TODO Rohit error handling
	}
	var content Content
	switch upstreamTypeSelected {
	case GithubRelease:
		githubStep := NewGithubStep(upstreamStep.ui)
		err := githubStep.Run()
		if err != nil {
			return err
		}
		content.GithubRelease = *githubStep
	}

	includedPaths, err := upstreamStep.getIncludedPaths()
	if err != nil {
		return err
	}
	content.IncludePaths = includedPaths
	content.Path = "."

	upstreamStep.Directories = []Directory{
		Directory{
			Path: "config",
			Contents: []Content{
				content,
			},
		},
	}

	return nil
}

func (upstreamStep UpstreamStep) getIncludedPaths() ([]string, error) {
	var includeEverything bool
	input, _ := upstreamStep.ui.AskForText("Does your package needs to include everything from the upstream(y/n)")
	for {
		var isValidInput bool
		includeEverything, isValidInput = common.ValidateInputYesOrNo(input)
		if isValidInput {
			break
		} else {
			input, _ = upstreamStep.ui.AskForText("Invalid input. (must be 'y','n','Y','N')")
		}
	}
	var paths []string
	var err error
	if includeEverything {

	} else {
		paths, err = upstreamStep.getPaths()
		if err != nil {
			return nil, err
		}
	}
	return paths, nil
}

func (upstreamStep UpstreamStep) getPaths() ([]string, error) {
	str := `Now, we need to enter the specific paths which we want to include as package content. More than one paths can be added with comma separator.`
	upstreamStep.ui.BeginLinef(str)

	path, err := upstreamStep.ui.AskForText("Enter the paths which needs to be included as part of this package")
	if err != nil {
		return nil, err
	}
	paths := strings.Split(path, ",")
	return paths, nil
}

func (upstreamStep *UpstreamStep) Run() error {
	err := upstreamStep.PreInteract()
	if err != nil {
		return err
	}
	err = upstreamStep.Interact()
	if err != nil {
		return err
	}
	err = upstreamStep.PostInteract()
	if err != nil {
		return err
	}
	return nil
}
