// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type AddOrUpdateOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string
	URL            string

	CreateRepository bool

	WaitFlags cmdcore.WaitFlags

	positionalNameArg bool
}

func NewAddOrUpdateOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, positionalNameArg bool) *AddOrUpdateOptions {
	return &AddOrUpdateOptions{ui: ui, depsFactory: depsFactory, logger: logger, positionalNameArg: positionalNameArg}
}

func NewAddCmd(o *AddOrUpdateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a package repository",
		RunE:  func(_ *cobra.Command, args []string) error { return o.Run(args) },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)

	if !o.positionalNameArg {
		cmd.Flags().StringVarP(&o.Name, "repository", "r", "", "Set package repository name")
	}

	// TODO consider how to support other repository types
	cmd.Flags().StringVar(&o.URL, "url", "", "OCI registry url for package repository bundle")
	cmd.MarkFlagRequired("url")

	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   5 * time.Minute,
	})

	// For `add` command create option will always be true
	o.CreateRepository = true

	return cmd
}

func NewUpdateCmd(o *AddOrUpdateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a package repository",
		RunE:  func(_ *cobra.Command, args []string) error { return o.Run(args) },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)

	if !o.positionalNameArg {
		cmd.Flags().StringVarP(&o.Name, "repository", "r", "", "Set package repository name")
	}

	cmd.Flags().StringVarP(&o.URL, "url", "", "", "OCI registry url for package repository bundle")
	cmd.MarkFlagRequired("url")

	cmd.Flags().BoolVar(&o.CreateRepository, "create", false, "Creates the package repository if it does not exist, optional")

	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   5 * time.Minute,
	})

	return cmd
}

func (o *AddOrUpdateOptions) Run(args []string) error {
	if o.positionalNameArg {
		o.Name = args[0]
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	existingRepository, err := client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Get(
		context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") && o.CreateRepository {
			return o.add(client)
		}
		return err
	}

	pkgRepository, err := o.updateExistingPackageRepository(existingRepository)
	if err != nil {
		return err
	}

	_, err = client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Update(
		context.Background(), pkgRepository, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	if o.WaitFlags.Enabled {
		o.ui.PrintLinef("Waiting for package repository to be updated")
		err = o.waitForPackageRepositoryInstallation(client)
	}

	return err
}

func (o *AddOrUpdateOptions) add(client kcclient.Interface) error {
	pkgRepository, err := o.newPackageRepository()
	if err != nil {
		return err
	}

	_, err = client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Create(
		context.Background(), pkgRepository, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	if o.WaitFlags.Enabled {
		o.ui.PrintLinef("Waiting for package repository to be added")
		err = o.waitForPackageRepositoryInstallation(client)
	}

	return err
}

func (o *AddOrUpdateOptions) newPackageRepository() (*v1alpha1.PackageRepository, error) {
	pkgr := &v1alpha1.PackageRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: o.NamespaceFlags.Name,
		},
	}

	pkgr.Spec = kappipkg.PackageRepositorySpec{
		Fetch: &kappipkg.PackageRepositoryFetch{
			ImgpkgBundle: &kappctrl.AppFetchImgpkgBundle{Image: o.URL},
		},
	}

	return o.updateExistingPackageRepository(pkgr)
}

func (o *AddOrUpdateOptions) updateExistingPackageRepository(pkgr *v1alpha1.PackageRepository) (*v1alpha1.PackageRepository, error) {

	pkgr = pkgr.DeepCopy()

	pkgr.Spec.Fetch.ImgpkgBundle.Image = o.URL

	ref, err := name.ParseReference(o.URL, name.WeakValidation)
	if err != nil {
		return pkgr, fmt.Errorf("Parsing OCI registry URL: %s", err)
	}

	tag := ref.Identifier()

	// the parser function sets the tag to "latest" if not specified, however we want it to be empty
	if tag == "latest" && !strings.HasSuffix(o.URL, ":"+"latest") {
		tag = ""
	}

	if tag == "" {
		pkgr.Spec.Fetch.ImgpkgBundle.TagSelection = &versions.VersionSelection{
			Semver: &versions.VersionSelectionSemver{},
		}
	}

	return pkgr, err
}

func (o *AddOrUpdateOptions) waitForPackageRepositoryInstallation(client kcclient.Interface) error {
	msgsUI := cmdcore.NewDedupingMessagesUI(cmdcore.NewPlainMessagesUI(o.ui))
	description := getPackageRepositoryDescription(o.Name, o.NamespaceFlags.Name)
	if err := wait.Poll(o.WaitFlags.CheckInterval, o.WaitFlags.Timeout, func() (done bool, err error) {
		pkgr, err := client.PackagingV1alpha1().PackageRepositories(
			o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if pkgr.Generation != pkgr.Status.ObservedGeneration {
			// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
			return false, nil
		}

		status := pkgr.Status.GenericStatus

		for _, condition := range status.Conditions {
			msgsUI.NotifySection("%s: %s", description, condition.Type)

			switch {
			case condition.Type == kappctrl.ReconcileSucceeded && condition.Status == corev1.ConditionTrue:
				return true, nil
			case condition.Type == kappctrl.ReconcileFailed && condition.Status == corev1.ConditionTrue:
				return false, fmt.Errorf("%s. %s", status.UsefulErrorMessage, status.FriendlyDescription)
			}
		}
		return false, nil
	}); err != nil {
		return fmt.Errorf("%s: Reconciliation failed: %s", description, err)
	}

	return nil
}

func getPackageRepositoryDescription(name string, namespace string) string {
	description := fmt.Sprintf("packagerepository/%s (packaging.carvel.dev/v1alpha1)", name)
	if len(namespace) > 0 {
		description += " namespace: " + namespace
	} else {
		description += " cluster"
	}
	return description
}