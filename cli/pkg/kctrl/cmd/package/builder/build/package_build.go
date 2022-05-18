package build

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	PkgBuildFileName = "package-build.yml"
)

type PackageBuild struct {
	metav1.TypeMeta `json:",inline"`
	Spec            Spec `json:"spec, omitempty"`
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
	//TODO we should start getting the data from pkgBuilder rather than create
	return *pkgBuilder.Spec.PkgMetadata
}

func (pkgBuilder PackageBuild) GetPackage() v1alpha1.Package {
	//TODO we should start getting the data from pkgBuilder rather than create
	return *pkgBuilder.Spec.Pkg
}

func (pkgBuilder PackageBuild) WriteToFile(dirPath string) error {
	content, err := yaml.Marshal(pkgBuilder)
	if err != nil {
		return err
	}
	fileLocation := filepath.Join(dirPath, PkgBuildFileName)
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
	// check if error is "file not exists"
	if os.IsNotExist(err) {
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
	err = json.Unmarshal(content, &packageBuild)
	if err != nil {
		return PackageBuild{}, err
	}
	return packageBuild, nil
}
