package build

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"sigs.k8s.io/yaml"
)

const defaultPackageBuildYAML = `
--- 
apiVersion: kctrl.carvel.dev/v1alpha1
kind: PackageBuild
spec: 
  imgpkg: 
    registryUrl: ~
  package: 
    apiVersion: data.packaging.carvel.dev/v1alpha1
    kind: Package
    metadata: 
      name: samplepackage.corp.com.1.0.0
      namespace: default
    spec: 
      licenses: 
        - "Apache 2.0"
      refName: samplepackage.corp.com
      releaseNotes: "Initial release"
      template: 
        spec: 
          deploy: 
            - kapp: {}
          fetch: 
            - imgpkgBundle: 
                image: ~
          template:
            - kbld: 
                paths: 
                  - "-"
                  - .imgpkg/images.yml
      version: "1.0.0"
  packageMetadata: 
    apiVersion: data.packaging.carvel.dev/v1alpha1
    kind: PackageMetadata
    metadata:
      name: samplepackage.corp.com
      namespace: default
    spec: 
      categories: 
        - demo
      displayName: ""
      longDescription: ""
      shortDescription: ""
  vendir: 
    apiVersion: vendir.k14s.io/v1alpha1
    kind: Config
    minimumRequiredVersion: "0.12.0"
`

func NewDefaultPackageBuild() (PackageBuild, error) {
	var packageBuild PackageBuild
	err := yaml.Unmarshal([]byte(defaultPackageBuildYAML), &packageBuild)
	if err != nil {
		return PackageBuild{}, err
	}
	return packageBuild, nil
}

func NewDefaultPackage() (*v1alpha1.Package, error) {
	pkgBuild, err := NewDefaultPackageBuild()
	if err != nil {
		return &v1alpha1.Package{}, err
	}
	return pkgBuild.Spec.Pkg, nil
}

func NewDefaultPackageMetadata() (*v1alpha1.PackageMetadata, error) {
	pkgBuild, err := NewDefaultPackageBuild()
	if err != nil {
		return &v1alpha1.PackageMetadata{}, err
	}
	return pkgBuild.Spec.PkgMetadata, nil
}

func NewDefaultVendir() (*vendirconf.Config, error) {
	pkgBuild, err := NewDefaultPackageBuild()
	if err != nil {
		return &vendirconf.Config{}, err
	}
	return pkgBuild.Spec.Vendir, nil
}
