package pkgrepository

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/scheme"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	appNs = "kapp-controller"
)

func NewApp(existingApp *v1alpha1.App, pkgRepository *v1alpha1.PkgRepository) (*v1alpha1.App, error) {
	desiredApp := existingApp.DeepCopy()

	desiredApp.Name = pkgRepository.Name
	desiredApp.Namespace = appNs

	err := controllerutil.SetControllerReference(pkgRepository, desiredApp, scheme.Scheme)
	if err != nil {
		return &v1alpha1.App{}, err
	}

	desiredApp.Spec = v1alpha1.AppSpec{
		// TODO since we are assuming that we are inside kapp-controller NS, use its SA
		ServiceAccountName: "kapp-controller-sa",
		Fetch: []v1alpha1.AppFetch{{
			Image:        pkgRepository.Spec.Fetch.Image,
			Git:          pkgRepository.Spec.Fetch.Git,
			HTTP:         pkgRepository.Spec.Fetch.HTTP,
			ImgpkgBundle: pkgRepository.Spec.Fetch.Bundle,
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

	if desiredApp.Spec.Fetch[0].ImgpkgBundle != nil {
		desiredApp.Spec.Template = append(desiredApp.Spec.Template,
			v1alpha1.AppTemplate{Kbld: &v1alpha1.AppTemplateKbld{Paths: []string{"-", ".imgpkg/images.yml"}}})
	}

	return desiredApp, nil
}
