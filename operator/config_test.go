package main

import (
	"io/ioutil"
	"os"
	"testing"

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
      memory: 20Mi`,
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
			config, err := ReadConfig(file.Name())

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
	_, err := ReadConfig("non-existing-file")
	require.Error(t, err)
	require.EqualError(t, err, "open non-existing-file: no such file or directory")
}
