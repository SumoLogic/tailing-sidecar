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
	"fmt"
	"strings"

	tailingsidecarv1 "github.com/SumoLogic/tailing-sidecar/operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	configRaw              = 2 // e.g. tailing-sidecar: <volume-name0>:<path-to-tail0>
	configRawWithContainer = 3 // e.g. tailing-sidecar: <container-name0>:<volume-name0>:<path-to-tail0>

	volumeIndex        = 0
	fileIndex          = 1
	containerNameIndex = 0

	volumeFileSeparator = ":"
	configSeparator     = ";"

	sidecarAnnotation = "tailing-sidecar"
)

type sidecarConfig struct {
	name string
	spec tailingsidecarv1.SidecarSpec
}

// getConfigs gets configurations from TailingSidecars and annotations
func getConfigs(annotations map[string]string, tailingSidecarConfigs []tailingsidecarv1.TailingSidecarConfig) ([]sidecarConfig, error) {
	crConfigs, err := convertTailingSidecarConfigs(tailingSidecarConfigs)
	if err != nil {
		return nil, err
	}

	configs := parseAnnotation(annotations)
	configs = append(configs, crConfigs...)

	if err = validateConfigs(configs); err != nil {
		return nil, err
	}
	return configs, nil
}

// parseAnnotation parses configurations from 'tailing-sidecar' annotation
func parseAnnotation(annotations map[string]string) []sidecarConfig {
	annotation, ok := annotations[sidecarAnnotation]
	if !ok {
		return nil
	}

	if annotation == "" {
		handlerLog.Info("Empty tailing-sidecar annotation",
			"annotation", annotation)
		return nil
	}

	configs := make([]sidecarConfig, 0)
	configElements := strings.Split(annotation, configSeparator)

	for _, configElement := range configElements {
		configParts := strings.Split(configElement, volumeFileSeparator)

		nonEmptyConfigParts := removeEmptyConfigs(configParts)

		switch len(nonEmptyConfigParts) {
		case configRaw:
			config := sidecarConfig{
				spec: tailingsidecarv1.SidecarSpec{
					Path: configParts[fileIndex],
					VolumeMount: corev1.VolumeMount{
						Name: configParts[volumeIndex],
					},
				},
			}
			configs = append(configs, config)
		case configRawWithContainer:
			config := sidecarConfig{
				name: configParts[containerNameIndex],
				spec: tailingsidecarv1.SidecarSpec{
					VolumeMount: corev1.VolumeMount{
						Name: configParts[containerNameIndex+1],
					},
					Path: configParts[containerNameIndex+2],
				},
			}
			configs = append(configs, config)
		default:
			handlerLog.Info("Incorrect format of 'tailing-sidecar' annotation",
				"annotation", annotation)
		}
	}
	return configs
}

// convertTailingSidecarConfigs converts configurations defined in TailingSidecarConfigs to sidecarConfig
func convertTailingSidecarConfigs(tailingSidecars []tailingsidecarv1.TailingSidecarConfig) ([]sidecarConfig, error) {
	sidecarNames := make(map[string]struct{}, len(tailingSidecars))
	configs := make([]sidecarConfig, len(tailingSidecars))

	for _, tailitailinSidecar := range tailingSidecars {
		for name, spec := range tailitailinSidecar.Spec.SidecarSpecs {
			if _, ok := sidecarNames[name]; ok {
				return nil, fmt.Errorf("not unique names for tailing sidecar containers in TailingSidecarConfigs, name: %s", name)
			}
			sidecarNames[name] = struct{}{}

			config := sidecarConfig{
				name: name,
				spec: spec,
			}
			configs = append(configs, config)
		}
	}
	return configs, nil
}

// removeEmptyConfigs removes empty elements from configuration e.g. when there is ":" in annotation
func removeEmptyConfigs(configParts []string) []string {
	nonEmptyConfigs := make([]string, 0)
	for _, configPart := range configParts {
		if configPart != "" {
			nonEmptyConfigs = append(nonEmptyConfigs, configPart)
		}
	}
	return nonEmptyConfigs
}

// validateConfigs validates configurations
// checks if container names provided in configurations have unique names
func validateConfigs(configs []sidecarConfig) error {
	containerNames := make(map[string]interface{})
	namesCount := 0
	for _, config := range configs {
		if config.name != "" {
			containerNames[config.name] = nil
			namesCount++
		}
	}

	if len(containerNames) != namesCount {
		return fmt.Errorf("names for tailing sidecar containers must be unique, configs: %+v", configs)
	}
	return nil
}
