package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/eduardolat/clancy/internal/config"
	"github.com/eduardolat/clancy/internal/loop"
	"github.com/eduardolat/clancy/internal/runner"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// Args defines command line arguments.
type Args struct {
	Config string `arg:"positional" default:"clancy.yaml" help:"Path to configuration file"`
	New    bool   `arg:"--new" help:"Generate a new configuration file"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	// Handle --new flag
	if args.New {
		if err := generateConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating config: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 1. Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config file '%s': %v\n", args.Config, err)
		os.Exit(1)
	}

	// 2. Resolve Input Prompt
	prompt, err := cfg.ResolvePrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving prompt: %v\n", err)
		os.Exit(1)
	}

	// 3. Initialize Runner
	r := runner.NewRealRunner()

	// 4. Run Loop
	fmt.Fprintf(os.Stderr, ">>> [Clancy] Starting loop. Config: %s, Steps: %d, Timeout: %s\n",
		args.Config, cfg.Loop.MaxSteps, cfg.Loop.Timeout)

	if err := loop.Run(cfg, r, prompt); err != nil {
		fmt.Fprintf(os.Stderr, ">>> [Clancy] Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, ">>> [Clancy] Success.\n")
}

func generateConfig() error {
	filename := "clancy.yaml"

	// Check if default exists
	if _, err := os.Stat(filename); err == nil {
		// Generate short NanoID using custom alphabet and length 6
		const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
		id, err := gonanoid.Generate(alphabet, 6)
		if err != nil {
			return fmt.Errorf("failed to generate ID: %w", err)
		}
		filename = fmt.Sprintf("clancy-%s.yaml", id)
	}

	content := `version: 1

agent:
  # The command to run. ${PROMPT} is replaced with the content from input.prompt.
  # Note: Ensure usage of quotes compatible with your shell.
  command: "opencode run '${PROMPT}'"
  env:
    # Optional environment variables
    NO_COLOR: "true"

loop:
  max_steps: 10          # Stop after 10 iterations
  timeout: "30m"         # Stop after 30 minutes
  stop_phrase: "DONE"    # The success signal
  stop_mode: "exact"     # "exact" or "contains"

input:
  # Can be a string literal or "file:path/to/prompt.md"
  prompt: "file:./tasks/task.md"
`
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Generated configuration file: %s\n", filename)
	return nil
}
