package sessionlog

import (
	"strings"
	"testing"
	"time"

	"github.com/hairglasses/open-ralphglasses/internal/chatevents"
)

func TestStoreRoundTripAndAnalyze(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 5, 27, 1, 2, 3, 0, time.UTC)
	ps := PersistedSession{
		Snapshot: Snapshot{
			ID:                "sess-example",
			Provider:          "codex",
			ProviderSessionID: "provider-session",
			Model:             "gpt-5",
		},
		Transcript: []chatevents.Event{
			{Kind: chatevents.KindOperatorMessage, SessionID: "sess-example", Text: "Inspect this repo", At: now},
			{Kind: chatevents.KindDelta, SessionID: "sess-example", Channel: chatevents.ChannelText, Text: "Done", At: now.Add(time.Second)},
			{Kind: chatevents.KindToolUseStart, SessionID: "sess-example", ToolName: "Read", At: now.Add(2 * time.Second)},
			{Kind: chatevents.KindEnd, SessionID: "sess-example", StopReason: "end_turn", At: now.Add(3 * time.Second)},
		},
	}
	store := Store{Root: root}
	if err := store.Save(ps); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := store.Load("sess-example")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	analysis := Analyze(loaded)
	if !analysis.ReplayReady || analysis.ToolCallCount != 1 || analysis.AssistantTextBytes != len("Done") {
		t.Fatalf("analysis=%+v", analysis)
	}
	replay := RenderReplayText(loaded.Transcript)
	if !strings.Contains(replay, "Operator:") || !strings.Contains(replay, "Assistant:") {
		t.Fatalf("replay=%q", replay)
	}
}

func TestListNewestFirst(t *testing.T) {
	root := t.TempDir()
	store := Store{Root: root}
	older := PersistedSession{Snapshot: Snapshot{ID: "older", LastActivity: time.Unix(10, 0)}}
	newer := PersistedSession{Snapshot: Snapshot{ID: "newer", LastActivity: time.Unix(20, 0)}}
	if err := store.Save(older); err != nil {
		t.Fatalf("save older: %v", err)
	}
	if err := store.Save(newer); err != nil {
		t.Fatalf("save newer: %v", err)
	}
	snapshots, err := store.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(snapshots) != 2 || snapshots[0].ID != "newer" {
		t.Fatalf("snapshots=%+v", snapshots)
	}
}

func TestEventStructuralSemantics(t *testing.T) {
	if (chatevents.Event{Kind: chatevents.KindDelta}).IsStructural() {
		t.Fatal("delta should not be structural")
	}
	if !(chatevents.Event{Kind: chatevents.KindToolUseStart}).IsStructural() {
		t.Fatal("tool start should be structural")
	}
}
