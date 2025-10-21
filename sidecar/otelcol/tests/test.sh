#!/usr/bin/env bash

set -e

[[ $(kubectl logs --tail 5 example-with-otelcol-tailing-sidecars sidecar1 | grep example1 | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs --tail 5 example-with-otelcol-tailing-sidecars sidecar2 | grep example2 | wc -l) -ne 5 ]] && exit 1

echo "ok"
exit 0
