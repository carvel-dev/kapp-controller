package upstream

import (
	"fmt"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/util"
	"sigs.k8s.io/yaml"

	"os"
	"strings"
)

const (
	GithubRelease int = iota
	HelmChart
	Image
)

type Content struct {
	Path          string
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
	Ui                     ui.UI       `json:"-"`
	PkgName                string      `json:"-"`
	PkgLocation            string      `json:"-"`
	PkgVersionLocation     string      `json:"-"`
}

func NewUpstreamStep(ui ui.UI, pkgName string, pkgLocation string, pkgVersionLocation string) *UpstreamStep {
	return &UpstreamStep{
		Ui:                 ui,
		PkgName:            pkgName,
		PkgLocation:        pkgLocation,
		PkgVersionLocation: pkgVersionLocation,
	}
}

func (u *UpstreamStep) PreInteract() error {
	str := `
# In Carvel, An upstream source is the location from where we want to sync the software configuration.
# Different types of upstream available are`
	u.Ui.PrintBlock([]byte(str))
	return nil
}

func (u *UpstreamStep) PostInteract() error {
	u.populateUpstreamMetadata()

	u.createVendirFile()

	u.printVendirFile()

	u.syncDataFromUpstream()

	u.printVendirLockFile()

	return nil
}

func (u *UpstreamStep) createVendirFile() error {
	str := `# We have all the information needed to sync the upstream.
# Lets build vendir.yml file with above inputs.
`
	u.Ui.PrintBlock([]byte(str))
	data, err := yaml.Marshal(&u)
	vendirFileLocation := u.PkgVersionLocation + "/bundle/vendir.yml"
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

func (u *UpstreamStep) printVendirFile() error {
	vendirFileLocation := u.PkgVersionLocation + "/bundle/vendir.yml"
	str := `cat vendir.yml`
	u.Ui.PrintBlock([]byte(str))
	fmt.Println()
	resp, err := util.Execute("cat", []string{vendirFileLocation})
	if err != nil {
		fmt.Println("Unable to read vendir.yaml file")
		return err
	}
	u.Ui.PrintBlock([]byte(resp))
	return nil
}

func (u *UpstreamStep) populateUpstreamMetadata() {
	u.ApiVersion = "vendir.k14s.io/v1alpha1"
	u.Kind = "Config"
	u.MinimumRequiredVersion = "0.12.0"
}

func (u *UpstreamStep) printVendirLockFile() error {
	vendirLockFileLocation := u.PkgVersionLocation + "/bundle/vendir.lock.yml"
	str := fmt.Sprintf(`
# After running vendir sync, there is one more file created i.e. bundle/vendir.lock.yml
# This lock file resolves the release tag to the specific GitHub release and declares that the config is the synchronization target path.
# Lets see its content
# cat %s`, vendirLockFileLocation)
	u.Ui.PrintBlock([]byte(str))
	output, err := util.Execute("cat", []string{vendirLockFileLocation})
	if err != nil {
		return err
	}
	u.Ui.PrintBlock([]byte(output))
	return nil
}

func (u *UpstreamStep) syncDataFromUpstream() error {
	bundleLocation := u.PkgVersionLocation + "/bundle"
	str := fmt.Sprintf(`
# Next step is to Run vendir to sync the data from upstream.
# Running vendir sync --chdir %s
`, bundleLocation)
	u.Ui.PrintBlock([]byte(str))
	resp, err := util.Execute("vendir", []string{"sync", "--chdir", bundleLocation})
	if err != nil {
		fmt.Printf("Error while running vendir sync. Error is: %s", err.Error())
		return err
	}
	u.Ui.PrintBlock([]byte(resp))
	configLocation := u.PkgVersionLocation + "/bundle/config"
	str = fmt.Sprintf(`# To ensure that data has been synced, lets do
# ls -l %s`, configLocation)
	u.Ui.PrintBlock([]byte(str))
	output, err := util.Execute("ls", []string{"-l", configLocation})
	if err != nil {
		return err
	}
	u.Ui.PrintBlock([]byte(output))
	return nil
}

func (u *UpstreamStep) Interact() error {
	upstreamTypeSelected, err := u.Ui.AskForChoice("Enter the upstream type", []string{"Github Release", "HelmChart", "Image"})
	if err != nil {

	}
	var content Content
	switch upstreamTypeSelected {
	case GithubRelease:
		githubStep := NewGithubStep(u.Ui)
		githubStep.Run()
		content.GithubRelease = *githubStep

	}

	includedPaths, err := u.getIncludedPaths()
	if err != nil {
		return err
	}
	content.IncludePaths = includedPaths
	content.Path = "."

	u.Directories = []Directory{
		Directory{
			Path: "config",
			Contents: []Content{
				content,
			},
		},
	}

	return nil
}

func (u UpstreamStep) getIncludedPaths() ([]string, error) {
	var includeEverything bool
	input, _ := u.Ui.AskForText("Does your package needs to include everything from the upstream(y/n)")
	for {
		var isValidInput bool
		includeEverything, isValidInput = common.ValidateInputYesOrNo(input)
		if isValidInput {
			break
		} else {
			input, _ = u.Ui.AskForText("Invalid input. (must be 'y','n','Y','N')")
		}
	}
	var paths []string
	var err error
	if includeEverything {

	} else {
		paths, err = u.getPaths()
		if err != nil {
			return nil, err
		}
	}
	return paths, nil
}

func (u UpstreamStep) getPaths() ([]string, error) {
	str := `# Now, we need to enter the specific paths which we want to include as package content. More than one paths can be added with comma separator.`
	u.Ui.PrintBlock([]byte(str))

	path, err := u.Ui.AskForText("Enter the paths which needs to be included as part of this package")
	if err != nil {
		return nil, err
	}
	paths := strings.Split(path, ",")
	return paths, nil
}

func (u *UpstreamStep) Run() error {
	u.PreInteract()
	u.Interact()
	u.PostInteract()
	return nil
}
