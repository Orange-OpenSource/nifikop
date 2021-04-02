package selfmanagerpki

import (
	"context"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	"k8s.io/apimachinery/pkg/runtime"
)

// TODO
func (s selfManager) ReconcileUserCertificate(ctx context.Context, user *v1alpha1.NifiUser, scheme *runtime.Scheme) (*pki.UserCertificate, error) {
	panic("implement me")
}

// TODO
func (s selfManager) FinalizeUserCertificate(ctx context.Context, user *v1alpha1.NifiUser) error {
	panic("implement me")
}
