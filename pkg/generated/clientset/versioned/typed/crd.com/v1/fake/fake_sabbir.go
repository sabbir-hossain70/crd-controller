/*
Copyright Sabbir Hossain.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "github.com/sabbir-hossain70/crd/pkg/apis/crd.com/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSabbirs implements SabbirInterface
type FakeSabbirs struct {
	Fake *FakeCrdV1
	ns   string
}

var sabbirsResource = v1.SchemeGroupVersion.WithResource("sabbirs")

var sabbirsKind = v1.SchemeGroupVersion.WithKind("Sabbir")

// Get takes name of the sabbir, and returns the corresponding sabbir object, and an error if there is any.
func (c *FakeSabbirs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Sabbir, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(sabbirsResource, c.ns, name), &v1.Sabbir{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Sabbir), err
}

// List takes label and field selectors, and returns the list of Sabbirs that match those selectors.
func (c *FakeSabbirs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.SabbirList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(sabbirsResource, sabbirsKind, c.ns, opts), &v1.SabbirList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.SabbirList{ListMeta: obj.(*v1.SabbirList).ListMeta}
	for _, item := range obj.(*v1.SabbirList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested sabbirs.
func (c *FakeSabbirs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(sabbirsResource, c.ns, opts))

}

// Create takes the representation of a sabbir and creates it.  Returns the server's representation of the sabbir, and an error, if there is any.
func (c *FakeSabbirs) Create(ctx context.Context, sabbir *v1.Sabbir, opts metav1.CreateOptions) (result *v1.Sabbir, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(sabbirsResource, c.ns, sabbir), &v1.Sabbir{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Sabbir), err
}

// Update takes the representation of a sabbir and updates it. Returns the server's representation of the sabbir, and an error, if there is any.
func (c *FakeSabbirs) Update(ctx context.Context, sabbir *v1.Sabbir, opts metav1.UpdateOptions) (result *v1.Sabbir, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(sabbirsResource, c.ns, sabbir), &v1.Sabbir{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Sabbir), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSabbirs) UpdateStatus(ctx context.Context, sabbir *v1.Sabbir, opts metav1.UpdateOptions) (*v1.Sabbir, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(sabbirsResource, "status", c.ns, sabbir), &v1.Sabbir{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Sabbir), err
}

// Delete takes name of the sabbir and deletes it. Returns an error if one occurs.
func (c *FakeSabbirs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(sabbirsResource, c.ns, name, opts), &v1.Sabbir{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSabbirs) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(sabbirsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1.SabbirList{})
	return err
}

// Patch applies the patch and returns the patched sabbir.
func (c *FakeSabbirs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Sabbir, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(sabbirsResource, c.ns, name, pt, data, subresources...), &v1.Sabbir{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Sabbir), err
}
