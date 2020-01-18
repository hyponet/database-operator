package mysql

import (
	databasev1 "database-operator/pkg/apis/database/v1"
	appsv1 "k8s.io/api/apps/v1"
)

type eventType string

const (
	EventTypeNormal  eventType = "Normal"
	EventTypeWarning           = "Warning"
)

type reasonType string

const (
	StsCreated         reasonType = "StatefulSetCreated"
	StsCreateFailed               = "StatefulSetCreateFailed"
	StsScaled                     = "StatefulSetScaled"
	ClusterInitialized            = "ClusterInitialized"
	ClusterReady                  = "ClusterReady"
	ClusterNotReady               = "ClusterNotReady"
)

func (r *ReconcileMySQL) recordMySqlInstanceEvent(instance *databasev1.MySQL, et eventType, reason reasonType, message string) {
	r.eventRecord.Event(instance, string(et), string(reason), message)
}

func (r *ReconcileMySQL) recordStsEvent(sts *appsv1.StatefulSet, et eventType, reason reasonType, message string) {
	r.eventRecord.Event(sts, string(et), string(reason), message)
}
