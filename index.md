# Tailing Sidecar Operator

[Tailing sidecar operator](https://github.com/SumoLogic/tailing-sidecar/tree/main/operator) adds
[streaming sidecar containers](https://kubernetes.io/docs/concepts/cluster-administration/logging/#streaming-sidecar-container)
which use [tailing sidecar image](https://github.com/SumoLogic/tailing-sidecar/tree/main/sidecar) to Pods based on configuration provided in annotation.

## Add Tailing Sidecar Operator Helm Repository

```sh
helm repo add tailing-sidecar https://sumologic.github.io/tailing-sidecar
helm repo update
```

## Install Tailing Sidecar Operator Helm Chart

```sh
helm upgrade --install tailing-sidecar tailing-sidecar/tailing-sidecar-operator \
  -n tailing-sidecar-system \
  --create-namespace
```

See [values.yaml](https://github.com/SumoLogic/tailing-sidecar/blob/main/helm/tailing-sidecar-operator/values.yaml)
file for the available Helm Chart configuration options.

## Tailing Sidecar Operator Configuration

Add tailing-sidecar annotation to Pod:

```sh
metadata:
  annotations:
    tailing-sidecar: <sidecar-name-0>:<volume-name-0>:<path-to-tail-0>;<sidecar-name-1>:<volume-name-1>:<path-to-tail-1>
```

For more details related to configuration please see [documentation](https://github.com/SumoLogic/tailing-sidecar/blob/main/operator/docs/configuration.md).
