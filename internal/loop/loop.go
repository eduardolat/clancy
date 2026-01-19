package loop

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eduardolat/clancy/internal/config"
	"github.com/eduardolat/clancy/internal/runner"
)

// ANSI Color Codes
const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
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
		// 1. HEADER (Cyan Box)
		if i > 1 {
			_, _ = fmt.Fprint(os.Stdout, "\n\n") // Visual separation from previous step
		}

		// Update Window Title (Passive Monitoring)
		_, _ = fmt.Fprintf(os.Stdout, "\033]0;🍩 Clancy: Step %d/%d\007", i, cfg.Loop.MaxSteps)

		printHeaderBox(i, cfg.Loop.MaxSteps)

		// Check Context before execution
		select {
		case <-ctx.Done():
			return fmt.Errorf("global timeout reached")
		default:
		}

		// 2. EXECUTION (With breathing room)
		_, _ = fmt.Fprintln(os.Stdout) // Blank line BEFORE agent output
		output, err := r.Run(cmd, cfg.Agent.Env)
		_, _ = fmt.Fprintln(os.Stdout) // Blank line AFTER agent output

		if err != nil {
			// CRITICAL ERROR (Red Box)
			printErrorBox(err)
		}

		// 3. CHECK CONDITION
		if CheckStopCondition(output, cfg.Loop.StopPhrase, cfg.Loop.StopMode) {
			// SUCCESS (Green Box)
			printSuccessBox(i)
			// Update Window Title to Done
			_, _ = fmt.Fprint(os.Stdout, "\033]0;✅ Clancy: Done\007")
			return nil
		}

		// 4. RETRY & DELAY
		if i < cfg.Loop.MaxSteps {
			// RETRY (Yellow Box)
			printRetryBox(i)

			if cfg.Loop.DelayDuration > 0 {
				// STATIC VISIBLE DELAY LOG (No animation)
				// Bold Yellow text to make it pop
				_, _ = fmt.Fprintf(os.Stdout, "\n%s%s⏳ COOLDOWN: Waiting %s before next step...%s\n",
					colorBold, colorYellow, cfg.Loop.Delay, colorReset)

				// Sleep with context check
				select {
				case <-ctx.Done():
					return fmt.Errorf("global timeout reached during delay")
				case <-time.After(cfg.Loop.DelayDuration):
					// Continue to next iteration
				}
			}
		}
	}

	return fmt.Errorf("max steps (%d) reached without success", cfg.Loop.MaxSteps)
}

// CheckStopCondition evaluates if the output meets the stop criteria.
func CheckStopCondition(output, phrase, mode string) bool {
	if mode == "exact" {
		return strings.TrimSpace(output) == phrase
	}
	return strings.Contains(output, phrase)
}

// --- Visual Helpers (Box System) ---

func printHeaderBox(step, total int) {
	c := colorCyan
	r := colorReset
	// Heavy box style for high visibility
	_, _ = fmt.Fprintf(os.Stdout, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", c, r)
	_, _ = fmt.Fprintf(os.Stdout, "%s  🍩 CLANCY LOOP | STEP %02d/%02d%s\n", c, step, total, r)
	_, _ = fmt.Fprintf(os.Stdout, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", c, r)
}

func printSuccessBox(step int) {
	g := colorGreen
	r := colorReset
	_, _ = fmt.Fprintf(os.Stdout, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", g, r)
	_, _ = fmt.Fprintf(os.Stdout, "%s  ✅ SUCCESS! Stop phrase found in step %02d%s\n", g, step, r)
	_, _ = fmt.Fprintf(os.Stdout, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", g, r)
}

func printRetryBox(step int) {
	y := colorYellow
	r := colorReset
	_, _ = fmt.Fprintf(os.Stdout, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", y, r)
	_, _ = fmt.Fprintf(os.Stdout, "%s  🔄 Stop phrase NOT found in step %02d. Retrying...%s\n", y, step, r)
	_, _ = fmt.Fprintf(os.Stdout, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", y, r)
}

func printErrorBox(err error) {
	red := colorRed
	r := colorReset
	// Format error string to fit roughly in the box
	// Truncate to 55 chars to allow for "..."
	errStr := fmt.Sprintf("%.55s...", err.Error())

	_, _ = fmt.Fprintf(os.Stderr, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", red, r)
	_, _ = fmt.Fprintf(os.Stderr, "%s  💥 CRITICAL: Agent execution failed!%s\n", red, r)
	_, _ = fmt.Fprintf(os.Stderr, "%s  %v%s\n", red, errStr, r)
	_, _ = fmt.Fprintf(os.Stderr, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", red, r)
}
