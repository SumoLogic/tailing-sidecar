# Tailing sidecar operator configuration

Configuration for tailing sidecar operator is provided through `tailing-sidecar` annotation added to Pod metadata

```yaml
metadata:
  annotations:
    tailing-sidecar: <tailing-sidecar-configuration>
```

Tailing sidecar container is joined with container containing logs by Volume.

Configuration for single tailing sidecar consists of:

- tailing sidecar container name (optional, if not specified container name will be automatically created and
  it will start with "tailing-sidecar" prefix)
- volume name
- path to file containing logs to tail

Configuration for single tailing sidecar is separated by `;`.

Configuration in annotations allows to define configurations for multiple tailing sidecars in following form:

```yaml
metadata:
  annotations:
    tailing-sidecar: <container-name0>:<volume-name0>:<path-to-tail0>;<volume-name1>:<path-to-tail1>;<volume-name2>:<path-to-tail2>
```

## Examples

Example configurations for Kubernetes resources can be found in [examples](../examples) directory.
