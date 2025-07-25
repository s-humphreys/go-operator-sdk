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
	ConditionFailed
)

const (
	ResourcesReady ConditionReason = iota
	ProgressingResources
	ResourcesFailed
)

var ConditionTypeMap = map[ConditionType]string{
	ConditionReady:       "Ready",
	ConditionProgressing: "Progressing",
	ConditionFailed:      "Failed",
}

var conditionReasonMap = map[ConditionReason]Condition{
	ResourcesReady: {
		conditionType: ConditionReady,
		reason:        "ResourcesReady",
		message:       "Recources all ready and in desired state",
	},
	ProgressingResources: {
		conditionType: ConditionProgressing,
		reason:        "ProgressingResources",
		message:       "Progressing resources to sync with the desired state",
	},
	ResourcesFailed: {
		conditionType: ConditionFailed,
		reason:        "ResourcesFailed",
		message:       "Failed to provision resources",
	},
}

// Creates a condition status for the CRD using a provided ConditionReason.
// This allows for a consistent way to create conditions based on predefined reasons
// across CRDs.
func NewStatusCondition(reason ConditionReason) metav1.Condition {
	return metav1.Condition{
		Type:    ConditionTypeMap[conditionReasonMap[reason].conditionType],
		Status:  metav1.ConditionTrue,
		Reason:  conditionReasonMap[reason].reason,
		Message: conditionReasonMap[reason].message,
	}
}
