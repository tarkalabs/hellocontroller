package v1alpha

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              DatabaseSpec   `json:"spec"`
	Status            DatabaseStatus `json:"status,omitempty"`
}

type DatabaseSpec struct {
	DatabaseName string `json:"databaseName"`
	SecretName   string `json:"secretName"`
}
type DatabaseStatus struct {
	Status string `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Database `json:"items"`
}
