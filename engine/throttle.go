package engine

import (
	"sync"
	"time"
)

// Throttler limits the rate at which a function is called.
// Supports both throttle (fire at most once per interval) and debounce (fire after quiet period).
type Throttler struct {
	mu       sync.Mutex
	interval time.Duration
	mode     string // "throttle" or "debounce"
	last     time.Time
	timer    *time.Timer
	pending  func() // debounce: the latest pending call
}

// NewThrottler creates a rate limiter.
// mode: "throttle" fires immediately then drops until interval passes.
// mode: "debounce" delays firing until interval passes with no new calls.
func NewThrottler(intervalMs int, mode string) *Throttler {
	return &Throttler{
		interval: time.Duration(intervalMs) * time.Millisecond,
		mode:     mode,
	}
}

// Call attempts to execute fn, respecting the rate limit.
func (t *Throttler) Call(fn func()) {
	t.mu.Lock()
	defer t.mu.Unlock()

	switch t.mode {
	case "debounce":
		// Reset the timer each call; fire only after quiet period
		t.pending = fn
		if t.timer != nil {
			t.timer.Stop()
		}
		t.timer = time.AfterFunc(t.interval, func() {
			t.mu.Lock()
			p := t.pending
			t.pending = nil
			t.mu.Unlock()
			if p != nil {
				p()
			}
		})

	default: // "throttle"
		now := time.Now()
		if now.Sub(t.last) >= t.interval {
			t.last = now
			fn()
		}
		// else: dropped
	}
}

// Stop cancels any pending debounce timer.
func (t *Throttler) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}
	t.pending = nil
}
