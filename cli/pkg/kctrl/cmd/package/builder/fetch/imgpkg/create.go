package imgpkg

import (
	"fmt"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/util"
	"path/filepath"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/fetch/imgpkg/upstream"
	"k8s.io/apimachinery/pkg/util/json"
)

type CreateImgPkgStep struct {
	ui          ui.UI
	image       string
	pkgLocation string
}

func NewCreateImgPkgStep(ui ui.UI, pkgLocation string) *CreateImgPkgStep {
	return &CreateImgPkgStep{
		ui:          ui,
		pkgLocation: pkgLocation,
	}
}

func (createImgPkgStep *CreateImgPkgStep) Run() error {
	err := createImgPkgStep.PreInteract()
	if err != nil {
		return err
	}
	err = createImgPkgStep.Interact()
	if err != nil {
		return err
	}
	err = createImgPkgStep.PostInteract()
	if err != nil {
		return err
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) PreInteract() error {
	err := createImgPkgStep.createBundleConfigDir()
	if err != nil {
		return err
	}

	err = createImgPkgStep.createBundleDotImgpkgDir()
	if err != nil {
		return err
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) createBundleConfigDir() error {
	bundleConfigLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle", "config")
	str := fmt.Sprintf(`Cool, lets create the imgpkg bundle first
Creating directory %s
	$ mkdir -p %s
`, bundleConfigLocation, bundleConfigLocation)
	createImgPkgStep.ui.BeginLinef(str)
	output, err := util.Execute("mkdir", []string{"-p", bundleConfigLocation})
	if err != nil {
		return err
	}
	createImgPkgStep.ui.BeginLinef(output)
	return nil
}

func (createImgPkgStep CreateImgPkgStep) createBundleDotImgpkgDir() error {
	bundleDotImgPkgLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle", ".imgpkg")
	str := fmt.Sprintf(`
Creating directory %s
	$ mkdir -p %s
`, bundleDotImgPkgLocation, bundleDotImgPkgLocation)
	createImgPkgStep.ui.BeginLinef(str)
	output, err := util.Execute("mkdir", []string{"-p", bundleDotImgPkgLocation})
	if err != nil {
		return err
	}
	createImgPkgStep.ui.BeginLinef(output)
	return nil
}

func (createImgPkgStep CreateImgPkgStep) Interact() error {
	upstreamStep := upstream.NewUpstreamStep(createImgPkgStep.ui, createImgPkgStep.pkgLocation)
	err := upstreamStep.Run()
	if err != nil {
		return err
	}
	/*
	   	str := `
	   # If you wish to use default values, then skip next step. Otherwise, we can use the ytt(a templating and overlay tool) to provide custom values.`
	   	createImgPkgStep.Ui.PrintBlock([]byte(str))
	   	var useYttAsTemplate bool
	   	for {
	   		input, err := createImgPkgStep.Ui.AskForText("Do you want to use ytt as a templating and overlay tool(y/n)")
	   		if err != nil {
	   			return err
	   		}
	   		var isValidInput bool
	   		useYttAsTemplate, isValidInput = common.ValidateInputYesOrNo(input)
	   		if isValidInput {
	   			break
	   		} else {
	   			input, _ = createImgPkgStep.Ui.AskForText("Invalid input (must be 'y','n','Y','N')")
	   		}
	   	}
	   	if useYttAsTemplate {
	   		yttPath, err := createImgPkgStep.Ui.AskForText("Enter the path where ytt files are located:")
	   		if err != nil {
	   			return err
	   		}
	   		configDirLocation := createImgPkgStep.PkgLocation + "/bundle/config"
	   		str = fmt.Sprintf(`# Copying the ytt files inside the package.
	   # cp -r %s %s`, yttPath, configDirLocation)
	   		createImgPkgStep.Ui.PrintBlock([]byte(str))
	   		util.Execute("cp", []string{"-r", yttPath, configDirLocation})
	   	}
	*/
	return nil
}

func (createImgPkgStep *CreateImgPkgStep) PostInteract() error {
	imagesFileLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle", ".imgpkg", "images.yml")
	bundleLocation := filepath.Join(createImgPkgStep.pkgLocation, "bundle")
	str := fmt.Sprintf(`
imgpkg bundle configuration is now complete. Let's use kbld to lock it down.
kbld allows to build the imgpkg bundle with immutable image references.
kbld scans a package configuration for any references to images and creates a mapping of image tags to a URL with a sha256 digest. 
This mapping will then be placed into an images.yml lock file in your bundle/.imgpkg directory. Running kbld now.
	$ kbld --file %s --imgpkg-lock-output %s`, bundleLocation, imagesFileLocation)
	createImgPkgStep.ui.BeginLinef(str)

	output, err := util.Execute("kbld", []string{"--file", bundleLocation, "--imgpkg-lock-output", imagesFileLocation})
	if err != nil {
		createImgPkgStep.ui.BeginLinef(err.Error())
		return err
	}

	str = fmt.Sprintf(`
Lets see how the images.yml file looks like:
Running cat %s
`, imagesFileLocation)
	createImgPkgStep.ui.BeginLinef(str)
	output, err = util.Execute("cat", []string{imagesFileLocation})
	if err != nil {
		return err
	}
	createImgPkgStep.ui.BeginLinef(output)

	var pushBundle bool
	for {
		input, _ := createImgPkgStep.ui.AskForText("Do you want to push the bundle to the registry(y/n)")
		var isValidInput bool
		pushBundle, isValidInput = common.ValidateInputYesOrNo(input)
		if isValidInput {
			break
		} else {
			input, _ = createImgPkgStep.ui.AskForText("Invalid input (must be 'y','n','Y','N')")
		}
	}
	if pushBundle {
		bundleURL, err := createImgPkgStep.pushImgpkgBundleToRegistry(bundleLocation)
		if err != nil {
			return err
		}
		createImgPkgStep.image = bundleURL
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) pushImgpkgBundleToRegistry(bundleLoc string) (string, error) {
	registryAuthDetails, err := createImgPkgStep.PopulateRegistryAuthDetails()
	if err != nil {
		return "", err
	}
	//Can repoName be empty?
	repoName, err := createImgPkgStep.ui.AskForText("Provide the repository name to which this bundle belong")
	if err != nil {
		return "", err
	}
	tagName, err := createImgPkgStep.ui.AskForText("Enter the tag for bundle")
	if err != nil {
		return "", err
	}
	pushURL := registryAuthDetails.RegistryURL + "/" + repoName + ":" + tagName
	str := fmt.Sprintf(`Running imgpkg to push the bundle directory and indicate what project name and tag to give it.
 	$ imgpkg push --bundle %s --file %s --json
`, pushURL, bundleLoc)
	createImgPkgStep.ui.BeginLinef(str)

	//TODO Rohit It is not showing the actual error
	output, err := util.Execute("imgpkg", []string{"push", "--bundle", pushURL, "--file", bundleLoc, "--registry-username", registryAuthDetails.Username, "--registry-password", registryAuthDetails.Password, "--json"})
	if err != nil {
		return "", err
	}
	createImgPkgStep.ui.BeginLinef(output)
	bundleURL := getBundleURL(output)
	return bundleURL, nil
}

type ImgpkgPushOutput struct {
	Lines  []string    `json:"Lines"`
	Tables interface{} `json:"Tables"`
	Blocks interface{} `json:"Blocks"`
}

func getBundleURL(output string) string {
	var imgPkgPushOutput ImgpkgPushOutput
	json.Unmarshal([]byte(output), &imgPkgPushOutput)
	for _, val := range imgPkgPushOutput.Lines {
		if strings.HasPrefix(val, "Pushed") {
			bundleURL := strings.Split(val, " ")[1]
			return strings.Trim(bundleURL, "'")
		}
	}
	return ""

}
