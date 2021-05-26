// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packages

import (
	"context"
	"strings"
	"time"

	installv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/packages"
	pkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/packages/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/packages/validation"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/watchers"
	installclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/registry/rest"
)

// PackageCRDREST is a rest implementation that proxies the rest endpoints provided by
// CRDs. This will allow us to introduce the api server without the
// complexities associated with custom storage options for now.
type PackageCRDREST struct {
	crdClient installclient.Interface
}

var (
	_ rest.StandardStorage    = &PackageCRDREST{}
	_ rest.ShortNamesProvider = &PackageCRDREST{}
)

const (
	packageName         = "Package"
	internalPackageName = "InternalPackage"
)

func NewPackageCRDREST(crdClient installclient.Interface) *PackageCRDREST {
	return &PackageCRDREST{crdClient}
}

func (r *PackageCRDREST) ShortNames() []string {
	return []string{"pkg"}
}

func (r *PackageCRDREST) New() runtime.Object {
	return &packages.Package{}
}

func (r *PackageCRDREST) NewList() runtime.Object {
	return &packages.PackageList{}
}

func (r *PackageCRDREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	if createValidation != nil {
		if err := createValidation(ctx, obj); err != nil {
			return nil, err
		}
	}

	pkg := obj.(*packages.Package)
	errs := validation.ValidatePackage(*pkg)
	if len(errs) != 0 {
		return nil, errors.NewInvalid(pkg.GroupVersionKind().GroupKind(), pkg.Name, errs)
	}

	intpkg := r.packageToInternalPackage(pkg)
	intpkg, err := r.crdClient.InstallV1alpha1().InternalPackages().Create(ctx, intpkg, *options)
	return r.internalPackageToPackage(intpkg), err
}

func (r *PackageCRDREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	pkg, err := r.Get(ctx, name, &metav1.GetOptions{})

	if errors.IsNotFound(err) {
		updatedPkg, err := objInfo.UpdatedObject(ctx, &packages.Package{})
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

	updatedObj, err := objInfo.UpdatedObject(ctx, pkg)
	if err != nil {
		return nil, false, err
	}

	updatedPkg := updatedObj.(*packages.Package)
	errList := validation.ValidatePackage(*updatedPkg)
	if len(errList) != 0 {
		return nil, false, errors.NewInvalid(updatedPkg.GroupVersionKind().GroupKind(), updatedPkg.Name, errList)
	}

	updatedIntPkg := r.packageToInternalPackage(updatedPkg)
	updatedIntPkg, err = r.crdClient.InstallV1alpha1().InternalPackages().Update(ctx, updatedIntPkg, *options)
	return r.internalPackageToPackage(updatedIntPkg), false, err
}

func (r *PackageCRDREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	intPkg, err := r.crdClient.InstallV1alpha1().InternalPackages().Get(ctx, name, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		return nil, true, err
	}

	if err != nil {
		return nil, false, err
	}

	if deleteValidation != nil {
		if err := deleteValidation(ctx, intPkg); err != nil {
			return nil, true, err
		}
	}

	err = r.crdClient.InstallV1alpha1().InternalPackages().Delete(ctx, name, *options)
	if err != nil {
		return nil, false, err
	}

	return nil, true, nil
}

func (r *PackageCRDREST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	objs, err := r.List(ctx, listOptions)
	if err != nil {
		return nil, err
	}

	var deletedPackages []packages.Package
	for _, obj := range objs.(*packages.PackageList).Items {
		_, _, err := r.Delete(ctx, obj.Name, deleteValidation, options)
		if err != nil {
			break
		}
		deletedPackages = append(deletedPackages, obj)
	}
	return &packages.PackageList{Items: deletedPackages}, err
}

func (r *PackageCRDREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	intpkg, err := r.crdClient.InstallV1alpha1().InternalPackages().Get(ctx, name, *options)
	return r.internalPackageToPackage(intpkg), err
}

func (r *PackageCRDREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	list, err := r.crdClient.InstallV1alpha1().InternalPackages().List(ctx, r.internalToMetaListOpts(*options))
	pkgList := packages.PackageList{
		TypeMeta: list.TypeMeta,
		ListMeta: list.ListMeta,
	}
	for _, intpkg := range list.Items {
		pkgList.Items = append(pkgList.Items, *r.internalPackageToPackage(&intpkg))
	}

	return &pkgList, err
}

func (r *PackageCRDREST) NamespaceScoped() bool {
	return false
}

func (r *PackageCRDREST) Watch(ctx context.Context, options *internalversion.ListOptions) (watch.Interface, error) {
	watcher, err := r.crdClient.InstallV1alpha1().InternalPackages().Watch(ctx, r.internalToMetaListOpts(*options))
	return watchers.NewTranslationWatcher(r.translateFunc(), r.filterFunc(), watcher), err
}

func (r *PackageCRDREST) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table
	fn := func(obj runtime.Object) error {
		pkg := obj.(*packages.Package)
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

func (r *PackageCRDREST) internalPackageToPackage(intpkg *installv1alpha1.InternalPackage) *packages.Package {
	pkg := (*packages.Package)(intpkg)
	for i := range pkg.ManagedFields {
		mf := pkg.ManagedFields[i]
		if mf.APIVersion == installv1alpha1.SchemeGroupVersion.Identifier() {
			mf.APIVersion = pkgv1alpha1.SchemeGroupVersion.Identifier()
		}
	}
	pkg.TypeMeta.Kind = packageName
	pkg.TypeMeta.APIVersion = pkgv1alpha1.SchemeGroupVersion.Identifier()
	return pkg
}

func (r *PackageCRDREST) packageToInternalPackage(pkg *packages.Package) *installv1alpha1.InternalPackage {
	intpkg := (*installv1alpha1.InternalPackage)(pkg)
	for i := range intpkg.ManagedFields {
		if intpkg.ManagedFields[i].APIVersion == pkgv1alpha1.SchemeGroupVersion.Identifier() {
			intpkg.ManagedFields[i].APIVersion = installv1alpha1.SchemeGroupVersion.Identifier()
		}
	}
	pkg.TypeMeta.Kind = internalPackageName
	pkg.TypeMeta.APIVersion = installv1alpha1.SchemeGroupVersion.Identifier()
	return intpkg
}

func (r *PackageCRDREST) translateFunc() func(evt watch.Event) watch.Event {
	return func(evt watch.Event) watch.Event {
		if intpkg, ok := evt.Object.(*installv1alpha1.InternalPackage); ok {
			evt.Object = r.internalPackageToPackage(intpkg)
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
