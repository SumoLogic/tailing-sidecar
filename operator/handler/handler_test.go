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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/nsf/jsondiff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gomodules.xyz/jsonpatch/v2"

	tailingsidecarv1 "github.com/SumoLogic/tailing-sidecar/operator/api/v1"
	admv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestPodExtender(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PodExtender Suite")
}

var _ = Describe("handler", func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	mountPropagationBidirectional := corev1.MountPropagationBidirectional

	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}
	cfg, err := testEnv.Start()
	It("starts test environment", func() {
		Expect(err).ToNot(HaveOccurred())
		Expect(cfg).ToNot(BeNil())
	})

	ctx := context.Background()

	Context("PodExtender.Handle without TailingSidecarConfig CRD installed", func() {
		k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
		It("creates a new client", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		decoder := admission.NewDecoder(scheme.Scheme)
		It("creates decoder without any errors", func() {
			Expect(err).ToNot(HaveOccurred())

		})

		podExtender := PodExtender{
			Client:              k8sClient,
			TailingSidecarImage: "tailing-sidecar-image:test",
			Decoder:             decoder,
		}

		When("Pod with raw configuration is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "annotations": {
								"tailing-sidecar": "varlog:/var/log/example0.log;varlog:/var/log/example1.log"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								   "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns empty patch and Internal Server Error as TailingSidecarConfig CRD needs to be available", func() {
				Expect(resp.Allowed).To(BeFalse())
				Expect(resp.Patches).To(BeEmpty())
				Expect(resp.Result.Code).Should(Equal(int32(http.StatusInternalServerError)))
			})
		})
	})

	Context("PodExtender.Handle", func() {
		err = tailingsidecarv1.AddToScheme(scheme.Scheme)
		It("adds TailingSidecarConfig to scheme", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
		It("creates a new client", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		decoder := admission.NewDecoder(scheme.Scheme)
		It("creates decoder without any errors", func() {
			Expect(err).ToNot(HaveOccurred())

		})

		podExtender := PodExtender{
			Client:              k8sClient,
			TailingSidecarImage: "tailing-sidecar-image:test",
			Decoder:             decoder,
		}

		podExtenderWithConfiguration := PodExtender{
			Client:              k8sClient,
			TailingSidecarImage: "tailing-sidecar-image:test",
			Decoder:             decoder,
			ConfigMapName:       "my-config-map",
			ConfigMountPath:     "my-custom-path",
			ConfigMapNamespace:  "tailing-sidecar-system",
		}

		namespace1 := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "tailing-sidecar-system",
			},
		}

		err = k8sClient.Create(ctx, namespace1)
		It("creates the first namespace", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		namespace2 := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "tailing-sidecar-system-different",
			},
		}

		err = k8sClient.Create(ctx, namespace2)
		It("creates the second namespace", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		When("request does not contain any object", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(``),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("rejects request as decoder returns an error", func() {
				Expect(resp.Allowed).To(BeFalse())
				Expect(resp.Patches).To(BeEmpty())
				Expect(resp.Result.Code).Should(Equal(int32(http.StatusBadRequest)))
			})
		})

		When("request contains empty json", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{}`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)

			It("returns empty patch as there is missing tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).To(BeEmpty())
			})
		})

		When("Pod with null metadata is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": null,
									"status": {},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox",
										  "resources": {}
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns empty patch as there is missing tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).To(BeEmpty())
			})
		})

		When("Pod with empty metadata is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {},
									"status": {},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox",
										  "resources": {}
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns empty patch as there is missing tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).To(BeEmpty())
			})
		})

		When("Pod with empty annotation is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {
										"creationTimestamp": null,
										"name": "pod-empty-annotations",
									  	"annotations": {}
									},
									"status": {},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox",
										  "resources": {}
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns empty patch as there is missing tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).To(BeEmpty())
			})
		})

		When("Pod with null annotation is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {
										"creationTimestamp": null,
										"name": "pod-with-null-annotation",
									  	"annotations": null
									},
									"status": {},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox",
										  "resources": {}
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("returns empty patch as there is missing tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).To(BeEmpty())
			})
		})

		When("Pod with ':' in tailing-sidecar annotation is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {
										"creationTimestamp": null,
										"name": "pod-with-colon-in-annotation",
										  "annotations": {
											"tailing-sidecar": ":"
										  }
									},
									"status": {},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox",
										  "resources": {}
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("returns patch empty patch", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).To(BeEmpty())
			})
		})

		When("Pod with empty string tailing-sidecar annotation is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {
										"creationTimestamp": null,
										"name": "pod-with-empty-string-in-annotation",
										  "annotations": {
											"tailing-sidecar": ""
										  }
									},
									"status": {},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox",
										  "resources": {}
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("returns patch empty patch", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).To(BeEmpty())
			})
		})

		When("Pod with TailingSidecarConfig without PodSelector", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system-different",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"sidecar": {
							Path: "/varconfig/log/example2.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig",
								MountPath: "/varconfig/log",
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "annotations": {
								"tailing-sidecar": "varlog:/var/log/example0.log;varlog:/var/log/example1.log"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								  "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns patch with tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())

				expectedPatches := loadJSONPatches("testdata/patch_with_2_tailing_sidecars.json")

				Expect(len(resp.Patches)).Should(Equal(len(expectedPatches)))

				for _, patch := range resp.Patches {
					Expect(isExpectedPatch(expectedPatches, patch)).To(BeTrue(), "cannot find patch in expected patches, patch: %+v", patch)
				}
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("Pod with TailingSidecarConfig with configuration in the same namespace without PodSelector", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system-different",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"sidecar": {
							Path: "/varconfig/log/example2.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig",
								MountPath: "/varconfig/log",
							},
						},
					},
				},
			}

			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-config-map",
					Namespace: "tailing-sidecar-system",
				},
			}

			err = k8sClient.Create(ctx, configMap)
			It("creates an exemplar sidecar configuration for operator", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Namespace: "tailing-sidecar-system",
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "annotations": {
								"tailing-sidecar": "varlog:/var/log/example0.log;varlog:/var/log/example1.log"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								  "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtenderWithConfiguration.Handle(ctx, request)
			It("returns patch with tailing sidecar containers for custom configuration handler", func() {
				fmt.Printf("%v\n", resp)
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())

				expectedPatches := loadJSONPatches("testdata/patch_with_2_tailing_sidecars_with_configuration.json")

				Expect(len(resp.Patches)).Should(Equal(len(expectedPatches)))

				for _, patch := range resp.Patches {
					Expect(isExpectedPatch(expectedPatches, patch)).To(BeTrue(), "cannot find patch in expected patches, patch: %+v", patch)
				}
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			err = k8sClient.Delete(ctx, configMap)
			It("deletes configMap", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("Pod with TailingSidecarConfig with configuration in different namespace without PodSelector", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system-different",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"sidecar": {
							Path: "/varconfig/log/example2.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig",
								MountPath: "/varconfig/log",
							},
						},
					},
				},
			}

			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-config-map",
					Namespace: "tailing-sidecar-system",
				},
			}

			err = k8sClient.Create(ctx, configMap)
			It("creates an exemplar sidecar configuration for operator", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Namespace: "tailing-sidecar-system-different",
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system-different",
							  "annotations": {
								"tailing-sidecar": "varlog:/var/log/example0.log;varlog:/var/log/example1.log"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								  "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtenderWithConfiguration.Handle(ctx, request)
			It("returns patch with tailing sidecar containers for custom configuration handler", func() {
				fmt.Printf("%v\n", resp)
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())

				expectedPatches := loadJSONPatches("testdata/patch_with_2_tailing_sidecars_with_configuration.json")

				Expect(len(resp.Patches)).Should(Equal(len(expectedPatches)))

				for _, patch := range resp.Patches {
					Expect(isExpectedPatch(expectedPatches, patch)).To(BeTrue(), "cannot find patch in expected patches, patch: %+v", patch)
				}
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			err = k8sClient.Delete(ctx, configMap)
			It("deletes configMap", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			configMap.Namespace = "tailing-sidecar-system-different"
			err = k8sClient.Delete(ctx, configMap)
			It("deletes created configMap", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("Pod with existing tailing sidecar", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"tailing-sidecar": "true",
						},
					},
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"tailing-sidecar-0": {
							Path: "/varconfig/log/example2.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig",
								MountPath: "/varconfig/log",
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "labels": {
								"tailing-sidecar": "true"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								  "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								},
								{
									"name": "tailing-sidecar-0",
									"image": "tailing-sidecar-image:test",
									"resources": {},
									"env": [
									  {
										"name": "PATH_TO_TAIL",
										"value": "/varconfig/log/example2.log"
									  },
									  {
										"name": "TAILING_SIDECAR",
										"value": "true"
									  }
									],
									"volumeMounts": [
									  {
										"name": "varlogconfig",
										"mountPath": "/varconfig/log"
									  },
									  {
										"name": "volume-sidecar-0",
										"mountPath": "/tailing-sidecar/var"
									  }
									]
								  }
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								},
								{
									"name": "volume-sidecar-0",
									"emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns patch without tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).To(BeEmpty())
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})

		})

		When("Pod with tailing sidecar configuration containing missing volume", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"tailing-sidecar": "true",
						},
					},
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"sidecar": {
							Path: "/varconfig/log/example2.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig-non-existing",
								MountPath: "/varconfig/log",
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "labels": {
								"tailing-sidecar": "true"
							  },
							  "annotations": {
								"tailing-sidecar": "varlog:/var/log/example0.log;varlog:/var/log/example1.log"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								  "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns patch with tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())

				expectedPatches := loadJSONPatches("testdata/patch_with_2_tailing_sidecars.json")

				Expect(len(resp.Patches)).Should(Equal(len(expectedPatches)))

				for _, patch := range resp.Patches {
					Expect(isExpectedPatch(expectedPatches, patch)).To(BeTrue(), "cannot find patch in expected patches, patch: %+v", patch)
				}
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("Pod with configuration in TailingSidecarConfigs in different namespaces", func() {
			tailingSidecar1 := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-1",
					Namespace: "tailing-sidecar-system-different",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"tailing-sidecar-0": "true",
						},
					},
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"sidecar-1": {
							Path: "/varconfig/log/example2.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig",
								MountPath: "/varconfig/log",
								ReadOnly:  true,
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar1)
			It("creates the first Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			tailingSidecar2 := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-2",
					Namespace: "tailing-sidecar-system-different",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"tailing-sidecar-1": "true",
						},
					},
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"sidecar-2": {
							Path: "/var/log/example1.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlog",
								MountPath: "/var/log",
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar2)
			It("creates the second Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "labels": {
								"tailing-sidecar-0": "true",
								"tailing-sidecar-1": "true"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								  "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns patch with tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())

				expectedPatches := loadJSONPatches("testdata/patch_with_2_tailing_sidecars_different_namespace.json")

				Expect(len(resp.Patches)).Should(Equal(len(expectedPatches)))

				for _, patch := range resp.Patches {
					Expect(isExpectedPatch(expectedPatches, patch)).To(BeTrue(), "cannot find patch in expected patches, patch: %+v", patch)
				}
			})

			err = k8sClient.Delete(ctx, tailingSidecar1)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			err = k8sClient.Delete(ctx, tailingSidecar2)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("Pod with raw and predefined configurations is created", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"tailing-sidecar": "true",
						},
					},
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"sidecar-0": {
							Path: "/varconfig/log/example2.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig",
								MountPath: "/varconfig/log",
								ReadOnly:  true,
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "labels": {
								"tailing-sidecar": "true"
							  },
							  "annotations": {
								"tailing-sidecar": "varlog:/var/log/example0.log;varlog:/var/log/example1.log"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								   "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns patch with tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())

				expectedPatches := loadJSONPatches("testdata/patch_with_3_tailing_sidecars_raw_and_predefined.json")

				Expect(len(resp.Patches)).Should(Equal(len(expectedPatches)))

				for _, patch := range resp.Patches {
					Expect(isExpectedPatch(expectedPatches, patch)).To(BeTrue(), "cannot find patch in expected patches, patch: %+v", patch)
				}
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})

		})

		When("Pod with named sidecars", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"tailing-sidecar": "true",
						},
					},
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"test-container-2": {
							Path: "/varconfig/log/example2.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig",
								MountPath: "/varconfig/log",
							},
							Annotations: map[string]string{
								"sourceCategory": "sourceCategory-1",
								"annotation-1":   "annotation-1",
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "labels": {
								"tailing-sidecar": "true"
							  },
							  "annotations": {
								"tailing-sidecar": "test-container-3:varlog:/var/log/example0.log;test-container-1:varlog:/var/log/example1.log"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								   "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns patch with tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())

				expectedPatches := loadJSONPatches("testdata/patch_with_3_named_tailing_sidecars.json")

				Expect(len(resp.Patches)).Should(Equal(len(expectedPatches)))

				for _, patch := range resp.Patches {
					Expect(isExpectedPatch(expectedPatches, patch)).To(BeTrue(), "cannot find patch in expected patches, patch: %+v", patch)
				}
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("Pod with named and not named sidecars", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"tailing-sidecar": "true",
						},
					},
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"test-container-2": {
							Path: "/varconfig/log/example2.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig",
								MountPath: "/varconfig/log",
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "labels": {
								"tailing-sidecar": "true"
							  },
							  "annotations": {
								"tailing-sidecar": "test-container-0:varlog:/var/log/example0.log;varlog:/var/log/example1.log"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								   "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns patch with tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())

				expectedPatches := loadJSONPatches("testdata/patch_with_3_named_not_named_tailing_sidecars.json")

				Expect(len(resp.Patches)).Should(Equal(len(expectedPatches)))

				for _, patch := range resp.Patches {
					Expect(isExpectedPatch(expectedPatches, patch)).To(BeTrue(), "cannot find patch in expected patches, patch: %+v", patch)
				}
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("Update Pod with one named sidecars and add not named", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"tailing-sidecar": "true",
						},
					},
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"test-container": {
							Path: "/varconfig/log/example0.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig",
								MountPath: "/varconfig/log",
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "labels": {
								"tailing-sidecar": "true"
							  },
							  "annotations": {
								"tailing-sidecar": "varlog:/var/log/example1.log"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								   "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								},
								{
									"name": "test-container",
									"image": "busybox",
									"resources": {},
									"env": [
										{
											"name": "PATH_TO_TAIL",
											"value": "/varconfig/log/example0.log"
										},
										{
											"name": "TAILING_SIDECAR",
											"value": "true"
										}
									],
									"volumeMounts": [
									  {
										"mountPath": "/tailing-sidecar/var",
										"name": "volume-sidecar-0"
									  },
									  {
										"name": "varlogconfig",
										"mountPath": "/varconfig/log"
									  }
									]
								  }
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								},
								{
								  "name": "volume-sidecar-0",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns patch with tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())

				expectedPatches := loadJSONPatches("testdata/patch_update_1_tailing_sidecar.json")

				Expect(len(resp.Patches)).Should(Equal(len(expectedPatches)))

				for _, patch := range resp.Patches {
					Expect(isExpectedPatch(expectedPatches, patch)).To(BeTrue(), "cannot find patch in expected patches, patch: %+v", patch)
				}
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("Pod with configurations containing the same names for tailing sidecar containers", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"tailing-sidecar": "true",
						},
					},
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"test-container": {
							Path: "/varconfig/log/example2.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig",
								MountPath: "/varconfig/log",
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a TailingSidecarConfig with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "labels": {
								"tailing-sidecar": "true"
							  },
							  "annotations": {
								"tailing-sidecar": "test-container:varlogconfig:/varconfig/log/example2.log"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								   "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns empty patch", func() {
				Expect(resp.Allowed).To(BeFalse())
				Expect(resp.Patches).To(BeEmpty())
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("Pod with configuration containing name of existing container", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"tailing-sidecar": "true",
						},
					},
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"test-container": {
							Path: "/varconfig/log/example2.log",
							VolumeMount: corev1.VolumeMount{
								Name:      "varlogconfig",
								MountPath: "/varconfig/log",
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "labels": {
								"tailing-sidecar": "true"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "test-container",
								  "image": "busybox",
								   "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns patch with tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeFalse())
				Expect(resp.Patches).To(BeEmpty())
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("Remove all tailing sidecars", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "annotations": {
								"tailing-sidecar": ""
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								   "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								},
								{
									"name": "tailing-sidecar-0",
									"image": "busybox",
									"resources": {},
									"env": [
										{
											"name": "PATH_TO_TAIL",
											"value": "/varconfig/log/example0.log"
										},
										{
											"name": "TAILING_SIDECAR",
											"value": "true"
										}
									],
									"volumeMounts": [
									  {
										"mountPath": "/tailing-sidecar/var",
										"name": "volume-sidecar-0"
									  },
									  {
										"name": "varlogconfig",
										"mountPath": "/varconfig/log"
									  }
									]
								  }
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								},
								{
								  "name": "volume-sidecar-0",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns patch with tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())

				expectedPatches := loadJSONPatches("testdata/patch_remove_tailing_sidecar.json")

				Expect(len(resp.Patches)).Should(Equal(len(expectedPatches)))

				for _, patch := range resp.Patches {
					Expect(isExpectedPatch(expectedPatches, patch)).To(BeTrue(), "cannot find patch in expected patches, patch: %+v", patch)
				}
			})
		})

		When("Update Pod and change volumeMount configuration", func() {
			tailingSidecar := &tailingsidecarv1.TailingSidecarConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tailing-sidecar-in-pod-namespace",
					Namespace: "tailing-sidecar-system",
				},
				Spec: tailingsidecarv1.TailingSidecarConfigSpec{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"tailing-sidecar": "true",
						},
					},
					SidecarSpecs: map[string]tailingsidecarv1.SidecarSpec{
						"test-container": {
							Path: "/varconfig/log/example0.log",
							VolumeMount: corev1.VolumeMount{
								Name:             "varlogconfig",
								MountPath:        "/varconfig/log",
								ReadOnly:         true,
								MountPropagation: &mountPropagationBidirectional,
							},
						},
					},
				},
			}

			err = k8sClient.Create(ctx, tailingSidecar)
			It("creates a Tailingsidecar with configuration", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "creationTimestamp": null,
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "labels": {
								"tailing-sidecar": "true"
							  }
							},
							"status": {},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								   "resources": {},
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								},
								{
									"name": "test-container",
									"image": "tailing-sidecar-image:test",
									"resources": {},
									"env": [
										{
											"name": "PATH_TO_TAIL",
											"value": "/varconfig/log/example0.log"
										},
										{
											"name": "TAILING_SIDECAR",
											"value": "true"
										}
									],
									"volumeMounts": [
									  {
										"name": "varlogconfig",
										"mountPath": "/varconfig/log"
									  },
									  {
										"mountPath": "/tailing-sidecar/var",
										"name": "volume-sidecar-0"
									  }
									]
								  }
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								},
								{
								  "name": "volume-sidecar-0",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(ctx, request)
			It("returns patch with tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())

				expectedPatches := loadJSONPatches("testdata/patch_update_volume.json")

				Expect(len(resp.Patches)).Should(Equal(len(expectedPatches)))

				for _, patch := range resp.Patches {
					Expect(isExpectedPatch(expectedPatches, patch)).To(BeTrue(), "cannot find patch in expected patches, patch: %+v", patch)
				}
			})

			err = k8sClient.Delete(ctx, tailingSidecar)
			It("deletes TailingSidecarConfig", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		err = k8sClient.Delete(ctx, namespace1)
		It("deletes the first Namespace", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		err = k8sClient.Delete(ctx, namespace2)
		It("deletes the second Namespace", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

func isExpectedPatch(expectedPatches []jsonpatch.JsonPatchOperation, patch jsonpatch.JsonPatchOperation) bool {
	diffOpts := jsondiff.DefaultConsoleOptions()
	jsonPatch, err := json.MarshalIndent(patch, "", "  ")
	Expect(err).ToNot(HaveOccurred())

	for _, expectedPatch := range expectedPatches {
		jsonExpectedPatch, err := json.MarshalIndent(expectedPatch, "", "  ")
		Expect(err).ToNot(HaveOccurred())
		res, _ := jsondiff.Compare([]byte(jsonExpectedPatch), []byte(jsonPatch), &diffOpts)
		if res == jsondiff.FullMatch {
			return true
		}
	}
	return false
}

func loadJSONPatches(filePath string) []jsonpatch.JsonPatchOperation {
	jsonFromFile, err := ioutil.ReadFile(filePath)
	Expect(err).To(BeNil(), "error loading patches, file path: %s", filePath)
	expectedPatches := make([]jsonpatch.JsonPatchOperation, 0)
	err = json.Unmarshal([]byte(jsonFromFile), &expectedPatches)
	Expect(err).To(BeNil(), "cannot unmarshal patches, file path: %s", filePath)
	return expectedPatches
}
