# sumologic-tailing-sidecar

**tailing sidecar** is a cluster-level logging agent for Kubernetes using
[Fluent Bit](https://fluentbit.io/) as the underlying logging component.

For more information about cluster-level logging architecture please read Kubernetes
[documentation](https://kubernetes.io/docs/concepts/cluster-administration/logging/#cluster-level-logging-architectures).

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

### License

This project is released under the [Apache 2.0 License](LICENSE).

### Contributing

Please share your thoughts about sumologic-tailing-sidecar by opening an issue.

To get started contributing, please refer to our [Contributing](CONTRIBUTING.md) documentation.
