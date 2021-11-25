// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

const (
	valuesFileKey = "values"
)

type CreateOrUpdateOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	pollInterval time.Duration
	pollTimeout  time.Duration
	wait         bool

	packageName        string
	version            string
	valuesFile         string
	serviceAccountName string
	createNewNamespace bool

	install bool

	Name               string
	NamespaceFlags     cmdcore.NamespaceFlags
	CreatedAnnotations *CreatedResourceAnnotations

	positionalNameArg bool
}

func NewCreateOrUpdateOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, positionalNameArg bool) *CreateOrUpdateOptions {
	return &CreateOrUpdateOptions{ui: ui, depsFactory: depsFactory, logger: logger, positionalNameArg: positionalNameArg}
}

func NewCreateCmd(o *CreateOrUpdateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Install package",
		RunE:  func(_ *cobra.Command, args []string) error { return o.RunCreate(args) },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)

	if !o.positionalNameArg {
		cmd.Flags().StringVarP(&o.Name, "package-install", "i", "", "Set installed package name")
	}

	cmd.Flags().StringVar(&o.packageName, "package-name", "", "Set package name")
	cmd.Flags().StringVar(&o.version, "version", "", "Set package version")
	cmd.Flags().StringVar(&o.serviceAccountName, "service-account-name", "", "Name of an existing service account used to install underlying package contents, optional")
	cmd.Flags().BoolVar(&o.createNewNamespace, "create-namespace", false, "Create namespace if the target namespace does not exist, optional")
	cmd.Flags().StringVar(&o.valuesFile, "values-file", "", "The path to the configuration values file, optional")

	cmd.Flags().DurationVar(&o.pollInterval, "poll-interval", 1*time.Second, "Time interval between consecutive polls while reconciling")
	cmd.Flags().DurationVar(&o.pollTimeout, "poll-timeout", 1*time.Minute, "Timeout for the reconciliation process")
	cmd.Flags().BoolVar(&o.wait, "wait", true, "Wait for reconcilation, default true")

	return cmd
}

func NewInstallCmd(o *CreateOrUpdateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install package",
		RunE:  func(_ *cobra.Command, args []string) error { return o.RunCreate(args) },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)

	if !o.positionalNameArg {
		cmd.Flags().StringVarP(&o.Name, "package-install", "i", "", "Set installed package name")
	}

	cmd.Flags().StringVar(&o.packageName, "package-name", "", "Set package name")
	cmd.Flags().StringVar(&o.version, "version", "", "Set package version")
	cmd.Flags().StringVar(&o.serviceAccountName, "service-account-name", "", "Name of an existing service account used to install underlying package contents, optional")
	cmd.Flags().BoolVar(&o.createNewNamespace, "create-namespace", false, "Create namespace if the target namespace does not exist, optional")
	cmd.Flags().StringVar(&o.valuesFile, "values-file", "", "The path to the configuration values file, optional")

	cmd.Flags().DurationVar(&o.pollInterval, "poll-interval", 1*time.Second, "Time interval between consecutive polls while reconciling")
	cmd.Flags().DurationVar(&o.pollTimeout, "poll-timeout", 1*time.Minute, "Timeout for the reconciliation process")
	cmd.Flags().BoolVar(&o.wait, "wait", true, "Wait for reconcilation, default true")

	return cmd
}

func NewUpdateCmd(o *CreateOrUpdateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update package",
		RunE:  func(_ *cobra.Command, args []string) error { return o.RunUpdate(args) },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)

	if !o.positionalNameArg {
		cmd.Flags().StringVarP(&o.Name, "package-install", "i", "", "Set installed package name")
	}

	cmd.Flags().StringVar(&o.packageName, "package-name", "", "Name of package install to be updated")
	cmd.Flags().StringVar(&o.version, "version", "", "Set package version")
	cmd.Flags().StringVar(&o.serviceAccountName, "service-account-name", "", "Name of an existing service account used to install underlying package contents, optional")
	cmd.Flags().BoolVar(&o.createNewNamespace, "create-namespace", false, "Create namespace if the target namespace does not exist, optional")
	cmd.Flags().StringVar(&o.valuesFile, "values-file", "", "The path to the configuration values file, optional")
	cmd.Flags().BoolVarP(&o.install, "install", "", false, "Install package if the installed package does not exist, optional")

	cmd.Flags().DurationVar(&o.pollInterval, "poll-interval", 1*time.Second, "Time interval between consecutive polls while reconciling")
	cmd.Flags().DurationVar(&o.pollTimeout, "poll-timeout", 1*time.Minute, "Timeout for the reconciliation process")
	cmd.Flags().BoolVar(&o.wait, "wait", true, "Wait for reconcilation, default true")

	return cmd
}

func (o *CreateOrUpdateOptions) RunCreate(args []string) error {
	if o.positionalNameArg {
		o.Name = args[0]
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

	o.CreatedAnnotations = NewCreatedResourceAnnotations(o.Name, o.NamespaceFlags.Name)

	// Fallback to update if resource exists
	if pkgInstall != nil && err == nil {
		o.ui.PrintLinef("Updating existing package install")
		err = o.update(client, kcClient, pkgInstall)
		if err != nil {
			return err
		}
		return nil
	}

	if o.createNewNamespace {
		o.ui.PrintLinef("Creating namespace '%s'", o.NamespaceFlags.Name)
		if err = o.createNamespace(client); err != nil {
			return err
		}
	} else if _, err = client.CoreV1().Namespaces().Get(context.Background(), o.NamespaceFlags.Name, metav1.GetOptions{}); err != nil {
		return err
	}

	err = o.create(client, kcClient)
	if err != nil {
		return err
	}

	return nil
}

func (o *CreateOrUpdateOptions) create(client kubernetes.Interface, kcClient versioned.Interface) error {
	isServiceAccountCreated, isSecretCreated, err := o.createRelatedResources(client)
	if err != nil {
		return err
	}

	o.ui.PrintLinef("Creating package install resource")
	if err = o.createPackageInstall(isServiceAccountCreated, isSecretCreated, kcClient); err != nil {
		return err
	}

	if o.wait {
		if err = o.waitForResourceInstallation(o.Name, o.NamespaceFlags.Name, o.pollInterval, o.pollTimeout, o.ui, kcClient); err != nil {
			return err
		}
	}

	return nil
}

func (o *CreateOrUpdateOptions) RunUpdate(args []string) error {
	if o.positionalNameArg {
		o.Name = args[0]
	}

	client, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	kcClient, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	o.CreatedAnnotations = NewCreatedResourceAnnotations(o.Name, o.NamespaceFlags.Name)

	o.ui.PrintLinef("Getting package install for '%s'", o.Name)
	pkgInstall, err := kcClient.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Get(
		context.Background(), o.Name, metav1.GetOptions{},
	)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		if !o.install {
			return fmt.Errorf("Package not installed")
		}
		o.ui.PrintLinef("Installing package '%s'", o.Name)

		err = o.create(client, kcClient)
		if err != nil {
			return err
		}

		return nil
	}

	err = o.update(client, kcClient, pkgInstall)
	if err != nil {
		return err
	}

	return nil
}

func (o CreateOrUpdateOptions) update(client kubernetes.Interface, kcClient versioned.Interface, pkgInstall *kcpkgv1alpha1.PackageInstall) error {
	updatedPkgInstall, changed, err := o.preparePackageInstallForUpdate(pkgInstall)
	if err != nil {
		return err
	}

	if o.valuesFile == "" && !changed {
		return err
	}

	isSecretCreated, err := o.createOrUpdateValuesSecret(updatedPkgInstall, client)
	if err != nil {
		return err
	}

	o.ui.PrintLinef("Updating package install for '%s'", o.Name)
	o.addCreatedResourceAnnotations(&pkgInstall.ObjectMeta, false, isSecretCreated)
	_, err = kcClient.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Update(
		context.Background(), updatedPkgInstall, metav1.UpdateOptions{},
	)
	if err != nil {
		err = fmt.Errorf("failed to update package '%s': %s", o.Name, err.Error())
		return err
	}

	if o.wait {
		if err = o.waitForResourceInstallation(o.Name, o.NamespaceFlags.Name, o.pollInterval, o.pollTimeout, o.ui, kcClient); err != nil {
			return err
		}
	}

	return nil
}

func (o *CreateOrUpdateOptions) createRelatedResources(client kubernetes.Interface) (bool, bool, error) {
	var (
		isServiceAccountCreated bool
		isSecretCreated         bool
		err                     error
	)

	if o.serviceAccountName == "" {

		o.ui.PrintLinef("Creating service account '%s'", o.CreatedAnnotations.ServiceAccountAnnValue())
		if isServiceAccountCreated, err = o.createOrUpdateServiceAccount(client); err != nil {
			return isServiceAccountCreated, isSecretCreated, err
		}

		o.ui.PrintLinef("Creating cluster admin role '%s'", o.CreatedAnnotations.ClusterRoleAnnValue())
		if err := o.createOrUpdateClusterAdminRole(client); err != nil {
			return isServiceAccountCreated, isSecretCreated, err
		}

		o.ui.PrintLinef("Creating cluster role binding '%s'", o.CreatedAnnotations.ClusterRoleBindingAnnValue())
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
			err = fmt.Errorf("failed to find service account '%s' in namespace '%s': %s", o.serviceAccountName, o.NamespaceFlags.Name, err.Error())
			return isServiceAccountCreated, isSecretCreated, err
		}
		if svcAccountAnnotation, ok := svcAccount.GetAnnotations()[KappPkgAnnotation]; ok {
			if svcAccountAnnotation != o.CreatedAnnotations.PackageAnnValue() {
				err = fmt.Errorf("provided service account '%s' is already used by another package in namespace '%s': %s", o.serviceAccountName, o.NamespaceFlags.Name, err.Error())
				return isServiceAccountCreated, isSecretCreated, err
			}
		}
	}

	if o.valuesFile != "" {
		o.ui.PrintLinef("Creating secret '%s'", o.CreatedAnnotations.SecretAnnValue())
		if isSecretCreated, err = o.createOrUpdateDataValuesSecret(client); err != nil {
			return isServiceAccountCreated, isSecretCreated, err
		}
	}

	return isServiceAccountCreated, isSecretCreated, nil
}

func (o *CreateOrUpdateOptions) createOrUpdateClusterAdminRole(client kubernetes.Interface) error {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:        o.CreatedAnnotations.ClusterRoleAnnValue(),
			Annotations: map[string]string{KappPkgAnnotation: o.CreatedAnnotations.PackageAnnValue()},
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
		svcAccount = o.CreatedAnnotations.ServiceAccountAnnValue()
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        o.CreatedAnnotations.ClusterRoleBindingAnnValue(),
			Annotations: map[string]string{KappPkgAnnotation: o.CreatedAnnotations.PackageAnnValue()},
		},
		Subjects: []rbacv1.Subject{{Kind: KindServiceAccount.AsString(), Name: svcAccount, Namespace: o.NamespaceFlags.Name}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Kind:     KindClusterRole.AsString(),
			Name:     o.CreatedAnnotations.ClusterRoleAnnValue(),
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

	dataValues[valuesFileKey], err = ioutil.ReadFile(o.valuesFile)
	if err != nil {
		return false, fmt.Errorf("failed to read from data values file '%s': %s", o.valuesFile, err.Error())
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        o.CreatedAnnotations.SecretAnnValue(),
			Namespace:   o.NamespaceFlags.Name,
			Annotations: map[string]string{KappPkgAnnotation: o.CreatedAnnotations.PackageAnnValue()},
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

func (o *CreateOrUpdateOptions) createNamespace(client kubernetes.Interface) error {

	ns := &corev1.Namespace{
		TypeMeta:   metav1.TypeMeta{Kind: KindNamespace.AsString()},
		ObjectMeta: metav1.ObjectMeta{Name: o.NamespaceFlags.Name},
	}

	_, err := client.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func (o *CreateOrUpdateOptions) createPackageInstall(serviceAccountCreated, secretCreated bool, kcClient kcclient.Interface) error {
	svcAccount := o.serviceAccountName
	if svcAccount == "" {
		svcAccount = o.CreatedAnnotations.ServiceAccountAnnValue()
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
					Name: o.CreatedAnnotations.SecretAnnValue(),
				},
			},
		}
	}

	o.addCreatedResourceAnnotations(&packageInstall.ObjectMeta, serviceAccountCreated, secretCreated)

	_, err := kcClient.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Create(context.Background(), packageInstall, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create PackageInstall resource: %s", err.Error())
	}

	return nil
}

func (o *CreateOrUpdateOptions) createOrUpdateServiceAccount(client kubernetes.Interface) (bool, error) {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        o.CreatedAnnotations.ServiceAccountAnnValue(),
			Namespace:   o.NamespaceFlags.Name,
			Annotations: map[string]string{KappPkgAnnotation: o.CreatedAnnotations.PackageAnnValue()},
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
		err = fmt.Errorf("failed to update package '%s' as no existing package reference/version was found in the package install", o.Name)
		return nil, false, err
	}

	// If o.PackageName is provided by the user (via --package-name flag), verify that the package name in PackageInstall matches it.
	// This will prevent the users from accidentally overwriting an installed package with another package content due to choosing a pre-existing name for the package isntall.
	// Otherwise if o.PackageName is not provided, fill it from the installed package spec
	if o.packageName != "" && updatedPkgInstall.Spec.PackageRef.RefName != o.packageName {
		err = fmt.Errorf("installed package '%s' is already associated with package '%s'", o.Name, updatedPkgInstall.Spec.PackageRef.RefName)
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

	secretName := o.CreatedAnnotations.SecretAnnValue()

	if secretName == pkgInstallToUpdate.GetAnnotations()[KappPkgAnnotation+"-"+KindSecret.AsString()] {
		o.ui.PrintLinef("Updating secret '%s'", secretName)
		if err = o.updateDataValuesSecret(client); err != nil {
			err = fmt.Errorf("failed to update secret based on values file: %s", err.Error())
			return false, err
		}
	} else {
		o.ui.PrintLinef("Creating secret '%s'", secretName)
		if secretCreated, err = o.createOrUpdateDataValuesSecret(client); err != nil {
			return secretCreated, fmt.Errorf("failed to create secret based on values file: %s", err.Error())
		}
	}

	pkgInstallToUpdate.Spec.Values = []kcpkgv1alpha1.PackageInstallValues{
		{SecretRef: &kcpkgv1alpha1.PackageInstallValuesSecretRef{Name: secretName}},
	}

	return secretCreated, nil
}

func (o *CreateOrUpdateOptions) updateDataValuesSecret(client kubernetes.Interface) error {
	var err error
	dataValues := make(map[string][]byte)
	secretName := o.CreatedAnnotations.SecretAnnValue()

	if dataValues[valuesFileKey], err = ioutil.ReadFile(o.valuesFile); err != nil {
		return fmt.Errorf("failed to read from data values file '%s': %s", o.valuesFile, err.Error())
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: o.NamespaceFlags.Name}, Data: dataValues,
	}

	_, err = client.CoreV1().Secrets(o.NamespaceFlags.Name).Update(context.Background(), secret, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update Secret resource: %s", err.Error())
	}

	return nil
}

func (o *CreateOrUpdateOptions) addCreatedResourceAnnotations(meta *metav1.ObjectMeta, createdSvcAccount, createdSecret bool) {
	if meta.Annotations == nil {
		meta.Annotations = make(map[string]string)
	}
	if createdSvcAccount {
		meta.Annotations[KappPkgAnnotation+"-"+KindClusterRole.AsString()] = o.CreatedAnnotations.ClusterRoleAnnValue()
		meta.Annotations[KappPkgAnnotation+"-"+KindClusterRoleBinding.AsString()] = o.CreatedAnnotations.ClusterRoleBindingAnnValue()
		meta.Annotations[KappPkgAnnotation+"-"+KindServiceAccount.AsString()] = o.CreatedAnnotations.ServiceAccountAnnValue()
	}
	if createdSecret {
		meta.Annotations[KappPkgAnnotation+"-"+KindSecret.AsString()] = o.CreatedAnnotations.SecretAnnValue()
	}
}

// waitForResourceInstallation waits until the package get installed successfully or a failure happen
func (o *CreateOrUpdateOptions) waitForResourceInstallation(name, namespace string, pollInterval, pollTimeout time.Duration, ui ui.UI, client kcclient.Interface) error {
	var (
		status             kcv1alpha1.GenericStatus
		reconcileSucceeded bool
	)
	ui.PrintLinef("Waiting for PackageInstall reconciliation for '%s'", name)
	err := wait.Poll(pollInterval, pollTimeout, func() (done bool, err error) {

		resource, err := client.PackagingV1alpha1().PackageInstalls(namespace).Get(context.Background(), name, metav1.GetOptions{})
		//resource, err := p.kappClient.GetPackageInstall(name, namespace)
		if err != nil {
			return false, err
		}
		if resource.Generation != resource.Status.ObservedGeneration {
			// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
			return false, nil
		}
		status = resource.Status.GenericStatus

		for _, condition := range status.Conditions {
			ui.PrintLinef("PackageInstall resource install status: %s", condition.Type)

			switch {
			case condition.Type == kcv1alpha1.ReconcileSucceeded && condition.Status == corev1.ConditionTrue:
				ui.PrintLinef("PackageInstall resource successfully reconciled")
				reconcileSucceeded = true
				return true, nil
			case condition.Type == kcv1alpha1.ReconcileFailed && condition.Status == corev1.ConditionTrue:
				return false, fmt.Errorf("resource reconciliation failed: %s. %s", status.UsefulErrorMessage, status.FriendlyDescription)
			}
		}
		return false, nil
	})
	if err != nil {
		return err
	}

	if !reconcileSucceeded {
		return fmt.Errorf("PackageInstall resource reconciliation failed")
	}

	return nil
}
