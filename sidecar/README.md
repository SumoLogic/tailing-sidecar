# Tailing sidecar

## Run tailing sidecar in Docker container

To run tailing sidecar in Docker container define following variables:

- `DIR_TO_TAIL` - path to directory with files to read
- `FLUENT_BIT_DB_DIR` - path to directory where Fluent Bit database will be created,
  Fluent Bit uses database to track monitored files and offsets, for details please see tail plugin
  [documentation](https://docs.fluentbit.io/manual/pipeline/inputs/tail)
- `FILES_PATTERN` - pattern to match files in directory specified as `DIR_TO_TAIL`
- `TAILING_SIDECAR_IMAGE` - tailing sidecar Docker image
- `LOG_LEVEL` - verbosity level, by default 'warning' is set,
  allowed values: error, warning, info, debug, trace

e.g.

```bash
export DIR_TO_TAIL="$PWD/examples"
export FLUENT_BIT_DB_DIR="$PWD/var"
export FILES_PATTERN="*.log"
export TAILING_SIDECAR_IMAGE="localhost:32000/sumologic/tailing-sidecar:demo"
export LOG_LEVEL="warning"
```

And run tailing sidecar in Docker container:

```bash
docker run --rm -it \
    -v ${DIR_TO_TAIL}:/tmp/host \
    -v ${FLUENT_BIT_DB_DIR}:/tailing-sidecar/var \
    --env "PATH_TO_TAIL=/tmp/host/${FILES_PATTERN}" \
    --env "LOG_LEVEL=${LOG_LEVEL}" ${TAILING_SIDECAR_IMAGE}
```

## Build and run tailing sidecar in Docker container

To build and run Docker container with tailing sidecar to tail files in `$PWD/examples`
which match pattern `*.log` and save Fluent Bit database in `$PWD/var`:

```bash
make run \
    TAG=sidecar:dev \
    DIR_TO_TAIL="$PWD/examples" \
    FILES_PATTERN="*.log" \
    FLUENT_BIT_DB_DIR="$PWD/var" \
    LOG_LEVEL="warning"
```

## Build and push tailing sidecar to Docker registry

To build Docker image with tailing sidecar:

```bash
make build TAG=<DOCKER_IMAGE_TAG>
```

e.g.

```bash
make build TAG=localhost:32000/sumologic/tailing-sidecar:demo
```

To push Docker image to container registry:

```bash
make push TAG=<DOCKER_IMAGE_TAG>
```

e.g.

```bash
make push TAG=localhost:32000/sumologic/tailing-sidecar:demo
```
