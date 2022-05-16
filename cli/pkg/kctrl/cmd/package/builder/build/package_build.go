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
	RegistryURL string `json:"registryUrl"`
	RepoName    string `json:"repoName"`
	Tag         string `json:"tag"`
}

type Spec struct {
	Pkg         *v1alpha1.Package         `json:"package, omitempty"`
	PkgMetadata *v1alpha1.PackageMetadata `json:"packageMetadata, omitempty"`
	Vendir      *vendirconf.Config        `json:"vendir, omitempty"`
	Imgpkg      *Imgpkg                   `json:"imgpkg, omitempty"`
}

func (pkgBuilder PackageBuild) GetPackageMetadata() v1alpha1.PackageMetadata {
	//TODO we should start getting the data from pkgBuilder rather than create
	return *pkgBuilder.Spec.PkgMetadata
}

func (pkgBuilder PackageBuild) GetPackage() v1alpha1.Package {
	//TODO we should start getting the data from pkgBuilder rather than create
	return *pkgBuilder.Spec.Pkg
}
