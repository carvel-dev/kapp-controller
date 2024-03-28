// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package datapackaging_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/internalpackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging"
	datapkgreg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/registry/datapackaging"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	cgtesting "k8s.io/client-go/testing"
)

// listing
func TestPackageListIncludesGlobalAndNamespaced(t *testing.T) {
	internalClient := fake.NewSimpleClientset(globalIntPackage(), namespacedIntPackage(), excludedNonGlobalIntPackage())
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgCRDREST := datapkgreg.NewPackageMetadataCRDREST(internalClient, fakeCoreClient, globalNamespace)

	pkgList, err := pkgCRDREST.List(namespacedCtx(nonGlobalNamespace), &internalversion.ListOptions{})
	if err != nil {
		t.Fatalf("Expected list operation to succeed, got: %v", err)
	}

	packageList, ok := pkgList.(*datapackaging.PackageMetadataList)
	if !ok {
		t.Fatalf("Expected list operation to return PackageMetadataList, but got: %v", reflect.TypeOf(pkgList))
	}

	expectedPkgs := []datapackaging.PackageMetadata{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "some.other.ns.carvel.dev",
				Name:      "global-package.carvel.dev",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "some.other.ns.carvel.dev",
				Name:      "namespaced-package.carvel.dev",
			},
		},
	}

	assertPackageListUnorderedEquals(packageList.Items, expectedPkgs, t)
}

func TestPackageListPrefersNamespacedOverGlobal(t *testing.T) {
	// seed client with resources
	internalClient := fake.NewSimpleClientset(globalIntPackage(), overrideIntPackage())
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgCRDREST := datapkgreg.NewPackageMetadataCRDREST(internalClient, fakeCoreClient, globalNamespace)

	// list package versions and verify all of them are there
	pkgList, err := pkgCRDREST.List(namespacedCtx(nonGlobalNamespace), &internalversion.ListOptions{})
	if err != nil {
		t.Fatalf("Expected list operation to succeed, got: %v", err)
	}

	packageList, ok := pkgList.(*datapackaging.PackageMetadataList)
	if !ok {
		t.Fatalf("Expected list operation to return PackageMetadataList, but got: %v", reflect.TypeOf(pkgList))
	}

	expectedPkgs := []datapackaging.PackageMetadata{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: nonGlobalNamespace,
				Name:      "global-package.carvel.dev",
			},
		},
	}

	assertPackageListUnorderedEquals(packageList.Items, expectedPkgs, t)
}

// getting

func TestPackageGetNotPresentInNS(t *testing.T) {
	globalPackage := globalIntPackage()
	name := globalPackage.Name

	internalClient := fake.NewSimpleClientset(globalPackage)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgCRDREST := datapkgreg.NewPackageMetadataCRDREST(internalClient, fakeCoreClient, globalNamespace)

	obj, err := pkgCRDREST.Get(namespacedCtx(nonGlobalNamespace), name, &metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Expected get operation to succeed, got: %v", err)
	}

	pkg, ok := obj.(*datapackaging.PackageMetadata)
	if !ok {
		t.Fatalf("Expected get operation to return PackageMetadata, but got: %v", reflect.TypeOf(pkg))
	}

	if pkg.Name != name || pkg.Spec.DisplayName != "GLOBAL" {
		t.Fatalf("Expected returned package to have name %s and namespace %s, got %s and %s", name, globalNamespace, pkg.Name, pkg.Namespace)
	}
}

func TestPackageGetPresentInOnlyNS(t *testing.T) {
	namespacedPackage := namespacedIntPackage()
	name := namespacedPackage.Name

	internalClient := fake.NewSimpleClientset(namespacedPackage)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgCRDREST := datapkgreg.NewPackageMetadataCRDREST(internalClient, fakeCoreClient, globalNamespace)

	obj, err := pkgCRDREST.Get(namespacedCtx(nonGlobalNamespace), name, &metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Expected get operation to succeed, got: %v", err)
	}

	pkg, ok := obj.(*datapackaging.PackageMetadata)
	if !ok {
		t.Fatalf("Expected get operation to return PackageMetadata, but got: %v", reflect.TypeOf(pkg))
	}

	if pkg.Name != name || pkg.Spec.DisplayName != "NAMESPACED" {
		t.Fatalf("Expected returned package to have name %s and namespace %s, got %s and %s", name, nonGlobalNamespace, pkg.Name, pkg.Namespace)
	}
}

func TestPackageGetNotFound(t *testing.T) {
	namespacedPackage := excludedNonGlobalIntPackage()
	name := namespacedPackage.Name

	internalClient := fake.NewSimpleClientset(namespacedPackage)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgCRDREST := datapkgreg.NewPackageMetadataCRDREST(internalClient, fakeCoreClient, globalNamespace)

	_, err := pkgCRDREST.Get(namespacedCtx(nonGlobalNamespace), name, &metav1.GetOptions{})
	if err == nil {
		t.Fatalf("Expected get operation to fail, but it didn't")
	}

	if !errors.IsNotFound(err) {
		t.Fatalf("Expected a not found error, got: %v", err)
	}

}

func TestPackageGetPreferNS(t *testing.T) {
	overridePackage := overrideIntPackage()
	name := overridePackage.Name

	internalClient := fake.NewSimpleClientset(overridePackage, globalIntPackage())
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgCRDREST := datapkgreg.NewPackageMetadataCRDREST(internalClient, fakeCoreClient, globalNamespace)

	obj, err := pkgCRDREST.Get(namespacedCtx(nonGlobalNamespace), name, &metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Expected get operation to succeed, got: %v", err)
	}

	pkg, ok := obj.(*datapackaging.PackageMetadata)
	if !ok {
		t.Fatalf("Expected get operation to return PackageMetadata, but got: %v", reflect.TypeOf(pkg))
	}

	if pkg.Name != name || pkg.Spec.DisplayName != "OVERRIDE" {
		t.Fatalf("Expected returned package to have name %s and namespace %s, got %s and %s", name, nonGlobalNamespace, pkg.Name, pkg.Namespace)
	}
}

// updating
func TestPackageUpdateDoesntUpdateGlobal(t *testing.T) {
	globalPackage := globalIntPackage()
	namespacedPackage := namespacedIntPackage()
	name := globalPackage.Name
	originalDisplayName := globalPackage.Spec.DisplayName
	newDisplayName := "im-new"

	internalClient := fake.NewSimpleClientset(globalPackage, namespacedPackage)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgCRDREST := datapkgreg.NewPackageMetadataCRDREST(internalClient, fakeCoreClient, globalNamespace)

	obj, created, err := pkgCRDREST.Update(namespacedCtx(nonGlobalNamespace), name, UpdatePackageTestImpl{updateDisplayNameFn(newDisplayName, name)}, nil, nil, false, &metav1.UpdateOptions{})
	if err != nil {
		t.Fatalf("Expected update operation to succeed, got: %v", err)
	}

	if !created {
		t.Fatalf("Expected object to be created")
	}

	pkg, ok := obj.(*datapackaging.PackageMetadata)
	if !ok {
		t.Fatalf("Expected get operation to return PackageMetadata, but got: %v", reflect.TypeOf(pkg))
	}

	if pkg.Name != name || pkg.Namespace != nonGlobalNamespace {
		t.Fatalf("Expected returned package to have name %s and namespace %s, got %s and %s", name, nonGlobalNamespace, pkg.Name, pkg.Namespace)
	}

	if pkg.Spec.DisplayName != newDisplayName {
		t.Fatalf("Expected display name of namespaced package to be updated")
	}

	if globalPackage.Spec.DisplayName != originalDisplayName {
		t.Fatalf("Expected update not to affect the global package, but DisplayName was updated")
	}
}

// scoped to ns, so if can't find create in ns
func TestPackageUpdateCreatesInNS(t *testing.T) {
	globalPackage := globalIntPackage()
	name := globalPackage.Name
	originalDisplayName := globalPackage.Spec.DisplayName
	newDisplayName := "im-new"

	internalClient := fake.NewSimpleClientset(globalPackage)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgCRDREST := datapkgreg.NewPackageMetadataCRDREST(internalClient, fakeCoreClient, globalNamespace)

	obj, created, err := pkgCRDREST.Update(namespacedCtx(nonGlobalNamespace), name, UpdatePackageTestImpl{updateDisplayNameFn(newDisplayName, name)}, nil, nil, false, &metav1.UpdateOptions{})
	if err != nil {
		t.Fatalf("Expected update operation to succeed, got: %v", err)
	}

	if !created {
		t.Fatalf("Expected object to be created")
	}

	pkg, ok := obj.(*datapackaging.PackageMetadata)
	if !ok {
		t.Fatalf("Expected get operation to return PackageMetadata, but got: %v", reflect.TypeOf(pkg))
	}

	if pkg.Name != name || pkg.Namespace != nonGlobalNamespace {
		t.Fatalf("Expected returned package to have name %s and namespace %s, got %s and %s", name, nonGlobalNamespace, pkg.Name, pkg.Namespace)
	}

	if globalPackage.Spec.DisplayName != originalDisplayName {
		t.Fatalf("Expected update not to affect the global package, but DisplayName was updated")
	}
}

// deleting
// scoped to ns, so if cant find in ns, don't do anything
func TestPackageDeleteExistsInNS(t *testing.T) {
	namespacedPackage := namespacedIntPackage()
	name := namespacedPackage.Name

	internalClient := fake.NewSimpleClientset(namespacedPackage)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgCRDREST := datapkgreg.NewPackageMetadataCRDREST(internalClient, fakeCoreClient, globalNamespace)

	_, _, err := pkgCRDREST.Delete(namespacedCtx(nonGlobalNamespace), name, nil, &metav1.DeleteOptions{})
	if err != nil {
		t.Fatalf("Expected delete operation to succeed, got: %v", err)
	}

	actions := internalClient.Actions()

	if len(actions) == 0 {
		t.Fatalf("Internal Client was not used for deletion")
	}

	for _, action := range actions {
		if deleteAction, ok := action.(cgtesting.DeleteActionImpl); ok {
			if deleteAction.GetNamespace() == nonGlobalNamespace && deleteAction.GetName() == name {
				return
			}
			t.Fatalf("Unexpected delete action: %#v", deleteAction)
		}
	}

	t.Fatalf("Did not find delete action for namespace %s", nonGlobalNamespace)
}

func TestPackageDeleteExistsGlobalNotInNS(t *testing.T) {
	globalPackage := globalIntPackage()
	name := globalPackage.Name

	internalClient := fake.NewSimpleClientset(globalPackage)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgCRDREST := datapkgreg.NewPackageMetadataCRDREST(internalClient, fakeCoreClient, globalNamespace)

	_, _, err := pkgCRDREST.Delete(namespacedCtx(nonGlobalNamespace), name, nil, &metav1.DeleteOptions{})
	if !errors.IsNotFound(err) {
		t.Fatalf("Expected delete operation to return not found error, got: %v", err)
	}

	actions := internalClient.Actions()

	if len(actions) == 0 {
		t.Fatalf("Internal Client was not used for deletion")
	}

	// we shouldn't send any deletes because package is in global ns
	for _, action := range actions {
		if action.GetVerb() == "DELETE" {
			t.Fatalf("Did not expect any delete actions, but got %#v", action)
		}
	}
}

// Helpers
func namespacedCtx(ns string) context.Context {
	return request.WithNamespace(request.NewContext(), ns)
}

func assertPackageListUnorderedEquals(actual, expected []datapackaging.PackageMetadata, t *testing.T) {
	if len(actual) != len(expected) {
		t.Fatalf("arrays had different length:\n actual: \n%#v\n\n expected: \n%#v", actual, expected)
	}

	pkgIdentifiers := make(map[string]struct{})
	for _, pkg := range actual {
		pkgIdentifiers[pkg.Namespace+"/"+pkg.Name] = struct{}{}
	}

	for _, pkg := range expected {
		identifier := pkg.Namespace + "/" + pkg.Name
		if _, ok := pkgIdentifiers[identifier]; !ok {
			t.Fatalf("Expected actual to contain pakcage %s, but it didnt", identifier)
		}
	}
}

func globalIntPackage() *v1alpha1.InternalPackageMetadata {
	return &v1alpha1.InternalPackageMetadata{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: globalNamespace,
			Name:      "global-package.carvel.dev",
		},
		Spec: datapackaging.PackageMetadataSpec{
			DisplayName: "GLOBAL",
		},
	}
}

func namespacedIntPackage() *v1alpha1.InternalPackageMetadata {
	return &v1alpha1.InternalPackageMetadata{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: nonGlobalNamespace,
			Name:      "namespaced-package.carvel.dev",
		},
		Spec: datapackaging.PackageMetadataSpec{
			DisplayName: "NAMESPACED",
		},
	}
}

func overrideIntPackage() *v1alpha1.InternalPackageMetadata {
	return &v1alpha1.InternalPackageMetadata{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: nonGlobalNamespace,
			Name:      "global-package.carvel.dev",
		},
		Spec: datapackaging.PackageMetadataSpec{
			DisplayName: "OVERRIDE",
		},
	}
}

func excludedNonGlobalIntPackage() *v1alpha1.InternalPackageMetadata {
	return &v1alpha1.InternalPackageMetadata{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: excludedNonGlobalNamespace,
			Name:      "excluded-package",
		},
		Spec: datapackaging.PackageMetadataSpec{
			DisplayName: "EXCLUDED",
		},
	}
}

func namespace() *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        nonGlobalNamespace,
			Annotations: make(map[string]string),
		},
	}
}

func updateDisplayNameFn(newDisplayName, resourceName string) func(pkg *datapackaging.PackageMetadata) *datapackaging.PackageMetadata {
	return func(pkg *datapackaging.PackageMetadata) *datapackaging.PackageMetadata {
		pkg.Spec.DisplayName = newDisplayName
		if pkg.Name == "" {
			pkg.Name = resourceName
		}
		return pkg
	}
}

type UpdatePackageTestImpl struct {
	updateFn func(pkg *datapackaging.PackageMetadata) *datapackaging.PackageMetadata
}

func (UpdatePackageTestImpl) Preconditions() *metav1.Preconditions {
	return nil
}

func (u UpdatePackageTestImpl) UpdatedObject(ctx context.Context, oldObj runtime.Object) (runtime.Object, error) {
	pkg := oldObj.(*datapackaging.PackageMetadata)
	return u.updateFn(pkg), nil
}
