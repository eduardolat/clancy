package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/eduardolat/clancy/internal/config"
	"github.com/eduardolat/clancy/internal/loop"
	"github.com/eduardolat/clancy/internal/runner"
)

// Args defines command line arguments.
type Args struct {
	Config string `arg:"--config,env:CLANCY_CONFIG" default:"clancy.yaml" help:"Path to configuration file"`
}

func main() {
	var args Args
	arg.MustParse(&args)

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
