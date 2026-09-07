# Support for sidecar configuration

## Why?

For now only Sidecar image can be specified during configuration of Tailing Sidecar Operator.
This comes with one, but problematic requirement. In order to use custom Sidecar provider or custom configuration,
custom image has to be built and provided.

Example issues which could be resolved by exposing configuration:

- [Increase buffer size for fluent bit - lines are too long, skipping file](https://github.com/SumoLogic/tailing-sidecar/issues/226)
- [Allow tailing directory instead of specific file](https://github.com/SumoLogic/tailing-sidecar/issues/276)

## Implementation

### Configuration options

In order to resolve this issue and following the planned architecture ([Expose fluent-bit configurations](https://github.com/SumoLogic/tailing-sidecar/issues/30))
We would like to expose the following properties of Sidecar containers:

- `image` - this is already exposed
- `commands` - this is necessary if someone would like to modify arguments of the Sidecar provider
- `configuration` - we would like to expose configuration as Config Map name, so everyone could override it

  There is one technical issue. We cannot mount configMap from different namespace than the application we are going to get logs from.
  Possible solution is to allow Operator to create configmaps in other namespaces than it is installed.
  It would require customer to create one configmap in Operator Namespace and Operator would take care of rest.

  This should consists of the following suboptions:

  - `configMapName` - name of Config Map where the configuration is
  - `mountPath` - path to the directory where the configuration should be mounted

  We may consider different options of providing configuration (Volumes, Secrets)

- `resources` - every Sidecar provider can require different resources to work correctly

### Providing Configuation

In order to provide default override values for operator I propose to add to Operator support for configuartion file.
In helm Chart it would be manged by configMap, and it would make it easier to handle complex configuartion structures
like `commands` and `resoures`.

### Extending TalingSidecarConfig Custom Resource Definition

I believe it should be possible to specify conifiguration per TailingSidecarConfig, but I would leave it for future.

### Steps

The following steps are required to implement this functionality:

- [ ] Support configuration via config file

  - [ ] in Operator
  - [ ] in helm chart

- [ ] Add support for `configuration`

  - [ ] in Operator
  - [ ] in helm chart

- [ ] Add support for `commands`

  - [ ] in Operator
  - [ ] in helm chart

- [ ] Add support for `resources`

  - [ ] in Operator
  - [ ] in helm chart
