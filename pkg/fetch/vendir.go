package fetch

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
)

const vendirEntireDirPath = "."

type Vendir struct{}

func NewVendir() *Vendir {
	return &Vendir{}
}

func (v *Vendir) InlineDirConf(fetchOpts v1alpha1.AppFetchInline, dirPath string) vendirconf.Directory {
	return v.dir(v.inlineConf(fetchOpts), dirPath)
}

func (v *Vendir) ImageDirConf(fetchOpts v1alpha1.AppFetchImage, dirPath string) vendirconf.Directory {
	return v.dir(v.imageConf(fetchOpts), dirPath)
}

func (v *Vendir) HTTPDirConf(fetchOpts v1alpha1.AppFetchHTTP, dirPath string) vendirconf.Directory {
	return v.dir(v.httpConf(fetchOpts), dirPath)
}

func (v *Vendir) GitDirConf(fetchOpts v1alpha1.AppFetchGit, dirPath string) vendirconf.Directory {
	return v.dir(v.gitConf(fetchOpts), dirPath)
}

func (v *Vendir) HelmChartDirConf(fetchOpts v1alpha1.AppFetchHelmChart, dirPath string) vendirconf.Directory {
	return v.dir(v.helmChartConf(fetchOpts), dirPath)
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
