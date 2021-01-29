build:
	docker build . --tag $(TAG)

run: build
	docker run --rm -it \
		-v $(DIR_TO_TAIL):/tmp/host \
		-v $(FLUENT_BIT_DB_DIR):/tailing-sidecar/var \
		--env "PATH_TO_TAIL=/tmp/host/$(FILES_PATTERN)" \
		--env "LOG_LEVEL=${LOG_LEVEL}" $(TAG)

push:
	docker push $(TAG)
