package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	content := `
agent:
  command: "echo 'hello'"
  env:
    FOO: "bar"
loop:
  max_steps: 5
  timeout: "1h"
  stop_phrase: "DONE"
  stop_mode: "exact"
input:
  prompt: "Do work"
`
	tmpfile := filepath.Join(t.TempDir(), "clancy.yaml")
	require.NoError(t, os.WriteFile(tmpfile, []byte(content), 0644))

	cfg, err := Load(tmpfile)
	require.NoError(t, err)

	require.Equal(t, "echo 'hello'", cfg.Agent.Command)
	require.Equal(t, "bar", cfg.Agent.Env["FOO"])
	require.Equal(t, 5, cfg.Loop.MaxSteps)
	require.Equal(t, "DONE", cfg.Loop.StopPhrase)
	require.Equal(t, "exact", cfg.Loop.StopMode)
	require.Equal(t, time.Hour, cfg.Loop.TimeoutDuration)
	require.Equal(t, "Do work", cfg.Input.Prompt)
}

func TestDefaultStopMode(t *testing.T) {
	content := `
agent:
  command: "echo"
loop:
  max_steps: 1
input:
  prompt: "foo"
`
	tmpfile := filepath.Join(t.TempDir(), "clancy_defaults.yaml")
	require.NoError(t, os.WriteFile(tmpfile, []byte(content), 0644))

	cfg, err := Load(tmpfile)
	require.NoError(t, err)
	require.Equal(t, "exact", cfg.Loop.StopMode)
}

func TestDelayParsing(t *testing.T) {
	content := `
agent:
  command: "echo"
loop:
  delay: "500ms"
input:
  prompt: "foo"
`
	tmpfile := filepath.Join(t.TempDir(), "clancy_delay.yaml")
	require.NoError(t, os.WriteFile(tmpfile, []byte(content), 0644))

	cfg, err := Load(tmpfile)
	require.NoError(t, err)
	require.Equal(t, 500*time.Millisecond, cfg.Loop.DelayDuration)
}

func TestResolvePrompt_String(t *testing.T) {
	cfg := &Config{
		Input: InputConfig{Prompt: "Just a string"},
	}
	p, err := cfg.ResolvePrompt()
	require.NoError(t, err)
	require.Equal(t, "Just a string", p)
}

func TestResolvePrompt_File(t *testing.T) {
	promptContent := "This is the task"
	tmpDir := t.TempDir()
	promptFile := filepath.Join(tmpDir, "task.txt")
	require.NoError(t, os.WriteFile(promptFile, []byte(promptContent), 0644))

	cfg := &Config{
		Input: InputConfig{Prompt: "file:" + promptFile},
	}
	p, err := cfg.ResolvePrompt()
	require.NoError(t, err)
	require.Equal(t, promptContent, p)
}

func TestResolvePrompt_FileMissing(t *testing.T) {
	cfg := &Config{
		Input: InputConfig{Prompt: "file:./nonexistent.txt"},
	}
	_, err := cfg.ResolvePrompt()
	require.Error(t, err)
}
