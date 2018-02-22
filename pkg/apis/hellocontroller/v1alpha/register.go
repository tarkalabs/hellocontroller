package v1alpha

import (
	"github.com/tarkalabs/hellocontroller/pkg/apis/hellocontroller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemaGroupVersion = schema.SchemaGroupVersion{Group: hellocontroller.GroupName, Version: "v1alpha"}

var (
	SchemeBuilder = runtime.NewSchemeBuilder()
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion, Database{}, DatabaseList{})
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}
