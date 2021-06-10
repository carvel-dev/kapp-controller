// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package datapackaging

import (
	"context"

	internalpkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/internalpackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/watchers"
	internalclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	packageName         = "Package"
	internalPackageName = "InternalPackage"
)

type PackageTranslator struct {
	namespace string
}

func NewPackageTranslator(namespace string) PackageTranslator {
	return PackageTranslator{namespace}
}

func (t PackageTranslator) ToExternalObj(intObj *internalpkgingv1alpha1.InternalPackage) *datapackaging.Package {
	if intObj == nil {
		return nil
	}

	obj := (*datapackaging.Package)(intObj)
	for i := range obj.ManagedFields {
		if obj.ManagedFields[i].APIVersion == internalpkgingv1alpha1.SchemeGroupVersion.Identifier() {
			obj.ManagedFields[i].APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}

	obj.TypeMeta.Kind = packageName
	obj.TypeMeta.APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
	if t.namespace != "" {
		obj.Namespace = t.namespace
	}

	// Self link is deprecated and planned for removal, so we don't translate it

	return obj
}

func (t PackageTranslator) ToInternalObj(extObj *datapackaging.Package) *internalpkgingv1alpha1.InternalPackage {
	if extObj == nil {
		return nil
	}

	intObj := (*internalpkgingv1alpha1.InternalPackage)(extObj)
	for i := range intObj.ManagedFields {
		if intObj.ManagedFields[i].APIVersion == datapkgingv1alpha1.SchemeGroupVersion.Identifier() {
			intObj.ManagedFields[i].APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}
	intObj.TypeMeta.Kind = internalPackageName
	intObj.TypeMeta.APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()

	return intObj
}

func (t PackageTranslator) ToExternalList(intObjList *internalpkgingv1alpha1.InternalPackageList) *datapackaging.PackageList {
	if intObjList == nil {
		return nil
	}

	externalObjList := datapackaging.PackageList{
		TypeMeta: intObjList.TypeMeta,
		ListMeta: intObjList.ListMeta,
	}

	for _, item := range intObjList.Items {
		externalObjList.Items = append(externalObjList.Items, *t.ToExternalObj(&item))
	}

	return &externalObjList
}

func (t PackageTranslator) ToExternalWatcher(intObjWatcher watch.Interface) watch.Interface {
	watchTranslation := func(evt watch.Event) watch.Event {
		if intpkg, ok := evt.Object.(*internalpkgingv1alpha1.InternalPackage); ok {
			evt.Object = t.ToExternalObj(intpkg)
		}
		return evt
	}

	watchFilter := func(evt watch.Event) bool {
		return true
	}

	return watchers.NewTranslationWatcher(watchTranslation, watchFilter, intObjWatcher)
}

func (t PackageTranslator) ToExternalError(err error) error {
	// TODO: implement
	return err
}

////////////////////////////////////////////
//==================Client================//
////////////////////////////////////////////

type PackageStorageClient struct {
	crdClient  internalclient.Interface
	translator PackageTranslator
}

func NewPackageStorageClient(crdClient internalclient.Interface, translator PackageTranslator) *PackageStorageClient {
	return &PackageStorageClient{crdClient, translator}
}

func (psc PackageStorageClient) Create(ctx context.Context, namespace string, obj *datapackaging.Package, opts metav1.CreateOptions) (*datapackaging.Package, error) {
	internalObj := psc.translator.ToInternalObj(obj)

	internalObj, err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).Create(ctx, internalObj, opts)

	return psc.translator.ToExternalObj(internalObj), psc.translator.ToExternalError(err)
}

func (psc PackageStorageClient) Get(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*datapackaging.Package, error) {
	intObj, err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).Get(ctx, name, opts)
	return psc.translator.ToExternalObj(intObj), psc.translator.ToExternalError(err)
}

func (psc PackageStorageClient) List(ctx context.Context, namespace string, opts metav1.ListOptions) (*datapackaging.PackageList, error) {
	intObjList, err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).List(ctx, opts)
	return psc.translator.ToExternalList(intObjList), psc.translator.ToExternalError(err)
}

func (psc PackageStorageClient) Update(ctx context.Context, namespace string, updatedObj *datapackaging.Package, opts metav1.UpdateOptions) (*datapackaging.Package, error) {
	intUpdatedObj := psc.translator.ToInternalObj(updatedObj)

	intUpdatedObj, err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).Update(ctx, intUpdatedObj, opts)
	return psc.translator.ToExternalObj(intUpdatedObj), psc.translator.ToExternalError(err)
}

func (psc PackageStorageClient) Delete(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).Delete(ctx, name, opts)
	return psc.translator.ToExternalError(err)
}

func (psc PackageStorageClient) Watch(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	watcher, err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).Watch(ctx, opts)
	return psc.translator.ToExternalWatcher(watcher), psc.translator.ToExternalError(err)
}
