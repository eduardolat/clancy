# Clancy

**Clancy** is a robust "Ralph Wiggum" loop orchestrator written in Go. It automates the execution of AI coding agents (like `claude code`) by running them in a persistent loop until a specific success criteria is met or safety limits are reached.

## Features

- **Automated Looping:** Keeps the agent running until it says the "Safe Word" (e.g., `<promise>DONE</promise>`).
- **Safety Limits:** Hard limits on maximum iterations and global timeout.
- **Input Resolution:** Supports reading prompts directly from configuration or external files.
- **Cross-Platform:** Works on Linux, macOS (via PTY), and Windows (via standard pipes).
- **Zero Config Start:** Generate default configuration easily with `--new`.

## Installation

### Quick Install

```bash
curl -sfL https://raw.githubusercontent.com/eduardolat/clancy/main/install.sh | sh
```

### Download Binaries

Download pre-built binaries from the [Releases](https://github.com/eduardolat/clancy/releases) page for Linux, macOS, or Windows.

### Go Install

```bash
go install github.com/eduardolat/clancy/cmd/clancy@latest
```

### Build from Source

```bash
git clone https://github.com/eduardolat/clancy.git
cd clancy
go build -o clancy ./cmd/clancy
```

## Usage

### Quick Start

Generate a new configuration file:
```bash
clancy --new
```
This creates `clancy.yaml`. If the file exists, it creates `clancy-{unique-id}.yaml`.

### Running

```bash
# Run with default clancy.yaml
./clancy

# Run with custom config
./clancy my-task.yaml
```

### Configuration (`clancy.yaml`)

```yaml
version: 1

agent:
  # The command to run. ${PROMPT} is replaced with the content from input.prompt.
  # Windows users: Use double quotes for arguments if needed.
  command: "opencode run '${PROMPT}'"
  env:
    # Optional environment variables
    NO_COLOR: "true"

loop:
  max_steps: 10          # Stop after 10 iterations
  timeout: "30m"         # Stop after 30 minutes
  stop_phrase: "<promise>DONE</promise>" # The success signal
  stop_mode: "exact"     # "exact" or "contains"

input:
  # Can be a string literal or "file:path/to/prompt.md"
  prompt: "file:./tasks/refactor.md"
```

> **Tip:** For reliable stopping, explicitly instruct your LLM in the prompt to output the safe word only when finished. For example:
> "Output `<promise>DONE</promise>` when complete without any explanation."

## How it Works

1. **Clancy** reads the config and the prompt.
2. It constructs the agent command, injecting the prompt.
3. It runs the command in a loop.
4. After each iteration, it checks the agent's output for the `stop_phrase`.
5. If found, it exits successfully. If not, it repeats until `max_steps` or `timeout`.

## Author & Support

Follow me on Twitter ([@eduardoolat](https://x.com/eduardoolat)) for more open source tools and updates!

## License

MIT License. See [LICENSE](LICENSE) for details.
