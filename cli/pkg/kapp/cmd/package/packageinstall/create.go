// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/k14s/kapp/pkg/kapp/cmd/core"
	"github.com/k14s/kapp/pkg/kapp/logger"
	"github.com/spf13/cobra"
	kcpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	pollInterval time.Duration
	pollTimeout  time.Duration
	wait         bool

	pkgiName           string
	packageName        string
	version            string
	valuesFile         string
	serviceAccountName string
	createNewNamespace bool

	NamespaceFlags cmdcore.NamespaceFlags
}

func NewCreateOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *CreateOptions {
	return &CreateOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewCreateCmd(o *CreateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Install package",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVar(&o.pkgiName, "name", "", "Name of PackageInstall")
	cmd.Flags().StringVar(&o.packageName, "package-name", "", "Name of package to be installed")
	cmd.Flags().StringVar(&o.version, "version", "", "Version of package to be installed")
	cmd.Flags().StringVar(&o.serviceAccountName, "service-account-name", "", "Name of an existing service account used to install underlying package contents, optional")
	cmd.Flags().BoolVar(&o.createNewNamespace, "create-namespace", false, "Create namespace if the target namespace does not exist, optional")
	cmd.Flags().StringVar(&o.valuesFile, "values-file", "", "The path to the configuration values file, optional")

	cmd.Flags().DurationVar(&o.pollInterval, "poll-interval", 1*time.Second, "Time interval between consecutive polls while reconciling")
	cmd.Flags().DurationVar(&o.pollTimeout, "poll-timeout", 1*time.Minute, "Timeout for the reconciliation process")
	cmd.Flags().BoolVar(&o.wait, "wait", true, "Wait for reconcilation, default true")

	return cmd
}

func (o *CreateOptions) Run() error {

	var (
		isServiceAccountCreated bool
		isSecretCreated         bool
	)

	client, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	kcClient, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	// TODO: Fallback to update if exists

	if o.createNewNamespace {
		o.ui.PrintLinef("Creating namespace '%s'", o.NamespaceFlags.Name)
		if err = o.createNamespace(); err != nil {
			return err
		}
	} else if _, err = client.CoreV1().Namespaces().Get(context.Background(), o.NamespaceFlags.Name, metav1.GetOptions{}); err != nil {
		return err
	}

	if isServiceAccountCreated, isSecretCreated, err = o.createRelatedResources(); err != nil {
		return err
	}

	o.ui.PrintLinef("Creating package resource")
	if err = o.createPackageInstall(isServiceAccountCreated, isSecretCreated); err != nil {
		return err
	}

	if o.wait {
		if err = waitForResourceInstallation(o.pkgiName, o.NamespaceFlags.Name, o.pollInterval, o.pollTimeout, o.ui, kcClient); err != nil {
			return err
		}
	}

	return nil
}

// TODO: Handle created resource names better. Reduce duplication of logic used to get the names
// TODO: Handle initialisation of clients better
func (o *CreateOptions) createRelatedResources() (bool, bool, error) {
	var (
		isServiceAccountCreated bool
		isSecretCreated         bool
		err                     error
	)

	if o.serviceAccountName == "" {

		o.ui.PrintLinef("Creating service account '%s'", fmt.Sprintf(ServiceAccountName, o.pkgiName, o.NamespaceFlags.Name))
		if isServiceAccountCreated, err = o.createOrUpdateServiceAccount(); err != nil {
			return isServiceAccountCreated, isSecretCreated, err
		}

		o.ui.PrintLinef("Creating cluster admin role '%s'", fmt.Sprintf(ClusterRoleName, o.pkgiName, o.NamespaceFlags.Name))
		if err := o.createOrUpdateClusterAdminRole(); err != nil {
			return isServiceAccountCreated, isSecretCreated, err
		}

		o.ui.PrintLinef("Creating cluster role binding '%s'", fmt.Sprintf(ClusterRoleBindingName, o.pkgiName, o.NamespaceFlags.Name))
		if err := o.createOrUpdateClusterRoleBinding(); err != nil {
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
			if svcAccountAnnotation != fmt.Sprintf("%s-%s", o.pkgiName, o.NamespaceFlags.Name) {
				err = fmt.Errorf("provided service account '%s' is already used by another package in namespace '%s': %s", o.serviceAccountName, o.NamespaceFlags.Name, err.Error())
				return isServiceAccountCreated, isSecretCreated, err
			}
		}
	}

	if o.valuesFile != "" {
		o.ui.PrintLinef("Creating secret '%s'", fmt.Sprintf(SecretName, o.pkgiName, o.NamespaceFlags.Name))
		if isSecretCreated, err = o.createOrUpdateDataValuesSecret(); err != nil {
			return isServiceAccountCreated, isSecretCreated, err
		}
	}

	return isServiceAccountCreated, isSecretCreated, nil
}

func (o *CreateOptions) createOrUpdateClusterAdminRole() error {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(ClusterRoleName, o.pkgiName, o.NamespaceFlags.Name),
			Annotations: map[string]string{KappPkgAnnotation: fmt.Sprintf("%s-%s", o.pkgiName, o.NamespaceFlags.Name)},
		},
		Rules: []rbacv1.PolicyRule{
			{APIGroups: []string{"*"}, Verbs: []string{"*"}, Resources: []string{"*"}},
		},
	}

	client, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	_, err = client.RbacV1().ClusterRoles().Create(context.Background(), clusterRole, metav1.CreateOptions{})
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

func (o *CreateOptions) createOrUpdateClusterRoleBinding() error {
	svcAccount := o.serviceAccountName
	if svcAccount == "" {
		svcAccount = fmt.Sprintf(ServiceAccountName, o.pkgiName, o.NamespaceFlags.Name)
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(ClusterRoleBindingName, o.pkgiName, o.NamespaceFlags.Name),
			Annotations: map[string]string{KappPkgAnnotation: fmt.Sprintf("%s-%s", o.pkgiName, o.NamespaceFlags.Name)},
		},
		Subjects: []rbacv1.Subject{{Kind: KindServiceAccount.AsString(), Name: svcAccount, Namespace: o.NamespaceFlags.Name}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Kind:     KindClusterRole.AsString(),
			Name:     fmt.Sprintf(ClusterRoleName, o.pkgiName, o.NamespaceFlags.Name),
		},
	}

	client, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	_, err = client.RbacV1().ClusterRoleBindings().Create(context.Background(), clusterRoleBinding, metav1.CreateOptions{})
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

func (o *CreateOptions) createOrUpdateDataValuesSecret() (bool, error) {
	var err error

	dataValues := make(map[string][]byte)

	dataValues[filepath.Base(o.valuesFile)], err = ioutil.ReadFile(o.valuesFile)
	if err != nil {
		return false, fmt.Errorf("failed to read from data values file '%s': %s", o.valuesFile, err.Error())
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(SecretName, o.pkgiName, o.NamespaceFlags.Name),
			Namespace:   o.NamespaceFlags.Name,
			Annotations: map[string]string{KappPkgAnnotation: fmt.Sprintf("%s-%s", o.pkgiName, o.NamespaceFlags.Name)},
		},
		Data: dataValues,
	}

	client, err := o.depsFactory.CoreClient()
	if err != nil {
		return false, err
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

func (o *CreateOptions) createNamespace() error {

	ns := &corev1.Namespace{
		TypeMeta:   metav1.TypeMeta{Kind: KindNamespace.AsString()},
		ObjectMeta: metav1.ObjectMeta{Name: o.NamespaceFlags.Name},
	}

	client, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	_, err = client.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func (o *CreateOptions) createPackageInstall(serviceAccountCreated, secretCreated bool) error {
	svcAccount := o.serviceAccountName
	if svcAccount == "" {
		svcAccount = fmt.Sprintf(ServiceAccountName, o.pkgiName, o.NamespaceFlags.Name)
	}

	// construct the PackageInstall CR
	packageInstall := &kcpkgv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{Name: o.pkgiName, Namespace: o.NamespaceFlags.Name},
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
					Name: fmt.Sprintf(SecretName, o.pkgiName, o.NamespaceFlags.Name),
				},
			},
		}
	}

	addCreatedResourceAnnotations(&packageInstall.ObjectMeta, serviceAccountCreated, secretCreated)

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	_, err = client.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Create(context.Background(), packageInstall, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create PackageInstall resource: %s", err.Error())
	}

	return nil
}

func (o *CreateOptions) createOrUpdateServiceAccount() (bool, error) {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(ServiceAccountName, o.pkgiName, o.NamespaceFlags.Name),
			Namespace:   o.NamespaceFlags.Name,
			Annotations: map[string]string{KappPkgAnnotation: fmt.Sprintf("%s-%s", o.pkgiName, o.NamespaceFlags.Name)}},
	}

	client, err := o.depsFactory.CoreClient()
	if err != nil {
		return false, err
	}

	_, err = client.CoreV1().ServiceAccounts(o.NamespaceFlags.Name).Create(context.Background(), serviceAccount, metav1.CreateOptions{})
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
