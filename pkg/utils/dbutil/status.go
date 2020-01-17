package dbutil

import (
	v1 "database-operator/pkg/apis/database/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsDatabaseClusterControllerCreated(conditions []v1.DatabaseCondition) bool {
	for _, c := range conditions {
		if c.Type == v1.ControllerCreated && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func SetDatabaseClusterControllerCreated(conditions []v1.DatabaseCondition) []v1.DatabaseCondition {
	now := metav1.Now()
	for i := range conditions {
		c := conditions[i]
		if c.Type == v1.ControllerCreated {
			if c.Status == v1.ConditionTrue {
				return conditions
			}
			c.Status = v1.ConditionTrue
			c.Time = &now
			return conditions
		}
	}
	conditions = append(conditions, v1.DatabaseCondition{
		Type:   v1.ControllerCreated,
		Status: v1.ConditionTrue,
		Time:   &now,
	})
	return conditions
}

func IsDatabaseClusterInitialized(conditions []v1.DatabaseCondition) bool {
	for _, c := range conditions {
		if c.Type == v1.ClusterInitialized && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func SetDatabaseClusterInitialized(conditions []v1.DatabaseCondition) []v1.DatabaseCondition {
	now := metav1.Now()
	for i := range conditions {
		c := conditions[i]
		if c.Type == v1.ClusterInitialized {
			if c.Status == v1.ConditionTrue {
				return conditions
			}
			c.Status = v1.ConditionTrue
			c.Time = &now
			return conditions
		}
	}
	conditions = append(conditions, v1.DatabaseCondition{
		Type:   v1.ClusterInitialized,
		Status: v1.ConditionTrue,
		Time:   &now,
	})
	return conditions
}

func SetDatabaseClusterInitError(conditions []v1.DatabaseCondition, message string) []v1.DatabaseCondition {
	now := metav1.Now()
	for i := range conditions {
		c := conditions[i]
		if c.Type == v1.ClusterInitialized {
			c.Status = v1.ConditionFalse
			c.Message = message
			c.Time = &now
			return conditions
		}
	}
	conditions = append(conditions, v1.DatabaseCondition{
		Type:    v1.ClusterInitialized,
		Status:  v1.ConditionFalse,
		Message: message,
		Time:    &now,
	})
	return conditions
}

func IsDatabaseClusterReady(conditions []v1.DatabaseCondition) bool {
	for _, c := range conditions {
		if c.Type == v1.ClusterReady && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func SetDatabaseClusterReady(conditions []v1.DatabaseCondition) []v1.DatabaseCondition {
	now := metav1.Now()
	for i := range conditions {
		c := conditions[i]
		if c.Type == v1.ClusterReady {
			if c.Status == v1.ConditionTrue {
				return conditions
			}
			c.Status = v1.ConditionTrue
			c.Time = &now
			return conditions
		}
	}
	conditions = append(conditions, v1.DatabaseCondition{
		Type:   v1.ClusterReady,
		Status: v1.ConditionTrue,
		Time:   &now,
	})
	return conditions
}

func SetDatabaseClusterNotReady(conditions []v1.DatabaseCondition, message string) []v1.DatabaseCondition {
	now := metav1.Now()
	for i := range conditions {
		c := conditions[i]
		if c.Type == v1.ClusterReady {
			c.Status = v1.ConditionFalse
			c.Message = message
			c.Time = &now
			return conditions
		}
	}
	conditions = append(conditions, v1.DatabaseCondition{
		Type:    v1.ClusterReady,
		Status:  v1.ConditionFalse,
		Message: message,
		Time:    &now,
	})
	return conditions
}
