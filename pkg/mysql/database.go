package mysql

import (
	v1 "database-operator/pkg/apis/database/v1"
	"database-operator/pkg/utils/dbutil"
	appsv1 "k8s.io/api/apps/v1"
)

type MySQL struct {
	clusterType string
	resource    *v1.MySQL
}

func (m MySQL) NewStsForCR() *appsv1.StatefulSet {
	return nil
}

func (m MySQL) UpdateStatus(crt *appsv1.StatefulSet) (*v1.MySQLStatus, error) {
	status := m.resource.Status
	if !dbutil.IsDatabaseClusterControllerCreated(status.Conditions) {
		status.Conditions = dbutil.SetDatabaseClusterControllerCreated(status.Conditions)
	}

	return nil, nil
}

func NewMySqlInstance(mysql *v1.MySQL) MySQL {
	return MySQL{
		clusterType: mysql.Spec.Type,
		resource:    mysql,
	}
}
