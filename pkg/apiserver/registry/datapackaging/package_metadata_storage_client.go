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
	packageMetadataName         = "PackageMetadata"
	internalPackageMetadataName = "InternalPackageMetadata"
)

type PackageMetadataTranslator struct {
	namespace string
}

func NewPackageMetadataTranslator(namespace string) PackageMetadataTranslator {
	return PackageMetadataTranslator{namespace}
}

func (t PackageMetadataTranslator) ToExternalObj(intObj *internalpkgingv1alpha1.InternalPackageMetadata) *datapackaging.PackageMetadata {
	if intObj == nil {
		return nil
	}

	obj := (*datapackaging.PackageMetadata)(intObj)
	for i := range obj.ManagedFields {
		if obj.ManagedFields[i].APIVersion == internalpkgingv1alpha1.SchemeGroupVersion.Identifier() {
			obj.ManagedFields[i].APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}

	obj.TypeMeta.Kind = packageMetadataName
	obj.TypeMeta.APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
	if t.namespace != "" {
		obj.Namespace = t.namespace
	}

	// Self link is deprecated and planned for removal, so we don't translate it

	return obj
}

func (t PackageMetadataTranslator) ToInternalObj(extObj *datapackaging.PackageMetadata) *internalpkgingv1alpha1.InternalPackageMetadata {
	if extObj == nil {
		return nil
	}

	intObj := (*internalpkgingv1alpha1.InternalPackageMetadata)(extObj)
	for i := range intObj.ManagedFields {
		if intObj.ManagedFields[i].APIVersion == datapkgingv1alpha1.SchemeGroupVersion.Identifier() {
			intObj.ManagedFields[i].APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}
	intObj.TypeMeta.Kind = internalPackageMetadataName
	intObj.TypeMeta.APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()

	return intObj
}

func (t PackageMetadataTranslator) ToExternalList(intObjList *internalpkgingv1alpha1.InternalPackageMetadataList) *datapackaging.PackageMetadataList {
	if intObjList == nil {
		return nil
	}

	externalObjList := datapackaging.PackageMetadataList{
		TypeMeta: intObjList.TypeMeta,
		ListMeta: intObjList.ListMeta,
	}

	for _, item := range intObjList.Items {
		externalObjList.Items = append(externalObjList.Items, *t.ToExternalObj(&item))
	}

	return &externalObjList
}

func (t PackageMetadataTranslator) ToExternalWatcher(intObjWatcher watch.Interface) watch.Interface {
	watchTranslation := func(evt watch.Event) watch.Event {
		if intpkg, ok := evt.Object.(*internalpkgingv1alpha1.InternalPackageMetadata); ok {
			evt.Object = t.ToExternalObj(intpkg)
		}
		return evt
	}

	watchFilter := func(evt watch.Event) bool {
		return true
	}

	return watchers.NewTranslationWatcher(watchTranslation, watchFilter, intObjWatcher)
}

func (t PackageMetadataTranslator) ToExternalError(err error) error {
	// TODO: implement
	return err
}

////////////////////////////////////////////
//==================Client================//
////////////////////////////////////////////

type PackageMetadataStorageClient struct {
	crdClient  internalclient.Interface
	translator PackageMetadataTranslator
}

func NewPackageMetadataStorageClient(crdClient internalclient.Interface, translator PackageMetadataTranslator) *PackageMetadataStorageClient {
	return &PackageMetadataStorageClient{crdClient, translator}
}

func (psc PackageMetadataStorageClient) Create(ctx context.Context, namespace string, obj *datapackaging.PackageMetadata, opts metav1.CreateOptions) (*datapackaging.PackageMetadata, error) {
	internalObj := psc.translator.ToInternalObj(obj)

	internalObj, err := psc.crdClient.InternalV1alpha1().InternalPackageMetadatas(namespace).Create(ctx, internalObj, opts)

	return psc.translator.ToExternalObj(internalObj), psc.translator.ToExternalError(err)
}

func (psc PackageMetadataStorageClient) Get(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*datapackaging.PackageMetadata, error) {
	intObj, err := psc.crdClient.InternalV1alpha1().InternalPackageMetadatas(namespace).Get(ctx, name, opts)
	return psc.translator.ToExternalObj(intObj), psc.translator.ToExternalError(err)
}

func (psc PackageMetadataStorageClient) List(ctx context.Context, namespace string, opts metav1.ListOptions) (*datapackaging.PackageMetadataList, error) {
	intObjList, err := psc.crdClient.InternalV1alpha1().InternalPackageMetadatas(namespace).List(ctx, opts)
	return psc.translator.ToExternalList(intObjList), psc.translator.ToExternalError(err)
}

func (psc PackageMetadataStorageClient) Update(ctx context.Context, namespace string, updatedObj *datapackaging.PackageMetadata, opts metav1.UpdateOptions) (*datapackaging.PackageMetadata, error) {
	intUpdatedObj := psc.translator.ToInternalObj(updatedObj)

	intUpdatedObj, err := psc.crdClient.InternalV1alpha1().InternalPackageMetadatas(namespace).Update(ctx, intUpdatedObj, opts)
	return psc.translator.ToExternalObj(intUpdatedObj), psc.translator.ToExternalError(err)
}

func (psc PackageMetadataStorageClient) Delete(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	err := psc.crdClient.InternalV1alpha1().InternalPackageMetadatas(namespace).Delete(ctx, name, opts)
	return psc.translator.ToExternalError(err)
}

func (psc PackageMetadataStorageClient) Watch(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	watcher, err := psc.crdClient.InternalV1alpha1().InternalPackageMetadatas(namespace).Watch(ctx, opts)
	return psc.translator.ToExternalWatcher(watcher), psc.translator.ToExternalError(err)
}
