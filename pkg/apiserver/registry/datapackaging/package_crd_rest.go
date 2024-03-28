// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package datapackaging

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/validation"
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
	"k8s.io/client-go/kubernetes"
)

// PackageCRDREST is a rest implementation that proxies the rest endpoints provided by
// CRDs. This will allow us to introduce the api server without the
// complexities associated with custom storage options for now.
type PackageCRDREST struct {
	crdClient       installclient.Interface
	nsClient        kubernetes.Interface
	globalNamespace string
	logger          logr.Logger
}

var (
	_ rest.StandardStorage    = &PackageCRDREST{}
	_ rest.ShortNamesProvider = &PackageCRDREST{}
)

// NewPackageCRDREST creates a new instance of the PackageCRDREST type
func NewPackageCRDREST(crdClient installclient.Interface, nsClient kubernetes.Interface, globalNS string, logger logr.Logger) *PackageCRDREST {
	return &PackageCRDREST{crdClient, nsClient, globalNS, logger}
}

func (r *PackageCRDREST) ShortNames() []string {
	return []string{"pkg"}
}

// GetSingularName returns singular name of resources.
func (r *PackageCRDREST) GetSingularName() string {
	return string("package")
}

func (r *PackageCRDREST) NamespaceScoped() bool {
	return true
}

func (r *PackageCRDREST) New() runtime.Object {
	return &datapackaging.Package{}
}

// Destroy satisfies the apiserver.registry.storage interface but it's a no-op
func (r *PackageCRDREST) Destroy() {
	return
}

func (r *PackageCRDREST) NewList() runtime.Object {
	return &datapackaging.PackageList{}
}

func (r *PackageCRDREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageStorageClient(r.crdClient, NewPackageTranslator(namespace, r.logger))

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

	return client.Create(ctx, namespace, pkg, *options)
}

func (r *PackageCRDREST) shouldFetchGlobal(ctx context.Context, namespace string) bool {
	ns, err := r.nsClient.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		return false
	}
	_, exclude := ns.ObjectMeta.Annotations[excludeGlobalPackagesAnn]
	return namespace != r.globalNamespace && namespace != "" && !exclude
}

func (r *PackageCRDREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageStorageClient(r.crdClient, NewPackageTranslator(namespace, r.logger))

	pkg, err := client.Get(ctx, namespace, name, *options)
	if errors.IsNotFound(err) && r.shouldFetchGlobal(ctx, namespace) {
		pkg, err = client.Get(ctx, r.globalNamespace, name, *options)
	}
	return pkg, err
}

func (r *PackageCRDREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageStorageClient(r.crdClient, NewPackageTranslator(namespace, r.logger))

	// field selector isnt supported by CRD's so reset it, we will apply it later
	fs := options.FieldSelector
	options.FieldSelector = fields.Everything()

	var pkgs []datapackaging.Package
	if r.shouldFetchGlobal(ctx, namespace) {
		globalPkgList, err := client.List(ctx, r.globalNamespace, r.internalToMetaListOpts(*options))
		if err != nil {
			return nil, err
		}
		pkgs = globalPkgList.Items
	}

	namespacedPkgList, err := client.List(ctx, namespace, r.internalToMetaListOpts(*options))
	if err != nil {
		return nil, err
	}
	namespacedPkgs := namespacedPkgList.Items

	pkgIndex := make(map[string]int)
	for i, pkg := range pkgs {
		identifier := pkg.Namespace + "/" + pkg.Spec.RefName + "." + pkg.Spec.Version
		pkgIndex[identifier] = i
	}

	// if exists as global, overwrite, else append
	for _, pkg := range namespacedPkgs {
		identifier := pkg.Namespace + "/" + pkg.Spec.RefName + "." + pkg.Spec.Version
		if index, found := pkgIndex[identifier]; found {
			pkgs[index] = pkg
		} else {
			pkgs = append(pkgs, pkg)
		}
	}

	pkgList := &datapackaging.PackageList{
		TypeMeta: namespacedPkgList.TypeMeta,
		ListMeta: namespacedPkgList.ListMeta,
		Items:    pkgs,
	}

	return r.applySelector(pkgList, fs), err
}

func (r *PackageCRDREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageStorageClient(r.crdClient, NewPackageTranslator(namespace, r.logger))

	pkg, err := client.Get(ctx, namespace, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		// Because kubetl does a get before sending an update, the presence
		// of a global package may cause it to send a patch request, even though
		// the package doesn't exist in the namespace. To service this, we must check
		// if the package exists globally and then patch that instead of patching an empty
		// package. If we try patching an empty obj the patch UpdatedObjectInfo will blow up.
		patchingGlobal := true
		pkg, err := client.Get(ctx, r.globalNamespace, name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			pkg = &datapackaging.Package{}
			patchingGlobal = false
		}

		updatedObj, err := objInfo.UpdatedObject(ctx, pkg)
		if err != nil {
			return nil, false, err
		}

		if createValidation != nil {
			if err := createValidation(ctx, updatedObj); err != nil {
				return nil, false, err
			}
		}

		updatedPkg := updatedObj.(*datapackaging.Package)
		if patchingGlobal {
			// we have to do this in case we are "patching" a global package
			annotations := updatedPkg.ObjectMeta.Annotations
			labels := updatedPkg.ObjectMeta.Labels
			updatedPkg.ObjectMeta = metav1.ObjectMeta{}
			updatedPkg.ObjectMeta.Name = name
			updatedPkg.ObjectMeta.Namespace = namespace
			updatedPkg.ObjectMeta.Annotations = annotations
			updatedPkg.ObjectMeta.Labels = labels
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

	updatedPkg := updatedObj.(*datapackaging.Package)
	errs := validation.ValidatePackage(*updatedPkg)
	if len(errs) != 0 {
		return nil, false, errors.NewInvalid(updatedPkg.GroupVersionKind().GroupKind(), updatedPkg.Name, errs)
	}

	updatedPkg, err = client.Update(ctx, namespace, updatedPkg, *options)
	return updatedPkg, false, err
}

func (r *PackageCRDREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageStorageClient(r.crdClient, NewPackageTranslator(namespace, r.logger))

	pkg, err := client.Get(ctx, namespace, name, metav1.GetOptions{})
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

	err = client.Delete(ctx, namespace, pkg.Name, *options)
	if err != nil {
		return nil, false, err
	}

	return pkg, true, nil
}

func (r *PackageCRDREST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageStorageClient(r.crdClient, NewPackageTranslator(namespace, r.logger))

	// clear unsupported field selectors
	fs := listOptions.FieldSelector
	listOptions.FieldSelector = fields.Everything()

	list, err := client.List(ctx, namespace, r.internalToMetaListOpts(*listOptions))
	if err != nil {
		return nil, err
	}

	list = r.applySelector(list, fs)

	deleteAllGlobal := false
	{
		filteredList, err := client.List(ctx, r.globalNamespace, r.internalToMetaListOpts(*listOptions))
		if err != nil {
			return nil, err
		}

		regularList, err := client.List(ctx, r.globalNamespace, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		deleteAllGlobal = len(regularList.Items) == len(filteredList.Items)
	}

	if deleteAllGlobal {
		err := r.deleteGlobalPackagesFromNS(ctx, namespace)
		if err != nil {
			return nil, errors.NewInternalError(fmt.Errorf("Removing global packages from ns '%s': %v", namespace, err))
		}
	}

	var deletedPackages []datapackaging.Package
	for _, pv := range list.Items {
		// use crd delete for validations
		_, _, err := r.Delete(ctx, pv.Name, deleteValidation, options)
		if err != nil {
			break
		}
		deletedPackages = append(deletedPackages, pv)
	}
	return &datapackaging.PackageList{Items: deletedPackages}, err
}

func (r *PackageCRDREST) Watch(ctx context.Context, options *internalversion.ListOptions) (watch.Interface, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageStorageClient(r.crdClient, NewPackageTranslator(namespace, r.logger))

	watcher, err := client.Watch(ctx, namespace, r.internalToMetaListOpts(*options))
	if errors.IsNotFound(err) && namespace != r.globalNamespace {
		watcher, err = client.Watch(ctx, r.globalNamespace, r.internalToMetaListOpts(*options))
	}

	return watcher, err
}

func (r *PackageCRDREST) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table
	fn := func(obj runtime.Object) error {
		pkg := obj.(*datapackaging.Package)
		table.Rows = append(table.Rows, metav1.TableRow{
			Cells: []interface{}{
				pkg.Name, pkg.Spec.RefName,
				pkg.Spec.Version, time.Since(pkg.ObjectMeta.CreationTimestamp.Time).Round(1 * time.Second).String(),
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
			{Name: "Name", Type: "string", Format: "name", Description: "PackageMetadata Version resource name"},
			{Name: "PackageMetadata Name", Type: "string", Format: "name", Description: "Associated PackageMetadata name"},
			{Name: "Version", Type: "string", Description: "Version"},
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

func (r *PackageCRDREST) deleteGlobalPackagesFromNS(ctx context.Context, ns string) error {
	namespace, err := r.nsClient.CoreV1().Namespaces().Get(ctx, ns, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if namespace.ObjectMeta.Annotations == nil {
		namespace.ObjectMeta.Annotations = make(map[string]string)
	}

	namespace.ObjectMeta.Annotations[excludeGlobalPackagesAnn] = ""
	_, err = r.nsClient.CoreV1().Namespaces().Update(ctx, namespace, metav1.UpdateOptions{})
	return err
}

func (r *PackageCRDREST) applySelector(list *datapackaging.PackageList, selector fields.Selector) *datapackaging.PackageList {
	if selector == nil || selector.Empty() {
		return list
	}

	filteredPVs := []datapackaging.Package{}
	for _, pv := range list.Items {
		fieldSet := fields.Set{"spec.refName": pv.Spec.RefName, "metadata.name": pv.Name, "metadata.namespace": pv.Namespace}
		if selector.Matches(fieldSet) {
			filteredPVs = append(filteredPVs, pv)
		}
	}

	list.Items = filteredPVs
	return list
}
