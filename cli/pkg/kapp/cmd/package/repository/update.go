// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"context"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/k14s/kapp/pkg/kapp/cmd/core"

	"github.com/k14s/kapp/pkg/kapp/logger"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type UpdateOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags   cmdcore.NamespaceFlags
	RepositoryName   string
	RepositoryURL    string
	CreateRepository bool
	CreateNamespace  bool
	Wait             bool
	PollInterval     time.Duration
	PollTimeout      time.Duration
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
	cmd.Flags().StringVarP(&o.RepositoryName, "repository", "R", "", "Add a package repository")
	cmd.Flags().StringVarP(&o.RepositoryURL, "url", "", "", "OCI registry url for package repository bundle")
	cmd.Flags().BoolVarP(&o.CreateRepository, "create", "", false, "Creates the package repository if it does not exist, optional")
	cmd.Flags().BoolVarP(&o.CreateNamespace, "create-namespace", "", false, "Create namespace if the target namespace does not exist, optional")
	cmd.Flags().BoolVarP(&o.Wait, "wait", "", true, "Wait for the package repository reconciliation to complete, optional. To disable wait, specify --wait=false")
	cmd.Flags().DurationVarP(&o.PollInterval, "poll-interval", "", 1*time.Second, "Time interval between subsequent polls of package repository reconciliation status, optional")
	cmd.Flags().DurationVarP(&o.PollTimeout, "poll-timeout", "", 5*time.Minute, "Timeout value for polls of package repository reconciliation status, optional")
	cmd.MarkFlagRequired("url")
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
			&corev1.Namespace{ObjectMeta: v1.ObjectMeta{Name: o.NamespaceFlags.Name}}, v1.CreateOptions{})
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}

	existingRepository, err := client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Get(
		context.Background(), o.RepositoryName, v1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) && o.CreateRepository {
			pkgRepository, err := newPackageRepository(o.RepositoryName, o.RepositoryURL, o.NamespaceFlags.Name)
			if err != nil {
				return err
			}
			_, err = client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Create(
				context.Background(), pkgRepository, v1.CreateOptions{})
			if err != nil {
				return err
			}
		}

		return err

	} else {
		pkgRepository, err := updateExistingPackageRepoository(existingRepository, o.RepositoryName, o.RepositoryURL, o.NamespaceFlags.Name)
		if err != nil {
			return err
		}
		_, err = client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Update(
			context.Background(), pkgRepository, v1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	if o.Wait {
		o.ui.PrintLinef("Waiting for package repository to be added/updated")
		err = waitForPackageRepositoryInstallation(o.PollInterval, o.PollTimeout, o.NamespaceFlags.Name, o.RepositoryName, client)
	}

	return err
}
