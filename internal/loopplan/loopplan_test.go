package loopplan

import (
	"strings"
	"testing"
)

func TestBuildDefaultLoopPlan(t *testing.T) {
	plan, err := Build(Options{
		RepoPath:      ".",
		Goal:          "add focused tests",
		Provider:      "codex",
		VerifyCommand: "go test ./...",
	})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if plan.Provider != "codex" || plan.MaxIterations != defaultMaxIterations {
		t.Fatalf("unexpected plan defaults: %+v", plan)
	}
	if plan.ExecutionStatus != "planned_only" || !plan.ReviewRequired {
		t.Fatalf("unexpected execution contract: %+v", plan)
	}
	if len(plan.VerifyCommands) != 1 || plan.VerifyCommands[0] != "go test ./..." {
		t.Fatalf("verify commands=%+v", plan.VerifyCommands)
	}
}

func TestBuildRejectsMissingGoal(t *testing.T) {
	_, err := Build(Options{RepoPath: "."})
	if err == nil || !strings.Contains(err.Error(), "goal is required") {
		t.Fatalf("err=%v", err)
	}
}

func TestBuildRejectsUnknownProvider(t *testing.T) {
	_, err := Build(Options{
		RepoPath: ".",
		Goal:     "inspect",
		Provider: "unknown",
	})
	if err == nil || !strings.Contains(err.Error(), "unknown provider") {
		t.Fatalf("err=%v", err)
	}
}

func TestBuildCapsIterations(t *testing.T) {
	_, err := Build(Options{
		RepoPath:      ".",
		Goal:          "inspect",
		MaxIterations: maxIterationsLimit + 1,
	})
	if err == nil || !strings.Contains(err.Error(), "max iterations") {
		t.Fatalf("err=%v", err)
	}
}
