//go:build !windows

package runner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

// Run executes a shell command in a pseudo-terminal.
// It streams output to os.Stdout and also returns the full captured output.
func (r *RealRunner) Run(command string, env map[string]string) (string, error) {
	// Create the command. We use "sh -c" to allow complex command strings.
	cmd := exec.Command("sh", "-c", command)

	// Build environment
	currentEnv := os.Environ()
	newEnv := make([]string, 0, len(currentEnv)+len(env))
	newEnv = append(newEnv, currentEnv...)
	for k, v := range env {
		newEnv = append(newEnv, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = newEnv

	// Start with PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to start pty: %w", err)
	}
	defer func() { _ = ptmx.Close() }() // Best effort close

	// Capture output while streaming
	var buf bytes.Buffer

	// MultiWriter to write to both Stdout and our buffer
	mw := io.MultiWriter(os.Stdout, &buf)

	// Copy content.
	_, _ = io.Copy(mw, ptmx)

	// Wait for the command to exit
	err = cmd.Wait()
	if err != nil {
		return buf.String(), err
	}

	return buf.String(), nil
}

// PrepareCommand injects the prompt into the command template using Bash escaping.
func PrepareCommand(tmpl string, prompt string) string {
	// Escape single quotes: ' -> '"'"'
	escapedPrompt := strings.ReplaceAll(prompt, "'", `'"'"'`)
	return strings.ReplaceAll(tmpl, "${PROMPT}", escapedPrompt)
}
