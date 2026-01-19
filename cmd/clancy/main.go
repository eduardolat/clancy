package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/eduardolat/clancy/internal/config"
	"github.com/eduardolat/clancy/internal/loop"
	"github.com/eduardolat/clancy/internal/runner"
	"github.com/eduardolat/clancy/internal/version"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

//go:embed template.yaml
var templateContent []byte

// Args defines command line arguments.
type Args struct {
	Config string `arg:"positional" default:"clancy.yaml" help:"Path to configuration file"`
	New    bool   `arg:"--new" help:"Generate a new configuration file"`
}

func (Args) Version() string {
	return fmt.Sprintf("Version %s (commit: %s, date: %s)",
		version.Version, version.Commit, version.Date)
}

func (Args) Description() string {
	return "Clancy - AI Agent Loop Orchestrator"
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

	if err := os.WriteFile(filename, templateContent, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "New configuration file generated: %s\n", filename)
	return nil
}
