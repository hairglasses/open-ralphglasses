package mcpadapter

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/hairglasses/open-ralphglasses/internal/chatevents"
	"github.com/hairglasses/open-ralphglasses/internal/sessionlog"
)

func TestCallBudgetEstimate(t *testing.T) {
	result := Call(context.Background(), "open_ralph_budget_estimate", map[string]any{
		"provider":      "codex",
		"input_tokens":  1000,
		"output_tokens": 500,
	})
	if !result.OK {
		t.Fatalf("result=%+v", result)
	}
	payload, ok := result.Payload.(map[string]any)
	if !ok || payload["estimate"] == nil {
		t.Fatalf("payload=%#v", result.Payload)
	}
}

func TestCallUnknownTool(t *testing.T) {
	result := Call(context.Background(), "missing_tool", nil)
	if result.OK || !strings.Contains(result.Error, "unknown public tool") {
		t.Fatalf("result=%+v", result)
	}
}

func TestCallProcessRun(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go executable not on PATH")
	}
	result := Call(context.Background(), "open_ralph_process_run", map[string]any{
		"repo":            ".",
		"timeout_seconds": 5,
		"command":         []any{"go", "version"},
	})
	if !result.OK {
		t.Fatalf("result=%+v", result)
	}
}

func TestCallSessionAnalyze(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 5, 27, 1, 2, 3, 0, time.UTC)
	err := (sessionlog.Store{Root: root}).Save(sessionlog.PersistedSession{
		Snapshot: sessionlog.Snapshot{ID: "sess-example", ProviderSessionID: "provider-session"},
		Transcript: []chatevents.Event{
			{Kind: chatevents.KindOperatorMessage, Text: "Inspect", At: now},
			{Kind: chatevents.KindDelta, Channel: chatevents.ChannelText, Text: "Done", At: now.Add(time.Second)},
		},
	})
	if err != nil {
		t.Fatalf("save transcript: %v", err)
	}
	result := Call(context.Background(), "open_ralph_session_analyze", map[string]any{
		"root": root,
		"id":   "sess-example",
	})
	if !result.OK {
		t.Fatalf("result=%+v", result)
	}
}
