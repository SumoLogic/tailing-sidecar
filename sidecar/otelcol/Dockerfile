FROM alpine:3.21.3


RUN mkdir -p /etc/otel /var/lib/otc /var/log
COPY ./config.yaml /etc/otel/config.yaml
COPY dist/otelcol-sidecar_linux_arm64_*/otelcol-sidecar /otelcol
RUN chmod +x /otelcol

ENTRYPOINT ["/otelcol"]
CMD ["--config", "/etc/otel/config.yaml"]
