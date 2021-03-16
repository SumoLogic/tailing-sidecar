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

Configuration can be provided either through simple configurations in annotations or named configurations
using TailingSidecar resources.

Configuration for single tailing sidecar is separated by `;`.

## Simple configurations in annotations

Simple configurations in annotations allow to define configurations for multiple tailing sidecars in following form:

```yaml
metadata:
  annotations:
    tailing-sidecar: <container-name0>:<volume-name0>:<path-to-tail0>;<volume-name1>:<path-to-tail1>;<volume-name2>:<path-to-tail2>
```

## Named configurations using TailingSidecar

![Immature Config Status](https://img.shields.io/badge/config-immature-important?style=for-the-badge)

**_NOTE: this configuration option is being reviewed and most probably will change. Use at your own risk._**

Named configurations using TailingSidecar allow to define named configurations for multiple tailing sidecars
in the following form:

```yaml
apiVersion: tailing-sidecar.sumologic.com/v1
kind: TailingSidecar
metadata:
  name: tailingsidecar-name
spec:
  configs:
    <config-name0>:
      container: <container-name0>
      volume: <volume-name0>
      file: <path-to-tail0>
    <config-name1>:
      volume: <volume-name1>
      file: <path-to-tail1>
    <config-name2>:
      container: <container-name2>
      volume: <volume-name2>
      file: <path-to-tail2>
```

Named configurations needs to be added to annotations:

```yaml
metadata:
  annotations:
    tailing-sidecar: <config-name0>;<config-name1>;<config-name2>
```

Named configurations can be mixed with simple configurations so following form of configuration is supported:

```yaml
metadata:
  annotations:
    tailing-sidecar: <config-name0>;<config-name1>;<config-name2>;<container-name3>:<volume-name3>:<path-to-tail3>;<volume-name4>:<path-to-tail4>
```

## Examples

Example configurations for Kubernetes resources can be found in [examples](../examples) directory.
