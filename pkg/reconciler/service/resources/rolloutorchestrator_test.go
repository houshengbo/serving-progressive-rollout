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

package resources

import (
	"reflect"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/ptr"
	v1 "knative.dev/serving-progressive-rollout/pkg/apis/serving/v1"
	"knative.dev/serving/pkg/apis/autoscaling"
	"knative.dev/serving/pkg/apis/serving"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func TestReadIntAnnotation(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		revision       *servingv1.Revision
		ExpectedResult *int32
	}{{
		name: "Test when both MinScaleAnnotationKey and MaxScaleAnnotationKey exist and valid",
		key:  autoscaling.MinScaleAnnotationKey,
		revision: &servingv1.Revision{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					autoscaling.MinScaleAnnotationKey: "10",
					autoscaling.MaxScaleAnnotationKey: "12",
				},
			},
		},
		ExpectedResult: ptr.Int32(int32(10)),
	}, {
		name: "Test when both MinScaleAnnotationKey and MaxScaleAnnotationKey exist and valid",
		key:  autoscaling.MaxScaleAnnotationKey,
		revision: &servingv1.Revision{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					autoscaling.MinScaleAnnotationKey: "10",
					autoscaling.MaxScaleAnnotationKey: "12",
				},
			},
		},
		ExpectedResult: ptr.Int32(int32(12)),
	}, {
		name: "Test when MinScaleAnnotationKey is missing",
		key:  autoscaling.MinScaleAnnotationKey,
		revision: &servingv1.Revision{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					autoscaling.MaxScaleAnnotationKey: "12",
				},
			},
		},
		ExpectedResult: nil,
	}, {
		name: "Test when MinScaleAnnotationKey is invalid",
		key:  autoscaling.MinScaleAnnotationKey,
		revision: &servingv1.Revision{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					autoscaling.MinScaleAnnotationKey: "12-invalid",
				},
			},
		},
		ExpectedResult: nil,
	}, {
		name: "Test when MaxScaleAnnotationKey is invalid",
		key:  autoscaling.MaxScaleAnnotationKey,
		revision: &servingv1.Revision{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					autoscaling.MaxScaleAnnotationKey: "12-invalid",
				},
			},
		},
		ExpectedResult: nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := ReadIntAnnotation(test.revision, test.key)
			if r == nil {
				if r != test.ExpectedResult {
					t.Fatalf("Result of ReadIntAnnotation() = %v, want %v", r, test.ExpectedResult)
				}
			} else {
				if *r != *test.ExpectedResult {
					t.Fatalf("Result of ReadIntAnnotation() = %v, want %v", *r, *test.ExpectedResult)
				}
			}
		})
	}
}

func TestReadIntServiceAnnotation(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		service        *servingv1.Service
		ExpectedResult *int32
	}{{
		name: "Test when both MinScaleAnnotationKey and MaxScaleAnnotationKey exist and valid",
		key:  autoscaling.MinScaleAnnotationKey,
		service: &servingv1.Service{
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								autoscaling.MinScaleAnnotationKey: "10",
								autoscaling.MaxScaleAnnotationKey: "12",
							},
						},
					},
				},
			},
		},
		ExpectedResult: ptr.Int32(int32(10)),
	}, {
		name: "Test when both MinScaleAnnotationKey and MaxScaleAnnotationKey exist and valid",
		key:  autoscaling.MaxScaleAnnotationKey,
		service: &servingv1.Service{
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								autoscaling.MinScaleAnnotationKey: "10",
								autoscaling.MaxScaleAnnotationKey: "12",
							},
						},
					},
				},
			},
		},
		ExpectedResult: ptr.Int32(int32(12)),
	}, {
		name: "Test when MinScaleAnnotationKey is missing",
		key:  autoscaling.MinScaleAnnotationKey,
		service: &servingv1.Service{
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								autoscaling.MaxScaleAnnotationKey: "12",
							},
						},
					},
				},
			},
		},
		ExpectedResult: nil,
	}, {
		name: "Test when MinScaleAnnotationKey is invalid",
		key:  autoscaling.MinScaleAnnotationKey,
		service: &servingv1.Service{
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								autoscaling.MinScaleAnnotationKey: "10-invalid",
							},
						},
					},
				},
			},
		},
		ExpectedResult: nil,
	}, {
		name: "Test when MaxScaleAnnotationKey is invalid",
		key:  autoscaling.MaxScaleAnnotationKey,
		service: &servingv1.Service{
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								autoscaling.MaxScaleAnnotationKey: "10-invalid",
							},
						},
					},
				},
			},
		},
		ExpectedResult: nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := ReadIntServiceAnnotation(test.service, test.key)
			if r == nil {
				if r != test.ExpectedResult {
					t.Fatalf("Result of ReadIntServiceAnnotation() = %v, want %v", r, test.ExpectedResult)
				}
			} else {
				if *r != *test.ExpectedResult {
					t.Fatalf("Result of ReadIntServiceAnnotation() = %v, want %v", *r, *test.ExpectedResult)
				}
			}
		})
	}
}

func TestReadIntRevisionRecord(t *testing.T) {
	tests := []struct {
		name           string
		RevisionRecord RevisionRecord
		ExpectedMin    *int32
		ExpectedMax    *int32
	}{{
		name: "Test when both MinScale and MaxScale exist and valid",
		RevisionRecord: RevisionRecord{
			MinScale: ptr.Int32(10),
			MaxScale: ptr.Int32(15),
		},
		ExpectedMin: ptr.Int32(10),
		ExpectedMax: ptr.Int32(15),
	}, {
		name:           "Test when both MinScale and MaxScale does not exist",
		RevisionRecord: RevisionRecord{},
		ExpectedMin:    nil,
		ExpectedMax:    nil,
	}, {
		name: "Test when MinScale is valid",
		RevisionRecord: RevisionRecord{
			MinScale: ptr.Int32(15),
		},
		ExpectedMin: ptr.Int32(15),
		ExpectedMax: nil,
	}, {
		name: "Test when MaxScale is valid",
		RevisionRecord: RevisionRecord{
			MaxScale: ptr.Int32(15),
		},
		ExpectedMin: nil,
		ExpectedMax: ptr.Int32(15),
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			minR, maxR := ReadIntRevisionRecord(test.RevisionRecord)
			if minR == nil {
				if minR != test.ExpectedMin {
					t.Fatalf("Result of ReadIntRevisionRecord() = %v, want %v", minR, test.ExpectedMin)
				}
			} else {
				if *minR != *test.ExpectedMin {
					t.Fatalf("Result of ReadIntRevisionRecord() = %v, want %v", *minR, *test.ExpectedMin)
				}
			}
			if maxR == nil {
				if maxR != test.ExpectedMax {
					t.Fatalf("Result of ReadIntRevisionRecord() = %v, want %v", maxR, test.ExpectedMax)
				}
			} else {
				if *maxR != *test.ExpectedMax {
					t.Fatalf("Result of ReadIntRevisionRecord() = %v, want %v", *maxR, *test.ExpectedMax)
				}
			}
		})
	}
}

func TestNewInitialFinalTargetRev(t *testing.T) {
	service := &servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-001",
			Namespace: "test-ns",
		},
	}
	tests := []struct {
		name                   string
		initialRevisionStatus  []v1.TargetRevision
		ultimateRevisionTarget []v1.TargetRevision
		service                *servingv1.Service
		ExpectedResult         *v1.RolloutOrchestrator
	}{{
		name:                   "Test the creation of RolloutOrchestrator with empty initial and final target revisions",
		initialRevisionStatus:  []v1.TargetRevision{},
		ultimateRevisionTarget: []v1.TargetRevision{},
		service:                service,
		ExpectedResult: &v1.RolloutOrchestrator{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-001",
				Namespace: "test-ns",
				Labels:    map[string]string{serving.ServiceLabelKey: "service-001"},
				OwnerReferences: []metav1.OwnerReference{
					*kmeta.NewControllerRef(service),
				},
			},
			Spec: v1.RolloutOrchestratorSpec{
				TargetRevisions:  []v1.TargetRevision{},
				InitialRevisions: []v1.TargetRevision{},
			},
		},
	}, {
		name: "Test the creation of RolloutOrchestrator with initial and final target revisions",
		initialRevisionStatus: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "rev-001",
					LatestRevision: ptr.Bool(false),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(15),
				MaxScale: ptr.Int32(15),
			},
		},
		ultimateRevisionTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "rev-002",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(15),
				MaxScale: ptr.Int32(15),
			},
		},
		service: service,
		ExpectedResult: &v1.RolloutOrchestrator{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-001",
				Namespace: "test-ns",
				Labels:    map[string]string{serving.ServiceLabelKey: "service-001"},
				OwnerReferences: []metav1.OwnerReference{
					*kmeta.NewControllerRef(service),
				},
			},
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "rev-001",
							LatestRevision: ptr.Bool(false),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(15),
						MaxScale: ptr.Int32(15),
					},
				},
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "rev-002",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(15),
						MaxScale: ptr.Int32(15),
					},
				},
			},
		},
	}, {
		name:                  "Test the creation of RolloutOrchestrator with final target revisions",
		initialRevisionStatus: nil,
		ultimateRevisionTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "rev-002",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(15),
				MaxScale: ptr.Int32(15),
			},
		},
		service: service,
		ExpectedResult: &v1.RolloutOrchestrator{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-001",
				Namespace: "test-ns",
				Labels:    map[string]string{serving.ServiceLabelKey: "service-001"},
				OwnerReferences: []metav1.OwnerReference{
					*kmeta.NewControllerRef(service),
				},
			},
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: nil,
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "rev-002",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(15),
						MaxScale: ptr.Int32(15),
					},
				},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewInitialFinalTargetRev(test.initialRevisionStatus, test.ultimateRevisionTarget, test.service)
			if !reflect.DeepEqual(r, test.ExpectedResult) {
				t.Fatalf("Result of NewInitialFinalTargetRev() = %v, want %v", r, test.ExpectedResult)
			}
		})
	}
}

func TestGetInitialFinalTargetRevision(t *testing.T) {
	tests := []struct {
		name                      string
		records                   map[string]RevisionRecord
		route                     *servingv1.Route
		service                   *servingv1.Service
		config                    *servingv1.Configuration
		ExpectedInitialTarget     []v1.TargetRevision
		ExpectedFinalTargetResult []v1.TargetRevision
	}{{
		name:    "Test the creation/update of RolloutOrchestrator with empty records and route",
		records: map[string]RevisionRecord{},
		route:   &servingv1.Route{},
		config: &servingv1.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 2,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
		},
		service: &servingv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 2,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								autoscaling.MinScaleAnnotationKey: "11",
								autoscaling.MaxScaleAnnotationKey: "12",
							},
						},
					},
				},
			},
		},
		ExpectedInitialTarget: nil,
		ExpectedFinalTargetResult: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00002",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(11),
				MaxScale: ptr.Int32(12),
			},
		},
	}, {
		name: "Test the creation/update of RolloutOrchestrator with valid records",
		records: map[string]RevisionRecord{
			"service-001-00001": {
				MinScale: ptr.Int32(5),
				MaxScale: ptr.Int32(10),
				Name:     "service-001-00001",
			},
			"service-001-00002": {
				MinScale: ptr.Int32(5),
				MaxScale: ptr.Int32(10),
				Name:     "service-001-00001",
			},
		},
		route: &servingv1.Route{},
		config: &servingv1.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 2,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
		},
		service: &servingv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 2,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								autoscaling.MinScaleAnnotationKey: "10",
							},
						},
					},
				},
			},
		},
		ExpectedInitialTarget: nil,
		ExpectedFinalTargetResult: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00002",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(5),
				MaxScale: ptr.Int32(10),
			},
		},
	}, {
		name: "Test the creation/update of RolloutOrchestrator with valid records, route and an existing generation",
		records: map[string]RevisionRecord{
			"service-001-00001": {
				MinScale: ptr.Int32(5),
				MaxScale: ptr.Int32(10),
				Name:     "service-001-00001",
			},
		},
		route: &servingv1.Route{
			Status: servingv1.RouteStatus{
				RouteStatusFields: servingv1.RouteStatusFields{
					Traffic: []servingv1.TrafficTarget{
						{
							RevisionName: "service-001-00001",
							Percent:      ptr.Int64(100),
						},
					},
				},
			},
		},
		config: &servingv1.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 2,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
		},
		service: &servingv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 2,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								autoscaling.MinScaleAnnotationKey: "8",
								autoscaling.MaxScaleAnnotationKey: "10",
							},
						},
					},
				},
			},
		},
		ExpectedInitialTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00001",
					LatestRevision: ptr.Bool(false),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(5),
				MaxScale: ptr.Int32(10),
			},
		},
		ExpectedFinalTargetResult: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00002",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
			},
		},
	}, {
		name: "Test the creation/update of RolloutOrchestrator with valid records, route and an existing generate",
		records: map[string]RevisionRecord{
			"service-001-00001": {
				MinScale: ptr.Int32(5),
				MaxScale: ptr.Int32(10),
				Name:     "service-001-00001",
			},
			"service-001-00002": {
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
				Name:     "service-001-00002",
			},
		},
		route: &servingv1.Route{
			Status: servingv1.RouteStatus{
				RouteStatusFields: servingv1.RouteStatusFields{
					Traffic: []servingv1.TrafficTarget{
						{
							RevisionName: "service-001-00001",
							Percent:      ptr.Int64(40),
						},
						{
							RevisionName: "service-001-00002",
							Percent:      ptr.Int64(60),
						},
					},
				},
			},
		},
		config: &servingv1.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 2,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
		},
		service: &servingv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					autoscaling.MinScaleAnnotationKey: "8",
					autoscaling.MaxScaleAnnotationKey: "10",
				},
				Generation: 2,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
		},
		ExpectedInitialTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00001",
					LatestRevision: ptr.Bool(false),
					Percent:        ptr.Int64(40),
				},
				MinScale: ptr.Int32(5),
				MaxScale: ptr.Int32(10),
			},
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00002",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(60),
				},
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
			},
		},
		ExpectedFinalTargetResult: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00002",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
			},
		},
	}, {
		name: "Test the creation/update of RolloutOrchestrator with valid records, route and a new generate",
		records: map[string]RevisionRecord{
			"service-001-00001": {
				MinScale: ptr.Int32(5),
				MaxScale: ptr.Int32(10),
				Name:     "service-001-00001",
			},
			"service-001-00002": {
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
				Name:     "service-001-00002",
			},
		},
		route: &servingv1.Route{
			Status: servingv1.RouteStatus{
				RouteStatusFields: servingv1.RouteStatusFields{
					Traffic: []servingv1.TrafficTarget{
						{
							RevisionName: "service-001-00001",
							Percent:      ptr.Int64(40),
						},
						{
							RevisionName: "service-001-00002",
							Percent:      ptr.Int64(60),
						},
					},
				},
			},
		},
		config: &servingv1.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 3,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
		},
		service: &servingv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 3,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								autoscaling.MinScaleAnnotationKey: "8",
								autoscaling.MaxScaleAnnotationKey: "10",
							},
						},
					},
				},
			},
		},
		ExpectedInitialTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00001",
					LatestRevision: ptr.Bool(false),
					Percent:        ptr.Int64(40),
				},
				MinScale: ptr.Int32(5),
				MaxScale: ptr.Int32(10),
			},
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00002",
					LatestRevision: ptr.Bool(false),
					Percent:        ptr.Int64(60),
				},
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
			},
		},
		ExpectedFinalTargetResult: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00003",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
			},
		},
	}, {
		name: "Test the creation/update of RolloutOrchestrator with valid records, route, and traffic defined in the service",
		records: map[string]RevisionRecord{
			"service-001-00001": {
				MinScale: ptr.Int32(5),
				MaxScale: ptr.Int32(10),
				Name:     "service-001-00001",
			},
			"service-001-00002": {
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
				Name:     "service-001-00002",
			},
		},
		route: &servingv1.Route{
			Status: servingv1.RouteStatus{
				RouteStatusFields: servingv1.RouteStatusFields{
					Traffic: []servingv1.TrafficTarget{
						{
							RevisionName: "service-001-00002",
							Percent:      ptr.Int64(100),
						},
					},
				},
			},
		},
		config: &servingv1.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 3,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
		},
		service: &servingv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 3,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								autoscaling.MinScaleAnnotationKey: "8",
								autoscaling.MaxScaleAnnotationKey: "10",
							},
						},
					},
				},
				RouteSpec: servingv1.RouteSpec{
					Traffic: []servingv1.TrafficTarget{{
						RevisionName:   "service-001-00002",
						Percent:        ptr.Int64(30),
						LatestRevision: ptr.Bool(false),
					}, {
						Percent:        ptr.Int64(70),
						LatestRevision: ptr.Bool(true),
					}},
				},
			},
		},
		ExpectedInitialTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00002",
					LatestRevision: ptr.Bool(false),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
			},
		},
		ExpectedFinalTargetResult: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00002",
					LatestRevision: ptr.Bool(false),
					Percent:        ptr.Int64(30),
				},
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
			},
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00003",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(70),
				},
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
			},
		},
	}, {
		name: "Test the creation/update of RolloutOrchestrator with valid records, route, and one traffic in the service",
		records: map[string]RevisionRecord{
			"service-001-00002": {
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
				Name:     "service-001-00002",
			},
		},
		route: &servingv1.Route{
			Status: servingv1.RouteStatus{
				RouteStatusFields: servingv1.RouteStatusFields{
					Traffic: []servingv1.TrafficTarget{
						{
							RevisionName: "service-001-00002",
							Percent:      ptr.Int64(100),
						},
					},
				},
			},
		},
		config: &servingv1.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 3,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
		},
		service: &servingv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Generation: 3,
				Namespace:  "test-ns",
				Name:       "service-001",
			},
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								autoscaling.MinScaleAnnotationKey: "8",
								autoscaling.MaxScaleAnnotationKey: "10",
							},
						},
					},
				},
				RouteSpec: servingv1.RouteSpec{
					Traffic: []servingv1.TrafficTarget{{
						Percent:        ptr.Int64(100),
						LatestRevision: ptr.Bool(true),
					}},
				},
			},
		},
		ExpectedInitialTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00002",
					LatestRevision: ptr.Bool(false),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
			},
		},
		ExpectedFinalTargetResult: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00003",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(8),
				MaxScale: ptr.Int32(10),
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initialTarget := GetInitialTargetRevision(test.service, test.config, test.records, test.route)
			finalTarget := GetFinalTargetRevision(test.service, test.config, test.records)
			if !reflect.DeepEqual(initialTarget, test.ExpectedInitialTarget) {
				t.Fatalf("Result of GetInitialFinalTargetRevision() = %v, want %v", initialTarget, test.ExpectedInitialTarget)
			}
			if !reflect.DeepEqual(finalTarget, test.ExpectedFinalTargetResult) {
				t.Fatalf("Result of GetInitialFinalTargetRevision() = %v, want %v", finalTarget, test.ExpectedFinalTargetResult)
			}
		})
	}
}

type MockEmptyDeploymentLister struct{}

func (lister *MockEmptyDeploymentLister) List(_ labels.Selector) (ret []*appsv1.Deployment, err error) {
	return nil, nil
}

func (lister *MockEmptyDeploymentLister) Deployments(_ string) appsv1listers.DeploymentNamespaceLister {
	return &MockEmptyDeploymentNamespaceLister{}
}

type MockEmptyDeploymentNamespaceLister struct{}

func (lister *MockEmptyDeploymentNamespaceLister) List(_ labels.Selector) (ret []*appsv1.Deployment, err error) {
	return nil, nil
}

func (lister *MockEmptyDeploymentNamespaceLister) Get(_ string) (*appsv1.Deployment, error) {
	return nil, nil
}

type MockDeploymentLister struct{}

func (lister *MockDeploymentLister) List(_ labels.Selector) (ret []*appsv1.Deployment, err error) {
	return nil, nil
}

func (lister *MockDeploymentLister) Deployments(_ string) appsv1listers.DeploymentNamespaceLister {
	return &MockDeploymentNamespaceLister{}
}

type MockDeploymentNamespaceLister struct{}

func (lister *MockDeploymentNamespaceLister) List(_ labels.Selector) (ret []*appsv1.Deployment, err error) {
	return []*appsv1.Deployment{{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mock-deployment",
			Namespace: "mock-deployment-ns",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"a": "b",
				},
			},
			Replicas: ptr.Int32(2),
		},
		Status: appsv1.DeploymentStatus{
			Replicas:          2,
			AvailableReplicas: 2,
		},
	}}, nil
}

func (lister *MockDeploymentNamespaceLister) Get(_ string) (*appsv1.Deployment, error) {
	return nil, nil
}

func TestUpdateInitialFinalTargetRev(t *testing.T) {
	testTime := apis.VolatileTime{
		Inner: metav1.NewTime(time.Now()),
	}
	targetRevisions := []v1.TargetRevision{
		{
			TrafficTarget: servingv1.TrafficTarget{
				RevisionName:   "service-001-00003",
				LatestRevision: ptr.Bool(true),
				Percent:        ptr.Int64(100),
			},
			MinScale: ptr.Int32(8),
			MaxScale: ptr.Int32(10),
		},
	}
	ro := &v1.RolloutOrchestrator{
		Spec: v1.RolloutOrchestratorSpec{
			TargetRevisions: targetRevisions,
			StageTarget: v1.StageTarget{
				StageTargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00002",
							LatestRevision: ptr.Bool(false),
							Percent:        ptr.Int64(20),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00003",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(80),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				TargetFinishTime: testTime,
			},
		},
	}

	tests := []struct {
		name             string
		ultimateTarget   []v1.TargetRevision
		ro               *v1.RolloutOrchestrator
		route            *servingv1.Route
		deploymentLister appsv1listers.DeploymentLister
		records          map[string]RevisionRecord
		ExpectedResult   *v1.RolloutOrchestrator
	}{{
		name:           "Test the UpdateFinalTargetRev with ultimateTarget",
		ultimateTarget: []v1.TargetRevision{},
		ro:             &v1.RolloutOrchestrator{},
		ExpectedResult: &v1.RolloutOrchestrator{},
	}, {
		name:           "Test the UpdateFinalTargetRev, when the updated target revisions are equal to the existing ones",
		ultimateTarget: targetRevisions,
		ro:             ro,
		ExpectedResult: ro,
	}, {
		name: "Test the UpdateFinalTargetRev, when the updated target revisions become different with no status or route",
		ultimateTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00004",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(9),
				MaxScale: ptr.Int32(19),
			},
		},
		ro: &v1.RolloutOrchestrator{
			Spec: v1.RolloutOrchestratorSpec{
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00003",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
					TargetFinishTime: testTime,
				},
			},
		},
		ExpectedResult: &v1.RolloutOrchestrator{
			Spec: v1.RolloutOrchestratorSpec{
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00004",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(9),
						MaxScale: ptr.Int32(19),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: nil,
					TargetFinishTime:     apis.VolatileTime{},
				},
			},
		},
	}, {
		name: "Test the UpdateFinalTargetRev with a new final target, no status or route",
		ultimateTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00004",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(9),
				MaxScale: ptr.Int32(19),
			},
		},
		ro: &v1.RolloutOrchestrator{
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00002",
							LatestRevision: ptr.Bool(false),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00003",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
					TargetFinishTime: testTime,
				},
			},
		},
		ExpectedResult: &v1.RolloutOrchestrator{
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: nil,
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00004",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(9),
						MaxScale: ptr.Int32(19),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: nil,
					TargetFinishTime:     apis.VolatileTime{},
				},
			},
		},
	}, {
		name: "Test the UpdateFinalTargetRev with a new final target with status stage ready",
		ultimateTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00004",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(9),
				MaxScale: ptr.Int32(19),
			},
		},
		ro: &v1.RolloutOrchestrator{
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00002",
							LatestRevision: ptr.Bool(false),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00003",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
					TargetFinishTime: testTime,
				},
			},
			Status: v1.RolloutOrchestratorStatus{
				Status: duckv1.Status{
					Conditions: duckv1.Conditions{{
						Type:   v1.SOStageReady,
						Status: corev1.ConditionTrue,
					}},
				},
				RolloutOrchestratorStatusFields: v1.RolloutOrchestratorStatusFields{
					StageRevisionStatus: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
				},
			},
		},
		ExpectedResult: &v1.RolloutOrchestrator{
			Status: v1.RolloutOrchestratorStatus{
				Status: duckv1.Status{
					Conditions: duckv1.Conditions{{
						Type:   v1.SOStageReady,
						Status: corev1.ConditionTrue,
					}},
				},
				RolloutOrchestratorStatusFields: v1.RolloutOrchestratorStatusFields{
					StageRevisionStatus: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
				},
			},
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00002",
							LatestRevision: ptr.Bool(false),
							Percent:        ptr.Int64(20),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00003",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(80),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00004",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(9),
						MaxScale: ptr.Int32(19),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: nil,
					TargetFinishTime:     apis.VolatileTime{},
				},
			},
		},
	}, {
		name: "Test the UpdateFinalTargetRev with a new final target with status stage not ready",
		ultimateTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00004",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(9),
				MaxScale: ptr.Int32(19),
			},
		},
		ro: &v1.RolloutOrchestrator{
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00002",
							LatestRevision: ptr.Bool(false),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00003",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
					TargetFinishTime: testTime,
				},
			},
			Status: v1.RolloutOrchestratorStatus{
				Status: duckv1.Status{
					Conditions: duckv1.Conditions{{
						Type:   v1.SOStageReady,
						Status: corev1.ConditionUnknown,
					}},
				},
				RolloutOrchestratorStatusFields: v1.RolloutOrchestratorStatusFields{
					StageRevisionStatus: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
				},
			},
		},
		ExpectedResult: &v1.RolloutOrchestrator{
			Status: v1.RolloutOrchestratorStatus{
				Status: duckv1.Status{
					Conditions: duckv1.Conditions{{
						Type:   v1.SOStageReady,
						Status: corev1.ConditionUnknown,
					}},
				},
				RolloutOrchestratorStatusFields: v1.RolloutOrchestratorStatusFields{
					StageRevisionStatus: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
				},
			},
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: nil,
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00004",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(9),
						MaxScale: ptr.Int32(19),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: nil,
					TargetFinishTime:     apis.VolatileTime{},
				},
			},
		},
	}, {
		name: "Test the UpdateFinalTargetRev with a new final target with route and status stage not ready, deployment not ready",
		ultimateTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00004",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(9),
				MaxScale: ptr.Int32(19),
			},
		},
		deploymentLister: &MockEmptyDeploymentLister{},
		route: &servingv1.Route{
			Spec: servingv1.RouteSpec{
				Traffic: []servingv1.TrafficTarget{{
					RevisionName:   "service-001-00002",
					LatestRevision: ptr.Bool(false),
					Percent:        ptr.Int64(80),
				}, {
					RevisionName:   "service-001-00003",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(20),
				}},
			},
			Status: servingv1.RouteStatus{
				RouteStatusFields: servingv1.RouteStatusFields{
					Traffic: []servingv1.TrafficTarget{{
						RevisionName:   "service-001-00002",
						LatestRevision: ptr.Bool(false),
						Percent:        ptr.Int64(80),
					}, {
						RevisionName:   "service-001-00003",
						LatestRevision: ptr.Bool(true),
						Percent:        ptr.Int64(20),
					}},
				},
			},
		},
		ro: &v1.RolloutOrchestrator{
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00002",
							LatestRevision: ptr.Bool(false),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00003",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(60),
							},
							MinScale:  ptr.Int32(8),
							MaxScale:  ptr.Int32(10),
							Direction: "down",
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(40),
							},
							MinScale:  ptr.Int32(8),
							MaxScale:  ptr.Int32(10),
							Direction: "up",
						},
					},
					TargetFinishTime: testTime,
				},
			},
			Status: v1.RolloutOrchestratorStatus{
				Status: duckv1.Status{
					Conditions: duckv1.Conditions{{
						Type:   v1.SOStageReady,
						Status: corev1.ConditionUnknown,
					}},
				},
				RolloutOrchestratorStatusFields: v1.RolloutOrchestratorStatusFields{
					StageRevisionStatus: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
				},
			},
		},
		ExpectedResult: &v1.RolloutOrchestrator{
			Status: v1.RolloutOrchestratorStatus{
				Status: duckv1.Status{
					Conditions: duckv1.Conditions{{
						Type:   v1.SOStageReady,
						Status: corev1.ConditionUnknown,
					}},
				},
				RolloutOrchestratorStatusFields: v1.RolloutOrchestratorStatusFields{
					StageRevisionStatus: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
				},
			},
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00002",
							LatestRevision: ptr.Bool(false),
							Percent:        ptr.Int64(80),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00003",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(20),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00004",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(9),
						MaxScale: ptr.Int32(19),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: nil,
					TargetFinishTime:     apis.VolatileTime{},
				},
			},
		},
	}, {
		name: "Test the UpdateFinalTargetRev with a new final target with route and status stage not ready, deployment ready",
		ultimateTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00004",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(9),
				MaxScale: ptr.Int32(19),
			},
		},
		deploymentLister: &MockDeploymentLister{},
		route: &servingv1.Route{
			Spec: servingv1.RouteSpec{
				Traffic: []servingv1.TrafficTarget{{
					RevisionName:   "service-001-00002",
					LatestRevision: ptr.Bool(false),
					Percent:        ptr.Int64(80),
				}, {
					RevisionName:   "service-001-00003",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(20),
				}},
			},
			Status: servingv1.RouteStatus{
				RouteStatusFields: servingv1.RouteStatusFields{
					Traffic: []servingv1.TrafficTarget{{
						RevisionName:   "service-001-00002",
						LatestRevision: ptr.Bool(false),
						Percent:        ptr.Int64(80),
					}, {
						RevisionName:   "service-001-00003",
						LatestRevision: ptr.Bool(true),
						Percent:        ptr.Int64(20),
					}},
				},
			},
		},
		ro: &v1.RolloutOrchestrator{
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00002",
							LatestRevision: ptr.Bool(false),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00003",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(60),
							},
							MinScale:  ptr.Int32(8),
							MaxScale:  ptr.Int32(10),
							Direction: "down",
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(40),
							},
							MinScale:  ptr.Int32(8),
							MaxScale:  ptr.Int32(10),
							Direction: "up",
						},
					},
					TargetFinishTime: testTime,
				},
			},
			Status: v1.RolloutOrchestratorStatus{
				Status: duckv1.Status{
					Conditions: duckv1.Conditions{{
						Type:   v1.SOStageReady,
						Status: corev1.ConditionUnknown,
					}},
				},
				RolloutOrchestratorStatusFields: v1.RolloutOrchestratorStatusFields{
					StageRevisionStatus: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
				},
			},
		},
		ExpectedResult: &v1.RolloutOrchestrator{
			Status: v1.RolloutOrchestratorStatus{
				Status: duckv1.Status{
					Conditions: duckv1.Conditions{{
						Type:   v1.SOStageReady,
						Status: corev1.ConditionUnknown,
					}},
				},
				RolloutOrchestratorStatusFields: v1.RolloutOrchestratorStatusFields{
					StageRevisionStatus: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
				},
			},
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00002",
							LatestRevision: ptr.Bool(false),
							Percent:        ptr.Int64(60),
						},
						MinScale:  ptr.Int32(8),
						MaxScale:  ptr.Int32(10),
						Direction: "down",
					},
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00003",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(40),
						},
						MinScale:  ptr.Int32(8),
						MaxScale:  ptr.Int32(10),
						Direction: "up",
					},
				},
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00004",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(9),
						MaxScale: ptr.Int32(19),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: nil,
					TargetFinishTime:     apis.VolatileTime{},
				},
			},
		},
	}, {
		name: "Test the UpdateFinalTargetRev with a new final target and route, but no status",
		ultimateTarget: []v1.TargetRevision{
			{
				TrafficTarget: servingv1.TrafficTarget{
					RevisionName:   "service-001-00004",
					LatestRevision: ptr.Bool(true),
					Percent:        ptr.Int64(100),
				},
				MinScale: ptr.Int32(9),
				MaxScale: ptr.Int32(19),
			},
		},
		route: &servingv1.Route{
			Spec: servingv1.RouteSpec{
				Traffic: []servingv1.TrafficTarget{{
					ConfigurationName: "service-001",
					LatestRevision:    ptr.Bool(true),
					Percent:           ptr.Int64(100),
				}},
			},
			Status: servingv1.RouteStatus{
				RouteStatusFields: servingv1.RouteStatusFields{
					Traffic: []servingv1.TrafficTarget{{
						RevisionName:   "service-001-00002",
						LatestRevision: ptr.Bool(true),
						Percent:        ptr.Int64(100),
					}},
				},
			},
		},
		ro: &v1.RolloutOrchestrator{
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00002",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00003",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: []v1.TargetRevision{
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00002",
								LatestRevision: ptr.Bool(false),
								Percent:        ptr.Int64(80),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
						{
							TrafficTarget: servingv1.TrafficTarget{
								RevisionName:   "service-001-00003",
								LatestRevision: ptr.Bool(true),
								Percent:        ptr.Int64(20),
							},
							MinScale: ptr.Int32(8),
							MaxScale: ptr.Int32(10),
						},
					},
					TargetFinishTime: testTime,
				},
			},
		},
		ExpectedResult: &v1.RolloutOrchestrator{
			Spec: v1.RolloutOrchestratorSpec{
				InitialRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00002",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(8),
						MaxScale: ptr.Int32(10),
					},
				},
				TargetRevisions: []v1.TargetRevision{
					{
						TrafficTarget: servingv1.TrafficTarget{
							RevisionName:   "service-001-00004",
							LatestRevision: ptr.Bool(true),
							Percent:        ptr.Int64(100),
						},
						MinScale: ptr.Int32(9),
						MaxScale: ptr.Int32(19),
					},
				},
				StageTarget: v1.StageTarget{
					StageTargetRevisions: nil,
					TargetFinishTime:     apis.VolatileTime{},
				},
			},
		},
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			UpdateInitialFinalTargetRev(test.ultimateTarget, test.ro, test.route, test.deploymentLister)
			if !reflect.DeepEqual(test.ro, test.ExpectedResult) {
				t.Fatalf("Result of UpdateFinalTargetRev() = %v, want %v", test.ro, test.ExpectedResult)
			}
		})
	}
}

func TestConsolidateTraffic(t *testing.T) {
	tests := []struct {
		name               string
		routeStatusTraffic []servingv1.TrafficTarget
		ExpectedResult     []servingv1.TrafficTarget
	}{{
		name: "Test the consolidateTraffic without repeated revision",
		routeStatusTraffic: []servingv1.TrafficTarget{{
			LatestRevision: ptr.Bool(true),
			RevisionName:   "rev-0003",
			Percent:        ptr.Int64(10),
		}, {
			LatestRevision: ptr.Bool(false),
			RevisionName:   "rev-0001",
			Percent:        ptr.Int64(70),
		}, {
			LatestRevision: ptr.Bool(false),
			RevisionName:   "rev-0002",
			Percent:        ptr.Int64(20),
		}},
		ExpectedResult: []servingv1.TrafficTarget{{
			LatestRevision: ptr.Bool(false),
			RevisionName:   "rev-0001",
			Percent:        ptr.Int64(70),
		}, {
			LatestRevision: ptr.Bool(false),
			RevisionName:   "rev-0002",
			Percent:        ptr.Int64(20),
		}, {
			LatestRevision: ptr.Bool(true),
			RevisionName:   "rev-0003",
			Percent:        ptr.Int64(10),
		}},
	}, {
		name: "Test the consolidateTraffic with repeated revision",
		routeStatusTraffic: []servingv1.TrafficTarget{{
			LatestRevision: ptr.Bool(true),
			RevisionName:   "rev-0003",
			Percent:        ptr.Int64(10),
		}, {
			LatestRevision: ptr.Bool(false),
			RevisionName:   "rev-0001",
			Percent:        ptr.Int64(40),
		}, {
			LatestRevision: ptr.Bool(false),
			RevisionName:   "rev-0001",
			Percent:        ptr.Int64(30),
		}, {
			LatestRevision: ptr.Bool(false),
			RevisionName:   "rev-0003",
			Percent:        ptr.Int64(10),
		}, {
			LatestRevision: ptr.Bool(false),
			RevisionName:   "rev-0002",
			Percent:        ptr.Int64(10),
		}},
		ExpectedResult: []servingv1.TrafficTarget{{
			LatestRevision: ptr.Bool(false),
			RevisionName:   "rev-0001",
			Percent:        ptr.Int64(70),
		}, {
			LatestRevision: ptr.Bool(false),
			RevisionName:   "rev-0002",
			Percent:        ptr.Int64(10),
		}, {
			LatestRevision: ptr.Bool(true),
			RevisionName:   "rev-0003",
			Percent:        ptr.Int64(20),
		}},
	}, {
		name: "Test the consolidateTraffic with repeated revision",
		routeStatusTraffic: []servingv1.TrafficTarget{{
			LatestRevision: ptr.Bool(false),
			RevisionName:   "demo-triton-v-predictor-00001",
			Percent:        ptr.Int64(50),
		}, {
			LatestRevision: ptr.Bool(false),
			RevisionName:   "demo-triton-v-predictor-00002",
			Percent:        ptr.Int64(25),
		}, {
			LatestRevision: ptr.Bool(true),
			RevisionName:   "demo-triton-v-predictor-00002",
			Percent:        ptr.Int64(25),
		}},
		ExpectedResult: []servingv1.TrafficTarget{{
			LatestRevision: ptr.Bool(false),
			RevisionName:   "demo-triton-v-predictor-00001",
			Percent:        ptr.Int64(50),
		}, {
			LatestRevision: ptr.Bool(true),
			RevisionName:   "demo-triton-v-predictor-00002",
			Percent:        ptr.Int64(50),
		}},
	}, {
		name:               "Test the consolidateTraffic with empty input",
		routeStatusTraffic: []servingv1.TrafficTarget{},
		ExpectedResult:     []servingv1.TrafficTarget{},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := consolidateTraffic(test.routeStatusTraffic)
			if !reflect.DeepEqual(res, test.ExpectedResult) {
				t.Fatalf("Result of consolidateTraffic() = %v, want %v", res, test.ExpectedResult)
			}
		})
	}
}
