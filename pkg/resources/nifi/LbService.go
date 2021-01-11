// Copyright 2020 Orange SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.package apis

package nifi

import (
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/resources/templates"
	corev1 "k8s.io/api/core/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO: To remove ? Or to redo
func (r *Reconciler) externalServices() []runtimeClient.Object {

	var services []runtimeClient.Object
	for _, eService := range r.NifiCluster.Spec.ExternalServices {

		var listeners []v1alpha1.InternalListenerConfig

		for _, port := range eService.Spec.PortConfigs {
			for  _, iListener := range r.NifiCluster.Spec.ListenersConfig.InternalListeners {
				if port.InternalListenerName == iListener.Name {
					listeners = append(listeners, iListener)
				}
			}
		}

		usedPorts := generateServicePortForInternalListeners(listeners)
		services = append(services, &corev1.Service{
			ObjectMeta: templates.ObjectMetaWithAnnotations(r.NifiCluster.Name, LabelsForNifi(r.NifiCluster.Name),
				r.NifiCluster.Spec.Service.Annotations, r.NifiCluster),
			Spec: corev1.ServiceSpec{
				Type:            corev1.ServiceTypeLoadBalancer,
				SessionAffinity: corev1.ServiceAffinityClientIP,
				Selector:        LabelsForNifi(r.NifiCluster.Name),
				Ports:           usedPorts,
			},
		})
	}

	return services
}
