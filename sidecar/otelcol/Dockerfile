# Prepare the necessary directories in an Alpine container where we have the necessary tooling
FROM alpine:3.19.0 as directories
RUN mkdir /etc/otel/
RUN mkdir /var/lib/otc
RUN touch /var/log/otelcol.log

FROM scratch

ARG USER_UID=10001
USER ${USER_UID}

COPY --from=directories --chown=${USER_UID}:${USER_UID} /etc/otel/ /etc/otel/
COPY --from=directories --chown=${USER_UID}:${USER_UID} /var/lib/otc /var/lib/otc
COPY --from=directories --chown=${USER_UID}:${USER_UID} /var/log/otelcol.log /var/log/otelcol.log

# copy the default tailing-sidecar configuration file
COPY ./config.yaml /etc/otel/config.yaml

COPY otelcol-sidecar /otelcol
ENTRYPOINT ["/otelcol"]
CMD ["--config", "/etc/otel/config.yaml"]
