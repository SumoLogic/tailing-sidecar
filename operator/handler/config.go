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
	configPredefined       = 1 // e.g. tailing-sidecar: <config-name0>

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
func getConfigs(annotations map[string]string, tailingSidecars []tailingsidecarv1.TailingSidecarConfig) []sidecarConfig {
	annotation, ok := annotations[sidecarAnnotation]
	if !ok {
		return nil
	}

	if annotation == "" {
		handlerLog.Info("Empty tailing-sidecar annotation",
			"annotation", annotation)
		return nil
	}

	sidecarSpecs := joinTailingSidecarSpecs(tailingSidecars)

	configs := combineConfigs(annotation, sidecarSpecs)
	return configs
}

// combineConfigs parses configurations from 'tailing-sidecar' annotation and joins them with configurations from TailingSidecarConfigs
func combineConfigs(annotation string, sidecarSpecs map[string]tailingsidecarv1.SidecarSpec) []sidecarConfig {
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
		case configPredefined:
			if _, ok := sidecarSpecs[configParts[containerNameIndex]]; !ok {
				handlerLog.Info("Missing configuration in TailingSidecarConfig configurations",
					"configuration name", configParts[containerNameIndex],
				)
				continue
			}
			config := sidecarConfig{
				name: configParts[containerNameIndex],
				spec: sidecarSpecs[configParts[containerNameIndex]],
			}
			configs = append(configs, config)
		default:
			handlerLog.Info("Incorrect configuration",
				"annotation", annotation)
		}
	}
	return configs
}

// joinTailingSidecarSpecs joins configurations defined in TailingSidecarConfig resources
func joinTailingSidecarSpecs(tailingSidecars []tailingsidecarv1.TailingSidecarConfig) map[string]tailingsidecarv1.SidecarSpec {
	sidecarSpecs := make(map[string]tailingsidecarv1.SidecarSpec, len(tailingSidecars))
	for _, tailitailinSidecar := range tailingSidecars {
		for name, spec := range tailitailinSidecar.Spec.SidecarSpecs {
			sidecarSpecs[name] = spec
		}
	}
	return sidecarSpecs
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
		return fmt.Errorf("names for tailing sidecar containers must be unique")
	}
	return nil
}
