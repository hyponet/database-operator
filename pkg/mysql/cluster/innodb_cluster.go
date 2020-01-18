package cluster

import (
	databasev1 "database-operator/pkg/apis/database/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewPodSpec(instance *databasev1.MySQL, dbVolumeName string, usePv bool) corev1.PodSpec {
	containers := []corev1.Container{
		createMySqlServerContainer(instance, dbVolumeName),
		//createMySqlRouteContainer(instance),
	}
	spec := corev1.PodSpec{Containers: containers}

	if !usePv {
		volumes := createEmptyDir(dbVolumeName)
		spec.Volumes = volumes
	}
	return spec
}

func NewService(instance *databasev1.MySQL, labels, selector map[string]string) *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "mysql",
					Port:       int32(MySqlPort),
					TargetPort: intstr.FromInt(MySqlPort),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Selector:  selector,
			ClusterIP: "None",
			Type:      corev1.ServiceTypeClusterIP,
		},
	}

	return svc
}

func createEmptyDir(dbDataVolumeName string) []corev1.Volume {
	return []corev1.Volume{
		{
			Name: dbDataVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}
}
