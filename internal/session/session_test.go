package session

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewAndStoreRoundTrip(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)
	sess, err := New(StartOptions{
		Provider: "codex",
		RepoPath: ".",
		Prompt:   "Inspect this repository",
		Now:      now,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if !strings.HasPrefix(sess.ID, "sess-20260527T120000Z") {
		t.Fatalf("session id missing timestamp prefix: %q", sess.ID)
	}

	store := Store{Root: root}
	if err := store.Append(sess); err != nil {
		t.Fatalf("Append: %v", err)
	}
	got, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 1 || got[0].ID != sess.ID {
		t.Fatalf("round trip mismatch: %+v", got)
	}
	if filepath.Base(store.Path()) != "sessions.jsonl" {
		t.Fatalf("unexpected store path: %s", store.Path())
	}
}

func TestNewRequiresPrompt(t *testing.T) {
	_, err := New(StartOptions{Provider: "codex", RepoPath: ".", Prompt: ""})
	if err == nil {
		t.Fatal("expected prompt validation error")
	}
}
