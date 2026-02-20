package history

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

// parseFile is a test helper that runs the same parsing logic as GetRecentCommands
// but on an explicit file path, bypassing OS/shell auto-detection.
func parseFile(path string, limit int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, ":") {
			if parts := strings.SplitN(line, ";", 2); len(parts) == 2 {
				line = parts[1]
			}
		}
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	start := len(lines) - limit
	if start < 0 {
		start = 0
	}
	return lines[start:], nil
}

func writeTempHistory(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "history")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("could not write temp file: %v", err)
	}
	return f.Name()
}

// TestParseFile_HappyPath verifies the last N commands are returned correctly.
func TestParseFile_HappyPath(t *testing.T) {
	path := writeTempHistory(t, "git status\nls -la\necho hello\ndocker ps\ngo build\n")
	result, err := parseFile(path, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 commands, got %d: %v", len(result), result)
	}
	expected := []string{"echo hello", "docker ps", "go build"}
	for i, cmd := range result {
		if cmd != expected[i] {
			t.Errorf("commands[%d]: want %q, got %q", i, expected[i], cmd)
		}
	}
}

// TestParseFile_ZshTimestamps verifies that zsh extended history timestamps are stripped.
func TestParseFile_ZshTimestamps(t *testing.T) {
	path := writeTempHistory(t, ": 1700000000:0;git pull\n: 1700000001:0;npm install\n")
	result, err := parseFile(path, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 commands, got %d: %v", len(result), result)
	}
	if result[0] != "git pull" || result[1] != "npm install" {
		t.Errorf("unexpected commands: %v", result)
	}
}

// TestParseFile_EmptyFile verifies that an empty file returns an empty slice.
func TestParseFile_EmptyFile(t *testing.T) {
	path := writeTempHistory(t, "")
	result, err := parseFile(path, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

// TestParseFile_FileNotFound verifies that a non-existent file returns an error.
func TestParseFile_FileNotFound(t *testing.T) {
	_, err := parseFile("/non/existent/path/history", 5)
	if err == nil {
		t.Error("expected an error for missing file, got nil")
	}
}

// TestParseFile_LimitGreaterThanLines verifies behaviour when limit > number of lines.
func TestParseFile_LimitGreaterThanLines(t *testing.T) {
	path := writeTempHistory(t, "cmd1\ncmd2\n")
	result, err := parseFile(path, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 commands, got %d", len(result))
	}
}
