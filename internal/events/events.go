// Package events provides a small in-memory event bus.
//
// ralphglasses uses typed events to decouple session lifecycle, MCP handlers,
// UI projections, and health checks. This public version keeps the same design
// idea without private event types or runtime state.
package events

import (
	"slices"
	"sync"
	"time"
)

// Type names the category of an event.
type Type string

const (
	SessionPlanned Type = "session.planned"
	SessionListed  Type = "session.listed"
	DoctorChecked  Type = "doctor.checked"
)

// Event is intentionally metadata-only. Public examples should not put prompts,
// credentials, or private paths into event details.
type Event struct {
	Type      Type              `json:"type"`
	Time      time.Time         `json:"time"`
	Details   map[string]string `json:"details,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

// Bus stores a bounded history so CLI, TUI, or MCP adapters can read the same
// recent facts without sharing mutable package globals.
type Bus struct {
	mu      sync.Mutex
	limit   int
	history []Event
}

// NewBus creates a bus with a bounded history ring.
func NewBus(limit int) *Bus {
	if limit <= 0 {
		limit = 100
	}
	return &Bus{limit: limit}
}

// Publish appends an event to history, trimming the oldest event when the ring
// is full.
func (b *Bus) Publish(event Event) {
	if event.Time.IsZero() {
		event.Time = time.Now().UTC()
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.history = append(b.history, event)
	if len(b.history) > b.limit {
		b.history = slices.Clone(b.history[len(b.history)-b.limit:])
	}
}

// History returns a copy so callers cannot mutate the bus internals.
func (b *Bus) History() []Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	return slices.Clone(b.history)
}
