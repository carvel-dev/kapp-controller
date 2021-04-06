// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/package/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePackages implements PackageInterface
type FakePackages struct {
	Fake *FakePackageV1alpha1
}

var packagesResource = schema.GroupVersionResource{Group: "package.carvel.dev", Version: "v1alpha1", Resource: "packages"}

var packagesKind = schema.GroupVersionKind{Group: "package.carvel.dev", Version: "v1alpha1", Kind: "Package"}

// Get takes name of the package, and returns the corresponding package object, and an error if there is any.
func (c *FakePackages) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Package, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(packagesResource, name), &v1alpha1.Package{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Package), err
}

// List takes label and field selectors, and returns the list of Packages that match those selectors.
func (c *FakePackages) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.PackageList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(packagesResource, packagesKind, opts), &v1alpha1.PackageList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.PackageList{ListMeta: obj.(*v1alpha1.PackageList).ListMeta}
	for _, item := range obj.(*v1alpha1.PackageList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested packages.
func (c *FakePackages) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(packagesResource, opts))
}

// Create takes the representation of a package and creates it.  Returns the server's representation of the package, and an error, if there is any.
func (c *FakePackages) Create(ctx context.Context, pkg *v1alpha1.Package, opts v1.CreateOptions) (result *v1alpha1.Package, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(packagesResource, pkg), &v1alpha1.Package{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Package), err
}

// Update takes the representation of a package and updates it. Returns the server's representation of the package, and an error, if there is any.
func (c *FakePackages) Update(ctx context.Context, pkg *v1alpha1.Package, opts v1.UpdateOptions) (result *v1alpha1.Package, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(packagesResource, pkg), &v1alpha1.Package{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Package), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePackages) UpdateStatus(ctx context.Context, pkg *v1alpha1.Package, opts v1.UpdateOptions) (*v1alpha1.Package, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(packagesResource, "status", pkg), &v1alpha1.Package{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Package), err
}

// Delete takes name of the package and deletes it. Returns an error if one occurs.
func (c *FakePackages) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(packagesResource, name), &v1alpha1.Package{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePackages) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(packagesResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.PackageList{})
	return err
}

// Patch applies the patch and returns the patched package.
func (c *FakePackages) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Package, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(packagesResource, name, pt, data, subresources...), &v1alpha1.Package{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Package), err
}
