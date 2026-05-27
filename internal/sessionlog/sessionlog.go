// Package sessionlog stores inspectable public session transcripts.
//
// This package deliberately sits beside internal/session's planned JSONL
// ledger. Planned sessions describe launch intent; transcript logs describe
// provider-neutral events captured by a future runner or imported by tests.
package sessionlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/hairglasses/open-ralphglasses/internal/chatevents"
)

const schemaVersion = 1

// Snapshot is the fast metadata for a persisted transcript.
type Snapshot struct {
	ID                string    `json:"id"`
	Provider          string    `json:"provider,omitempty"`
	ProviderSessionID string    `json:"provider_session_id,omitempty"`
	Model             string    `json:"model,omitempty"`
	RepoPath          string    `json:"repo_path,omitempty"`
	Status            string    `json:"status,omitempty"`
	StartedAt         time.Time `json:"started_at,omitempty"`
	LastActivity      time.Time `json:"last_activity,omitempty"`
}

// PersistedSession is the on-disk transcript artifact.
type PersistedSession struct {
	SchemaVersion int                `json:"schema_version"`
	Snapshot      Snapshot           `json:"snapshot"`
	Transcript    []chatevents.Event `json:"transcript"`
}

// Store reads and writes transcript JSON under Root.
type Store struct {
	Root string
}

// Analysis is a deterministic summary of a transcript.
type Analysis struct {
	Snapshot             Snapshot       `json:"snapshot"`
	EventCount           int            `json:"event_count"`
	KindCounts           map[string]int `json:"kind_counts"`
	OperatorMessageCount int            `json:"operator_message_count"`
	AssistantDeltaCount  int            `json:"assistant_delta_count"`
	AssistantTextBytes   int            `json:"assistant_text_bytes"`
	ToolCallCount        int            `json:"tool_call_count"`
	ToolNames            []string       `json:"tool_names,omitempty"`
	ErrorCount           int            `json:"error_count"`
	Errors               []string       `json:"errors,omitempty"`
	InputTokens          int            `json:"input_tokens"`
	OutputTokens         int            `json:"output_tokens"`
	FirstEventAt         *time.Time     `json:"first_event_at,omitempty"`
	LastEventAt          *time.Time     `json:"last_event_at,omitempty"`
	ReplayReady          bool           `json:"replay_ready"`
	ReplayBlockers       []string       `json:"replay_blockers,omitempty"`
}

// Save writes a transcript atomically.
func (s Store) Save(ps PersistedSession) error {
	if strings.TrimSpace(ps.Snapshot.ID) == "" {
		return fmt.Errorf("snapshot id is required")
	}
	ps.SchemaVersion = schemaVersion
	if ps.Snapshot.Status == "" {
		ps.Snapshot.Status = inferStatus(ps.Transcript)
	}
	fillActivity(&ps)
	path := s.Path(ps.Snapshot.ID)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create transcript dir: %w", err)
	}
	body, err := json.MarshalIndent(ps, "", "  ")
	if err != nil {
		return fmt.Errorf("encode transcript: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, body, 0o600); err != nil {
		return fmt.Errorf("write transcript tmp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename transcript tmp: %w", err)
	}
	return nil
}

// Load reads one transcript by id.
func (s Store) Load(id string) (PersistedSession, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return PersistedSession{}, fmt.Errorf("session id is required")
	}
	path := s.Path(id)
	body, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return PersistedSession{}, fmt.Errorf("session transcript %q not found", id)
		}
		return PersistedSession{}, fmt.Errorf("read transcript: %w", err)
	}
	var ps PersistedSession
	if err := json.Unmarshal(body, &ps); err != nil {
		return PersistedSession{}, fmt.Errorf("decode transcript: %w", err)
	}
	return ps, nil
}

// List returns transcript snapshots newest first.
func (s Store) List() ([]Snapshot, error) {
	root := s.dir()
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read transcript dir: %w", err)
	}
	out := make([]Snapshot, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".json")
		ps, err := s.Load(id)
		if err != nil {
			continue
		}
		out = append(out, ps.Snapshot)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].LastActivity.After(out[j].LastActivity)
	})
	return out, nil
}

// Path returns the transcript JSON path for id.
func (s Store) Path(id string) string {
	return filepath.Join(s.dir(), id+".json")
}

func (s Store) dir() string {
	root := strings.TrimSpace(s.Root)
	if root == "" {
		root = "."
	}
	return filepath.Join(root, ".open-ralph", "transcripts")
}

// Analyze summarizes a persisted transcript.
func Analyze(ps PersistedSession) Analysis {
	out := Analysis{
		Snapshot:   ps.Snapshot,
		KindCounts: make(map[string]int),
	}
	toolSeen := make(map[string]bool)
	for _, ev := range ps.Transcript {
		out.EventCount++
		out.KindCounts[string(ev.Kind)]++
		if !ev.At.IsZero() {
			at := ev.At
			if out.FirstEventAt == nil || at.Before(*out.FirstEventAt) {
				out.FirstEventAt = &at
			}
			if out.LastEventAt == nil || at.After(*out.LastEventAt) {
				out.LastEventAt = &at
			}
		}
		if ev.InputTokens > 0 {
			out.InputTokens = ev.InputTokens
		}
		if ev.OutputTokens > 0 {
			out.OutputTokens = ev.OutputTokens
		}
		switch ev.Kind {
		case chatevents.KindOperatorMessage:
			out.OperatorMessageCount++
		case chatevents.KindDelta:
			if ev.Channel == "" || ev.Channel == chatevents.ChannelText {
				out.AssistantDeltaCount++
				out.AssistantTextBytes += len(ev.Text)
			}
		case chatevents.KindToolUseStart:
			out.ToolCallCount++
			if ev.ToolName != "" && !toolSeen[ev.ToolName] {
				toolSeen[ev.ToolName] = true
				out.ToolNames = append(out.ToolNames, ev.ToolName)
			}
		case chatevents.KindError:
			out.ErrorCount++
			if ev.Error != "" {
				out.Errors = append(out.Errors, ev.Error)
			}
		case chatevents.KindEnd:
			if ev.Error != "" {
				out.ErrorCount++
				out.Errors = append(out.Errors, ev.Error)
			}
		}
	}
	slices.Sort(out.ToolNames)
	out.ReplayBlockers = replayBlockers(ps.Snapshot, out)
	out.ReplayReady = len(out.ReplayBlockers) == 0
	return out
}

// RenderReplayText returns a compact human-readable replay.
func RenderReplayText(events []chatevents.Event) string {
	var b strings.Builder
	for _, ev := range events {
		switch ev.Kind {
		case chatevents.KindOperatorMessage:
			appendBlock(&b, "Operator", ev.Text)
		case chatevents.KindDelta:
			if ev.Channel == "" || ev.Channel == chatevents.ChannelText {
				appendBlock(&b, "Assistant", ev.Text)
			}
		case chatevents.KindToolUseStart:
			appendBlock(&b, "Tool start", firstNonEmpty(ev.ToolName, ev.ToolUseID))
		case chatevents.KindToolUseEnd:
			label := "Tool end"
			if !ev.ToolOK {
				label = "Tool error"
			}
			appendBlock(&b, label, firstNonEmpty(ev.ToolName, ev.ToolUseID))
		case chatevents.KindError:
			appendBlock(&b, "Error", ev.Error)
		case chatevents.KindEnd:
			appendBlock(&b, "End", ev.StopReason)
		}
	}
	return strings.TrimSpace(b.String())
}

func replayBlockers(snapshot Snapshot, analysis Analysis) []string {
	var blockers []string
	if analysis.EventCount == 0 {
		blockers = append(blockers, "empty_transcript")
	}
	if analysis.OperatorMessageCount == 0 {
		blockers = append(blockers, "missing_operator_message")
	}
	if analysis.AssistantDeltaCount == 0 {
		blockers = append(blockers, "missing_assistant_text")
	}
	if snapshot.ProviderSessionID == "" {
		blockers = append(blockers, "missing_provider_session_id")
	}
	if analysis.ErrorCount > 0 {
		blockers = append(blockers, "contains_errors")
	}
	return blockers
}

func fillActivity(ps *PersistedSession) {
	for _, ev := range ps.Transcript {
		if ev.At.IsZero() {
			continue
		}
		if ps.Snapshot.StartedAt.IsZero() || ev.At.Before(ps.Snapshot.StartedAt) {
			ps.Snapshot.StartedAt = ev.At
		}
		if ps.Snapshot.LastActivity.IsZero() || ev.At.After(ps.Snapshot.LastActivity) {
			ps.Snapshot.LastActivity = ev.At
		}
	}
}

func inferStatus(events []chatevents.Event) string {
	for i := len(events) - 1; i >= 0; i-- {
		switch events[i].Kind {
		case chatevents.KindError:
			return "error"
		case chatevents.KindEnd:
			return "ended"
		}
	}
	return "captured"
}

func appendBlock(b *strings.Builder, label, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	if b.Len() > 0 {
		b.WriteString("\n\n")
	}
	b.WriteString(label)
	b.WriteString(":\n")
	b.WriteString(text)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
