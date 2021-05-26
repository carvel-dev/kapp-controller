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
	"k8s.io/apiserver/pkg/registry/rest"
)

// PackageVersionCRDREST is a rest implementation that proxies the rest endpoints provided by
// CRDs. This will allow us to introduce the api server without the
// complexities associated with custom storage options for now.
type PackageVersionCRDREST struct {
	crdClient installclient.Interface
}

var (
	_ rest.StandardStorage    = &PackageVersionCRDREST{}
	_ rest.ShortNamesProvider = &PackageVersionCRDREST{}
)

const (
	packageVersionName         = "PackageVersion"
	internalPackageVersionName = "InternalPackageVersion"
)

func NewPackageVersionCRDREST(crdClient installclient.Interface) *PackageVersionCRDREST {
	return &PackageVersionCRDREST{crdClient}
}

func (r *PackageVersionCRDREST) ShortNames() []string {
	return []string{"pkgv"}
}

func (r *PackageVersionCRDREST) New() runtime.Object {
	return &datapackaging.PackageVersion{}
}

func (r *PackageVersionCRDREST) NewList() runtime.Object {
	return &datapackaging.PackageVersionList{}
}

func (r *PackageVersionCRDREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
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
	ipv, err := r.crdClient.InternalV1alpha1().InternalPackageVersions("").Create(ctx, ipv, *options)
	return r.internalPackageVersionToPackageVersion(ipv), err
}

func (r *PackageVersionCRDREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	pv, err := r.Get(ctx, name, &metav1.GetOptions{})

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
	updatedIpv, err = r.crdClient.InternalV1alpha1().InternalPackageVersions("").Update(ctx, updatedIpv, *options)
	return r.internalPackageVersionToPackageVersion(updatedIpv), false, err
}

func (r *PackageVersionCRDREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	ipv, err := r.crdClient.InternalV1alpha1().InternalPackageVersions("").Get(ctx, name, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		return nil, true, err
	}

	if err != nil {
		return nil, false, err
	}

	if deleteValidation != nil {
		if err := deleteValidation(ctx, ipv); err != nil {
			return nil, true, err
		}
	}

	err = r.crdClient.InternalV1alpha1().InternalPackageVersions("").Delete(ctx, name, *options)
	if err != nil {
		return nil, false, err
	}

	return nil, true, nil
}

func (r *PackageVersionCRDREST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	objs, err := r.List(ctx, listOptions)
	if err != nil {
		return nil, err
	}

	var deletedPackages []datapackaging.PackageVersion
	for _, obj := range objs.(*datapackaging.PackageVersionList).Items {
		_, _, err := r.Delete(ctx, obj.Name, deleteValidation, options)
		if err != nil {
			break
		}
		deletedPackages = append(deletedPackages, obj)
	}
	return &datapackaging.PackageVersionList{Items: deletedPackages}, err
}

func (r *PackageVersionCRDREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	ipv, err := r.crdClient.InternalV1alpha1().InternalPackageVersions("").Get(ctx, name, *options)
	return r.internalPackageVersionToPackageVersion(ipv), err
}

func (r *PackageVersionCRDREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	// field selector isnt supported by CRD's so reset it, we will apply it later
	fs := options.FieldSelector
	options.FieldSelector = fields.Everything()

	// Label selectors and other options will be applied here
	list, err := r.crdClient.InternalV1alpha1().InternalPackageVersions("").List(ctx, r.internalToMetaListOpts(*options))
	pkgList := datapackaging.PackageVersionList{
		TypeMeta: list.TypeMeta,
		ListMeta: list.ListMeta,
	}
	for _, ipv := range list.Items {
		pkgList.Items = append(pkgList.Items, *r.internalPackageVersionToPackageVersion(&ipv))
	}

	filteredList := r.applySelector(pkgList, fs)
	return &filteredList, err
}

func (r *PackageVersionCRDREST) NamespaceScoped() bool {
	return false
}

func (r *PackageVersionCRDREST) Watch(ctx context.Context, options *internalversion.ListOptions) (watch.Interface, error) {
	fs := options.FieldSelector
	options.FieldSelector = fields.Everything()
	watcher, err := r.crdClient.InternalV1alpha1().InternalPackageVersions("").Watch(ctx, r.internalToMetaListOpts(*options))
	return watchers.NewTranslationWatcher(r.translateFunc(), r.filterFunc(fs), watcher), err
}

func (r *PackageVersionCRDREST) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table
	fn := func(obj runtime.Object) error {
		pkgVersion := obj.(*datapackaging.PackageVersion)
		table.Rows = append(table.Rows, metav1.TableRow{
			Cells:  []interface{}{pkgVersion.Name, pkgVersion.Spec.PackageName, pkgVersion.Spec.Version, time.Since(pkgVersion.ObjectMeta.CreationTimestamp.Time).Round(1 * time.Second).String()},
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
			{Name: "Version", Type: "string", Description: "Version"},
			{Name: "Age", Type: "date", Description: "Time since resource creation"},
		}
	}
	return &table, nil
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
	pv := (*datapackaging.PackageVersion)(ipv)
	for i := range pv.ManagedFields {
		mf := pv.ManagedFields[i]
		if mf.APIVersion == internalpkgingv1alpha1.SchemeGroupVersion.Identifier() {
			mf.APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}
	pv.TypeMeta.Kind = packageVersionName
	pv.TypeMeta.APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
	return pv
}

func (r *PackageVersionCRDREST) packageVersionToInternalPackageVersion(pv *datapackaging.PackageVersion) *internalpkgingv1alpha1.InternalPackageVersion {
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
		fieldSet := fields.Set{"spec.packageName": pv.Spec.PackageName}
		if selector.Matches(fieldSet) {
			filteredPVs = append(filteredPVs, pv)
		}
	}

	list.Items = filteredPVs
	return list
}
