/*
Copyright 2022.

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

package v1beta1

import (
	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DesignateWorkerSpec defines the desired state of DesignateWorker
type DesignateWorkerSpec struct {
  // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
  // Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=designate
	// ServiceUser - optional username used for this service to register in designate
	ServiceUser string `json:"serviceUser"`

	// +kubebuilder:validation:Required
	// Designate Container Image URL
	ContainerImage string `json:"containerImage"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Maximum=32
	// +kubebuilder:validation:Minimum=0
	// Replicas of designate Worker to run
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Optional
	// DatabaseHostname - Cinder Database Hostname
	DatabaseHostname string `json:"databaseHostname,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=designate
	// DatabaseUser - optional username used for designate DB, defaults to designate
	// TODO: -> implement needs work in mariadb-operator, right now only designate
	DatabaseUser string `json:"databaseUser"`

	// +kubebuilder:validation:Required
	// Secret containing OpenStack password information for designate DesignateDatabasePassword, AdminPassword
	Secret string `json:"secret,omitempty"`

	// +kubebuilder:validation:Optional
	// Secret containing RabbitMq transport URL
	TransportURLSecret string `json:"transportURLSecret,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={database: DesignateDatabasePassword, service: DesignatePassword}
	// PasswordSelectors - Selectors to identify the DB and AdminUser password from the Secret
	PasswordSelectors PasswordSelector `json:"passwordSelectors,omitempty"`

	// +kubebuilder:validation:Optional
	// NodeSelector to target subset of worker nodes running this service
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	// ConfigOverwrite - interface to overwrite default config files like e.g. logging.conf or policy.json.
	// But can also be used to add additional files. Those get added to the service config dir in /etc/<service> .
	// TODO: -> implement
	DefaultConfigOverwrite map[string]string `json:"defaultConfigOverwrite,omitempty"`

	// +kubebuilder:validation:Optional
	// Resources - Compute Resources required by this service (Limits/Requests).
	// https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// DesignateWorkerStatus defines the observed state of DesignateWorker
type DesignateWorkerStatus struct {
	// ReadyCount of designate central instances
	ReadyCount int32 `json:"readyCount,omitempty"`

	// Map of hashes to track e.g. job status
	Hash map[string]string `json:"hash,omitempty"`

	// Conditions
	Conditions condition.Conditions `json:"conditions,omitempty" optional:"true"`

	// Designate Database Hostname
	DatabaseHostname string `json:"databaseHostname,omitempty"`

	// ServiceID - the ID of the registered service in keystone
	ServiceID string `json:"serviceID,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[0].status",description="Status"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[0].message",description="Message"

// DesignateWorker is the Schema for the designatecentrals
type DesignateWorker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DesignateWorkerSpec   `json:"spec,omitempty"`
	Status DesignateWorkerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DesignateWorkerList contains a list of DesignateWorker
type DesignateWorkerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DesignateWorker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DesignateWorker{}, &DesignateWorkerList{})
}

// IsReady - returns true if service is ready to serve requests
func (instance CinderAPI) IsReady() bool {
	return instance.Status.ReadyCount == instance.Spec.Replicas
}