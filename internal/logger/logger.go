// Package logger provides a centralized logrus logger for shell-sage.
// Logs are written to ~/.ssage.log in JSON format so they can be parsed
// by any log aggregation tool later.
package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Log is the shared application logger instance.
var Log *logrus.Logger

func init() {
	Log = logrus.New()

	// Write logs to ~/.ssage.log
	home, err := os.UserHomeDir()
	if err == nil {
		logPath := filepath.Join(home, ".ssage.log")
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			Log.SetOutput(f)
		}
	}

	// JSON format for structured, machine-readable logs
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)
}
