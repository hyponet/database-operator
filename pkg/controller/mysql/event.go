package mysql

import (
	databasev1 "database-operator/pkg/apis/database/v1"
	appsv1 "k8s.io/api/apps/v1"
)

type eventType string

const (
	stsCreated      eventType = "StatefulSetCreated"
	stsScaled                 = "StatefulSetScaled"
	ClusterReady              = "ClusterReady"
	ClusterNotReady           = "ClusterNotReady"
)

func (r *ReconcileMySQL) recordMySqlInstanceEvent(instance *databasev1.MySQL, et eventType, message string) {
	r.eventRecord.Event(instance, string(et), message, message)
}

func (r *ReconcileMySQL) recordStsEvent(sts *appsv1.StatefulSet, et eventType, message string) {
	r.eventRecord.Event(sts, string(et), message, message)
}
