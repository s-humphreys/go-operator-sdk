package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConditionType int
type ConditionReason int

type Condition struct {
	conditionType ConditionType
	reason        string
	message       string
}

const (
	ConditionReady ConditionType = iota
	ConditionProgressing
)

const (
	ResourcesCreated ConditionReason = iota
	ResourcesReconciled
	ResourcesOutOfSync
)

var ConditionTypeMap = map[ConditionType]string{
	ConditionReady:       "Ready",
	ConditionProgressing: "Progressing",
}

var conditionReasonMap = map[ConditionReason]Condition{
	ResourcesCreated: {
		conditionType: ConditionReady,
		reason:        "ResourcesCreated",
		message:       "Resources have been created successfully",
	},
	ResourcesReconciled: {
		conditionType: ConditionReady,
		reason:        "ResourcesReconciled",
		message:       "Recources have been reconciled successfully",
	},
	ResourcesOutOfSync: {
		conditionType: ConditionProgressing,
		reason:        "ResourcesOutOfSync",
		message:       "At least one resource is out of sync with the desired state",
	},
}

// Creates a condition status for the CRD using a provided ConditionReason.
// This allows for a consistent way to create conditions based on predefined reasons
// across CRDs.
func NewStatusCondition(reason ConditionReason) *metav1.Condition {
	return &metav1.Condition{
		Type:    ConditionTypeMap[conditionReasonMap[reason].conditionType],
		Status:  metav1.ConditionTrue,
		Reason:  conditionReasonMap[reason].reason,
		Message: conditionReasonMap[reason].message,
	}
}
