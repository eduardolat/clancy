package runner

// AgentRunner defines the interface for executing agent commands.
// This allows mocking the execution logic for testing.
type AgentRunner interface {
	Run(command string, env map[string]string) (output string, err error)
}

// RealRunner implements AgentRunner using actual system processes.
type RealRunner struct{}

// NewRealRunner creates a new instance of RealRunner.
func NewRealRunner() *RealRunner {
	return &RealRunner{}
}
