// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeInstalledPackages implements InstalledPackageInterface
type FakeInstalledPackages struct {
	Fake *FakeInstallV1alpha1
	ns   string
}

var installedpackagesResource = schema.GroupVersionResource{Group: "install.package.carvel.dev", Version: "v1alpha1", Resource: "installedpackages"}

var installedpackagesKind = schema.GroupVersionKind{Group: "install.package.carvel.dev", Version: "v1alpha1", Kind: "InstalledPackage"}

// Get takes name of the installedPackage, and returns the corresponding installedPackage object, and an error if there is any.
func (c *FakeInstalledPackages) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.InstalledPackage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(installedpackagesResource, c.ns, name), &v1alpha1.InstalledPackage{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.InstalledPackage), err
}

// List takes label and field selectors, and returns the list of InstalledPackages that match those selectors.
func (c *FakeInstalledPackages) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.InstalledPackageList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(installedpackagesResource, installedpackagesKind, c.ns, opts), &v1alpha1.InstalledPackageList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.InstalledPackageList{ListMeta: obj.(*v1alpha1.InstalledPackageList).ListMeta}
	for _, item := range obj.(*v1alpha1.InstalledPackageList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested installedPackages.
func (c *FakeInstalledPackages) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(installedpackagesResource, c.ns, opts))

}

// Create takes the representation of a installedPackage and creates it.  Returns the server's representation of the installedPackage, and an error, if there is any.
func (c *FakeInstalledPackages) Create(ctx context.Context, installedPackage *v1alpha1.InstalledPackage, opts v1.CreateOptions) (result *v1alpha1.InstalledPackage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(installedpackagesResource, c.ns, installedPackage), &v1alpha1.InstalledPackage{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.InstalledPackage), err
}

// Update takes the representation of a installedPackage and updates it. Returns the server's representation of the installedPackage, and an error, if there is any.
func (c *FakeInstalledPackages) Update(ctx context.Context, installedPackage *v1alpha1.InstalledPackage, opts v1.UpdateOptions) (result *v1alpha1.InstalledPackage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(installedpackagesResource, c.ns, installedPackage), &v1alpha1.InstalledPackage{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.InstalledPackage), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeInstalledPackages) UpdateStatus(ctx context.Context, installedPackage *v1alpha1.InstalledPackage, opts v1.UpdateOptions) (*v1alpha1.InstalledPackage, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(installedpackagesResource, "status", c.ns, installedPackage), &v1alpha1.InstalledPackage{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.InstalledPackage), err
}

// Delete takes name of the installedPackage and deletes it. Returns an error if one occurs.
func (c *FakeInstalledPackages) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(installedpackagesResource, c.ns, name), &v1alpha1.InstalledPackage{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeInstalledPackages) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(installedpackagesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.InstalledPackageList{})
	return err
}

// Patch applies the patch and returns the patched installedPackage.
func (c *FakeInstalledPackages) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.InstalledPackage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(installedpackagesResource, c.ns, name, pt, data, subresources...), &v1alpha1.InstalledPackage{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.InstalledPackage), err
}
