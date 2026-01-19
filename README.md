<p align="center">
  <h1 align="center">Clancy Wiggum</h1>
  <p align="center">
    <img align="center" src="./gopher.png" height="250px" width="auto">
  </p>
  <p align="center">
    <b>Clancy</b> (named after Chief Clancy Wiggum) is a robust <a href="https://awesomeclaude.ai/ralph-wiggum">"Ralph Wiggum"</a> loop orchestrator written in Go.
    <br/>
    It automates the execution of AI coding agents (like opencode, claude code, etc) by running them in a persistent loop until a specific success criteria is met or safety limits are reached.
  </p>
</p>

<p align="center">
  <a href="https://github.com/eduardolat/clancy/actions">
    <img src="https://github.com/eduardolat/clancy/actions/workflows/ci.yaml/badge.svg" alt="CI status"/>
  </a>
  <a href="https://pkg.go.dev/github.com/eduardolat/clancy">
    <img src="https://pkg.go.dev/badge/github.com/eduardolat/clancy" alt="Go Reference"/>
  </a>
  <a href="https://goreportcard.com/report/eduardolat/clancy">
    <img src="https://goreportcard.com/badge/eduardolat/clancy" alt="Go Report Card"/>
  </a>
  <a href="https://github.com/eduardolat/clancy/releases/latest">
    <img src="https://img.shields.io/github/release/eduardolat/clancy.svg" alt="Release Version"/>
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/github/license/eduardolat/clancy.svg" alt="License"/>
  </a>
  <a href="https://github.com/eduardolat/clancy">
    <img src="https://img.shields.io/github/stars/eduardolat/clancy?style=flat&label=github+stars"/>
  </a>
</p>

## Features

- **Automated Looping:** Keeps the agent running until it says the "Safe Word" (e.g., `<promise>DONE</promise>`).
- **Safety Limits:** Hard limits on maximum iterations and global timeout.
- **Input Resolution:** Supports reading prompts directly from configuration or external files.
- **Cross-Platform:** Works on Linux, macOS (via PTY), and Windows (via standard pipes).
- **Zero Config Start:** Generate default configuration easily with `--new`.

## Installation

### Quick Install Script (recommended)

You can also run this script to update to the latest version.

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
    FOO: "bar"

loop:
  max_steps: 10 # Stop after 10 iterations
  timeout: "30m" # Global timeout. Units: ns, us, ms, s, m, h
  delay: "10s" # Optional wait time between iterations
  stop_phrase: "<promise>DONE</promise>" # The success signal
  stop_mode: "exact" # "exact" or "contains"

input:
  # Can be a string literal or "file:path/to/prompt.md"
  prompt: "file:./task.md"
```

> **Tip:** For reliable stopping, explicitly instruct your LLM in the prompt to output the safe word only when finished. For example:
> "Output `<promise>DONE</promise>` when complete without any explanation."

## Example Prompts

The `prompts/` directory contains ready-to-use prompt templates for common automation tasks. These are generic examples that should be adapted to your project's specific needs. For more specialized use cases, create your own custom prompt.

- **[coverage-hunter.md](prompts/coverage-hunter.md)** - Generate, execute, and verify unit tests to systematically increase test coverage
- **[docs-syncer.md](prompts/docs-syncer.md)** - Synchronize documentation with actual source code across any language
- **[i18n-extractor.md](prompts/i18n-extractor.md)** - Extract hardcoded strings and replace with localization keys for internationalization
- **[inline-documenter.md](prompts/inline-documenter.md)** - Add standard inline documentation (JSDoc, GoDoc, PyDoc) to exported functions and classes
- **[log-cleanup.md](prompts/log-cleanup.md)** - Remove debug print statements or upgrade to structured logging
- **[refactor-auditor.md](prompts/refactor-auditor.md)** - Eliminate technical debt, fix linter errors, or perform code migrations safely
- **[security-auditor.md](prompts/security-auditor.md)** - Detect and fix potential security vulnerabilities like hardcoded secrets

## Security & Best Practices

⚠️ **Important:** Clancy gives AI agents full autonomy to execute commands repeatedly until completion. For maximum safety and to take full advantage of this autonomy:

- **Run in Containers:** Execute Clancy inside Docker containers or similar isolated environments
- **Controlled Environment:** Use dedicated development environments or VMs with proper resource limits
- **Permission Management:** Grant agents the necessary permissions for your specific task (file system access, network, etc.) within the controlled boundary
- **Data Safety:** Ensure no sensitive data or critical production systems are accessible from the agent's execution environment

Running in a controlled, isolated environment allows you to confidently give agents full autonomy while keeping your host system safe.

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
