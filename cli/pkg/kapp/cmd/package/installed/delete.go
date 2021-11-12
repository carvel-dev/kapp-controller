// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
)

type DeleteOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	pollInterval time.Duration
	pollTimeout  time.Duration
	wait         bool
}

func NewDeleteOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *DeleteOptions {
	return &DeleteOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewDeleteCmd(o *DeleteOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"d"},
		Short:   "Uninstall installed package",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVar(&o.Name, "name", "", "Name of PackageInstall")

	cmd.Flags().DurationVar(&o.pollInterval, "poll-interval", 1*time.Second, "Time interval between consecutive polls while reconciling")
	cmd.Flags().DurationVar(&o.pollTimeout, "poll-timeout", 1*time.Minute, "Timeout for the reconciliation process")
	cmd.Flags().BoolVar(&o.wait, "wait", true, "Wait for reconcilation, default true")

	return cmd
}

func (o *DeleteOptions) Run() error {
	o.ui.PrintLinef("Delete package install '%s' from namespace '%s'", o.Name, o.NamespaceFlags.Name)

	err := o.ui.AskForConfirmation()
	if err != nil {
		return err
	}

	kcClient, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return nil
	}

	//TODO: Read warnings flag. Is it needed?
	dynamicClient, err := o.depsFactory.DynamicClient(cmdcore.DynamicClientOpts{})
	if err != nil {
		return nil
	}

	pkgi, err := kcClient.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Get(
		context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

		o.ui.PrintLinef("Could not find package install '%s' in namespace '%s'. Cleaning up created resources.", o.Name, o.NamespaceFlags.Name)

		return o.cleanUpIfInstallNotFound(dynamicClient)
	}

	o.ui.PrintLinef("Deleting package install '%s' from namespace '%s'", o.Name, o.NamespaceFlags.Name)

	err = kcClient.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Delete(
		context.Background(), o.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	// TODO dont think we can support no-wait for deletion since we have to clean up other resources
	if !o.wait {
		return nil
	}

	o.ui.PrintLinef("Waiting for deletion of package install '%s' from namespace '%s'", o.Name, o.NamespaceFlags.Name)

	err = o.waitForResourceDelete(kcClient)
	if err != nil {
		return err
	}

	return o.deleteInstallCreatedResources(pkgi, dynamicClient)
}

// deletePkgPluginCreatedResources deletes the associated resources which were installed upon installation of the PackageInstall CR
func (o *DeleteOptions) deleteInstallCreatedResources(pkgInstall *kcpkgv1alpha1.PackageInstall, dynamicClient dynamic.Interface) error {
	for k, resourceName := range pkgInstall.GetAnnotations() {
		split := strings.Split(k, "/")
		if len(split) <= 1 {
			continue
		}

		resourceKind := CreatedResourceKind(strings.TrimPrefix(split[1], KappPkgAnnotationPrefix+"-"))

		var apiGroup, version, namespace string
		if resourceKind == KindClusterRole || resourceKind == KindClusterRoleBinding {
			apiGroup = rbacv1.SchemeGroupVersion.Group
			version = rbacv1.SchemeGroupVersion.Version
		} else {
			apiGroup = corev1.SchemeGroupVersion.Group
			version = corev1.SchemeGroupVersion.Version
			namespace = o.NamespaceFlags.Name
		}

		o.ui.PrintLinef("Deleting '%s': %s", resourceKind, resourceName)

		err := o.deleteResourceUsingGVR(schema.GroupVersionResource{
			Group:    apiGroup,
			Version:  version,
			Resource: resourceKind.Resource(),
		}, resourceName, namespace, dynamicClient)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *DeleteOptions) cleanUpIfInstallNotFound(dynamicClient dynamic.Interface) error {
	err := o.deleteIfExistsAndOwned(
		schema.GroupVersionResource{
			Group:    corev1.SchemeGroupVersion.Group,
			Version:  corev1.SchemeGroupVersion.Version,
			Resource: KindServiceAccount.Resource(),
		}, KindServiceAccount.Name(o.Name, o.NamespaceFlags.Name), o.NamespaceFlags.Name, dynamicClient)
	if err != nil {
		return err
	}

	err = o.deleteIfExistsAndOwned(
		schema.GroupVersionResource{
			Group:    corev1.SchemeGroupVersion.Group,
			Version:  corev1.SchemeGroupVersion.Version,
			Resource: KindServiceAccount.Resource(),
		}, KindClusterRole.Name(o.Name, o.NamespaceFlags.Name), o.NamespaceFlags.Name, dynamicClient)
	if err != nil {
		return err
	}

	err = o.deleteIfExistsAndOwned(
		schema.GroupVersionResource{
			Group:    rbacv1.SchemeGroupVersion.Group,
			Version:  rbacv1.SchemeGroupVersion.Version,
			Resource: KindClusterRole.Resource(),
		}, KindClusterRole.Name(o.Name, o.NamespaceFlags.Name), "", dynamicClient)
	if err != nil {
		return err
	}

	err = o.deleteIfExistsAndOwned(
		schema.GroupVersionResource{
			Group:    rbacv1.SchemeGroupVersion.Group,
			Version:  rbacv1.SchemeGroupVersion.Version,
			Resource: KindClusterRoleBinding.Resource(),
		}, KindClusterRoleBinding.Name(o.Name, o.NamespaceFlags.Name), "", dynamicClient)
	if err != nil {
		return err
	}

	return nil
}

func (o *DeleteOptions) deleteIfExistsAndOwned(groupVersionResource schema.GroupVersionResource, name string, namespace string, dynamicClient dynamic.Interface) error {
	resource, err := dynamicClient.Resource(groupVersionResource).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			// Ignoring NotFound errors
			return err
		}
		return nil
	}

	annotations := resource.GetAnnotations()
	pkgiIdentifier := fmt.Sprintf("%s-%s", o.Name, o.NamespaceFlags.Name)

	val, found := annotations[KappPkgAnnotation]
	if !found || val != pkgiIdentifier {
		// Do not delete if the resource is not owned by the package, but no need to error out
		return nil
	}

	o.ui.PrintLinef("Deleting '%s': %s", groupVersionResource.Resource, name)
	err = dynamicClient.Resource(groupVersionResource).Namespace(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (o *DeleteOptions) deleteResourceUsingGVR(groupVersionResource schema.GroupVersionResource, name string, namespace string, dynamicClient dynamic.Interface) error {
	err := dynamicClient.Resource(groupVersionResource).Namespace(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (o *DeleteOptions) waitForResourceDelete(kcClient kcclient.Interface) error {
	err := wait.Poll(o.pollInterval, o.pollTimeout, func() (bool, error) {
		resource, err := kcClient.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Get(
			context.Background(), o.Name, metav1.GetOptions{},
		)
		if err != nil {
			if errors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		if resource.Generation != resource.Status.ObservedGeneration {
			// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
			return false, nil
		}
		status := resource.Status.GenericStatus
		for _, cond := range status.Conditions {
			o.ui.PrintLinef("'PackageInstall' resource deletion status: %s", cond.Type)
			if cond.Type == kcv1alpha1.DeleteFailed && cond.Status == corev1.ConditionTrue {
				return false, fmt.Errorf("resource deletion failed: %s. %s", status.UsefulErrorMessage, status.FriendlyDescription)
			}
		}
		return false, nil
	})

	if err != nil {
		return err
	}

	return nil
}
