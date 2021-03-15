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
)

const (
	configRaw              = 2 // e.g. tailing-sidecar: <volume-name0>:<path-to-tail0>
	configRawWithContainer = 3 // e.g. tailing-sidecar: <container-name0>:<volume-name0>:<path-to-tail0>
	configPredefined       = 1 // e.g. tailing-sidecar: <config-name0>

	volumeIndex     = 0
	fileIndex       = 1
	containerIndex  = 0
	configNameIndex = 0

	volumeFileSeparator = ":"
	configSeparator     = ";"

	sidecarAnnotation = "tailing-sidecar"
)

// getConfigs gets configurations from TailingSidecars and annotations
func getConfigs(annotations map[string]string, tailingSidecars []tailingsidecarv1.TailingSidecar) []tailingsidecarv1.SidecarConfig {
	annotation, ok := annotations[sidecarAnnotation]
	if !ok {
		return []tailingsidecarv1.SidecarConfig{}
	}

	if annotation == "" {
		handlerLog.Info("Empty tailing-sidecar annotation",
			"annotation", annotation)
		return []tailingsidecarv1.SidecarConfig{}
	}

	sidecarConfigs := joinTailingSidecarConfigs(tailingSidecars)

	configs := parseAnnotation(annotation, sidecarConfigs)
	return configs
}

// parseAnnotation parses configurations from 'tailing-sidecar' annotation and joins them with configurations from TailingSidecars
func parseAnnotation(annotation string, sidecarConfigs map[string]tailingsidecarv1.SidecarConfig) []tailingsidecarv1.SidecarConfig {
	configs := make([]tailingsidecarv1.SidecarConfig, 0)
	configElements := strings.Split(annotation, configSeparator)

	for _, configElement := range configElements {
		configParts := strings.Split(configElement, volumeFileSeparator)

		nonEmptyConfigParts := removeEmptyConfigs(configParts)

		switch len(nonEmptyConfigParts) {
		case configRaw:
			config := tailingsidecarv1.SidecarConfig{
				File:   configParts[fileIndex],
				Volume: configParts[volumeIndex],
			}
			configs = append(configs, config)
		case configRawWithContainer:
			config := tailingsidecarv1.SidecarConfig{
				Container: configParts[containerIndex],
				Volume:    configParts[containerIndex+1],
				File:      configParts[containerIndex+2],
			}
			configs = append(configs, config)
		case configPredefined:
			if _, ok := sidecarConfigs[configParts[configNameIndex]]; !ok {
				handlerLog.Info("Missing configuration in TailingSidecar configurations",
					"configuration name", configParts[configNameIndex],
				)
				continue
			}
			configs = append(configs, sidecarConfigs[configParts[configNameIndex]])
		default:
			handlerLog.Info("Incorrect configuration",
				"annotation", annotation)
		}
	}
	return configs
}

// joinTailingSidecarConfigs joins configurations defined in TailingSidecar resources
func joinTailingSidecarConfigs(tailingSidecars []tailingsidecarv1.TailingSidecar) map[string]tailingsidecarv1.SidecarConfig {
	sidecarConfigs := make(map[string]tailingsidecarv1.SidecarConfig, len(tailingSidecars))
	for _, tailitailinSidecar := range tailingSidecars {
		for name, config := range tailitailinSidecar.Spec.Configs {
			sidecarConfigs[name] = config
		}
	}
	return sidecarConfigs
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
func validateConfigs(configs []tailingsidecarv1.SidecarConfig) error {
	containerNames := make(map[string]interface{})
	namesCount := 0
	for _, config := range configs {
		if config.Container != "" {
			containerNames[config.Container] = nil
			namesCount++
		}
	}

	if len(containerNames) != namesCount {
		return fmt.Errorf("names for tailing sidecars must be unique")
	}
	return nil
}
