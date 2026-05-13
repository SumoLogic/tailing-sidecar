# Tailing Sidecar Operator configuration

## Configuration in annotation

Configuration for Tailing Sidecar Operator can be provided through `tailing-sidecar` annotation added to Pod metadata

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

Example configurations in annotations for Kubernetes resources can be found in [examples](../examples) directory.

**Notice**: Only basic options can be configured in annotations, for extended configuration options please
see [Configuration in TailingSidecarConfig](#configuration-in-tailingsidecarconfig).

## Configuration in TailingSidecarConfig

Configuration for Tailing Sidecar Operator can be provided through `TailingSidecarConfig` which is a Custom Resource
used by the operator.

Example definitions of `TailingSidecarConfig` are available in [samples](../config/samples) directory and Pod adjusted to this configuration
in available [pod-with-tailing-sidecar-config.yaml](../examples/pod_with_tailing_sidecar_config.yaml).

To try example configuration use following commands:

```bash
kubectl apply -f https://raw.githubusercontent.com/SumoLogic/tailing-sidecar/release-v0.5/operator/config/samples/tailing-sidecar_v1_tailingsidecar.yaml
kubectl apply -f https://raw.githubusercontent.com/SumoLogic/tailing-sidecar/release-v0.5/operator/examples/pod_with_tailing_sidecar_config.yaml
```

For details related to `TailingSidecarConfig` definition please see subsections below.

### TailingSidecarConfig

| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata | Metadata for TailingSidecarConfig | [metav1.ObjectMeta][metav1.ObjectMeta] |
| spec | Spec defines specification of TailingSidecarConfig | [tailingsidecarv1.TailingSidecarConfigSpec](#tailingSidecarConfigSpec) |

[metav1.ObjectMeta]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta

### TailingSidecarConfigSpec

| Field | Description | Scheme |
| ----- | ----------- | ------ |
| annotationsPrefix | AnnotationsPrefix defines prefix for per container annotations. | [metav1.LabelSelector][metav1.LabelSelector] |
| podSelector | PodSelector selects Pods to which this tailing sidecar configuration applies. | [metav1.LabelSelector][metav1.LabelSelector] |
| SidecarSpecs | SidecarSpecs defines specifications for tailing sidecar containers, map key indicates name of tailing sidecar container. | [map\[string\]tailingsidecarv1.SidecarSpec](#sidecarspec) |

[metav1.LabelSelector]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#labelselector-v1-meta

### SidecarSpec

| Field       | Description                                                                                                                                                                                                     | Scheme |
|-------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------| ------ |
| annotations | Annotations defines tailing sidecar container annotations. Tailing sidecar annotations are added in following form `<annotationsPrefix>/<tailing-sidecar-containter-name>.<annotation-key>:<annotation-value>`  | map\[string\]string |
| path        | Path defines path to a file containing logs to tail within a tailing sidecar container.                                                                                                                         | string |
| volumeMount | VolumeMount describes a mounting of a volume within a tailing sidecar container. This volume joins tailing sidecar container with container containing logs to tail and provide access to file with logs.       | [corev1.VolumeMount][corev1.VolumeMount] |
| resources   | resources describes the compute resource requirements for a tailing sidecar container.  | [corev1.ResourceRequirements][corev1.ResourceRequirements] |
[corev1.VolumeMount]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#volumemount-v1-core
[corev1.ResourceRequirements]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core

