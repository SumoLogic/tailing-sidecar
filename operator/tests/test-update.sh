#!/usr/bin/env bash

set -e

# Check Pod logs
[[ $(kubectl logs pod-with-annotations tailing-sidecar-0 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs pod-with-annotations named-sidecar -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs pod-with-annotations tailing-sidecar-1 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1

# Check Deployment logs
readonly DEPLOYMENT_POD_NAME="$(kubectl get pod -l app=deployment-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[0].metadata.name}")"
[[ $(kubectl logs ${DEPLOYMENT_POD_NAME} tailing-sidecar-0 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1

echo "ok"
exit 0
