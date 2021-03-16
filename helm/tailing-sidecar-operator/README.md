# Tailing Sidecar Operator

![Project Status](https://img.shields.io/badge/status-alpha-important?style=for-the-badge)

[Tailing Sidecar Operator](https://github.com/SumoLogic/tailing-sidecar/tree/main/operator) makes it easy to add
[streaming sidecar containers](https://kubernetes.io/docs/concepts/cluster-administration/logging/#streaming-sidecar-container)
to pods in your cluster by adding annotations on the pods.

## Prerequisites

Before installing this chart, ensure the following prerequisites are satisfied in your cluster:

- [cert-manager](https://cert-manager.io/docs/installation/) is installed
- [admission webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#prerequisites)
  are enabled

Unfortunately, due to a limitation of `cert-manager` chart, it cannot be used as a dependency and needs to be installed separately.
See issue [jetstack/cert-manager#3062](https://github.com/jetstack/cert-manager/issues/3062) for details.

## Installing

Having satisfied the [Prerequisites](#prerequisites), run the following to install the chart:

```sh
helm repo add TBD TBD
helm install tailing-sidecar-operator TBD
```

## Uninstalling

```sh
helm uninstall tailing-sidecar-operator
```

## Configuration

TBD
