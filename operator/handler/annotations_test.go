/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handler

import (
	tailingsidecarv1 "github.com/SumoLogic/tailing-sidecar/operator/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("annotations", func() {
	Context("Pod with non empty configuration for Tailing Sidecar Operator", func() {
		mountPropagationBidirectional := corev1.MountPropagationBidirectional

		When("configuration for annotations is not provided", func() {
			annotations := make(map[string]string, 0)
			config := sidecarConfig{
				name: "sidecar-0",
				spec: tailingsidecarv1.SidecarSpec{
					Path: "/varconfig/log/example2.log",
					VolumeMount: corev1.VolumeMount{
						Name:      "varlogconfig",
						MountPath: "/varconfig/log",
					},
				},
			}

			extendedAnnotations := addAnnotations(annotations, config)

			It("returns extended annotations", func() {
				Expect(len(extendedAnnotations)).Should(Equal(0))
			})

		})

		When("configuration for annotations is provided", func() {
			annotations := make(map[string]string, 0)
			config := sidecarConfig{
				name: "sidecar-1",
				spec: tailingsidecarv1.SidecarSpec{
					Path: "/var/log/example_dir/example.log",
					VolumeMount: corev1.VolumeMount{
						Name:             "example_volume",
						MountPath:        "/var/log/example_dir",
						ReadOnly:         true,
						SubPath:          "example_dir",
						MountPropagation: &mountPropagationBidirectional,
						SubPathExpr:      "example_dir",
					},
					Annotations: map[string]string{
						"sourceCategory": "sourceCategory-1",
						"annotation-1":   "annotatation-value-1",
					},
				},
			}

			extendedAnnotations := addAnnotations(annotations, config)

			It("returns extended annotations", func() {
				Expect(len(extendedAnnotations)).Should(Equal(2))

				Expect(extendedAnnotations["tailing-sidecar.sumologic.com/sidecar-1.annotation-1"]).
					Should(Equal("annotatation-value-1"))

				Expect(extendedAnnotations["tailing-sidecar.sumologic.com/sidecar-1.sourceCategory"]).
					Should(Equal("sourceCategory-1"))
			})

		})

		When("only sidecar name is provided", func() {
			annotations := make(map[string]string, 0)
			config := sidecarConfig{
				name: "sidecar-2",
			}

			extendedAnnotations := addAnnotations(annotations, config)

			It("returns extended annotations", func() {
				Expect(len(extendedAnnotations)).Should(Equal(0))
			})
		})

		When("custom prefix is provided and configuration for annotations is provided", func() {
			annotations := make(map[string]string, 0)
			config := sidecarConfig{
				annotationsPrefix: "customPrefix",
				name:              "sidecar-3",
				spec: tailingsidecarv1.SidecarSpec{
					Path: "/var/log/example_dir/example.log",
					VolumeMount: corev1.VolumeMount{
						Name:             "example_volume",
						MountPath:        "/var/log/example_dir",
						ReadOnly:         true,
						SubPath:          "example_dir",
						MountPropagation: &mountPropagationBidirectional,
						SubPathExpr:      "example_dir",
					},
					Annotations: map[string]string{
						"sourceCategory": "sourceCategory-3",
						"annotation-3":   "annotatation-value-3",
					},
				},
			}

			extendedAnnotations := addAnnotations(annotations, config)

			It("returns extended annotations", func() {
				Expect(len(extendedAnnotations)).Should(Equal(2))

				Expect(extendedAnnotations["customPrefix/sidecar-3.annotation-3"]).
					Should(Equal("annotatation-value-3"))

				Expect(extendedAnnotations["customPrefix/sidecar-3.sourceCategory"]).
					Should(Equal("sourceCategory-3"))
			})
		})

		When("custom prefix is provided and configuration for annotations is provided", func() {
			var annotations map[string]string
			config := sidecarConfig{
				annotationsPrefix: "customPrefix",
				name:              "sidecar-4",
				spec: tailingsidecarv1.SidecarSpec{
					Path: "/var/log/example_dir/example.log",
					VolumeMount: corev1.VolumeMount{
						Name:             "example_volume",
						MountPath:        "/var/log/example_dir",
						ReadOnly:         true,
						SubPath:          "example_dir",
						MountPropagation: &mountPropagationBidirectional,
						SubPathExpr:      "example_dir",
					},
					Annotations: map[string]string{
						"sourceCategory": "sourceCategory-4",
						"annotation-4":   "annotatation-value-4",
					},
				},
			}

			extendedAnnotations := addAnnotations(annotations, config)

			It("returns extended annotations", func() {
				Expect(len(extendedAnnotations)).Should(Equal(2))

				Expect(extendedAnnotations["customPrefix/sidecar-4.annotation-4"]).
					Should(Equal("annotatation-value-4"))

				Expect(extendedAnnotations["customPrefix/sidecar-4.sourceCategory"]).
					Should(Equal("sourceCategory-4"))
			})
		})
	})
})
