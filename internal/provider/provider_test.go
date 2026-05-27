package provider

import "testing"

func TestLookupKnownProvider(t *testing.T) {
	p, err := Lookup("Codex")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}
	if p.ID != "codex" || p.Command == "" {
		t.Fatalf("unexpected provider: %+v", p)
	}
}

func TestLookupRejectsUnknownProvider(t *testing.T) {
	if _, err := Lookup("example"); err == nil {
		t.Fatal("expected unknown provider error")
	}
}
