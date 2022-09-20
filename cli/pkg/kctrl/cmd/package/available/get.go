// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package available

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	pkgclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

type GetOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	ValuesSchema      bool
	DefaultValuesFile string

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewGetOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *GetOptions {
	return &GetOptions{ui: ui, depsFactory: depsFactory, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewGetCmd(o *GetOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get details for an available package or the openAPI schema of a package with a specific version",
		Args:    cobra.ExactArgs(1),
		RunE:    func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Get details about an available package",
				[]string{"package", "available", "get", "-p", "cert-manager.community.tanzu.vmware.com"},
			},
			cmdcore.Example{"Get the values schema for a particular version of the package",
				[]string{"package", "available", "get", "-p", "cert-manager.community.tanzu.vmware.com/1.0.0", "--values-schema"}},
		}.Description("-p", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations: map[string]string{"table": "",
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value},
	}

	o.NamespaceFlags.SetWithPackageCommandTreeOpts(cmd, flagsFactory, o.pkgCmdTreeOpts)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "package", "p", "", "Set package name (required)")
	} else {
		cmd.Use = "get PACKAGE_NAME or PACKAGE_NAME/VERSION"
		cmd.Args = cobra.ExactArgs(1)
	}

	cmd.Flags().BoolVar(&o.ValuesSchema, "values-schema", false, "Values schema of the package (optional)")
	cmd.Flags().StringVar(&o.DefaultValuesFile, "default-values-file-output", "", "File path to save default values (optional)")
	return cmd
}

func (o *GetOptions) Run(args []string) error {
	var pkgName, pkgVersion string

	if o.pkgCmdTreeOpts.PositionalArgs {
		o.Name = args[0]
	}

	if len(o.Name) == 0 {
		return fmt.Errorf("Expected package name to be non-empty")
	}

	pkgNameVersion := strings.Split(o.Name, "/")
	if len(pkgNameVersion) == 2 {
		pkgName = pkgNameVersion[0]
		pkgVersion = pkgNameVersion[1]
	} else if len(pkgNameVersion) == 1 {
		pkgName = pkgNameVersion[0]
	} else {
		return fmt.Errorf("Package name should be of the format 'name' or 'name/version'")
	}

	client, err := o.depsFactory.PackageClient()
	if err != nil {
		return err
	}

	if o.ValuesSchema {
		if pkgVersion == "" {
			return fmt.Errorf("Package version is required when --values-schema flag is declared (hint: to specify a particular version use the format: '-p <package-name>/<version>')")
		}
		return o.showValuesSchema(client, pkgName, pkgVersion)
	}

	return o.show(client, pkgName, pkgVersion)
}

func (o *GetOptions) show(client pkgclient.Interface, pkgName, pkgVersion string) error {
	// Name is always present
	headers := []uitable.Header{uitable.NewHeader("Name")}
	row := []uitable.Value{uitable.NewValueString(pkgName)}

	var pkgList *v1alpha1.PackageList

	pkgMetadata, err := client.DataV1alpha1().PackageMetadatas(
		o.NamespaceFlags.Name).Get(context.Background(), pkgName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		pkgMetadata = nil
	}

	// PackageMetadata record is not required to be present
	if pkgMetadata != nil {
		headers = append(headers, []uitable.Header{
			uitable.NewHeader("Display name"),

			uitable.NewHeader("Categories"),
			uitable.NewHeader("Short description"),
			uitable.NewHeader("Long description"),

			uitable.NewHeader("Provider"),
			uitable.NewHeader("Maintainers"),
			uitable.NewHeader("Support description"),
		}...)

		row = append(row, []uitable.Value{
			uitable.NewValueString(pkgMetadata.Spec.DisplayName),

			uitable.NewValueInterface(pkgMetadata.Spec.Categories),
			uitable.NewValueString(wordwrap.WrapString(pkgMetadata.Spec.ShortDescription, 80)),
			uitable.NewValueString(wordwrap.WrapString(pkgMetadata.Spec.LongDescription, 80)),

			uitable.NewValueString(pkgMetadata.Spec.ProviderName),
			uitable.NewValueInterface(pkgMetadata.Spec.Maintainers),
			uitable.NewValueString(wordwrap.WrapString(pkgMetadata.Spec.SupportDescription, 80)),
		}...)
	}

	if pkgVersion != "" {
		// TODO should we use --field-selector?
		pkg, err := client.DataV1alpha1().Packages(o.NamespaceFlags.Name).Get(
			context.Background(), fmt.Sprintf("%s.%s", pkgName, pkgVersion), metav1.GetOptions{})
		if err != nil {
			return err
		}

		if len(o.DefaultValuesFile) > 0 {
			o.saveDefaultValuesFileOutput(pkg)
		}

		headers = append(headers, []uitable.Header{
			uitable.NewHeader("Version"),
			uitable.NewHeader("Released at"),
			uitable.NewHeader("Min capacity requirements"),
			uitable.NewHeader("Release notes"),
			uitable.NewHeader("Licenses"),
		}...)

		row = append(row, []uitable.Value{
			uitable.NewValueString(pkg.Spec.Version),
			uitable.NewValueString(formatTimestamp(pkg.Spec.ReleasedAt)),
			uitable.NewValueString(wordwrap.WrapString(pkg.Spec.CapactiyRequirementsDescription, 80)),
			uitable.NewValueString(wordwrap.WrapString(pkg.Spec.ReleaseNotes, 80)),
			uitable.NewValueStrings(pkg.Spec.Licenses),
		}...)
	} else {
		if len(o.DefaultValuesFile) > 0 {
			return fmt.Errorf("Package version is required when --default-values-file-output flag is declared (hint: to specify a particular version use the format: '-p <package-name>/<version>')")
		}
		listOpts := metav1.ListOptions{}
		if len(o.Name) > 0 {
			listOpts.FieldSelector = fields.Set{"spec.refName": o.Name}.String()
		}

		pkgList, err = client.DataV1alpha1().Packages(
			o.NamespaceFlags.Name).List(context.Background(), listOpts)
		if err != nil {
			return err
		}

		if pkgMetadata == nil && len(pkgList.Items) == 0 {
			return fmt.Errorf("Package '%s' not found in namespace '%s'", o.Name, o.NamespaceFlags.Name)
		}
	}

	table := uitable.Table{
		Transpose: true,
		Header:    headers,
		Rows:      [][]uitable.Value{row},
	}

	o.ui.PrintTable(table)

	if pkgVersion == "" {
		return o.showVersions(pkgList)
	}

	return nil
}

func (o *GetOptions) showVersions(pkgList *v1alpha1.PackageList) error {
	table := uitable.Table{
		Header: []uitable.Header{
			uitable.NewHeader("Version"),
			uitable.NewHeader("Released at"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
		},
	}

	for _, pkg := range pkgList.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(pkg.Spec.Version),
			uitable.NewValueString(formatTimestamp(pkg.Spec.ReleasedAt)),
		})
	}

	o.ui.PrintTable(table)

	return nil
}

func (o *GetOptions) showValuesSchema(client pkgclient.Interface, pkgName, pkgVersion string) error {
	pkg, err := client.DataV1alpha1().Packages(o.NamespaceFlags.Name).Get(
		context.Background(), fmt.Sprintf("%s.%s", pkgName, pkgVersion), metav1.GetOptions{})
	if err != nil {
		return err
	}

	if len(pkg.Spec.ValuesSchema.OpenAPIv3.Raw) == 0 {
		o.ui.PrintLinef("Package '%s/%s' does not have any user configurable values in the '%s' namespace", pkgName, pkgVersion, o.NamespaceFlags.Name)
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
		Title: fmt.Sprintf("Values schema for '%s/%s'", pkgName, pkgVersion),

		Header: []uitable.Header{
			uitable.NewHeader("Key"),
			uitable.NewHeader("Default"),
			uitable.NewHeader("Type"),
			uitable.NewHeader("Description"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
		},
	}

	for _, v := range parsedProperties {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(v.Key),
			uitable.NewValueInterface(v.Default),
			uitable.NewValueString(v.Type),
			uitable.NewValueString(v.Description),
		})
	}

	o.ui.PrintTable(table)

	return err
}

func (o *GetOptions) saveDefaultValuesFileOutput(pkg *v1alpha1.Package) error {
	if len(pkg.Spec.ValuesSchema.OpenAPIv3.Raw) == 0 {
		o.ui.PrintLinef("Package '%s/%s' does not have any user configurable values in the '%s' namespace", pkg.Spec.RefName, pkg.Spec.Version, o.NamespaceFlags.Name)
		return nil
	}

	s := PackageSchema{pkg.Spec.ValuesSchema.OpenAPIv3.Raw}
	defaultValues, err := s.DefaultValues()
	if err != nil {
		return err
	}

	err = os.WriteFile(o.DefaultValuesFile, defaultValues, 0600)
	if err != nil {
		return fmt.Errorf("Writing default values: %s", err)
	}

	o.ui.PrintLinef("Created default values file at %s", o.DefaultValuesFile)

	return nil
}
