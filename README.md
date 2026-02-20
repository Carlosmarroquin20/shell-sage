# 🧙‍♂️ Shell Sage (ssage)

**Shell Sage** is an AI-powered CLI assistant that brings the power of local LLMs (via Ollama) directly to your terminal. It helps you understand commands, fix errors, and analyze logs with a professional, hacker-style aesthetic.

![Shell Sage](https://placehold.co/600x400?text=Shell+Sage+Demo)

## 🚀 Features

- **`ssage explain "<command>"`**: Get a concise, flag-by-flag explanation of any shell command.
- **`ssage fix`**: The "Crown Jewel". Scans your recent shell history, detects failed commands, and suggests fixes using AI.
- **`ssage analyze <file>`**: Reads error logs and provides a summary of critical issues.

## 🛠️ Prerequisites

Before using Shell Sage, ensure you have the following installed:

1.  **Go 1.21+**: [Download and Install Go](https://go.dev/dl/)
2.  **Ollama**: [Download Ollama](https://ollama.com/)
    - Running locally on `localhost:11434`.
    - Recommended models: `llama3` or `mistral`.
    - Pull a model: `ollama pull llama3`

## 📦 Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/shell-sage.git
cd shell-sage

# Install dependencies
go mod tidy

# Build the binary
go build -o ssage main.go

# (Optional) Move to your PATH
mv ssage /usr/local/bin/
```

## 🎮 Usage

### 1. Explain a Command
Don't know what `tar -xzvf` does? Ask the sage!

```bash
ssage explain "tar -xzvf archive.tar.gz"
```

### 2. Fix a Mistake
Did your last command fail?

```bash
ssage fix
```
*Note: This reads from your shell history (`.zsh_history` or `.bash_history`).*

### 3. Analyze Logs
Have a huge error log? Get the gist of it.

```bash
ssage analyze ./server.log
```

## 🏗️ Project Structure

```
shell-sage/
├── cmd/            # Cobra CLI commands
│   ├── root.go     # Root command
│   ├── explain.go  # Explain logic
│   ├── fix.go      # Fix logic
│   └── analyze.go  # Analyze logic
├── internal/       # Internal packages
│   ├── ollama/     # Ollama API client
│   └── history/    # Shell history parser
├── main.go         # Entry point
└── go.mod          # Go module definition
```

## 🤝 Contributing

Pull requests are welcome! Please make sure to update tests as appropriate.

## 📄 License

[MIT](https://choosealicense.com/licenses/mit/)
