// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"context"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/k14s/kapp/pkg/kapp/cmd/core"
	"github.com/k14s/kapp/pkg/kapp/logger"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger
	pkgiName    string
	//valuesFile  string

	NamespaceFlags cmdcore.NamespaceFlags
	//AllNamespaces  bool
}

func NewGetOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *GetOptions {
	return &GetOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewGetCmd(o *GetOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get installed Package",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVar(&o.pkgiName, "name", "", "Name of PackageInstall")
	//cmd.Flags().StringVar(&o.valuesFile, "values-file", "", "File path for exporting configuration values file")
	//cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List installed Packages in all namespaces")
	return cmd
}

func (o *GetOptions) Run() error {

	// if o.AllNamespaces {
	// 	o.NamespaceFlags.Name = ""
	// }

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	pkgi, err := client.PackagingV1alpha1().PackageInstalls(
		o.NamespaceFlags.Name).Get(context.Background(), o.pkgiName, metav1.GetOptions{})
	if err != nil {
		//Handle IsNotFound error?
		return err
	}

	// if o.valuesFile != "" {
	// 	f, err := os.Create(o.valuesFile)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer f.Close()
	// 	w := bufio.NewWriter(f)

	// 	coreClient := o.depsFactory.CoreClient()

	// 	dataValue := ""
	// 	for _, value := range pkgi.Spec.Values {
	// 		if value.SecretRef == nil {
	// 			continue
	// 		}
	// 		s, err := coreClient.Secrets(o.NamespaceFlags.Name).Get(context.Background(), value.SecretRef.Name, metav1.GetOptions{})
	// 		if err != nil {
	// 			return err
	// 		}

	// 		var data []byte
	// 		for _, value := range s.Data {
	// 			if len(string(value)) < 3 {
	// 				data = append(data, tkgpackagedatamodel.YamlSeparator...)
	// 				data = append(data, "\n"...)
	// 			}
	// 			if len(string(value)) >= 3 && string(value)[:3] != tkgpackagedatamodel.YamlSeparator {
	// 				data = append(data, tkgpackagedatamodel.YamlSeparator...)
	// 				data = append(data, "\n"...)
	// 			}
	// 			data = append(data, value...)
	// 		}

	// 		if len(string(s)) < 3 {
	// 			dataValue += tkgpackagedatamodel.YamlSeparator
	// 			dataValue += "\n"
	// 		}
	// 		if len(string(s)) >= 3 && string(s)[:3] != tkgpackagedatamodel.YamlSeparator {
	// 			dataValue += tkgpackagedatamodel.YamlSeparator
	// 			dataValue += "\n"
	// 		}
	// 		dataValue += string(s)
	// 	}
	// 	if _, err = fmt.Fprintf(w, "%s", dataValue); err != nil {
	// 		return err
	// 	}
	// 	w.Flush()
	// 	return nil
	// }

	o.ui.PrintLinef(
		`
NAME:                 %s
PACKAGE NAME:         %s
PACKAGE VERSION:      %s
STATUS:               %s
CONDITIONS:           %s
USEFUL ERROR MESSAGE: %s
	`,
		pkgi.Name, pkgi.Spec.PackageRef.RefName, pkgi.Status.Version, pkgi.Status.FriendlyDescription,
		pkgi.Status.Conditions, pkgi.Status.UsefulErrorMessage,
	)

	return nil
}
