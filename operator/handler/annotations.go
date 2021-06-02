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
	"regexp"
)

const defaultAnnotationKeyPrefix = "tailing-sidecar.sumologic.com"

var whiteSpaces = regexp.MustCompile(`\s+`)

// addAnnotations adds per tailing sidecar container annotations to Pod specification
// according to configuration defined in TailingSidecarConfig
// format of annotation: <annotation-prefix>/<sidecar-container-name>.<annotation-key>: <annotation-value>
func addAnnotations(annotations map[string]string, config sidecarConfig) map[string]string {
	if annotations == nil {
		annotations = make(map[string]string, 0)
	}

	if config.name == "" {
		handlerLog.Info(
			"Missing tailing sidecar container name, cannot create per sidecar annotations",
			"config", config,
		)
		return annotations
	}

	if config.annotationsPrefix == "" {
		config.annotationsPrefix = defaultAnnotationKeyPrefix
	}

	keyPrefix := fmt.Sprintf("%s/%s.", config.annotationsPrefix, config.name)

	for k, v := range config.spec.Annotations {
		annotations[keyPrefix+k] = v
	}
	return annotations
}
