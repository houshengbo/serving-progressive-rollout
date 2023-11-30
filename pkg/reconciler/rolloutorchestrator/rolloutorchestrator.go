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

package rolloutorchestrator

import (
	"context"
	"fmt"
	"math"

	"k8s.io/apimachinery/pkg/api/equality"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/ptr"
	pkgreconciler "knative.dev/pkg/reconciler"
	v1 "knative.dev/serving-progressive-rollout/pkg/apis/serving/v1"
	clientset "knative.dev/serving-progressive-rollout/pkg/client/clientset/versioned"
	roreconciler "knative.dev/serving-progressive-rollout/pkg/client/injection/reconciler/serving/v1/rolloutorchestrator"
	listers "knative.dev/serving-progressive-rollout/pkg/client/listers/serving/v1"
	"knative.dev/serving/pkg/apis/serving"
)

// Reconciler implements controller.Reconciler for RolloutOrchestrator resources.
type Reconciler struct {
	client clientset.Interface

	// lister indexes properties about StagePodAutoscaler
	stagePodAutoscalerLister listers.StagePodAutoscalerLister
}

// Check that our Reconciler implements roreconciler.Interface
var _ roreconciler.Interface = (*Reconciler)(nil)

// createOrUpdateSPARevDown create or update the StagePodAutoscaler, based on the specific (Stage)TargetRevision
// defined in the RolloutOrchestrator for the revision scaling down.
func (r *Reconciler) createOrUpdateSPARevDown(ctx context.Context, ro *v1.RolloutOrchestrator,
	targetRev *v1.TargetRevision, scaleUpReady bool) (*v1.StagePodAutoscaler, error) {
	spa, err := r.stagePodAutoscalerLister.StagePodAutoscalers(ro.Namespace).Get(targetRev.RevisionName)
	if apierrs.IsNotFound(err) {
		return r.createStagePA(ctx, ro, targetRev, scaleUpReady, UpdateSPAForRevDown)
	}
	if err != nil {
		return spa, err
	}
	return r.client.ServingV1().StagePodAutoscalers(ro.Namespace).Update(ctx,
		UpdateSPAForRevDown(spa, targetRev, scaleUpReady), metav1.UpdateOptions{})
}

// createOrUpdateSPARevUp create or update the StagePodAutoscaler, based on the specific (Stage)TargetRevision
// defined in the RolloutOrchestrator for the revision scaling up.
func (r *Reconciler) createOrUpdateSPARevUp(ctx context.Context, ro *v1.RolloutOrchestrator,
	targetRev *v1.TargetRevision) (*v1.StagePodAutoscaler, error) {
	spa, err := r.stagePodAutoscalerLister.StagePodAutoscalers(ro.Namespace).Get(targetRev.RevisionName)
	if apierrs.IsNotFound(err) {
		return r.createStagePA(ctx, ro, targetRev, true, UpdateSPAForRevUp)
	}
	if err != nil {
		return spa, err
	}
	return r.client.ServingV1().StagePodAutoscalers(ro.Namespace).Update(ctx,
		UpdateSPAForRevUp(spa, targetRev, true), metav1.UpdateOptions{})
}

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, ro *v1.RolloutOrchestrator) pkgreconciler.Event {
	ctx, cancel := context.WithTimeout(ctx, pkgreconciler.DefaultTimeout)
	defer cancel()

	// If spec.StageRevisionStatus is nil, do nothing.
	if len(ro.Spec.StageTargetRevisions) == 0 {
		return nil
	}

	// Spec.StageTargetRevisions in the RolloutOrchestrator defines what the current stage looks like, in terms
	// of the available revisions, and their name, traffic percentage, target number of replicas, whether it
	// scales up or down, min and max scales defined by the Knative Service.
	stageTargetRevisions := ro.Spec.StageTargetRevisions
	revScalingUp, revScalingDown, err := RetrieveRevsUpDown(stageTargetRevisions)
	if err != nil {
		return err
	}

	// Create or update the StagePodAutoscaler for the revision to be scaled up
	//_, err = c.createOrUpdateEachSPAForRev(ctx, ro, revScalingUp, false)
	_, err = r.createOrUpdateSPARevUp(ctx, ro, revScalingUp)
	if err != nil {
		return err
	}

	// If spec.StageRevisionStatus is nil, check on if the number of replicas meets the conditions.
	if ro.IsStageInProgress() {
		spa, err := r.stagePodAutoscalerLister.StagePodAutoscalers(ro.Namespace).Get(revScalingUp.RevisionName)
		if err != nil {
			return err
		}
		// spa.IsStageScaleInReady() returns true, as long as both DesireScale and ActualScale are available.
		if !spa.IsStageScaleInReady() || !IsStageScaleUpReady(spa, revScalingUp) {
			// Create the stage pod autoscaler with the new maxScale set to
			// maxScale defined in the revision traffic, because scale up phase is not over, we cannot
			// scale down the old revision.
			// Create or update the stagePodAutoscaler for the revision to be scaled down, eve if the scaling up
			// phase is not over.
			_, err = r.createOrUpdateSPARevDown(ctx, ro, revScalingDown, false)
			if err != nil {
				return err
			}
			return nil
		}

		ro.Status.MarkStageRevisionScaleUpReady()

		// Create the stage pod autoscaler with the new maxScale set to targetScale defined
		// in the revision traffic. Scaling up phase is over, we are able to scale down.
		// Create or update the stagePodAutoscaler for the revision to be scaled down.
		_, err = r.createOrUpdateSPARevDown(ctx, ro, revScalingDown, true)
		if err != nil {
			return err
		}

		spa, err = r.stagePodAutoscalerLister.StagePodAutoscalers(ro.Namespace).Get(revScalingDown.RevisionName)
		if err != nil {
			return err
		}
		if !IsStageScaleDownReady(spa, revScalingDown) {
			return nil
		}

		ro.Status.MarkStageRevisionScaleDownReady()

		// Clean up and set the status of the StageRevision. It means the orchestrator has accomplished this stage.
		stageCleaned := RemoveNonTrafficRev(stageTargetRevisions)
		ro.Status.SetStageRevisionStatus(stageCleaned)
		ro.Status.MarkStageRevisionReady()

		if LastStageComplete(ro.Status.StageRevisionStatus, ro.Spec.TargetRevisions) {
			ro.Status.MarkLastStageRevisionComplete()
			return nil
		}
		ro.Status.MarkLastStageRevisionInComplete()
		return nil
	}

	if ro.IsStageReady() && ro.IsInProgress() && !LastStageComplete(ro.Status.StageRevisionStatus,
		stageTargetRevisions) {
		// Start to move to the next stage.
		ro.Status.LaunchNewStage()
	}

	return nil
}

type updateSPAForRev func(*v1.StagePodAutoscaler, *v1.TargetRevision, bool) *v1.StagePodAutoscaler

func (r *Reconciler) createStagePA(ctx context.Context, ro *v1.RolloutOrchestrator, revision *v1.TargetRevision,
	scaleUpReady bool, fn updateSPAForRev) (*v1.StagePodAutoscaler, error) {
	spa := CreateBaseStagePodAutoscaler(ro, revision)
	spa = fn(spa, revision, scaleUpReady)
	return r.client.ServingV1().StagePodAutoscalers(ro.Namespace).Create(ctx, spa, metav1.CreateOptions{})
}

// RemoveNonTrafficRev removes the redundant TargetRevision from the list of TargetRevisions.
func RemoveNonTrafficRev(ts []v1.TargetRevision) []v1.TargetRevision {
	result := make([]v1.TargetRevision, 0)
	for _, r := range ts {
		if r.Percent != nil && *r.Percent != 0 {
			result = append(result, r)
		}
	}
	return result
}

func targetRevisionEqual(currentStatusRevisions, finalTargetRevisions []v1.TargetRevision) bool {
	if *finalTargetRevisions[0].Percent != 100 {
		return false
	}
	for _, r := range currentStatusRevisions {
		if *r.Percent == 100 && r.RevisionName == finalTargetRevisions[0].RevisionName {
			return true
		}
	}
	return false
}

// UpdateSPAForRevUp update the SPA(StagePodAutoscaler) for the revision scaling up, based on the TargetReplicas
// min & max scales defined in the Knative Service.
func UpdateSPAForRevUp(spa *v1.StagePodAutoscaler, revision *v1.TargetRevision, _ bool) *v1.StagePodAutoscaler {
	// For revisions scaling up, the StageMaxScale is always set to the final MaxScale.
	spa.Spec.StageMaxScale = revision.MaxScale
	min := getMinScale(revision)

	// TargetReplicas is nil, when either the revision has 0% traffic assigned or 100% traffic assigned.
	// It will be either the old revision scales down to 0 or the new revision scale up to replace the old revision.
	// In either case, the spa pick up the StageMinScale and StageMaxScale scales from the knative service.
	if revision.TargetReplicas == nil {
		spa.Spec.StageMinScale = revision.MinScale
		return spa
	}
	targetReplicas := *revision.TargetReplicas
	if targetReplicas < min && *revision.Percent < int64(100) {
		// If the less than of 100% traffic is assigned to this revision or targetReplicas is less than minscale,
		// set StageMinScale directly to targetReplicas.
		spa.Spec.StageMinScale = ptr.Int32(targetReplicas)
	} else {
		// If the 100% traffic is assigned to this revision or targetReplicas is equal to or greater than min,
		// set StageMinScale directly to the final MinScale from the knative service.
		spa.Spec.StageMinScale = revision.MinScale
	}
	return spa
}

// UpdateSPAForRevDown update the SPA(StagePodAutoscaler) for the revision scaling down, based on the TargetReplicas
// min & max scales defined in the Knative Service, if the scaleUpReady is true.
//
// If the scaleUpReady is false, no change to the SPA(StagePodAutoscaler).
func UpdateSPAForRevDown(spa *v1.StagePodAutoscaler, revision *v1.TargetRevision,
	scaleUpReady bool) *v1.StagePodAutoscaler {
	if !scaleUpReady {
		return spa
	}
	min := getMinScale(revision)
	max := getMaxScale(revision)

	if revision.TargetReplicas == nil {
		spa.Spec.StageMinScale = revision.MinScale
		spa.Spec.StageMaxScale = revision.MaxScale
		return spa
	}
	targetReplicas := *revision.TargetReplicas

	// If targetReplicas is equal to or greater than maxScale, StageMinScale and StageMaxScale are set to the final
	// MinScale and MaxScale.
	if targetReplicas >= max {
		spa.Spec.StageMaxScale = revision.MaxScale
		spa.Spec.StageMinScale = revision.MinScale
		return spa
	}

	// If targetReplicas is less than maxScale, StageMaxScale is set to targetReplicas.
	spa.Spec.StageMaxScale = ptr.Int32(targetReplicas)

	// If targetReplicas is less than minScale, StageMinScale is set to targetReplicas.
	if targetReplicas < min {
		spa.Spec.StageMinScale = ptr.Int32(targetReplicas)
		return spa
	}

	// If targetReplicas is larger than or equal to minScale, StageMinScale is set to final MinScale.
	spa.Spec.StageMinScale = revision.MinScale
	return spa
}

// RetrieveRevsUpDown returns the old and the new revision in the list of the TargetRevisions.
func RetrieveRevsUpDown(targetRevs []v1.TargetRevision) (*v1.TargetRevision, *v1.TargetRevision, error) {
	upIndex, downIndex := -1, -1
	for i, rev := range targetRevs {
		if rev.IsRevScalingUp() {
			upIndex = i
		} else if rev.IsRevScalingDown() {
			downIndex = i
		}
	}
	if upIndex == -1 || downIndex == -1 {
		return nil, nil, fmt.Errorf("unable to find the revision to scale up or down in the target revisions %v", targetRevs)
	}
	return &targetRevs[upIndex], &targetRevs[downIndex], nil
}

// LastStageComplete decides whether the last stage of the progressive upgrade is complete or not.
func LastStageComplete(stageRevisionStatus, finalTargetRevs []v1.TargetRevision) bool {
	return equality.Semantic.DeepEqual(stageRevisionStatus, finalTargetRevs) ||
		targetRevisionEqual(stageRevisionStatus, finalTargetRevs)
}

func actualScaleBetweenMinMax(spa *v1.StagePodAutoscaler, min, max int32) bool {
	return *spa.Status.DesiredScale == *spa.Status.ActualScale && *spa.Status.ActualScale >= min && *spa.Status.ActualScale <= max
}

// IsStageScaleUpReady decides whether the scaling up has completed or on the way for the current stage, based
// on the revision and the spa(StagePodAutoscaler).
func IsStageScaleUpReady(spa *v1.StagePodAutoscaler, revision *v1.TargetRevision) bool {
	if spa.Status.DesiredScale == nil || spa.Status.ActualScale == nil {
		return false
	}
	min := getMinScale(revision)
	max := getMaxScale(revision)
	if revision.TargetReplicas == nil {
		// For revision scaling up without TargetReplicas, it means this revision will be assigned 100% of the traffic.
		return actualScaleBetweenMinMax(spa, min, max)
	}

	// There are two modes to scale up and down the replicas of the revisions:
	// 1. No traffic. Knative Service specifies both min and max scales for the revision. We need to control the
	// StageMinScale and StageMaxScale to make sure the replicas increase or decrease. In this case, we need to
	// precisely make sure both DesiredScale and ActualScale are equal to or greater than TargetReplicas to determine
	// scaling up phase is over. TargetReplicas is no larger than minScale, because revision runs at the number of
	// minScale. We need to first scale up the new revision, make sure it run at the correct number, and scale down the
	// old revision. Shifting traffic does not change anything in terms of the number of replicas.
	// 2. Traffic driven. In this case, TargetReplicas is larger than minScale and lower than or equal to maxScale
	// for the revision. We need to change the number of the replicas by shifting the traffic. As long as we know the
	// new revision is on the way of scaling up, we are able to start the scaling down phase as well.
	if *spa.Status.DesiredScale >= *revision.TargetReplicas && *spa.Status.ActualScale >= *revision.TargetReplicas {
		// This is for the first mode.
		return true
	}
	if *spa.Status.DesiredScale >= min && *spa.Status.DesiredScale == *spa.Status.ActualScale {
		// This is for the second mode.
		return true
	}

	return false
}

// IsStageScaleDownReady decides whether the scaling down has completed for the current stage, based
// on the revision and the spa(StagePodAutoscaler).
func IsStageScaleDownReady(spa *v1.StagePodAutoscaler, revision *v1.TargetRevision) bool {
	if spa.Status.DesiredScale == nil || spa.Status.ActualScale == nil {
		return false
	}
	if revision.TargetReplicas == nil {
		// For revision scaling up without TargetReplicas, it means this revision will be assigned 0% of the traffic.
		max := getMaxScale(revision)
		return bothValuesUnderTargetValue(*spa.Status.DesiredScale, *spa.Status.ActualScale, max)
	}
	return bothValuesUnderTargetValue(*spa.Status.DesiredScale, *spa.Status.ActualScale, *revision.TargetReplicas)
}

func bothValuesUnderTargetValue(desire, actual, target int32) bool {
	return desire <= target && actual <= target
}

// CreateBaseStagePodAutoscaler returns the basic spa(StagePodAutoscaler), base
// on the RolloutOrchestrator and the revision.
func CreateBaseStagePodAutoscaler(ro *v1.RolloutOrchestrator, revision *v1.TargetRevision) (spa *v1.StagePodAutoscaler) {
	spa = &v1.StagePodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      revision.RevisionName,
			Namespace: ro.Namespace,
			Labels:    map[string]string{serving.RevisionLabelKey: revision.RevisionName},
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(ro),
			},
		},
		Spec: v1.StagePodAutoscalerSpec{
			StageMinScale: revision.MinScale,
			StageMaxScale: revision.MaxScale,
		},
	}
	return
}

func getMinScale(revision *v1.TargetRevision) (min int32) {
	if revision.MinScale != nil {
		min = *revision.MinScale
	}
	return
}

func getMaxScale(revision *v1.TargetRevision) (max int32) {
	max = int32(math.MaxInt32)
	if revision.MaxScale != nil {
		max = *revision.MaxScale
	}
	return
}