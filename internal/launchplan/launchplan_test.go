package launchplan

import (
	"strings"
	"testing"
)

func TestBuildCodexPlan(t *testing.T) {
	plan, err := Build(Options{
		Provider:       "codex",
		RepoPath:       ".",
		Prompt:         "add tests",
		Model:          "gpt-5",
		PermissionMode: "read-only",
		MaxTurns:       3,
		BudgetUSD:      1.25,
		Agent:          "tester",
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if plan.Command != "codex" || plan.Provider != "codex" {
		t.Fatalf("unexpected provider plan: %+v", plan)
	}
	joined := strings.Join(plan.Args, " ")
	if !strings.Contains(joined, "--sandbox read-only") {
		t.Fatalf("codex args missing sandbox: %v", plan.Args)
	}
	if !strings.Contains(joined, "tester") {
		t.Fatalf("codex prompt missing agent hint: %v", plan.Args)
	}
	if len(plan.Unsupported) != 2 {
		t.Fatalf("expected codex budget/max-turn unsupported notes: %+v", plan.Unsupported)
	}
	if !plan.ReviewRequired || plan.ExecutionStatus != "planned_only" {
		t.Fatalf("plan should be review-only: %+v", plan)
	}
}

func TestBuildClaudePlan(t *testing.T) {
	plan, err := Build(Options{
		Provider:       "claude",
		RepoPath:       ".",
		Prompt:         "summarize",
		PermissionMode: "workspace-write",
		MaxTurns:       2,
		BudgetUSD:      0.75,
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	joined := strings.Join(plan.Args, " ")
	for _, want := range []string{"--max-budget-usd 0.75", "--max-turns 2", "--permission-mode workspace-write"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("claude args missing %q: %v", want, plan.Args)
		}
	}
}

func TestBuildRejectsHighRiskPermissionMode(t *testing.T) {
	_, err := Build(Options{
		Provider:       "codex",
		RepoPath:       ".",
		Prompt:         "do work",
		PermissionMode: "full-access",
	})
	if err == nil {
		t.Fatal("expected high-risk permission mode error")
	}
}
