//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	versionsv1alpha1 "carvel.dev/vendir/pkg/vendir/versions/v1alpha1"
	kappctrlv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInstall) DeepCopyInto(out *PackageInstall) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInstall.
func (in *PackageInstall) DeepCopy() *PackageInstall {
	if in == nil {
		return nil
	}
	out := new(PackageInstall)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageInstall) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInstallList) DeepCopyInto(out *PackageInstallList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PackageInstall, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInstallList.
func (in *PackageInstallList) DeepCopy() *PackageInstallList {
	if in == nil {
		return nil
	}
	out := new(PackageInstallList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageInstallList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInstallSpec) DeepCopyInto(out *PackageInstallSpec) {
	*out = *in
	if in.Cluster != nil {
		in, out := &in.Cluster, &out.Cluster
		*out = new(kappctrlv1alpha1.AppCluster)
		(*in).DeepCopyInto(*out)
	}
	if in.PackageRef != nil {
		in, out := &in.PackageRef, &out.PackageRef
		*out = new(PackageRef)
		(*in).DeepCopyInto(*out)
	}
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]PackageInstallValues, len(*in))
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

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInstallSpec.
func (in *PackageInstallSpec) DeepCopy() *PackageInstallSpec {
	if in == nil {
		return nil
	}
	out := new(PackageInstallSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInstallStatus) DeepCopyInto(out *PackageInstallStatus) {
	*out = *in
	in.GenericStatus.DeepCopyInto(&out.GenericStatus)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInstallStatus.
func (in *PackageInstallStatus) DeepCopy() *PackageInstallStatus {
	if in == nil {
		return nil
	}
	out := new(PackageInstallStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInstallValues) DeepCopyInto(out *PackageInstallValues) {
	*out = *in
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(PackageInstallValuesSecretRef)
		**out = **in
	}
	if in.TemplateSteps != nil {
		in, out := &in.TemplateSteps, &out.TemplateSteps
		*out = make([]int, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInstallValues.
func (in *PackageInstallValues) DeepCopy() *PackageInstallValues {
	if in == nil {
		return nil
	}
	out := new(PackageInstallValues)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInstallValuesSecretRef) DeepCopyInto(out *PackageInstallValuesSecretRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInstallValuesSecretRef.
func (in *PackageInstallValuesSecretRef) DeepCopy() *PackageInstallValuesSecretRef {
	if in == nil {
		return nil
	}
	out := new(PackageInstallValuesSecretRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRef) DeepCopyInto(out *PackageRef) {
	*out = *in
	if in.VersionSelection != nil {
		in, out := &in.VersionSelection, &out.VersionSelection
		*out = new(versionsv1alpha1.VersionSelectionSemver)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRef.
func (in *PackageRef) DeepCopy() *PackageRef {
	if in == nil {
		return nil
	}
	out := new(PackageRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepository) DeepCopyInto(out *PackageRepository) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepository.
func (in *PackageRepository) DeepCopy() *PackageRepository {
	if in == nil {
		return nil
	}
	out := new(PackageRepository)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageRepository) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepositoryFetch) DeepCopyInto(out *PackageRepositoryFetch) {
	*out = *in
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(kappctrlv1alpha1.AppFetchImage)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTP != nil {
		in, out := &in.HTTP, &out.HTTP
		*out = new(kappctrlv1alpha1.AppFetchHTTP)
		(*in).DeepCopyInto(*out)
	}
	if in.Git != nil {
		in, out := &in.Git, &out.Git
		*out = new(kappctrlv1alpha1.AppFetchGit)
		(*in).DeepCopyInto(*out)
	}
	if in.ImgpkgBundle != nil {
		in, out := &in.ImgpkgBundle, &out.ImgpkgBundle
		*out = new(kappctrlv1alpha1.AppFetchImgpkgBundle)
		(*in).DeepCopyInto(*out)
	}
	if in.Inline != nil {
		in, out := &in.Inline, &out.Inline
		*out = new(kappctrlv1alpha1.AppFetchInline)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepositoryFetch.
func (in *PackageRepositoryFetch) DeepCopy() *PackageRepositoryFetch {
	if in == nil {
		return nil
	}
	out := new(PackageRepositoryFetch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepositoryList) DeepCopyInto(out *PackageRepositoryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PackageRepository, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepositoryList.
func (in *PackageRepositoryList) DeepCopy() *PackageRepositoryList {
	if in == nil {
		return nil
	}
	out := new(PackageRepositoryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageRepositoryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepositorySpec) DeepCopyInto(out *PackageRepositorySpec) {
	*out = *in
	if in.SyncPeriod != nil {
		in, out := &in.SyncPeriod, &out.SyncPeriod
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Fetch != nil {
		in, out := &in.Fetch, &out.Fetch
		*out = new(PackageRepositoryFetch)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepositorySpec.
func (in *PackageRepositorySpec) DeepCopy() *PackageRepositorySpec {
	if in == nil {
		return nil
	}
	out := new(PackageRepositorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepositoryStatus) DeepCopyInto(out *PackageRepositoryStatus) {
	*out = *in
	if in.Fetch != nil {
		in, out := &in.Fetch, &out.Fetch
		*out = new(kappctrlv1alpha1.AppStatusFetch)
		(*in).DeepCopyInto(*out)
	}
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(kappctrlv1alpha1.AppStatusTemplate)
		(*in).DeepCopyInto(*out)
	}
	if in.Deploy != nil {
		in, out := &in.Deploy, &out.Deploy
		*out = new(kappctrlv1alpha1.AppStatusDeploy)
		(*in).DeepCopyInto(*out)
	}
	in.GenericStatus.DeepCopyInto(&out.GenericStatus)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepositoryStatus.
func (in *PackageRepositoryStatus) DeepCopy() *PackageRepositoryStatus {
	if in == nil {
		return nil
	}
	out := new(PackageRepositoryStatus)
	in.DeepCopyInto(out)
	return out
}
