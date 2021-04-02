package selfmanagerpki

import (
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s selfManager) ReconcilePKI(ctx context.Context, logger logr.Logger, scheme *runtime.Scheme, externalHostnames []string) error {
	panic("implement me")
	// TODO Generate all certs
	// TODO Setup all secrets from certs
}

// TODO
func (s selfManager) FinalizePKI(ctx context.Context, logger logr.Logger) error {
	panic("implement me")
}
