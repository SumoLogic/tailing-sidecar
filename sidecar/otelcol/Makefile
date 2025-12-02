DOCKER ?= docker

.PHONY: build
build: TARGET ?= local
build: export TAG = $(shell git describe --tags --always)
build: export VERSION = $(TAG:v%=%)
build:
	$(DOCKER) bake $(TARGET)

.PHONY: release-dev
release-dev: TAG = $(shell git describe --tags --always)
release-dev:
	make build TARGET=dev

.PHONY: release
release:
	make build TARGET=prod

.PHONY: build-test-image
build-test-image:
	make build TARGET=test

.PHONY: print-tag
print-tag: TAG = $(shell git describe --tags --always)
print-tag:
	@echo $(TAG)
