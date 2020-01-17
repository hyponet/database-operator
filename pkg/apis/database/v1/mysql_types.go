package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

// MySQLSpec defines the desired state of MySQL
type MySQLSpec struct {
	// +kubebuilder:validation:Enum=Cluster
	Type string `json:"type"`

	// +kubebuilder:validation:Minimum=1
	Members int `json:"members"`

	Auth DatabaseAuth `json:"auth"`

	// TODO support more version and backup
	//Version    string `json:"version"`
	//BackupCron string `json:"backup_cron"`

	VolumeClaimTemplate v1.PersistentVolumeClaim `json:"volume_claim_template,omitempty"`
}

// MySQLStatus defines the observed state of MySQL
type MySQLStatus struct {
	Members         int                 `json:"members"`
	ReadyMembers    int                 `json:"ready_members"`
	NotReadyMembers int                 `json:"not_ready_members"`
	Conditions      []DatabaseCondition `json:"conditions,omitempty"`
	StartTime       *metav1.Time        `json:"start_time,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQL is the Schema for the mysqls API
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.members,statuspath=.status.members
// +kubebuilder:resource:path=mysqls,scope=Namespaced
type MySQL struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MySQLSpec   `json:"spec,omitempty"`
	Status MySQLStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQLList contains a list of MySQL
type MySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MySQL `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MySQL{}, &MySQLList{})
}
