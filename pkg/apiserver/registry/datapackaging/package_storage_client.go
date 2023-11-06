// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package datapackaging

import (
	"context"
	"encoding/base32"
	"fmt"
	"strings"

	internalpkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apis/internalpackaging/v1alpha1"
	"carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging"
	datapkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"carvel.dev/kapp-controller/pkg/apiserver/watchers"
	internalclient "carvel.dev/kapp-controller/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
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

func (t PackageTranslator) ToInternalName(name string) string {
	return strings.ToLower(base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(name)))

}

func (t PackageTranslator) ToExternalName(name string) (string, error) {
	decodedBytes, err := base32.HexEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(name))
	return string(decodedBytes), err
}

func (t PackageTranslator) ToExternalObj(intObj *internalpkgingv1alpha1.InternalPackage) (*datapackaging.Package, error) {
	if intObj == nil {
		return nil, nil
	}

	obj := (*datapackaging.Package)(intObj.DeepCopy())
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

	// base64 deccode the name to recover the original name
	var err error
	obj.Name, err = t.ToExternalName(intObj.Name)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("decoding internal obj name '%s': %v", intObj.Name, err))
	}

	// Self link is deprecated and planned for removal, so we don't translate it

	return obj, nil
}

func (t PackageTranslator) ToInternalObj(extObj *datapackaging.Package) *internalpkgingv1alpha1.InternalPackage {
	if extObj == nil {
		return nil
	}

	intObj := (*internalpkgingv1alpha1.InternalPackage)(extObj.DeepCopy())
	for i := range intObj.ManagedFields {
		if intObj.ManagedFields[i].APIVersion == datapkgingv1alpha1.SchemeGroupVersion.Identifier() {
			intObj.ManagedFields[i].APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}
	intObj.TypeMeta.Kind = internalPackageName
	intObj.TypeMeta.APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()

	// base64 encode the name since it will be packageName.version and version
	// could potentially include invalid characters
	intObj.Name = t.ToInternalName(extObj.Name)

	return intObj
}

func (t PackageTranslator) ToExternalList(intObjList *internalpkgingv1alpha1.InternalPackageList) (*datapackaging.PackageList, error) {
	if intObjList == nil {
		return nil, nil
	}

	externalObjList := datapackaging.PackageList{
		TypeMeta: intObjList.TypeMeta,
		ListMeta: intObjList.ListMeta,
	}

	for _, item := range intObjList.Items {
		extObj, err := t.ToExternalObj(&item)
		if err != nil {
			return nil, err
		}
		externalObjList.Items = append(externalObjList.Items, *extObj)
	}

	return &externalObjList, nil
}

func (t PackageTranslator) ToExternalWatcher(intObjWatcher watch.Interface, fieldSelector fields.Selector) watch.Interface {
	watchTranslation := func(evt watch.Event) watch.Event {
		if intpkg, ok := evt.Object.(*internalpkgingv1alpha1.InternalPackage); ok {
			var err error
			evt.Object, err = t.ToExternalObj(intpkg)
			if err != nil {
				var status metav1.Status
				if statusErr, ok := err.(*errors.StatusError); ok {
					status = statusErr.Status()
				} else {
					status = errors.NewInternalError(err).Status()
				}
				return watch.Event{Type: watch.Error, Object: &status}
			}
		}
		return evt
	}

	watchFilter := func(evt watch.Event) bool {
		if fieldSelector == nil || fieldSelector.Empty() {
			return true
		}

		if pv, ok := evt.Object.(*datapackaging.Package); ok {
			fieldSet := fields.Set{"spec.refName": pv.Spec.RefName}
			if fieldSelector.Matches(fieldSet) {
				return true
			}
			return false
		}
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
	intObj := psc.translator.ToInternalObj(obj)

	intObj, err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).Create(ctx, intObj, opts)
	if err != nil {
		return nil, psc.translator.ToExternalError(err)
	}

	return psc.translator.ToExternalObj(intObj)
}

func (psc PackageStorageClient) Get(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*datapackaging.Package, error) {
	name = psc.translator.ToInternalName(name)
	intObj, err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, psc.translator.ToExternalError(err)
	}

	return psc.translator.ToExternalObj(intObj)
}

func (psc PackageStorageClient) List(ctx context.Context, namespace string, opts metav1.ListOptions) (*datapackaging.PackageList, error) {
	intObjList, err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).List(ctx, opts)
	if err != nil {
		return nil, psc.translator.ToExternalError(err)
	}
	return psc.translator.ToExternalList(intObjList)
}

func (psc PackageStorageClient) Update(ctx context.Context, namespace string, updatedObj *datapackaging.Package, opts metav1.UpdateOptions) (*datapackaging.Package, error) {
	intUpdatedObj := psc.translator.ToInternalObj(updatedObj)

	intUpdatedObj, err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).Update(ctx, intUpdatedObj, opts)
	if err != nil {
		return nil, psc.translator.ToExternalError(err)
	}

	return psc.translator.ToExternalObj(intUpdatedObj)
}

func (psc PackageStorageClient) Delete(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	name = psc.translator.ToInternalName(name)
	err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).Delete(ctx, name, opts)
	return psc.translator.ToExternalError(err)
}

func (psc PackageStorageClient) Watch(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	fieldSelector, err := fields.ParseSelector(opts.FieldSelector)
	if err != nil {
		return nil, err
	}
	// clear this because internal package versions dont support field selectors
	opts.FieldSelector = fields.Everything().String()

	watcher, err := psc.crdClient.InternalV1alpha1().InternalPackages(namespace).Watch(ctx, opts)
	return psc.translator.ToExternalWatcher(watcher, fieldSelector), psc.translator.ToExternalError(err)
}
