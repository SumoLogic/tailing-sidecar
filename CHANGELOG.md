# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.5.3] - 2022-09-12

- feat: add scc configuration [#381]

[v0.5.3]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.5.3
[#381]: https://github.com/SumoLogic/tailing-sidecar/pull/381

## [v0.3.4] - 2022-09-12

- fix: fix descriptions and summaries in catalog.redhat.com [#364]
- fix(chart): add permissions for leases [#371]
- feat: enable privileged mode for container [#377]
- feat: add scc configuration [#379]

[v0.3.4]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.3.4
[#364]: https://github.com/SumoLogic/tailing-sidecar/pull/364
[#371]: https://github.com/SumoLogic/tailing-sidecar/pull/371
[#377]: https://github.com/SumoLogic/tailing-sidecar/pull/377
[#379]: https://github.com/SumoLogic/tailing-sidecar/pull/379

## [v0.3.3] - 2022-07-29

- chore: replace deprecated APIs and update dependencies to support Kubernetes 1.23 [#351]
  - change apiextensions.k8s.io/v1beta1 to apiextensions.k8s.io/v1
  - change admissionregistration.k8s.io/v1beta1 to admissionregistration.k8s.io/v1
  - change k8s.io/api/admission/v1beta1 to k8s.io/api/admission/v1
  - update sigs.k8s.io/controller-tools/cmd/controller-gen to 0.4.1
  - update sigs.k8s.io/kustomize/kustomize to 3.8.3
  - update sigs.k8s.io/controller-runtime to v0.8.3
  - update k8s.io/apimachinery to v0.20.2
  - update k8s.io/client-go to v0.20.2
  - update k8s.io/api to v0.20.2
  - update gomodules.xyz/jsonpatch/v2 to 2.1.0
  - update github.com/go-logr/logr v0.3.0
  - adjust code to updated dependencies
  - change Kuberentes version in Vagrant environment to 1.23
  - change Cert Manager version in Vagrant environment to v1.5.0

[v0.3.3]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.3.3
[#351]: https://github.com/SumoLogic/tailing-sidecar/pull/351

## [v0.5.2] - 2022-07-19

- fix(chart): add permissions for leases [#345]
- fix(operator): prevent the `Failed to prepare volume` error logs [#347]
- chore: upgrade Fluent Bit from 1.8.12 to 1.9.6
- chore: upgrade kube-rbac-proxy from 0.5.0 to 0.11.0 [#280]
- chore: change container repository for kube-rbac-proxy to quay.io/brancz/kube-rbac-proxy [#280]
- chore: upgrade Golang from 1.17.6 to 1.18.4

[v0.5.2]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.5.2
[#280]: https://github.com/SumoLogic/tailing-sidecar/pull/280
[#345]: https://github.com/SumoLogic/tailing-sidecar/pull/345
[#347]: https://github.com/SumoLogic/tailing-sidecar/pull/347

## [v0.5.1] - 2022-02-02

- chore: update Golang from 1.16.5 to 1.17.6 [#212] [#213] [#231] [#232] [#247] [#248] [#256] [#257]
- chore: update Fluent Bit from 1.7.8 to 1.8.12 [#249] [#261]
- chore: update k8s.io/apimachinery from 0.21.1 to 0.22.4 [#215] [#237]

[v0.5.1]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.5.1
[#212]: https://github.com/SumoLogic/tailing-sidecar/pull/212
[#213]: https://github.com/SumoLogic/tailing-sidecar/pull/213
[#231]: https://github.com/SumoLogic/tailing-sidecar/pull/232
[#232]: https://github.com/SumoLogic/tailing-sidecar/pull/232
[#247]: https://github.com/SumoLogic/tailing-sidecar/pull/247
[#248]: https://github.com/SumoLogic/tailing-sidecar/pull/248
[#256]: https://github.com/SumoLogic/tailing-sidecar/pull/256
[#257]: https://github.com/SumoLogic/tailing-sidecar/pull/257
[#249]: https://github.com/SumoLogic/tailing-sidecar/pull/249
[#261]: https://github.com/SumoLogic/tailing-sidecar/pull/261
[#215]: https://github.com/SumoLogic/tailing-sidecar/pull/215
[#237]: https://github.com/SumoLogic/tailing-sidecar/pull/237

## [v0.3.2] - 2022-02-02

- feat: add Dockerfile and Makefile target to build UBI based container image [#193]
- feat: add Dockerfile and Makefile target to build UBI based sidecar container image [#194]
- chore: upgrade Golang from 1.16.2 to 1.17.6 [#265]
- chore: upgrade Fluent Bit from 1.7.2 to 1.8.12 [#266]

[v0.3.2]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.3.2
[#193]: https://github.com/SumoLogic/tailing-sidecar/pull/193
[#194]: https://github.com/SumoLogic/tailing-sidecar/pull/194
[#265]: https://github.com/SumoLogic/tailing-sidecar/pull/265
[#266]: https://github.com/SumoLogic/tailing-sidecar/pull/266

## [v0.3.1] - 2021-07-05

- Container image for kube-rbac-proxy configurable in Helm Chart (#172)

## [v0.5.0] - 2021-06-14

- Changes in Custom Resource for Tailing Sidecar Operator
  - Rename TailingSidecar to TailingSidecarConfig (#142)
  - Rename SidecarConfig to SidecarSpec in TailingSidecarConfig (#144)
  - Rename Config to SidecarSpecs in TailingSidecarConfig  (#144)
  - Sidecar container name defined as key in SidecarSpecs (#145)
  - Add PodSelector to TailingSidecarConfig (#146)
  - Add per tailing sidecar container annotations (#147)

- Replace hostPath volume added to tailing sidecars with emptytDir volume (#160)

- Make kube-rbac-proxy image configurable in values.yaml (#161)

- Replace deprecated APIs (#152)
  - Change apiextensions.k8s.io/v1beta1 to apiextensions.k8s.io/v1
  - Change admissionregistration.k8s.io/v1beta1 to admissionregistration.k8s.io/v1
  - Change cert-manager.io/v1alpha2 to cert-manager.io

## [v0.4.0] - 2021-05-11

- Change prefix for default tailing sidecar container name to "taling-sidecar-" (#112)
- Change prefix for tailing sidecar volume name to "volume-sidecar-" (#114)
- Changes in Custom Resource for Tailing Sidecar Operator
  - Rename "file" to "path"' in CRD definition (#117)
  - Rename "volume" to 'volumeMount' (#118)
  - Change type for 'volumeMount' from "string" to "VolumeMount" (#119)
- Set default tag for container images to `.Chart.AppVersion`(#111)
- Change `reinvocationPolicy` for Mutating Webhook to `Never` from `ifNeeded` (#119)
- Add startupProbe and livenessProbe for webhook server (#124, #125)
- Add explicit non blocking handling of Pod deletion (#122)
- Add tests for update of resources modified by operator (#123)
- Add metadata for Helm Chart (#110)

## [v0.3.0] - 2021-03-23

- Expose mutating webhook configuration in Helm Chart #107:
  - failurePolicy
  - reinvocationPolicy
  - objectSelector
  - namespaceSelector

## [v0.2.0] - 2021-03-22

- Make cert-manager as an optional dependency for Helm Chart #88
- Rename defined templates in Helm Chart #101
- Set proper versions for Helm chart and container images #96
- Improve documentation #93, #95, #98, #100

## [v0.1.0] - 2021-03-19

Initial version of tailing sidecar

[v0.3.1]: https://github.com/SumoLogic/tailing-sidecar/releases/tag/v0.3.1
[v0.5.0]: https://github.com/SumoLogic/tailing-sidecar/releases/tag/v0.5.0
[v0.4.0]: https://github.com/SumoLogic/tailing-sidecar/releases/tag/v0.4.0
[v0.3.0]: https://github.com/SumoLogic/tailing-sidecar/releases/tag/v0.3.0
[v0.2.0]: https://github.com/SumoLogic/tailing-sidecar/releases/tag/v0.2.0
[v0.1.0]: https://github.com/SumoLogic/tailing-sidecar/releases/tag/v0.1.0
