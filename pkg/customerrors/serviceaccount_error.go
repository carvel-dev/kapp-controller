// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package customerrors

const (
	ServiceAccountNotFound = 1
)

type ServiceAccountError struct {
	errorMsg string
	kind     int
}

func NewServiceAccountError(errorMsg string, kind int) *ServiceAccountError {
	return &ServiceAccountError{errorMsg, kind}
}

func (s ServiceAccountError) Error() string {
	return s.errorMsg
}

func (s ServiceAccountError) Kind() int {
	return s.kind
}
