# Tailing sidecar container image

**tailing sidecar container image** is a an image which can be used to manually extend Pods by [streaming sidecar containers](https://kubernetes.io/docs/concepts/cluster-administration/logging/#streaming-sidecar-container).

## Getting Started

To understand benefits of using tailing sidecar see example below.

### Extend Pod by adding tailing sidecars

Assuming that container writes logs to two different files and Pod has this specification:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example-with-tailling-sidecars
spec:
  containers:
  - name: count
    image: busybox
    args:
    - /bin/sh
    - -c
    - >
      i=0;
      while true;
      do
        echo "example1: $i $(date)" >> /var/log/example1.log;
        echo "example2: $i $(date)" >> /var/log/example2.log;
        i=$((i+1));
        sleep 1;
      done
    volumeMounts:
    - name: varlog
      mountPath: /var/log
  volumes:
  - name: varlog
    emptyDir: {}
```

Pod can be extended by adding tailing sidecar containers for easier log access:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example-with-tailling-sidecars
spec:
  containers:
  - name: count
    image: busybox
    args:
    - /bin/sh
    - -c
    - >
      i=0;
      while true;
      do
        echo "example1: $i $(date)" >> /var/log/example1.log;
        echo "example2: $i $(date)" >> /var/log/example2.log;
        i=$((i+1));
        sleep 1;
      done
    volumeMounts:
    - name: varlog
      mountPath: /var/log
  - name: sidecar1
    image: localhost:32000/sumologic/tailing-sidecar:demo
    env:
    - name: PATH_TO_TAIL
      value: /var/log/example1.log
    - name: LOG_LEVEL
      value: warning
    volumeMounts:
    - name: varlog
      mountPath: /var/log
    - name: volume-sidecar1
      mountPath: /tailing-sidecar/var
  - name: sidecar2
    image: localhost:32000/sumologic/tailing-sidecar:demo
    env:
    - name: PATH_TO_TAIL
      value: /var/log/example2.log
    - name: LOG_LEVEL
      value: warning
    volumeMounts:
    - name: varlog
      mountPath: /var/log
    - name: volume-sidecar2
      mountPath: /tailing-sidecar/var
  volumes:
  - name: varlog
    emptyDir: {}
  - name: volume-sidecar1
    hostPath:
      path: /var/log/sidecar1
      type: DirectoryOrCreate
  - name: volume-sidecar2
    hostPath:
      path: /var/log/sidecar2
      type: DirectoryOrCreate
```

Notice that tailing sidecar containers are configured through two environmental variables:

- `PATH_TO_TAIL` - pattern specifying a log file or multiple ones through the use of common wildcards,
  multiple patterns separated by commas are also allowed
- `LOG_LEVEL` - verbosity level, by default 'warning' is set,
  allowed values: error, warning, info, debug, trace

Try it!

```bash
kubectl apply -f sidecar/examples/pod_with_tailing_sidecars.yaml
```

And check logs:

```bash
$ kubectl logs example-with-tailling-sidecars sidecar1
example1: 0 Wed Jan 27 11:59:28 UTC 2021
example1: 1 Wed Jan 27 11:59:29 UTC 2021
example1: 2 Wed Jan 27 11:59:30 UTC 2021
example1: 3 Wed Jan 27 11:59:31 UTC 2021
...
```

```bash
$ kubectl logs example-with-tailling-sidecars sidecar2
example2: 0 Wed Jan 27 11:59:28 UTC 2021
example2: 1 Wed Jan 27 11:59:29 UTC 2021
example2: 2 Wed Jan 27 11:59:30 UTC 2021
example2: 3 Wed Jan 27 11:59:31 UTC 2021
...
```

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

## Testing in Vagrant environment

### Prerequisites

Please install the following:

- [VirtualBox](https://www.virtualbox.org/)
- [Vagrant](https://www.vagrantup.com/)
- [vagrant-disksize](https://github.com/sprotheroe/vagrant-disksize) plugin

### Setting up

Start and provision the Vagrant environment:

```bash
vagrant up
```

Connect to virtual machine:

```bash
vagrant ssh
```

### Build and run tailing sidecar

Build and push docker image to local container registry:

```bash
/tailing-sidecar/sidecar/Makefile
```

Run example Pod:

```bash
kubectl apply -f /tailing-sidecar/sidecar/examples/pod_with_tailing_sidecars.yaml
```
