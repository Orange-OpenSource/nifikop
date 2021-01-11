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
	"fmt"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/resources/templates"
	"github.com/Orange-OpenSource/nifikop/pkg/util"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func (r *Reconciler) service(id int32, log logr.Logger) runtimeClient.Object {

	usedPorts := generateServicePortForInternalListeners(r.NifiCluster.Spec.ListenersConfig.InternalListeners)

	return &corev1.Service{
		ObjectMeta: templates.ObjectMeta(fmt.Sprintf("%s-%d", r.NifiCluster.Name, id),
			util.MergeLabels(
				LabelsForNifi(r.NifiCluster.Name),
				map[string]string{"nodeId": fmt.Sprintf("%d", id)},
			),
			r.NifiCluster),
		Spec: corev1.ServiceSpec{
			Type:            corev1.ServiceTypeClusterIP,
			SessionAffinity: corev1.ServiceAffinityNone,
			Selector:        util.MergeLabels(LabelsForNifi(r.NifiCluster.Name), map[string]string{"nodeId": fmt.Sprintf("%d", id)}),
			Ports:           usedPorts,
		},
	}
}

//
func generateServicePortForInternalListeners(listeners []v1alpha1.InternalListenerConfig) []corev1.ServicePort {
	var usedPorts []corev1.ServicePort

	for _, iListeners := range listeners {
		usedPorts = append(usedPorts, corev1.ServicePort{
			Name:       strings.ReplaceAll(iListeners.Name, "_", ""),
			Port:       iListeners.ContainerPort,
			TargetPort: intstr.FromInt(int(iListeners.ContainerPort)),
			Protocol:   corev1.ProtocolTCP,
		})
	}

	return usedPorts
}
