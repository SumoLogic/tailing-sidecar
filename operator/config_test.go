package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
				},
			},
			expectedError: nil,
		},
		{
			name: "overwrite defaults",
			content: `
sidecar:
  image: my-new-image`,
			expected: Config{
				Sidecar: SidecarConfig{
					Image: "my-new-image",
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
