package pkgrepository

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

const (
	appNs = "kapp-controller"
)

func NewApp(existingApp *v1alpha1.App, pkgRepository *v1alpha1.PkgRepository) (*v1alpha1.App, error) {
	desiredApp := existingApp.DeepCopy()

	desiredApp.Name = pkgRepository.Name
	desiredApp.Namespace = appNs

	desiredApp.Spec = v1alpha1.AppSpec{
		// TODO since we are assuming that we are inside kapp-controller NS, use its SA
		ServiceAccountName: "kapp-controller-sa",
		Fetch: []v1alpha1.AppFetch{{
			Image: pkgRepository.Spec.Fetch.Image,
			Git:   pkgRepository.Spec.Fetch.Git,
			HTTP:  pkgRepository.Spec.Fetch.HTTP,
		}},
		Template: []v1alpha1.AppTemplate{{
			Ytt: &v1alpha1.AppTemplateYtt{
				IgnoreUnknownComments: true,
			},
		}},
		Deploy: []v1alpha1.AppDeploy{{
			Kapp: &v1alpha1.AppDeployKapp{},
		}},
	}

	return desiredApp, nil
}
