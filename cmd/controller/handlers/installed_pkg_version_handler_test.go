// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package handlers_test

import (
	"reflect"
	"testing"

	"github.com/go-logr/logr"

	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/vmware-tanzu/carvel-kapp-controller/cmd/controller/handlers"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestOnlyEligiblePackagesAreEnqueued(t *testing.T) {
	q := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	eligibleInstalledPkg := pkgingv1alpha1.InstalledPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg",
		},
		Spec: pkgingv1alpha1.InstalledPackageSpec{
			PackageVersionRef: &pkgingv1alpha1.PackageVersionRef{
				PackageName: "expec-pkg",
				VersionSelection: &v1alpha1.VersionSelectionSemver{
					Constraints: ">=1.0.0",
				},
			},
		},
	}

	ineligibleInstalledPkg := pkgingv1alpha1.InstalledPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg-ineligible",
		},
		Spec: pkgingv1alpha1.InstalledPackageSpec{
			PackageVersionRef: &pkgingv1alpha1.PackageVersionRef{
				PackageName: "expec-pkg",
				VersionSelection: &v1alpha1.VersionSelectionSemver{
					Constraints: "<1.0.0",
				},
			},
		},
	}

	// Load installed package into fake client
	kappcs := fake.NewSimpleClientset(&eligibleInstalledPkg, &ineligibleInstalledPkg)
	ipvh := handlers.NewInstalledPkgVersionHandler(kappcs, &EmptyLog{})

	event := event.GenericEvent{
		Object: &datapkgingv1alpha1.PackageVersion{
			Spec: datapkgingv1alpha1.PackageVersionSpec{
				PackageName: "expec-pkg",
				Version:     "1.5.0",
			},
		},
	}

	ipvh.Generic(event, q)

	if q.Len() != 1 {
		t.Fatalf("Expected queue to have length of 1, got %d", q.Len())
	}

	expectedRequest := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "", Name: "expected-pkg"}}
	if obj, _ := q.Get(); !reflect.DeepEqual(obj, expectedRequest) {
		t.Fatalf("Expected queue to contain the installed package eligible for upgrade, but contained:\n\n%#v\n", obj)
	}

}

type EmptyLog struct{}

func (l *EmptyLog) Info(msg string, keysAndValues ...interface{}) {}

func (l *EmptyLog) Enabled() bool { return false }

func (l *EmptyLog) Error(err error, msg string, keysAndValues ...interface{}) {}

func (l *EmptyLog) V(level int) logr.InfoLogger { return l }

func (l *EmptyLog) WithValues(keysAndValues ...interface{}) logr.Logger { return l }

func (l *EmptyLog) WithName(name string) logr.Logger { return l }
