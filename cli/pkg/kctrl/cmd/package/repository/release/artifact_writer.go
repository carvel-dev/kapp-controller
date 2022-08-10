package release

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/yaml"
)

type ArtifactWriter struct {
	PackageRepoName string
	TargetDir       string
}

func NewArtifactWriter(pkgRepoName string, directory string) *ArtifactWriter {
	return &ArtifactWriter{PackageRepoName: pkgRepoName, TargetDir: directory}
}

func (w *ArtifactWriter) WritePackageRepositoryFile(imgpkgBundleLocation string) error {
	packageRepository := v1alpha1.PackageRepository{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PackageRepository",
			APIVersion: "packaging.carvel.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              w.PackageRepoName,
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
		Spec: v1alpha1.PackageRepositorySpec{
			Fetch: &v1alpha1.PackageRepositoryFetch{
				ImgpkgBundle: &v1alpha12.AppFetchImgpkgBundle{
					Image: imgpkgBundleLocation,
				},
			},
		},
	}

	packageRepoBytes, err := yaml.Marshal(packageRepository)
	if err != nil {
		return err
	}
	path := filepath.Join(w.TargetDir, PkgRepositoryFileName)

	return w.createOrOverwriteFile(path, packageRepoBytes)
}

func (w *ArtifactWriter) createOrOverwriteFile(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
