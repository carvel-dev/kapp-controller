// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package available

import (
	"context"
	"fmt"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string
}

func NewGetOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *GetOptions {
	return &GetOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewGetCmd(o *GetOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get details for an available package or the openAPI schema of a package with a specific version",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "package", "p", "", "Set package name")
	return cmd
}

func (o *GetOptions) Run() error {
	var pkgName, pkgVersion string
	pkgNameVersion := strings.Split(o.Name, "/")
	if len(pkgNameVersion) == 2 {
		pkgName = pkgNameVersion[0]
		pkgVersion = pkgNameVersion[1]
	} else if len(pkgNameVersion) == 1 {
		pkgName = pkgNameVersion[0]
	} else {
		return fmt.Errorf("Package name should be of the format 'name' or 'name/version'")
	}

	headers := []uitable.Header{
		uitable.NewHeader("name"),
		uitable.NewHeader("display-name"),
		uitable.NewHeader("short-description"),
		uitable.NewHeader("package-provider"),
		uitable.NewHeader("long-description"),
		uitable.NewHeader("maintainers"),
		uitable.NewHeader("support"),
		uitable.NewHeader("category"),
	}

	client, err := o.depsFactory.PackageClient()
	if err != nil {
		return err
	}

	pkgMetadata, err := client.DataV1alpha1().PackageMetadatas(
		o.NamespaceFlags.Name).Get(context.Background(), pkgName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	row := []uitable.Value{
		uitable.NewValueString(pkgMetadata.Name),
		uitable.NewValueString(pkgMetadata.Spec.DisplayName),
		uitable.NewValueString(pkgMetadata.Spec.ShortDescription),
		uitable.NewValueString(pkgMetadata.Spec.ProviderName),
		uitable.NewValueString(pkgMetadata.Spec.LongDescription),
		uitable.NewValueInterface(pkgMetadata.Spec.Maintainers),
		uitable.NewValueString(pkgMetadata.Spec.SupportDescription),
		uitable.NewValueStrings(pkgMetadata.Spec.Categories),
	}

	if pkgVersion != "" {
		// TODO should we use --field-selector?
		pkg, err := client.DataV1alpha1().Packages(o.NamespaceFlags.Name).Get(
			context.Background(), fmt.Sprintf("%s.%s", pkgName, pkgVersion), metav1.GetOptions{})
		if err != nil {
			return err
		}

		headers = append(headers, []uitable.Header{
			uitable.NewHeader("version"),
			uitable.NewHeader("released-at"),
			uitable.NewHeader("minimum-capacity-requirements"),
			uitable.NewHeader("release-notes"),
			uitable.NewHeader("license"),
		}...)

		row = append(row, []uitable.Value{
			uitable.NewValueString(pkg.Spec.Version),
			uitable.NewValueString(pkg.Spec.ReleasedAt.String()),
			uitable.NewValueString(pkg.Spec.CapactiyRequirementsDescription),
			uitable.NewValueString(pkg.Spec.ReleaseNotes),
			uitable.NewValueStrings(pkg.Spec.Licenses),
		}...)
	}

	table := uitable.Table{
		// TODO better title? should it be different for package vs packagemetadata
		Title:     fmt.Sprintf("Package details for '%s'", pkgName),
		Transpose: true,

		Header: headers,
		Rows:   [][]uitable.Value{row},
	}

	o.ui.PrintTable(table)

	return nil
}
