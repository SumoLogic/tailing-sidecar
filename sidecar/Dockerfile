FROM golang:1.16.5 as go-builder
RUN mkdir /build
ADD ./out_gstdout /build/
WORKDIR /build
RUN make all

FROM fluent/fluent-bit:1.8.1
ENV LOG_LEVEL=warning

COPY --from=go-builder \
  /build/out_gstdout.so \
  /tailing-sidecar/lib/

COPY conf/fluent-bit.conf \
  conf/plugins.conf \
  /fluent-bit/etc/

CMD ["/fluent-bit/bin/fluent-bit", "-c", "/fluent-bit/etc/fluent-bit.conf", "--quiet"]
