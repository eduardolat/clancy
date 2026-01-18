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

// AgentRunner defines the interface for executing agent commands.
// This allows mocking the execution logic for testing.
type AgentRunner interface {
	Run(command string, env map[string]string) (output string, err error)
}

// RealRunner implements AgentRunner using actual system processes and PTY.
type RealRunner struct{}

// NewRealRunner creates a new instance of RealRunner.
func NewRealRunner() *RealRunner {
	return &RealRunner{}
}

// Run executes a shell command in a pseudo-terminal.
// It streams output to os.Stdout and also returns the full captured output.
func (r *RealRunner) Run(command string, env map[string]string) (string, error) {
	// Create the command. We use "sh -c" to allow complex command strings.
	// Using "sh" ensures compatibility with most Unix-like environments.
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
	// We use os.Stdout directly to ensure user sees real-time progress
	mw := io.MultiWriter(os.Stdout, &buf)

	// Copy content. This blocks until the command finishes and closes stdout.
	_, _ = io.Copy(mw, ptmx)

	// Wait for the command to exit
	err = cmd.Wait()
	if err != nil {
		// If it's just an exit code error, we still want the output
		// But strictly speaking, if the agent fails (exit 1), we might return error.
		// However, Ralph loops often rely on output even if exit code is non-zero?
		// Usually agents exit 0 on success.
		// For now, we return the error so the loop knows something went wrong technically.
		return buf.String(), err
	}

	return buf.String(), nil
}

// PrepareCommand is a helper to inject the prompt into the command template.
// It handles escaping of single quotes in the prompt to prevent breaking the shell command if wrapped in quotes.
// NOTE: This is a naive implementation. For robust shell usage, consider how the user configures quotes.
// The spec example is: command: "claude code --print-output --prompt '${PROMPT}'"
func PrepareCommand(tmpl string, prompt string) string {
	// Escape single quotes in the prompt if the user wrapped ${PROMPT} in single quotes.
	// Common shell pattern: '... '${PROMPT}' ...'
	// If prompt contains ', it breaks. We replace ' with '"'"' (close quote, literal quote, open quote)
	// This is standard bash escaping for single quoted strings.

	// However, we don't know for sure if the user used single quotes in the config.
	// But assuming the spec example:
	escapedPrompt := strings.ReplaceAll(prompt, "'", `'"'"'`)
	return strings.ReplaceAll(tmpl, "${PROMPT}", escapedPrompt)
}
