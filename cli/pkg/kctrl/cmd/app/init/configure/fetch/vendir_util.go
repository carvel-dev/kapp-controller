// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch

import (
	"os"

	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"sigs.k8s.io/yaml"
)

func ReadVendirConfig() (vendirconf.Config, error) {
	var vendirConfig vendirconf.Config
	exists, err := common.IsFileExists(VendirFileName)
	if err != nil {
		return vendirconf.Config{}, err
	}

	if exists {
		vendirConfig, err = VendirConfigFromExistingFile(VendirFileName)
		if err != nil {
			return vendirconf.Config{}, err
		}
	} else {
		vendirConfig = NewDefaultVendirConfig()
	}
	return vendirConfig, nil
}

func VendirConfigFromExistingFile(filePath string) (vendirconf.Config, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return vendirconf.Config{}, err
	}
	vendirConfig := vendirconf.Config{}
	err = yaml.Unmarshal(content, &vendirConfig)
	if err != nil {
		return vendirconf.Config{}, err
	}
	return vendirConfig, nil
}

func NewDefaultVendirConfig() vendirconf.Config {
	config := vendirconf.Config{
		APIVersion: "vendir.k14s.io/v1alpha1", // TODO: use constant from vendir package
		Kind:       "Config",                  // TODO: use constant from vendir package

	}
	SaveVendir(config)
	return config
}

func getPreviousFetchOptionFromVendir(config vendirconf.Config, ishelmTemplateExist bool) string {
	if config.Directories == nil || config.Directories[0].Contents == nil {
		return ""
	}
	content := config.Directories[0].Contents[0]
	var selectedVendirOption string
	switch {
	case content.GithubRelease != nil:
		selectedVendirOption = FetchReleaseArtifactFromGithub
	case content.HelmChart != nil:
		selectedVendirOption = FetchChartFromHelmRepo
	case content.Directory != nil:
		selectedVendirOption = FetchFromLocalDirectory
	case content.Git != nil:
		if ishelmTemplateExist {
			selectedVendirOption = FetchChartFromGithub
		}
		selectedVendirOption = FetchManifestFromGithub
	}
	return selectedVendirOption
}

func SaveVendir(config vendirconf.Config) error {
	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	common.WriteFile(VendirFileName, content)
	return nil
}
