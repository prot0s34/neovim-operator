/*
Copyright 2024.

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

// NeovimSpec defines the desired state of Neovim
type NeovimSpec struct {
	// Username is the name of the user for whom this Neovim instance is created.
	Username string `json:"username,omitempty"`

	// Version is the version of Neovim to be deployed.
	Version string `json:"version,omitempty"`
}

// NeovimStatus defines the observed state of Neovim
type NeovimStatus struct {
	// Ready indicates whether the Neovim instance is ready to be used.
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// Neovim is the Schema for the neovims API
type Neovim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NeovimSpec   `json:"spec,omitempty"`
	Status NeovimStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type NeovimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Neovim `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Neovim{}, &NeovimList{})
}
