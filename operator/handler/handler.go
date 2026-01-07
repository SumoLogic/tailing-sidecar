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
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	tailingsidecarv1 "github.com/SumoLogic/tailing-sidecar/operator/api/v1"
	admv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/add-tailing-sidecars-v1-pod,mutating=true,failurePolicy=ignore,groups="",resources=pods,verbs=create;update;delete,versions=v1,name=tailing-sidecar.sumologic.com,sideEffects=none,admissionReviewVersions={v1,v1beta1}

const (
	sidecarEnvPath                   = "PATH_TO_TAIL"
	sidecarOtelFileStoragePathEnv    = "OTEL_FILE_STORAGE_PATH"
	sidecarOtelFileStoragePath       = "/var/lib/otc/tailing-sidecar-%d"
	sidecarOtelFileStorageVolumeName = "tailing-sidecar-otel-file-storage-tailing-sidecar-%d"
	sidecarOtelLogsPathEnv           = "SIDECAR_OTEL_LOG_PATH"
	sidecarOtelLogsPath              = "/var/log/tailing-sidecar-%d"
	sidecarOtelLogsVolumeName        = "tailing-sidecar-otel-logs-tailing-sidecar-%d"
	sidecarEnvMarker                 = "TAILING_SIDECAR"
	sidecarEnvMarkerVal              = "true"

	sidecarContainerName    = "tailing-sidecar-%d"
	sidecarContainerNameEnv = "SIDECAR_CONTAINER_NAME"
	sidecarContainerPrefix  = "tailing-sidecar-"

	sidecarVolumeName   = "volume-sidecar-%d"
	sidecarVolumePrefix = "volume-sidecar-"
	sidecarMountPath    = "/tailing-sidecar/var"

	sidecarConfigurationName = "tailing-sidecar-configuration"

	deletionMessage = "Tailing Sidecar Operator does not block Pod deletion"
)

var handlerLog = ctrl.Log.WithName("tailing-sidecar.operator.handler.PodExtender")

// PodExtender extends Pods by tailling sidecar containers
type PodExtender struct {
	Client                  client.Client
	TailingSidecarImage     string
	TailingSidecarResources corev1.ResourceRequirements
	Decoder                 admission.Decoder
	ConfigMapName           string
	ConfigMapNamespace      string
	ConfigMountPath         string
}

// Handle handles requests to create/update Pod and extends it by adding tailing sidecars
func (e *PodExtender) Handle(ctx context.Context, req admission.Request) admission.Response {
	if req.Operation == "" {
		return admission.Allowed("Received startupProbe/livenessProbe")
	}

	if req.Operation == admv1.Delete {
		return e.handleDelete(ctx, req)
	}

	pod := &corev1.Pod{}
	err := e.Decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	tailingSidecarConfigs, err := e.getTailingSidecarConfigs(ctx, pod.ObjectMeta.Labels)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	if _, ok := pod.ObjectMeta.Annotations[sidecarAnnotation]; !ok && len(tailingSidecarConfigs) == 0 {
		return admission.Allowed("Configuration for Tailing Sidecar Operator is not provided")
	}

	handlerLog.Info("Handling request",
		"Name", req.Name,
		"Namespace", req.Namespace,
		"Kind", req.Kind,
		"GenerateName", pod.ObjectMeta.GenerateName,
		"Operation", req.Operation,
	)

	if err := e.extendPod(ctx, pod, tailingSidecarConfigs, req); err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	if err := validateContainers(pod.Spec.Containers); err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// extendPod extends Pod by adding tailing sidecars according to configuration in annotation
func (e PodExtender) extendPod(ctx context.Context, pod *corev1.Pod, tailingSidecarConfigs []tailingsidecarv1.TailingSidecarConfig, req admission.Request) error {
	// Get number of existing tailing sidecars
	sidecarsCount := len(getTailingSidecars(pod.Spec.Containers))

	// Get configurations from TailingSidecars and annotations
	configs, err := getConfigs(pod.ObjectMeta.Annotations, tailingSidecarConfigs)
	if err != nil {
		handlerLog.Error(err, "Incorrect configuration")
		return err
	}

	namespace := req.Namespace

	if len(configs) == 0 && sidecarsCount == 0 {
		handlerLog.Info("Pod does not need to be configured",
			"Name", req.Name,
			"Namespace", namespace,
			"Kind", req.Kind,
			"GenerateName", pod.ObjectMeta.GenerateName,
		)
		return nil
	}

	handlerLog.Info("Found configuration for Pod",
		"Name", req.Name,
		"Namespace", namespace,
		"Kind", req.Kind,
		"GenerateName", pod.ObjectMeta.GenerateName,
	)

	containers := make([]corev1.Container, 0)
	for _, config := range configs {

		err := prepareVolume(pod.Spec.Containers, &config.spec.VolumeMount)
		if err != nil {
			handlerLog.Error(err,
				"Failed to prepare volume",
				"Name", req.Name,
				"Namespace", namespace,
				"Kind", req.Kind,
				"GenerateName", pod.ObjectMeta.GenerateName,
			)
			continue
		}

		if isSidecarAvailable(pod.Spec.Containers, config) {
			// Do not add tailing sidecar if tailing sidecar with specific configuration exists
			handlerLog.Info("Tailing sidecar exists",
				"config", config,
			)
			continue
		}

		volumeName := fmt.Sprintf(sidecarVolumeName, sidecarsCount)
		if config.name == "" {
			config.name = fmt.Sprintf(sidecarContainerName, sidecarsCount)
		}

		pod.Spec.Volumes = append(pod.Spec.Volumes,
			corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{}},
			})

		otelLogsVolumneName := fmt.Sprintf(sidecarOtelLogsVolumeName, sidecarsCount)
		otelCollectorLogsPath := fmt.Sprintf(sidecarOtelLogsPath, sidecarsCount)

		otelFileStorageVolumeName := fmt.Sprintf(sidecarOtelFileStorageVolumeName, sidecarsCount)
		otelFileStoragePath := fmt.Sprintf(sidecarOtelFileStoragePath, sidecarsCount)

		pod.Spec.Volumes = append(pod.Spec.Volumes,
			corev1.Volume{
				Name: otelLogsVolumneName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{}},
			})

		pod.Spec.Volumes = append(pod.Spec.Volumes,
			corev1.Volume{
				Name: otelFileStorageVolumeName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{}},
			})

		// check if sidecar need add default resources
		if config.spec.Resources.Requests == nil {
			config.spec.Resources.Requests = e.TailingSidecarResources.Requests
		}
		if config.spec.Resources.Limits == nil {
			config.spec.Resources.Limits = e.TailingSidecarResources.Limits
		}

		volumeMounts := []corev1.VolumeMount{
			config.spec.VolumeMount,
			{
				Name:      volumeName,
				MountPath: sidecarMountPath,
			},
			{
				Name:      otelFileStorageVolumeName,
				MountPath: otelFileStoragePath,
			},

			{
				Name:      otelLogsVolumneName,
				MountPath: otelCollectorLogsPath,
			},
		}

		if e.ConfigMapName != "" && e.ConfigMountPath != "" && e.ConfigMapNamespace != "" {
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      sidecarConfigurationName,
				MountPath: e.ConfigMountPath,
			})
		}

		container := corev1.Container{
			Image: e.TailingSidecarImage,
			Name:  config.name,
			Env: []corev1.EnvVar{
				{
					Name:  sidecarEnvPath,
					Value: config.spec.Path,
				},
				{
					Name:  sidecarEnvMarker,
					Value: sidecarEnvMarkerVal,
				},
				{
					Name:  sidecarOtelFileStoragePathEnv,
					Value: otelFileStoragePath,
				},

				{
					Name:  sidecarOtelLogsPathEnv,
					Value: otelCollectorLogsPath,
				},
				{
					Name:  sidecarContainerNameEnv,
					Value: config.name,
				},
			},
			VolumeMounts: volumeMounts,
			Resources:    config.spec.Resources,
		}
		containers = append(containers, container)
		pod.ObjectMeta.Annotations = addAnnotations(pod.ObjectMeta.Annotations, config)
		sidecarsCount++
	}
	podContainers := removeDeletedSidecars(pod.Spec.Containers, configs)

	pod.Spec.Containers = append(podContainers, containers...)

	if e.ConfigMapName != "" && e.ConfigMountPath != "" && e.ConfigMapNamespace != "" {
		err = e.createSidecarConfigMap(ctx, namespace)
		if err != nil {
			return err
		}

		pod.Spec.Volumes = append(pod.Spec.Volumes,
			corev1.Volume{
				Name: sidecarConfigurationName,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: e.ConfigMapName,
						},
					},
				},
			})
	}

	pod.Spec.Volumes = filterUnusedVolumes(pod.Spec.Volumes, pod.Spec.Containers)
	return nil
}

func (e PodExtender) getTailingSidecarConfigs(ctx context.Context, podLabels map[string]string) ([]tailingsidecarv1.TailingSidecarConfig, error) {
	tailingSidecarConfigList := &tailingsidecarv1.TailingSidecarConfigList{}
	tailingSidecarConfigListOpts := []client.ListOption{}

	if err := e.Client.List(ctx, tailingSidecarConfigList, tailingSidecarConfigListOpts...); err != nil {
		handlerLog.Error(err, "Failed to get list of TailingSidecarConfigs")
		return nil, err
	}

	tailingSidcarConfigs := make([]tailingsidecarv1.TailingSidecarConfig, 0)
	for _, tailingSidcarConfig := range tailingSidecarConfigList.Items {
		selector, err := metav1.LabelSelectorAsSelector(tailingSidcarConfig.Spec.PodSelector)
		if err != nil {
			return nil, fmt.Errorf("invalid label selector in TailingSidecarConfig: %v", err)
		}
		// TailingSidecarConfig with a nil or empty selector should match nothing
		if selector.Empty() || !selector.Matches(labels.Set(podLabels)) {
			continue
		}
		tailingSidcarConfigs = append(tailingSidcarConfigs, tailingSidcarConfig)
	}
	return tailingSidcarConfigs, nil
}

func (e PodExtender) createSidecarConfigMap(ctx context.Context, namespace string) error {
	// get exemplar configMap
	mainConfigMap := corev1.ConfigMap{}
	key := types.NamespacedName{
		Namespace: e.ConfigMapNamespace,
		Name:      e.ConfigMapName,
	}

	err := e.Client.Get(ctx, key, &mainConfigMap)
	if err != nil {
		return errors.Join(fmt.Errorf("cannot find exemplar (%s) configuration in `%s` namespace", e.ConfigMapName, e.ConfigMapNamespace), err)
	}

	// create configMap for non-exemplar namespace
	if namespace != e.ConfigMapNamespace {
		key := types.NamespacedName{
			Namespace: namespace,
			Name:      e.ConfigMapName,
		}
		configMap := corev1.ConfigMap{}
		err = e.Client.Get(ctx, key, &configMap)
		if err == nil {
			configMap.Data = mainConfigMap.Data
			err = e.Client.Update(ctx, &configMap)
			if err != nil {
				return errors.Join(fmt.Errorf("error updating tailing sidecar configuration in `%s` namespace", namespace), err)
			}
		} else {
			// prepare configMap
			configMap.SetNamespace(namespace)
			configMap.Data = mainConfigMap.Data
			configMap.SetAnnotations(mainConfigMap.GetAnnotations())
			configMap.SetLabels(mainConfigMap.GetLabels())
			configMap.SetName(mainConfigMap.GetName())

			err = e.Client.Create(ctx, &configMap)
			if err != nil {
				return errors.Join(fmt.Errorf("cannot create tailing sidecar configuration in `%s` namespace", namespace), err)
			}
		}
	}

	return nil
}

// validateContainers validates containers
// checks if there is container names conflict
// potential conflict e.g. when container not managed by operator has name with prefix "tailing-sidecar"
func validateContainers(containers []corev1.Container) error {
	containerNames := make(map[string]interface{})
	for _, container := range containers {
		containerNames[container.Name] = nil
	}
	if len(containerNames) != len(containers) {
		return fmt.Errorf("container names not unique, when name is not configured for tailing sidecar container it starts with 'tailing-sidecar' prefix")
	}
	return nil
}

// removeDeletedSidecars removes deleted tailing sidecar containers from Pod specification
func removeDeletedSidecars(containers []corev1.Container, configs []sidecarConfig) []corev1.Container {
	podContainers := make([]corev1.Container, 0)
	for _, container := range containers {
		if !isSidecarEnvAvailable(container.Env, sidecarEnvMarker, sidecarEnvMarkerVal) {
			podContainers = append(podContainers, container)
		} else {
			for _, config := range configs {
				if ((config.name == "" && strings.HasPrefix(container.Name, sidecarContainerPrefix)) || config.name == container.Name) &&
					isSidecarEnvAvailable(container.Env, sidecarEnvPath, config.spec.Path) &&
					isVolumeMountAvailable(container.VolumeMounts, config.spec.VolumeMount) {
					podContainers = append(podContainers, container)
				}
			}
		}
	}
	return podContainers
}

// filterUnusedVolumes filters out unused volumes, previously assigned to tailing sidecars from the provided slice.
// Each of tailing-sidecars has its own volume to store Fluent Bit database.
// When sidecar container is removed volume is no longer needed.
func filterUnusedVolumes(volumes []corev1.Volume, containers []corev1.Container) []corev1.Volume {
	podVolumes := make([]corev1.Volume, 0)
	for _, volume := range volumes {
		if !strings.HasPrefix(volume.Name, sidecarVolumePrefix) {
			// name of volumes assigned to tailing sidecar starts with 'volume-sidecar' prefix
			// when volumes starts with different prefix it should not be filtered out
			podVolumes = append(podVolumes, volume)
			continue
		}
		found := false
		for _, container := range containers {
			for _, volumeMount := range container.VolumeMounts {
				if volumeMount.Name == volume.Name {
					podVolumes = append(podVolumes, volume)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}
	return podVolumes
}

// isSidecarAvailable checks if tailing sidecar container with given configuration exists in Pod specification
func isSidecarAvailable(containers []corev1.Container, config sidecarConfig) bool {
	for _, container := range containers {
		if ((config.name == "" && strings.HasPrefix(container.Name, sidecarContainerPrefix)) || config.name == container.Name) &&
			isSidecarEnvAvailable(container.Env, sidecarEnvPath, config.spec.Path) &&
			isSidecarEnvAvailable(container.Env, sidecarEnvMarker, sidecarEnvMarkerVal) &&
			isVolumeMountAvailable(container.VolumeMounts, config.spec.VolumeMount) {
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

// isVolumeMountAvailable checks if volume is available as volume mounted to the container
func isVolumeMountAvailable(volumeMounts []corev1.VolumeMount, volume corev1.VolumeMount) bool {
	for _, volumeMount := range volumeMounts {
		if reflect.DeepEqual(volumeMount, volume) {
			return true
		}
	}
	return false
}

// prepareVolume returns volume with given name
func prepareVolume(containers []corev1.Container, sidecarVolume *corev1.VolumeMount) error {
	for _, container := range containers {
		for _, volume := range container.VolumeMounts {
			if volume.Name == sidecarVolume.Name {
				if sidecarVolume.MountPath == "" {
					// mount volume at the same path as for container with log file
					// by default mountPath does not need to be provided in configuration
					sidecarVolume.MountPath = volume.MountPath
				}
				return nil
			}
		}
	}
	return fmt.Errorf("volume provided in configuration is not mounted to any container, volume name: %s", sidecarVolume.Name)
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

// handle pod deletion:
// - remove sidecar configmap if it is not used by any pod in the req namespace
func (e *PodExtender) handleDelete(ctx context.Context, req admission.Request) admission.Response {
	// check if configmap is configured
	if e.ConfigMapName == "" || e.ConfigMapNamespace == "" || e.ConfigMountPath == "" {
		return admission.Allowed(deletionMessage)
	}

	// eliminates hanging kubectl apply -f command
	// kube-apiserver server waits for response from operator on DELETE request
	pod := &corev1.Pod{}
	err := e.Decoder.DecodeRaw(req.OldObject, pod)
	if err != nil {
		return admission.Allowed(fmt.Sprintf("Error ocurred (%v); %s", err, deletionMessage))
	}

	// do not try to remove exemplar configMap
	if req.Namespace == e.ConfigMapNamespace {
		return admission.Allowed(deletionMessage)
	}

	// check if theres is sidecar configuration configMap in the namespace
	volumes := pod.Spec.Volumes
	ret := true
	for _, v := range volumes {
		if v.Name == sidecarConfigurationName {
			ret = false
		}
	}
	if ret {
		return admission.Allowed(deletionMessage)
	}

	// get sidecar configmap from the namespace
	key := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      e.ConfigMapName,
	}
	configMap := corev1.ConfigMap{}
	err = e.Client.Get(ctx, key, &configMap)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return admission.Allowed(deletionMessage)
		}
		return admission.Allowed(fmt.Sprintf("Error while getting configMap to clean up (%v); %s", err, deletionMessage))
	}

	// check if the configMap is used by any pod in the namespace
	podList := &corev1.PodList{}

	listOptions := []client.ListOption{
		client.InNamespace(req.Namespace),
	}
	err = e.Client.List(ctx, podList, listOptions...)
	if err != nil {
		return admission.Allowed(fmt.Sprintf("Error while getting list of pods (%v); %s", err, deletionMessage))
	}

	for _, p := range podList.Items {
		// skip current pod as it is going to be removed anyway
		if p.Name == pod.Name && p.Namespace == pod.Namespace {
			continue
		}

		// check if sidecar configuration is attached to the p pod
		for _, volume := range p.Spec.Volumes {
			if volume.ConfigMap != nil && volume.ConfigMap.Name == e.ConfigMapName {
				// do not anything in case volume is used
				return admission.Allowed(deletionMessage)
			}
		}
	}

	// delete configMap as it is not used by any pod
	err = e.Client.Delete(ctx, &configMap)
	if err != nil {
		return admission.Allowed(fmt.Sprintf("Error while cleaning up configMap (%v); %s", err, deletionMessage))
	}

	return admission.Allowed(deletionMessage)
}
