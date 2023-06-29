package main

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/yaml"
)

type Config struct {
	Sidecar        SidecarConfig        `yaml:"sidecar,omitempty"`
	LeaderElection LeaderElectionConfig `yaml:"leaderElection,omitempty"`
}

type SidecarConfig struct {
	Image     string                      `yaml:"image,omitempty"`
	Resources corev1.ResourceRequirements `yaml:"resources,omitempty"`
	Config    SidecarConfigConfig         `yaml:"config,omitempty"`
}

type LeaderElectionConfig struct {
	LeaseDuration Duration `yaml:"leaseDuration,omitempty"`
	RenewDeadline Duration `yaml:"renewDeadline,omitempty"`
	RetryPeriod   Duration `yaml:"retryPeriod,omitempty"`
}

type SidecarConfigConfig struct {
	ConfigMapName string `yaml:"name,omitempty"`
	MountPath     string `yaml:"mountPath,omitempty"`
	NamespaceName string `yaml:"namespace,omitempty"`
}

// Duration sigs.k8s.io/yaml not support time.Duration:https://github.com/kubernetes-sigs/yaml/issues/64
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

func ReadConfig(configPath string, config *Config) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		return err
	}
	return err
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
		// reference for values: https://github.com/open-telemetry/opentelemetry-operator/blob/a8653601cd6a6e2b35fd7f3e1a28b4e9608fb794/main.go#L181
		LeaderElection: LeaderElectionConfig{
			LeaseDuration: Duration(time.Second * 137),
			RenewDeadline: Duration(time.Second * 107),
			RetryPeriod:   Duration(time.Second * 26),
		},
	}
}
