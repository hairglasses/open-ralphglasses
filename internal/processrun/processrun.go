// Package processrun executes explicit commands with conservative public bounds.
//
// This is intentionally smaller than the private process supervisor. It never
// invokes a shell, does not manage credentials, and captures capped output from a
// single foreground command in a caller-selected repository directory.
package processrun

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultTimeout     = 30 * time.Second
	maxTimeout         = 5 * time.Minute
	defaultOutputLimit = 64 * 1024
)

// Options describes one explicit process run.
type Options struct {
	RepoPath    string
	Command     []string
	Timeout     time.Duration
	OutputLimit int
}

// Result captures the bounded process outcome.
type Result struct {
	Command       []string `json:"command"`
	RepoPath      string   `json:"repo_path"`
	ExitCode      int      `json:"exit_code"`
	TimedOut      bool     `json:"timed_out"`
	DurationMS    int64    `json:"duration_ms"`
	Stdout        string   `json:"stdout,omitempty"`
	Stderr        string   `json:"stderr,omitempty"`
	StdoutTrimmed bool     `json:"stdout_trimmed,omitempty"`
	StderrTrimmed bool     `json:"stderr_trimmed,omitempty"`
}

// Run executes Command without a shell and returns a bounded result. A non-zero
// exit code is reported in Result rather than returned as an error.
func Run(ctx context.Context, opts Options) (Result, error) {
	command := normalizeCommand(opts.Command)
	if len(command) == 0 {
		return Result{}, fmt.Errorf("command is required")
	}
	repoPath := strings.TrimSpace(opts.RepoPath)
	if repoPath == "" {
		return Result{}, fmt.Errorf("repo path is required")
	}
	absRepo, err := filepath.Abs(repoPath)
	if err != nil {
		return Result{}, fmt.Errorf("resolve repo path: %w", err)
	}
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}
	if timeout < 0 || timeout > maxTimeout {
		return Result{}, fmt.Errorf("timeout must be between 0 and %s", maxTimeout)
	}
	outputLimit := opts.OutputLimit
	if outputLimit == 0 {
		outputLimit = defaultOutputLimit
	}
	if outputLimit < 0 {
		return Result{}, fmt.Errorf("output limit must be non-negative")
	}

	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(runCtx, command[0], command[1:]...)
	cmd.Dir = absRepo
	var stdout, stderr cappedBuffer
	stdout.limit = outputLimit
	stderr.limit = outputLimit
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err = cmd.Run()
	result := Result{
		Command:       command,
		RepoPath:      absRepo,
		ExitCode:      0,
		TimedOut:      errors.Is(runCtx.Err(), context.DeadlineExceeded),
		DurationMS:    time.Since(start).Milliseconds(),
		Stdout:        stdout.String(),
		Stderr:        stderr.String(),
		StdoutTrimmed: stdout.trimmed,
		StderrTrimmed: stderr.trimmed,
	}
	if err == nil {
		return result, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
		return result, nil
	}
	if result.TimedOut {
		result.ExitCode = -1
		return result, nil
	}
	return result, fmt.Errorf("start process: %w", err)
}

func normalizeCommand(command []string) []string {
	out := make([]string, 0, len(command))
	for _, part := range command {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

type cappedBuffer struct {
	limit   int
	data    []byte
	trimmed bool
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	if b.limit == 0 {
		b.trimmed = b.trimmed || len(p) > 0
		return len(p), nil
	}
	remaining := b.limit - len(b.data)
	if remaining <= 0 {
		b.trimmed = b.trimmed || len(p) > 0
		return len(p), nil
	}
	if len(p) > remaining {
		b.data = append(b.data, p[:remaining]...)
		b.trimmed = true
		return len(p), nil
	}
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *cappedBuffer) String() string {
	return string(b.data)
}
