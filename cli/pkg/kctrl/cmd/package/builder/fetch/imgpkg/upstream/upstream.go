package upstream

import (
	"fmt"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	"sigs.k8s.io/yaml"
)

const (
	vendirAPIVersion             = "vendir.k14s.io/v1alpha1"
	vendirKind                   = "Config"
	vendirMinimumRequiredVersion = "0.12.0"
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
	VendirConfig vendirconf.Config
	ui           ui.UI
	PkgLocation  string
	pkgBuild     *pkgbuilder.PackageBuild
}

func NewUpstreamStep(ui ui.UI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *UpstreamStep {
	return &UpstreamStep{
		ui:          ui,
		PkgLocation: pkgLocation,
		pkgBuild:    pkgBuild,
	}
}

func (upstreamStep *UpstreamStep) PreInteract() error {
	str := `
In Carvel, An upstream source is the location from where we want to sync the software configuration.
Different types of upstream available are`
	upstreamStep.ui.BeginLinef(str)
	return nil
}

func (upstreamStep *UpstreamStep) Interact() error {
	upstreamOptions := []string{"Github Release", "HelmChart", "Image"}
	upstreamTypeSelected, err := upstreamStep.ui.AskForChoice("Enter the upstream type", upstreamOptions)
	if err != nil {
		//TODO Rohit error handling
	}
	contents := []vendirconf.DirectoryContents{}

	switch upstreamOptions[upstreamTypeSelected] {
	case "Github Release":
		content := vendirconf.DirectoryContents{}
		githubStep := NewGithubStep(upstreamStep.ui)
		err := githubStep.Run()
		if err != nil {
			return err
		}
		includedPaths, err := upstreamStep.getIncludedPaths()
		if err != nil {
			return err
		}
		content.IncludePaths = includedPaths
		content.Path = "."
		content.GithubRelease = githubStep.GithubRelease
		contents = append(contents, content)
	}

	directory := vendirconf.Directory{
		Path:     "config",
		Contents: contents,
	}
	directories := []vendirconf.Directory{}
	upstreamStep.VendirConfig.Directories = append(directories, directory)
	return nil
}

func (upstreamStep *UpstreamStep) PostInteract() error {
	upstreamStep.populateUpstreamMetadata()
	upstreamStep.pkgBuild.Spec.Vendir = upstreamStep.VendirConfig
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
	err = upstreamStep.printVendirLockFile()
	if err != nil {
		return err
	}
	return nil
}

func (upstreamStep *UpstreamStep) createVendirFile() error {
	vendirFileLocation := filepath.Join(upstreamStep.PkgLocation, "bundle", "vendir.yml")
	str := fmt.Sprintf(`We have all the information needed to sync the upstream.
To create an imgpkg bundle, data has to be synced from upstream to local. 
To sync the data from upstream to local, we will use vendir.
Vendir allows to declaratively state what should be in a directory and sync any number of data sources into it.
Lets use our inputs to create vendir.yml file.
Creating vendir.yml file in directory %s
`, vendirFileLocation)
	upstreamStep.ui.BeginLinef(str)
	data, err := yaml.Marshal(&upstreamStep.VendirConfig)
	if err != nil {
		upstreamStep.ui.ErrorLinef("Unable to create vendir.yml")
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
	str := `Our vendir.yml is created. This file looks like this
	$ cat vendir.yml`
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
	upstreamStep.VendirConfig.APIVersion = vendirAPIVersion
	upstreamStep.VendirConfig.Kind = vendirKind
	upstreamStep.VendirConfig.MinimumRequiredVersion = vendirMinimumRequiredVersion
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

func (upstreamStep *UpstreamStep) runVendirSync() error {
	bundleLocation := filepath.Join(upstreamStep.PkgLocation, "bundle")
	str := fmt.Sprintf(`
Next step is to run vendir to sync the data from upstream. Running 'vendir sync'
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

func (upstreamStep UpstreamStep) getIncludedPaths() ([]string, error) {
	str := `Now, we need to enter the specific paths which we want to include as package content. More than one paths can be added with comma separator. 
To include everything from the upstream, leave it empty`
	upstreamStep.ui.BeginLinef(str)
	path, err := upstreamStep.ui.AskForText("Enter the paths which need to be included as part of this package")
	if err != nil {
		return nil, err
	}
	if len(path) == 0 {
		return []string{}, nil
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
