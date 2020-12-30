module github.com/Orange-OpenSource/nifikop

go 1.15

require (
	emperror.dev/errors v0.4.2
	github.com/antihax/optional v1.0.0
	github.com/banzaicloud/k8s-objectmatcher v1.3.3
	github.com/erdrix/nigoapi v0.0.0-20200824133217-ce90b74151a2
	github.com/go-logr/logr v0.3.0
	github.com/imdario/mergo v0.3.10
	github.com/jarcoal/httpmock v1.0.6
	github.com/jetstack/cert-manager v0.15.1
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/operator-framework/operator-sdk v1.3.0
	github.com/pavel-v-chernykh/keystore-go v2.1.0+incompatible
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v0.19.4
	sigs.k8s.io/controller-runtime v0.7.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
)