package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterIdentifierSpec defines the desired state of ClusterIdentifier
type ClusterIdentifierSpec struct {
	Name string `json:"name"`
}

// ClusterIdentifierStatus defines the observed state of ClusterIdentifier
type ClusterIdentifierStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// ClusterIdentifier is the Schema for the clusteridentifiers API
type ClusterIdentifier struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterIdentifierSpec   `json:"spec,omitempty"`
	Status ClusterIdentifierStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterIdentifierList contains a list of ClusterIdentifier
type ClusterIdentifierList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterIdentifier `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterIdentifier{}, &ClusterIdentifierList{})
}
