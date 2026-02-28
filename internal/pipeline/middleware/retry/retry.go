// Package retry provides a pipeline.Middleware that automatically retries
// failed AI generation requests with exponential backoff.
//
// For non-streaming requests all attempts are transparent to the caller.
//
// For streaming requests a retry is only performed if no tokens have been
// delivered to the caller yet. Once the first token has been sent, the caller
// has already begun rendering output and a retry would produce duplicate
// content, so the error is returned as-is.
package retry

import (
	"time"

	"github.com/shell-sage/internal/pipeline"
)

// Middleware retries failed requests up to maxAttempts times.
type Middleware struct {
	maxAttempts int
}

// New creates a retry Middleware.
// maxAttempts is the total number of attempts including the first one, so
// New(3) means one initial attempt followed by up to two retries.
// Values less than 1 are clamped to 1 (no retry).
func New(maxAttempts int) *Middleware {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	return &Middleware{maxAttempts: maxAttempts}
}

// Wrap retries the non-streaming handler on any error, sleeping between
// attempts using exponential backoff (500 ms, 1 s, 2 s, …).
func (m *Middleware) Wrap(next pipeline.Handler) pipeline.Handler {
	return func(req pipeline.Request) (string, error) {
		var (
			resp string
			err  error
		)
		for attempt := 0; attempt < m.maxAttempts; attempt++ {
			if attempt > 0 {
				time.Sleep(backoff(attempt - 1))
			}
			resp, err = next(req)
			if err == nil {
				return resp, nil
			}
		}
		return resp, err
	}
}

// WrapStream retries the streaming handler but only when no tokens have been
// received yet. Once the first token arrives, the caller has started rendering
// so we cannot restart the stream without corrupting the output.
func (m *Middleware) WrapStream(next pipeline.StreamHandler) pipeline.StreamHandler {
	return func(req pipeline.Request, onChunk func(string)) (string, error) {
		var (
			resp     string
			err      error
			received bool
		)
		for attempt := 0; attempt < m.maxAttempts; attempt++ {
			if attempt > 0 {
				time.Sleep(backoff(attempt - 1))
			}
			received = false

			// Guard onChunk to detect whether the stream has started.
			guarded := func(token string) {
				received = true
				onChunk(token)
			}

			resp, err = next(req, guarded)
			if err == nil {
				return resp, nil
			}
			if received {
				// Tokens already sent downstream — cannot retry cleanly.
				return resp, err
			}
		}
		return resp, err
	}
}

// backoff returns the sleep duration before the given retry attempt (0-based).
// Sequence: 500 ms, 1 s, 2 s, 4 s, …
func backoff(attempt int) time.Duration {
	d := 500 * time.Millisecond
	for i := 0; i < attempt; i++ {
		d *= 2
	}
	return d
}
