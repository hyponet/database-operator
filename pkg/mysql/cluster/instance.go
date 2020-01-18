package cluster

import (
	"bytes"
	databasev1 "database-operator/pkg/apis/database/v1"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
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
		Command: getMySqlServerCommand(instance),
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
	return "mysql/mysql-server:8.0.19"
}

const runTemplate = `
 # Set baseServerID
 base=1000

 # Finds the replica index from the hostname, and uses this to define
 # a unique server id for this instance.
 index=$(cat /etc/hostname | grep -o '[^-]*$')
 /entrypoint.sh --server-id=$(expr $base + $index) --datadir=/var/lib/mysql --user=mysql --gtid-mode=ON --binlog-checksum=NONE --enforce-gtid-consistency=ON --log-slave-updates=ON --binlog-format=ROW --master-info-repository=TABLE --relay-log-info-repository=TABLE --transaction-write-set-extraction=XXHASH64 --log-error-verbosity=3 %s
`

func getMySqlServerCommand(instance *databasev1.MySQL) []string {
	instanceName := instance.Name
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(`--log-bin="%s-${index}-bin"`, instanceName))
	buf.WriteString(" ")
	buf.WriteString(fmt.Sprintf(`--report-host="%s-${index}.%s"`, instanceName, instanceName))
	buf.WriteString(" ")
	buf.WriteString("--relay-log-recovery=ON")
	buf.WriteString(" ")
	buf.WriteString(fmt.Sprintf(`--relay-log=%s-${index}-relay-bin`, instanceName))
	buf.WriteString(" ")

	// group replication config
	buf.WriteString("--plugin-load=group_replication.so")
	buf.WriteString(" ")
	buf.WriteString(fmt.Sprintf("--group-replication-group-name=%s", instance.UID))
	buf.WriteString(" ")
	buf.WriteString(`--group-replication-local-address="$PODIP:13306"`)
	buf.WriteString(" ")

	address := make([]string, instance.Spec.Members)
	for i := int32(0); i < instance.Spec.Members; i++ {
		address[i] = fmt.Sprintf("%s-%d.%s:13306", instanceName, i, instanceName)
	}

	buf.WriteString("--group-replication-group-seeds=")
	buf.WriteString(strings.Join(address, ","))
	buf.WriteString(" ")
	buf.WriteString("--group-replication-start-on-boot=OFF")
	buf.WriteString(" ")
	buf.WriteString("--group-replication-bootstrap-group=OFF")
	buf.WriteString(" ")
	buf.WriteString("--group-replication-single-primary-mode=ON")
	buf.WriteString(" ")
	buf.WriteString("--group-replication-enforce-update-everywhere-checks=OFF")

	run := fmt.Sprintf(runTemplate, buf.String())
	return []string{
		"/bin/bash",
		"-ecx",
		run,
	}
}
