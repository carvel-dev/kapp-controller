package imgpkg

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/fetch/imgpkg/upstream"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/util"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type ImgpkgStep struct {
	Ui           ui.UI
	PkgLocation  string
	ImgpkgBundle v1alpha1.AppFetchImgpkgBundle
}

func NewImgPkgStep(ui ui.UI, pkgLocation string) *ImgpkgStep {
	imgpkg := ImgpkgStep{
		Ui:          ui,
		PkgLocation: pkgLocation,
	}
	return &imgpkg
}

func (imgpkg ImgpkgStep) PreInteract() error {
	return nil
}

func (imgpkg *ImgpkgStep) Interact() error {
	var isImgpkgCreated bool
	input, _ := imgpkg.Ui.AskForText("Is the imgpkg bundle already created(y/n)")

	for {
		var isValidInput bool
		isImgpkgCreated, isValidInput = common.ValidateInputYesOrNo(input)
		if isValidInput {
			break
		} else {
			input, _ = imgpkg.Ui.AskForText("Invalid input (must be 'y','n','Y','N')")
		}
	}

	if isImgpkgCreated {

	} else {
		createImgPkgStep := NewCreateImgPkgStep(imgpkg.Ui, imgpkg.PkgLocation)
		err := createImgPkgStep.Run()
		if err != nil {
			return err
		}
		imgpkg.ImgpkgBundle.Image = createImgPkgStep.Image
	}
	return nil
}

func (imgpkg ImgpkgStep) PostInteract() error {
	return nil
}

func (imgpkg *ImgpkgStep) Run() error {
	err := imgpkg.PreInteract()
	if err != nil {
		return err
	}
	err = imgpkg.Interact()
	if err != nil {
		return err
	}
	err = imgpkg.PostInteract()
	if err != nil {
		return err
	}
	return nil
}

type CreateImgPkgStep struct {
	Ui          ui.UI
	Image       string
	PkgLocation string
}

func NewCreateImgPkgStep(ui ui.UI, pkgLocation string) *CreateImgPkgStep {
	return &CreateImgPkgStep{
		Ui:          ui,
		PkgLocation: pkgLocation,
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
	bundleConfigLocation := fmt.Sprintf("%s/bundle/config", createImgPkgStep.PkgLocation)
	str := fmt.Sprintf(`# Cool, lets build the imgpkg bundle first
# Creating directory %s
# 	$ mkdir -p %s
`, bundleConfigLocation, bundleConfigLocation)
	createImgPkgStep.Ui.PrintBlock([]byte(str))
	output, err := util.Execute("mkdir", []string{"-p", bundleConfigLocation})
	if err != nil {
		return err
	}
	createImgPkgStep.Ui.PrintBlock([]byte(output))
	return nil
}

func (createImgPkgStep CreateImgPkgStep) createBundleDotImgpkgDir() error {
	bundleDotImgPkgLocation := fmt.Sprintf("%s/bundle/.imgpkg", createImgPkgStep.PkgLocation)
	str := fmt.Sprintf(`
# Creating directory %s
# 	$ mkdir -p %s
`, bundleDotImgPkgLocation, bundleDotImgPkgLocation)
	createImgPkgStep.Ui.PrintBlock([]byte(str))
	output, err := util.Execute("mkdir", []string{"-p", bundleDotImgPkgLocation})
	if err != nil {
		return err
	}
	createImgPkgStep.Ui.PrintBlock([]byte(output))
	return nil
}

func (createImgPkgStep CreateImgPkgStep) Interact() error {
	upstreamStep := upstream.NewUpstreamStep(createImgPkgStep.Ui, createImgPkgStep.PkgLocation)
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
	imagesFileLocation := fmt.Sprintf("%s/bundle/.imgpkg/images.yml", createImgPkgStep.PkgLocation)
	str := fmt.Sprintf(`
# imgpkg bundle configuration is now complete. Let's use kbld to lock it down.
# kbld allows to build the imgpkg bundle with immutable image references.
# kbld scans a package configuration for any references to images and creates a mapping of image tags to a URL with a sha256 digest. 
# This mapping will then be placed into an images.yml lock file in your bundle/.imgpkg directory. Running kbld now.
# 	$ kbld --file %s/bundle --imgpkg-lock-output %s`, createImgPkgStep.PkgLocation, imagesFileLocation)
	createImgPkgStep.Ui.PrintBlock([]byte(str))

	output, err := util.Execute("kbld", []string{"--file", createImgPkgStep.PkgLocation + "/bundle", "--imgpkg-lock-output", imagesFileLocation})
	if err != nil {
		createImgPkgStep.Ui.PrintBlock([]byte(err.Error()))
		return err
	}

	str = fmt.Sprintf(`
# Lets see how the images.yml file looks like:
# Running cat %s
`, imagesFileLocation)
	createImgPkgStep.Ui.PrintBlock([]byte(str))
	output, err = util.Execute("cat", []string{imagesFileLocation})
	if err != nil {
		return err
	}
	createImgPkgStep.Ui.PrintBlock([]byte(output))

	var pushBundle bool
	for {
		input, _ := createImgPkgStep.Ui.AskForText("Do you want to push the bundle to the registry(y/n)")
		var isValidInput bool
		pushBundle, isValidInput = common.ValidateInputYesOrNo(input)
		if isValidInput {
			break
		} else {
			input, _ = createImgPkgStep.Ui.AskForText("Invalid input (must be 'y','n','Y','N')")
		}
	}
	if pushBundle {
		bundleURL, err := createImgPkgStep.pushImgpkgBundleToRegistry(createImgPkgStep.PkgLocation + "/bundle")
		if err != nil {
			return err
		}
		createImgPkgStep.Image = bundleURL
	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) pushImgpkgBundleToRegistry(bundleLoc string) (string, error) {
	registryAuthDetails, err := createImgPkgStep.PopulateRegistryAuthDetails()
	if err != nil {
		return "", err
	}
	//Can repoName be empty?
	repoName, err := createImgPkgStep.Ui.AskForText("Provide the repository name to which this bundle belong")
	if err != nil {
		return "", err
	}
	tagName, err := createImgPkgStep.Ui.AskForText("Enter the tag for bundle")
	if err != nil {
		return "", err
	}
	pushURL := registryAuthDetails.RegistryURL + "/" + repoName + ":" + tagName
	str := fmt.Sprintf(`# Running imgpkg to push the bundle directory and indicate what project name and tag to give it.
# 	$ imgpkg push --bundle %s --file %s --json
`, pushURL, bundleLoc)
	createImgPkgStep.Ui.PrintBlock([]byte(str))

	//TODO Rohit It is not showing the actual error
	output, err := util.Execute("imgpkg", []string{"push", "--bundle", pushURL, "--file", bundleLoc, "--registry-username", registryAuthDetails.Username, "--registry-password", registryAuthDetails.Password, "--json"})
	if err != nil {
		return "", err
	}
	createImgPkgStep.Ui.PrintBlock([]byte(output))
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
