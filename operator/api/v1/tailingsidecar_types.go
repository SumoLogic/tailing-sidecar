/*
Copyright 2021.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type SidecarConfig struct {
	File   string `json:"file,omitempty"`
	Volume string `json:"volume,omitempty"`
}

// TailingSidecarSpec defines the desired state of TailingSidecar
type TailingSidecarSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Configs contains configurations for tailing sidecars
	// map key indicates name of configuration which can be used in annotation
	Configs map[string]SidecarConfig `json:"configs,omitempty"`
}

// TailingSidecarStatus defines the observed state of TailingSidecar
type TailingSidecarStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TailingSidecar is the Schema for the tailingsidecars API
type TailingSidecar struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TailingSidecarSpec   `json:"spec,omitempty"`
	Status TailingSidecarStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TailingSidecarList contains a list of TailingSidecar
type TailingSidecarList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TailingSidecar `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TailingSidecar{}, &TailingSidecarList{})
}
