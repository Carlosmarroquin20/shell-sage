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
// It attempts to detect the shell environment and find the appropriate history file.
func GetRecentCommands(limit int) ([]string, error) {
	historyFile := getHistoryFilePath()
	if historyFile == "" {
		return nil, fmt.Errorf("could not determine history file path")
	}

	file, err := os.Open(historyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Simple cleaning: remove timestamps if present (zsh specific extended history)
		// This is a naive implementation; robust parsing depends on specific shell config.
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

	// Get the last 'limit' lines
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

	// Priority 1: Check SHELL env var
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		return filepath.Join(home, ".zsh_history")
	}
	if strings.Contains(shell, "bash") {
		return filepath.Join(home, ".bash_history")
	}

	// Fallback/Windows logic (Powershell history is more complex, using simple check for now)
	if runtime.GOOS == "windows" {
		// Attempt to find git-bash or generic history if available, 
		// otherwise might return empty or need PowerShell specific reader.
		// For this MVP, we will try to look for common unix-like history files even on Windows
		// assuming the user might be using Git Bash or WSL.
		
		zshPath := filepath.Join(home, ".zsh_history")
		if _, err := os.Stat(zshPath); err == nil {
			return zshPath
		}
		bashPath := filepath.Join(home, ".bash_history")
		if _, err := os.Stat(bashPath); err == nil {
			return bashPath
		}
	}

	// Default fallback
	return filepath.Join(home, ".bash_history")
}
