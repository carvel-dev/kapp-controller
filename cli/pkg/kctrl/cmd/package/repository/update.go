// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"context"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type UpdateOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string
	URL            string

	CreateRepository bool
	CreateNamespace  bool

	Wait         bool
	PollInterval time.Duration
	PollTimeout  time.Duration
}

func NewUpdateOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *UpdateOptions {
	return &UpdateOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewUpdateCmd(o *UpdateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a package repository",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)

	cmd.Flags().StringVarP(&o.Name, "repository", "r", "", "Set package repository name")
	cmd.Flags().StringVarP(&o.URL, "url", "", "", "OCI registry url for package repository bundle")
	cmd.MarkFlagRequired("url")

	cmd.Flags().BoolVar(&o.CreateRepository, "create", false, "Creates the package repository if it does not exist, optional")
	cmd.Flags().BoolVar(&o.CreateNamespace, "create-namespace", false, "Create namespace if the target namespace does not exist, optional")

	cmd.Flags().BoolVar(&o.Wait, "wait", true, "Wait for the package repository reconciliation to complete, optional. To disable wait, specify --wait=false")
	cmd.Flags().DurationVar(&o.PollInterval, "poll-interval", 1*time.Second, "Time interval between subsequent polls of package repository reconciliation status, optional")
	cmd.Flags().DurationVar(&o.PollTimeout, "poll-timeout", 5*time.Minute, "Timeout value for polls of package repository reconciliation status, optional")

	return cmd
}

func (o *UpdateOptions) Run() error {
	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	kappClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	if o.CreateNamespace {
		_, err := kappClient.CoreV1().Namespaces().Create(context.Background(),
			&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: o.NamespaceFlags.Name}}, metav1.CreateOptions{})
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}

	existingRepository, err := client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Get(
		context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) && o.CreateRepository {
			pkgRepository, err := newPackageRepository(o.Name, o.URL, o.NamespaceFlags.Name)
			if err != nil {
				return err
			}

			_, err = client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Create(
				context.Background(), pkgRepository, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		}

		return err
	}

	pkgRepository, err := updateExistingPackageRepository(existingRepository, o.URL)
	if err != nil {
		return err
	}

	_, err = client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Update(
		context.Background(), pkgRepository, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	if o.Wait {
		o.ui.PrintLinef("Waiting for package repository to be added/updated")
		err = waitForPackageRepositoryInstallation(o.PollInterval, o.PollTimeout, o.NamespaceFlags.Name, o.Name, client)
	}

	return err
}
