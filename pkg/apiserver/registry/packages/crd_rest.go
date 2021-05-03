// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packages

import (
	"context"
	"time"

	installv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/packages"
	pkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/packages/v1alpha1"
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

// CRDREST is a rest implementation that proxies the rest endpoints provided by
// CRDs. This will allow us to introduce the api server without the
// complexities associated with custom storage options for now.
type CRDREST struct {
	crdClient installclient.Interface
}

var (
	_ rest.StandardStorage    = &CRDREST{}
	_ rest.ShortNamesProvider = &CRDREST{}
)

const (
	packageName         = "Package"
	internalPackageName = "InternalPackage"
)

func NewCRDREST(crdClient installclient.Interface) *CRDREST {
	return &CRDREST{crdClient}
}

func (r *CRDREST) ShortNames() []string {
	return []string{"pkg"}
}

func (r *CRDREST) New() runtime.Object {
	return &packages.Package{}
}

func (r *CRDREST) NewList() runtime.Object {
	return &packages.PackageList{}
}

func (r *CRDREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	if createValidation != nil {
		if err := createValidation(ctx, obj); err != nil {
			return nil, err
		}
	}

	pkg := obj.(*packages.Package)
	intpkg := r.packageToInternalPackage(pkg)
	intpkg, err := r.crdClient.InstallV1alpha1().InternalPackages().Create(ctx, intpkg, *options)
	return r.internalPackageToPackage(intpkg), err
}

func (r *CRDREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
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

		updatedIntPkg := r.packageToInternalPackage(updatedPkg.(*packages.Package))
		updatedIntPkg, err = r.crdClient.InstallV1alpha1().InternalPackages().Update(ctx, updatedIntPkg, *options)
		return r.internalPackageToPackage(updatedIntPkg), true, nil
	}

	if err != nil {
		return nil, false, err
	}

	updatedPkg, err := objInfo.UpdatedObject(ctx, pkg)
	if err != nil {
		return nil, false, err
	}
	updatedIntPkg := r.packageToInternalPackage(updatedPkg.(*packages.Package))
	updatedIntPkg, err = r.crdClient.InstallV1alpha1().InternalPackages().Update(ctx, updatedIntPkg, *options)
	return r.internalPackageToPackage(updatedIntPkg), false, err
}

func (r *CRDREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
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

func (r *CRDREST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
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

func (r *CRDREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	intpkg, err := r.crdClient.InstallV1alpha1().InternalPackages().Get(ctx, name, *options)
	return r.internalPackageToPackage(intpkg), err
}

func (r *CRDREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
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

func (r *CRDREST) NamespaceScoped() bool {
	return false
}

func (r *CRDREST) Watch(ctx context.Context, options *internalversion.ListOptions) (watch.Interface, error) {
	watcher, err := r.crdClient.InstallV1alpha1().InternalPackages().Watch(ctx, r.internalToMetaListOpts(*options))
	return watchers.NewTranslationWatcher(r.translate, watcher), err
}

func (r *CRDREST) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table
	fn := func(obj runtime.Object) error {
		pkg := obj.(*packages.Package)
		table.Rows = append(table.Rows, metav1.TableRow{
			Cells:  []interface{}{pkg.Name, pkg.Spec.PublicName, pkg.Spec.Version, time.Since(pkg.ObjectMeta.CreationTimestamp.Time).Round(1 * time.Second).String()},
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
			{Name: "PublicName", Type: "string", Description: "User facing package name"},
			{Name: "Version", Type: "string", Description: "Package version"},
			{Name: "Age", Type: "date", Description: "Time since resource creation"},
		}
	}
	return &table, nil
}

func (r *CRDREST) internalToMetaListOpts(options internalversion.ListOptions) metav1.ListOptions {
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

func (r *CRDREST) internalPackageToPackage(intpkg *installv1alpha1.InternalPackage) *packages.Package {
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

func (r *CRDREST) packageToInternalPackage(pkg *packages.Package) *installv1alpha1.InternalPackage {
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

func (r *CRDREST) translate(evt watch.Event) watch.Event {
	if evt.Object != nil {
		intpkg := evt.Object.(*installv1alpha1.InternalPackage)
		evt.Object = r.internalPackageToPackage(intpkg)
	}
	return evt
}
