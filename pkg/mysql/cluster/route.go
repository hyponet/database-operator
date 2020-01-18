package cluster

import (
	databasev1 "database-operator/pkg/apis/database/v1"
	corev1 "k8s.io/api/core/v1"
)

const defaultWaitForDbReadyTime = 60

func createMySqlRouteContainer(instance *databasev1.MySQL) corev1.Container {
	envs := []corev1.EnvVar{
		{Name: "MYSQL_HOST", Value: "localhost"},
		{Name: "MYSQL_PORT", Value: "3306"},
		{Name: "MYSQL_USER", Value: "root"},
	}
	auth := instance.Spec.Auth
	envs = append(envs, corev1.EnvVar{
		Name: "MYSQL_PASSWORD",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: auth.RootPasswordSecret.Name},
				Key:                  "password",
			},
		},
	})

	return corev1.Container{
		Name:      "mysql-route",
		Image:     getMySqlRouteImage(),
		Env:       envs,
	}
}

func getMySqlRouteImage() string {
	return "mysql/mysql-router"
}
