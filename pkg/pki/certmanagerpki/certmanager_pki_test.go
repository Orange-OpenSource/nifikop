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
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/errorfactory"
	certutil "github.com/Orange-OpenSource/nifikop/pkg/util/cert"
	pkicommon "github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("testing")

func newNodeServerSecret(nodeId int32) *corev1.Secret {
	secret := &corev1.Secret{}
	secret.Name = fmt.Sprintf(pkicommon.NodeServerCertTemplate, "test", nodeId)
	secret.Namespace = "test-namespace"
	cert, key, _, _ := certutil.GenerateTestCert()
	secret.Data = map[string][]byte{
		corev1.TLSCertKey:       cert,
		corev1.TLSPrivateKeyKey: key,
		v1alpha1.CoreCACertKey:  cert,
	}
	return secret
}

func newControllerSecret() *corev1.Secret {
	secret := &corev1.Secret{}
	secret.Name = fmt.Sprintf(pkicommon.NodeControllerTemplate, "test")
	secret.Namespace = "test-namespace"
	cert, key, _, _ := certutil.GenerateTestCert()
	secret.Data = map[string][]byte{
		corev1.TLSCertKey:       cert,
		corev1.TLSPrivateKeyKey: key,
		v1alpha1.CoreCACertKey:  cert,
	}
	return secret
}

func newCASecret() *corev1.Secret {
	secret := &corev1.Secret{}
	secret.Name = fmt.Sprintf(pkicommon.NodeCACertTemplate, "test")
	secret.Namespace = "cert-manager"
	cert, key, _, _ := certutil.GenerateTestCert()
	secret.Data = map[string][]byte{
		corev1.TLSCertKey:       cert,
		corev1.TLSPrivateKeyKey: key,
		v1alpha1.CoreCACertKey:  cert,
	}
	return secret
}

func newPreCreatedSecret() *corev1.Secret {
	secret := &corev1.Secret{}
	secret.Name = "test-controller"
	secret.Namespace = "test-namespace"
	cert, key, _, _ := certutil.GenerateTestCert()
	secret.Data = map[string][]byte{
		v1alpha1.CAPrivateKeyKey: key,
		v1alpha1.CACertKey:       cert,
	}
	return secret
}

func TestFinalizePKI(t *testing.T) {
	manager := newMock(newMockCluster())

	if err := manager.FinalizePKI(context.Background(), log); err != nil {
		t.Error("Expected no error on finalize, got:", err)
	}
}

func TestReconcilePKI(t *testing.T) {
	cluster := newMockCluster()
	manager := newMock(cluster)
	ctx := context.Background()

	for _, node := range cluster.Spec.Nodes {
		manager.client.Create(ctx, newNodeServerSecret(node.Id))
		if err := manager.ReconcilePKI(ctx, log, scheme.Scheme, []string{}); err != nil {
			if reflect.TypeOf(err) != reflect.TypeOf(errorfactory.ResourceNotReady{}) {
				t.Error("Expected not ready error, got:", reflect.TypeOf(err))
			}
		}
	}

	manager.client.Create(ctx, newControllerSecret())
	if err := manager.ReconcilePKI(ctx, log, scheme.Scheme, []string{}); err != nil {
		if reflect.TypeOf(err) != reflect.TypeOf(errorfactory.ResourceNotReady{}) {
			t.Error("Expected not ready error, got:", reflect.TypeOf(err))
		}
	}

	manager.client.Create(ctx, newCASecret())
	if err := manager.ReconcilePKI(ctx, log, scheme.Scheme, []string{}); err != nil {
		t.Error("Expected successful reconcile, got:", err)
	}

	cluster.Spec.ListenersConfig.SSLSecrets.Create = false
	manager = newMock(cluster)
	if err := manager.ReconcilePKI(ctx, log, scheme.Scheme, []string{}); err == nil {
		t.Error("Expected error got nil")
	} else if reflect.TypeOf(err) != reflect.TypeOf(errorfactory.ResourceNotReady{}) {
		t.Error("Expected not ready error, got:", reflect.TypeOf(err))
	}
	manager.client.Create(ctx, newPreCreatedSecret())
	if err := manager.ReconcilePKI(ctx, log, scheme.Scheme, []string{}); err != nil {
		t.Error("Expected successful reconcile, got:", err)
	}
}
