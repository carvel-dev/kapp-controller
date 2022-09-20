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
	cmdapp "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcpkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AddOrUpdateOptions struct {
	ui          ui.UI
	statusUI    cmdcore.StatusLoggingUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags       cmdcore.NamespaceFlags
	SecureNamespaceFlags cmdcore.SecureNamespaceFlags
	Name                 string
	URL                  string

	CreateRepository bool

	WaitFlags cmdcore.WaitFlags

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewAddOrUpdateOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *AddOrUpdateOptions {
	return &AddOrUpdateOptions{ui: ui, statusUI: cmdcore.NewStatusLoggingUI(ui), depsFactory: depsFactory, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewAddCmd(o *AddOrUpdateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a package repository",
		RunE:  func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Add a package repository",
				[]string{"package", "repository", "add", "-r", "tce", "--url", "projects.registry.vmware.com/tce/main:0.9.1"}},
		}.Description("-r", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations: map[string]string{cmdapp.TTYByDefaultKey: "",
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value},
	}

	o.NamespaceFlags.SetWithPackageCommandTreeOpts(cmd, flagsFactory, o.pkgCmdTreeOpts)
	o.SecureNamespaceFlags.Set(cmd)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "repository", "r", "", "Set package repository name (required)")
	} else {
		cmd.Use = "add REPOSITORY_NAME --url REPOSITORY_URL"
		cmd.Args = cobra.ExactArgs(1)
	}

	// TODO consider how to support other repository types
	cmd.Flags().StringVar(&o.URL, "url", "", "OCI registry url for package repository bundle (required)")

	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   5 * time.Minute,
	})

	o.CreateRepository = true

	return cmd
}

func NewUpdateCmd(o *AddOrUpdateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a package repository",
		RunE:  func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Update a package repository with a new URL",
				[]string{"package", "repository", "update", "-r", "tce", "--url", "projects.registry.vmware.com/tce/main:0.9.2"}},
		}.Description("-r", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations: map[string]string{cmdapp.TTYByDefaultKey: "",
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value},
	}

	o.NamespaceFlags.SetWithPackageCommandTreeOpts(cmd, flagsFactory, o.pkgCmdTreeOpts)
	o.SecureNamespaceFlags.Set(cmd)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "repository", "r", "", "Set package repository name (required)")
	} else {
		cmd.Use = "update REPOSITORY_NAME --url REPOSITORY_URL"
		cmd.Args = cobra.ExactArgs(1)
	}

	cmd.Flags().StringVarP(&o.URL, "url", "", "", "OCI registry url for package repository bundle (required)")

	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   5 * time.Minute,
	})

	return cmd
}

func (o *AddOrUpdateOptions) Run(args []string) error {
	if o.pkgCmdTreeOpts.PositionalArgs {
		o.Name = args[0]
	}

	if len(o.Name) == 0 {
		return fmt.Errorf("Expected package repository name to be non-empty")
	}

	if len(o.URL) == 0 {
		return fmt.Errorf("Expected package repository url to be non-empty")
	}

	err := o.SecureNamespaceFlags.CheckForDisallowedSharedNamespaces(o.NamespaceFlags.Name)
	if err != nil {
		return err
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	existingRepository, err := client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Get(
		context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) && o.CreateRepository {
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

func (o *AddOrUpdateOptions) newPackageRepository() (*kcpkg.PackageRepository, error) {
	pkgr := &kcpkg.PackageRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: o.NamespaceFlags.Name,
		},
		Spec: kcpkg.PackageRepositorySpec{},
	}

	return o.updateExistingPackageRepository(pkgr)
}

func (o *AddOrUpdateOptions) updateExistingPackageRepository(pkgr *kcpkg.PackageRepository) (*kcpkg.PackageRepository, error) {

	pkgr = pkgr.DeepCopy()

	pkgr.Spec.Fetch = &kcpkg.PackageRepositoryFetch{
		ImgpkgBundle: &kappctrl.AppFetchImgpkgBundle{Image: o.URL},
	}

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
	o.statusUI.PrintMessagef("Waiting for package repository reconciliation for '%s'", o.Name)
	repoWatcher := NewRepoTailer(o.NamespaceFlags.Name, o.Name, o.ui, client)

	err := repoWatcher.TailRepoStatus()
	if err != nil {
		return err
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
