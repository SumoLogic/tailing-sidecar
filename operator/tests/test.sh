#!/usr/bin/env bash

set -e

readonly ROOT_DIR="$(dirname "$(dirname "${0}")")"
source "${ROOT_DIR}"/tests/functions.sh

readonly NAMESPACE="tailing-sidecar-system"
readonly TIME=60

wait_for_all_pods_running ${NAMESPACE} ${TIME}
kubectl describe pods -n ${NAMESPACE}

# Check Pod logs
readonly POD="pod-with-annotations"
wait_for_pod ${NAMESPACE} ${POD} ${TIME}
[[ $(kubectl logs ${POD} tailing-sidecar-0 -n ${NAMESPACE} --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs ${POD} named-container -n ${NAMESPACE} --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1

# Check StatefulSet logs
readonly STATEFULSET_POD_NAME="statefulset-with-annotations-0"
wait_for_pod ${NAMESPACE} ${STATEFULSET_POD_NAME} ${TIME}
[[ $(kubectl logs ${STATEFULSET_POD_NAME} my-named-sidecar -n ${NAMESPACE} --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs ${STATEFULSET_POD_NAME} tailing-sidecar-1 -n ${NAMESPACE} --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1

# Check Deployment logs
readonly DEPLOYMENT_POD_NAME="$(kubectl get pod -l app=deployment-with-annotations -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")"
wait_for_pod ${NAMESPACE} ${DEPLOYMENT_POD_NAME} ${TIME}
[[ $(kubectl logs ${DEPLOYMENT_POD_NAME} tailing-sidecar-0 -n ${NAMESPACE} --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs ${DEPLOYMENT_POD_NAME} tailing-sidecar-1 -n ${NAMESPACE} --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1

# Check Daemonset logs
readonly DAEMONSET_POD_NAME="$(kubectl get pod -l app=daemonset-with-annotations -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")"
wait_for_pod ${NAMESPACE} ${DAEMONSET_POD_NAME} ${TIME}
[[ $(kubectl logs ${DAEMONSET_POD_NAME} tailing-sidecar-0 -n ${NAMESPACE} --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs ${DAEMONSET_POD_NAME} tailing-sidecar-1 -n ${NAMESPACE} --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1

# Test Pod with configuration in CRD
readonly POD_WITH_CRD="pod-with-annotations-crd"
wait_for_pod ${NAMESPACE} ${POD_WITH_CRD} ${TIME}
[[ $(kubectl logs ${POD_WITH_CRD} tailing-sidecar-0 -n ${NAMESPACE} --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1
[[ $(kubectl logs ${POD_WITH_CRD} tailing-sidecar-1 -n ${NAMESPACE} --tail 5 | grep example | wc -l) -ne 5 ]] && exit 1

echo "ok"
exit 0
