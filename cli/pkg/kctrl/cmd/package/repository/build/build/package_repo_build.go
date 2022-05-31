package build

import (
	"fmt"
	"os"
	"path/filepath"

	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	PackageRepositoryBuildFileName   = "pkgrepo-build.yml"
	PackageRepositoryKind            = "PackageRepository"
	PackageRepositoryAPIVersion      = "packaging.carvel.dev/v1alpha1"
	PackageRepositoryBuildKind       = "PackageRepositoryBuild"
	PackageRepositoryBuildAPIVersion = "kctrl.carvel.dev/v1alpha1"
)

type PackageRepositoryBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Spec `json:"spec, omitempty"`
}

type Spec struct {
	PkgRepo *v1alpha12.PackageRepository `json:"packageRepository, omitempty"`
}

func (pkgBuilder PackageRepositoryBuild) WriteToFile(dirPath string) error {
	content, err := yaml.Marshal(pkgBuilder)
	if err != nil {
		return err
	}
	fileLocation := filepath.Join(dirPath, PackageRepositoryBuildFileName)
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

func GeneratePackageRepositoryBuild(pkgRepoBuildFilePath string) (PackageRepositoryBuild, error) {
	var pkgRepoBuild PackageRepositoryBuild
	exists, err := isPackageRepositoryBuildFileExists(pkgRepoBuildFilePath)
	if err != nil {
		return PackageRepositoryBuild{}, err
	}
	if exists {
		pkgRepoBuild, err = NewPackageRepositoryBuildFromFile(pkgRepoBuildFilePath)
		if err != nil {
			return PackageRepositoryBuild{}, err
		}
	} else {
		pkgRepoBuild = PackageRepositoryBuild{
			TypeMeta: metav1.TypeMeta{
				Kind:       PackageRepositoryBuildKind,
				APIVersion: PackageRepositoryBuildAPIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: Spec{
				PkgRepo: &v1alpha12.PackageRepository{
					TypeMeta: metav1.TypeMeta{
						Kind:       PackageRepositoryKind,
						APIVersion: PackageRepositoryAPIVersion,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "samplepackagerepository",
					},
				},
			},
		}
	}
	return pkgRepoBuild, nil
}

// Check if file exists
func isPackageRepositoryBuildFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to check for the existence of file. Error is: %s", err.Error())
	}
}

func NewPackageRepositoryBuildFromFile(filePath string) (PackageRepositoryBuild, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return PackageRepositoryBuild{}, err
	}
	var pkgRepoBuild PackageRepositoryBuild
	err = yaml.Unmarshal(content, &pkgRepoBuild)
	if err != nil {
		return PackageRepositoryBuild{}, err
	}
	return pkgRepoBuild, nil
}
