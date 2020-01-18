package dbutil

import (
	v1 "database-operator/pkg/apis/database/v1"
	"reflect"
	"testing"
)

func TestGetDatabaseCondition(t *testing.T) {
	conditions := make([]v1.DatabaseCondition, 0)
	first := v1.DatabaseCondition{
		Type:   v1.ClusterReady,
		Status: v1.ConditionTrue,
	}
	conditions = append(conditions, first)

	type args struct {
		dt         v1.DatabaseConditionType
		conditions []v1.DatabaseCondition
	}
	tests := []struct {
		name string
		args args
		want *v1.DatabaseCondition
	}{
		{
			name: "getCondition",
			args: args{dt: v1.ClusterReady, conditions: conditions},
			want: &conditions[0],
		},
		{
			name: "getNil",
			args: args{dt: v1.ClusterInitialized, conditions: conditions},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDatabaseCondition(tt.args.dt, tt.args.conditions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDatabaseCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}
