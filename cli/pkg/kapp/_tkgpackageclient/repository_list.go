// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"github.com/pkg/errors"

	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

func (p *pkgClient) ListRepositories(o *tkgpackagedatamodel.RepositoryOptions) (*kappipkg.PackageRepositoryList, error) {
	packageRepositoryList, err := p.kappClient.ListPackageRepositories(o.Namespace)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list existing package repositories in the cluster")
	}

	return packageRepositoryList, nil
}
