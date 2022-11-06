// Package v1 contains API Schema definitions for the opslevel v1 API group
//+kubebuilder:object:generate=true
//+groupName=opslevel.com
package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "opslevel.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
)

func AddToScheme(scheme *runtime.Scheme) {
	SchemeBuilder.Register(&ClusterIdentifier{}, &ClusterIdentifierList{})
	SchemeBuilder.AddToScheme(scheme)
}
