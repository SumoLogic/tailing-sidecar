module github.com/SumoLogic/tailing-sidecar/operator

go 1.13

require (
	github.com/go-logr/logr v0.4.0
	github.com/google/uuid v1.2.0 // indirect
	github.com/nsf/jsondiff v0.0.0-20200515183724-f29ed568f4ce
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	gomodules.xyz/jsonpatch/v2 v2.2.0
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.8.3
)
