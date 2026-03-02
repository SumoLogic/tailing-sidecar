# tailing-sidecar

![Project Status](https://img.shields.io/badge/status-alpha-important?style=for-the-badge)

## TL;DR

```sh
helm repo add tailing-sidecar https://sumologic.github.io/tailing-sidecar
helm repo update
```

```sh
helm upgrade --install tailing-sidecar tailing-sidecar/tailing-sidecar-operator \
  -n tailing-sidecar-system \
  --create-namespace
```

Add `tailing-sidecar` annotation to Pod:

```yaml
metadata:
  annotations:
    tailing-sidecar: <sidecar-name-0>:<volume-name-0>:<path-to-tail-0>;<sidecar-name-1>:<volume-name-1>:<path-to-tail-1>
```

Tailing Sidecar Operator configuration is described [here](operator/docs/configuration.md).

## Tailing Sidecar

**tailing sidecar** is a [streaming sidecar container](https://kubernetes.io/docs/concepts/cluster-administration/logging/#streaming-sidecar-container),
the cluster-level logging agent for Kubernetes.

It helps when your application inside the Pod cannot write to standard output and/or standard error stream
or when it outputs additional logs to a file instead (eg. the gc.log).

It [tails](https://en.wikipedia.org/wiki/Tail_(Unix)) the files inside Kubernetes Pods,
handling situations like the file not being there when tailing starts, tailing multiple files, rotating files, etc.

It uses [Sumologic Collector](https://www.sumologic.com/help/docs/send-data/opentelemetry-collector/) version 0.19.0 onwards.
Before that [Fluent Bit](https://fluentbit.io/) was used.

For more information about cluster-level logging architecture please read Kubernetes
[documentation](https://kubernetes.io/docs/concepts/cluster-administration/logging/#cluster-level-logging-architectures).

The project consists of two parts:

- [tailing sidecar container image](sidecar/) which can be used to manually extend Pods by tailing sidecars
- [tailing sidecar operator](operator/) which automatically adds tailing sidecars to Pods based on configuration
  provided in annotation

## License

This project is released under the [Apache 2.0 License](LICENSE).

## Contributing

Please share your thoughts about tailing sidecar by opening an [issue](https://github.com/SumoLogic/tailing-sidecar/issues/new).

To get started contributing, please refer to our [Contributing](CONTRIBUTING.md) documentation.
