package processrun

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestRunCapturesGoVersion(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go executable not on PATH")
	}
	result, err := Run(context.Background(), Options{
		RepoPath: ".",
		Command:  []string{"go", "version"},
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.ExitCode != 0 || !strings.Contains(result.Stdout, "go version") {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestRunReportsNonZeroExit(t *testing.T) {
	result, err := Run(context.Background(), Options{
		RepoPath: ".",
		Command:  []string{"go", "version", "--badflag"},
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.ExitCode == 0 {
		t.Fatalf("expected non-zero exit: %+v", result)
	}
}

func TestRunRequiresCommand(t *testing.T) {
	_, err := Run(context.Background(), Options{RepoPath: "."})
	if err == nil || !strings.Contains(err.Error(), "command is required") {
		t.Fatalf("err=%v", err)
	}
}

func TestRunCapsOutput(t *testing.T) {
	result, err := Run(context.Background(), Options{
		RepoPath:    ".",
		Command:     []string{"go", "env"},
		Timeout:     5 * time.Second,
		OutputLimit: 8,
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !result.StdoutTrimmed || len(result.Stdout) > 8 {
		t.Fatalf("expected capped stdout: %+v", result)
	}
}
