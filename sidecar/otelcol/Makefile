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

CONTAINER_REGISTRY := ghcr.io/sumologic/tailing-sidecar-otel

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
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

.PHONY: install-builder
install-builder:
	@mkdir -p "$(BUILDER_BIN_PATH)"
	curl -L -o "$(BUILDER_BIN_PATH)/$(BUILDER_BIN_NAME)" https://$(BUILDER_REPO)/releases/download/cmd/builder/v$(OTELCOL_VERSION)/ocb_$(OTELCOL_VERSION)_$(OS)_$(ARCH)$(BUILDER_BIN_EXT)
	@chmod +x "$(BUILDER_BIN_PATH)/$(BUILDER_BIN_NAME)"
	@$(MAKE) ensure-correct-builder-version

.PHONY: install-goreleaser
install-goreleaser:
	$(call go-get-tool,$(GORELEASER),github.com/goreleaser/goreleaser,latest)


.PHONY: ensure-correct-builder-version
ensure-correct-builder-version:
ifneq ($(lastword $(INSTALLED_VERSION)),$(OTELCOL_VERSION))
	@$(error Installed opentelemetry-collector-builder version \
		"$(INSTALLED_VERSION)" \
		does not match the requested "$(OTELCOL_VERSION)" \
		Please check if "$(BUILDER_BIN_PATH)" can be found in your PATH \
		and if not, then install it using 'make install-builder' from otelcolbuilder's directory\
	)
endif

.PHONY: build
build: ensure-correct-builder-version install-goreleaser
	@$(GORELEASER) release --snapshot --clean

.PHONY: release-dev
release-dev: TAG = $(shell git describe --tags --always)
release-dev: ensure-correct-builder-version install-goreleaser
	git tag "$(TAG)"
	GORELEASER_CURRENT_TAG=$(TAG) \
		ECR_IMAGE=public.ecr.aws/sumologic/tailing-sidecar-otel-dev \
		GHCR_IMAGE=ghcr.io/sumologic/tailing-sidecar-otel \
		make release
	git tag -d "$(TAG)"

.PHONY: release
release: ensure-correct-builder-version install-goreleaser
	@$(GORELEASER) release --clean

.PHONY: build-test-image
build-test-image: build
	docker tag $(CONTAINER_REGISTRY):latest-$(ARCH) $(TAG)

.PHONY: generate-sources
generate-sources:
	$(BUILDER) \
		--go $(GO) \
		--version "$(VERSION)" \
		--config .otelcol-builder.yaml \
		--skip-compilation=true

