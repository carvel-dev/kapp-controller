// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"os"

	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"sigs.k8s.io/yaml"
)

func ReadVendirConfig() (vendirconf.Config, error) {
	var vendirConfig vendirconf.Config
	exists, err := IsFileExists(VendirFileName)
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
		APIVersion: "vendir.k14s.io/v1alpha1",
		Kind:       "Config",
	}
	return config
}

func GetFetchOptionFromVendir(config vendirconf.Config, ishelmTemplateExist bool) string {
	if config.Directories == nil || config.Directories[0].Contents == nil {
		return ""
	}
	if len(config.Directories) > 1 || len(config.Directories[0].Contents) > 1 {
		return MultipleFetchOptionsSelected
	}
	content := config.Directories[0].Contents[0]
	var selectedVendirOption string
	switch {
	case content.GithubRelease != nil:
		selectedVendirOption = FetchFromGithubRelease
	case content.HelmChart != nil:
		selectedVendirOption = FetchFromHelmRepo
	case content.Directory != nil:
		selectedVendirOption = FetchFromLocalDirectory
	case content.Git != nil:
		if ishelmTemplateExist {
			selectedVendirOption = FetchChartFromGit
		} else {
			selectedVendirOption = FetchFromGit
		}
	}
	return selectedVendirOption
}

func SaveVendir(config vendirconf.Config) error {
	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	WriteFile(VendirFileName, content)
	return nil
}
