# Development

This document contains information helpful for developers.

## Development Builds

Development builds are created on each merge to `main` branch.

### Container Images

Container images for [tailing sidecar](../sidecar) and [tailing sidecar operator](../operator) are pushed to
Github Container Registry. The latest development images have `main` tag.

To pull the latest tailing sidecar image use:

```sh
docker pull ghcr.io/sumologic/tailing-sidecar:main
```

To pull the latest tailing sidecar operator image use:

```sh
docker pull ghcr.io/sumologic/tailing-sidecar-operator:main
```

### Helm Charts

Helm Charts are available on [gh-pages](https://github.com/SumoLogic/tailing-sidecar/tree/gh-pages) branch.
Development Helm Charts are stored in [dev](https://github.com/SumoLogic/tailing-sidecar/tree/gh-pages/dev) directory.

To add development repository:

```sh
helm repo add tailing-sidecar-dev https://sumologic.github.io/tailing-sidecar/dev
helm repo update
```

To install development Helm chart with the latest development container images:

```sh
helm upgrade --install tailing-sidecar tailing-sidecar-dev/tailing-sidecar-operator \
   --namespace tailing-sidecar-system \
   --create-namespace \
   --set operator.image.tag=main \
   --set sidecar.image.tag=main \
   --version <CHART_VERSION>
```

e.g.

```sh
helm upgrade --install tailing-sidecar tailing-sidecar-dev/tailing-sidecar-operator \
   --namespace tailing-sidecar-system \
   --create-namespace \
   --set operator.image.tag=main \
   --set sidecar.image.tag=main \
   --version 0.1.0-13-g177189e057df0180b46232ebea53f60fa93d242f
```
