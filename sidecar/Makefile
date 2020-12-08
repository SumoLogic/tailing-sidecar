build:
	docker build . --tag sidecar

run: build
	docker run --rm -it -v $(shell pwd):/tmp/host sidecar /fluent-bit/bin/fluent-bit -i tail -p path=/tmp/host/test.log -e /opt/out_gstdout.so  -o gstdout --quiet
