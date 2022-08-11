// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reftracker

import "fmt"

type RefKey struct {
	kind, refName, namespace string
}

func NewRefKey(kind, refName, namespace string) RefKey {
	return RefKey{kind, refName, namespace}
}

func NewSecretKey(refName, namespace string) RefKey {
	return NewRefKey("secret", refName, namespace)
}

func NewConfigMapKey(refName, namespace string) RefKey {
	return NewRefKey("configmap", refName, namespace)
}

func NewAppKey(refName, namespace string) RefKey {
	return NewRefKey("app", refName, namespace)
}

func NewPackageRepositoryKey(refName, namespace string) RefKey {
	return NewRefKey("packagerepository", refName, namespace)
}

func (r RefKey) Kind() string {
	return r.kind
}

func (r RefKey) RefName() string {
	return r.refName
}

func (r RefKey) Namespace() string {
	return r.namespace
}

func (r RefKey) Description() string {
	return fmt.Sprintf("%s:%s:%s", r.kind, r.refName, r.namespace)
}
