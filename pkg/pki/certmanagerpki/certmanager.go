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

package certmanagerpki

import (
	"github.com/Orange-OpenSource/nifikop/pkg/apis/nifi/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CertManager interface {
	pki.Manager
}

// certManager implements a PKIManager using cert-manager as the backend
type certManager struct {
	client  client.Client
	cluster *v1alpha1.NifiCluster
}

func New(client client.Client, cluster *v1alpha1.NifiCluster) CertManager {
	return &certManager{client: client, cluster: cluster}
}
