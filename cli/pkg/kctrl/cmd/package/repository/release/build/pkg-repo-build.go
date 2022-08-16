package build

import (
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	PkgRepoBuildAPIVersion = "kctrl.carvel.dev/v1alpha1"
	PkgRepoBuildKind       = "PackageRepositoryBuild"
	PkgRepoBuildFileName   = "pkgrepo-build.yml"
)

type PackageRepoBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              PackageRepoBuildSpec `json:"spec,omitempty"`
}

type PackageRepoBuildSpec struct {
	Export *PackageRepoBuildExport `json:"export,omitempty"`
}

type PackageRepoBuildExport struct {
	ImgpkgBundle *PackageRepoBuildExportImgpkgBundle `json:"imgpkgBundle,omitempty"`
}

type PackageRepoBuildExportImgpkgBundle struct {
	Image string `json:"image,omitempty"`
}

func (pkgBuilder PackageRepoBuild) WriteToFile() error {
	content, err := yaml.Marshal(pkgBuilder)
	if err != nil {
		return err
	}
	fileLocation := filepath.Join(PkgRepoBuildFileName)
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
