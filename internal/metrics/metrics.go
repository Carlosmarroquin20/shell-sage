// Package metrics provides local persistent metrics for shell-sage commands.
// Stats are stored in ~/.ssage_metrics.json and updated after each command run.
package metrics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// CommandStats holds aggregated statistics for a single command.
type CommandStats struct {
	Runs        int       `json:"runs"`
	Failures    int       `json:"failures"`
	TotalTimeMs int64     `json:"total_time_ms"`
	AvgTimeMs   int64     `json:"avg_time_ms"`
	LastRun     time.Time `json:"last_run"`
	LastError   string    `json:"last_error,omitempty"`
}

// Store holds stats for every command keyed by command name.
type Store map[string]*CommandStats

func metricsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ssage_metrics.json")
}

// Load reads the metrics file from disk, or returns an empty store if it doesn't exist.
func Load() Store {
	data, err := os.ReadFile(metricsPath())
	if err != nil {
		return make(Store)
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return make(Store)
	}
	return s
}

// Save persists the metrics store to disk.
func (s Store) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(metricsPath(), data, 0644)
}

// Record updates stats for the given command after a run.
// elapsed is the time the full command took. errMsg is empty on success.
func Record(cmd string, elapsed time.Duration, errMsg string) {
	s := Load()

	if s[cmd] == nil {
		s[cmd] = &CommandStats{}
	}
	stat := s[cmd]
	stat.Runs++
	stat.TotalTimeMs += elapsed.Milliseconds()
	stat.AvgTimeMs = stat.TotalTimeMs / int64(stat.Runs)
	stat.LastRun = time.Now()

	if errMsg != "" {
		stat.Failures++
		stat.LastError = errMsg
	}

	_ = s.Save() // Best-effort — don't crash if we can't write metrics
}
