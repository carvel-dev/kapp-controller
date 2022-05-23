package builder

import (
	"os"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ValuesSchema struct {
	ui          ui.UI
	pkgLocation string
	pkgBuild    *build.PackageBuild
}

func NewValuesSchemaStep(ui ui.UI, pkgLocation string, pkgBuild *build.PackageBuild) *ValuesSchema {
	return &ValuesSchema{
		ui:          ui,
		pkgLocation: pkgLocation,
		pkgBuild:    pkgBuild,
	}
}

func (createStep CreateStep) getValueSchema() (v1alpha1.ValuesSchema, error) {
	valuesSchema := v1alpha1.ValuesSchema{}
	var isValueSchemaSpecified bool
	var isValidInput bool
	input, err := createStep.pkgAuthoringUI.AskForText(ui.TextOpts{
		Label: "Do you want to specify the values Schema(y/n)",
	})
	if err != nil {
		return valuesSchema, err
	}
	for {
		isValueSchemaSpecified, isValidInput = common.ValidateInputYesOrNo(input)
		if !isValidInput {
			input, err = createStep.pkgAuthoringUI.AskForText(ui.TextOpts{
				Label: "Invalid input. (must be 'y','n','Y','N')",
			})
			if err != nil {
				return valuesSchema, err
			}
			continue
		}
		if isValueSchemaSpecified {
			valuesSchemaFileLocation, err := createStep.pkgAuthoringUI.AskForText(ui.TextOpts{
				Label: "Enter the values schema file location",
			})
			if err != nil {
				return valuesSchema, err
			}
			valuesSchemaData, err := readDataFromFile(valuesSchemaFileLocation)
			if err != nil {
				return valuesSchema, err
			}
			valuesSchema = v1alpha1.ValuesSchema{
				OpenAPIv3: runtime.RawExtension{
					Raw: valuesSchemaData,
				},
			}
		} else {
			break
		}
	}
	return valuesSchema, nil

}
func readDataFromFile(fileLocation string) ([]byte, error) {
	//TODO should we read it in a buffer
	data, err := os.ReadFile(fileLocation)
	if err != nil {
		return nil, err
	}
	return data, nil
}
