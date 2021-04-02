package selfmanagerpki

import (
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *SelfManager) ReconcilePKI(ctx context.Context, logger logr.Logger, scheme *runtime.Scheme, externalHostnames []string) error {

	logger.Info("Reconciling cert-manager PKI")

	// TODO Generate all certs
	// TODO Setup all secrets from certs
}

// TODO
func (s SelfManager) FinalizePKI(ctx context.Context, logger logr.Logger) error {
	panic("implement me")
}
