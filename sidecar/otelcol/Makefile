BINARY_NAME ?= otelcol-sidecar
OTELCOL_VERSION ?= 0.90.1
BUILDER_REPO ?= github.com/open-telemetry/opentelemetry-collector
GO ?= go

OS ?= $(shell uname -s | tr A-Z a-z)
ARCH ?= $(shell uname -m | sed s/aarch64/arm64/ | sed s/x86_64/amd64/)
LOCALBIN ?= $(shell pwd)/bin

GORELEASER ?= $(LOCALBIN)/goreleaser

BUILDER_BIN_NAME ?= opentelemetry-collector-builder$(BUILDER_BIN_EXT)
BUILDER_BIN_PATH ?= $(LOCALBIN)
BUILDER=$(BUILDER_BIN_PATH)/$(BUILDER_BIN_NAME)

INSTALLED_VERSION := $(shell $(BUILDER) version 2>&1)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

DOCKERFILE ?= Dockerfile
GHCR_IMAGE = ghcr.io/sumologic/tailing-sidecar-otel
RELEASE_NUMBER ?= 1
VERSION ?= ""

define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
go get -d $(2)@$(3) ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

build:
	docker build \
	--build-arg RELEASE_NUMBER=$(RELEASE_NUMBER) \
	--build-arg VERSION=$(VERSION) \
	--tag $(TAG)  \
	--file ${DOCKERFILE}

push:
	docker push $(TAG)


build-push-multiplatform:
	docker buildx build \
		--push \
		--platform linux/amd64,linux/arm64 \
		--tag ${TAG} \
		.

build-push-ubi:
	$(MAKE) build DOCKERFILE=${DOCKERFILE}.ubi TAG=${TAG}-ubi
	$(MAKE) push DOCKERFILE=${DOCKERFILE}.ubi TAG=${TAG}-ubi



.PHONY: build-test-image
build-test-image: build

