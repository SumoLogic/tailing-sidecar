#!/usr/bin/env bash

function wait_for_pod() {
  local namespace="$1"
  local pod="$2"
  local time="$3"
  local state="Running"

  for i in $(seq 0 $time)
  do
    get="$(kubectl get pods -n $namespace $pod)"
    if [[ ! $get =~ $state ]]; then
      echo "Waiting for pod $pod in $i interation"
      sleep 1
    else
      echo "Found pod $pod"
      break
    fi
  done
}

function wait_for_all_pods_running() {
  local namespace="$1"
  local time="$2"

  for i in $(seq 0 $time)
  do
    pods="$(kubectl get pods -n $namespace | tail -n +2 | wc -l)"
    running_pods="$(kubectl get pods -n $namespace | grep Running | wc -l)"
    if [[ $pods -eq $running_pods ]]; then
      echo "All Pods are in Running state"
      break
    fi
    echo "$(kubectl get pods -n $namespace)"
    sleep 1
  done
}
