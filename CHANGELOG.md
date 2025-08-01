# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.18.0] - 2025-07-18

- chore: update Go to 1.24.5 for fluentbit images [#807]
- chore: update Go to 1.24.5 for operator images [#808]

[#807]: https://github.com/SumoLogic/tailing-sidecar/pull/807
[#808]: https://github.com/SumoLogic/tailing-sidecar/pull/808

## [v0.17.0] - 2025-05-07

- chore: update base ubi image from 9.4 to 9.5 and golang version from 1.22.5 to 1.24.3 [#802](https://github.com/SumoLogic/tailing-sidecar/pull/802)

## [v0.16.0] - 2024-10-10

- build(deps): bump docker/setup-buildx-action from 3.4.0 to 3.7.1 [#745]
- build(deps): bump github.com/onsi/gomega in /operator [#757]
- build(deps): bump fluent/fluent-bit in /sidecar/fluentbit [#763]
- build(deps): bump golang from 1.22.4 to 1.23.2 in /operator [#764]

[#745]: https://github.com/SumoLogic/tailing-sidecar/pull/745
[#757]: https://github.com/SumoLogic/tailing-sidecar/pull/757
[#763]: https://github.com/SumoLogic/tailing-sidecar/pull/763
[#764]: https://github.com/SumoLogic/tailing-sidecar/pull/764

## [v0.15.0] - 2024-06-27

- chore(sidecar): update Fluent Bit UBI image to 3.0.7 by @sumo-drosiek in #725

[#725]: https://github.com/SumoLogic/tailing-sidecar/pull/725

## [v0.14.1] - 2024-06-21

- deps: upgrade kube-rbac-proxy to v0.18.0 [#723]

[#723]: https://github.com/SumoLogic/tailing-sidecar/pull/723

## [v0.14.0] - 2024-06-04

- build(deps): bump ubi8/ubi-minimal from 8.9 to 8.10 [#706]
- build(deps): bump fluent/fluent-bit from 3.0.4 to 3.0.6 [#708]
- chore: set version and release_number labels in ubi image [#711]
- build: downgrade sidecar to Go 1.20 [#710]
- deps: upgrade controller-runtime to 0.18.3 [#715]
- feat(helm): set tolerations for the operator [#716]

[#706]: https://github.com/SumoLogic/tailing-sidecar/pull/706
[#708]: https://github.com/SumoLogic/tailing-sidecar/pull/708
[#710]: https://github.com/SumoLogic/tailing-sidecar/pull/710
[#711]: https://github.com/SumoLogic/tailing-sidecar/pull/711
[#715]: https://github.com/SumoLogic/tailing-sidecar/pull/715
[#716]: https://github.com/SumoLogic/tailing-sidecar/pull/716

## [v0.13.0] - 2024-05-21

- build(deps): bump fluent/fluent-bit from 3.0.2 to 3.0.4 in /sidecar/fluentbit. [#703]
- build(deps): bump golang from 1.22.2 to 1.22.3 in /operator [#699]

[#703]: https://github.com/SumoLogic/tailing-sidecar/pull/703
[#699]: https://github.com/SumoLogic/tailing-sidecar/pull/699
[v0.13.0]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.13.0

## [v0.12.0] - 2024-04-19

- build(deps): update fluent/fluent-bit in /sidecar/fluentbit from from 3.0.1 to 3.0.2. [#680]
- feat: push container images to Docker Hub [#687]

[#680]: https://github.com/SumoLogic/tailing-sidecar/pull/680
[#687]: https://github.com/SumoLogic/tailing-sidecar/pull/687
[v0.12.0]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.12.0

## [v0.11.0] - 2024-03-25

- feat: build ubi images [#654]
- add affinity config for operator pod [#670]
- build(deps): bump Fluent Bit from 2.2.2 to 3.0.0 [#672]

[#654]: https://github.com/SumoLogic/tailing-sidecar/pull/654
[#670]: https://github.com/SumoLogic/tailing-sidecar/pull/670
[#672]: https://github.com/SumoLogic/tailing-sidecar/pull/672
[v0.11.0]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.11.0

## [v0.10.0] - 2024-02-14

- feat(operator): add proper health checks [#608]
- build(deps): bump Fluent Bit from 2.1.7 to 2.2.2 [#569] [#581] [#589] [#622] [#636] [#639]
- deps(operator): upgrade controller-runtime to 0.17.1 [#603] [#644] [#648]

[#569]: https://github.com/SumoLogic/tailing-sidecar/pull/569
[#581]: https://github.com/SumoLogic/tailing-sidecar/pull/581
[#589]: https://github.com/SumoLogic/tailing-sidecar/pull/589
[#603]: https://github.com/SumoLogic/tailing-sidecar/pull/603
[#608]: https://github.com/SumoLogic/tailing-sidecar/pull/608
[#622]: https://github.com/SumoLogic/tailing-sidecar/pull/622
[#636]: https://github.com/SumoLogic/tailing-sidecar/pull/636
[#639]: https://github.com/SumoLogic/tailing-sidecar/pull/639
[#644]: https://github.com/SumoLogic/tailing-sidecar/pull/644
[#648]: https://github.com/SumoLogic/tailing-sidecar/pull/648
[v0.10.0]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.10.0

## [v0.9.0] - 2023-07-20

- chore(sidecar): update fluent-bit to 2.1.7 [#563]
- feat(helm): add resource setting for kubeRbacProxy [#562]
- build(deps): bump golang from 1.20.5 to 1.20.6 in /sidecar [#559]
- build(deps): bump golang from 1.20.5 to 1.20.6 in /operator [#558]

[#558]: https://github.com/SumoLogic/tailing-sidecar/pull/558
[#559]: https://github.com/SumoLogic/tailing-sidecar/pull/559
[#562]: https://github.com/SumoLogic/tailing-sidecar/pull/562
[#563]: https://github.com/SumoLogic/tailing-sidecar/pull/563
[v0.9.0]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.9.0

## [v0.8.0] - 2023-07-07

- build(deps): bump fluent/fluent-bit from 2.0.10 to 2.0.11 in /sidecar [#512]
- build(deps): bump golang from 1.20.3 to 1.20.5 in /sidecar [#521] [#538]
- build(deps): bump golang from 1.20.3 to 1.20.5 in /operator [#520] [#539]
- feat(operator): add support for configuration file [#470]
- feat(helm): use configmap to configure tailing-sidecar-operator [#534]
- chore(sidecar): use debian bookworm for package base [#553]
- feat(operator): support config sidecar pod resources [#552]
- feat(operator): support config leader election [#550]
- feat(operator): allow to override sidecar configuration [#551]
- feat(operator): support config replicaCount [#555]

[#470]: https://github.com/SumoLogic/tailing-sidecar/pull/470
[#512]: https://github.com/SumoLogic/tailing-sidecar/pull/512
[#520]: https://github.com/SumoLogic/tailing-sidecar/pull/520
[#521]: https://github.com/SumoLogic/tailing-sidecar/pull/521
[#534]: https://github.com/SumoLogic/tailing-sidecar/pull/534
[#538]: https://github.com/SumoLogic/tailing-sidecar/pull/538
[#539]: https://github.com/SumoLogic/tailing-sidecar/pull/539
[#550]: https://github.com/SumoLogic/tailing-sidecar/pull/550
[#551]: https://github.com/SumoLogic/tailing-sidecar/pull/551
[#552]: https://github.com/SumoLogic/tailing-sidecar/pull/552
[#553]: https://github.com/SumoLogic/tailing-sidecar/pull/553
[#555]: https://github.com/SumoLogic/tailing-sidecar/pull/555
[v0.8.0]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.8.0

## [v0.7.0] - 2023-04-05

- build(deps): bump fluent/fluent-bit from 1.9.9 to 2.0.10 in /sidecar [#492]

[v0.7.0]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.7.0
[#492]: https://github.com/SumoLogic/tailing-sidecar/pull/492

## [v0.6.0] - 2023-03-23

- feat: add support for configuring livenessProbe and startupProbe [#484]
- chore: change container registry to public ECR [#500]

[v0.6.0]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.6.0
[#484]: https://github.com/SumoLogic/tailing-sidecar/pull/484
[#500]: https://github.com/SumoLogic/tailing-sidecar/pull/500

## [v0.5.6] - 2023-01-09

No significant changes since [v0.5.5]. This release contains ARM docker images.

[v0.5.6]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.5.6

## [v0.5.5] - 2022-11-03

- fix: remove update operation from pod from mutatingWebhookConfigurations [#413]

[v0.5.5]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.5.5
[#413]: https://github.com/SumoLogic/tailing-sidecar/pull/413

## [v0.5.4] - 2022-09-26

- chore: upgrade libc and zlib in the Dockerfile [#392]
- chore: upgrade client-go [#391]

[v0.5.4]: https://github.com/SumoLogic/tailing-sidecar/releases/v0.5.4
[#391]: https://github.com/SumoLogic/tailing-sidecar/pull/391
[#392]: https://github.com/SumoLogic/tailing-sidecar/pull/392

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
