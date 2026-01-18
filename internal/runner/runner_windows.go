//go:build windows

package runner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Run executes a command using cmd.exe on Windows.
// Note: PTY support is limited/absent here, so we use standard pipes.
func (r *RealRunner) Run(command string, env map[string]string) (string, error) {
	// Use cmd /C to execute the command string
	cmd := exec.Command("cmd", "/C", command)

	// Build environment
	currentEnv := os.Environ()
	newEnv := make([]string, 0, len(currentEnv)+len(env))
	newEnv = append(newEnv, currentEnv...)
	for k, v := range env {
		newEnv = append(newEnv, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = newEnv

	// Capture output
	var buf bytes.Buffer

	// Stream to stdout/stderr
	cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return buf.String(), err
	}

	return buf.String(), nil
}

// PrepareCommand injects the prompt into the command template using Windows escaping logic.
// It assumes the user wraps ${PROMPT} in double quotes in the config for Windows.
func PrepareCommand(tmpl string, prompt string) string {
	// Escape double quotes: " -> \"
	// This is a common convention for CLI args on Windows.
	escapedPrompt := strings.ReplaceAll(prompt, "\"", `\"`)
	return strings.ReplaceAll(tmpl, "${PROMPT}", escapedPrompt)
}
