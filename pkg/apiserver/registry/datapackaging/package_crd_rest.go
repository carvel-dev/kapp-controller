// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package datapackaging

import (
	"context"
	"strings"
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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

// PackageCRDREST is a rest implementation that proxies the rest endpoints provided by
// CRDs. This will allow us to introduce the api server without the
// complexities associated with custom storage options for now.
type PackageCRDREST struct {
	crdClient       installclient.Interface
	globalNamespace string
}

var (
	_ rest.StandardStorage    = &PackageCRDREST{}
	_ rest.ShortNamesProvider = &PackageCRDREST{}
)

const (
	packageName         = "Package"
	internalPackageName = "InternalPackage"
)

func NewPackageCRDREST(crdClient installclient.Interface, globalNS string) *PackageCRDREST {
	return &PackageCRDREST{crdClient, globalNS}
}

func (r *PackageCRDREST) ShortNames() []string {
	return []string{"pkg"}
}

func (r *PackageCRDREST) New() runtime.Object {
	return &datapackaging.Package{}
}

func (r *PackageCRDREST) NewList() runtime.Object {
	return &datapackaging.PackageList{}
}

func (r *PackageCRDREST) NamespaceScoped() bool {
	return true
}

func (r *PackageCRDREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	if createValidation != nil {
		if err := createValidation(ctx, obj); err != nil {
			return nil, err
		}
	}

	pkg := obj.(*datapackaging.Package)
	errs := validation.ValidatePackage(*pkg)
	if len(errs) != 0 {
		return nil, errors.NewInvalid(pkg.GroupVersionKind().GroupKind(), pkg.Name, errs)
	}

	intpkg := r.packageToInternalPackage(pkg)
	intpkg, err := r.crdClient.InternalV1alpha1().InternalPackages(namespace).Create(ctx, intpkg, *options)
	return r.internalPackageToPackage(intpkg, namespace), err
}

func (r *PackageCRDREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	pkg, err := r.namespacedGet(ctx, namespace, name, options)
	if errors.IsNotFound(err) && namespace != r.globalNamespace && namespace != "" {
		pkg, err = r.namespacedGet(ctx, r.globalNamespace, name, options)
	}
	return r.internalPackageToPackage(pkg, namespace), err
}

func (r *PackageCRDREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)

	namespacedPackagesList, err := r.namespacedList(ctx, namespace, options)
	if err != nil {
		return nil, err
	}

	itemList := namespacedPackagesList.Items
	var globalPackagesList *internalpkgingv1alpha1.InternalPackageList
	if namespace != "" && namespace != r.globalNamespace {
		globalPackagesList, err = r.namespacedList(ctx, r.globalNamespace, options)
		if err != nil {
			return nil, err
		}

		itemList = append(itemList, globalPackagesList.Items...)
	}

	pkgList := datapackaging.PackageList{
		TypeMeta: namespacedPackagesList.TypeMeta,
		ListMeta: namespacedPackagesList.ListMeta,
	}

	// dedup so that namespace packages take precedence over global
	packageNames := make(map[string]struct{})
	for _, pkg := range itemList {
		// all namespaced pkgs come first in the list, so if we have seen one don't append
		if _, seen := packageNames[pkg.Name]; !seen {
			pkgList.Items = append(pkgList.Items, *r.internalPackageToPackage(&pkg, namespace))
			packageNames[pkg.Name] = struct{}{}
		}
	}

	return &pkgList, err
}

func (r *PackageCRDREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	namespace := request.NamespaceValue(ctx)
	pkg, err := r.namespacedGet(ctx, namespace, name, &metav1.GetOptions{})

	if errors.IsNotFound(err) {
		updatedPkg, err := objInfo.UpdatedObject(ctx, &datapackaging.Package{})
		if err != nil {
			return nil, false, err
		}

		if createValidation != nil {
			if err := createValidation(ctx, updatedPkg); err != nil {
				return nil, false, err
			}
		}

		obj, err := r.Create(ctx, updatedPkg, createValidation, &metav1.CreateOptions{TypeMeta: options.TypeMeta, DryRun: options.DryRun, FieldManager: options.FieldManager})
		if err != nil {
			return nil, true, err
		}

		return obj, true, nil
	}

	if err != nil {
		return nil, false, err
	}

	updatedObj, err := objInfo.UpdatedObject(ctx, r.internalPackageToPackage(pkg, pkg.Namespace))
	if err != nil {
		return nil, false, err
	}

	updatedPkg := updatedObj.(*datapackaging.Package)
	errList := validation.ValidatePackage(*updatedPkg)
	if len(errList) != 0 {
		return nil, false, errors.NewInvalid(updatedPkg.GroupVersionKind().GroupKind(), updatedPkg.Name, errList)
	}

	updatedIntPkg := r.packageToInternalPackage(updatedPkg)
	updatedIntPkg, err = r.crdClient.InternalV1alpha1().InternalPackages(namespace).Update(ctx, updatedIntPkg, *options)
	return r.internalPackageToPackage(updatedIntPkg, namespace), false, err
}

func (r *PackageCRDREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	namespace := request.NamespaceValue(ctx)
	pkg, err := r.namespacedGet(ctx, namespace, name, &metav1.GetOptions{})

	if errors.IsNotFound(err) {
		return nil, true, err
	}

	if err != nil {
		return nil, false, err
	}

	if deleteValidation != nil {
		if err := deleteValidation(ctx, pkg); err != nil {
			return nil, true, err
		}
	}

	err = r.crdClient.InternalV1alpha1().InternalPackages(namespace).Delete(ctx, name, *options)
	if err != nil {
		return nil, false, err
	}

	return r.internalPackageToPackage(pkg, namespace), true, nil
}

func (r *PackageCRDREST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	list, err := r.namespacedList(ctx, namespace, listOptions)
	if err != nil {
		return nil, err
	}

	var deletedPackages []datapackaging.Package
	for _, pkg := range list.Items {
		_, _, err := r.Delete(ctx, pkg.Name, deleteValidation, options)
		if err != nil {
			break
		}
		// Although intPkg => Pkg can return nil when intpkg is nil, intpkg here will
		// never be nil, thanks to the type system, so we are ok to deref directly
		deletedPackages = append(deletedPackages, *r.internalPackageToPackage(&pkg, namespace))
	}
	return &datapackaging.PackageList{Items: deletedPackages}, err
}

func (r *PackageCRDREST) Watch(ctx context.Context, options *internalversion.ListOptions) (watch.Interface, error) {
	namespace := request.NamespaceValue(ctx)
	watcher, err := r.namespacedWatch(ctx, namespace, options)
	if errors.IsNotFound(err) && namespace != r.globalNamespace {
		watcher, err = r.namespacedWatch(ctx, r.globalNamespace, options)
	}

	if err != nil {
		return nil, err
	}

	return watchers.NewTranslationWatcher(r.translateFunc(namespace), r.filterFunc(), watcher), nil
}

func (r *PackageCRDREST) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table
	fn := func(obj runtime.Object) error {
		pkg := obj.(*datapackaging.Package)
		table.Rows = append(table.Rows, metav1.TableRow{
			Cells: []interface{}{
				pkg.Name, pkg.Spec.DisplayName,
				r.format(strings.Join(pkg.Spec.Categories, ",")),
				r.format(pkg.Spec.ShortDescription),
				time.Since(pkg.ObjectMeta.CreationTimestamp.Time).Round(1 * time.Second).String(),
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
			{Name: "Name", Type: "string", Format: "name", Description: "Package resource name"},
			{Name: "Display Name", Type: "string", Description: "User facing package name"},
			{Name: "Categories", Type: "string", Description: "Package description"},
			{Name: "Short Description", Type: "array", Description: "Package categories"},
			{Name: "Age", Type: "date", Description: "Time since resource creation"},
		}
	}
	return &table, nil
}

func (r *PackageCRDREST) namespacedList(ctx context.Context, namespace string, options *internalversion.ListOptions) (*internalpkgingv1alpha1.InternalPackageList, error) {
	return r.crdClient.InternalV1alpha1().InternalPackages(namespace).List(ctx, r.internalToMetaListOpts(*options))
}

func (r *PackageCRDREST) namespacedGet(ctx context.Context, namespace, name string, options *metav1.GetOptions) (*internalpkgingv1alpha1.InternalPackage, error) {
	return r.crdClient.InternalV1alpha1().InternalPackages(namespace).Get(ctx, name, *options)
}

func (r *PackageCRDREST) namespacedWatch(ctx context.Context, namespace string, options *internalversion.ListOptions) (watch.Interface, error) {
	return r.crdClient.InternalV1alpha1().InternalPackages(namespace).Watch(ctx, r.internalToMetaListOpts(*options))
}

func (r *PackageCRDREST) internalToMetaListOpts(options internalversion.ListOptions) metav1.ListOptions {
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

func (r *PackageCRDREST) internalPackageToPackage(intpkg *internalpkgingv1alpha1.InternalPackage, namespace string) *datapackaging.Package {
	if intpkg == nil {
		return nil
	}

	pkg := (*datapackaging.Package)(intpkg)
	for i := range pkg.ManagedFields {
		if pkg.ManagedFields[i].APIVersion == internalpkgingv1alpha1.SchemeGroupVersion.Identifier() {
			pkg.ManagedFields[i].APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}
	pkg.TypeMeta.Kind = packageName
	pkg.TypeMeta.APIVersion = datapkgingv1alpha1.SchemeGroupVersion.Identifier()
	if namespace != "" {
		pkg.Namespace = namespace
	}

	return pkg
}

func (r *PackageCRDREST) packageToInternalPackage(pkg *datapackaging.Package) *internalpkgingv1alpha1.InternalPackage {
	if pkg == nil {
		return nil
	}

	intpkg := (*internalpkgingv1alpha1.InternalPackage)(pkg)
	for i := range intpkg.ManagedFields {
		if intpkg.ManagedFields[i].APIVersion == datapkgingv1alpha1.SchemeGroupVersion.Identifier() {
			intpkg.ManagedFields[i].APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()
		}
	}
	pkg.TypeMeta.Kind = internalPackageName
	pkg.TypeMeta.APIVersion = internalpkgingv1alpha1.SchemeGroupVersion.Identifier()
	return intpkg
}

func (r *PackageCRDREST) internalPackageListToPackageList(list *internalpkgingv1alpha1.InternalPackageList, namespace string) *datapackaging.PackageList {
	if list == nil {
		return nil
	}

	pList := datapackaging.PackageList{
		TypeMeta: list.TypeMeta,
		ListMeta: list.ListMeta,
	}

	for _, item := range list.Items {
		pList.Items = append(pList.Items, *r.internalPackageToPackage(&item, namespace))
	}

	return &pList
}

func (r *PackageCRDREST) translateFunc(namespace string) func(evt watch.Event) watch.Event {
	return func(evt watch.Event) watch.Event {
		if intpkg, ok := evt.Object.(*internalpkgingv1alpha1.InternalPackage); ok {
			evt.Object = r.internalPackageToPackage(intpkg, namespace)
		}
		return evt
	}
}

func (r *PackageCRDREST) filterFunc() func(evt watch.Event) bool {
	return func(evt watch.Event) bool {
		return true
	}
}

func (r *PackageCRDREST) format(in string) string {
	if len(in) > 50 {
		return in[:47] + "..."
	}
	return in
}
