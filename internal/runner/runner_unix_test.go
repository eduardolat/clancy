//go:build !windows

package runner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrepareCommand(t *testing.T) {
	tests := []struct {
		name     string
		tmpl     string
		prompt   string
		expected string
	}{
		{
			name:     "Simple substitution",
			tmpl:     "echo '${PROMPT}'",
			prompt:   "hello world",
			expected: "echo 'hello world'",
		},
		{
			name:     "Substitution with single quotes",
			tmpl:     "echo '${PROMPT}'",
			prompt:   "I'm a developer",
			expected: "echo 'I'\"'\"'m a developer'",
		},
		{
			name:     "No substitution needed",
			tmpl:     "ls -la",
			prompt:   "whatever",
			expected: "ls -la",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PrepareCommand(tt.tmpl, tt.prompt)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestRealRunner_Run_Echo(t *testing.T) {
	// This test uses the real OS execution, assuming 'echo' exists.
	r := NewRealRunner()
	output, err := r.Run("echo 'hello from runner'", nil)
	require.NoError(t, err)
	require.Contains(t, output, "hello from runner")
}

func TestRealRunner_Run_Env(t *testing.T) {
	r := NewRealRunner()
	env := map[string]string{"TEST_VAR": "custom_value"}
	// We use 'env' command to print environment variables
	output, err := r.Run("env", env)
	require.NoError(t, err)
	require.Contains(t, output, "TEST_VAR=custom_value")
}
