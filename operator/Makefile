# Current Operator version
VERSION ?= 0.0.1
# Default bundle image tag
BUNDLE_IMG ?= controller-bundle:$(VERSION)
# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# Image URL to use all building/pushing image targets
IMG ?= localhost:32000/sumologic/tailing-sidecar-operator:latest

# Image URL for tailing sidecar
TAILING_SIDECAR_IMG ?= localhost:32000/sumologic/tailing-sidecar:demo

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Run tests
# Set default shell as bash
SHELL := /bin/bash
ENVTEST_ASSETS_DIR = $(shell pwd)/testbin

all: docker-build docker-push deploy

test: generate fmt vet manifests
	mkdir -p $(ENVTEST_ASSETS_DIR)
	test -f $(ENVTEST_ASSETS_DIR)/setup-envtest.sh || curl -sSLo $(ENVTEST_ASSETS_DIR)/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.6.3/hack/setup-envtest.sh
	source $(ENVTEST_ASSETS_DIR)/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test ./... -coverprofile cover.out

clean:
	git checkout -- config/default/manager_args_patch.yaml
	rm -f config/default/manager_args_patch.yaml.backup
	rm -rf testbin

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests kustomize
	sed -i.backup "s#sumologic/tailing-sidecar:latest#${TAILING_SIDECAR_IMG}#g" config/default/manager_args_patch.yaml
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

# Remove resources created with make deploy
teardown:
	$(KUSTOMIZE) build config/default | kubectl delete -f -

# Deploy examples
deploy-examples:
	kubectl apply -f config/samples/tailing-sidecar_v1_tailingsidecar.yaml -n tailing-sidecar-system
	kubectl apply -f examples/pod_with_annotations.yaml
	kubectl apply -f examples/deployment_with_annotations.yaml
	kubectl apply -f examples/pod_without_annotations.yaml

# Deploy examples update
deploy-examples-update:
	kubectl apply -f examples/pod_with_annotations_update.yaml --force
	kubectl apply -f examples/deployment_with_annotations_update.yaml --force

# Remove resources created with make deploy-examples
teardown-examples:
	kubectl delete -n tailing-sidecar-system pod pod-without-annotations
	kubectl delete -n tailing-sidecar-system pod pod-with-annotations
	kubectl delete -n tailing-sidecar-system deployment deployment-with-annotations
	kubectl delete -n tailing-sidecar-system TailingSidecar tailingsidecar-sample

# Check resources create with make deploy-examples
check-examples:
	kubectl get all -n tailing-sidecar-system

	kubectl logs pod-with-annotations tailing-sidecar0 -n tailing-sidecar-system --tail 5
	kubectl logs pod-with-annotations tailing-sidecar1 -n tailing-sidecar-system --tail 5
	kubectl logs pod-with-annotations tailing-sidecar2 -n tailing-sidecar-system --tail 5

	kubectl logs $(shell kubectl get pod -l app=deployment-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[0].metadata.name}") \
		tailing-sidecar0 -n tailing-sidecar-system --tail 5
	kubectl logs $(shell kubectl get pod -l app=deployment-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[0].metadata.name}") \
		tailing-sidecar1 -n tailing-sidecar-system --tail 5
	kubectl logs $(shell kubectl get pod -l app=deployment-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[0].metadata.name}") \
		tailing-sidecar2 -n tailing-sidecar-system --tail 5

	kubectl logs $(shell kubectl get pod -l app=deployment-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[1].metadata.name}") \
		tailing-sidecar0 -n tailing-sidecar-system --tail 5
	kubectl logs $(shell kubectl get pod -l app=deployment-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[1].metadata.name}") \
		tailing-sidecar1 -n tailing-sidecar-system --tail 5
	kubectl logs $(shell kubectl get pod -l app=deployment-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[1].metadata.name}") \
		tailing-sidecar2 -n tailing-sidecar-system --tail 5

	kubectl logs $(shell kubectl get pod -l app=deployment-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[2].metadata.name}") \
		tailing-sidecar0 -n tailing-sidecar-system --tail 5
	kubectl logs $(shell kubectl get pod -l app=deployment-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[2].metadata.name}") \
		tailing-sidecar1 -n tailing-sidecar-system --tail 5
	kubectl logs $(shell kubectl get pod -l app=deployment-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[2].metadata.name}") \
		tailing-sidecar2 -n tailing-sidecar-system --tail 5

# Check resources create with make deploy-examples-update
check-examples-update:
	kubectl get all -n tailing-sidecar-system
	kubectl logs pod-with-annotations tailing-sidecar3 -n tailing-sidecar-system --tail 5

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

kustomize:
ifeq (, $(shell which kustomize))
	@{ \
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4 ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
KUSTOMIZE=$(GOBIN)/kustomize
else
KUSTOMIZE=$(shell which kustomize)
endif

# Generate bundle manifests and metadata, then validate generated files.
.PHONY: bundle
bundle: manifests kustomize
	operator-sdk generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | operator-sdk generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	operator-sdk bundle validate ./bundle

# Build the bundle image.
.PHONY: bundle-build
bundle-build:
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .
