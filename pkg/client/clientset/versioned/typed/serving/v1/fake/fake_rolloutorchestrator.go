/*
Copyright 2023 The Knative Authors

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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1 "knative.dev/serving-progressive-rollout/pkg/apis/serving/v1"
)

// FakeRolloutOrchestrators implements RolloutOrchestratorInterface
type FakeRolloutOrchestrators struct {
	Fake *FakeServingV1
	ns   string
}

var rolloutorchestratorsResource = v1.SchemeGroupVersion.WithResource("rolloutorchestrators")

var rolloutorchestratorsKind = v1.SchemeGroupVersion.WithKind("RolloutOrchestrator")

// Get takes name of the rolloutOrchestrator, and returns the corresponding rolloutOrchestrator object, and an error if there is any.
func (c *FakeRolloutOrchestrators) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.RolloutOrchestrator, err error) {
	emptyResult := &v1.RolloutOrchestrator{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(rolloutorchestratorsResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.RolloutOrchestrator), err
}

// List takes label and field selectors, and returns the list of RolloutOrchestrators that match those selectors.
func (c *FakeRolloutOrchestrators) List(ctx context.Context, opts metav1.ListOptions) (result *v1.RolloutOrchestratorList, err error) {
	emptyResult := &v1.RolloutOrchestratorList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(rolloutorchestratorsResource, rolloutorchestratorsKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.RolloutOrchestratorList{ListMeta: obj.(*v1.RolloutOrchestratorList).ListMeta}
	for _, item := range obj.(*v1.RolloutOrchestratorList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested rolloutOrchestrators.
func (c *FakeRolloutOrchestrators) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(rolloutorchestratorsResource, c.ns, opts))

}

// Create takes the representation of a rolloutOrchestrator and creates it.  Returns the server's representation of the rolloutOrchestrator, and an error, if there is any.
func (c *FakeRolloutOrchestrators) Create(ctx context.Context, rolloutOrchestrator *v1.RolloutOrchestrator, opts metav1.CreateOptions) (result *v1.RolloutOrchestrator, err error) {
	emptyResult := &v1.RolloutOrchestrator{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(rolloutorchestratorsResource, c.ns, rolloutOrchestrator, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.RolloutOrchestrator), err
}

// Update takes the representation of a rolloutOrchestrator and updates it. Returns the server's representation of the rolloutOrchestrator, and an error, if there is any.
func (c *FakeRolloutOrchestrators) Update(ctx context.Context, rolloutOrchestrator *v1.RolloutOrchestrator, opts metav1.UpdateOptions) (result *v1.RolloutOrchestrator, err error) {
	emptyResult := &v1.RolloutOrchestrator{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(rolloutorchestratorsResource, c.ns, rolloutOrchestrator, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.RolloutOrchestrator), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRolloutOrchestrators) UpdateStatus(ctx context.Context, rolloutOrchestrator *v1.RolloutOrchestrator, opts metav1.UpdateOptions) (result *v1.RolloutOrchestrator, err error) {
	emptyResult := &v1.RolloutOrchestrator{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(rolloutorchestratorsResource, "status", c.ns, rolloutOrchestrator, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.RolloutOrchestrator), err
}

// Delete takes name of the rolloutOrchestrator and deletes it. Returns an error if one occurs.
func (c *FakeRolloutOrchestrators) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(rolloutorchestratorsResource, c.ns, name, opts), &v1.RolloutOrchestrator{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRolloutOrchestrators) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(rolloutorchestratorsResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1.RolloutOrchestratorList{})
	return err
}

// Patch applies the patch and returns the patched rolloutOrchestrator.
func (c *FakeRolloutOrchestrators) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.RolloutOrchestrator, err error) {
	emptyResult := &v1.RolloutOrchestrator{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(rolloutorchestratorsResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.RolloutOrchestrator), err
}
