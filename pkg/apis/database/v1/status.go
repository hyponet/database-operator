package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type DatabaseConditionType string

const (
	ControllerCreated  DatabaseConditionType = "ControllerCreated"
	ClusterInitialized                       = "ClusterInitialized"
	ClusterReady                             = "ClusterReady"
)

type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse                   = "False"
	ConditionUnknown                 = "Unknown"
)

type DatabaseCondition struct {
	Type    DatabaseConditionType `json:"type"`
	Status  ConditionStatus       `json:"status"`
	Message string                `json:"message,omitempty"`
	Time    *metav1.Time          `json:"time"`
}
