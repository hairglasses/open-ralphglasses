package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hairglasses/open-ralphglasses/internal/chatevents"
	"github.com/hairglasses/open-ralphglasses/internal/sessionlog"
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

func TestRunBudgetEstimate(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{
		"budget", "estimate",
		"--provider", "gemini",
		"--input-tokens", "1000000",
		"--output-tokens", "500000",
		"--budget", "1",
		"--spent", "0.90",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{`"provider": "gemini"`, `"estimated_usd": 0.45`, `"should_stop": true`} {
		if !strings.Contains(output, want) {
			t.Fatalf("budget output missing %q: %s", want, output)
		}
	}
}

func TestRunLaunchPlan(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{
		"launch", "plan",
		"--provider", "codex",
		"--repo", ".",
		"--prompt", "summarize",
		"--permission-mode", "read-only",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{`"provider": "codex"`, `"execution_status": "planned_only"`, `"--sandbox"`} {
		if !strings.Contains(output, want) {
			t.Fatalf("launch plan output missing %q: %s", want, output)
		}
	}
}

func TestRunHookCheck(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{
		"hook", "check",
		"--event", "PreToolUse",
		"--tool", "Bash",
		"--input", "git status --short",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{`"event": "PreToolUse"`, `"tool": "Bash"`, `"verdict": "allow"`} {
		if !strings.Contains(output, want) {
			t.Fatalf("hook output missing %q: %s", want, output)
		}
	}
}

func TestRunLoopPlan(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{
		"loop", "plan",
		"--repo", ".",
		"--goal", "add tests",
		"--provider", "codex",
		"--verify", "go test ./...",
		"--max-iterations", "2",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{`"goal": "add tests"`, `"provider": "codex"`, `"execution_status": "planned_only"`} {
		if !strings.Contains(output, want) {
			t.Fatalf("loop output missing %q: %s", want, output)
		}
	}
}

func TestRunProcessRun(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go executable not on PATH")
	}
	var stdout, stderr bytes.Buffer
	code := Run([]string{
		"process", "run",
		"--repo", ".",
		"--timeout-seconds", "5",
		"--",
		"go", "version",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{`"exit_code": 0`, `"timed_out": false`, "go version"} {
		if !strings.Contains(output, want) {
			t.Fatalf("process output missing %q: %s", want, output)
		}
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

func TestRunMCPCall(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{
		"mcp", "call", "open_ralph_budget_estimate",
		"--param", "provider=codex",
		"--param", "input_tokens=1000",
		"--param", "output_tokens=500",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{`"ok": true`, `"tool": "open_ralph_budget_estimate"`, `"estimate"`} {
		if !strings.Contains(output, want) {
			t.Fatalf("mcp call output missing %q: %s", want, output)
		}
	}
}

func TestRunSessionAnalyzeTranscript(t *testing.T) {
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
	var stdout, stderr bytes.Buffer
	code := Run([]string{"session", "analyze", "--root", root, "--id", "sess-example"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"replay_ready": true`) {
		t.Fatalf("analysis output: %s", stdout.String())
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
