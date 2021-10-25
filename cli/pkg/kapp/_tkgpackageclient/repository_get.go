// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"fmt"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

func (p *pkgClient) GetRepository(o *tkgpackagedatamodel.RepositoryOptions) (*kappipkg.PackageRepository, error) {
	packageRepository, err := p.kappClient.GetPackageRepository(o.RepositoryName, o.Namespace)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Warningf("package repository '%s' does not exist in namespace '%s'", o.RepositoryName, o.Namespace)
			return nil, nil
		}
		return nil, errors.Wrap(err, fmt.Sprintf("failed to find package repository '%s' in namespace '%s'", o.RepositoryName, o.Namespace))
	}

	return packageRepository, nil
}
