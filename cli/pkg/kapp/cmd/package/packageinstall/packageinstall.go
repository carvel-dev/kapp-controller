// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"context"
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type CreatedResourceKind string

const (
	// For use with packages created using kapp-ctrl-cli
	KappPkgAnnotation       = "packaging.carvel.dev/package"
	KappPkgAnnotationPrefix = "package"

	// For use with packages created with Tanzu CLI
	// KappPkgAnnotation       = "tkg.tanzu.vmware.com/tanzu-package"
	// KappPkgAnnotationPrefix = "tanzu-package"

	KindClusterRole        CreatedResourceKind = "ClusterRole"
	KindClusterRoleBinding CreatedResourceKind = "ClusterRoleBinding"
	KindSecret             CreatedResourceKind = "Secret"
	KindServiceAccount     CreatedResourceKind = "ServiceAccount"
	KindNamespace          CreatedResourceKind = "Namespace"

	ClusterRoleBindingName = "%s-%s-cluster-rolebinding"
	ClusterRoleName        = "%s-%s-cluster-role"
	SecretName             = "%s-%s-values"
	ServiceAccountName     = "%s-%s-sa"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "installed",
		Aliases: []string{"pkgi"},
		Short:   "PackageInstall",
	}
	return cmd
}

func (k CreatedResourceKind) Resource() string {
	switch k {
	case KindClusterRole:
		return "clusterroles"
	case KindClusterRoleBinding:
		return "clusterrolebindings"
	case KindSecret:
		return "secrets"
	case KindServiceAccount:
		return "serviceaccounts"
	}
	return ""
}

func (k CreatedResourceKind) Name(pkgiName string, pkgiNamespace string) string {
	switch k {
	case KindClusterRole:
		return fmt.Sprintf(ClusterRoleName, pkgiName, pkgiNamespace)
	case KindClusterRoleBinding:
		return fmt.Sprintf(ClusterRoleBindingName, pkgiName, pkgiNamespace)
	case KindSecret:
		return fmt.Sprintf(SecretName, pkgiName, pkgiNamespace)
	case KindServiceAccount:
		return fmt.Sprintf(ServiceAccountName, pkgiName, pkgiNamespace)
	}
	return ""
}

func (k CreatedResourceKind) AsString() string {
	return string(k)
}

// waitForResourceInstallation waits until the package get installed successfully or a failure happen
func waitForResourceInstallation(name, namespace string, pollInterval, pollTimeout time.Duration, ui ui.UI, client kcclient.Interface) error {
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

func addCreatedResourceAnnotations(meta *metav1.ObjectMeta, createdSvcAccount, createdSecret bool) {
	if meta.Annotations == nil {
		meta.Annotations = make(map[string]string)
	}
	if createdSvcAccount {
		meta.Annotations[KappPkgAnnotation+"-"+KindClusterRole.AsString()] = fmt.Sprintf(ClusterRoleName, meta.Name, meta.Namespace)
		meta.Annotations[KappPkgAnnotation+"-"+KindClusterRoleBinding.AsString()] = fmt.Sprintf(ClusterRoleBindingName, meta.Name, meta.Namespace)
		meta.Annotations[KappPkgAnnotation+"-"+KindServiceAccount.AsString()] = fmt.Sprintf(ServiceAccountName, meta.Name, meta.Namespace)
	}
	if createdSecret {
		meta.Annotations[KappPkgAnnotation+"-"+KindSecret.AsString()] = fmt.Sprintf(SecretName, meta.Name, meta.Namespace)
	}
}
