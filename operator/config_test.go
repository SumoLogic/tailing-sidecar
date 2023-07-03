package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestReadConfig(t *testing.T) {
	testCases := []struct {
		name          string
		content       string
		expected      Config
		expectedError error
	}{
		{
			name:    "empty file",
			content: ``,
			expected: Config{
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
				LeaderElection: LeaderElectionConfig{
					LeaseDuration: Duration(time.Second * 137),
					RenewDeadline: Duration(time.Second * 107),
					RetryPeriod:   Duration(time.Second * 26),
				},
			},
			expectedError: nil,
		},
		{
			name: "defaults",
			content: `
sidecar:
  commands:
    - test
    - command`,
			expected: Config{
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
				LeaderElection: LeaderElectionConfig{
					LeaseDuration: Duration(time.Second * 137),
					RenewDeadline: Duration(time.Second * 107),
					RetryPeriod:   Duration(time.Second * 26),
				},
			},
			expectedError: nil,
		},
		{
			name: "overwrite defaults",
			content: `
sidecar:
  image: my-new-image
  resources:
    limits:
      cpu: 400m
      memory: 400Mi
    requests:
      cpu: 20m
      memory: 20Mi
leaderElection:
  leaseDuration: 10s
  renewDeadline: 10s
  retryPeriod: 10s`,
			expected: Config{
				Sidecar: SidecarConfig{
					Image: "my-new-image",
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("400m"),
							corev1.ResourceMemory: resource.MustParse("400Mi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("20m"),
							corev1.ResourceMemory: resource.MustParse("20Mi"),
						},
					},
				},
				LeaderElection: LeaderElectionConfig{
					LeaseDuration: Duration(time.Second * 10),
					RenewDeadline: Duration(time.Second * 10),
					RetryPeriod:   Duration(time.Second * 10),
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			file, err := ioutil.TempFile(".", "prefix")
			require.NoError(t, err)
			defer os.Remove(file.Name())

			_, err = file.WriteString(tt.content)
			require.NoError(t, err)
			config := GetDefaultConfig()
			err = ReadConfig(file.Name(), &config)
			require.NoError(t, err)

			if tt.expectedError != nil {
				require.Error(t, tt.expectedError, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, config)
		})
	}
}

func TestReadConfigInvalidFile(t *testing.T) {
	config := GetDefaultConfig()
	err := ReadConfig("non-existing-file", &config)
	require.Error(t, err)
	require.EqualError(t, err, "open non-existing-file: no such file or directory")
}
