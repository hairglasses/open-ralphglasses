package cli

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestRunReposScan(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "example")
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir git dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".open-ralphrc"), []byte("provider=codex\n"), 0o644); err != nil {
		t.Fatalf("write opt-in file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"repos", "scan", "--root", root, "--depth", "2"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, `"name": "example"`) || !strings.Contains(output, `"enabled": true`) {
		t.Fatalf("scan output missing repo: %s", output)
	}
}
