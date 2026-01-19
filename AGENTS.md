# Agent Context for Clancy

## Summary

**Clancy** is a robust loop orchestrator written in Go that automates the execution of AI coding agents (like `claude code` or other AI tools). It runs agents in a persistent loop until a specific success criteria (a "safe word" or stop phrase) is met or safety limits are reached. The tool provides automated looping, safety limits (max steps, global timeout), input resolution from files, and cross-platform support (Linux, macOS via PTY, Windows via standard pipes).

## Maintaining this Document

After completing any task, review this file and update it if you made structural changes or discovered patterns worth documenting. Only add information that helps understand how to work with the project. Avoid implementation details, file listings, or trivial changes. This is a general guide, not a changelog.

When updating this document, do so with the context of the entire document in mind; do not simply add new sections at the end, but place them where they make the most sense within the context of the document.

## Tech Stack

- **Language**: Go 1.25
- **Key Libraries**:
  - `github.com/alexflint/go-arg` - Command line argument parsing
  - `github.com/creack/pty` - PTY (pseudo-terminal) support for Unix systems
  - `github.com/matoous/go-nanoid/v2` - ID generation for config files
  - `gopkg.in/yaml.v3` - YAML configuration parsing
  - `github.com/stretchr/testify` - Testing framework and mocking
- **Build Tool**: Task (Taskfile.yml)
- **CI/CD**: GitHub Actions with devcontainers
- **Linter**: golangci-lint

## Operational Commands

### Build & Run
- **Build**: `task build` (Builds to `./dist/clancy`)
- **Run**: `task run` (Builds and executes)
- **Manual Build**: `go build -o ./dist/clancy ./cmd/clancy/.`

### Testing & Quality
- **Test**: `task test` or `go test -v ./...` - Runs all tests with verbose output
- **Lint**: `task lint` or `golangci-lint run` - Runs golangci-lint
- **Format**: `task fmt` or `go fmt ./...` - Formats Go code

### CI/CD
- **Full CI**: `task ci` - Runs format, lint, test, and build in sequence
- **CI in GitHub Actions**: Uses devcontainers with `task ci` command
- **Cross-Platform Build**: `task build-all` - Builds binaries for Linux (amd64/arm64), macOS (amd64/arm64), and Windows (amd64/arm64)
- **Clean Artifacts**: `task clean` - Removes all build artifacts from dist/

## General Instructions

- **Context**: Always review the project structure via `tree -L 2` before making changes. This is a clean architecture Go project with clear separation between entry point (`cmd/`) and business logic (`internal/`).
- **Tools**: Always use Task (`task <taskname>`) for build, test, lint, and format operations. Do not use direct `make` or shell commands.
- **Verification**: After any changes, run `task ci` to ensure format, lint, tests, and build all pass before considering the task complete.
- **Cross-Platform**: Code must work on Linux, macOS, and Windows. Platform-specific implementations are in `internal/runner/` (e.g., `runner_unix.go`, `runner_windows.go`).
- **Interface-Based Design**: Core components use interfaces (e.g., `AgentRunner`) to enable mocking and testing. Follow this pattern for new components.
- **Configuration**: All behavior is driven by YAML configuration files. The template is embedded in `cmd/clancy/template.yaml`.
- **Error Handling**: Use idiomatic Go error handling with wrapped errors (`fmt.Errorf` with `%w`). Return descriptive errors that help users understand what went wrong.
- **Idiomatic Go**:
  - Follow standard Go project layout
  - Use descriptive, self-documenting names
  - Prefer explicit returns over named returns
  - Use `go generate` comments for embedded files (see `//go:embed` in main.go)
  - Keep functions focused and small
- **Testing**: Write tests for all new functionality using `testing` and `testify` for assertions and mocks. Test files should be in the same package as the code they test (e.g., `loop_test.go` alongside `loop.go`).

## Architecture & Organization

### Root Layout
- `go.mod`: Go module definition and dependencies
- `go.sum`: Go dependency checksums
- `Taskfile.yml`: Task runner configuration (build, test, lint, format, cross-platform builds)
- `install.sh`: Installation script for downloading and installing pre-built binaries
- `clancy`: Compiled binary (in git for distribution purposes)
- `clancy*.yaml`: Configuration files (examples and generated configs)
- `README.md`: User documentation and usage guide
- `LICENSE`: MIT license
- `.github/workflows/ci.yaml`: CI/CD pipeline configuration using devcontainers
- `.github/workflows/release.yaml`: Release workflow for publishing releases

### `cmd/clancy/` (Entry Point)
- **Role**: Application entry point and CLI interface
- **Key Files**:
  - `main.go`: Main function, argument parsing, and orchestration
  - `template.yaml`: Embedded configuration template (used with `//go:embed`)
- **Responsibilities**:
  - Parse command-line arguments using `go-arg`
  - Handle `--new` flag to generate config files
  - Load configuration from YAML
  - Initialize components (runner, loop)
  - Coordinate the main execution flow

### `internal/` (Business Logic)
- **Role**: Private internal packages containing core business logic
- **Key Directories**:
  - `config/`: Configuration loading, parsing, and validation
    - `config.go`: Configuration structs, YAML loading, prompt resolution
    - `config_test.go`: Configuration parsing tests
  - `loop/`: Loop orchestration and stop condition checking
    - `loop.go`: Main loop logic with timeout and max steps handling
    - `loop_test.go`: Comprehensive loop behavior tests with mocks
  - `runner/`: Command execution and platform-specific implementations
    - `runner.go`: `AgentRunner` interface and real runner implementation
    - `runner_unix.go`: Unix-specific execution with PTY support
    - `runner_windows.go`: Windows-specific execution via standard pipes
    - `runner_unix_test.go`: Unix runner tests
  - `version/`: Version information populated at build time via ldflags
    - `version.go`: Version, commit hash, and build date variables

### `dist/` (Build Artifacts)
- **Role**: Compiled output directory
- **Key Files**:
  - `clancy`: Built binary (target of `task build`)
  - `clancy-{os}-{arch}`: Cross-platform binaries from `task build-all`
  - `checksums.txt`: SHA256 checksums for all binaries

### `spec/` (Specifications)
- **Role**: Project specifications and requirements (currently empty)
- **Purpose**: Intended for detailed specifications of features and behaviors

## Testing Strategy

- **Framework**: Standard Go `testing` package with `github.com/stretchr/testify` for assertions and mocks
- **Conventions**:
  - Test files are named `<package>_test.go` and placed alongside the code they test
  - Tests use table-driven testing patterns where applicable
  - Mock interfaces using `testify/mock` for unit testing
  - Test functions are descriptive: `Test<Component>_<Scenario>_<ExpectedResult>`
  - Use `require.NoError` and `require.Error` for assertions to fail fast
  - Mock expectations are asserted with `mock.AssertExpectations(t)`
- **Coverage Areas**:
  - Configuration parsing and validation (different stop modes, defaults)
  - Loop logic (success on first step, multi-step, max steps, timeout)
  - Stop condition checking (exact vs contains mode)
  - Platform-specific runner behavior
- **Command**: `task test` runs all tests with verbose output

## Design Patterns & Conventions

### Clean Architecture
- The project follows clean architecture principles with clear separation:
  - **Entry Point** (`cmd/clancy`): CLI and orchestration only
  - **Business Logic** (`internal/`): Core functionality without dependencies
  - **Interfaces**: Define contracts for external dependencies (e.g., `AgentRunner`)

### Configuration-Driven Behavior
- All execution parameters are read from YAML configuration
- Use `//go:embed` for static resources (configuration templates)
- Support both literal values and external file references (via `file:` prefix)

### Error Handling
- Wrap errors with context using `fmt.Errorf` with `%w`
- Provide clear, actionable error messages to users
- Distinguish between recoverable errors (agent failure) and fatal errors (config parsing)

### Cross-Platform Implementation
- Use build tags for platform-specific code: `//go:build unix` or `//go:build windows`
- Share common interface, implement platform-specific behaviors
- Test on target platforms when possible

## Key Behaviors

### Loop Execution
1. Load configuration from YAML file
2. Resolve input prompt (from literal or file)
3. Prepare command by substituting `${PROMPT}` placeholder
4. Execute loop:
   - Check timeout before each iteration
   - Run agent command
   - Capture output (even on non-zero exit)
   - Check for stop phrase in output
   - Repeat until stop phrase found, max steps reached, or timeout expires
5. Exit with success (stop phrase found) or error (limits exceeded)

### Stop Phrase Modes
- **exact**: Output must exactly match the stop phrase.
- **contains**: Output must contain the stop phrase anywhere.
- **suffix**: Output must **end** with the stop phrase. (Recommended for agents that "think" before acting).

> **Note:** Clancy applies `strings.TrimSpace()` to the agent's output **before** any check in all modes.

### Safety Limits
- `max_steps`: Maximum number of loop iterations (default: 10)
- `timeout`: Global timeout for entire loop execution (default: 30m)
- Both limits are enforced and the loop stops gracefully when reached

### Configuration Generation
- Running with `--new` flag generates a new config file
- If `clancy.yaml` exists, generates `clancy-{nanoid}.yaml` instead
- Uses embedded template from `cmd/clancy/template.yaml`
- NanoID uses lowercase alphanumeric alphabet, 6 characters
