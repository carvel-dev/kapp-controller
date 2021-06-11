// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package datapackaging_test

import (
	"context"
	"encoding/base32"
	"reflect"
	"strings"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/internalpackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging"
	datapkgreg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/registry/datapackaging"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	cgtesting "k8s.io/client-go/testing"
)

const globalNamespace = "global.packaging.kapp-controller.carvel.dev"
const nonGlobalNamespace = "some.other.ns.carvel.dev"
const excludedNonGlobalNamespace = "excluded-from-lists.ns.carvel.dev"

// listing
func TestPackageVersionListIncludesGlobalAndNamespaced(t *testing.T) {
	internalClient := fake.NewSimpleClientset(globalIntPackageVersion(), namespacedIntPackageVersion(), excludedNonGlobalIntPackageVersion())
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgvCRDREST := datapkgreg.NewPackageCRDREST(internalClient, fakeCoreClient, globalNamespace)

	pkgvList, err := pkgvCRDREST.List(namespacedCtx(nonGlobalNamespace), &internalversion.ListOptions{})
	if err != nil {
		t.Fatalf("Expected list operation to succeed, got: %v", err)
	}

	packageVersionList, ok := pkgvList.(*datapackaging.PackageList)
	if !ok {
		t.Fatalf("Expected list operation to return PackageList, but got: %v", reflect.TypeOf(pkgvList))
	}

	expectedpkgvs := []datapackaging.Package{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "some.other.ns.carvel.dev",
				Name:      "global-package-version.carvel.dev.1.0.0",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "some.other.ns.carvel.dev",
				Name:      "namespaced-package-version.carvel.dev.1.0.0",
			},
		},
	}

	assertPVListUnorderedEquals(packageVersionList.Items, expectedpkgvs, t)
}

func TestPackageVersionListPrefersNamespacedOverGlobal(t *testing.T) {
	// seed client with resources
	internalClient := fake.NewSimpleClientset(globalIntPackageVersion(), overrideIntPackageVersion())
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgvCRDREST := datapkgreg.NewPackageCRDREST(internalClient, fakeCoreClient, globalNamespace)

	// list package versions and verify all of them are there
	pkgvList, err := pkgvCRDREST.List(namespacedCtx(nonGlobalNamespace), &internalversion.ListOptions{})
	if err != nil {
		t.Fatalf("Expected list operation to succeed, got: %v", err)
	}

	packageVersionList, ok := pkgvList.(*datapackaging.PackageList)
	if !ok {
		t.Fatalf("Expected list operation to return PackageList, but got: %v", reflect.TypeOf(pkgvList))
	}

	expectedpkgvs := []datapackaging.Package{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: nonGlobalNamespace,
				Name:      "global-package-version.carvel.dev.mismatch",
			},
		},
	}

	assertPVListUnorderedEquals(packageVersionList.Items, expectedpkgvs, t)
}

// getting

func TestPackageVersionGetNotPresentInNS(t *testing.T) {
	globalPackageVersion := globalIntPackageVersion()
	name := globalPackageName
	releaseNotes := globalPackageVersion.Spec.ReleaseNotes

	internalClient := fake.NewSimpleClientset(globalPackageVersion)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgvCRDREST := datapkgreg.NewPackageCRDREST(internalClient, fakeCoreClient, globalNamespace)

	obj, err := pkgvCRDREST.Get(namespacedCtx(nonGlobalNamespace), name, &metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Expected get operation to succeed, got: %v", err)
	}

	pkgv, ok := obj.(*datapackaging.Package)
	if !ok {
		t.Fatalf("Expected get operation to return Package, but got: %v", reflect.TypeOf(pkgv))
	}

	if pkgv.Name != name || pkgv.Spec.ReleaseNotes != releaseNotes {
		t.Fatalf("Expected returned package to have name %s and namespace %s, got %s and %s", name, globalNamespace, pkgv.Name, pkgv.Namespace)
	}
}

func TestPackageVersionGetPresentInOnlyNS(t *testing.T) {
	namespacedPackageVersion := namespacedIntPackageVersion()
	name := namespacedPackageName
	releaseNotes := namespacedPackageVersion.Spec.ReleaseNotes

	internalClient := fake.NewSimpleClientset(namespacedPackageVersion)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgvCRDREST := datapkgreg.NewPackageCRDREST(internalClient, fakeCoreClient, globalNamespace)

	obj, err := pkgvCRDREST.Get(namespacedCtx(nonGlobalNamespace), name, &metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Expected get operation to succeed, got: %v", err)
	}

	pkgv, ok := obj.(*datapackaging.Package)
	if !ok {
		t.Fatalf("Expected get operation to return Package, but got: %v", reflect.TypeOf(pkgv))
	}

	if pkgv.Name != name || pkgv.Spec.ReleaseNotes != releaseNotes {
		t.Fatalf("Expected returned package to have name %s and namespace %s, got %s and %s", name, nonGlobalNamespace, pkgv.Name, pkgv.Namespace)
	}
}

func TestPackageVersionGetNotFound(t *testing.T) {
	namespacedPackageVersion := excludedNonGlobalIntPackageVersion()
	name := namespacedPackageName

	internalClient := fake.NewSimpleClientset(namespacedPackageVersion)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgvCRDREST := datapkgreg.NewPackageCRDREST(internalClient, fakeCoreClient, globalNamespace)

	_, err := pkgvCRDREST.Get(namespacedCtx(nonGlobalNamespace), name, &metav1.GetOptions{})
	if err == nil {
		t.Fatalf("Expected get operation to fail, but it didn't")
	}

	if !errors.IsNotFound(err) {
		t.Fatalf("Expected a not found error, got: %v", err)
	}

}

func TestPackageVersionGetPreferNS(t *testing.T) {
	overridePackageVersion := overrideIntPackageVersion()
	name := overridePackageName
	releaseNotes := overridePackageVersion.Spec.ReleaseNotes

	internalClient := fake.NewSimpleClientset(overridePackageVersion, globalIntPackageVersion())
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgvCRDREST := datapkgreg.NewPackageCRDREST(internalClient, fakeCoreClient, globalNamespace)

	obj, err := pkgvCRDREST.Get(namespacedCtx(nonGlobalNamespace), name, &metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Expected get operation to succeed, got: %v", err)
	}

	pkgv, ok := obj.(*datapackaging.Package)
	if !ok {
		t.Fatalf("Expected get operation to return Package, but got: %v", reflect.TypeOf(pkgv))
	}

	if pkgv.Name != name || pkgv.Spec.ReleaseNotes != releaseNotes {
		t.Fatalf("Expected returned package to have name %s and namespace %s, got %s and %s", name, nonGlobalNamespace, pkgv.Name, pkgv.Namespace)
	}
}

// updating
func TestPackageVersionUpdateDoesntUpdateGlobal(t *testing.T) {
	globalPackageVersion := globalIntPackageVersion()
	namespacedPackageVersion := namespacedIntPackageVersion()
	name := globalPackageName
	packageName := globalPackageVersion.Spec.RefName
	version := globalPackageVersion.Spec.Version
	originalReleaseNotes := globalPackageVersion.Spec.ReleaseNotes
	newReleaseNotes := "im-new"

	internalClient := fake.NewSimpleClientset(globalPackageVersion, namespacedPackageVersion)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgvCRDREST := datapkgreg.NewPackageCRDREST(internalClient, fakeCoreClient, globalNamespace)

	obj, created, err := pkgvCRDREST.Update(namespacedCtx(nonGlobalNamespace), name, UpdatePackageVersionTestImpl{updateReleaseNotesFn(newReleaseNotes, name, packageName, version)}, nil, nil, false, &metav1.UpdateOptions{})
	if err != nil {
		t.Fatalf("Expected update operation to succeed, got: %v", err)
	}

	if !created {
		t.Fatalf("Expected object to be created")
	}

	pkgv, ok := obj.(*datapackaging.Package)
	if !ok {
		t.Fatalf("Expected get operation to return Package, but got: %v", reflect.TypeOf(pkgv))
	}

	if pkgv.Name != name || pkgv.Namespace != nonGlobalNamespace {
		t.Fatalf("Expected returned package to have name %s and namespace %s, got %s and %s", name, nonGlobalNamespace, pkgv.Name, pkgv.Namespace)
	}

	if pkgv.Spec.ReleaseNotes != newReleaseNotes {
		t.Fatalf("Expected release notes of namespaced pacakge to be updated")
	}

	if globalPackageVersion.Spec.ReleaseNotes != originalReleaseNotes {
		t.Fatalf("Expected update not to affect the global package, but ReleaseNotes was updated")
	}
}

// scoped to ns, so if can't find create in ns
func TestPackageVersionUpdateCreatesInNS(t *testing.T) {
	globalPackageVersion := globalIntPackageVersion()
	name := globalPackageName
	packageName := globalPackageVersion.Spec.RefName
	version := globalPackageVersion.Spec.Version
	originalReleaseNotes := globalPackageVersion.Spec.ReleaseNotes
	newReleaseNotes := "im-new"

	internalClient := fake.NewSimpleClientset(globalPackageVersion)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgvCRDREST := datapkgreg.NewPackageCRDREST(internalClient, fakeCoreClient, globalNamespace)

	obj, created, err := pkgvCRDREST.Update(namespacedCtx(nonGlobalNamespace), name, UpdatePackageVersionTestImpl{updateReleaseNotesFn(newReleaseNotes, name, packageName, version)}, nil, nil, false, &metav1.UpdateOptions{})
	if err != nil {
		t.Fatalf("Expected update operation to succeed, got: %v", err)
	}

	if !created {
		t.Fatalf("Expected object to be created")
	}

	pkgv, ok := obj.(*datapackaging.Package)
	if !ok {
		t.Fatalf("Expected get operation to return Package, but got: %v", reflect.TypeOf(pkgv))
	}

	if pkgv.Name != name || pkgv.Namespace != nonGlobalNamespace {
		t.Fatalf("Expected returned package to have name %s and namespace %s, got %s and %s", name, nonGlobalNamespace, pkgv.Name, pkgv.Namespace)
	}

	if globalPackageVersion.Spec.ReleaseNotes != originalReleaseNotes {
		t.Fatalf("Expected update not to affect the global package, but ReleaseNotes was updated")
	}
}

// deleting
// scoped to ns, so if cant find in ns, don't do anything
func TestPackageVersionDeleteExistsInNS(t *testing.T) {
	namespacedPackageVersion := namespacedIntPackageVersion()
	name := namespacedPackageName

	internalClient := fake.NewSimpleClientset(namespacedPackageVersion)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgvCRDREST := datapkgreg.NewPackageCRDREST(internalClient, fakeCoreClient, globalNamespace)

	_, _, err := pkgvCRDREST.Delete(namespacedCtx(nonGlobalNamespace), name, nil, &metav1.DeleteOptions{})
	if err != nil {
		t.Fatalf("Expected delete operation to succeed, got: %v", err)
	}

	actions := internalClient.Actions()

	if len(actions) == 0 {
		t.Fatalf("Internal Client was not used for deletion")
	}

	for _, action := range actions {
		if deleteAction, ok := action.(cgtesting.DeleteActionImpl); ok {
			if deleteAction.GetNamespace() == nonGlobalNamespace && deleteAction.GetName() == encodeName(name) {
				return
			}
			t.Fatalf("Unexpected delete action: %#v", deleteAction)
		}
	}

	t.Fatalf("Did not find delete action for namespace %s", nonGlobalNamespace)
}

func TestPackageVersionDeleteExistsGlobalNotInNS(t *testing.T) {
	globalPackageVersion := globalIntPackageVersion()
	name := globalPackageName

	internalClient := fake.NewSimpleClientset(globalPackageVersion)
	fakeCoreClient := k8sfake.NewSimpleClientset(namespace())

	pkgvCRDREST := datapkgreg.NewPackageCRDREST(internalClient, fakeCoreClient, globalNamespace)

	_, _, err := pkgvCRDREST.Delete(namespacedCtx(nonGlobalNamespace), name, nil, &metav1.DeleteOptions{})
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
func assertPVListUnorderedEquals(actual, expected []datapackaging.Package, t *testing.T) {
	if len(actual) != len(expected) {
		t.Fatalf("arrays had different length:\n actual: \n%#v\n\n expected: \n%#v", actual, expected)
	}

	pkgvIdentifiers := make(map[string]struct{})
	for _, pkgv := range actual {
		pkgvIdentifiers[pkgv.Namespace+"/"+pkgv.Name] = struct{}{}
	}

	for _, pkgv := range expected {
		identifier := pkgv.Namespace + "/" + pkgv.Name
		if _, ok := pkgvIdentifiers[identifier]; !ok {
			t.Fatalf("Expected actual to contain pakcage %s, but it didnt", identifier)
		}
	}
}

const (
	globalPackageName     = "global-package-version.carvel.dev.1.0.0"
	namespacedPackageName = "namespaced-package-version.carvel.dev.1.0.0"
	overridePackageName   = "global-package-version.carvel.dev.mismatch"
)

func globalIntPackageVersion() *v1alpha1.InternalPackage {
	return &v1alpha1.InternalPackage{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: globalNamespace,
			Name:      encodeName("global-package-version.carvel.dev.1.0.0"),
		},
		Spec: datapackaging.PackageSpec{
			Version:      "1.0.0",
			RefName:      "global-package-version.carvel.dev",
			ReleaseNotes: "GLOBAL",
		},
	}
}

func namespacedIntPackageVersion() *v1alpha1.InternalPackage {
	return &v1alpha1.InternalPackage{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: nonGlobalNamespace,
			Name:      encodeName("namespaced-package-version.carvel.dev.1.0.0"),
		},
		Spec: datapackaging.PackageSpec{
			Version:      "1.0.0",
			RefName:      "namespaced-package-version.carvel.dev",
			ReleaseNotes: "NAMESPACED",
		},
	}
}

// Override is determined by packageMetadataName and version instead of just name
func overrideIntPackageVersion() *v1alpha1.InternalPackage {
	return &v1alpha1.InternalPackage{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: nonGlobalNamespace,
			Name:      encodeName("global-package-version.carvel.dev.mismatch"),
		},
		Spec: datapackaging.PackageSpec{
			Version:      "1.0.0",
			RefName:      "global-package-version.carvel.dev",
			ReleaseNotes: "OVERRIDE",
		},
	}
}

func excludedNonGlobalIntPackageVersion() *v1alpha1.InternalPackage {
	return &v1alpha1.InternalPackage{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: excludedNonGlobalNamespace,
			Name:      encodeName("excluded-package-version.carvel.dev.1.0.0"),
		},
		Spec: datapackaging.PackageSpec{
			Version:      "1.0.0",
			RefName:      "excluded-package-version.carvel.dev",
			ReleaseNotes: "EXCLUDED",
		},
	}
}

func encodeName(name string) string {
	return strings.ToLower(base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(name)))
}

func updateReleaseNotesFn(newNote, resourceName, packageName, version string) func(pkgv *datapackaging.Package) *datapackaging.Package {
	return func(pkgv *datapackaging.Package) *datapackaging.Package {
		pkgv.Spec.ReleaseNotes = newNote
		if pkgv.Name == "" {
			pkgv.Name = resourceName
		}
		if pkgv.Spec.RefName == "" {
			pkgv.Spec.RefName = packageName
		}
		if pkgv.Spec.Version == "" {
			pkgv.Spec.Version = version
		}
		return pkgv
	}
}

type UpdatePackageVersionTestImpl struct {
	updateFn func(pkgv *datapackaging.Package) *datapackaging.Package
}

func (UpdatePackageVersionTestImpl) Preconditions() *metav1.Preconditions {
	return nil
}

func (u UpdatePackageVersionTestImpl) UpdatedObject(ctx context.Context, oldObj runtime.Object) (runtime.Object, error) {
	pkgv := oldObj.(*datapackaging.Package)
	return u.updateFn(pkgv), nil
}
