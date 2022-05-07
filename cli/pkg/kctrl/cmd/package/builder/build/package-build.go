package build

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PackageBuild struct {
	metav1.TypeMeta `json:",inline"`
	Spec            Spec `json:"spec, omitempty"`
}

type Imgpkg struct {
	metav1.TypeMeta  `json:",inline"`
	RegistryURL      string `json:"registryUrl"`
	RegistryUserName string `json:"registryUserName,omitempty"`
	RegistryPassword string `json:"registryPassword,omitempty"`
	RepoName         string `json:"repoName"`
	Tag              string `json:"tag"`
}

type Spec struct {
	Pkg         v1alpha1.Package         `json:"package"`
	PkgMetadata v1alpha1.PackageMetadata `json:"packageMetadata"`
	Vendir      vendirconf.Config        `json:"vendir"`
	Imgpkg      Imgpkg                   `json:"imgpkg"`
}

func (pkgBuilder PackageBuild) GetPackageMetadata() v1alpha1.PackageMetadata {
	//TODO we should start getting the data from pkgBuilder rather than create
	return pkgBuilder.Spec.PkgMetadata
}

func (pkgBuilder PackageBuild) GetPackage() v1alpha1.Package {
	//TODO we should start getting the data from pkgBuilder rather than create
	return pkgBuilder.Spec.Pkg
}
