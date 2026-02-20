package history

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrNoHistoryFile is returned when no recognized history file can be found.
type ErrNoHistoryFile struct {
	OS   string
	Tips []string
}

func (e *ErrNoHistoryFile) Error() string {
	msg := fmt.Sprintf("could not find a shell history file on %s.\n", e.OS)
	msg += "  Possible causes and fixes:\n"
	for _, tip := range e.Tips {
		msg += "    • " + tip + "\n"
	}
	return msg
}

// GetRecentCommands reads the last n commands from the shell history file.
// It supports bash, zsh, and PowerShell (via PSReadLine).
func GetRecentCommands(limit int) ([]string, error) {
	historyFile, err := getHistoryFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(historyFile)
	if err != nil {
		return nil, fmt.Errorf("found history file at '%s' but could not open it: %w\n  Try: check file permissions with 'ls -la %s'", historyFile, err, historyFile)
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

	start := len(lines) - limit
	if start < 0 {
		start = 0
	}

	return lines[start:], nil
}

func getHistoryFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	if runtime.GOOS == "windows" {
		// Priority 1: PowerShell PSReadLine history
		appData := os.Getenv("APPDATA")
		if appData != "" {
			psHistory := filepath.Join(appData, "Microsoft", "Windows", "PowerShell", "PSReadLine", "ConsoleHost_history.txt")
			if _, err := os.Stat(psHistory); err == nil {
				return psHistory, nil
			}
		}
		// Priority 2: Git Bash / WSL bash history
		bashPath := filepath.Join(home, ".bash_history")
		if _, err := os.Stat(bashPath); err == nil {
			return bashPath, nil
		}
		// Priority 3: Zsh (rare on Windows but possible via WSL tools)
		zshPath := filepath.Join(home, ".zsh_history")
		if _, err := os.Stat(zshPath); err == nil {
			return zshPath, nil
		}

		return "", &ErrNoHistoryFile{
			OS: "Windows",
			Tips: []string{
				"Make sure PowerShell PSReadLine module is installed: Install-Module PSReadLine",
				"Or install Git for Windows which includes bash history support",
				"Run a few commands in PowerShell first — history is created after the first session",
			},
		}
	}

	// Unix/macOS: check SHELL env var
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		return filepath.Join(home, ".zsh_history"), nil
	}
	if strings.Contains(shell, "bash") {
		return filepath.Join(home, ".bash_history"), nil
	}

	// Fallback: try both
	for _, f := range []string{".zsh_history", ".bash_history"} {
		p := filepath.Join(home, f)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", &ErrNoHistoryFile{
		OS: runtime.GOOS,
		Tips: []string{
			"Set the SHELL environment variable (e.g. export SHELL=/bin/zsh)",
			"Make sure you have run some commands so the history file is created",
			"Check ~/.bash_history or ~/.zsh_history exist and are readable",
		},
	}
}
