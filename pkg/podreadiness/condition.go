package podreadiness

import (
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RTEConditionType is a valid value for PodCondition.Type
type RTEConditionType string

// These are valid conditions of RTE pod.
const (
	// ResourcesScanned means that resources scanned successfully.
	ResourcesScanned RTEConditionType = "ResourcesScanned"
	// TopologyUpdated means that noderesourcetopology objects updated successfully.
	TopologyUpdated RTEConditionType = "TopologyUpdated"
)

func NewConditionTemplate(condType RTEConditionType, status v1.ConditionStatus) (condition v1.PodCondition) {
	return v1.PodCondition{
		Type:               v1.PodConditionType(condType),
		Status:             status,
		LastTransitionTime: metav1.Time{Time: time.Now()},
	}
}
