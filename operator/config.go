package main

import (
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/yaml"
)

type Config struct {
	Sidecar SidecarConfig `yaml:"sidecar,omitempty"`
}

type SidecarConfig struct {
	Image     string                      `yaml:"image,omitempty"`
	Resources corev1.ResourceRequirements `yaml:"resources,omitempty"`
}

func ReadConfig(configPath string, defaultConfig Config) (Config, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return defaultConfig, err
	}
	err = yaml.Unmarshal(content, &defaultConfig)
	if err != nil {
		return defaultConfig, err
	}
	return defaultConfig, err
}

func (c *Config) Validate() error {
	return nil
}

func GetDefaultConfig() Config {
	return Config{
		Sidecar: SidecarConfig{
			Image: "sumologic/tailing-sidecar:latest",
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("500Mi"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("100m"),
					corev1.ResourceMemory: resource.MustParse("200Mi"),
				},
			},
		},
	}
}
