/*
Copyright 2025.

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

// SamtestSpec defines the desired state of Samtest.
type SamtestSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default:=false
	Suspend bool `json:"suspend,omitempty"`

	// +kubebuilder:validation:Pattern=`^(.*):(.*)$`
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// +kubebuilder:validation:Required
	Repliacas int `json:"replicas"`
}

// SamtestStatus defines the observed state of Samtest.
type SamtestStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Samtest is the Schema for the samtests API.
type Samtest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SamtestSpec   `json:"spec,omitempty"`
	Status SamtestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SamtestList contains a list of Samtest.
type SamtestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Samtest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Samtest{}, &SamtestList{})
}
