ARG ALPINE_VERSION=3.22.1
ARG GO_VERSION=1.24.5
ARG OCB_VERSION=0.130.0
ARG YQ_VERSION=4.46.1

#################################################################################
# Set version information in builder config.
#################################################################################
FROM mikefarah/yq:${YQ_VERSION} AS config

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG VERSION
ARG OCB_VERSION

RUN echo "I am running on ${BUILDPLATFORM}, building for ${TARGETPLATFORM}"

# Temporarily switch to root to install required packages
USER root
# hadolint ignore=DL3018
RUN apk add --no-cache bash
USER yq

WORKDIR /workdir

# Copy the builder template and generate builder config script to the stage
COPY builder-template.yaml builder-template.yaml
COPY generate-builder-config.sh /usr/local/bin/generate-builder-config.sh

ENV DIST_VERSION=${VERSION}
ENV OT_VERSION=${OCB_VERSION}

RUN generate-builder-config.sh

#################################################################################
# Generate the collector source.
#################################################################################
FROM otel/opentelemetry-collector-builder:${OCB_VERSION} AS source

# Temporarily switch to root to install required packages
USER root
# hadolint ignore=DL3018
RUN apk add --no-cache git
USER ocb

WORKDIR /build

# Copy the builder config
COPY --from=config /workdir/builder-config.yaml /build/config.yaml

RUN ocb --config=config.yaml --skip-compilation=true

#################################################################################
# Build the collector.
#################################################################################
FROM golang:${GO_VERSION} AS builder

WORKDIR /go

# Copy the collector source code
COPY --from=source --chown=root:root /build/cmd .

ENV CGO_ENABLED=0

RUN go build .

#################################################################################
# Prepare the necessary directories in an Alpine container where we have the
# necessary tooling.
#################################################################################
FROM alpine:${ALPINE_VERSION} AS directories
RUN mkdir /etc/otel /var/lib/otc && touch /var/log/otelcol.log

#################################################################################
# Build a minimal image from scratch, copying over the binary & directories from
# the previous stages.
#################################################################################
FROM scratch

ARG USER_UID=10001
USER ${USER_UID}

COPY --from=directories --chown=${USER_UID}:${USER_UID} /etc/otel/ /etc/otel/
COPY --from=directories --chown=${USER_UID}:${USER_UID} /var/lib/otc /var/lib/otc
COPY --from=directories --chown=${USER_UID}:${USER_UID} /var/log/otelcol.log /var/log/otelcol.log

COPY --from=builder /go/sidecar /otelcol

# Copy the default tailing-sidecar configuration file
COPY ./config.yaml /etc/otel/config.yaml

ENTRYPOINT ["/otelcol"]
CMD ["--config", "/etc/otel/config.yaml"]
