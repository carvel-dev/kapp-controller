//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	versionsv1alpha1 "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *App) DeepCopyInto(out *App) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new App.
func (in *App) DeepCopy() *App {
	if in == nil {
		return nil
	}
	out := new(App)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *App) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppCluster) DeepCopyInto(out *AppCluster) {
	*out = *in
	if in.KubeconfigSecretRef != nil {
		in, out := &in.KubeconfigSecretRef, &out.KubeconfigSecretRef
		*out = new(AppClusterKubeconfigSecretRef)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppCluster.
func (in *AppCluster) DeepCopy() *AppCluster {
	if in == nil {
		return nil
	}
	out := new(AppCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppClusterKubeconfigSecretRef) DeepCopyInto(out *AppClusterKubeconfigSecretRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppClusterKubeconfigSecretRef.
func (in *AppClusterKubeconfigSecretRef) DeepCopy() *AppClusterKubeconfigSecretRef {
	if in == nil {
		return nil
	}
	out := new(AppClusterKubeconfigSecretRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppDeploy) DeepCopyInto(out *AppDeploy) {
	*out = *in
	if in.Kapp != nil {
		in, out := &in.Kapp, &out.Kapp
		*out = new(AppDeployKapp)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppDeploy.
func (in *AppDeploy) DeepCopy() *AppDeploy {
	if in == nil {
		return nil
	}
	out := new(AppDeploy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppDeployKapp) DeepCopyInto(out *AppDeployKapp) {
	*out = *in
	if in.MapNs != nil {
		in, out := &in.MapNs, &out.MapNs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RawOptions != nil {
		in, out := &in.RawOptions, &out.RawOptions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Inspect != nil {
		in, out := &in.Inspect, &out.Inspect
		*out = new(AppDeployKappInspect)
		(*in).DeepCopyInto(*out)
	}
	if in.Delete != nil {
		in, out := &in.Delete, &out.Delete
		*out = new(AppDeployKappDelete)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppDeployKapp.
func (in *AppDeployKapp) DeepCopy() *AppDeployKapp {
	if in == nil {
		return nil
	}
	out := new(AppDeployKapp)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppDeployKappDelete) DeepCopyInto(out *AppDeployKappDelete) {
	*out = *in
	if in.RawOptions != nil {
		in, out := &in.RawOptions, &out.RawOptions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppDeployKappDelete.
func (in *AppDeployKappDelete) DeepCopy() *AppDeployKappDelete {
	if in == nil {
		return nil
	}
	out := new(AppDeployKappDelete)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppDeployKappInspect) DeepCopyInto(out *AppDeployKappInspect) {
	*out = *in
	if in.RawOptions != nil {
		in, out := &in.RawOptions, &out.RawOptions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppDeployKappInspect.
func (in *AppDeployKappInspect) DeepCopy() *AppDeployKappInspect {
	if in == nil {
		return nil
	}
	out := new(AppDeployKappInspect)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppFetch) DeepCopyInto(out *AppFetch) {
	*out = *in
	if in.Inline != nil {
		in, out := &in.Inline, &out.Inline
		*out = new(AppFetchInline)
		(*in).DeepCopyInto(*out)
	}
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(AppFetchImage)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTP != nil {
		in, out := &in.HTTP, &out.HTTP
		*out = new(AppFetchHTTP)
		(*in).DeepCopyInto(*out)
	}
	if in.Git != nil {
		in, out := &in.Git, &out.Git
		*out = new(AppFetchGit)
		(*in).DeepCopyInto(*out)
	}
	if in.HelmChart != nil {
		in, out := &in.HelmChart, &out.HelmChart
		*out = new(AppFetchHelmChart)
		(*in).DeepCopyInto(*out)
	}
	if in.ImgpkgBundle != nil {
		in, out := &in.ImgpkgBundle, &out.ImgpkgBundle
		*out = new(AppFetchImgpkgBundle)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppFetch.
func (in *AppFetch) DeepCopy() *AppFetch {
	if in == nil {
		return nil
	}
	out := new(AppFetch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppFetchGit) DeepCopyInto(out *AppFetchGit) {
	*out = *in
	if in.RefSelection != nil {
		in, out := &in.RefSelection, &out.RefSelection
		*out = new(versionsv1alpha1.VersionSelection)
		(*in).DeepCopyInto(*out)
	}
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(AppFetchLocalRef)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppFetchGit.
func (in *AppFetchGit) DeepCopy() *AppFetchGit {
	if in == nil {
		return nil
	}
	out := new(AppFetchGit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppFetchHTTP) DeepCopyInto(out *AppFetchHTTP) {
	*out = *in
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(AppFetchLocalRef)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppFetchHTTP.
func (in *AppFetchHTTP) DeepCopy() *AppFetchHTTP {
	if in == nil {
		return nil
	}
	out := new(AppFetchHTTP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppFetchHelmChart) DeepCopyInto(out *AppFetchHelmChart) {
	*out = *in
	if in.Repository != nil {
		in, out := &in.Repository, &out.Repository
		*out = new(AppFetchHelmChartRepo)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppFetchHelmChart.
func (in *AppFetchHelmChart) DeepCopy() *AppFetchHelmChart {
	if in == nil {
		return nil
	}
	out := new(AppFetchHelmChart)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppFetchHelmChartRepo) DeepCopyInto(out *AppFetchHelmChartRepo) {
	*out = *in
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(AppFetchLocalRef)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppFetchHelmChartRepo.
func (in *AppFetchHelmChartRepo) DeepCopy() *AppFetchHelmChartRepo {
	if in == nil {
		return nil
	}
	out := new(AppFetchHelmChartRepo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppFetchImage) DeepCopyInto(out *AppFetchImage) {
	*out = *in
	if in.TagSelection != nil {
		in, out := &in.TagSelection, &out.TagSelection
		*out = new(versionsv1alpha1.VersionSelection)
		(*in).DeepCopyInto(*out)
	}
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(AppFetchLocalRef)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppFetchImage.
func (in *AppFetchImage) DeepCopy() *AppFetchImage {
	if in == nil {
		return nil
	}
	out := new(AppFetchImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppFetchImgpkgBundle) DeepCopyInto(out *AppFetchImgpkgBundle) {
	*out = *in
	if in.TagSelection != nil {
		in, out := &in.TagSelection, &out.TagSelection
		*out = new(versionsv1alpha1.VersionSelection)
		(*in).DeepCopyInto(*out)
	}
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(AppFetchLocalRef)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppFetchImgpkgBundle.
func (in *AppFetchImgpkgBundle) DeepCopy() *AppFetchImgpkgBundle {
	if in == nil {
		return nil
	}
	out := new(AppFetchImgpkgBundle)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppFetchInline) DeepCopyInto(out *AppFetchInline) {
	*out = *in
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.PathsFrom != nil {
		in, out := &in.PathsFrom, &out.PathsFrom
		*out = make([]AppFetchInlineSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppFetchInline.
func (in *AppFetchInline) DeepCopy() *AppFetchInline {
	if in == nil {
		return nil
	}
	out := new(AppFetchInline)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppFetchInlineSource) DeepCopyInto(out *AppFetchInlineSource) {
	*out = *in
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(AppFetchInlineSourceRef)
		**out = **in
	}
	if in.ConfigMapRef != nil {
		in, out := &in.ConfigMapRef, &out.ConfigMapRef
		*out = new(AppFetchInlineSourceRef)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppFetchInlineSource.
func (in *AppFetchInlineSource) DeepCopy() *AppFetchInlineSource {
	if in == nil {
		return nil
	}
	out := new(AppFetchInlineSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppFetchInlineSourceRef) DeepCopyInto(out *AppFetchInlineSourceRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppFetchInlineSourceRef.
func (in *AppFetchInlineSourceRef) DeepCopy() *AppFetchInlineSourceRef {
	if in == nil {
		return nil
	}
	out := new(AppFetchInlineSourceRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppFetchLocalRef) DeepCopyInto(out *AppFetchLocalRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppFetchLocalRef.
func (in *AppFetchLocalRef) DeepCopy() *AppFetchLocalRef {
	if in == nil {
		return nil
	}
	out := new(AppFetchLocalRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppList) DeepCopyInto(out *AppList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]App, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppList.
func (in *AppList) DeepCopy() *AppList {
	if in == nil {
		return nil
	}
	out := new(AppList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppSpec) DeepCopyInto(out *AppSpec) {
	*out = *in
	if in.Cluster != nil {
		in, out := &in.Cluster, &out.Cluster
		*out = new(AppCluster)
		(*in).DeepCopyInto(*out)
	}
	if in.Fetch != nil {
		in, out := &in.Fetch, &out.Fetch
		*out = make([]AppFetch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = make([]AppTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Deploy != nil {
		in, out := &in.Deploy, &out.Deploy
		*out = make([]AppDeploy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SyncPeriod != nil {
		in, out := &in.SyncPeriod, &out.SyncPeriod
		*out = new(v1.Duration)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppSpec.
func (in *AppSpec) DeepCopy() *AppSpec {
	if in == nil {
		return nil
	}
	out := new(AppSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppStatus) DeepCopyInto(out *AppStatus) {
	*out = *in
	if in.Fetch != nil {
		in, out := &in.Fetch, &out.Fetch
		*out = new(AppStatusFetch)
		(*in).DeepCopyInto(*out)
	}
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(AppStatusTemplate)
		(*in).DeepCopyInto(*out)
	}
	if in.Deploy != nil {
		in, out := &in.Deploy, &out.Deploy
		*out = new(AppStatusDeploy)
		(*in).DeepCopyInto(*out)
	}
	if in.Inspect != nil {
		in, out := &in.Inspect, &out.Inspect
		*out = new(AppStatusInspect)
		(*in).DeepCopyInto(*out)
	}
	in.GenericStatus.DeepCopyInto(&out.GenericStatus)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppStatus.
func (in *AppStatus) DeepCopy() *AppStatus {
	if in == nil {
		return nil
	}
	out := new(AppStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppStatusDeploy) DeepCopyInto(out *AppStatusDeploy) {
	*out = *in
	in.StartedAt.DeepCopyInto(&out.StartedAt)
	in.UpdatedAt.DeepCopyInto(&out.UpdatedAt)
	if in.KappDeployStatus != nil {
		in, out := &in.KappDeployStatus, &out.KappDeployStatus
		*out = new(KappDeployStatus)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppStatusDeploy.
func (in *AppStatusDeploy) DeepCopy() *AppStatusDeploy {
	if in == nil {
		return nil
	}
	out := new(AppStatusDeploy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppStatusFetch) DeepCopyInto(out *AppStatusFetch) {
	*out = *in
	in.StartedAt.DeepCopyInto(&out.StartedAt)
	in.UpdatedAt.DeepCopyInto(&out.UpdatedAt)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppStatusFetch.
func (in *AppStatusFetch) DeepCopy() *AppStatusFetch {
	if in == nil {
		return nil
	}
	out := new(AppStatusFetch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppStatusInspect) DeepCopyInto(out *AppStatusInspect) {
	*out = *in
	in.UpdatedAt.DeepCopyInto(&out.UpdatedAt)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppStatusInspect.
func (in *AppStatusInspect) DeepCopy() *AppStatusInspect {
	if in == nil {
		return nil
	}
	out := new(AppStatusInspect)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppStatusTemplate) DeepCopyInto(out *AppStatusTemplate) {
	*out = *in
	in.UpdatedAt.DeepCopyInto(&out.UpdatedAt)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppStatusTemplate.
func (in *AppStatusTemplate) DeepCopy() *AppStatusTemplate {
	if in == nil {
		return nil
	}
	out := new(AppStatusTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplate) DeepCopyInto(out *AppTemplate) {
	*out = *in
	if in.Ytt != nil {
		in, out := &in.Ytt, &out.Ytt
		*out = new(AppTemplateYtt)
		(*in).DeepCopyInto(*out)
	}
	if in.Kbld != nil {
		in, out := &in.Kbld, &out.Kbld
		*out = new(AppTemplateKbld)
		(*in).DeepCopyInto(*out)
	}
	if in.HelmTemplate != nil {
		in, out := &in.HelmTemplate, &out.HelmTemplate
		*out = new(AppTemplateHelmTemplate)
		(*in).DeepCopyInto(*out)
	}
	if in.Kustomize != nil {
		in, out := &in.Kustomize, &out.Kustomize
		*out = new(AppTemplateKustomize)
		**out = **in
	}
	if in.Jsonnet != nil {
		in, out := &in.Jsonnet, &out.Jsonnet
		*out = new(AppTemplateJsonnet)
		**out = **in
	}
	if in.Sops != nil {
		in, out := &in.Sops, &out.Sops
		*out = new(AppTemplateSops)
		(*in).DeepCopyInto(*out)
	}
	if in.Cue != nil {
		in, out := &in.Cue, &out.Cue
		*out = new(AppTemplateCue)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplate.
func (in *AppTemplate) DeepCopy() *AppTemplate {
	if in == nil {
		return nil
	}
	out := new(AppTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateCue) DeepCopyInto(out *AppTemplateCue) {
	*out = *in
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ValuesFrom != nil {
		in, out := &in.ValuesFrom, &out.ValuesFrom
		*out = make([]AppTemplateValuesSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateCue.
func (in *AppTemplateCue) DeepCopy() *AppTemplateCue {
	if in == nil {
		return nil
	}
	out := new(AppTemplateCue)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateHelmTemplate) DeepCopyInto(out *AppTemplateHelmTemplate) {
	*out = *in
	if in.ValuesFrom != nil {
		in, out := &in.ValuesFrom, &out.ValuesFrom
		*out = make([]AppTemplateValuesSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.KubernetesVersion != nil {
		in, out := &in.KubernetesVersion, &out.KubernetesVersion
		*out = new(Version)
		**out = **in
	}
	if in.KubernetesAPIs != nil {
		in, out := &in.KubernetesAPIs, &out.KubernetesAPIs
		*out = new(Version)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateHelmTemplate.
func (in *AppTemplateHelmTemplate) DeepCopy() *AppTemplateHelmTemplate {
	if in == nil {
		return nil
	}
	out := new(AppTemplateHelmTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateJsonnet) DeepCopyInto(out *AppTemplateJsonnet) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateJsonnet.
func (in *AppTemplateJsonnet) DeepCopy() *AppTemplateJsonnet {
	if in == nil {
		return nil
	}
	out := new(AppTemplateJsonnet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateKbld) DeepCopyInto(out *AppTemplateKbld) {
	*out = *in
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateKbld.
func (in *AppTemplateKbld) DeepCopy() *AppTemplateKbld {
	if in == nil {
		return nil
	}
	out := new(AppTemplateKbld)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateKustomize) DeepCopyInto(out *AppTemplateKustomize) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateKustomize.
func (in *AppTemplateKustomize) DeepCopy() *AppTemplateKustomize {
	if in == nil {
		return nil
	}
	out := new(AppTemplateKustomize)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateSops) DeepCopyInto(out *AppTemplateSops) {
	*out = *in
	if in.PGP != nil {
		in, out := &in.PGP, &out.PGP
		*out = new(AppTemplateSopsPGP)
		(*in).DeepCopyInto(*out)
	}
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Age != nil {
		in, out := &in.Age, &out.Age
		*out = new(AppTemplateSopsAge)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateSops.
func (in *AppTemplateSops) DeepCopy() *AppTemplateSops {
	if in == nil {
		return nil
	}
	out := new(AppTemplateSops)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateSopsAge) DeepCopyInto(out *AppTemplateSopsAge) {
	*out = *in
	if in.PrivateKeysSecretRef != nil {
		in, out := &in.PrivateKeysSecretRef, &out.PrivateKeysSecretRef
		*out = new(AppTemplateSopsPrivateKeysSecretRef)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateSopsAge.
func (in *AppTemplateSopsAge) DeepCopy() *AppTemplateSopsAge {
	if in == nil {
		return nil
	}
	out := new(AppTemplateSopsAge)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateSopsPGP) DeepCopyInto(out *AppTemplateSopsPGP) {
	*out = *in
	if in.PrivateKeysSecretRef != nil {
		in, out := &in.PrivateKeysSecretRef, &out.PrivateKeysSecretRef
		*out = new(AppTemplateSopsPrivateKeysSecretRef)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateSopsPGP.
func (in *AppTemplateSopsPGP) DeepCopy() *AppTemplateSopsPGP {
	if in == nil {
		return nil
	}
	out := new(AppTemplateSopsPGP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateSopsPrivateKeysSecretRef) DeepCopyInto(out *AppTemplateSopsPrivateKeysSecretRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateSopsPrivateKeysSecretRef.
func (in *AppTemplateSopsPrivateKeysSecretRef) DeepCopy() *AppTemplateSopsPrivateKeysSecretRef {
	if in == nil {
		return nil
	}
	out := new(AppTemplateSopsPrivateKeysSecretRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateValuesDownwardAPI) DeepCopyInto(out *AppTemplateValuesDownwardAPI) {
	*out = *in
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AppTemplateValuesDownwardAPIItem, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateValuesDownwardAPI.
func (in *AppTemplateValuesDownwardAPI) DeepCopy() *AppTemplateValuesDownwardAPI {
	if in == nil {
		return nil
	}
	out := new(AppTemplateValuesDownwardAPI)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateValuesDownwardAPIItem) DeepCopyInto(out *AppTemplateValuesDownwardAPIItem) {
	*out = *in
	if in.KubernetesVersion != nil {
		in, out := &in.KubernetesVersion, &out.KubernetesVersion
		*out = new(Version)
		**out = **in
	}
	if in.KappControllerVersion != nil {
		in, out := &in.KappControllerVersion, &out.KappControllerVersion
		*out = new(Version)
		**out = **in
	}
	if in.KubernetesAPIs != nil {
		in, out := &in.KubernetesAPIs, &out.KubernetesAPIs
		*out = new(GroupVersion)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateValuesDownwardAPIItem.
func (in *AppTemplateValuesDownwardAPIItem) DeepCopy() *AppTemplateValuesDownwardAPIItem {
	if in == nil {
		return nil
	}
	out := new(AppTemplateValuesDownwardAPIItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateValuesSource) DeepCopyInto(out *AppTemplateValuesSource) {
	*out = *in
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(AppTemplateValuesSourceRef)
		**out = **in
	}
	if in.ConfigMapRef != nil {
		in, out := &in.ConfigMapRef, &out.ConfigMapRef
		*out = new(AppTemplateValuesSourceRef)
		**out = **in
	}
	if in.DownwardAPI != nil {
		in, out := &in.DownwardAPI, &out.DownwardAPI
		*out = new(AppTemplateValuesDownwardAPI)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateValuesSource.
func (in *AppTemplateValuesSource) DeepCopy() *AppTemplateValuesSource {
	if in == nil {
		return nil
	}
	out := new(AppTemplateValuesSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateValuesSourceRef) DeepCopyInto(out *AppTemplateValuesSourceRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateValuesSourceRef.
func (in *AppTemplateValuesSourceRef) DeepCopy() *AppTemplateValuesSourceRef {
	if in == nil {
		return nil
	}
	out := new(AppTemplateValuesSourceRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppTemplateYtt) DeepCopyInto(out *AppTemplateYtt) {
	*out = *in
	if in.Inline != nil {
		in, out := &in.Inline, &out.Inline
		*out = new(AppFetchInline)
		(*in).DeepCopyInto(*out)
	}
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.FileMarks != nil {
		in, out := &in.FileMarks, &out.FileMarks
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ValuesFrom != nil {
		in, out := &in.ValuesFrom, &out.ValuesFrom
		*out = make([]AppTemplateValuesSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppTemplateYtt.
func (in *AppTemplateYtt) DeepCopy() *AppTemplateYtt {
	if in == nil {
		return nil
	}
	out := new(AppTemplateYtt)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AssociatedResources) DeepCopyInto(out *AssociatedResources) {
	*out = *in
	if in.Namespaces != nil {
		in, out := &in.Namespaces, &out.Namespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.GroupKinds != nil {
		in, out := &in.GroupKinds, &out.GroupKinds
		*out = make([]v1.GroupKind, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AssociatedResources.
func (in *AssociatedResources) DeepCopy() *AssociatedResources {
	if in == nil {
		return nil
	}
	out := new(AssociatedResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *Condition) DeepCopy() *Condition {
	if in == nil {
		return nil
	}
	out := new(Condition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GenericStatus) DeepCopyInto(out *GenericStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]Condition, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericStatus.
func (in *GenericStatus) DeepCopy() *GenericStatus {
	if in == nil {
		return nil
	}
	out := new(GenericStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GroupVersion) DeepCopyInto(out *GroupVersion) {
	*out = *in
	if in.GroupVersions != nil {
		in, out := &in.GroupVersions, &out.GroupVersions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GroupVersion.
func (in *GroupVersion) DeepCopy() *GroupVersion {
	if in == nil {
		return nil
	}
	out := new(GroupVersion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KappDeployStatus) DeepCopyInto(out *KappDeployStatus) {
	*out = *in
	in.AssociatedResources.DeepCopyInto(&out.AssociatedResources)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KappDeployStatus.
func (in *KappDeployStatus) DeepCopy() *KappDeployStatus {
	if in == nil {
		return nil
	}
	out := new(KappDeployStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Version) DeepCopyInto(out *Version) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Version.
func (in *Version) DeepCopy() *Version {
	if in == nil {
		return nil
	}
	out := new(Version)
	in.DeepCopyInto(out)
	return out
}
