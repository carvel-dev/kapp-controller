// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package datapackaging

import (
	"context"
	"time"

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
	client := NewPackageVersionStorageClient(r.crdClient, NewPackageVersionTranslator(namespace))

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

	return client.Create(ctx, namespace, pkgVersion, *options)
}

func (r *PackageVersionCRDREST) shouldFetchGlobal(namespace string) bool {
	return namespace != r.globalNamespace && namespace != ""
}

func (r *PackageVersionCRDREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageVersionStorageClient(r.crdClient, NewPackageVersionTranslator(namespace))

	pkgVersion, err := client.Get(ctx, namespace, name, *options)
	if errors.IsNotFound(err) && r.shouldFetchGlobal(namespace) {
		pkgVersion, err = client.Get(ctx, r.globalNamespace, name, *options)
	}
	return pkgVersion, err
}

func (r *PackageVersionCRDREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageVersionStorageClient(r.crdClient, NewPackageVersionTranslator(namespace))

	// field selector isnt supported by CRD's so reset it, we will apply it later
	fs := options.FieldSelector
	options.FieldSelector = fields.Everything()

	namespacedPkgVersionList, err := client.List(ctx, namespace, r.internalToMetaListOpts(*options))
	if err != nil {
		return nil, err
	}
	namespacedPkgVersions := namespacedPkgVersionList.Items

	var globalPkgVersions []datapackaging.PackageVersion
	if r.shouldFetchGlobal(namespace) {
		globalPkgVersionList, err := client.List(ctx, r.globalNamespace, r.internalToMetaListOpts(*options))
		if err != nil {
			return nil, err
		}
		globalPkgVersions = globalPkgVersionList.Items
	}

	packageVersionsMap := make(map[string]datapackaging.PackageVersion)
	for _, pkgVersion := range globalPkgVersions {
		identifier := pkgVersion.Namespace + "/" + pkgVersion.Spec.PackageName + "." + pkgVersion.Spec.Version
		packageVersionsMap[identifier] = pkgVersion
	}

	for _, pkgVersion := range namespacedPkgVersions {
		identifier := pkgVersion.Namespace + "/" + pkgVersion.Spec.PackageName + "." + pkgVersion.Spec.Version
		packageVersionsMap[identifier] = pkgVersion
	}

	pkgVersionList := &datapackaging.PackageVersionList{
		TypeMeta: namespacedPkgVersionList.TypeMeta,
		ListMeta: namespacedPkgVersionList.ListMeta,
	}

	for _, v := range packageVersionsMap {
		pkgVersionList.Items = append(pkgVersionList.Items, v)
	}

	return r.applySelector(pkgVersionList, fs), err
}

func (r *PackageVersionCRDREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageVersionStorageClient(r.crdClient, NewPackageVersionTranslator(namespace))

	pkgVersion, err := client.Get(ctx, namespace, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		// Because kubetl does a get before sending an update, the presence
		// of a global package may cause it to send a patch request, even though
		// the package doesn't exist in the namespace. To service this, we must check
		// if the package exists globally and then patch that instead of patching an empty
		// package. If we try patching an empty obj the patch UpdatedObjectInfo will blow up.
		patchingGlobal := true
		pkgVersion, err := client.Get(ctx, r.globalNamespace, name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			pkgVersion = &datapackaging.PackageVersion{}
			patchingGlobal = false
		}

		updatedObj, err := objInfo.UpdatedObject(ctx, pkgVersion)
		if err != nil {
			return nil, false, err
		}

		if createValidation != nil {
			if err := createValidation(ctx, updatedObj); err != nil {
				return nil, false, err
			}
		}

		updatedPkgVersion := updatedObj.(*datapackaging.PackageVersion)
		if patchingGlobal {
			// we have to do this in case we are "patching" a global package
			annotations := updatedPkgVersion.ObjectMeta.Annotations
			labels := updatedPkgVersion.ObjectMeta.Labels
			updatedPkgVersion.ObjectMeta = metav1.ObjectMeta{}
			updatedPkgVersion.ObjectMeta.Name = name
			updatedPkgVersion.ObjectMeta.Namespace = namespace
			updatedPkgVersion.ObjectMeta.Annotations = annotations
			updatedPkgVersion.ObjectMeta.Labels = labels
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

	updatedObj, err := objInfo.UpdatedObject(ctx, pkgVersion)
	if err != nil {
		return nil, false, err
	}

	updatedPkgVersion := updatedObj.(*datapackaging.PackageVersion)
	errs := validation.ValidatePackageVersion(*updatedPkgVersion)
	if len(errs) != 0 {
		return nil, false, errors.NewInvalid(updatedPkgVersion.GroupVersionKind().GroupKind(), updatedPkgVersion.Name, errs)
	}

	updatedPkgVersion, err = client.Update(ctx, namespace, updatedPkgVersion, *options)
	return updatedPkgVersion, false, err
}

func (r *PackageVersionCRDREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageVersionStorageClient(r.crdClient, NewPackageVersionTranslator(namespace))

	pkgVersion, err := client.Get(ctx, namespace, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, true, err
	}

	if err != nil {
		return nil, false, err
	}

	if deleteValidation != nil {
		if err := deleteValidation(ctx, pkgVersion); err != nil {
			return nil, true, err
		}
	}

	err = client.Delete(ctx, namespace, pkgVersion.Name, *options)
	if err != nil {
		return nil, false, err
	}

	return pkgVersion, true, nil
}

func (r *PackageVersionCRDREST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageVersionStorageClient(r.crdClient, NewPackageVersionTranslator(namespace))

	list, err := client.List(ctx, namespace, r.internalToMetaListOpts(*listOptions))
	if err != nil {
		return nil, err
	}

	var deletedPackages []datapackaging.PackageVersion
	for _, pv := range list.Items {
		// use crd delete for validations
		_, _, err := r.Delete(ctx, pv.Name, deleteValidation, options)
		if err != nil {
			break
		}
		deletedPackages = append(deletedPackages, pv)
	}
	return &datapackaging.PackageVersionList{Items: deletedPackages}, err
}

func (r *PackageVersionCRDREST) Watch(ctx context.Context, options *internalversion.ListOptions) (watch.Interface, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageVersionStorageClient(r.crdClient, NewPackageVersionTranslator(namespace))

	watcher, err := client.Watch(ctx, namespace, r.internalToMetaListOpts(*options))
	if errors.IsNotFound(err) && namespace != r.globalNamespace {
		watcher, err = client.Watch(ctx, r.globalNamespace, r.internalToMetaListOpts(*options))
	}

	return watcher, err
}

func (r *PackageVersionCRDREST) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table
	fn := func(obj runtime.Object) error {
		pkgVersion := obj.(*datapackaging.PackageVersion)
		table.Rows = append(table.Rows, metav1.TableRow{
			Cells: []interface{}{
				pkgVersion.Name, pkgVersion.Spec.PackageName,
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

func (r *PackageVersionCRDREST) applySelector(list *datapackaging.PackageVersionList, selector fields.Selector) *datapackaging.PackageVersionList {
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
