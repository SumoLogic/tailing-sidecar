#!/usr/bin/make -f

# mkfile_path is absolute path of this file
# The intention is to be able to run this file from any location
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))

# root_dir is a root directory of the project (github repo)
root_dir := $(dir $(abspath $(mkfile_path)/..))

TAG ?= "localhost:32000/sumologic/tailing-sidecar:demo"

all: build push

build:
	docker build --tag $(TAG) --file ${root_dir}sidecar/Dockerfile ${root_dir}sidecar/

run: build
	docker run --rm -it \
		-v $(DIR_TO_TAIL):/tmp/host \
		-v $(FLUENT_BIT_DB_DIR):/tailing-sidecar/var \
		--env "PATH_TO_TAIL=/tmp/host/$(FILES_PATTERN)" \
		--env "LOG_LEVEL=${LOG_LEVEL}" $(TAG)

push:
	docker push $(TAG)
