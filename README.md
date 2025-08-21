# Snippety

AI-powered Go CLI tool that generates commit messages from staged git changes using Ollama.

## Features

- 🤖 **AI-Generated Messages**: Uses Ollama to generate meaningful commit messages with both title and detailed description
- 🎫 **Smart Branch Detection**: Automatically detects ticket prefixes from branch names (e.g., `BP-1234-feature` → `BP-1234: commit title`)
- 📝 **Conventional Commits**: Follows conventional commit format (Add, Fix, Update, Remove)
- 🎭 **Customizable Tone**: Choose from professional, fun, pirate, or serious tones
- 🤝 **Interactive Mode**: Optionally confirm before creating and pushing commits with generated messages
- 📁 **Auto-staging**: Automatically stages all changes with `git add -A` before analysis (can be disabled)
- 🔄 **Fallback Support**: Falls back to rule-based generation if Ollama is unavailable
- 🌊 **Smart Push Handling**: Automatically sets upstream for new branches during push
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

### Basic Usage
```bash
$ ./snippety
Staging all changes...
Generating commit message with Ollama...
Generated commit message:
Title: Add user authentication middleware
Description: Implemented JWT-based authentication middleware for API routes. Added validation for bearer tokens and user session management. Includes unit tests for authentication edge cases.
```

### With Branch Ticket Detection
```bash
# On branch: FEAT-1234-auth-middleware
$ ./snippety
Staging all changes...
Generating commit message with Ollama...
Generated commit message:
Title: FEAT-1234: Add user authentication middleware
Description: Implemented JWT-based authentication middleware for API routes. Added validation for bearer tokens and user session management. Includes unit tests for authentication edge cases.
```

### Different Tones
```bash
$ ./snippety --tone fun
Generated commit message:
Title: ✨ Add shiny new user auth middleware 🔐
Description: Whipped up some awesome JWT magic for our API routes! Now we've got bearer token validation and user sessions that actually work. Added tests because we're responsible developers! 🎉

$ ./snippety --tone pirate
Generated commit message:
Title: Hoist new authentication middleware aboard! ⚓
Description: Arrr! We've plundered the finest JWT treasures and secured our API routes from scurvy hackers. Added proper token validation and session management, plus tests to keep the crew honest, matey!
```

### Interactive Mode
```bash
$ ./snippety --interactive
Staging all changes...
Generated commit message:
Title: FEAT-1234: Add user authentication middleware
Description: Implemented JWT-based authentication middleware for API routes. Added validation for bearer tokens and user session management. Includes unit tests for authentication edge cases.

Do you want to create a commit with this message? (y/N): y
✅ Commit created successfully!
🚀 Commit pushed successfully!
```

### Branch Without Upstream
```bash
$ ./snippety --interactive
# ... commit creation ...
✅ Commit created successfully!
🚀 Commit pushed successfully!
# Automatically runs: git push --set-upstream origin feat/NEW-456-feature
```

## How It Works

1. **Branch Detection**: Detects current branch and extracts ticket prefixes (e.g., `FEAT-1234-feature` → `FEAT-1234:`)
2. **Auto-staging**: Automatically runs `git add -A` to stage all changes (unless disabled)
3. **Git Diff Analysis**: Retrieves staged changes using `git diff --staged`
4. **AI Processing**: Sends the diff to Ollama with a specialized prompt for title and description generation
5. **Commit Generation**: Returns a structured commit message with title and detailed description
6. **Prefix Integration**: Automatically prepends ticket prefix to commit title
7. **Interactive Confirmation**: Optionally prompts user to create commit and push
8. **Smart Push**: Automatically sets upstream for new branches during push
9. **Fallback**: Uses rule-based analysis if Ollama is unavailable

## Commit Message Format

Snippety generates commits using the dual `-m` flag format:

```bash
git commit -m "TICKET-123: Short descriptive title" -m "Detailed description with technical implementation details, changes made, and any test cases added."
```

This creates commits with both a concise title (≤50 characters) and a comprehensive description for better commit history.

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
- **"has no upstream branch"**: Automatically handled - will set upstream during push

### Branch Naming Patterns

Snippety automatically detects and extracts ticket prefixes from these branch patterns:

- `FEAT-1234-description` → `FEAT-1234: commit title`
- `BP-5678-fix-bug` → `BP-5678: commit title` 
- `feature/DEVOPS-999` → `DEVOPS-999: commit title`
- `chore/PROJ-123-cleanup` → `PROJ-123: commit title`

For branches like `main`, `master`, or unrecognized patterns, no prefix is added.

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
