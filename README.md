# Snippety

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/tahcohcat/snippety?style=flat-square)](https://github.com/tahcohcat/snippety/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/tahcohcat/snippety/build.yml?style=flat-square)](https://github.com/tahcohcat/snippety/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/tahcohcat/snippety?style=flat-square)](https://goreportcard.com/report/github.com/tahcohcat/snippety)

AI-powered Go CLI tool that generates commit messages from staged git changes using Ollama.

## Features

- 🤖 **AI-Generated Messages**: Uses Ollama to generate meaningful commit messages
- 📝 **Conventional Commits**: Follows conventional commit format (Add, Fix, Update, Remove)
- 🎭 **Customizable Tone**: Choose from professional, fun, pirate, or serious tones
- 🤝 **Interactive Mode**: Optionally confirm before creating and pushing commits with generated messages
- 📁 **Auto-staging**: Automatically stages all changes with `git add -A` before analysis (can be disabled)
- 🔄 **Fallback Support**: Falls back to rule-based generation if Ollama is unavailable
- ⚙️ **Configurable**: Supports custom Ollama endpoints and models
- 🚀 **Fast & Lightweight**: Built with Go and Cobra CLI framework

## Prerequisites

- [Ollama](https://ollama.ai/) installed and running
- Go 1.21 or later
- Git repository with staged changes

## Installation

```bash
# Clone the repository
git clone https://github.com/tahcohcat/snippety/snippety
cd snippety

# Build the binary
go mod tidy
go build -o snippety ./cmd/snippety

# Optional: Install globally
go install ./cmd/snippety
```

## Quick Start

1. **Start Ollama**:
```bash
ollama serve
```

2. **Pull a model** (if not already available):
```bash
ollama pull llama3.2
```

3. **Generate commit message** (auto-stages changes by default):
```bash
./snippety

# Or with interactive confirmation:
./snippety --interactive
```

## Usage

### Basic Usage
```bash
./snippety
```

### Custom Ollama Configuration
```bash
# Custom Ollama server URL
./snippety --ollama-url http://localhost:11434

# Different model
./snippety --model llama3.1

# Custom tone
./snippety --tone fun

# Interactive mode
./snippety --interactive

# Disable auto-staging (manual git add required)
./snippety --auto-stage=false

# Combined options
./snippety --ollama-url http://remote-server:11434 --model codellama --tone pirate --interactive
```

### Tone Options
```bash
# Professional tone (default)
./snippety --tone professional

# Fun tone with emojis and creative language
./snippety --tone fun

# Pirate speak with nautical terminology
./snippety --tone pirate

# Haiku poem with 5-7-5 syllable structure (single line)
./snippety --tone haiku

# Serious, formal tone with technical precision
./snippety --tone serious
```

### Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `--ollama-url` | `http://localhost:11434` | Ollama server URL |
| `--model` | `llama3.2` | Ollama model to use for generation |
| `--tone` | `professional` | Tone for commit messages (professional, fun, pirate, haiku, serious) |
| `--interactive` | `false` | Interactively confirm before creating and pushing the git commit |
| `--auto-stage` | `true` | Automatically stage all changes with 'git add -A' before analysis |

## Example Output

```bash
$ ./snippety
Staging all changes...
Generating commit message with Ollama...
Making request to: http://localhost:11434/api/generate
Generated commit message:
Add user authentication middleware

$ ./snippety --tone fun
Generating commit message with Ollama...
Generated commit message:
✨ Add shiny new user auth middleware 🔐

$ ./snippety --tone pirate
Generating commit message with Ollama...
Generated commit message:
Hoist new authentication middleware aboard! ⚓

$ ./snippety --tone haiku
Generating commit message with Ollama...
Generated commit message:
Auth middleware flows / Through the codebase like fresh streams / Security blooms bright

$ ./snippety --interactive
Staging all changes...
Generating commit message with Ollama...
Generated commit message:
Add user authentication middleware

Do you want to create a commit with this message? (y/N): y
✅ Commit created successfully!
🚀 Commit pushed successfully!
```

## How It Works

1. **Auto-staging**: Automatically runs `git add -A` to stage all changes (unless disabled)
2. **Git Diff Analysis**: Retrieves staged changes using `git diff --staged`
3. **AI Processing**: Sends the diff to Ollama with a specialized prompt
4. **Commit Generation**: Returns a conventional commit message
5. **Interactive Confirmation**: Optionally prompts user to create commit and push
6. **Fallback**: Uses rule-based analysis if Ollama is unavailable

## Supported Models

Any Ollama model should work, but these are recommended:

- `llama3.2` (default) - Good balance of speed and quality
- `llama3.1` - Larger model for better accuracy
- `codellama` - Specialized for code analysis
- `mistral` - Fast and efficient
- `qwen2.5-coder` - Excellent for code understanding

## Troubleshooting

### Ollama Connection Issues

```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags

# Start Ollama if not running
ollama serve

# Verify model availability
ollama list
```

### Common Errors

- **"connection refused"**: Ollama service is not running
- **"405 Method Not Allowed"**: Wrong URL or port
- **"No staged changes found"**: Run `git add` to stage your changes first

## Development

### Project Structure
```
snippety/
├── cmd/main.go              # Entry point
├── internal/
│   ├── cmd/                 # Cobra CLI commands
│   │   ├── root.go         # Root command setup
│   │   └── generate.go     # Core generation logic
│   └── ollama/             # Ollama client
│       └── client.go       # HTTP client for Ollama API
├── go.mod                  # Go module definition
└── README.md              # This file
```

### Building from Source

```bash
# Get dependencies
go mod tidy

# Run tests (if available)
go test ./...

# Build binary
go build -o snippety ./cmd

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o snippety-linux ./cmd
GOOS=windows GOARCH=amd64 go build -o snippety.exe ./cmd
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Ollama](https://ollama.ai/) for the local LLM runtime
- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- [Conventional Commits](https://www.conventionalcommits.org/) for the commit format standard
