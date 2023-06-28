#!/usr/bin/env bash

set -e

readonly ROOT_DIR="$(dirname "$(dirname "${0}")")"
source "${ROOT_DIR}"/tests/functions.sh

readonly NAMESPACE="tailing-sidecar-system"
readonly TIME=60

wait_for_all_pods_running ${NAMESPACE} ${TIME}

# Check Sidecar resources configuration in tailingsidecar-sample
readonly POD="pod-with-tailing-sidecar-config-resources"
if [ $(kubectl get pod ${POD} -n ${NAMESPACE} -o jsonpath='{.spec.containers[?(@.name=="sidecar-1")].resources.requests.cpu}') != "100m" ];then
  exit 1
fi
if [ $(kubectl get pod ${POD} -n ${NAMESPACE} -o jsonpath='{.spec.containers[?(@.name=="sidecar-1")].resources.requests.memory}') != "100Mi" ];then
  exit 1
fi
if [ $(kubectl get pod ${POD} -n ${NAMESPACE} -o jsonpath='{.spec.containers[?(@.name=="sidecar-1")].resources.limits.cpu}') != "200m" ];then
  exit 1
fi
if [ $(kubectl get pod ${POD} -n ${NAMESPACE} -o jsonpath='{.spec.containers[?(@.name=="sidecar-1")].resources.limits.memory}') != "200Mi" ];then
  exit 1
fi

echo "ok"
exit 0