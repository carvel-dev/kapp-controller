package fetch

import (
	"bytes"
	"fmt"
	"io"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kyaml "sigs.k8s.io/yaml"
)

const vendirEntireDirPath = "."

type Vendir struct {
	nsName     string
	coreClient kubernetes.Interface
	config     vendirconf.Config
}

func NewVendir(nsName string, coreClient kubernetes.Interface) *Vendir {
	return &Vendir{
		nsName:     nsName,
		coreClient: coreClient,
		config: vendirconf.Config{
			APIVersion: "vendir.k14s.io/v1alpha1", // TODO: use constant from vendir package
			Kind:       "Config",                  // TODO: use constant from vendir package
		}}
}

func (v *Vendir) AddDir(fetch v1alpha1.AppFetch, dirPath string) error {
	switch {
	case fetch.Inline != nil:
		v.config.Directories = append(v.config.Directories, v.dir(v.inlineConf(*fetch.Inline), dirPath))
	case fetch.Image != nil:
		v.config.Directories = append(v.config.Directories, v.dir(v.imageConf(*fetch.Image), dirPath))
	case fetch.HTTP != nil:
		v.config.Directories = append(v.config.Directories, v.dir(v.httpConf(*fetch.HTTP), dirPath))
	case fetch.Git != nil:
		v.config.Directories = append(v.config.Directories, v.dir(v.gitConf(*fetch.Git), dirPath))
	case fetch.HelmChart != nil:
		v.config.Directories = append(v.config.Directories, v.dir(v.helmChartConf(*fetch.HelmChart), dirPath))
	case fetch.ImgpkgBundle != nil:
		v.config.Directories = append(v.config.Directories, v.dir(v.imgpkgBundleConf(*fetch.ImgpkgBundle), dirPath))
	default:
		return fmt.Errorf("Unsupported way to fetch templates")
	}

	return nil
}

func (v *Vendir) dir(contents vendirconf.DirectoryContents, dirPath string) vendirconf.Directory {
	return vendirconf.Directory{
		Path:     dirPath,
		Contents: []vendirconf.DirectoryContents{contents},
	}
}

func (v *Vendir) inlineConf(inline v1alpha1.AppFetchInline) vendirconf.DirectoryContents {
	var inlineSources []vendirconf.DirectoryContentsInlineSource
	for _, source := range inline.PathsFrom {
		inlineSources = append(inlineSources, v.inlineSourceConf(source))
	}
	return vendirconf.DirectoryContents{
		Path: vendirEntireDirPath,
		Inline: &vendirconf.DirectoryContentsInline{
			Paths:     inline.Paths,
			PathsFrom: inlineSources,
		}}
}

func (v *Vendir) imageConf(image v1alpha1.AppFetchImage) vendirconf.DirectoryContents {
	return vendirconf.DirectoryContents{
		Path:        vendirEntireDirPath,
		NewRootPath: image.SubPath,
		Image: &vendirconf.DirectoryContentsImage{
			URL:       image.URL,
			SecretRef: v.localRefConf(image.SecretRef),
		},
	}
}

func (v *Vendir) imgpkgBundleConf(imgpkgBundle v1alpha1.AppFetchImgpkgBundle) vendirconf.DirectoryContents {
	return vendirconf.DirectoryContents{
		Path:        vendirEntireDirPath,
		NewRootPath: imgpkgBundle.SubPath,
		ImgpkgBundle: &vendirconf.DirectoryContentsImgpkgBundle{
			Image:     imgpkgBundle.Image,
			SecretRef: v.localRefConf(imgpkgBundle.SecretRef),
		},
	}
}

func (v *Vendir) httpConf(http v1alpha1.AppFetchHTTP) vendirconf.DirectoryContents {
	return vendirconf.DirectoryContents{
		Path: vendirEntireDirPath,
		HTTP: &vendirconf.DirectoryContentsHTTP{
			URL:       http.URL,
			SHA256:    http.SHA256,
			SecretRef: v.localRefConf(http.SecretRef),
		},
		NewRootPath: http.SubPath,
	}
}

func (v *Vendir) gitConf(git v1alpha1.AppFetchGit) vendirconf.DirectoryContents {
	return vendirconf.DirectoryContents{
		Path:        vendirEntireDirPath,
		NewRootPath: git.SubPath,
		Git: &vendirconf.DirectoryContentsGit{
			URL:           git.URL,
			Ref:           git.Ref,
			SecretRef:     v.localRefConf(git.SecretRef),
			LFSSkipSmudge: git.LFSSkipSmudge,
		},
	}
}

func (v *Vendir) helmChartConf(chart v1alpha1.AppFetchHelmChart) vendirconf.DirectoryContents {
	return vendirconf.DirectoryContents{
		Path: vendirEntireDirPath,
		HelmChart: &vendirconf.DirectoryContentsHelmChart{
			Name:       chart.Name,
			Version:    chart.Version,
			Repository: v.helmRepoConf(chart.Repository),
		},
	}
}

func (v *Vendir) inlineSourceConf(src v1alpha1.AppFetchInlineSource) vendirconf.DirectoryContentsInlineSource {
	return vendirconf.DirectoryContentsInlineSource{
		SecretRef:    v.inlineSourceRefConf(src.SecretRef),
		ConfigMapRef: v.inlineSourceRefConf(src.ConfigMapRef),
	}
}

func (v *Vendir) inlineSourceRefConf(ref *v1alpha1.AppFetchInlineSourceRef) *vendirconf.DirectoryContentsInlineSourceRef {
	if ref == nil {
		return nil
	}

	return &vendirconf.DirectoryContentsInlineSourceRef{
		DirectoryPath:             ref.DirectoryPath,
		DirectoryContentsLocalRef: vendirconf.DirectoryContentsLocalRef{Name: ref.LocalObjectReference.Name},
	}
}

func (v *Vendir) helmRepoConf(repo *v1alpha1.AppFetchHelmChartRepo) *vendirconf.DirectoryContentsHelmChartRepo {
	if repo == nil {
		return nil
	}

	return &vendirconf.DirectoryContentsHelmChartRepo{
		URL:       repo.URL,
		SecretRef: v.localRefConf(repo.SecretRef),
	}
}

func (v *Vendir) localRefConf(ref *v1alpha1.AppFetchLocalRef) *vendirconf.DirectoryContentsLocalRef {
	if ref == nil {
		return nil
	}

	return &vendirconf.DirectoryContentsLocalRef{
		Name: ref.LocalObjectReference.Name,
	}
}

func (v *Vendir) ConfigReader() (io.Reader, error) {
	var resourcesYaml [][]byte
	for _, dir := range v.config.Directories {
		for _, contents := range dir.Contents {
			yamlBytes, err := v.requiredResourcesYaml(contents)
			if err != nil {
				return nil, err
			}

			resourcesYaml = append(resourcesYaml, yamlBytes...)
		}
	}

	vendirConfBytes, err := v.config.AsBytes()
	if err != nil {
		return nil, err
	}

	finalConfig := bytes.Join(append(resourcesYaml, vendirConfBytes), []byte("---\n"))

	return bytes.NewReader(finalConfig), nil
}

func (v *Vendir) requiredResourcesYaml(contents vendirconf.DirectoryContents) ([][]byte, error) {
	switch {
	case contents.Inline != nil:
		return v.inlineResources(*contents.Inline)
	case contents.Image != nil:
		return v.imageResources(*contents.Image)
	case contents.HTTP != nil:
		return v.httpResources(*contents.HTTP)
	case contents.Git != nil:
		return v.gitResources(*contents.Git)
	case contents.HelmChart != nil:
		return v.helmChartResources(*contents.HelmChart)
	case contents.ImgpkgBundle != nil:
		return v.imgpkgBundleResources(*contents.ImgpkgBundle)
	}

	return nil, fmt.Errorf("Unknown fetch type: %v", contents)
}

func (v *Vendir) inlineResources(inline vendirconf.DirectoryContentsInline) ([][]byte, error) {
	var resourcesYamlBytes [][]byte
	for _, source := range inline.PathsFrom {
		switch {
		case source.SecretRef != nil:
			bytes, err := v.secretBytes(source.SecretRef.DirectoryContentsLocalRef)
			if err != nil {
				return nil, err
			}

			resourcesYamlBytes = append(resourcesYamlBytes, bytes)

		case source.ConfigMapRef != nil:
			bytes, err := v.configMapBytes(source.ConfigMapRef.DirectoryContentsLocalRef)
			if err != nil {
				return nil, err
			}

			resourcesYamlBytes = append(resourcesYamlBytes, bytes)
		}
	}

	return resourcesYamlBytes, nil
}

func (v *Vendir) imageResources(image vendirconf.DirectoryContentsImage) ([][]byte, error) {
	if image.SecretRef == nil {
		return nil, nil
	}

	resBytes, err := v.secretBytes(*image.SecretRef)
	if err != nil {
		return nil, err
	}

	return [][]byte{resBytes}, nil
}

func (v *Vendir) imgpkgBundleResources(imgpkgBundle vendirconf.DirectoryContentsImgpkgBundle) ([][]byte, error) {
	if imgpkgBundle.SecretRef == nil {
		return nil, nil
	}

	resBytes, err := v.secretBytes(*imgpkgBundle.SecretRef)
	if err != nil {
		return nil, err
	}

	return [][]byte{resBytes}, nil
}

func (v *Vendir) httpResources(http vendirconf.DirectoryContentsHTTP) ([][]byte, error) {
	if http.SecretRef == nil {
		return nil, nil
	}

	resBytes, err := v.secretBytes(*http.SecretRef)
	if err != nil {
		return nil, err
	}

	return [][]byte{resBytes}, nil
}

func (v *Vendir) gitResources(git vendirconf.DirectoryContentsGit) ([][]byte, error) {
	if git.SecretRef == nil {
		return nil, nil
	}

	resBytes, err := v.secretBytes(*git.SecretRef)
	if err != nil {
		return nil, err
	}

	return [][]byte{resBytes}, nil
}

func (v *Vendir) helmChartResources(helmChart vendirconf.DirectoryContentsHelmChart) ([][]byte, error) {
	if helmChart.Repository == nil || helmChart.Repository.SecretRef == nil {
		return nil, nil
	}

	resBytes, err := v.secretBytes(*helmChart.Repository.SecretRef)
	if err != nil {
		return nil, err
	}

	return [][]byte{resBytes}, nil
}

func (v *Vendir) secretBytes(secretRef vendirconf.DirectoryContentsLocalRef) ([]byte, error) {
	secret, err := v.coreClient.CoreV1().Secrets(v.nsName).Get(secretRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// typed clients drop GVK or resource (https://github.com/kubernetes/kubernetes/issues/80609)
	secret.TypeMeta.Kind = "Secret"
	secret.TypeMeta.APIVersion = "v1"

	return kyaml.Marshal(secret)
}

func (v *Vendir) configMapBytes(configMapRef vendirconf.DirectoryContentsLocalRef) ([]byte, error) {
	configMap, err := v.coreClient.CoreV1().ConfigMaps(v.nsName).Get(configMapRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// typed clients drop GVK or resource (https://github.com/kubernetes/kubernetes/issues/80609)
	configMap.TypeMeta.Kind = "ConfigMap"
	configMap.TypeMeta.APIVersion = "v1"

	return kyaml.Marshal(configMap)
}
