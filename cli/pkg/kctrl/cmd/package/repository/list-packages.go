package repository

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	goexec "os/exec"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	lclcfg "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
)

type ListPackagesOptions struct {
	ui             ui.UI
	logger         logger.Logger
	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
	URL            string
}

func NewListPackagesOptions(ui ui.UI, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *ListPackagesOptions {
	return &ListPackagesOptions{ui: ui, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewListPackagesCmd(o *ListPackagesOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-packages",
		Aliases: []string{"lp", "ls-p"},
		Short:   "List packages from a package repository from the given URL",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
		Example: cmdcore.Examples{
			cmdcore.Example{"List package repositories",
				[]string{"package", "repository", "list-packages", "--url", "<package-repository-url>"},
			},
		}.Description("", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations: map[string]string{"table": "",
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value},
	}

	cmd.Flags().StringVarP(&o.URL, "url", "u", "", "Package repository URL")
	cmd.MarkFlagRequired(o.URL)

	return cmd
}

func (o *ListPackagesOptions) Run() error {
	if len(o.URL) == 0 {
		return fmt.Errorf("Expected package repository url to be non-empty")
	}

	configs, err := LoadImgpkgBundleToConfigs(o.ui, o.URL)
	if err != nil {
		return err
	}

	tableTitle := fmt.Sprintf("Listing Packages from Package Repository '%s'", o.URL)
	table := uitable.Table{
		Title: tableTitle,

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Version"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	for _, pkg := range configs.Pkgs {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(pkg.Spec.RefName),
			uitable.NewValueString(pkg.Spec.Version),
		})
	}

	o.ui.PrintTable(table)
	return nil
}

func LoadImgpkgBundleToConfigs(ui ui.UI, bundleURL string) (lclcfg.Configs, error) {
	var configs lclcfg.Configs
	tmpDir, err := os.MkdirTemp(".", fmt.Sprintf("bundle-%s-*", strings.Replace(bundleURL, "/", "-", -1)))
	if err != nil {
		return configs, err
	}
	defer os.RemoveAll(tmpDir)

	cmd := goexec.Command("imgpkg", "pull", "-b", bundleURL, "-o", tmpDir, "--tty=true")
	ui.PrintLinef(fmt.Sprintf("$ %s", strings.Join(cmd.Args, " ")))

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		return configs, fmt.Errorf("%s", stderrBuf.String())
	}
	ui.PrintLinef(stdoutBuf.String())

	filePaths := []string{}
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		filePaths = append(filePaths, path)
		return nil
	})

	if err != nil {
		return configs, err
	}

	configs, err = lclcfg.NewConfigFromFiles(filePaths)
	return configs, nil
}
