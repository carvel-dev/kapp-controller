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

	ValuesSchema bool
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

	cmd.Flags().BoolVar(&o.ValuesSchema, "values-schema", false, "Values schema of the package, optional")
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

	if o.ValuesSchema {
		if pkgVersion == "" {
			return fmt.Errorf("version is required when --values-schema flag is declared. Please specify <PACKAGE-NAME>/<VERSION>")
		}
		return o.GetValuesSchema(pkgName, pkgVersion)
	}

	headers := []uitable.Header{
		uitable.NewHeader("NAME"),
		uitable.NewHeader("DISPLAY-NAME"),
		uitable.NewHeader("SHORT-DESCRIPTION"),
		uitable.NewHeader("PACKAGE-PROVIDER"),
		uitable.NewHeader("LONG-DESCRIPTION"),
		uitable.NewHeader("MAINTAINERS"),
		uitable.NewHeader("SUPPORT"),
		uitable.NewHeader("CATEGORY"),
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
			uitable.NewHeader("VERSION"),
			uitable.NewHeader("RELEASED-AT"),
			uitable.NewHeader("MINIMUM-CAPACITY-REQUIREMENTS"),
			uitable.NewHeader("RELEASE-NOTES"),
			uitable.NewHeader("LICENSE"),
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

func (o *GetOptions) GetValuesSchema(pkgName, pkgVersion string) error {
	client, err := o.depsFactory.PackageClient()
	if err != nil {
		return err
	}

	pkg, err := client.DataV1alpha1().Packages(o.NamespaceFlags.Name).Get(
		context.Background(), fmt.Sprintf("%s.%s", pkgName, pkgVersion), metav1.GetOptions{})
	if err != nil {
		return err
	}

	if len(pkg.Spec.ValuesSchema.OpenAPIv3.Raw) == 0 {
		o.ui.PrintLinef("package '%s/%s' does not have any user configurable values in the '%s' namespace", pkgName, pkgVersion, o.NamespaceFlags.Name)
		return nil
	}

	dataValuesSchemaParser, err := NewValuesSchemaParser(pkg.Spec.ValuesSchema)
	if err != nil {
		return err
	}
	parsedProperties, err := dataValuesSchemaParser.ParseProperties()
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title: fmt.Sprintf("Values schema for '%s'", pkgName),

		Header: []uitable.Header{
			uitable.NewHeader("KEY"),
			uitable.NewHeader("DEFAULT"),
			uitable.NewHeader("TYPE"),
			uitable.NewHeader("DESCRIPTION"),
		},
	}

	for _, v := range parsedProperties {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(v.Key),
			uitable.NewValueInterface(v.Default),
			uitable.NewValueInterface(v.Type),
			uitable.NewValueInterface(v.Description),
		})
	}

	o.ui.PrintTable(table)

	return err
}
