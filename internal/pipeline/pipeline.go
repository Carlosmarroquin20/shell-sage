// Package pipeline provides a composable middleware chain for AI generation
// requests. Middlewares wrap the core provider call to add cross-cutting
// behavior such as prompt enhancement, response caching, and retry logic.
//
// Usage:
//
//	p, _ := provider.New("ollama", "llama3")
//	pipe := pipeline.New(p,
//	    enhancer.New(),
//	    cache.New(24*time.Hour, "tip"),
//	    retry.New(3),
//	)
//	response, err := pipe.RunStream(prompt, "explain", func(token string) {
//	    fmt.Print(token)
//	})
package pipeline

import "github.com/shell-sage/internal/provider"

// Request carries the prompt and metadata through the middleware chain.
type Request struct {
	// Prompt is the full text sent to the AI backend. Middlewares (e.g.
	// enhancer) may modify it before it reaches the provider.
	Prompt string

	// Command is the ssage command name that initiated this request
	// (e.g. "explain", "fix", "analyze", "tip"). Middlewares may use it
	// to apply command-specific logic (e.g. cache skip-lists).
	Command string
}

// Handler is the function type for non-streaming invocations.
// Each middleware wraps the next Handler in the chain.
type Handler func(req Request) (string, error)

// StreamHandler is the function type for streaming invocations.
// onChunk is called once per token received from the backend.
type StreamHandler func(req Request, onChunk func(string)) (string, error)

// Middleware adds cross-cutting behavior around a Handler and StreamHandler.
// Implementations must be stateless (or thread-safe) as a single instance is
// shared across all requests.
type Middleware interface {
	// Wrap returns a new Handler that wraps next with additional behavior.
	Wrap(next Handler) Handler

	// WrapStream returns a new StreamHandler that wraps next with additional
	// behavior.
	WrapStream(next StreamHandler) StreamHandler
}

// Pipeline chains a provider with an ordered list of middlewares.
// The first middleware in the slice is the outermost layer (runs first on
// ingress, last on egress). This is the standard onion/Russian-doll model.
type Pipeline struct {
	handler       Handler
	streamHandler StreamHandler
}

// New builds a Pipeline from a provider and zero or more middlewares.
// Calling New with no middlewares creates a direct pass-through to the provider.
func New(p provider.Provider, middlewares ...Middleware) *Pipeline {
	// Terminal handlers that delegate directly to the provider.
	baseH := Handler(func(req Request) (string, error) {
		return p.Generate(req.Prompt)
	})
	baseSH := StreamHandler(func(req Request, onChunk func(string)) (string, error) {
		return p.GenerateStream(req.Prompt, onChunk)
	})

	// Wrap from last to first so that middlewares[0] is outermost.
	h := baseH
	sh := baseSH
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i].Wrap(h)
		sh = middlewares[i].WrapStream(sh)
	}

	return &Pipeline{handler: h, streamHandler: sh}
}

// Run executes the full middleware chain for a non-streaming request and
// returns the complete response.
func (p *Pipeline) Run(prompt, command string) (string, error) {
	return p.handler(Request{Prompt: prompt, Command: command})
}

// RunStream executes the full middleware chain for a streaming request.
// onChunk is called once per token. Returns the full accumulated response.
func (p *Pipeline) RunStream(prompt, command string, onChunk func(string)) (string, error) {
	return p.streamHandler(Request{Prompt: prompt, Command: command}, onChunk)
}
