package build

import (
	"fmt"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/common"
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type PackageBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:",inline"`
	Spec              Spec `json:"spec, omitempty"`
}

type Imgpkg struct {
	RegistryURL string `json:"registryUrl,omitempty"`
}

type Spec struct {
	Pkg         *v1alpha1.Package         `json:"package, omitempty"`
	PkgMetadata *v1alpha1.PackageMetadata `json:"packageMetadata, omitempty"`
	Vendir      *vendirconf.Config        `json:"vendir, omitempty"`
	Imgpkg      *Imgpkg                   `json:"imgpkg, omitempty"`
}

func (pkgBuilder PackageBuild) GetPackageMetadata() v1alpha1.PackageMetadata {
	return *pkgBuilder.Spec.PkgMetadata
}

func (pkgBuilder PackageBuild) GetPackage() v1alpha1.Package {
	return *pkgBuilder.Spec.Pkg
}

func (pkgBuilder PackageBuild) WriteToFile() error {
	content, err := yaml.Marshal(pkgBuilder)
	if err != nil {
		return err
	}
	fileLocation := filepath.Join(common.PkgBuildFileName)
	file, err := os.Create(fileLocation)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		return err
	}
	return nil
}

func GeneratePackageBuild(pkgBuildFilePath string) (PackageBuild, error) {
	var pkgBuild PackageBuild
	exists, err := isPkgBuildFileExists(pkgBuildFilePath)
	if err != nil {
		return PackageBuild{}, err
	}
	if exists {
		pkgBuild, err = NewPackageBuildFromFile(pkgBuildFilePath)
		if err != nil {
			return PackageBuild{}, err
		}

		//In case user has manually removed the pkg section from the package-build
		if pkgBuild.Spec.Pkg == nil {
			defaultPkg, err := NewDefaultPackage()
			if err != nil {
				return PackageBuild{}, err
			}
			pkgBuild.Spec.Pkg = defaultPkg
		}

		//In case user has manually removed the pkg-metadata section from the package-build
		if pkgBuild.Spec.PkgMetadata == nil {
			defaultPkgMetadata, err := NewDefaultPackageMetadata()
			if err != nil {
				return PackageBuild{}, err
			}
			pkgBuild.Spec.PkgMetadata = defaultPkgMetadata
		}

		//In case user has manually removed the vendir config section from the package-build
		if pkgBuild.Spec.Vendir == nil {
			defaultVendirConfig, err := NewDefaultVendir()
			if err != nil {
				return PackageBuild{}, err
			}
			pkgBuild.Spec.Vendir = defaultVendirConfig
		}
	} else {
		pkgBuild, err = NewDefaultPackageBuild()
		if err != nil {
			return PackageBuild{}, err
		}
	}
	return pkgBuild, nil
}

// Check if file exists
func isPkgBuildFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to check for the existence of file. Error is: %s", err.Error())
	}
}

func NewPackageBuildFromFile(filePath string) (PackageBuild, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return PackageBuild{}, err
	}
	var packageBuild PackageBuild
	err = yaml.Unmarshal(content, &packageBuild)
	if err != nil {
		return PackageBuild{}, err
	}
	return packageBuild, nil
}
