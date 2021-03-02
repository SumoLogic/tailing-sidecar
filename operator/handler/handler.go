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
	"net/http"
	"strings"

	tailingsidecarv1 "github.com/SumoLogic/tailing-sidecar/operator/api/v1"
	guuid "github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/add-tailing-sidecars-v1-pod,mutating=true,failurePolicy=ignore,groups="",resources=pods,verbs=create;update,versions=v1,name=tailing-sidecar.sumologic.com

const (
	sidecarEnvPath      = "PATH_TO_TAIL"
	sidecarEnvMarker    = "TAILING_SIDECAR"
	sidecarEnvMarkerVal = "true"

	sidecarContainerName   = "tailing-sidecar%d"
	sidecarContainerPrefix = "tailing-sidecar"

	hostPathDirPath    = "/var/log/tailing-sidecar-fluentbit/%s/%s"
	hostPathVolumeName = "volume-sidecar%d"
	hostPathMountPath  = "/tailing-sidecar/var"
)

var (
	handlerLog   = ctrl.Log.WithName("tailing-sidecar.operator.handler.PodExtender")
	hostPathType = corev1.HostPathDirectoryOrCreate
)

// PodExtender extends Pods by tailling sidecar containers
type PodExtender struct {
	Client              client.Client
	TailingSidecarImage string
	decoder             *admission.Decoder
}

// Handle handles requests to create/update Pod and extends it by adding tailing sidecars
func (e *PodExtender) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	err := e.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if _, ok := pod.ObjectMeta.Annotations[sidecarAnnotation]; !ok {
		return admission.Allowed("missing tailing-sidecar annotation, tailing sidecars not added")
	}

	handlerLog.Info("Handling request for Pod",
		"Name", pod.ObjectMeta.Name,
		"Namespace", pod.ObjectMeta.Namespace,
		"GenerateName", pod.ObjectMeta.GenerateName,
		"Operation", req.Operation,
	)

	if err := e.extendPod(ctx, pod); err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// InjectDecoder injects the decoder
func (e *PodExtender) InjectDecoder(d *admission.Decoder) error {
	e.decoder = d
	return nil
}

// extendPod extends Pod by adding tailing sidecars according to configuration in annotation
func (e PodExtender) extendPod(ctx context.Context, pod *corev1.Pod) error {
	// Get TailingSidecars from namespace
	tailingSidecarList := &tailingsidecarv1.TailingSidecarList{}
	tailingSidecarListOpts := []client.ListOption{
		client.InNamespace(pod.ObjectMeta.Namespace),
	}

	if err := e.Client.List(ctx, tailingSidecarList, tailingSidecarListOpts...); err != nil {
		handlerLog.Error(err,
			"Failed to get list of TailingSidecars in namespace",
			"namespace", pod.ObjectMeta.Namespace,
		)
		return err
	}

	// Join configurations from TailingSidecars
	tailingSidecarConfigs := joinTailinSidecarConfigs(tailingSidecarList.Items)

	// Parse configurations from annotation and join them with configurations from TailingSidecars
	configs := getConfigs(pod.ObjectMeta.Annotations, tailingSidecarConfigs)

	if len(configs) == 0 {
		handlerLog.Info("Missing configuration for Pod",
			"Name", pod.ObjectMeta.Name,
			"Namespace", pod.ObjectMeta.Namespace,
			"GenerateName", pod.ObjectMeta.GenerateName,
		)
		return nil
	}

	handlerLog.Info("Found configuration for Pod",
		"Pod Name", pod.ObjectMeta.Name,
		"Namespace", pod.ObjectMeta.Namespace,
		"GenerateName", pod.ObjectMeta.GenerateName,
	)

	containers := make([]corev1.Container, 0)
	hostPathDir := getHostPath(pod)
	sidecarID := len(getTailingSidecars(pod.Spec.Containers))

	for _, config := range configs {
		if isSidecarAvailable(pod.Spec.Containers, config) {
			// Do not add tailing sidecar if tailing sidecar with specific configuration exists
			handlerLog.Info("Tailing sidecar exists",
				"file", config.File,
				"volume", config.Volume,
				"container", config.Container,
			)
			continue
		}

		volume, err := getVolume(pod.Spec.Containers, config.Volume)
		if err != nil {
			handlerLog.Error(err,
				"Failed to find volume",
				"Pod Name", pod.ObjectMeta.Name,
				"Namespace", pod.ObjectMeta.Namespace,
				"GenerateName", pod.ObjectMeta.GenerateName,
			)
			continue
		}

		volumeName := fmt.Sprintf(hostPathVolumeName, sidecarID)
		if config.Container == "" {
			config.Container = fmt.Sprintf(sidecarContainerName, sidecarID)
		}

		hostPath := fmt.Sprintf("%s/%s", hostPathDir, config.Container)
		pod.Spec.Volumes = append(pod.Spec.Volumes,
			corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: hostPath,
						Type: &hostPathType,
					},
				},
			})

		container := corev1.Container{
			Image: e.TailingSidecarImage,
			Name:  config.Container,
			Env: []corev1.EnvVar{
				{
					Name:  sidecarEnvPath,
					Value: config.File,
				},
				{
					Name:  sidecarEnvMarker,
					Value: sidecarEnvMarkerVal,
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				volume,
				{
					Name:      volumeName,
					MountPath: hostPathMountPath,
				},
			},
		}
		containers = append(containers, container)
		sidecarID++
	}
	podContainers := removeDeletedSidecars(pod.Spec.Containers, configs)
	pod.Spec.Containers = append(podContainers, containers...)
	return nil
}

// removeDeletedSidecars removes deleted tailing sidecar containers from Pod specification
func removeDeletedSidecars(containers []corev1.Container, configs []tailingsidecarv1.SidecarConfig) []corev1.Container {
	podContainers := make([]corev1.Container, 0)
	for _, container := range containers {
		if !isSidecarEnvAvailable(container.Env, sidecarEnvMarker, sidecarEnvMarkerVal) {
			podContainers = append(podContainers, container)
		} else {
			for _, config := range configs {
				if ((config.Container == "" && strings.HasPrefix(container.Name, sidecarContainerPrefix)) || config.Container == container.Name) &&
					isSidecarEnvAvailable(container.Env, sidecarEnvPath, config.File) &&
					isVolumeMountAvailable(container.VolumeMounts, config.Volume) {
					podContainers = append(podContainers, container)
				}
			}
		}
	}
	return podContainers
}

// joinTailinSidecarConfigs joins configurations defined in TailingSidecar resources
func joinTailinSidecarConfigs(tailinSidecars []tailingsidecarv1.TailingSidecar) map[string]tailingsidecarv1.SidecarConfig {
	sidecarConfigs := make(map[string]tailingsidecarv1.SidecarConfig, 0)
	for _, tailitailinSidecar := range tailinSidecars {
		for name, config := range tailitailinSidecar.Spec.Configs {
			sidecarConfigs[name] = config
		}
	}
	return sidecarConfigs
}

// isSidecarAvailable checks if tailing sidecar container with given configuration exists in Pod specification
func isSidecarAvailable(containers []corev1.Container, config tailingsidecarv1.SidecarConfig) bool {
	for _, container := range containers {
		if ((config.Container == "" && strings.HasPrefix(container.Name, sidecarContainerPrefix)) || config.Container == container.Name) &&
			isSidecarEnvAvailable(container.Env, sidecarEnvPath, config.File) &&
			isSidecarEnvAvailable(container.Env, sidecarEnvMarker, sidecarEnvMarkerVal) &&
			isVolumeMountAvailable(container.VolumeMounts, config.Volume) {
			return true
		}
	}
	return false
}

// isSidecarEnvAvailable checks if env is defined and has specific value
func isSidecarEnvAvailable(envs []corev1.EnvVar, envName string, envValue string) bool {
	for _, env := range envs {
		if env.Name == envName && env.Value == envValue {
			return true
		}
	}
	return false
}

// isVolumeMountAvailable checks is volume with given name is available as volume mounted to the container
func isVolumeMountAvailable(volumeMounts []corev1.VolumeMount, volumeName string) bool {
	for _, volumeMount := range volumeMounts {
		if volumeMount.Name == volumeName {
			return true
		}
	}
	return false
}

// getVolume returns volume with given name
func getVolume(containers []corev1.Container, volumeName string) (corev1.VolumeMount, error) {
	for _, container := range containers {
		for _, volume := range container.VolumeMounts {
			if volume.Name == volumeName {
				return volume, nil
			}
		}
	}
	return corev1.VolumeMount{}, fmt.Errorf("volume was not found, volume: %s", volumeName)
}

// getTailingSidecars returns tailing sidecar containers,
// tailing sidecar containers have environmental variable TAILING_SIDECAR=true
func getTailingSidecars(containers []corev1.Container) []corev1.Container {
	tailingSidecars := make([]corev1.Container, 0)
	for _, container := range containers {
		if isSidecarEnvAvailable(container.Env, sidecarEnvMarker, sidecarEnvMarkerVal) {
			tailingSidecars = append(tailingSidecars, container)
		}
	}
	return tailingSidecars
}

// getHostPath returns path to host path directory for Fluent Bit database
func getHostPath(pod *corev1.Pod) string {
	if pod.ObjectMeta.Namespace != "" && pod.ObjectMeta.Name != "" {
		return fmt.Sprintf(hostPathDirPath, pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
	}
	return fmt.Sprintf(hostPathDirPath, strings.TrimRight(pod.GenerateName, "-"), guuid.New().String())
}
