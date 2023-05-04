/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Overrides struct {
	ContainerPorts []int32  `json:"ports,omitempty"`
	Command        []string `json:"command,omitempty"`
	Args           []string `json:"args,omitempty"`
}

// KontrollerSpec defines the desired state of Kontroller.
type KontrollerSpec struct {
	Image     string `json:"image,omitempty"`
	Overrides `json:",inline"`
}

// KontrollerStatus defines the observed state of Kontroller.
type KontrollerStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Kontroller is the Schema for the Kontrollers API.
type Kontroller struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KontrollerSpec   `json:"spec,omitempty"`
	Status KontrollerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KontrollerList contains a list of Kontroller.
type KontrollerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kontroller `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kontroller{}, &KontrollerList{})
}
