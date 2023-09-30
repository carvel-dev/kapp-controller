// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	pkginstall "github.com/vmware-tanzu/carvel-kapp-controller/pkg/packageinstall"
)

func TestOnlyEligiblePackagesAreEnqueued(t *testing.T) {
	q := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	eligibleInstalledPkg := pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg",
		},
		Spec: pkgingv1alpha1.PackageInstallSpec{
			PackageRef: &pkgingv1alpha1.PackageRef{
				RefName: "expec-pkg",
				VersionSelection: &v1alpha1.VersionSelectionSemver{
					Constraints: ">=1.0.0",
				},
			},
		},
	}

	ineligibleInstalledPkg := pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg-ineligible",
		},
		Spec: pkgingv1alpha1.PackageInstallSpec{
			PackageRef: &pkgingv1alpha1.PackageRef{
				RefName: "expec-pkg",
				VersionSelection: &v1alpha1.VersionSelectionSemver{
					Constraints: "<1.0.0",
				},
			},
		},
	}

	// Load installed package into fake client
	kappcs := fake.NewSimpleClientset(&eligibleInstalledPkg, &ineligibleInstalledPkg)
	ipvh := pkginstall.NewPackageInstallVersionHandler(kappcs, "", testr.New(t))

	event := event.GenericEvent{
		Object: &datapkgingv1alpha1.Package{
			Spec: datapkgingv1alpha1.PackageSpec{
				RefName: "expec-pkg",
				Version: "1.5.0",
			},
		},
	}

	ipvh.Generic(context.Background(), event, q)

	if q.Len() != 1 {
		t.Fatalf("Expected queue to have length of 1, got %d", q.Len())
	}

	expectedRequest := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "", Name: "expected-pkg"}}
	if obj, _ := q.Get(); !reflect.DeepEqual(obj, expectedRequest) {
		t.Fatalf("Expected queue to contain the installed package eligible for upgrade, but contained:\n\n%#v\n", obj)
	}

}
