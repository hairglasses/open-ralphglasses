package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunProviders(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"providers"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "codex") || !strings.Contains(stdout.String(), "gemini") {
		t.Fatalf("providers output missing catalog: %s", stdout.String())
	}
}

func TestRunMCPManifest(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"mcp", "manifest"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "open_ralph_session_plan") {
		t.Fatalf("manifest output missing session tool: %s", stdout.String())
	}
}
