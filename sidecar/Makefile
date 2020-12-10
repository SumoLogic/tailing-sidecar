build:
	docker build . --tag sidecar:dev

run: build
	docker run --rm -it -v $(shell pwd):/tmp/host --env "PATH_TO_TAIL=/tmp/host/*.log" sidecar:dev
