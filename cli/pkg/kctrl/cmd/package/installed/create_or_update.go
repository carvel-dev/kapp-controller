// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
	cmdapp "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	pkgclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

const (
	valuesFileKey        = "values.yaml"
	yttOverlayPrefix     = "ext.packaging.carvel.dev/ytt-paths-from-secret-name"
	yttOverlayAnnotation = yttOverlayPrefix + ".kctrl-ytt-overlays"
)

type CreateOrUpdateOptions struct {
	ui          ui.UI
	statusUI    cmdcore.StatusLoggingUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	WaitFlags       cmdcore.WaitFlags
	YttOverlayFlags YttOverlayFlags

	packageName        string
	version            string
	valuesFile         string
	values             bool
	serviceAccountName string

	install bool

	Name                 string
	NamespaceFlags       cmdcore.NamespaceFlags
	SecureNamespaceFlags cmdcore.SecureNamespaceFlags
	createdAnnotations   *CreatedResourceAnnotations

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewCreateOrUpdateOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *CreateOrUpdateOptions {
	return &CreateOrUpdateOptions{ui: ui, statusUI: cmdcore.NewStatusLoggingUI(ui), depsFactory: depsFactory, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewCreateCmd(o *CreateOrUpdateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Install package",
		RunE:  func(_ *cobra.Command, args []string) error { return o.RunCreate(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Install a package",
				[]string{"package", "installed", "create", "-i", "cert-man", "-p", "cert-manager.community.tanzu.vmware.com", "--version", "1.6.1"},
			},
			cmdcore.Example{"Install package with values file",
				[]string{"package", "installed", "create", "-i", "cert-man", "-p", "cert-manager.community.tanzu.vmware.com", "--version", "1.6.1", "--values-file", "values.yml"},
			},
			cmdcore.Example{"Install package and ask it to use an existing service account",
				[]string{"package", "installed", "create", "-i", "cert-man", "-p", "cert-manager.community.tanzu.vmware.com", "--version", "1.6.1", "--service-account-name", "existing-sa"}},
		}.Description("-i", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations: map[string]string{cmdapp.TTYByDefaultKey: "",
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value},
	}
	o.NamespaceFlags.SetWithPackageCommandTreeOpts(cmd, flagsFactory, o.pkgCmdTreeOpts)
	o.SecureNamespaceFlags.Set(cmd)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "package-install", "i", "", "Set installed package name (required)")
	} else {
		cmd.Use = "create INSTALLED_PACKAGE_NAME --package PACKAGE_NAME --version VERSION"
		cmd.Args = cobra.ExactArgs(1)
	}

	cmd.Flags().StringVarP(&o.packageName, "package", "p", "", "Set package name (required)")
	cmd.Flags().StringVarP(&o.version, "version", "v", "", "Set package version (required)")
	cmd.Flags().StringVar(&o.serviceAccountName, "service-account-name", "", "Name of an existing service account used to install underlying package contents, optional")
	cmd.Flags().StringVar(&o.valuesFile, "values-file", "", "The path to the configuration values file, optional")
	cmd.Flags().BoolVar(&o.values, "values", true, "Add or keep values supplied to package install, optional")

	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   30 * time.Minute,
	})
	o.YttOverlayFlags.Set(cmd)

	return cmd
}

func NewInstallCmd(o *CreateOrUpdateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install package",
		RunE:  func(_ *cobra.Command, args []string) error { return o.RunCreate(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Install a package",
				[]string{"package", "install", "-i", "cert-man", "-p", "cert-manager.community.tanzu.vmware.com", "--version", "1.6.1"},
			},
			cmdcore.Example{"Install package with values file",
				[]string{"package", "install", "-i", "cert-man", "-p", "cert-manager.community.tanzu.vmware.com", "--version", "1.6.1", "--values-file", "values.yml"},
			},
			cmdcore.Example{"Install package and ask it to use an existing service account",
				[]string{"package", "install", "-i", "cert-man", "-p", "cert-manager.community.tanzu.vmware.com", "--version", "1.6.1", "--service-account-name", "existing-sa"}},
		}.Description("-i", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations: map[string]string{cmdapp.TTYByDefaultKey: "",
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value},
	}
	o.NamespaceFlags.SetWithPackageCommandTreeOpts(cmd, flagsFactory, o.pkgCmdTreeOpts)
	o.SecureNamespaceFlags.Set(cmd)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "package-install", "i", "", "Set installed package name (required)")
	} else {
		cmd.Use = "install INSTALLED_PACKAGE_NAME --package PACKAGE_NAME --version VERSION"
		cmd.Args = cobra.ExactArgs(1)
	}

	cmd.Flags().StringVarP(&o.packageName, "package", "p", "", "Set package name (required)")
	cmd.Flags().StringVarP(&o.version, "version", "v", "", "Set package version (required)")
	cmd.Flags().StringVar(&o.serviceAccountName, "service-account-name", "", "Name of an existing service account used to install underlying package contents, optional")
	cmd.Flags().StringVar(&o.valuesFile, "values-file", "", "The path to the configuration values file, optional")
	cmd.Flags().BoolVar(&o.values, "values", true, "Add or keep values supplied to package install, optional")

	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   30 * time.Minute,
	})
	o.YttOverlayFlags.Set(cmd)

	return cmd
}

func NewUpdateCmd(o *CreateOrUpdateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update package",
		RunE:  func(_ *cobra.Command, args []string) error { return o.RunUpdate(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Upgrade package install to a newer version",
				[]string{"package", "installed", "update", "-i", "cert-man", "--version", "1.6.2"},
			},
			cmdcore.Example{"Update package install with new values file",
				[]string{"package", "installed", "update", "-i", "cert-man", "--values-file", "values.yml"}},
			cmdcore.Example{"Update package install to stop consuming supplied values",
				[]string{"package", "installed", "update", "-i", "cert-man", "--values", "false"}},
		}.Description("-i", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations: map[string]string{cmdapp.TTYByDefaultKey: "",
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value},
	}
	o.NamespaceFlags.SetWithPackageCommandTreeOpts(cmd, flagsFactory, o.pkgCmdTreeOpts)
	o.SecureNamespaceFlags.Set(cmd)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "package-install", "i", "", "Set installed package name")
	} else {
		cmd.Use = "update INSTALLED_PACKAGE_NAME"
		cmd.Args = cobra.ExactArgs(1)
	}

	cmd.Flags().StringVarP(&o.packageName, "package", "p", "", "Name of package install to be updated")
	cmd.Flags().StringVarP(&o.version, "version", "v", "", "Set package version")
	cmd.Flags().StringVar(&o.valuesFile, "values-file", "", "The path to the configuration values file, optional")
	cmd.Flags().BoolVar(&o.values, "values", true, "Add or keep values supplied to package install, optional")

	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   30 * time.Minute,
	})
	o.YttOverlayFlags.Set(cmd)

	return cmd
}

func (o *CreateOrUpdateOptions) RunCreate(args []string) error {
	if o.pkgCmdTreeOpts.PositionalArgs {
		o.Name = args[0]
	}

	if len(o.Name) == 0 {
		return fmt.Errorf("Expected package install name to be non empty")
	}

	if len(o.packageName) == 0 {
		return fmt.Errorf("Expected package name to be non empty")
	}

	err := o.SecureNamespaceFlags.CheckForDisallowedSharedNamespaces(o.NamespaceFlags.Name)
	if err != nil {
		return err
	}

	if len(o.version) == 0 {
		pkgClient, err := o.depsFactory.PackageClient()
		if err != nil {
			return err
		}

		err = o.showVersions(pkgClient)
		if err != nil {
			return err
		}
		return fmt.Errorf("Expected package version to be non empty")
	}

	client, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	kcClient, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	pkgInstall, err := kcClient.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Get(
		context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	o.createdAnnotations = NewCreatedResourceAnnotations(o.Name, o.NamespaceFlags.Name)

	// Fallback to update if resource exists
	if pkgInstall != nil && err == nil {
		err = o.update(client, kcClient, pkgInstall)
		if err != nil {
			return err
		}
		return nil
	}

	err = o.create(client, kcClient)
	if err != nil {
		return err
	}

	return nil
}

func (o *CreateOrUpdateOptions) create(client kubernetes.Interface, kcClient kcclient.Interface) error {
	isServiceAccountCreated, isSecretCreated, err := o.createRelatedResources(client)
	if err != nil {
		return err
	}

	overlaySecretName := ""
	if o.YttOverlayFlags.yttOverlays {
		overlaysSecret, err := o.createOrUpdateYttOverlaySecrets(nil, client)
		if err != nil {
			return err
		}
		if overlaysSecret != nil {
			overlaySecretName = overlaysSecret.Name
		}
	}

	o.statusUI.PrintMessagef("Creating package install resource")
	if err = o.createPackageInstall(isServiceAccountCreated, isSecretCreated, overlaySecretName, kcClient); err != nil {
		return err
	}

	if o.WaitFlags.Enabled {
		if err = o.waitForResourceInstallation(o.Name, o.NamespaceFlags.Name, o.WaitFlags.CheckInterval, o.WaitFlags.Timeout, kcClient); err != nil {
			return err
		}
	}

	return nil
}

func (o *CreateOrUpdateOptions) RunUpdate(args []string) error {
	if o.pkgCmdTreeOpts.PositionalArgs {
		o.Name = args[0]
	}

	if len(o.Name) == 0 {
		return fmt.Errorf("Expected package install to be non-empty")
	}

	switch {
	case len(o.version) != 0: // Continue if version is changed
	case len(o.valuesFile) != 0: // Continue if supplied values might have changed
	case !o.values: // Continue if supplied values might be getting removed
	case !o.YttOverlayFlags.yttOverlays: // Continue if supplied overlays might be getting removed
	case len(o.YttOverlayFlags.yttOverlayFiles) != 0: // Continue if supplied overlays might be getting changed
	default:
		return fmt.Errorf("Expected either package version, values file or overlays to update the package")
	}

	err := o.SecureNamespaceFlags.CheckForDisallowedSharedNamespaces(o.NamespaceFlags.Name)
	if err != nil {
		return err
	}

	client, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	kcClient, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	o.createdAnnotations = NewCreatedResourceAnnotations(o.Name, o.NamespaceFlags.Name)
	pkgInstall, err := kcClient.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Get(
		context.Background(), o.Name, metav1.GetOptions{},
	)
	if err != nil {
		return err
	}

	err = o.update(client, kcClient, pkgInstall)
	if err != nil {
		return err
	}

	return nil
}

func (o CreateOrUpdateOptions) update(client kubernetes.Interface, kcClient kcclient.Interface, pkgInstall *kcpkgv1alpha1.PackageInstall) error {
	updatedPkgInstall, changed, err := o.preparePackageInstallForUpdate(pkgInstall)
	if err != nil {
		return err
	}

	switch {
	case changed: // Continue if package install resource is changed
	case len(o.valuesFile) != 0: // Continue if supplied values might have changed
	case !o.values: // Continue if supplied values might be getting removed
	case !o.YttOverlayFlags.yttOverlays: // Continue if supplied overlays might be getting removed
	case len(o.YttOverlayFlags.yttOverlayFiles) != 0: // Continue if supplied overlays might be getting changed
	default:
		o.statusUI.PrintMessagef("No changes to package install '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
		return nil
	}

	// Pause reconciliation so that kctrl can tail status if an existing values secret is updated
	// Or when existing ytt overlay has to be updated
	reconciliationPaused := false
	if (o.valuesFile != "" && len(pkgInstall.Spec.Values) > 0) ||
		(len(o.YttOverlayFlags.yttOverlayFiles) > 0 && hasYttOverlays(pkgInstall)) {
		updatedPkgInstall, err = o.pauseReconciliation(kcClient)
		if err != nil {
			return err
		}
		err = o.waitForAppPause(kcClient)
		if err != nil {
			return err
		}
		reconciliationPaused = true
	}

	isSecretCreated, err := o.createOrUpdateValuesSecret(updatedPkgInstall, client)
	if err != nil {
		return err
	}

	isSecretDeleted := false
	if !o.values && len(updatedPkgInstall.Spec.Values) > 0 {
		isSecretDeleted, err = o.dropValuesSecret(client)
		if err != nil {
			return fmt.Errorf("Deleting values secret: %s", err.Error())
		}

		if isSecretDeleted {
			o.statusUI.PrintMessagef("Removing values secret reference from package install '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
			o.removeValuesSecretReference(updatedPkgInstall)
			changed = true
		}
	} else if !o.values && len(updatedPkgInstall.Spec.Values) == 0 {
		o.statusUI.PrintMessagef("No values have been supplied to installation '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
		if !changed {
			return nil
		}
	}

	if o.YttOverlayFlags.yttOverlays {
		overlaysSecret, err := o.createOrUpdateYttOverlaySecrets(updatedPkgInstall, client)
		if err != nil {
			return err
		}
		_, found := updatedPkgInstall.Annotations[yttOverlayAnnotation]
		if !found && overlaysSecret != nil {
			updatedPkgInstall.Annotations[yttOverlayAnnotation] = overlaysSecret.Name
			changed = true
		}
	}

	if !o.YttOverlayFlags.yttOverlays && hasYttOverlays(updatedPkgInstall) {
		err = o.dropYttOverlaySecrets(updatedPkgInstall, client)
		if err != nil {
			return err
		}
		o.removeYttOverlaysAnnotation(updatedPkgInstall)
		changed = true
	}

	if isSecretCreated || changed {
		o.statusUI.PrintMessagef("Updating package install for '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
		o.addCreatedResourceAnnotations(&updatedPkgInstall.ObjectMeta, false, isSecretCreated, isSecretDeleted)
		_, err = kcClient.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Update(
			context.Background(), updatedPkgInstall, metav1.UpdateOptions{},
		)
		if err != nil {
			err = fmt.Errorf("Updating package '%s': %s", o.Name, err.Error())
			return err
		}
	}

	if reconciliationPaused {
		err = o.unpauseReconciliation(kcClient)
		if err != nil {
			return err
		}
	}

	if o.WaitFlags.Enabled {
		if err = o.waitForResourceInstallation(o.Name, o.NamespaceFlags.Name, o.WaitFlags.CheckInterval, o.WaitFlags.Timeout, kcClient); err != nil {
			return err
		}
	}

	return nil
}

func (o *CreateOrUpdateOptions) dropValuesSecret(client kubernetes.Interface) (bool, error) {
	secretName := o.createdAnnotations.SecretAnnValue()
	pkgiIdentifier := o.createdAnnotations.PackageAnnValue()

	valuesSecret, err := client.CoreV1().Secrets(o.NamespaceFlags.Name).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			o.statusUI.PrintMessagef("Values secret '%s' not found in namespace '%s'", secretName, o.NamespaceFlags.Name)
			return true, nil
		}
		return false, fmt.Errorf("Getting values secret: %s", err.Error())
	}

	// Do not delete secrets not created by kctrl
	// TODO: Should we not error out and continue with the rest of the changes?
	annotations := valuesSecret.GetAnnotations()
	val, found := annotations[KctrlPkgAnnotation]
	if !found || val != pkgiIdentifier {
		// To support older version of Tanzu CLI. To be deprecated
		val, found = annotations[TanzuPkgAnnotation]
		if !found || val != pkgiIdentifier {
			return false, fmt.Errorf("Secret was not created by kctrl")
		}
	}

	o.statusUI.PrintMessagef("Deleting values secret '%s'", secretName)
	err = client.CoreV1().Secrets(o.NamespaceFlags.Name).Delete(context.Background(), secretName, metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (o *CreateOrUpdateOptions) removeValuesSecretReference(pkgi *kcpkgv1alpha1.PackageInstall) {
	for i, valueRef := range pkgi.Spec.Values {
		if valueRef.SecretRef.Name == o.createdAnnotations.SecretAnnValue() {
			pkgi.Spec.Values = append(pkgi.Spec.Values[:i], pkgi.Spec.Values[i+1:]...)
			return
		}
	}
}

func (o *CreateOrUpdateOptions) createRelatedResources(client kubernetes.Interface) (bool, bool, error) {
	var (
		isServiceAccountCreated bool
		isSecretCreated         bool
		err                     error
	)

	if o.serviceAccountName == "" {

		o.statusUI.PrintMessagef("Creating service account '%s'", o.createdAnnotations.ServiceAccountAnnValue())
		if isServiceAccountCreated, err = o.createOrUpdateServiceAccount(client); err != nil {
			return isServiceAccountCreated, isSecretCreated, err
		}

		o.statusUI.PrintMessagef("Creating cluster admin role '%s'", o.createdAnnotations.ClusterRoleAnnValue())
		if err := o.createOrUpdateClusterAdminRole(client); err != nil {
			return isServiceAccountCreated, isSecretCreated, err
		}

		o.statusUI.PrintMessagef("Creating cluster role binding '%s'", o.createdAnnotations.ClusterRoleBindingAnnValue())
		if err := o.createOrUpdateClusterRoleBinding(client); err != nil {
			return isServiceAccountCreated, isSecretCreated, err
		}
	} else {
		client, err := o.depsFactory.CoreClient()
		if err != nil {
			return isServiceAccountCreated, isSecretCreated, err
		}
		svcAccount, err := client.CoreV1().ServiceAccounts(o.NamespaceFlags.Name).Get(context.Background(), o.serviceAccountName, metav1.GetOptions{})
		if err != nil {
			err = fmt.Errorf("Finding service account '%s' in namespace '%s': %s", o.serviceAccountName, o.NamespaceFlags.Name, err.Error())
			return isServiceAccountCreated, isSecretCreated, err
		}

		svcAccountAnnotation, ok := svcAccount.GetAnnotations()[KctrlPkgAnnotation]

		// To support older versions of Tanzu CLI. To be deprecated
		if !ok {
			svcAccountAnnotation, ok = svcAccount.GetAnnotations()[TanzuPkgAnnotation]
		}

		if ok {
			if svcAccountAnnotation != o.createdAnnotations.PackageAnnValue() {
				err = fmt.Errorf("Provided service account '%s' is already used by another package in namespace '%s': %s", o.serviceAccountName, o.NamespaceFlags.Name, err.Error())
				return isServiceAccountCreated, isSecretCreated, err
			}
		}
	}

	if o.valuesFile != "" && o.values {
		o.statusUI.PrintMessagef("Creating secret '%s'", o.createdAnnotations.SecretAnnValue())
		if isSecretCreated, err = o.createOrUpdateDataValuesSecret(client); err != nil {
			return isServiceAccountCreated, isSecretCreated, err
		}
	}

	return isServiceAccountCreated, isSecretCreated, nil
}

func (o *CreateOrUpdateOptions) createOrUpdateClusterAdminRole(client kubernetes.Interface) error {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: o.createdAnnotations.ClusterRoleAnnValue(),
			Annotations: map[string]string{
				KctrlPkgAnnotation: o.createdAnnotations.PackageAnnValue(),
				TanzuPkgAnnotation: o.createdAnnotations.PackageAnnValue(),
			},
		},
		Rules: []rbacv1.PolicyRule{
			{APIGroups: []string{"*"}, Verbs: []string{"*"}, Resources: []string{"*"}},
		},
	}

	_, err := client.RbacV1().ClusterRoles().Create(context.Background(), clusterRole, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err := client.RbacV1().ClusterRoles().Update(context.Background(), clusterRole, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (o *CreateOrUpdateOptions) createOrUpdateClusterRoleBinding(client kubernetes.Interface) error {
	svcAccount := o.serviceAccountName
	if svcAccount == "" {
		svcAccount = o.createdAnnotations.ServiceAccountAnnValue()
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: o.createdAnnotations.ClusterRoleBindingAnnValue(),
			Annotations: map[string]string{
				KctrlPkgAnnotation: o.createdAnnotations.PackageAnnValue(),
				TanzuPkgAnnotation: o.createdAnnotations.PackageAnnValue(),
			},
		},
		Subjects: []rbacv1.Subject{{Kind: KindServiceAccount.AsString(), Name: svcAccount, Namespace: o.NamespaceFlags.Name}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Kind:     KindClusterRole.AsString(),
			Name:     o.createdAnnotations.ClusterRoleAnnValue(),
		},
	}

	_, err := client.RbacV1().ClusterRoleBindings().Create(context.Background(), clusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err = client.RbacV1().ClusterRoleBindings().Update(context.Background(), clusterRoleBinding, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (o *CreateOrUpdateOptions) createOrUpdateDataValuesSecret(client kubernetes.Interface) (bool, error) {
	var err error

	dataValues := make(map[string][]byte)

	dataValues[valuesFileKey], err = cmdcore.NewInputFile(o.valuesFile).Bytes()
	if err != nil {
		return false, fmt.Errorf("Reading data values file '%s': %s", o.valuesFile, err.Error())
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.createdAnnotations.SecretAnnValue(),
			Namespace: o.NamespaceFlags.Name,
			Annotations: map[string]string{
				KctrlPkgAnnotation: o.createdAnnotations.PackageAnnValue(),
				TanzuPkgAnnotation: o.createdAnnotations.PackageAnnValue(),
			},
		},
		Data: dataValues,
	}

	_, err = client.CoreV1().Secrets(o.NamespaceFlags.Name).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err := client.CoreV1().Secrets(o.NamespaceFlags.Name).Update(context.Background(), secret, metav1.UpdateOptions{})
			if err != nil {
				return false, err
			}
		} else {
			return false, err
		}
	}

	return true, nil
}

func (o *CreateOrUpdateOptions) createPackageInstall(serviceAccountCreated, secretCreated bool, overlaysSecretName string, kcClient kcclient.Interface) error {
	svcAccount := o.serviceAccountName
	if svcAccount == "" {
		svcAccount = o.createdAnnotations.ServiceAccountAnnValue()
	}

	// construct the PackageInstall CR
	packageInstall := &kcpkgv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{Name: o.Name, Namespace: o.NamespaceFlags.Name},
		Spec: kcpkgv1alpha1.PackageInstallSpec{
			ServiceAccountName: svcAccount,
			PackageRef: &kcpkgv1alpha1.PackageRef{
				RefName: o.packageName,
				VersionSelection: &versions.VersionSelectionSemver{
					Constraints: o.version,
					Prereleases: &versions.VersionSelectionSemverPrereleases{},
				},
			},
		},
	}

	// if configuration data file was provided, reference the secret name in the PackageInstall
	if secretCreated {
		packageInstall.Spec.Values = []kcpkgv1alpha1.PackageInstallValues{
			{
				SecretRef: &kcpkgv1alpha1.PackageInstallValuesSecretRef{
					Name: o.createdAnnotations.SecretAnnValue(),
				},
			},
		}
	}

	// Add reference to ytt overlaty annotation if overlay secret has been created
	if overlaysSecretName != "" {
		packageInstall.Annotations = map[string]string{
			yttOverlayAnnotation: overlaysSecretName,
		}
	}

	o.addCreatedResourceAnnotations(&packageInstall.ObjectMeta, serviceAccountCreated, secretCreated, false)

	_, err := kcClient.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Create(context.Background(), packageInstall, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("Creating PackageInstall resource: %s", err.Error())
	}

	return nil
}

func (o *CreateOrUpdateOptions) createOrUpdateServiceAccount(client kubernetes.Interface) (bool, error) {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.createdAnnotations.ServiceAccountAnnValue(),
			Namespace: o.NamespaceFlags.Name,
			Annotations: map[string]string{
				KctrlPkgAnnotation: o.createdAnnotations.PackageAnnValue(),
				TanzuPkgAnnotation: o.createdAnnotations.PackageAnnValue(),
			},
		},
	}

	_, err := client.CoreV1().ServiceAccounts(o.NamespaceFlags.Name).Create(context.Background(), serviceAccount, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err := client.CoreV1().ServiceAccounts(o.NamespaceFlags.Name).Update(context.Background(), serviceAccount, metav1.UpdateOptions{})
			if err != nil {
				return false, err
			}
		} else {
			return false, err
		}
	}

	return true, nil
}

func (o *CreateOrUpdateOptions) preparePackageInstallForUpdate(pkgInstall *kcpkgv1alpha1.PackageInstall) (*kcpkgv1alpha1.PackageInstall, bool, error) {
	var (
		changed bool
		err     error
	)

	updatedPkgInstall := pkgInstall.DeepCopy()

	if updatedPkgInstall.Spec.PackageRef == nil || updatedPkgInstall.Spec.PackageRef.VersionSelection == nil {
		err = fmt.Errorf("Failed to update package '%s' as no existing package reference/version was found in the package install", o.Name)
		return nil, false, err
	}

	// If o.PackageName is provided by the user (via --package flag), verify that the package name in PackageInstall matches it.
	// This will prevent the users from accidentally overwriting an installed package with another package content due to choosing a pre-existing name for the package isntall.
	// Otherwise if o.PackageName is not provided, fill it from the installed package spec
	if o.packageName != "" && updatedPkgInstall.Spec.PackageRef.RefName != o.packageName {
		err = fmt.Errorf("Installed package '%s' is already associated with package '%s'", o.Name, updatedPkgInstall.Spec.PackageRef.RefName)
		return nil, false, err
	}
	o.packageName = updatedPkgInstall.Spec.PackageRef.RefName

	// If o.Version is provided by the user (via --version flag), set the version in PackageInstall to this version
	// Otherwise if o.Version is not provided, fill it from the installed package spec
	if o.version != "" {
		if updatedPkgInstall.Spec.PackageRef.VersionSelection.Constraints != o.version {
			changed = true
			updatedPkgInstall.Spec.PackageRef.VersionSelection.Constraints = o.version
		}
	} else {
		o.version = updatedPkgInstall.Spec.PackageRef.VersionSelection.Constraints
	}

	return updatedPkgInstall, changed, nil
}

func (o *CreateOrUpdateOptions) createOrUpdateValuesSecret(pkgInstallToUpdate *kcpkgv1alpha1.PackageInstall, client kubernetes.Interface) (bool, error) {
	var (
		secretCreated bool
		err           error
	)

	if o.valuesFile == "" {
		return false, nil
	}

	secretName := o.createdAnnotations.SecretAnnValue()

	if len(pkgInstallToUpdate.Spec.Values) > 1 {
		return false, fmt.Errorf(`Expected package install to have one or no value references while updating. Please delete and install the package install with appropriate values.`)
	}

	if len(pkgInstallToUpdate.Spec.Values) == 1 && pkgInstallToUpdate.Spec.Values[0].SecretRef.Name != "" {
		secretName = pkgInstallToUpdate.Spec.Values[0].SecretRef.Name
		o.statusUI.PrintMessagef("Updating secret '%s'", secretName)
		err := o.updateDataValuesSecret(client, secretName)
		if err != nil {
			return false, fmt.Errorf("Failed to update manually referenced secret based on values file: %s", err.Error())
		}
		return secretCreated, nil
	}

	// Second condition supports older versions of Tanzu CLI. To be deprecated
	if secretName == pkgInstallToUpdate.GetAnnotations()[KctrlPkgAnnotation+"-"+KindSecret.AsString()] ||
		secretName == pkgInstallToUpdate.GetAnnotations()[TanzuPkgAnnotation+"-"+KindSecret.AsString()] {
		o.statusUI.PrintMessagef("Updating secret '%s'", secretName)
		if err = o.updateDataValuesSecret(client, secretName); err != nil {
			return false, fmt.Errorf("Failed to update secret based on values file: %s", err.Error())
		}
	} else {
		o.statusUI.PrintMessagef("Creating secret '%s'", secretName)
		if secretCreated, err = o.createOrUpdateDataValuesSecret(client); err != nil {
			return secretCreated, fmt.Errorf("Failed to create secret based on values file: %s", err.Error())
		}
	}

	pkgInstallToUpdate.Spec.Values = []kcpkgv1alpha1.PackageInstallValues{
		{SecretRef: &kcpkgv1alpha1.PackageInstallValuesSecretRef{Name: secretName}},
	}

	return secretCreated, nil
}

func (o *CreateOrUpdateOptions) updateDataValuesSecret(client kubernetes.Interface, secretName string) error {
	var err error
	dataValues := make(map[string][]byte)

	createdSecret, err := client.CoreV1().Secrets(o.NamespaceFlags.Name).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Could not find manually referenced secret '%s' in namespace '%s'", secretName, o.NamespaceFlags.Name)
	}

	if len(createdSecret.Data) > 1 {
		return fmt.Errorf("Could not safely update manually referenced secret '%s' in namespace '%s' as it has more than one data keys", secretName, o.NamespaceFlags.Name)
	}

	dataKey := valuesFileKey
	if len(createdSecret.Data) == 1 {
		for key := range createdSecret.Data {
			dataKey = key
		}
	}

	dataValues[dataKey], err = cmdcore.NewInputFile(o.valuesFile).Bytes()
	if err != nil {
		return fmt.Errorf("Reading data values file '%s': %s", o.valuesFile, err.Error())
	}

	createdSecret.Data = dataValues
	_, err = client.CoreV1().Secrets(o.NamespaceFlags.Name).Update(context.Background(), createdSecret, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Updating Secret resource: %s", err.Error())
	}

	return nil
}

func (o *CreateOrUpdateOptions) pauseReconciliation(client kcclient.Interface) (*kcpkgv1alpha1.PackageInstall, error) {
	pausePatch := []map[string]interface{}{
		{
			"op":    "add",
			"path":  "/spec/paused",
			"value": true,
		},
	}

	patchJSON, err := json.Marshal(pausePatch)
	if err != nil {
		return nil, err
	}

	o.statusUI.PrintMessagef("Pausing reconciliation for package installation '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
	pkgi, err := client.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return pkgi, nil
}

func (o *CreateOrUpdateOptions) unpauseReconciliation(client kcclient.Interface) error {
	unpausePatch := []map[string]interface{}{
		{
			"op":   "remove",
			"path": "/spec/paused",
		},
	}

	patchJSON, err := json.Marshal(unpausePatch)
	if err != nil {
		return err
	}

	o.statusUI.PrintMessagef("Resuming reconciliation for package installation '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
	_, err = client.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Waits for the App CR created by the package installation to pick up it's paused status
func (o *CreateOrUpdateOptions) waitForAppPause(client kcclient.Interface) error {
	if err := wait.Poll(o.WaitFlags.CheckInterval, o.WaitFlags.Timeout, func() (done bool, err error) {
		appResource, err := client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if appResource.Generation != appResource.Status.ObservedGeneration {
			return false, nil
		}
		if appResource.Status.FriendlyDescription == "Canceled/paused" {
			return true, nil
		}
		return false, nil
	}); err != nil {
		return fmt.Errorf("Waiting for app '%s' in namespace '%s' to be paused: %s", o.Name, o.NamespaceFlags.Name, err)
	}
	return nil
}

func (o *CreateOrUpdateOptions) addCreatedResourceAnnotations(meta *metav1.ObjectMeta, createdSvcAccount, createdSecret bool, deletedSecret bool) {
	if meta.Annotations == nil {
		meta.Annotations = make(map[string]string)
	}
	if createdSvcAccount {
		meta.Annotations[KctrlPkgAnnotation+"-"+KindClusterRole.AsString()] = o.createdAnnotations.ClusterRoleAnnValue()
		meta.Annotations[KctrlPkgAnnotation+"-"+KindClusterRoleBinding.AsString()] = o.createdAnnotations.ClusterRoleBindingAnnValue()
		meta.Annotations[KctrlPkgAnnotation+"-"+KindServiceAccount.AsString()] = o.createdAnnotations.ServiceAccountAnnValue()

		// To support older versions of Tanzu CLI. To be deprecated
		meta.Annotations[TanzuPkgAnnotation+"-"+KindClusterRole.AsString()] = o.createdAnnotations.ClusterRoleAnnValue()
		meta.Annotations[TanzuPkgAnnotation+"-"+KindClusterRoleBinding.AsString()] = o.createdAnnotations.ClusterRoleBindingAnnValue()
		meta.Annotations[TanzuPkgAnnotation+"-"+KindServiceAccount.AsString()] = o.createdAnnotations.ServiceAccountAnnValue()
	}
	if createdSecret {
		meta.Annotations[KctrlPkgAnnotation+"-"+KindSecret.AsString()] = o.createdAnnotations.SecretAnnValue()

		// To support older versions of Tanzu CLI. To be deprecated
		meta.Annotations[TanzuPkgAnnotation+"-"+KindSecret.AsString()] = o.createdAnnotations.SecretAnnValue()
	}

	if deletedSecret {
		delete(meta.Annotations, KctrlPkgAnnotation+"-"+KindSecret.AsString())

		// To support older versions of Tanzu CLI. To be deprecated
		delete(meta.Annotations, TanzuPkgAnnotation+"-"+KindSecret.AsString())
	}
}

// waitForResourceInstallation waits until the package get installed successfully or a failure happens
func (o *CreateOrUpdateOptions) waitForResourceInstallation(name, namespace string, pollInterval, pollTimeout time.Duration, client kcclient.Interface) error {
	o.statusUI.PrintMessagef("Waiting for PackageInstall reconciliation for '%s'", name)
	msgsUI := cmdcore.NewDedupingMessagesUI(cmdcore.NewPlainMessagesUI(o.ui))
	description := getPackageInstallDescription(o.Name, o.NamespaceFlags.Name)

	appStatusTailErrored := false
	tailAppStatusOutput := func(tailErrored *bool) {
		appWatcher := cmdapp.NewAppTailer(o.NamespaceFlags.Name, o.Name, o.ui, client, cmdapp.AppTailerOpts{
			IgnoreNotExists: true,
		})

		err := appWatcher.TailAppStatus()
		if err != nil {
			o.statusUI.PrintMessagef("Error tailing app: %s\n", err.Error())
			*tailErrored = true
		}
	}
	go tailAppStatusOutput(&appStatusTailErrored)

	if err := wait.Poll(pollInterval, pollTimeout, func() (done bool, err error) {

		resource, err := client.PackagingV1alpha1().PackageInstalls(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if resource.Generation != resource.Status.ObservedGeneration {
			// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
			return false, nil
		}
		status := resource.Status.GenericStatus

		for _, condition := range status.Conditions {
			if appStatusTailErrored {
				msgsUI.NotifySection("%s: %s", description, condition.Type)
			}

			switch {
			case condition.Type == kcv1alpha1.ReconcileSucceeded && condition.Status == corev1.ConditionTrue:
				return true, nil
			case condition.Type == kcv1alpha1.ReconcileFailed && condition.Status == corev1.ConditionTrue:
				return false, fmt.Errorf("%s. %s", status.UsefulErrorMessage, status.FriendlyDescription)
			}
		}
		return false, nil
	}); err != nil {
		return fmt.Errorf("%s: Reconciling: %s", description, err)
	}

	return nil
}

func (o *CreateOrUpdateOptions) showVersions(client pkgclient.Interface) error {
	listOpts := metav1.ListOptions{}
	if len(o.packageName) > 0 {
		listOpts.FieldSelector = fields.Set{"spec.refName": o.packageName}.String()
	}

	pkgList, err := client.DataV1alpha1().Packages(
		o.NamespaceFlags.Name).List(context.Background(), listOpts)
	if err != nil {
		return err
	}

	if len(pkgList.Items) == 0 {
		return fmt.Errorf("No versions of package '%s' found in namespace '%s'", o.packageName, o.NamespaceFlags.Name)
	}

	table := uitable.Table{
		Title: fmt.Sprintf("Available Versions of %s", o.packageName),
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
			uitable.NewValueString(pkg.Spec.ReleasedAt.String()),
		})
	}

	o.ui.PrintTable(table)

	return nil
}

func (o *CreateOrUpdateOptions) createOrUpdateYttOverlaySecrets(pkgi *kcpkgv1alpha1.PackageInstall, client kubernetes.Interface) (*corev1.Secret, error) {
	o.statusUI.PrintMessage("Creating overlay secrets")
	if len(o.YttOverlayFlags.yttOverlayFiles) == 0 {
		return nil, nil
	}

	if pkgi != nil {
		for annotation := range pkgi.Annotations {
			// Ensure that kctrl does not clobber existing overlays
			if strings.HasPrefix(annotation, yttOverlayPrefix) && annotation != yttOverlayAnnotation {
				return nil, fmt.Errorf("Package install has manually supplied overlays")
			}
		}
	}

	secret, err := NewYttOverlays(o.YttOverlayFlags.yttOverlayFiles, o.Name, o.NamespaceFlags.Name).OverlaysSecret()
	if err != nil {
		return nil, fmt.Errorf("Structuring overlays secret: %s", err.Error())
	}

	var createdOrUpdatedSecret *corev1.Secret
	createdOrUpdatedSecret, err = client.CoreV1().Secrets(o.NamespaceFlags.Name).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			o.statusUI.PrintMessagef("Updating existing overlay secret '%s' in namespace '%s'", secret.Name, secret.Name)
			return nil, fmt.Errorf("Creating overlays secret: %s", err.Error())
		}
		createdOrUpdatedSecret, err = client.CoreV1().Secrets(o.NamespaceFlags.Name).Update(context.Background(), secret, metav1.UpdateOptions{})
		if err != nil {
			return nil, fmt.Errorf("Updating overlays secret: %s", err.Error())
		}
	}

	return createdOrUpdatedSecret, nil
}

func (o *CreateOrUpdateOptions) dropYttOverlaySecrets(pkgi *kcpkgv1alpha1.PackageInstall, client kubernetes.Interface) error {
	o.statusUI.PrintMessage("Dropping overlay secrets")
	overlaySecretName, hasOverlay := pkgi.Annotations[yttOverlayAnnotation]
	if !hasOverlay {
		return nil
	}

	overlaySecret, err := client.CoreV1().Secrets(o.NamespaceFlags.Name).Get(context.Background(), overlaySecretName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			o.statusUI.PrintMessagef("Overlay secret '%s' not found in namespace '%s'", overlaySecretName, o.NamespaceFlags.Name)
			return nil
		}
		return fmt.Errorf("Getting overlay secret: %s", err.Error())
	}
	annValue, found := overlaySecret.Annotations[KctrlPkgAnnotation]
	if !found || annValue != o.createdAnnotations.PackageAnnValue() {
		return fmt.Errorf("Overlay secret was not created by kctrl")
	}

	err = client.CoreV1().Secrets(o.NamespaceFlags.Name).Delete(context.Background(), overlaySecretName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Deleting overlays secret: %s", err)
	}

	return nil
}

func (o *CreateOrUpdateOptions) removeYttOverlaysAnnotation(pkgi *kcpkgv1alpha1.PackageInstall) {
	delete(pkgi.Annotations, yttOverlayAnnotation)
}

func getPackageInstallDescription(name string, namespace string) string {
	description := fmt.Sprintf("packageinstall/%s (packaging.carvel.dev/v1alpha1)", name)
	if len(namespace) > 0 {
		description += " namespace: " + namespace
	} else {
		description += " cluster"
	}
	return description
}

func hasYttOverlays(pkgi *kcpkgv1alpha1.PackageInstall) bool {
	for annotation := range pkgi.Annotations {
		if strings.HasPrefix(annotation, yttOverlayPrefix) {
			return true
		}
	}
	return false
}
