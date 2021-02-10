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
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	admv1 "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestPodExtender(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "PodExtender Suite")
}

var _ = Describe("PodExtender", func() {

	Context("PodExtender Handler", func() {

		decoder, err := admission.NewDecoder(runtime.NewScheme())
		It("creates decoder without any errors", func() {
			Expect(err).To(BeNil())

		})

		podExtender := PodExtender{
			Client:              testclient.NewFakeClient(),
			TailingSidecarImage: "tailing-sidecar-image:test",
			decoder:             decoder,
		}

		When("When request does not contain any object", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(``),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("rejects request as decoder returns an error", func() {
				Expect(resp.Allowed).To(BeFalse())
				Expect(resp.Patch).To(BeEmpty())
				Expect(resp.Result.Code).Should(Equal(int32(http.StatusBadRequest)))
			})
		})

		When("When request contains empty json", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{}`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)

			It("returns empty patch as extendPod function didn't find tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patch).To(BeEmpty())
			})
		})

		When("When Pod with null metadata is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": null,
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox"
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("eturns empty patch as extendPod function didn't find tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patch).To(BeEmpty())
			})
		})

		When("When Pod with empty metadata is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox",
										  "args": [
											"sleep",
											"1000000"
										  ]
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("returns empty patch as extendPod function didn't find tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patch).To(BeEmpty())
			})
		})

		When("When Pod with empty annotation is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {
										"name": "simple",
									  	"annotations": {}
									},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox"
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("returns empty patch as extendPod function didn't find tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patch).To(BeEmpty())
			})
		})

		When("When Pod with null annotation is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {
										"name": "simple",
									  	"annotations": null
									},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox"
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("returns empty patch as extendPod function didn't find tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patch).To(BeEmpty())
			})
		})
	})
})
