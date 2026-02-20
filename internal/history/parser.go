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
// It supports bash, zsh, and PowerShell history on Windows.
func GetRecentCommands(limit int) ([]string, error) {
	historyFile := getHistoryFilePath()
	if historyFile == "" {
		return nil, fmt.Errorf("could not determine history file path")
	}

	file, err := os.Open(historyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open history file at %s: %w", historyFile, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Clean zsh extended history format (": <timestamp>:<elapsed>;<command>")
		if strings.HasPrefix(line, ":") {
			parts := strings.SplitN(line, ";", 2)
			if len(parts) == 2 {
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

	// Return the last 'limit' lines
	start := len(lines) - limit
	if start < 0 {
		start = 0
	}

	return lines[start:], nil
}

func getHistoryFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Windows: check PowerShell history first, then fall back to unix-style
	if runtime.GOOS == "windows" {
		// PowerShell history file location (PSReadLine module)
		appData := os.Getenv("APPDATA")
		if appData != "" {
			psHistory := filepath.Join(appData, "Microsoft", "Windows", "PowerShell", "PSReadLine", "ConsoleHost_history.txt")
			if _, err := os.Stat(psHistory); err == nil {
				return psHistory
			}
		}
		// Fallback: Git Bash / WSL-style history
		zshPath := filepath.Join(home, ".zsh_history")
		if _, err := os.Stat(zshPath); err == nil {
			return zshPath
		}
		bashPath := filepath.Join(home, ".bash_history")
		if _, err := os.Stat(bashPath); err == nil {
			return bashPath
		}
		return ""
	}

	// Unix/macOS: check SHELL env var first
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		return filepath.Join(home, ".zsh_history")
	}
	if strings.Contains(shell, "bash") {
		return filepath.Join(home, ".bash_history")
	}

	// Default fallback
	return filepath.Join(home, ".bash_history")
}
