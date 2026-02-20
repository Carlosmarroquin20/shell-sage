package history

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetRecentCommands reads the last n commands from the shell history file.
// It supports bash, zsh, and PowerShell (via PSReadLine).
func GetRecentCommands(limit int) ([]string, error) {
	historyFile, err := getHistoryFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(historyFile)
	if err != nil {
		return nil, fmt.Errorf("found history at '%s' but could not open it: %w\n  → Try: check permissions with 'ls -la %s'", historyFile, err, historyFile)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Strip zsh extended history format: ": <timestamp>:<elapsed>;<command>"
		if strings.HasPrefix(line, ":") {
			if parts := strings.SplitN(line, ";", 2); len(parts) == 2 {
				line = parts[1]
			}
		}
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading history file: %w", err)
	}

	start := len(lines) - limit
	if start < 0 {
		start = 0
	}
	return lines[start:], nil
}

// getHistoryFilePath detects the user's shell and returns the appropriate
// history file path. Returns a descriptive, shell-specific error if not found.
func getHistoryFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	if runtime.GOOS == "windows" {
		return windowsHistoryPath(home)
	}
	return unixHistoryPath(home)
}

// windowsHistoryPath checks PowerShell and bash-compatible history on Windows.
func windowsHistoryPath(home string) (string, error) {
	// Priority 1: PowerShell PSReadLine
	appData := os.Getenv("APPDATA")
	if appData != "" {
		psHistory := filepath.Join(appData, "Microsoft", "Windows", "PowerShell", "PSReadLine", "ConsoleHost_history.txt")
		if _, err := os.Stat(psHistory); err == nil {
			return psHistory, nil
		}
	}

	// Priority 2: Git Bash history (~/.bash_history)
	bashPath := filepath.Join(home, ".bash_history")
	if _, err := os.Stat(bashPath); err == nil {
		return bashPath, nil
	}

	// Priority 3: Zsh (rare on Windows but possible)
	zshPath := filepath.Join(home, ".zsh_history")
	if _, err := os.Stat(zshPath); err == nil {
		return zshPath, nil
	}

	return "", fmt.Errorf(
		"no shell history found on Windows.\n" +
			"  PowerShell users:\n" +
			"    → Install PSReadLine: Install-Module PSReadLine -Force\n" +
			"    → Restart PowerShell and run a few commands to create history\n" +
			"  Git Bash users:\n" +
			"    → Make sure Git for Windows is installed (git-scm.com)\n" +
			"    → ~/.bash_history is created automatically after the first session",
	)
}

// unixHistoryPath checks zsh, bash, fish or common fallback history on Linux/macOS.
func unixHistoryPath(home string) (string, error) {
	shell := os.Getenv("SHELL")

	// Zsh
	if strings.Contains(shell, "zsh") {
		p := filepath.Join(home, ".zsh_history")
		if _, err := os.Stat(p); err != nil {
			return "", fmt.Errorf(
				"zsh detected but ~/.zsh_history not found.\n" +
					"  → Check it exists: ls -la ~/.zsh_history\n" +
					"  → Make sure HISTFILE is set in your ~/.zshrc: echo $HISTFILE\n" +
					"  → If empty, add: echo 'HISTFILE=~/.zsh_history' >> ~/.zshrc && source ~/.zshrc",
			)
		}
		return p, nil
	}

	// Bash
	if strings.Contains(shell, "bash") {
		p := filepath.Join(home, ".bash_history")
		if _, err := os.Stat(p); err != nil {
			return "", fmt.Errorf(
				"bash detected but ~/.bash_history not found.\n" +
					"  → Check it exists: ls -la ~/.bash_history\n" +
					"  → Make sure HISTFILE is set: echo 'HISTFILE=~/.bash_history' >> ~/.bashrc && source ~/.bashrc\n" +
					"  → Run a few commands and then try again",
			)
		}
		return p, nil
	}

	// Fish shell — history stored as ~/.local/share/fish/fish_history (YAML-like)
	if strings.Contains(shell, "fish") {
		p := filepath.Join(home, ".local", "share", "fish", "fish_history")
		if _, err := os.Stat(p); err != nil {
			return "", fmt.Errorf(
				"fish shell detected but history file not found.\n" +
					"  → Expected path: ~/.local/share/fish/fish_history\n" +
					"  → Run a few commands in fish and try again",
			)
		}
		return p, nil
	}

	// Unknown $SHELL — try common files as fallback in order of popularity
	for _, f := range []string{".zsh_history", ".bash_history", ".local/share/fish/fish_history"} {
		p := filepath.Join(home, f)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf(
		"could not detect your shell or history file.\n" +
			"  → Set the SHELL environment variable: export SHELL=$(which zsh)\n" +
			"  → Or create the history file manually: touch ~/.bash_history\n" +
			"  → Supported shells: bash, zsh, fish, PowerShell (Windows)",
	)
}
