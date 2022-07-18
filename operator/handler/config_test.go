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
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("config", func() {
	DescribeTable("convertTailingSidecarConfigs",
		func(
			input []tailingsidecarv1.TailingSidecarConfig,
			expectedOutputLength int,
		) {
			converted, err := convertTailingSidecarConfigs(input)

			Expect(err).NotTo(HaveOccurred())
			Expect(converted).To(HaveLen(expectedOutputLength))
		},

		Entry(
			"When there are no configs",
			[]tailingsidecarv1.TailingSidecarConfig{},
			0,
		),
		Entry(
			"When there is one config and one sidecar",
			[]tailingsidecarv1.TailingSidecarConfig{
				{
					Spec: tailingsidecarv1.TailingSidecarConfigSpec{
						SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
							"config1": {
								Path: "/var/log/file1.log",
								VolumeMount: corev1.VolumeMount{
									Name:      "logs-dir",
									MountPath: "/var/log",
								},
							},
						},
					},
				},
			},
			// actual output length: 2
			1,
		),
		Entry(
			"When there is one config and two sidecars",
			[]tailingsidecarv1.TailingSidecarConfig{
				{
					Spec: tailingsidecarv1.TailingSidecarConfigSpec{
						SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
							"config1": {
								Path: "/var/log/file1.log",
								VolumeMount: corev1.VolumeMount{
									Name:      "logs-dir",
									MountPath: "/var/log",
								},
							},
							"config2": {
								Path: "/var/log/file2.log",
								VolumeMount: corev1.VolumeMount{
									Name:      "logs-dir",
									MountPath: "/var/log",
								},
							},
						},
					},
				},
			},
			// actual output length: 3
			2,
		),
		Entry(
			"When there are two configs and two sidecars",
			[]tailingsidecarv1.TailingSidecarConfig{
				{
					Spec: tailingsidecarv1.TailingSidecarConfigSpec{
						SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
							"config1": {
								Path: "/var/log/file1.log",
								VolumeMount: corev1.VolumeMount{
									Name:      "logs-dir",
									MountPath: "/var/log",
								},
							},
						},
					},
				},
				{
					Spec: tailingsidecarv1.TailingSidecarConfigSpec{
						SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
							"config2": {
								Path: "/var/log/file2.log",
								VolumeMount: corev1.VolumeMount{
									Name:      "logs-dir",
									MountPath: "/var/log",
								},
							},
						},
					},
				},
			},
			// actual output length: 4
			2,
		),
		Entry(
			"When there are two configs and five sidecars",
			[]tailingsidecarv1.TailingSidecarConfig{
				{
					Spec: tailingsidecarv1.TailingSidecarConfigSpec{
						SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
							"config1": {
								Path: "/var/log/file1.log",
								VolumeMount: corev1.VolumeMount{
									Name:      "logs-dir",
									MountPath: "/var/log",
								},
							},
							"config2": {
								Path: "/var/log/file2.log",
								VolumeMount: corev1.VolumeMount{
									Name:      "logs-dir",
									MountPath: "/var/log",
								},
							},
						},
					},
				},
				{
					Spec: tailingsidecarv1.TailingSidecarConfigSpec{
						SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
							"config3": {
								Path: "/var/log/file3.log",
								VolumeMount: corev1.VolumeMount{
									Name:      "logs-dir",
									MountPath: "/var/log",
								},
							},
							"config4": {
								Path: "/var/log/file4.log",
								VolumeMount: corev1.VolumeMount{
									Name:      "logs-dir",
									MountPath: "/var/log",
								},
							},
							"config5": {
								Path: "/var/log/file5.log",
								VolumeMount: corev1.VolumeMount{
									Name:      "logs-dir",
									MountPath: "/var/log",
								},
							},
						},
					},
				},
			},
			// actual output length: 7
			5,
		),
	)
})
