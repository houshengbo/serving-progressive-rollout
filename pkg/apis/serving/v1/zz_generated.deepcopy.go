//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceOrchestrator) DeepCopyInto(out *ServiceOrchestrator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceOrchestrator.
func (in *ServiceOrchestrator) DeepCopy() *ServiceOrchestrator {
	if in == nil {
		return nil
	}
	out := new(ServiceOrchestrator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceOrchestrator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceOrchestratorList) DeepCopyInto(out *ServiceOrchestratorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceOrchestrator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceOrchestratorList.
func (in *ServiceOrchestratorList) DeepCopy() *ServiceOrchestratorList {
	if in == nil {
		return nil
	}
	out := new(ServiceOrchestratorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceOrchestratorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceOrchestratorSpec) DeepCopyInto(out *ServiceOrchestratorSpec) {
	*out = *in
	in.StageTarget.DeepCopyInto(&out.StageTarget)
	if in.TargetRevisions != nil {
		in, out := &in.TargetRevisions, &out.TargetRevisions
		*out = make([]TargetRevision, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InitialRevisions != nil {
		in, out := &in.InitialRevisions, &out.InitialRevisions
		*out = make([]TargetRevision, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceOrchestratorSpec.
func (in *ServiceOrchestratorSpec) DeepCopy() *ServiceOrchestratorSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceOrchestratorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceOrchestratorStatus) DeepCopyInto(out *ServiceOrchestratorStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
	in.ServiceOrchestratorStatusFields.DeepCopyInto(&out.ServiceOrchestratorStatusFields)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceOrchestratorStatus.
func (in *ServiceOrchestratorStatus) DeepCopy() *ServiceOrchestratorStatus {
	if in == nil {
		return nil
	}
	out := new(ServiceOrchestratorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceOrchestratorStatusFields) DeepCopyInto(out *ServiceOrchestratorStatusFields) {
	*out = *in
	if in.StageRevisionStatus != nil {
		in, out := &in.StageRevisionStatus, &out.StageRevisionStatus
		*out = make([]TargetRevision, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceOrchestratorStatusFields.
func (in *ServiceOrchestratorStatusFields) DeepCopy() *ServiceOrchestratorStatusFields {
	if in == nil {
		return nil
	}
	out := new(ServiceOrchestratorStatusFields)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StagePodAutoscaler) DeepCopyInto(out *StagePodAutoscaler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StagePodAutoscaler.
func (in *StagePodAutoscaler) DeepCopy() *StagePodAutoscaler {
	if in == nil {
		return nil
	}
	out := new(StagePodAutoscaler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StagePodAutoscaler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StagePodAutoscalerList) DeepCopyInto(out *StagePodAutoscalerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StagePodAutoscaler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StagePodAutoscalerList.
func (in *StagePodAutoscalerList) DeepCopy() *StagePodAutoscalerList {
	if in == nil {
		return nil
	}
	out := new(StagePodAutoscalerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StagePodAutoscalerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StagePodAutoscalerSpec) DeepCopyInto(out *StagePodAutoscalerSpec) {
	*out = *in
	if in.MinScale != nil {
		in, out := &in.MinScale, &out.MinScale
		*out = new(int32)
		**out = **in
	}
	if in.MaxScale != nil {
		in, out := &in.MaxScale, &out.MaxScale
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StagePodAutoscalerSpec.
func (in *StagePodAutoscalerSpec) DeepCopy() *StagePodAutoscalerSpec {
	if in == nil {
		return nil
	}
	out := new(StagePodAutoscalerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StagePodAutoscalerStatus) DeepCopyInto(out *StagePodAutoscalerStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
	if in.DesiredScale != nil {
		in, out := &in.DesiredScale, &out.DesiredScale
		*out = new(int32)
		**out = **in
	}
	if in.ActualScale != nil {
		in, out := &in.ActualScale, &out.ActualScale
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StagePodAutoscalerStatus.
func (in *StagePodAutoscalerStatus) DeepCopy() *StagePodAutoscalerStatus {
	if in == nil {
		return nil
	}
	out := new(StagePodAutoscalerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StageTarget) DeepCopyInto(out *StageTarget) {
	*out = *in
	if in.StageTargetRevisions != nil {
		in, out := &in.StageTargetRevisions, &out.StageTargetRevisions
		*out = make([]TargetRevision, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.TargetFinishTime.DeepCopyInto(&out.TargetFinishTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StageTarget.
func (in *StageTarget) DeepCopy() *StageTarget {
	if in == nil {
		return nil
	}
	out := new(StageTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetRevision) DeepCopyInto(out *TargetRevision) {
	*out = *in
	if in.TargetReplicas != nil {
		in, out := &in.TargetReplicas, &out.TargetReplicas
		*out = new(int32)
		**out = **in
	}
	if in.IsLatestRevision != nil {
		in, out := &in.IsLatestRevision, &out.IsLatestRevision
		*out = new(bool)
		**out = **in
	}
	if in.Percent != nil {
		in, out := &in.Percent, &out.Percent
		*out = new(int64)
		**out = **in
	}
	if in.MinScale != nil {
		in, out := &in.MinScale, &out.MinScale
		*out = new(int32)
		**out = **in
	}
	if in.MaxScale != nil {
		in, out := &in.MaxScale, &out.MaxScale
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetRevision.
func (in *TargetRevision) DeepCopy() *TargetRevision {
	if in == nil {
		return nil
	}
	out := new(TargetRevision)
	in.DeepCopyInto(out)
	return out
}
