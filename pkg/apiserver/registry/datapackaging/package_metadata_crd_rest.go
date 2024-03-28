// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package datapackaging

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/validation"
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
	"k8s.io/client-go/kubernetes"
)

const excludeGlobalPackagesAnn = "kapp-controller.carvel.dev/exclude-global-packages"

// PackageMetadataCRDREST is a rest implementation that proxies the rest endpoints provided by
// CRDs. This will allow us to introduce the api server without the
// complexities associated with custom storage options for now.
type PackageMetadataCRDREST struct {
	crdClient       installclient.Interface
	nsClient        kubernetes.Interface
	globalNamespace string
}

var (
	_ rest.StandardStorage    = &PackageMetadataCRDREST{}
	_ rest.ShortNamesProvider = &PackageMetadataCRDREST{}
)

func NewPackageMetadataCRDREST(crdClient installclient.Interface, nsClient kubernetes.Interface, globalNS string) *PackageMetadataCRDREST {
	return &PackageMetadataCRDREST{crdClient, nsClient, globalNS}
}

func (r *PackageMetadataCRDREST) ShortNames() []string {
	return []string{"pkgm"}
}

// GetSingularName returns singular name of resources.
func (r *PackageMetadataCRDREST) GetSingularName() string {
	return string("packagemetadata")
}

func (r *PackageMetadataCRDREST) New() runtime.Object {
	return &datapackaging.PackageMetadata{}
}

// Destroy satisfies the apiserver.registry.rest.Storage interface but it's a no-op
func (r *PackageMetadataCRDREST) Destroy() {
	return
}

func (r *PackageMetadataCRDREST) NewList() runtime.Object {
	return &datapackaging.PackageMetadataList{}
}

func (r *PackageMetadataCRDREST) NamespaceScoped() bool {
	return true
}

func (r *PackageMetadataCRDREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageMetadataStorageClient(r.crdClient, NewPackageMetadataTranslator(namespace))

	// Run Validations
	if createValidation != nil {
		if err := createValidation(ctx, obj); err != nil {
			return nil, err
		}
	}

	pkg := obj.(*datapackaging.PackageMetadata)
	errs := validation.ValidatePackageMetadata(*pkg)
	if len(errs) != 0 {
		return nil, errors.NewInvalid(pkg.GroupVersionKind().GroupKind(), pkg.Name, errs)
	}

	// Update the data store
	return client.Create(ctx, namespace, pkg, *options)
}

func (r *PackageMetadataCRDREST) shouldFetchGlobal(ctx context.Context, namespace string) bool {
	ns, err := r.nsClient.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		return false
	}
	_, exclude := ns.ObjectMeta.Annotations[excludeGlobalPackagesAnn]
	return namespace != r.globalNamespace && namespace != "" && !exclude
}

func (r *PackageMetadataCRDREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageMetadataStorageClient(r.crdClient, NewPackageMetadataTranslator(namespace))

	// Check targeted namespace
	pkg, err := client.Get(ctx, namespace, name, *options)

	if errors.IsNotFound(err) && r.shouldFetchGlobal(ctx, namespace) {
		// check global namespace
		pkg, err = client.Get(ctx, r.globalNamespace, name, *options)
	}
	return pkg, err
}

func (r *PackageMetadataCRDREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageMetadataStorageClient(r.crdClient, NewPackageMetadataTranslator(namespace))

	var pkgMetas []datapackaging.PackageMetadata
	if r.shouldFetchGlobal(ctx, namespace) {
		globalPackagesList, err := client.List(ctx, r.globalNamespace, r.internalToMetaListOpts(*options))
		if err != nil {
			return nil, err
		}
		pkgMetas = globalPackagesList.Items
	}

	// fetch list of namespaced packages (ns could be "")
	namespacedPackageMetaList, err := client.List(ctx, namespace, r.internalToMetaListOpts(*options))
	if err != nil {
		return nil, err
	}
	namespacedPackageMetas := namespacedPackageMetaList.Items

	pkgMetaIndex := make(map[string]int)
	for i, pkgMeta := range pkgMetas {
		// identifier for package will be namespace/name
		identifier := pkgMeta.Namespace + "/" + pkgMeta.Name
		pkgMetaIndex[identifier] = i
	}

	for _, pkgMeta := range namespacedPackageMetas {
		// identifier for package will be namespace/name
		identifier := pkgMeta.Namespace + "/" + pkgMeta.Name
		if index, found := pkgMetaIndex[identifier]; found {
			pkgMetas[index] = pkgMeta
		} else {
			pkgMetas = append(pkgMetas, pkgMeta)
		}
	}

	packageList := &datapackaging.PackageMetadataList{
		TypeMeta: namespacedPackageMetaList.TypeMeta,
		ListMeta: namespacedPackageMetaList.ListMeta,
		Items:    pkgMetas,
	}

	return packageList, err
}

func (r *PackageMetadataCRDREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageMetadataStorageClient(r.crdClient, NewPackageMetadataTranslator(namespace))

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
			pkg = &datapackaging.PackageMetadata{}
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

		updatedPkg := updatedObj.(*datapackaging.PackageMetadata)
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

	updatedPkg := updatedObj.(*datapackaging.PackageMetadata)
	errList := validation.ValidatePackageMetadata(*updatedPkg)
	if len(errList) != 0 {
		return nil, false, errors.NewInvalid(updatedPkg.GroupVersionKind().GroupKind(), updatedPkg.Name, errList)
	}

	updatedPkg, err = client.Update(ctx, namespace, updatedPkg, *options)
	return updatedPkg, false, err
}

func (r *PackageMetadataCRDREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageMetadataStorageClient(r.crdClient, NewPackageMetadataTranslator(namespace))

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

	err = client.Delete(ctx, namespace, name, *options)
	if err != nil {
		return nil, false, err
	}

	return pkg, true, nil
}

func (r *PackageMetadataCRDREST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageMetadataStorageClient(r.crdClient, NewPackageMetadataTranslator(namespace))

	list, err := client.List(ctx, namespace, r.internalToMetaListOpts(*listOptions))
	if err != nil {
		return nil, err
	}

	// check to see if we are deleting all the global packages. This isnt a great way to do this
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
			return nil, errors.NewInternalError(fmt.Errorf("Removing global packages: %v", err))
		}
	}

	var deletedPackages []datapackaging.PackageMetadata
	for _, pkg := range list.Items {
		// use crd delete for validations
		_, _, err := r.Delete(ctx, pkg.Name, deleteValidation, options)
		if err != nil && !errors.IsNotFound(err) {
			break
		}
		deletedPackages = append(deletedPackages, pkg)
	}

	return &datapackaging.PackageMetadataList{Items: deletedPackages}, err
}

func (r *PackageMetadataCRDREST) Watch(ctx context.Context, options *internalversion.ListOptions) (watch.Interface, error) {
	namespace := request.NamespaceValue(ctx)
	client := NewPackageMetadataStorageClient(r.crdClient, NewPackageMetadataTranslator(namespace))

	watcher, err := client.Watch(ctx, namespace, r.internalToMetaListOpts(*options))
	if errors.IsNotFound(err) && r.shouldFetchGlobal(ctx, namespace) {
		watcher, err = client.Watch(ctx, r.globalNamespace, r.internalToMetaListOpts(*options))
	}

	return watcher, err
}

func (r *PackageMetadataCRDREST) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table
	fn := func(obj runtime.Object) error {
		pkg := obj.(*datapackaging.PackageMetadata)
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
			{Name: "Name", Type: "string", Format: "name", Description: "PackageMetadata resource name"},
			{Name: "Display Name", Type: "string", Description: "User facing package name"},
			{Name: "Categories", Type: "string", Description: "PackageMetadata description"},
			{Name: "Short Description", Type: "array", Description: "PackageMetadata categories"},
			{Name: "Age", Type: "date", Description: "Time since resource creation"},
		}
	}
	return &table, nil
}

func (r *PackageMetadataCRDREST) internalToMetaListOpts(options internalversion.ListOptions) metav1.ListOptions {
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

func (r *PackageMetadataCRDREST) deleteGlobalPackagesFromNS(ctx context.Context, ns string) error {
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

func (r *PackageMetadataCRDREST) format(in string) string {
	if len(in) > 50 {
		return in[:47] + "..."
	}
	return in
}
