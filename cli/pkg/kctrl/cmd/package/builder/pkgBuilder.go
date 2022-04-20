package builder

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PkgBuilder struct {
	metav1.TypeMeta `json:",inline"`
	Config          Config `json:"config, omitempty"`
}

type Config struct {
	Pkg         v1alpha1.Package         `json:"package"`
	PkgMetadata v1alpha1.PackageMetadata `json:"packageMetadata"`
	//Pending Add imgpkgConfiguration
}
