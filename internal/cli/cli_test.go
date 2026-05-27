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

func TestRunMCPManifest(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"mcp", "manifest"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "open_ralph_launch_plan") {
		t.Fatalf("manifest output missing launch tool: %s", stdout.String())
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
