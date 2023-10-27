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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PtemplateSpec defines the desired state of Ptemplate
type PtemplateSpec struct {
	ImagePullSecret *v1.LocalObjectReference `json:"imagePullSecret,omitempty"`
	DefaultMsg      string                   `json:"defaultmsg,omitempty"`
	MaxDev          int64                    `json:"maxdev,omitempty"`
	DevicePlugin    string                   `json:"deviceplugin,omitempty"`
	Selector        map[string]string        `json:"selector"`
	ConsumerImage   string                   `json:"consumer,omitempty"`
	RequiredDev     int64                    `json:"requiredDevices,omitempty"`
}

// PtemplateStatus defines the observed state of Ptemplate
type PtemplateStatus struct {
	Consumers []string `json:"consumers"`
	//Consumers []string `json:"consumers" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=consumers"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Ptemplate is the Schema for the ptemplates API
type Ptemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PtemplateSpec   `json:"spec,omitempty"`
	Status PtemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PtemplateList contains a list of Ptemplate
type PtemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ptemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Ptemplate{}, &PtemplateList{})
}
