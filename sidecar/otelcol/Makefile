BINARY_NAME ?= otelcol-sidecar
OTELCOL_VERSION ?= 0.137.0
BUILDER_REPO ?= github.com/open-telemetry/opentelemetry-collector-releases
GO ?= go

OS ?= $(shell uname -s | tr A-Z a-z)
ARCH ?= $(shell uname -m | sed s/aarch64/arm64/ | sed s/x86_64/amd64/)
LOCALBIN ?= $(shell pwd)/bin


BUILDER_BIN_NAME ?= opentelemetry-collector-builder$(BUILDER_BIN_EXT)
BUILDER_BIN_PATH ?= $(LOCALBIN)
BUILDER=$(BUILDER_BIN_PATH)/$(BUILDER_BIN_NAME)
GORELEASER ?= $(BUILDER_BIN_PATH)/goreleaser


INSTALLED_VERSION := $(shell $(BUILDER) version 2>&1)

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
	@echo "Installing GoReleaser..."
	@mkdir -p $(BUILDER_BIN_PATH)
	@VERSION=$$(curl -sf https://goreleaser.com/static/latest); \
	OS=$$(uname -s ); \
	ARCH=$$(uname -m); \
	echo "Detected OS: $$OS"; \
	echo "Detected ARCH: $$ARCH"; \
    TAR_FILE="goreleaser_${OS}_${ARCH}.tar.gz";\
	RELEASES_URL="https://github.com/goreleaser/goreleaser/releases"; \
	echo "Downloading GoReleaser $$VERSION..."; \
	echo "Download URL: $$RELEASES_URL/download/$$VERSION/$$TAR_FILE";\
	curl -sfLO "$$RELEASES_URL/download/$$VERSION/$$TAR_FILE"; \
	tar -xzf "$$TAR_FILE" -C $(BUILDER_BIN_PATH) ; \
	rm "$$TAR_FILE"; \
	echo "GoReleaser binary installed at $(BUILDER_BIN_PATH)/goreleaser"



.PHONY: ensure-correct-builder-version
ensure-correct-builder-version:
ifeq (,$(findstring $(OTELCOL_VERSION),$(lastword $(INSTALLED_VERSION))))
	@$(error Installed ocb version \
		"$(INSTALLED_VERSION)" \
		does not match the requested "$(OTELCOL_VERSION)". The correct version of \
		ocb can be installed by running 'make install-builder' from the \
		otelcolbuilder directory.\
	)
else
	@printf "Installed ocb (%s) is at the correct version %s\n" \
		"$(shell type $(OCB) | cut -d' ' -f 3)" \
		"$(INSTALLED_VERSION)"
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
		DOCKERHUB_IMAGE=sumologic/tailing-sidecar-otel-dev \
		make release
	git tag -d "$(TAG)"

.PHONY: release
release: ensure-correct-builder-version install-goreleaser
	@$(GORELEASER) release --clean

.PHONY: build-test-image
build-test-image: GHCR_IMAGE = ghcr.io/sumologic/tailing-sidecar-otel
build-test-image: build
	docker tag $(GHCR_IMAGE):main-$(ARCH) $(TAG)

.PHONY: generate-sources
generate-sources:
	$(BUILDER) \
		--config .otelcol-builder.yaml \
		--skip-compilation=true

