# Tailing Sidecar Operator

![Project Status](https://img.shields.io/badge/status-alpha-important?style=for-the-badge)

[Tailing Sidecar Operator](../../operator/README.md) makes it easy to add
[streaming sidecar containers](https://kubernetes.io/docs/concepts/cluster-administration/logging/#streaming-sidecar-container)
to pods in your cluster by adding annotations on the pods.

## Prerequisites

Before installing this chart, ensure the following prerequisites are satisfied in your cluster:

- [admission webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#prerequisites)
  are enabled

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

See [values.yaml](./values.yaml) file for the available configuration options.

### Using `cert-manager` to manage operator's certificates

By default, TLS certificates for the Tailing Sidecar Operator's API webhook
are created using Helm's functions `genCA` and `genSignedCert`.
The generated certificate is valid for 365 days after issuing, i.e. after chart installation.

If you have [cert-manager](https://cert-manager.io/) installed in your cluster,
you can make the chart use it for certificate management by setting the property `useCertManager` to `true`.
