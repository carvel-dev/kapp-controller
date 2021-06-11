// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
)

func addConversionFuncs(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Package"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.refName", "metadata.name", "metadata.namespace":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label %q not supported for Package", label)
			}
		})
}
