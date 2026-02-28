// Package enhancer provides a pipeline.Middleware that prepends a system
// context block to every prompt before it reaches the AI backend.
//
// The injected block looks like:
//
//	[System context: OS=linux, Arch=amd64, Shell=/bin/zsh]
//
// This helps the model give OS-aware and shell-specific answers without
// requiring each command to manually build context strings.
package enhancer

import (
	"fmt"
	"os"
	"runtime"

	"github.com/shell-sage/internal/pipeline"
)

// Middleware injects runtime OS/shell context into every prompt.
type Middleware struct{}

// New returns a ready-to-use enhancer Middleware.
func New() *Middleware {
	return &Middleware{}
}

// Wrap prepends system context to the prompt before calling next.
func (m *Middleware) Wrap(next pipeline.Handler) pipeline.Handler {
	return func(req pipeline.Request) (string, error) {
		req.Prompt = inject(req.Prompt)
		return next(req)
	}
}

// WrapStream prepends system context to the prompt before calling next.
func (m *Middleware) WrapStream(next pipeline.StreamHandler) pipeline.StreamHandler {
	return func(req pipeline.Request, onChunk func(string)) (string, error) {
		req.Prompt = inject(req.Prompt)
		return next(req, onChunk)
	}
}

// inject builds and prepends the context prefix string.
func inject(prompt string) string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "unknown"
	}
	prefix := fmt.Sprintf(
		"[System context: OS=%s, Arch=%s, Shell=%s]\n",
		runtime.GOOS, runtime.GOARCH, shell,
	)
	return prefix + prompt
}
