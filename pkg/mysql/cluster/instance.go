package cluster

import (
	databasev1 "database-operator/pkg/apis/database/v1"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	ProbeInitialDelaySeconds int32 = 60
	ProbeTimeoutSeconds            = 10
	ProbePeriodSeconds             = 30
	MySqlPort                int   = 3306
)

func createMySqlServerContainer(instance *databasev1.MySQL, dbVolumeName string) corev1.Container {
	container := corev1.Container{
		Name:    "mysql-server",
		Image:   getMySqlServerImage(),
		Command: getMySqlServerCommand(instance.Name),
	}

	// add env config
	envs := []corev1.EnvVar{
		{Name: "MYSQL_ROOT_HOST", Value: "%"},
		{Name: "MYSQL_LOG_CONSOLE", Value: "true"},
	}
	auth := instance.Spec.Auth
	envs = append(envs, corev1.EnvVar{
		Name: "MYSQL_ROOT_PASSWORD",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: auth.RootPasswordSecret.Name},
				Key:                  "password",
			},
		},
	})
	envs = append(envs, corev1.EnvVar{
		Name: "PODIP",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "status.podIP"},
		},
	})
	container.Env = envs

	// data volume
	dataVolume := corev1.VolumeMount{
		Name:      dbVolumeName,
		MountPath: "/var/lib/mysql",
		SubPath:   "mysql",
	}
	container.VolumeMounts = []corev1.VolumeMount{dataVolume}

	// init LivenessProbe & ReadinessProbe
	probe := &corev1.Probe{
		Handler: corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(MySqlPort),
			},
		},
		InitialDelaySeconds: ProbeInitialDelaySeconds,
		TimeoutSeconds:      ProbeTimeoutSeconds,
		PeriodSeconds:       ProbePeriodSeconds,
	}
	container.LivenessProbe = probe
	container.ReadinessProbe = probe

	// add cluster init hook
	//lifeCycle := &corev1.Lifecycle{
	//	PostStart: &corev1.Handler{
	//		Exec: &corev1.ExecAction{Command: []string{"mysqlsh", "--no-password", "--py", "--file=/script/init_cluster.py"}},
	//	},
	//}
	//container.Lifecycle = lifeCycle

	return container
}

func getMySqlServerImage() string {
	return "mysql/mysql-server:8.0"
}

const runTemplate = `
 # Set baseServerID
 base=1000

 # Finds the replica index from the hostname, and uses this to define
 # a unique server id for this instance.
 index=$(cat /etc/hostname | grep -o '[^-]*$')
 /entrypoint.sh --server_id=$(expr $base + $index) --datadir=/var/lib/mysql --user=mysql --gtid_mode=ON --log-bin --binlog_checksum=NONE --enforce_gtid_consistency=ON --log-slave-updates=ON --binlog-format=ROW --master-info-repository=TABLE --relay-log-info-repository=TABLE --transaction-write-set-extraction=XXHASH64 --relay-log=%s-${index}-relay-bin --report-host="%s-${index}.%s" --log-error-verbosity=3
`

func getMySqlServerCommand(instanceName string) []string {
	run := fmt.Sprintf(runTemplate, instanceName, instanceName, instanceName)
	return []string{
		"/bin/bash",
		"-ecx",
		run,
	}
}
