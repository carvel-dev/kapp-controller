// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"k8s.io/client-go/kubernetes"
)

// ValuesFactory abstracts the factories for fetching the values away from the template factory
type ValuesFactory struct {
	fetchFactory               fetch.Factory
	coreClient                 kubernetes.Interface
	additionaDownwardAPIValues AdditionalDownwardAPIValues
}

// NewValues returns a Values struct based on the app template source and the app context
func (vf ValuesFactory) NewValues(valuesFrom []v1alpha1.AppTemplateValuesSource, appContext AppContext) Values {
	return Values{ValuesFrom: valuesFrom, appContext: appContext, coreClient: vf.coreClient}
}
