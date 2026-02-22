# 🧙‍♂️ Shell Sage (ssage)

[![Go Version](https://img.shields.io/github/go-mod/go-version/Carlosmarroquin20/shell-sage)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Ollama Powered](https://img.shields.io/badge/Powered%20By-Ollama-orange.svg)](https://ollama.com/)

**Shell Sage** is your AI-powered terminal companion. Built for developers who live in the CLI, it leverages local LLMs (via Ollama) to bring intelligent command explanation, error fixing, and log analysis directly to your fingertips—without your data ever leaving your machine.

---

## 🚀 Core Features

### 🔍 `ssage explain "[command]"`
Unpack complex one-liners. Shell Sage breaks down flags and arguments into human-readable steps.
```bash
ssage explain "tar -xzvf archive.tar.gz"
```

### 🛠️ `ssage fix`
The "Crown Jewel". Scans your recent history, detects the last failed command, and provides an AI-suggested fix with an explanation.
```bash
ssage fix
```

### 📊 `ssage analyze [file]`
Don't drown in logs. Point the Sage at an error log, and it will summarize the root cause and suggest potential solutions.
```bash
ssage analyze ./build.log
```

### 💡 `ssage tip`
Feeling lucky? Get a random, high-productivity terminal tip or trick to level up your shell game.
```bash
ssage tip
```

### 📈 `ssage stats`
Track your growth. View metrics on how many commands you've explained, fixed, and analyzed.
```bash
ssage stats
```

---

## ⚙️ Global Power-Ups

Customize how the Sage speaks to you using persistent flags:

- **`--model, -m`**: Choose your brain. Works with any model you've pulled in Ollama (e.g., `llama3`, `mistral`, `codellama`).
- **`--lang, -l`**: Prefer another language? Set it globally (e.g., `--lang es` for Spanish, `--lang fr` for French).

---

## 💻 Shell Compatibility

Shell Sage plays well with your favorite environment:

| Shell | Status | Note |
| :--- | :--- | :--- |
| **Zsh** | ✅ Full Support | Uses `~/.zsh_history` |
| **Bash** | ✅ Full Support | Uses `~/.bash_history` |
| **Fish** | ✅ Full Support | Uses `~/.local/share/fish/fish_history` |
| **PowerShell** | ✅ Full Support | Windows only, requires `PSReadLine` |

---

## 🛠️ Installation

### Prerequisites
1. **Go 1.21+**
2. **Ollama**: [Download here](https://ollama.com/).
   - Ensure Ollama is running (`ollama serve`).
   - Pull a model: `ollama pull llama3`.

### Setup
```bash
# Clone and enter
git clone https://github.com/Carlosmarroquin20/shell-sage.git
cd shell-sage

# Build
go build -o ssage main.go

# Add to your PATH (Example for Mac/Linux)
mv ssage /usr/local/bin/
```

---

## 🤝 Contributing

Contributions are what make the open-source community an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

<p align="center">
  <i>"May your terminal always be wise and your errors few."</i>
</p>
