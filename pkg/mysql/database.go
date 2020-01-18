package mysql

import (
	databasev1 "database-operator/pkg/apis/database/v1"
	"database-operator/pkg/mysql/cluster"
	"database-operator/pkg/utils/dbutil"
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const InnoDBCluster = "Cluster"

type MySQL struct {
	clusterType string
	resource    *databasev1.MySQL
}

func (m MySQL) NewStsForCR() (*appsv1.StatefulSet, *corev1.Service, error) {
	replica := m.resource.Spec.Members
	// start first instance and init cluster
	if !dbutil.IsDatabaseClusterInitialized(m.resource.Status.Conditions) {
		replica = 1
	}

	var volumeClaimTemplates []corev1.PersistentVolumeClaim
	volumeName := defaultVolumeName
	usePv := false
	if m.resource.Spec.VolumeClaimTemplate != nil {
		usePv = true
		dbVolumeTemplate := m.resource.Spec.VolumeClaimTemplate
		volumeName = dbVolumeTemplate.Name
		volumeClaimTemplates = []corev1.PersistentVolumeClaim{*dbVolumeTemplate}
	}

	labels, selector := m.getDefaultLabelAndSelector()

	var podSpec corev1.PodSpec
	var svc *corev1.Service
	switch m.clusterType {
	case InnoDBCluster:
		podSpec = cluster.NewPodSpec(m.resource, volumeName, usePv)
		svc = cluster.NewService(m.resource, labels, selector.MatchLabels)
	default:
		return nil, nil, fmt.Errorf("create mysql cluster failed: unsupported type %s ", m.clusterType)
	}

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.resource.Name,
			Namespace: m.resource.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replica,
			Selector: selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: podSpec,
			},
			VolumeClaimTemplates: volumeClaimTemplates,
			ServiceName:          m.resource.Name,
			PodManagementPolicy:  appsv1.OrderedReadyPodManagement,
			UpdateStrategy:       appsv1.StatefulSetUpdateStrategy{},
		},
	}

	return sts, svc, nil
}

func (m MySQL) UpdateStatus(crt *appsv1.StatefulSet) (*databasev1.MySQLStatus, error) {
	if crt == nil {
		return nil, errors.New("update status error with nil sts")
	}
	status := m.resource.Status.DeepCopy()
	if !dbutil.IsDatabaseClusterControllerCreated(status.Conditions) {
		status.Conditions = dbutil.SetDatabaseClusterControllerCreated(status.Conditions)
	}

	status.Members = m.resource.Spec.Members
	status.ReadyMembers = crt.Status.ReadyReplicas
	status.NotReadyMembers = status.Members - status.ReadyMembers
	if status.NotReadyMembers < 0 {
		status.NotReadyMembers = 0
	}

	if !dbutil.IsDatabaseClusterInitialized(status.Conditions) {
		fmt.Printf("ready member count: %d, isReady: %v\n", status.ReadyMembers, status.ReadyMembers > 0)
		if status.ReadyMembers > 0 {
			status.Conditions = dbutil.SetDatabaseClusterInitialized(status.Conditions)
		} else {
			status.Conditions = dbutil.SetDatabaseClusterInitError(status.Conditions, "No member is ready")
		}
	}

	if status.Members == status.ReadyMembers {
		status.Conditions = dbutil.SetDatabaseClusterReady(status.Conditions)
	} else {
		message := fmt.Sprintf("%d member(s) not ready", m.resource.Spec.Members-status.ReadyMembers)
		status.Conditions = dbutil.SetDatabaseClusterNotReady(status.Conditions, message)
	}

	return status, nil
}

func (m MySQL) getDefaultLabelAndSelector() (map[string]string, *metav1.LabelSelector) {
	labels := make(map[string]string)
	for k, v := range m.resource.Labels {
		labels[k] = v
	}
	labels[InstanceNameKey] = m.resource.Name
	labels[InstanceClusterTypeKey] = m.clusterType

	selector := &metav1.LabelSelector{
		MatchLabels: labels,
	}

	return labels, selector
}

func NewMySqlInstance(mysql *databasev1.MySQL) MySQL {
	return MySQL{
		clusterType: mysql.Spec.Type,
		resource:    mysql,
	}
}
