// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	valuesFileOutput string
	values           bool

	binaryName        string
	positionalNameArg bool
}

func NewGetOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, binaryName string, positionalNameArg bool) *GetOptions {
	return &GetOptions{ui: ui, depsFactory: depsFactory, logger: logger, binaryName: binaryName, positionalNameArg: positionalNameArg}
}

func NewGetCmd(o *GetOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get details for installed package",
		Args:    cobra.ExactArgs(1),
		RunE:    func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Get details for a package install",
				[]string{"package", "installed", "get", "-i", "cert-man"},
			},
			cmdcore.Example{"View values being used by package install",
				[]string{"package", "installed", "get", "-i", "cert-man", "--values"},
			},
			cmdcore.Example{"Download values being used by package install",
				[]string{"package", "installed", "get", "-i", "cert-man", "--values-file-output", "values.yml"}},
		}.Description(o.binaryName, "-i", o.positionalNameArg),
		SilenceUsage: true,
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)

	if !o.positionalNameArg {
		cmd.Flags().StringVarP(&o.Name, "package-install", "i", "", "Set installed package name (required)")
	} else {
		cmd.Use = "get INSTALLED_PACKAGE_NAME"
	}

	cmd.Flags().StringVar(&o.valuesFileOutput, "values-file-output", "", "File path for exporting configuration values file")
	cmd.Flags().BoolVar(&o.values, "values", false, "Get values data for pacakge install")
	return cmd
}

func (o *GetOptions) Run(args []string) error {
	if o.positionalNameArg {
		o.Name = args[0]
	}

	if len(o.Name) == 0 {
		return fmt.Errorf("Expected package install name to be non empty")
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	pkgi, err := client.PackagingV1alpha1().PackageInstalls(
		o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if o.valuesFileOutput != "" {
		err := o.downloadValuesData(pkgi)
		if err != nil {
			return err
		}
		return nil
	}

	if o.values {
		err := o.showValuesData(pkgi)
		if err != nil {
			return err
		}
		return nil
	}

	table := uitable.Table{
		Transpose: true,

		Header: []uitable.Header{
			uitable.NewHeader("Namespace"),
			uitable.NewHeader("Name"),
			uitable.NewHeader("Package name"),
			uitable.NewHeader("Package version"),
			uitable.NewHeader("Description"),
			uitable.NewHeader("Conditions"),
			uitable.NewHeader("Useful error message"),
		},

		Rows: [][]uitable.Value{{
			uitable.NewValueString(pkgi.Namespace),
			uitable.NewValueString(pkgi.Name),
			uitable.NewValueString(pkgi.Spec.PackageRef.RefName),
			uitable.NewValueString(pkgi.Status.Version),
			uitable.NewValueString(pkgi.Status.FriendlyDescription),
			uitable.NewValueInterface(pkgi.Status.Conditions),
			uitable.NewValueString(pkgi.Status.UsefulErrorMessage),
		}},
	}

	o.ui.PrintTable(table)

	return nil
}

func (o *GetOptions) getSecretData(pkgi *kcpkgv1alpha1.PackageInstall) ([]byte, error) {

	if len(pkgi.Spec.Values) != 1 {
		return nil, fmt.Errorf("Expected 1 values reference, found %d", len(pkgi.Spec.Values))
	}

	if pkgi.Spec.Values[0].SecretRef == nil {
		return nil, fmt.Errorf("Values do not reference a Secret")
	}

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return nil, err
	}

	secretName := pkgi.Spec.Values[0].SecretRef.Name
	valuesSecret, err := coreClient.CoreV1().Secrets(o.NamespaceFlags.Name).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if len(valuesSecret.Data) > 1 {
		return nil, fmt.Errorf("Manually created values Secret has multiple keys")
	}

	// To get values data from any single key that is present
	var dataKey string
	for key := range valuesSecret.Data {
		dataKey = key
	}

	data, ok := valuesSecret.Data[dataKey]
	if !ok {
		return nil, fmt.Errorf("Could not find key with values data in referenced secret")
	}

	return data, nil
}

func (o *GetOptions) downloadValuesData(pkgi *kcpkgv1alpha1.PackageInstall) error {
	data, err := o.getSecretData(pkgi)
	if err != nil {
		return err
	}

	f, err := os.Create(o.valuesFileOutput)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	w.Write(data)
	w.Flush()

	return nil
}

func (o *GetOptions) showValuesData(pkgi *kcpkgv1alpha1.PackageInstall) error {
	data, err := o.getSecretData(pkgi)
	if err != nil {
		return err
	}

	o.ui.PrintBlock(data)

	return nil
}
