package app

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation"
)

const (
	deleteFinalizerName = "finalizers.kapp-ctrl.k14s.io/delete"
)

func init() {
	if errs := validation.IsQualifiedName(deleteFinalizerName); len(errs) > 0 {
		panic(fmt.Sprintf("Expected '%s' to be a valid finalizer name: %#v", deleteFinalizerName, errs))
	}
}
