#!/usr/bin/env bash

set -e

# Check Pod logs
[[ $(kubectl logs pod-with-annotations tailing-sidecar0 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs pod-with-annotations tailing-sidecar1 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs pod-with-annotations named-container -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1

# Check StatefulSet logs
[[ $(kubectl logs statefulset-with-annotations-0 tailing-sidecar0 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs statefulset-with-annotations-0 my-named-sidecar -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs statefulset-with-annotations-0 tailing-sidecar2 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1

# Check Deployment logs
readonly DEPLOYMENT_POD_NAME="$(kubectl get pod -l app=deployment-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[0].metadata.name}")"
[[ $(kubectl logs ${DEPLOYMENT_POD_NAME} tailing-sidecar0 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs ${DEPLOYMENT_POD_NAME} tailing-sidecar1 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs ${DEPLOYMENT_POD_NAME} tailing-sidecar2 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1

# Check Daemonset logs
readonly DAEMONSET_POD_NAME="$(kubectl get pod -l app=daemonset-with-annotations -n tailing-sidecar-system -o jsonpath="{.items[0].metadata.name}")"
[[ $(kubectl logs ${DAEMONSET_POD_NAME} tailing-sidecar0 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs ${DAEMONSET_POD_NAME} tailing-sidecar1 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs ${DAEMONSET_POD_NAME} tailing-sidecar2 -n tailing-sidecar-system --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1

echo "ok"
exit 0
