package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sidecar SidecarConfig `yaml:"sidecar,omitempty"`
}

type SidecarConfig struct {
	Image string `yaml:"image,omitempty"`
}

func ReadConfig(configPath string) (Config, error) {
	// Set default values
	config := Config{
		Sidecar: SidecarConfig{
			Image: "sumologic/tailing-sidecar:latest",
		},
	}

	content, err := os.Readfile(configPath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}
	return config, err
}

func (c *Config) Validate() error {
	return nil
}
