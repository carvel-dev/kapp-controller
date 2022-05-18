package upstream

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/util"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"sigs.k8s.io/yaml"
)

type UpstreamStep struct {
	ui          ui.UI
	PkgLocation string
	pkgBuild    *pkgbuilder.PackageBuild
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
	vendirDirectories := upstreamStep.pkgBuild.Spec.Vendir.Directories
	if len(vendirDirectories) == 0 {
		upstreamStep.initializeVendirDirectoryConf()
	}

	upstreamOptions := []string{"Github Release", "HelmChart", "Image"}
	upstreamTypeSelected, err := upstreamStep.ui.AskForChoice("Enter the upstream type", upstreamOptions)
	if err != nil {
		//TODO Rohit error handling
	}

	switch upstreamOptions[upstreamTypeSelected] {
	case "Github Release":
		githubStep := NewGithubStep(upstreamStep.ui, upstreamStep.PkgLocation, upstreamStep.pkgBuild)

		err := common.Run(githubStep)
		if err != nil {
			return err
		}

	}
	includedPaths, err := upstreamStep.getIncludedPaths()
	if err != nil {
		return err
	}
	upstreamStep.pkgBuild.Spec.Vendir.Directories[0].Contents[0].IncludePaths = includedPaths
	upstreamStep.pkgBuild.WriteToFile(upstreamStep.PkgLocation)
	return nil
}

func (upstreamStep *UpstreamStep) initializeVendirDirectoryConf() {
	var directory vendirconf.Directory
	directory = vendirconf.Directory{
		Path: "config",
		Contents: []vendirconf.DirectoryContents{
			{
				Path: ".",
			},
		},
	}
	directories := []vendirconf.Directory{}
	upstreamStep.pkgBuild.Spec.Vendir.Directories = append(directories, directory)
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
	data, err := yaml.Marshal(&upstreamStep.pkgBuild.Spec.Vendir)
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
	err := upstreamStep.printFile(vendirFileLocation)
	if err != nil {
		fmt.Println("Unable to read vendir.yaml file")
		return err
	}
	return nil
}

func (upstreamStep *UpstreamStep) printFile(filePath string) error {
	result := util.Execute("cat", []string{filePath})
	if result.Error != nil {
		upstreamStep.ui.ErrorLinef("Error printing file %s.Error is: %s", filePath, result.ErrorStr())
		return result.Error
	}
	upstreamStep.ui.PrintBlock([]byte(result.Stdout))
	return nil
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
	err := upstreamStep.printFile(vendirLockFileLocation)
	if err != nil {
		return err
	}
	return nil
}

func (upstreamStep *UpstreamStep) runVendirSync() error {
	bundleLocation := filepath.Join(upstreamStep.PkgLocation, "bundle")
	str := fmt.Sprintf(`
Next step is to run vendir to sync the data from upstream. Running 'vendir sync'
	$ vendir sync --chdir %s
`, bundleLocation)
	upstreamStep.ui.BeginLinef(str)
	result := util.Execute("vendir", []string{"sync", "--chdir", bundleLocation})
	if result.Error != nil {
		upstreamStep.ui.ErrorLinef("Error while running vendir sync. Error is: %s", result.Stderr)
		return result.Error
	}
	configLocation := filepath.Join(upstreamStep.PkgLocation, "bundle", "config")
	str = fmt.Sprintf(`To ensure that data has been synced, lets do
	$ ls -l %s
`, configLocation)
	upstreamStep.ui.BeginLinef(str)
	err := upstreamStep.listFiles(configLocation)
	if err != nil {
		return err
	}
	return nil
}

func (upstreamStep UpstreamStep) listFiles(dir string) error {
	result := util.Execute("ls", []string{"-l", dir})
	if result.Error != nil {
		upstreamStep.ui.ErrorLinef("Error while listing files. Error is: %s", result.ErrorStr())
		return result.Error
	}
	upstreamStep.ui.PrintBlock([]byte(result.Stdout))
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
