FROM golang:1.20.5 as go-builder
RUN mkdir /build
ADD ./out_gstdout /build/
WORKDIR /build
RUN make all

# Manually upgrade libc and zlib to address https://www.cve.org/CVERecord?id=CVE-2022-37434 and https://www.cve.org/CVERecord?id=CVE-2021-3999
# TODO: Revert this after fluent-bit 1.9.9 is released
FROM debian:bookworm-slim as deb-extractor

# We download all debs locally then extract them into a directory we can use as the root for distroless.
# This is directly copied from Fluent Bit's upstream Dockerfile
WORKDIR /tmp
RUN apt-get update && \
    apt-get download \
        zlib1g \
        libc6 && \
    mkdir -p /dpkg/var/lib/dpkg/status.d/ && \
    for deb in *.deb; do \
        package_name=$(dpkg-deb -I ${deb} | awk '/^ Package: .*$/ {print $2}'); \
        echo "Processing: ${package_name}"; \
        dpkg --ctrl-tarfile $deb | tar -Oxf - ./control > /dpkg/var/lib/dpkg/status.d/${package_name}; \
        dpkg --extract $deb /dpkg || exit 10; \
    done

# Remove unnecessary files extracted from deb packages like man pages and docs etc.
RUN find /dpkg/ -type d -empty -delete && \
    rm -r /dpkg/usr/share/doc/

FROM fluent/fluent-bit:2.0.11
ENV LOG_LEVEL=warning

# Copy the libraries from the extractor stage into root
COPY --from=deb-extractor /dpkg /

COPY --from=go-builder \
  /build/out_gstdout.so \
  /tailing-sidecar/lib/

COPY conf/fluent-bit.conf \
  conf/plugins.conf \
  /fluent-bit/etc/

CMD ["/fluent-bit/bin/fluent-bit", "-c", "/fluent-bit/etc/fluent-bit.conf", "--quiet"]
