package datapackaging

import (
	"context"

	internalpkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/internalpackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/watchers"
	internalclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	packageVersionName         = "PackageVersion"
	internalPackageVersionName = "InternalPackageVersion"
)

type PackageVersionTranslator struct {
	namespace string
}

func NewPackageVersionTranslator(namespace string) PackageVersionTranslator {
	return PackageVersionTranslator{namespace}
}

func (t PackageVersionTranslator) ToExternalObj(intObj *internalpkgingv1alpha1.InternalPackageVersion) *datapackaging.PackageVersion {
	if intObj == nil {
		return nil
	}

	obj := (*datapackaging.PackageVersion)(intObj)
	for i := range obj.ManagedFields {
		if obj.ManagedFields[i].APIVersion == internalpkgingv1alpha1.SchemeGroupVersion.Identifier() {
			obj.ManagedFields[i].APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}

	obj.TypeMeta.Kind = packageVersionName
	obj.TypeMeta.APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
	if t.namespace != "" {
		obj.Namespace = t.namespace
	}

	// Self link is deprecated and planned for removal, so we don't translate it

	return obj
}

func (t PackageVersionTranslator) ToInternalObj(extObj *datapackaging.PackageVersion) *internalpkgingv1alpha1.InternalPackageVersion {
	if extObj == nil {
		return nil
	}

	intObj := (*internalpkgingv1alpha1.InternalPackageVersion)(extObj)
	for i := range intObj.ManagedFields {
		if intObj.ManagedFields[i].APIVersion == datapkgingv1alpha1.SchemeGroupVersion.Identifier() {
			intObj.ManagedFields[i].APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}
	intObj.TypeMeta.Kind = internalPackageVersionName
	intObj.TypeMeta.APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()

	return intObj
}

func (t PackageVersionTranslator) ToExternalList(intObjList *internalpkgingv1alpha1.InternalPackageVersionList) *datapackaging.PackageVersionList {
	if intObjList == nil {
		return nil
	}

	externalObjList := datapackaging.PackageVersionList{
		TypeMeta: intObjList.TypeMeta,
		ListMeta: intObjList.ListMeta,
	}

	for _, item := range intObjList.Items {
		externalObjList.Items = append(externalObjList.Items, *t.ToExternalObj(&item))
	}

	return &externalObjList
}

func (t PackageVersionTranslator) ToExternalWatcher(intObjWatcher watch.Interface, fieldSelector fields.Selector) watch.Interface {
	watchTranslation := func(evt watch.Event) watch.Event {
		if intpkg, ok := evt.Object.(*internalpkgingv1alpha1.InternalPackageVersion); ok {
			evt.Object = t.ToExternalObj(intpkg)
		}
		return evt
	}

	watchFilter := func(evt watch.Event) bool {
		if fieldSelector == nil || fieldSelector.Empty() {
			return true
		}

		if pv, ok := evt.Object.(*datapackaging.PackageVersion); ok {
			fieldSet := fields.Set{"spec.packageName": pv.Spec.PackageName}
			if fieldSelector.Matches(fieldSet) {
				return true
			}
			return false
		}
		return true
	}

	return watchers.NewTranslationWatcher(watchTranslation, watchFilter, intObjWatcher)
}

func (t PackageVersionTranslator) ToExternalError(err error) error {
	// TODO: implement
	return err
}

////////////////////////////////////////////
//==================Client================//
////////////////////////////////////////////

type PackageVersionStorageClient struct {
	crdClient  internalclient.Interface
	translator PackageVersionTranslator
}

func NewPackageVersionStorageClient(crdClient internalclient.Interface, translator PackageVersionTranslator) *PackageVersionStorageClient {
	return &PackageVersionStorageClient{crdClient, translator}
}

func (psc PackageVersionStorageClient) Create(ctx context.Context, namespace string, obj *datapackaging.PackageVersion, opts metav1.CreateOptions) (*datapackaging.PackageVersion, error) {
	internalObj := psc.translator.ToInternalObj(obj)

	internalObj, err := psc.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).Create(ctx, internalObj, opts)

	return psc.translator.ToExternalObj(internalObj), psc.translator.ToExternalError(err)
}

func (psc PackageVersionStorageClient) Get(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*datapackaging.PackageVersion, error) {
	intObj, err := psc.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).Get(ctx, name, opts)
	return psc.translator.ToExternalObj(intObj), psc.translator.ToExternalError(err)
}

func (psc PackageVersionStorageClient) List(ctx context.Context, namespace string, opts metav1.ListOptions) (*datapackaging.PackageVersionList, error) {
	intObjList, err := psc.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).List(ctx, opts)
	return psc.translator.ToExternalList(intObjList), psc.translator.ToExternalError(err)
}

func (psc PackageVersionStorageClient) Update(ctx context.Context, namespace string, updatedObj *datapackaging.PackageVersion, opts metav1.UpdateOptions) (*datapackaging.PackageVersion, error) {
	intUpdatedObj := psc.translator.ToInternalObj(updatedObj)

	intUpdatedObj, err := psc.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).Update(ctx, intUpdatedObj, opts)
	return psc.translator.ToExternalObj(intUpdatedObj), psc.translator.ToExternalError(err)
}

func (psc PackageVersionStorageClient) Delete(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	err := psc.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).Delete(ctx, name, opts)
	return psc.translator.ToExternalError(err)
}

func (psc PackageVersionStorageClient) Watch(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	fieldSelector, err := fields.ParseSelector(opts.FieldSelector)
	if err != nil {
		return nil, err
	}
	// clear this because internal package versions dont support field selectors
	opts.FieldSelector = fields.Everything().String()

	watcher, err := psc.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).Watch(ctx, opts)
	return psc.translator.ToExternalWatcher(watcher, fieldSelector), psc.translator.ToExternalError(err)
}
