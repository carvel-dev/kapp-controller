// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package datapackaging

import (
	"context"
	"time"

	internalpkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/internalpackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/validation"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/watchers"
	installclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

// PackageVersionCRDREST is a rest implementation that proxies the rest endpoints provided by
// CRDs. This will allow us to introduce the api server without the
// complexities associated with custom storage options for now.
type PackageVersionCRDREST struct {
	crdClient       installclient.Interface
	globalNamespace string
}

var (
	_ rest.StandardStorage    = &PackageVersionCRDREST{}
	_ rest.ShortNamesProvider = &PackageVersionCRDREST{}
)

const (
	packageVersionName         = "PackageVersion"
	internalPackageVersionName = "InternalPackageVersion"
)

func NewPackageVersionCRDREST(crdClient installclient.Interface, globalNS string) *PackageVersionCRDREST {
	return &PackageVersionCRDREST{crdClient, globalNS}
}

func (r *PackageVersionCRDREST) ShortNames() []string {
	return []string{"pkgv"}
}

func (r *PackageVersionCRDREST) NamespaceScoped() bool {
	return true
}

func (r *PackageVersionCRDREST) New() runtime.Object {
	return &datapackaging.PackageVersion{}
}

func (r *PackageVersionCRDREST) NewList() runtime.Object {
	return &datapackaging.PackageVersionList{}
}

func (r *PackageVersionCRDREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	if createValidation != nil {
		if err := createValidation(ctx, obj); err != nil {
			return nil, err
		}
	}

	pkgVersion := obj.(*datapackaging.PackageVersion)

	errs := validation.ValidatePackageVersion(*pkgVersion)
	if len(errs) != 0 {
		return nil, errors.NewInvalid(pkgVersion.GroupVersionKind().GroupKind(), pkgVersion.Name, errs)
	}

	ipv := r.packageVersionToInternalPackageVersion(pkgVersion)
	ipv, err := r.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).Create(ctx, ipv, *options)
	return r.internalPackageVersionToPackageVersion(ipv), err
}

func (r *PackageVersionCRDREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	pkgv, err := r.namespacedGet(ctx, namespace, name, options)
	if errors.IsNotFound(err) && namespace != r.globalNamespace && namespace != "" {
		pkgv, err = r.namespacedGet(ctx, r.globalNamespace, name, options)
	}
	return pkgv, err
}

func (r *PackageVersionCRDREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	// field selector isnt supported by CRD's so reset it, we will apply it later
	fs := options.FieldSelector
	options.FieldSelector = fields.Everything()

	namespace := request.NamespaceValue(ctx)

	namespacedPVList, err := r.namespacedList(ctx, namespace, options)
	if err != nil {
		return nil, err
	}

	itemList := namespacedPVList.Items
	var globalPVList *datapackaging.PackageVersionList
	if namespace != "" && namespace != r.globalNamespace {
		globalPVList, err = r.namespacedList(ctx, r.globalNamespace, options)
		if err != nil {
			return nil, err
		}

		itemList = append(itemList, globalPVList.Items...)
	}

	pvList := datapackaging.PackageVersionList{
		TypeMeta: namespacedPVList.TypeMeta,
		ListMeta: namespacedPVList.ListMeta,
	}

	// dedup so that namespace packages take precedence over global
	packageVersionIdentifiers := make(map[string]struct{})
	for _, pv := range itemList {
		// all namespaced pkgs come first in the list, so if we have seen one don't append
		identifier := pv.Spec.PackageName + "/" + pv.Spec.Version
		if _, seen := packageVersionIdentifiers[identifier]; !seen {
			pvList.Items = append(pvList.Items, pv)
			packageVersionIdentifiers[identifier] = struct{}{}
		}
	}

	filteredList := r.applySelector(pvList, fs)
	return &filteredList, err
}

func (r *PackageVersionCRDREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	namespace := request.NamespaceValue(ctx)
	pv, err := r.namespacedGet(ctx, namespace, name, &metav1.GetOptions{})

	if errors.IsNotFound(err) {
		updatedPkgVersion, err := objInfo.UpdatedObject(ctx, &datapackaging.PackageVersion{})
		if err != nil {
			return nil, false, err
		}

		if createValidation != nil {
			if err := createValidation(ctx, updatedPkgVersion); err != nil {
				return nil, false, err
			}
		}

		obj, err := r.Create(ctx, updatedPkgVersion, createValidation, &metav1.CreateOptions{TypeMeta: options.TypeMeta, DryRun: options.DryRun, FieldManager: options.FieldManager})
		if err != nil {
			return nil, true, err
		}

		return obj, true, nil
	}

	if err != nil {
		return nil, false, err
	}

	updatedObj, err := objInfo.UpdatedObject(ctx, pv)
	if err != nil {
		return nil, false, err
	}

	updatedPkgVersion := updatedObj.(*datapackaging.PackageVersion)
	errs := validation.ValidatePackageVersion(*updatedPkgVersion)
	if len(errs) != 0 {
		return nil, false, errors.NewInvalid(updatedPkgVersion.GroupVersionKind().GroupKind(), updatedPkgVersion.Name, errs)
	}

	updatedIpv := r.packageVersionToInternalPackageVersion(updatedPkgVersion)
	updatedIpv, err = r.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).Update(ctx, updatedIpv, *options)
	return r.internalPackageVersionToPackageVersion(updatedIpv), false, err
}

func (r *PackageVersionCRDREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	namespace := request.NamespaceValue(ctx)
	pv, err := r.namespacedGet(ctx, namespace, name, &metav1.GetOptions{})

	if errors.IsNotFound(err) {
		return nil, true, err
	}

	if err != nil {
		return nil, false, err
	}

	if deleteValidation != nil {
		if err := deleteValidation(ctx, pv); err != nil {
			return nil, true, err
		}
	}

	err = r.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).Delete(ctx, name, *options)
	if err != nil {
		return nil, false, err
	}

	return nil, true, nil
}

func (r *PackageVersionCRDREST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	list, err := r.namespacedList(ctx, namespace, listOptions)
	if err != nil {
		return nil, err
	}

	var deletedPackages []datapackaging.PackageVersion
	for _, pv := range list.Items {
		_, _, err := r.Delete(ctx, pv.Name, deleteValidation, options)
		if err != nil {
			break
		}
		// Although intPkgV => PkgV can return nil when intpkg is nil, intpkgv here will
		// never be nil, thanks to the type system, so we are ok to deref directly
		deletedPackages = append(deletedPackages, pv)
	}
	return &datapackaging.PackageVersionList{Items: deletedPackages}, err
}

func (r *PackageVersionCRDREST) Watch(ctx context.Context, options *internalversion.ListOptions) (watch.Interface, error) {
	namespace := request.NamespaceValue(ctx)

	watcher, err := r.namespacedWatch(ctx, namespace, options)
	if errors.IsNotFound(err) && namespace != r.globalNamespace {
		watcher, err = r.namespacedWatch(ctx, r.globalNamespace, options)
	}

	return watcher, err
}

func (r *PackageVersionCRDREST) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table
	fn := func(obj runtime.Object) error {
		pkgVersion := obj.(*datapackaging.PackageVersion)
		table.Rows = append(table.Rows, metav1.TableRow{
			Cells: []interface{}{
				pkgVersion.Name, pkgVersion.Spec.PackageName, pkgVersion.Namespace == r.globalNamespace,
				pkgVersion.Spec.Version, time.Since(pkgVersion.ObjectMeta.CreationTimestamp.Time).Round(1 * time.Second).String(),
			},
			Object: runtime.RawExtension{Object: obj},
		})
		return nil
	}
	switch {
	case meta.IsListType(obj):
		if err := meta.EachListItem(obj, fn); err != nil {
			return nil, err
		}
	default:
		if err := fn(obj); err != nil {
			return nil, err
		}
	}
	if m, err := meta.ListAccessor(obj); err == nil {
		table.ResourceVersion = m.GetResourceVersion()
		table.SelfLink = m.GetSelfLink()
		table.Continue = m.GetContinue()
		table.RemainingItemCount = m.GetRemainingItemCount()
	} else {
		if m, err := meta.CommonAccessor(obj); err == nil {
			table.ResourceVersion = m.GetResourceVersion()
			table.SelfLink = m.GetSelfLink()
		}
	}
	if opt, ok := tableOptions.(*metav1.TableOptions); !ok || !opt.NoHeaders {
		table.ColumnDefinitions = []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string", Format: "name", Description: "Package Version resource name"},
			{Name: "Package Name", Type: "string", Format: "name", Description: "Associated Package name"},
			{Name: "Global", Type: "boolean", Description: "If package version is global"},
			{Name: "Version", Type: "string", Description: "Version"},
			{Name: "Age", Type: "date", Description: "Time since resource creation"},
		}
	}
	return &table, nil
}

func (r *PackageVersionCRDREST) namespacedList(ctx context.Context, namespace string, options *internalversion.ListOptions) (*datapackaging.PackageVersionList, error) {
	list, err := r.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).List(ctx, r.internalToMetaListOpts(*options))
	return r.internalPackageVersionListToPackageVersionList(list), err
}

func (r *PackageVersionCRDREST) namespacedGet(ctx context.Context, namespace, name string, options *metav1.GetOptions) (*datapackaging.PackageVersion, error) {
	ipv, err := r.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).Get(ctx, name, *options)
	return r.internalPackageVersionToPackageVersion(ipv), err
}

func (r *PackageVersionCRDREST) namespacedWatch(ctx context.Context, namespace string, options *internalversion.ListOptions) (watch.Interface, error) {
	fs := options.FieldSelector
	options.FieldSelector = fields.Everything()

	internalWatcher, err := r.crdClient.InternalV1alpha1().InternalPackageVersions(namespace).Watch(ctx, r.internalToMetaListOpts(*options))
	if err != nil {
		return nil, err
	}

	return watchers.NewTranslationWatcher(r.translateFunc(), r.filterFunc(fs), internalWatcher), nil
}

func (r *PackageVersionCRDREST) internalToMetaListOpts(options internalversion.ListOptions) metav1.ListOptions {
	lo := metav1.ListOptions{
		TypeMeta:             options.TypeMeta,
		Watch:                options.Watch,
		AllowWatchBookmarks:  options.AllowWatchBookmarks,
		ResourceVersion:      options.ResourceVersion,
		ResourceVersionMatch: options.ResourceVersionMatch,
		TimeoutSeconds:       options.TimeoutSeconds,
		Limit:                options.Limit,
		Continue:             options.Continue,
	}

	if options.LabelSelector != nil {
		lo.LabelSelector = options.LabelSelector.String()
	}

	if options.FieldSelector != nil {
		lo.FieldSelector = options.FieldSelector.String()
	}
	return lo
}

func (r *PackageVersionCRDREST) internalPackageVersionToPackageVersion(ipv *internalpkgingv1alpha1.InternalPackageVersion) *datapackaging.PackageVersion {
	if ipv == nil {
		return nil
	}

	pv := (*datapackaging.PackageVersion)(ipv)
	for i := range pv.ManagedFields {
		if pv.ManagedFields[i].APIVersion == internalpkgingv1alpha1.SchemeGroupVersion.Identifier() {
			pv.ManagedFields[i].APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}
	pv.TypeMeta.Kind = packageVersionName
	pv.TypeMeta.APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
	return pv
}

func (r *PackageVersionCRDREST) packageVersionToInternalPackageVersion(pv *datapackaging.PackageVersion) *internalpkgingv1alpha1.InternalPackageVersion {
	if pv == nil {
		return nil
	}

	ipv := (*internalpkgingv1alpha1.InternalPackageVersion)(pv)
	for i := range ipv.ManagedFields {
		if ipv.ManagedFields[i].APIVersion == datapkgingv1alpha1.SchemeGroupVersion.Identifier() {
			ipv.ManagedFields[i].APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}
	ipv.TypeMeta.Kind = internalPackageVersionName
	ipv.TypeMeta.APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()
	return ipv
}

func (r *PackageVersionCRDREST) internalPackageVersionListToPackageVersionList(list *internalpkgingv1alpha1.InternalPackageVersionList) *datapackaging.PackageVersionList {
	if list == nil {
		return nil
	}

	pvList := datapackaging.PackageVersionList{
		TypeMeta: list.TypeMeta,
		ListMeta: list.ListMeta,
	}

	for _, item := range list.Items {
		pvList.Items = append(pvList.Items, *r.internalPackageVersionToPackageVersion(&item))
	}

	return &pvList
}

func (r *PackageVersionCRDREST) translateFunc() func(watch.Event) watch.Event {
	return func(evt watch.Event) watch.Event {
		if pv, ok := evt.Object.(*internalpkgingv1alpha1.InternalPackageVersion); ok {
			evt.Object = r.internalPackageVersionToPackageVersion(pv)
		}
		return evt
	}
}

func (r *PackageVersionCRDREST) filterFunc(fs fields.Selector) func(evt watch.Event) bool {
	dontFilter := func(evt watch.Event) bool {
		return true
	}

	filter := func(evt watch.Event) bool {
		if pv, ok := evt.Object.(*datapackaging.PackageVersion); ok {
			fieldSet := fields.Set{"spec.packageName": pv.Spec.PackageName}
			if fs.Matches(fieldSet) {
				return true
			}
			return false
		}
		return true
	}

	if fs == nil || fs.Empty() {
		return dontFilter
	}

	return filter
}

func (r *PackageVersionCRDREST) applySelector(list datapackaging.PackageVersionList, selector fields.Selector) datapackaging.PackageVersionList {
	if selector == nil || selector.Empty() {
		return list
	}

	filteredPVs := []datapackaging.PackageVersion{}
	for _, pv := range list.Items {
		fieldSet := fields.Set{"spec.packageName": pv.Spec.PackageName, "metadata.name": pv.Name, "metadata.namespace": pv.Namespace}
		if selector.Matches(fieldSet) {
			filteredPVs = append(filteredPVs, pv)
		}
	}

	list.Items = filteredPVs
	return list
}
