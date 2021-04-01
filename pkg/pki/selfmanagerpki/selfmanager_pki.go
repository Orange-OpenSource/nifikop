package selfmanagerpki

import (
	"context"
	"crypto/tls"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s selfManager) ReconcilePKI(ctx context.Context, logger logr.Logger, scheme *runtime.Scheme, externalHostnames []string) error {
	panic("implement me")
	// TODO Generate all certs
	// TODO Setup all secrets from certs
}

// TODO nil or ?
func (s selfManager) FinalizePKI(ctx context.Context, logger logr.Logger) error {
	panic("implement me")
}

func (s selfManager) ReconcileUserCertificate(ctx context.Context, user *v1alpha1.NifiUser, scheme *runtime.Scheme) (*pki.UserCertificate, error) {
	panic("implement me")
}

func (s selfManager) FinalizeUserCertificate(ctx context.Context, user *v1alpha1.NifiUser) error {
	panic("implement me")
}

func (s selfManager) GetControllerTLSConfig() (*tls.Config, error) {
	panic("implement me")
}
