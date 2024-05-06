// +groupName=crd.com
package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ContainerSpec struct {
	Image string   `json:"image,omitempty"`
	Ports []string `json:"ports,omitempty"`
}

type SabbirSpec struct {
	Name      string        `json:"name,omitempty"`
	Replicas  *int32        `json:"replicas,omitempty"`
	Container ContainerSpec `json:"container,omitempty"`
}

type SabbirStatus struct {
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Sabbir struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SabbirSpec   `json:"spec,omitempty"`
	Status            SabbirStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SabbirList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Sabbir `json:"items"`
}
