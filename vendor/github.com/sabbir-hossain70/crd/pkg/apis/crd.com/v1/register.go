package v1

import (
	crd_com "github.com/sabbir-hossain70/crd/pkg/apis/crd.com"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{Group: crd_com.GroupName, Version: "v1"}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder      = runtime.NewSchemeBuilder(AddKnownTypes)
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

func init() {
	localSchemeBuilder.Register(AddKnownTypes)
}

func AddKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Sabbir{},
		&SabbirList{},
	)
	scheme.AddKnownTypes(SchemeGroupVersion, &metav1.Status{})
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
