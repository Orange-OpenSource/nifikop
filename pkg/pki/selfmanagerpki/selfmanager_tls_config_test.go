package selfmanagerpki

import (
	"context"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/errorfactory"
	certutil "github.com/Orange-OpenSource/nifikop/pkg/util/cert"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"testing"
)

func newMockControllerSecret(valid bool) *corev1.Secret {
	secret := &corev1.Secret{}
	secret.Name = "test-controller"
	secret.Namespace = "test-namespace"
	cert, key, _, _ := certutil.GenerateTestCert()
	if valid {
		secret.Data = map[string][]byte{
			corev1.TLSCertKey:       cert,
			corev1.TLSPrivateKeyKey: key,
			v1alpha1.CoreCACertKey:  cert,
		}
	}
	return secret
}

func TestGetControllerTLSConfig(t *testing.T) {
	manager, err := newMock(newMockCluster())
	if err != nil {
		t.Error("Expected no error from New, got:", err)
	}

	// Test good controller secret
	manager.client.Create(context.TODO(), newMockControllerSecret(true))
	if _, err := manager.GetControllerTLSConfig(); err != nil {
		t.Error("Expected no error, got:", err)
	}

	manager, err = newMock(newMockCluster())
	if err != nil {
		t.Error("Expected no error from New, got:", err)
	}

	// Test non-existent controller secret
	if _, err := manager.GetControllerTLSConfig(); err == nil {
		t.Error("Expected error got nil")
	} else if reflect.TypeOf(err) != reflect.TypeOf(errorfactory.ResourceNotReady{}) {
		t.Error("Expected not ready error, got:", reflect.TypeOf(err))
	}

	// Test invalid controller secret
	manager.client.Create(context.TODO(), newMockControllerSecret(false))
	if _, err := manager.GetControllerTLSConfig(); err == nil {
		t.Error("Expected error got nil")
	} else if reflect.TypeOf(err) != reflect.TypeOf(errorfactory.InternalError{}) {
		t.Error("Expected internal error, got:", reflect.TypeOf(err))
	}
}
