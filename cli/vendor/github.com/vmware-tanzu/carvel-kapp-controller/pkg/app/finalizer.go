// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

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

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
