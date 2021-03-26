# Tailing sidecar operator

*Tailing sidecar operator* automatically adds
[streaming sidecar containers](https://kubernetes.io/docs/concepts/cluster-administration/logging/#streaming-sidecar-container)
which use [tailing sidecar image](../sidecar/) to Pods based on configuration provided in annotation.

Configuration for tailing sidecar operator is described [here](docs/configuration.md).

To quickly see benefits of using tailing sidecar operator try it in prepared [Vagrant environment](#testing-in-Vagrant-environment).

## Deploy tailing sidecar operator

### Prerequisities

- [cert-manager](https://cert-manager.io/docs/installation/)
- [admission webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#prerequisites)
  enabled
- [kustomize](https://kustomize.io/)

### Set images (optional)

Set tailing sidecar image:

```bash
export TAILING_SIDECAR_IMG="<some-registry>/<project-name>:tag"
sed -i.backup "s#sumologic/tailing-sidecar:latest#${TAILING_SIDECAR_IMG}#g" config/default/manager_patch.yaml
```

Set tailing-sidecar operator image:

```bash
export TAILING_SIDECAR_OPERATOR_IMG="<some-registry>/<project-name>:tag"
(cd config/manager && kustomize edit set image controller="${TAILING_SIDECAR_OPERATOR_IMG}")
```

### Deploy operator

```bash
kustomize build config/default | kubectl apply -f -
```

### Test operator

Deploy TailingSidecar with configuration e.g.

```bash
kubectl apply -f config/samples/tailing-sidecar_v1_tailingsidecar.yaml -n tailing-sidecar-system
```

to learn more about configuration see [this](docs/configuration.md).

Deploy Pod with `tailing-sidecar` annotation e.g.

```bash
kubectl apply -f examples/pod_with_annotations.yaml
```

Check logs from tailing sidecar e.g.

```bash
kubectl logs pod-with-annotations tailing-sidecar-0  --tail 5 -n tailing-sidecar-system
```

## Build and push tailing sidecar operator image to container registry

To build tailing sidecar operator image:

```bash
make docker-build IMG="<some-registry>/<project-name>:tag"
```

To push tailing sidecar operator image to container registry:

```bash
make docker-push IMG="<some-registry>/<project-name>:tag"
```

## Testing in Vagrant environment

Start and provision the Vagrant environment:

```bash
vagrant up
```

Connect to virtual machine:

```bash
vagrant ssh
```

Build and push tailing sidecar image to local container registry:

```bash
/tailing-sidecar/sidecar/Makefile
```

Go to operator directory:

```bash
cd /tailing-sidecar/operator
```

Deploy operator:

```bash
make
```

Deploy examples:

```bash
make deploy-examples
```

Check that operator added tailing sidecars to example resources:

```bash
make check-examples
```
