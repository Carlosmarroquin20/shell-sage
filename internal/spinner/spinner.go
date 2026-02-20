package spinner

import (
	"fmt"
	"time"
)

// Spinner shows an animated spinner in the terminal while waiting for a result.
type Spinner struct {
	frames []string
	stop   chan struct{}
	label  string
}

// New creates a new Spinner with the given label text.
func New(label string) *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stop:   make(chan struct{}),
		label:  label,
	}
}

// Start begins the spinner animation in a background goroutine.
func (s *Spinner) Start() {
	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				// Clear the spinner line on exit
				fmt.Print("\r\033[K")
				return
			default:
				fmt.Printf("\r\033[36m%s\033[0m %s", s.frames[i%len(s.frames)], s.label)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// Stop terminates the spinner and clears its line.
func (s *Spinner) Stop() {
	s.stop <- struct{}{}
}
