// Package cache provides a pipeline.Middleware that persists AI responses to
// disk and replays them on subsequent identical requests, avoiding redundant
// network calls to the AI backend.
//
// Cache entries are stored in ~/.ssage_cache/<sha256-of-prompt>.json and
// expire after a configurable TTL. All disk operations fail silently so the
// middleware degrades gracefully to a transparent pass-through when the
// filesystem is unavailable.
//
// Commands in the skip-list (e.g. "tip") always bypass the cache because
// their prompt text is constant but the expected output should vary.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/shell-sage/internal/pipeline"
)

// entry is the on-disk representation of a single cached response.
type entry struct {
	Response  string    `json:"response"`
	CreatedAt time.Time `json:"created_at"`
}

// Middleware implements SHA256-keyed disk caching with a configurable TTL.
type Middleware struct {
	dir      string
	ttl      time.Duration
	skipCmds map[string]bool
}

// New creates a cache Middleware.
//   - ttl controls how long entries remain valid. After expiry the entry is
//     evicted on the next read and the provider is called again.
//   - skipCommands lists command names whose results must never be cached.
//     Useful for commands like "tip" whose prompt is identical every run but
//     whose output is expected to vary.
func New(ttl time.Duration, skipCommands ...string) *Middleware {
	skip := make(map[string]bool, len(skipCommands))
	for _, cmd := range skipCommands {
		skip[cmd] = true
	}

	dir := cacheDir()
	// Best-effort: if the directory can't be created, all cache operations
	// become no-ops (reads return miss, writes are silently dropped).
	_ = os.MkdirAll(dir, 0700)

	return &Middleware{dir: dir, ttl: ttl, skipCmds: skip}
}

// Wrap caches the response of non-streaming requests.
func (m *Middleware) Wrap(next pipeline.Handler) pipeline.Handler {
	return func(req pipeline.Request) (string, error) {
		if m.skipCmds[req.Command] {
			return next(req)
		}
		key := hashKey(req.Prompt)
		if cached, ok := m.load(key); ok {
			return cached, nil
		}
		resp, err := next(req)
		if err == nil {
			m.store(key, resp)
		}
		return resp, err
	}
}

// WrapStream caches the full response of streaming requests. On a cache hit
// the stored response is replayed as a single onChunk call so that command
// output logic (spinner, box drawing, clipboard) behaves identically.
func (m *Middleware) WrapStream(next pipeline.StreamHandler) pipeline.StreamHandler {
	return func(req pipeline.Request, onChunk func(string)) (string, error) {
		if m.skipCmds[req.Command] {
			return next(req, onChunk)
		}
		key := hashKey(req.Prompt)
		if cached, ok := m.load(key); ok {
			onChunk(cached)
			return cached, nil
		}
		resp, err := next(req, onChunk)
		if err == nil {
			m.store(key, resp)
		}
		return resp, err
	}
}

// hashKey returns the lowercase hex SHA-256 of the prompt, used as the
// cache filename (without extension).
func hashKey(prompt string) string {
	sum := sha256.Sum256([]byte(prompt))
	return hex.EncodeToString(sum[:])
}

// cacheDir returns the path to the cache directory.
func cacheDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), ".ssage_cache")
	}
	return filepath.Join(home, ".ssage_cache")
}

// load reads and validates a cache entry.
// Returns ("", false) on any failure (missing file, JSON error, TTL expired).
func (m *Middleware) load(key string) (string, bool) {
	path := filepath.Join(m.dir, key+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	var e entry
	if err := json.Unmarshal(data, &e); err != nil {
		return "", false
	}
	if time.Since(e.CreatedAt) > m.ttl {
		_ = os.Remove(path) // evict expired entry
		return "", false
	}
	return e.Response, true
}

// store writes a response entry to disk. Errors are silently dropped.
func (m *Middleware) store(key, response string) {
	e := entry{Response: response, CreatedAt: time.Now()}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return
	}
	path := filepath.Join(m.dir, key+".json")
	_ = os.WriteFile(path, data, 0600)
}
