# Snippety

AI-powered Go CLI tool that generates commit messages from staged git changes using Ollama.

## Features

- 🤖 **AI-Generated Messages**: Uses Ollama to generate meaningful commit messages
- 📝 **Conventional Commits**: Follows conventional commit format (Add, Fix, Update, Remove)
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
git clone https://github.com/snippety/snippety
cd snippety

# Build the binary
go mod tidy
go build -o snippety ./cmd

# Optional: Install globally
go install ./cmd
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

3. **Stage your changes**:
```bash
git add .
```

4. **Generate commit message**:
```bash
./snippety
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

# Both custom URL and model
./snippety --ollama-url http://remote-server:11434 --model codellama
```

### Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `--ollama-url` | `http://localhost:11434` | Ollama server URL |
| `--model` | `llama3.2` | Ollama model to use for generation |

## Example Output

```bash
$ ./snippety
Generating commit message with Ollama...
Making request to: http://localhost:11434/api/generate
Generated commit message:
Add user authentication middleware
```

## How It Works

1. **Git Diff Analysis**: Retrieves staged changes using `git diff --staged`
2. **AI Processing**: Sends the diff to Ollama with a specialized prompt
3. **Commit Generation**: Returns a conventional commit message
4. **Fallback**: Uses rule-based analysis if Ollama is unavailable

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
