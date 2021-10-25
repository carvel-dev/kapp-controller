// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"

	kappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

// UninstallPackage uninstalls the PackageInstall and its associated resources from the cluster
func (p *pkgClient) UninstallPackage(o *tkgpackagedatamodel.PackageOptions, progress *tkgpackagedatamodel.PackageProgress) {
	var (
		pkgInstall *kappipkg.PackageInstall
		err        error
	)

	defer func() {
		progressCleanup(err, progress)
	}()

	progress.ProgressMsg <- fmt.Sprintf("Getting package install for '%s'", o.PkgInstallName)
	pkgInstall, err = p.kappClient.GetPackageInstall(o.PkgInstallName, o.Namespace)
	if err != nil {
		if apierrors.IsNotFound(err) {
			if err := p.deletePreviouslyInstalledResources(o); err != nil {
				return
			}
			err = &tkgpackagedatamodel.PackagePluginNonCriticalError{Reason: tkgpackagedatamodel.ErrPackageNotInstalled}
		} else {
			err = errors.Wrap(err, fmt.Sprintf("\nfailed to find installed package '%s' in namespace '%s'", o.PkgInstallName, o.Namespace))
		}
		return
	}

	progress.ProgressMsg <- fmt.Sprintf("Deleting package install '%s' from namespace '%s'", o.PkgInstallName, o.Namespace)

	if err = p.deletePackageInstall(o); err != nil {
		return
	}

	if err = p.waitForResourceDeletion(o.PkgInstallName, o.Namespace, o.PollInterval, o.PollTimeout, progress.ProgressMsg, tkgpackagedatamodel.ResourceTypePackageInstall); err != nil {
		return
	}

	if err = p.deletePkgPluginCreatedResources(pkgInstall, progress.ProgressMsg); err != nil {
		return
	}
}

// deletePkgPluginCreatedResources deletes the associated resources which were installed upon installation of the PackageInstall CR
func (p *pkgClient) deletePkgPluginCreatedResources(pkgInstall *kappipkg.PackageInstall, progress chan string) error { //nolint:gocyclo
	for k, v := range pkgInstall.GetAnnotations() {
		split := strings.Split(k, "/")
		if len(split) <= 1 {
			continue
		}
		resourceKind := strings.Split(split[1], tkgpackagedatamodel.TanzuPkgPluginPrefix+"-")
		if len(resourceKind) <= 1 {
			continue
		}

		var obj runtime.Object
		objMeta := metav1.ObjectMeta{Name: v, Namespace: pkgInstall.Namespace}

		switch resourceKind[1] {
		case tkgpackagedatamodel.KindSecret:
			if progress != nil {
				progress <- fmt.Sprintf("Deleting secret '%s'", objMeta.Name)
			}
			obj = &corev1.Secret{ObjectMeta: objMeta, TypeMeta: metav1.TypeMeta{Kind: tkgpackagedatamodel.KindSecret}}
			if err := p.kappClient.GetClient().Delete(context.Background(), obj); err != nil {
				return errors.Wrap(err, "failed to delete Secret resource")
			}
		case tkgpackagedatamodel.KindServiceAccount:
			if progress != nil {
				progress <- fmt.Sprintf("Deleting service account '%s'", objMeta.Name)
			}
			obj = &corev1.ServiceAccount{ObjectMeta: objMeta, TypeMeta: metav1.TypeMeta{Kind: tkgpackagedatamodel.KindServiceAccount}}
			if err := p.kappClient.GetClient().Delete(context.Background(), obj); err != nil {
				return errors.Wrap(err, "failed to delete ServiceAccount resource")
			}
		case tkgpackagedatamodel.KindClusterRole:
			if progress != nil {
				progress <- fmt.Sprintf("Deleting admin role '%s'", objMeta.Name)
			}
			obj = &rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{Name: v},
				TypeMeta:   metav1.TypeMeta{Kind: tkgpackagedatamodel.KindClusterRole},
			}
			if err := p.kappClient.GetClient().Delete(context.Background(), obj); err != nil {
				return errors.Wrap(err, "failed to delete ClusterRole resource")
			}
		case tkgpackagedatamodel.KindClusterRoleBinding:
			if progress != nil {
				progress <- fmt.Sprintf("Deleting role binding '%s'", objMeta.Name)
			}
			obj = &rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{Name: v},
				TypeMeta:   metav1.TypeMeta{Kind: tkgpackagedatamodel.KindClusterRoleBinding},
			}
			if err := p.kappClient.GetClient().Delete(context.Background(), obj); err != nil {
				return errors.Wrap(err, "failed to delete ClusterRoleBinding resource")
			}
		}
	}

	return nil
}

// deletePackageInstall deletes the PackageInstall CR
func (p *pkgClient) deletePackageInstall(o *tkgpackagedatamodel.PackageOptions) error {
	obj := &kappipkg.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.PkgInstallName,
			Namespace: o.Namespace,
		},
		TypeMeta: metav1.TypeMeta{Kind: tkgpackagedatamodel.KindPackageInstall},
	}

	if err := p.kappClient.GetClient().Delete(context.Background(), obj); err != nil {
		if !apierrors.IsNotFound(err) {
			return errors.Wrap(err, "failed to delete PackageInstall resource")
		}
	}

	return nil
}

// deletePreviouslyInstalledResources deletes the related resources if previously installed through the package plugin
func (p *pkgClient) deletePreviouslyInstalledResources(o *tkgpackagedatamodel.PackageOptions) error {
	var objMeta metav1.ObjectMeta
	resourceAnnotation := fmt.Sprintf(tkgpackagedatamodel.TanzuPkgPluginResource, o.PkgInstallName, o.Namespace)

	objMeta = metav1.ObjectMeta{
		Name: fmt.Sprintf(tkgpackagedatamodel.ClusterRoleBindingName, o.PkgInstallName, o.Namespace),
	}
	if err := p.deleteAnnotatedResource(&rbacv1.ClusterRoleBinding{}, crtclient.ObjectKey{Name: objMeta.Name}, resourceAnnotation); err != nil {
		return err
	}

	objMeta = metav1.ObjectMeta{
		Name: fmt.Sprintf(tkgpackagedatamodel.ClusterRoleName, o.PkgInstallName, o.Namespace),
	}
	if err := p.deleteAnnotatedResource(&rbacv1.ClusterRole{}, crtclient.ObjectKey{Name: objMeta.Name}, resourceAnnotation); err != nil {
		return err
	}

	objMeta = metav1.ObjectMeta{
		Name:      fmt.Sprintf(tkgpackagedatamodel.ServiceAccountName, o.PkgInstallName, o.Namespace),
		Namespace: o.Namespace,
	}
	if err := p.deleteAnnotatedResource(&corev1.ServiceAccount{}, crtclient.ObjectKey{Name: objMeta.Name, Namespace: o.Namespace}, resourceAnnotation); err != nil {
		return err
	}

	objMeta = metav1.ObjectMeta{
		Name:      fmt.Sprintf(tkgpackagedatamodel.SecretName, o.PkgInstallName, o.Namespace),
		Namespace: o.Namespace,
	}
	if err := p.deleteAnnotatedResource(&corev1.Secret{}, crtclient.ObjectKey{Name: objMeta.Name, Namespace: o.Namespace}, resourceAnnotation); err != nil {
		return err
	}

	return nil
}

// deleteAnnotatedResource deletes the corresponding resource to the installed package name & namespace in case it has the package plugin annotation
func (p *pkgClient) deleteAnnotatedResource(obj runtime.Object, objKey crtclient.ObjectKey, resourceAnnotation string) error {
	if err := p.kappClient.GetClient().Get(context.Background(), objKey, obj); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
	} else {
		o, err := meta.Accessor(obj)
		if err != nil {
			return err
		}
		for k, v := range o.GetAnnotations() {
			split := strings.Split(k, "/")
			if len(split) <= 1 || split[1] != tkgpackagedatamodel.TanzuPkgPluginPrefix {
				continue
			}
			if v != resourceAnnotation {
				continue
			}
			if err := p.kappClient.GetClient().Delete(context.Background(), obj); err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func progressCleanup(err error, progress *tkgpackagedatamodel.PackageProgress) {
	if err != nil {
		progress.Err <- err
	}
	close(progress.ProgressMsg)
	close(progress.Done)
}

// waitForResourceDeletion waits until the CR gets deleted successfully or a failure happens
func (p *pkgClient) waitForResourceDeletion(name, namespace string, pollInterval, pollTimeout time.Duration, progress chan string, rscType tkgpackagedatamodel.ResourceType) error {
	var status kappctrl.GenericStatus
	if err := wait.Poll(pollInterval, pollTimeout, func() (done bool, err error) {
		switch rscType {
		case tkgpackagedatamodel.ResourceTypePackageRepository:
			resource, err := p.kappClient.GetPackageRepository(name, namespace)
			if err != nil {
				if apierrors.IsNotFound(err) {
					return true, nil
				}
				return false, err
			}
			if resource.Generation != resource.Status.ObservedGeneration {
				// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
				return false, nil
			}
			status = resource.Status.GenericStatus
		case tkgpackagedatamodel.ResourceTypePackageInstall:
			resource, err := p.kappClient.GetPackageInstall(name, namespace)
			if err != nil {
				if apierrors.IsNotFound(err) {
					return true, nil
				}
				return false, err
			}
			if resource.Generation != resource.Status.ObservedGeneration {
				// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
				return false, nil
			}
			status = resource.Status.GenericStatus
		}
		for _, cond := range status.Conditions {
			if progress != nil {
				progress <- fmt.Sprintf("'%s' resource deletion status: %s", rscType.String(), cond.Type)
			}
			if cond.Type == kappctrl.DeleteFailed && cond.Status == corev1.ConditionTrue {
				return false, fmt.Errorf("resource deletion failed: %s. %s", status.UsefulErrorMessage, status.FriendlyDescription)
			}
		}

		return false, nil
	}); err != nil {
		return err
	}

	return nil
}
