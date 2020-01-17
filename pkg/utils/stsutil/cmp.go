package stsutil

import (
	appsv1 "k8s.io/api/apps/v1"
	"reflect"
)

func DeepEqual(set1 *appsv1.StatefulSet, set2 *appsv1.StatefulSet) bool {
	if reflect.DeepEqual(set1.Spec, set2.Spec) {
		return true
	}

	// compare pod config

	return true
}
