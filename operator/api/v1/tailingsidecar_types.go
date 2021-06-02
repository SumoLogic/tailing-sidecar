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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type SidecarSpec struct {
	// Annotations defines tailing sidecar container annotations.
	Annotations map[string]string `json:"annotations,omitempty"`

	// Path defines path to a file containing logs to tail within a tailing sidecar container.
	Path string `json:"path,omitempty"`

	// VolumeMount describes a mounting of a volume within a tailing sidecar container.
	VolumeMount corev1.VolumeMount `json:"volumeMount,omitempty"`
}

// TailingSidecarConfigSpec defines the desired state of TailingSidecarConfig
type TailingSidecarConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// AnnotationsPrefix defines prefix for per container annotations.
	AnnotationsPrefix string `json:"annotationsPrefix,omitempty"`

	// SidecarSpecs defines specifications for tailing sidecar containers,
	// map key indicates name of tailing sidecar container
	SidecarSpecs map[string]SidecarSpec `json:"configs,omitempty"`

	// PodSelector selects Pods to which this tailing sidecar configuration applies.
	PodSelector *metav1.LabelSelector `json:"podSelector,omitempty"`
}

// TailingSidecarConfigStatus defines the observed state of TailingSidecarConfig
type TailingSidecarConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TailingSidecarConfig is the Schema for the tailingsidecars API
type TailingSidecarConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TailingSidecarConfigSpec   `json:"spec,omitempty"`
	Status TailingSidecarConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TailingSidecarConfigList contains a list of TailingSidecarConfig
type TailingSidecarConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TailingSidecarConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TailingSidecarConfig{}, &TailingSidecarConfigList{})
}
