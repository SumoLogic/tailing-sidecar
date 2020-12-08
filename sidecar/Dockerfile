FROM golang:1.14 as go-builder
RUN mkdir /build
ADD ./src/go /build/
WORKDIR /build
RUN set -x \
  && go get github.com/fluent/fluent-bit-go/output \
  && go build -buildmode=c-shared -o out_gstdout.so .


FROM fluent/fluent-bit:1.6.8

COPY --from=go-builder \
  /build/out_gstdout.so \
  /opt/
