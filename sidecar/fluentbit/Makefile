#!/usr/bin/make -f

# mkfile_path is absolute path of this file
# The intention is to be able to run this file from any location
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))

# root_dir is a root directory of the project (github repo)
root_dir := $(dir $(abspath $(mkfile_path)/..))

TAG ?= "localhost:32000/sumologic/tailing-sidecar:latest"
DOCKERFILE ?= Dockerfile
RELEASE_NUMBER ?= 1
VERSION ?= ""

all: build push

build:
	docker build \
	--build-arg RELEASE_NUMBER=$(RELEASE_NUMBER) \
	--build-arg VERSION=$(VERSION) \
	--tag $(TAG) \
	--file ${DOCKERFILE} .

build-test-image: build

run: build
	docker run --rm -it \
		-v $(DIR_TO_TAIL):/tmp/host \
		-v $(FLUENT_BIT_DB_DIR):/tailing-sidecar/var \
		--env "PATH_TO_TAIL=/tmp/host/$(FILES_PATTERN)" \
		--env "LOG_LEVEL=${LOG_LEVEL}" $(TAG)

build-push-multiplatform:
	docker buildx build \
		--push \
		--platform linux/amd64,linux/arm/v7,linux/arm64 \
		--tag ${TAG} \
		.

build-push-ubi:
	$(MAKE) build DOCKERFILE=${DOCKERFILE}.ubi TAG=${TAG}-ubi
	$(MAKE) push DOCKERFILE=${DOCKERFILE}.ubi TAG=${TAG}-ubi

push:
	docker push $(TAG)

deploy-examples:
	kubectl apply -f examples/pod_with_tailing_sidecars.yaml
	kubectl wait --for=condition=ready --timeout 60s pod example-with-tailling-sidecars

test-examples:
	tests/test.sh

teardown-examples:
	kubectl delete pod example-with-tailling-sidecars
