// Package chatevents defines provider-neutral transcript events.
//
// Private ralphglasses uses a richer live event bus. This public package keeps
// the reusable transcript contract: structural events are never safe to drop,
// while additive delta events may be merged by a future streaming layer.
package chatevents

import "time"

// Kind discriminates the event variant.
type Kind string

const (
	KindStart           Kind = "start"
	KindDelta           Kind = "delta"
	KindToolUseStart    Kind = "tool_use_start"
	KindToolUseEnd      Kind = "tool_use_end"
	KindUsage           Kind = "usage"
	KindEnd             Kind = "end"
	KindError           Kind = "error"
	KindOperatorMessage Kind = "operator_message"
)

// Channel labels additive text streams.
type Channel string

const (
	ChannelText     Channel = "text"
	ChannelThinking Channel = "thinking"
)

// Event is one normalized transcript event.
type Event struct {
	Kind              Kind      `json:"kind"`
	SessionID         string    `json:"session_id,omitempty"`
	ProviderSessionID string    `json:"provider_session_id,omitempty"`
	Provider          string    `json:"provider,omitempty"`
	Model             string    `json:"model,omitempty"`
	Channel           Channel   `json:"channel,omitempty"`
	Text              string    `json:"text,omitempty"`
	ToolUseID         string    `json:"tool_use_id,omitempty"`
	ToolName          string    `json:"tool_name,omitempty"`
	ToolInput         string    `json:"tool_input,omitempty"`
	ToolOK            bool      `json:"tool_ok,omitempty"`
	ToolOutput        string    `json:"tool_output,omitempty"`
	InputTokens       int       `json:"input_tokens,omitempty"`
	OutputTokens      int       `json:"output_tokens,omitempty"`
	StopReason        string    `json:"stop_reason,omitempty"`
	Error             string    `json:"error,omitempty"`
	At                time.Time `json:"at,omitempty"`
}

// IsStructural reports whether dropping this event would corrupt replay.
func (e Event) IsStructural() bool {
	switch e.Kind {
	case KindDelta, KindUsage:
		return false
	default:
		return true
	}
}
