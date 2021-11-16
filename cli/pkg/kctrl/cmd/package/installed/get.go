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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	valuesFile string
}

func NewGetOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *GetOptions {
	return &GetOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewGetCmd(o *GetOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get details for installed package",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVar(&o.Name, "package-install", "i", "Set installed package name")
	cmd.Flags().StringVar(&o.valuesFile, "values-file", "", "File path for exporting configuration values file")
	return cmd
}

func (o *GetOptions) Run() error {
	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	pkgi, err := client.PackagingV1alpha1().PackageInstalls(
		o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if o.valuesFile != "" {
		f, err := os.Create(o.valuesFile)
		if err != nil {
			return err
		}
		defer f.Close()
		w := bufio.NewWriter(f)

		coreClient, err := o.depsFactory.CoreClient()
		if err != nil {
			return err
		}

		if len(pkgi.Spec.Values) != 1 {
			return fmt.Errorf("Expected 1 values reference, found %d", len(pkgi.Spec.Values))
		}

		if pkgi.Spec.Values[0].SecretRef == nil {
			return fmt.Errorf("Values do not reference a Secret")
		}

		secretName := pkgi.Spec.Values[0].SecretRef.Name
		valuesSecret, err := coreClient.CoreV1().Secrets(o.NamespaceFlags.Name).Get(context.Background(), secretName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		data, ok := valuesSecret.Data[valuesFileKey]
		if !ok {
			// TODO: Add hint saying that install was not created by this CLI
			return fmt.Errorf("Expected to find key")
		}

		w.Write(data)
		w.Flush()
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
