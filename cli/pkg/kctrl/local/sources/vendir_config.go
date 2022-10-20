package sources

import (
	"os"

	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"sigs.k8s.io/yaml"
)

type VendirConfig struct {
	path string

	Config *vendirconf.Config
}

func NewVendirConfig(path string) *VendirConfig {
	return &VendirConfig{path: path}
}

func (c *VendirConfig) Load() error {
	_, err := os.Stat(vendirFileName)
	if err != nil {
		if os.IsNotExist(err) {
			c.Config = &vendirconf.Config{
				APIVersion: "vendir.k14s.io/v1alpha1",
				Kind:       "Config",
			}
			return nil
		}
		return err
	}

	content, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	vendirConfig := vendirconf.Config{}
	err = yaml.Unmarshal(content, &vendirConfig)
	if err != nil {
		return err
	}
	c.Config = &vendirConfig
	return nil
}

func (c *VendirConfig) FetchMode(ishelmTemplateExist bool) string {
	if c.Config.Directories == nil || c.Config.Directories[0].Contents == nil {
		return ""
	}
	if len(c.Config.Directories) > 1 || len(c.Config.Directories[0].Contents) > 1 {
		return MultipleFetchOptionsSelected
	}
	content := c.Config.Directories[0].Contents[0]
	var selectedVendirOption string
	switch {
	case content.GithubRelease != nil:
		selectedVendirOption = GithubRelease
	case content.HelmChart != nil:
		selectedVendirOption = HelmRepo
	case content.Directory != nil:
		selectedVendirOption = LocalDirectory
	case content.Git != nil:
		if ishelmTemplateExist {
			selectedVendirOption = ChartFromGit
		} else {
			selectedVendirOption = Git
		}
	}
	return selectedVendirOption
}

func (c *VendirConfig) save() error {
	content, err := yaml.Marshal(c.Config)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.path, content, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (c *VendirConfig) Save() error {
	return c.save()
}

func (c *VendirConfig) InitiatliseDirectories() error {
	directories := []vendirconf.Directory{vendirconf.Directory{Path: VendirSyncDirectory, Contents: []vendirconf.DirectoryContents{{Path: "."}}}}
	c.Config.Directories = directories
	err := c.save()
	if err != nil {
		return err
	}
	return nil
}

func (c *VendirConfig) Contents() []vendirconf.DirectoryContents {
	return c.Config.Directories[0].Contents
}

func (c *VendirConfig) SetContents(contents []vendirconf.DirectoryContents) {
	c.Config.Directories[0].Contents = contents
}

func (c *VendirConfig) Directories() []vendirconf.Directory {
	return c.Config.Directories
}
