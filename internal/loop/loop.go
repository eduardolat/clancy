package loop

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/eduardolat/clancy/internal/config"
	"github.com/eduardolat/clancy/internal/runner"
)

// Run executes the Ralph loop based on the provided configuration.
func Run(cfg *config.Config, r runner.AgentRunner, prompt string) error {
	ctx := context.Background()
	if cfg.Loop.TimeoutDuration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Loop.TimeoutDuration)
		defer cancel()
	}

	cmd := runner.PrepareCommand(cfg.Agent.Command, prompt)

	for i := 1; i <= cfg.Loop.MaxSteps; i++ {
		// Check context before starting
		select {
		case <-ctx.Done():
			return fmt.Errorf("global timeout reached after %s", cfg.Loop.Timeout)
		default:
		}

		fmt.Fprintf(os.Stderr, ">>> [Clancy] Step %d/%d\n", i, cfg.Loop.MaxSteps)

		// Execute
		output, err := r.Run(cmd, cfg.Agent.Env)
		if err != nil {
			// If agent fails, we usually log it.
			// Should we stop? The spec doesn't explicitly say to stop on error,
			// but usually if the tool crashes, the loop might be broken.
			// However, Ralph often fixes broken builds.
			// But if the *command execution* itself fails (e.g. command not found), we should stop.
			// Runner returns error on non-zero exit code too.
			// Let's print the error and continue, assuming the agent is "trying".
			// Unless it's a critical system error.
			// For now, warn and check stop phrase anyway (maybe the error message contains the stop phrase? unlikely).
			fmt.Fprintf(os.Stderr, ">>> [Clancy] Agent exited with error: %v\n", err)
		}

		// Check Stop Phrase
		if CheckStopCondition(output, cfg.Loop.StopPhrase, cfg.Loop.StopMode) {
			fmt.Fprintf(os.Stderr, ">>> [Clancy] Stop phrase '%s' detected. Loop Complete.\n", cfg.Loop.StopPhrase)
			return nil
		}
	}

	return fmt.Errorf("max steps (%d) reached without success", cfg.Loop.MaxSteps)
}

// CheckStopCondition evaluates if the output meets the stop criteria.
func CheckStopCondition(output, phrase, mode string) bool {
	if mode == "exact" {
		return strings.TrimSpace(output) == phrase
	}
	// Default to contains
	return strings.Contains(output, phrase)
}
