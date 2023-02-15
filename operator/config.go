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

	f, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)

	if err != nil && err.Error() == "EOF" {
		return config, nil
	}
	return config, err
}

func (c *Config) Validate() error {
	return nil
}
