package loop

import (
	"testing"
	"time"

	"github.com/eduardolat/clancy/internal/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRunner is a mock implementation of AgentRunner
type MockRunner struct {
	mock.Mock
}

func (m *MockRunner) Run(command string, env map[string]string) (string, error) {
	args := m.Called(command, env)
	return args.String(0), args.Error(1)
}

func TestRun_Success(t *testing.T) {
	// Scenario: Loop runs once, outputs "RALPH_DONE" immediately.
	cfg := &config.Config{
		Agent: config.AgentConfig{Command: "echo '${PROMPT}'"},
		Loop: config.LoopConfig{
			MaxSteps:        5,
			StopPhrase:      "RALPH_DONE",
			StopMode:        "contains",
			Timeout:         "1m",
			TimeoutDuration: time.Minute,
		},
	}
	prompt := "do work"

	mockRunner := new(MockRunner)
	// Expectation: Run called once.
	// Note: Command will have prompt injected. "echo 'do work'"
	mockRunner.On("Run", "echo 'do work'", cfg.Agent.Env).Return("Work complete. RALPH_DONE", nil).Times(1)

	err := Run(cfg, mockRunner, prompt)
	require.NoError(t, err)
	mockRunner.AssertExpectations(t)
}

func TestRun_MultiStep_Success(t *testing.T) {
	// Scenario: Loop runs twice. First time "working...", second time "RALPH_DONE".
	cfg := &config.Config{
		Agent: config.AgentConfig{Command: "cmd"},
		Loop: config.LoopConfig{
			MaxSteps:        5,
			StopPhrase:      "DONE",
			StopMode:        "exact",
			TimeoutDuration: time.Minute,
		},
	}

	mockRunner := new(MockRunner)
	// Call 1
	mockRunner.On("Run", "cmd", cfg.Agent.Env).Return("working...", nil).Once()
	// Call 2
	mockRunner.On("Run", "cmd", cfg.Agent.Env).Return("DONE", nil).Once()

	err := Run(cfg, mockRunner, "p")
	require.NoError(t, err)
	mockRunner.AssertExpectations(t)
}

func TestRun_MaxSteps_Reached(t *testing.T) {
	// Scenario: Never returns stop phrase.
	cfg := &config.Config{
		Agent: config.AgentConfig{Command: "cmd"},
		Loop: config.LoopConfig{
			MaxSteps:        3,
			StopPhrase:      "DONE",
			StopMode:        "contains",
			TimeoutDuration: time.Minute,
		},
	}

	mockRunner := new(MockRunner)
	mockRunner.On("Run", "cmd", cfg.Agent.Env).Return("still working", nil).Times(3)

	err := Run(cfg, mockRunner, "p")
	require.Error(t, err)
	require.Contains(t, err.Error(), "max steps (3) reached")
	mockRunner.AssertExpectations(t)
}

func TestRun_Timeout(t *testing.T) {
	// Scenario: Timeout triggers.
	// We set a very short timeout.
	cfg := &config.Config{
		Agent: config.AgentConfig{Command: "cmd"},
		Loop: config.LoopConfig{
			MaxSteps:        10,
			StopPhrase:      "DONE",
			TimeoutDuration: 1 * time.Millisecond,
		},
	}

	mockRunner := new(MockRunner)
	// The loop checks timeout BEFORE running.
	// But getting the mock to be slow is hard without a sleep inside the mock.
	// However, if we set timeout to 1ms, it likely expires before the first run or during it.
	// Let's make the Mock sleep slightly to force timeout.

	mockRunner.On("Run", "cmd", cfg.Agent.Env).Run(func(args mock.Arguments) {
		time.Sleep(10 * time.Millisecond)
	}).Return("working", nil)

	err := Run(cfg, mockRunner, "p")
	require.Error(t, err)
	// It could be "global timeout reached" or the loop just finished 1 step and then timed out on next check.
	// If the sleep happens inside Run, Run returns, then loop checks ctx.Done().

	// Wait, if Run takes long, the ctx might be done when Run returns.
	// Our loop check is at the top of the loop.
	// So:
	// 1. Loop start.
	// 2. Check ctx. Not done.
	// 3. Run() (Sleeps 10ms). Timeout (1ms) expires during sleep.
	// 4. Run returns.
	// 5. Check Stop (False).
	// 6. Loop increments.
	// 7. Check ctx. Done! -> Return Error.

	require.Contains(t, err.Error(), "global timeout reached")
}

func TestCheckStopCondition(t *testing.T) {
	require.True(t, CheckStopCondition(" foo DONE bar ", "DONE", "contains"))
	require.False(t, CheckStopCondition(" foo done bar ", "DONE", "contains")) // Case sensitive usually

	require.True(t, CheckStopCondition("DONE", "DONE", "exact"))
	require.True(t, CheckStopCondition("  DONE  \n", "DONE", "exact")) // Trimmed
	require.False(t, CheckStopCondition("DONE.", "DONE", "exact"))

	require.True(t, CheckStopCondition("Working on task... DONE", "DONE", "suffix"))
	require.True(t, CheckStopCondition("  DONE  \n", "DONE", "suffix")) // Trimmed, also exact match
	require.False(t, CheckStopCondition("DONE in the middle", "DONE", "suffix"))
	require.False(t, CheckStopCondition("DONEX", "DONE", "suffix"))

	require.True(t, CheckStopCondition("  thinking... DONE  \n", "DONE", "suffix")) // Trimmed + suffix
}

func TestRun_WithDelay(t *testing.T) {
	// Scenario: Loop runs twice with a small delay.
	cfg := &config.Config{
		Agent: config.AgentConfig{Command: "cmd"},
		Loop: config.LoopConfig{
			MaxSteps:        2,
			StopPhrase:      "DONE",
			StopMode:        "exact",
			TimeoutDuration: time.Minute,
			Delay:           "100ms",
			DelayDuration:   100 * time.Millisecond,
		},
	}

	mockRunner := new(MockRunner)
	// Call 1
	mockRunner.On("Run", "cmd", cfg.Agent.Env).Return("working...", nil).Once()
	// Call 2
	mockRunner.On("Run", "cmd", cfg.Agent.Env).Return("DONE", nil).Once()

	start := time.Now()
	err := Run(cfg, mockRunner, "p")
	duration := time.Since(start)

	require.NoError(t, err)
	mockRunner.AssertExpectations(t)

	// Should have waited at least 100ms
	require.GreaterOrEqual(t, duration, 100*time.Millisecond)
}
