// Package provider defines the Provider interface and a global registry that
// decouples command code from any concrete AI backend implementation.
//
// Backends (e.g. Ollama) register themselves via their own init() functions
// using the database/sql driver pattern:
//
//	func init() {
//	    provider.Register("ollama", func(model string) (provider.Provider, error) {
//	        return NewClient(model), nil
//	    })
//	}
//
// Call sites create providers through the registry:
//
//	p, err := provider.New("ollama", "llama3")
package provider

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/shell-sage/internal/config"
)

// Provider is the interface all AI backends must implement.
type Provider interface {
	// Generate sends a prompt and returns the full response synchronously.
	Generate(prompt string) (string, error)

	// GenerateStream sends a prompt and calls onChunk for each token received.
	// Returns the full accumulated response string so callers can use it for
	// clipboard copy or caching.
	GenerateStream(prompt string, onChunk func(string)) (string, error)

	// Name returns the unique identifier of this backend (e.g. "ollama").
	Name() string
}

// Factory is a constructor function that creates a Provider for a given model.
type Factory func(model string) (Provider, error)

var (
	mu       sync.RWMutex
	registry = make(map[string]Factory)
)

// Register adds a provider factory under the given name.
// It panics if the same name is registered twice, mirroring the contract used
// by database/sql and the image/* packages in the standard library.
func Register(name string, factory Factory) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := registry[name]; dup {
		panic(fmt.Sprintf("provider: provider %q already registered", name))
	}
	registry[name] = factory
}

// New creates a Provider by name and model string, applying the standard
// priority chain: nameOverride > SSAGE_PROVIDER env > config file > "ollama".
// The model string is forwarded to the factory unchanged.
func New(nameOverride, model string) (Provider, error) {
	name := resolveName(nameOverride)
	mu.RLock()
	factory, ok := registry[name]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf(
			"provider: unknown provider %q — available: %v\n"+
				"  → Import the provider package or register a factory first",
			name, Available(),
		)
	}
	return factory(model)
}

// Available returns a sorted list of all registered provider names.
func Available() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// resolveName applies the priority chain and returns the effective provider name.
func resolveName(override string) string {
	if override != "" {
		return override
	}
	if env := os.Getenv("SSAGE_PROVIDER"); env != "" {
		return env
	}
	if cfg, err := config.Load(); err == nil && cfg.Provider != "" {
		return cfg.Provider
	}
	return "ollama"
}
