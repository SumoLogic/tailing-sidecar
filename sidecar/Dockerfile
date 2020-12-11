FROM golang:1.14 as go-builder
RUN mkdir /build
ADD ./src/go /build/

WORKDIR /build
RUN set -x \
  && go get github.com/fluent/fluent-bit-go/output \
  && go build -buildmode=c-shared -o out_gstdout.so .

RUN set -x \
  && mkdir -p tailing-sidecar/lib \
  && mkdir -p tailing-sidecar/var \
  && cp out_gstdout.so tailing-sidecar/lib

FROM fluent/fluent-bit:1.6.8
ENV LOG_LEVEL=warning

COPY --from=go-builder \
  /build/tailing-sidecar/ \
  /tailing-sidecar/

COPY conf/fluent-bit.conf \
  conf/plugins.conf \
  /fluent-bit/etc/

CMD ["/fluent-bit/bin/fluent-bit", "-c", "/fluent-bit/etc/fluent-bit.conf", "--quiet"]
