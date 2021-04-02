package selfmanagerpki

import (
	"context"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	pkicommon "github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *SelfManager) ReconcilePKI(ctx context.Context, logger logr.Logger, scheme *runtime.Scheme, externalHostnames []string) error {
	logger.Info("Reconciling selfmanager PKI")

	resources := s.fullPKI(s.cluster, scheme, externalHostnames)

	for _, o := range resources {
		if err := reconcile(ctx, logger, s.client, o, s.cluster); err != nil {
			return err
		}
	}

	return nil
	// TODO Generate all certs
	// TODO Setup all secrets from certs
}

func (s *SelfManager) fullPKI(cluster *v1alpha1.NifiCluster, scheme *runtime.Scheme, externalHostnames []string) []runtime.Object {
	var objects []runtime.Object

	// TODO no need ?
	//objects = append(objects, []runtime.Object{
	//	// A self-signer for the CA Certificate
	//	selfSignerForNamespace(cluster, scheme),
	//	// The CA Certificate
	//	caCertForNamespace(cluster, scheme),
	//	// A issuer backed by the CA certificate - so it can provision secrets
	//	// in this namespace
	//	mainIssuerForNamespace(cluster, scheme),
	//}...,
	//)

	objects = append(objects, pkicommon.ControllerUserForCluster(cluster))
	// Node "users"
	for _, user := range pkicommon.NodeUsersForCluster(cluster, externalHostnames) {
		objects = append(objects, user)
	}
	return objects
}

// TODO
func (s SelfManager) FinalizePKI(ctx context.Context, logger logr.Logger) error {
	panic("implement me")
}
