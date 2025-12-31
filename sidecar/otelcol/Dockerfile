FROM public.ecr.aws/sumologic/sumologic-otel-collector:0.141.0-sumo-0 AS collector

FROM alpine:3.21.3

RUN mkdir -p /etc/otel /var/lib/otc /var/log

RUN addgroup -g 10001 otelcol && \
    adduser -D -u 10001 -G otelcol otelcol

COPY --from=collector / /
COPY ./config.yaml /etc/otel/config.yaml
RUN chown -R otelcol:otelcol /etc/otel /var/lib/otc /var/log
USER 10001

ENTRYPOINT ["/otelcol-sumo"]
CMD ["--config", "/etc/otel/config.yaml"]
