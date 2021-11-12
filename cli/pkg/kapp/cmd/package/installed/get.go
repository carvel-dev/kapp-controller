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

	valuesFile string
}

func NewGetOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *GetOptions {
	return &GetOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewGetCmd(o *GetOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get installed package",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVar(&o.Name, "name", "", "Name of PackageInstall")
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

	//TODO: Verify and enhance how we export values file
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

		dataValue := ""
		for _, value := range pkgi.Spec.Values {
			if value.SecretRef == nil {
				continue
			}

			_, err := coreClient.CoreV1().Secrets(o.NamespaceFlags.Name).Get(context.Background(), value.SecretRef.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			//TODO: Add secret yaml to a buffer which can be written to a file
		}
		if _, err = fmt.Fprintf(w, "%s", dataValue); err != nil {
			return err
		}
		w.Flush()
		return nil
	}

	tableTitle := "Package Information"
	table := uitable.Table{
		Title:     tableTitle,
		Content:   "PackageInstalls",
		Transpose: true,

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Package Name"),
			uitable.NewHeader("Package Version"),
			uitable.NewHeader("Status"),
			uitable.NewHeader("Conditions"),
			uitable.NewHeader("Useful Error Message"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	table.Rows = append(table.Rows, []uitable.Value{
		uitable.NewValueString(pkgi.Name),
		uitable.NewValueString(pkgi.Spec.PackageRef.RefName),
		uitable.NewValueString(pkgi.Status.Version),
		uitable.NewValueString(pkgi.Status.FriendlyDescription),
		uitable.NewValueInterface(pkgi.Status.Conditions),
		uitable.NewValueString(pkgi.Status.UsefulErrorMessage),
	})

	o.ui.PrintTable(table)

	return nil
}
