// Package loopplan builds review-only plans for iterative agent work.
//
// The private system can run perpetual loops, observe regressions, and relaunch
// providers. This public package keeps a smaller contract: describe the loop,
// verification gate, and stop conditions in JSON before any runner exists.
package loopplan

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hairglasses/open-ralphglasses/internal/provider"
)

const (
	defaultMaxIterations = 3
	maxIterationsLimit   = 20
)

// Options describes the review-only loop intent.
type Options struct {
	RepoPath      string
	Goal          string
	Provider      string
	VerifyCommand string
	MaxIterations int
}

// Plan is an inspectable loop plan. It is never executed by this package.
type Plan struct {
	Goal            string   `json:"goal"`
	RepoPath        string   `json:"repo_path"`
	Provider        string   `json:"provider"`
	MaxIterations   int      `json:"max_iterations"`
	VerifyCommands  []string `json:"verify_commands,omitempty"`
	ReviewRequired  bool     `json:"review_required"`
	ExecutionStatus string   `json:"execution_status"`
	SafetyNotes     []string `json:"safety_notes"`
	Steps           []Step   `json:"steps"`
	StopConditions  []string `json:"stop_conditions"`
}

// Step is one phase in the planned loop.
type Step struct {
	Order       int    `json:"order"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Build validates loop intent and returns a deterministic plan.
func Build(opts Options) (Plan, error) {
	repoPath := strings.TrimSpace(opts.RepoPath)
	if repoPath == "" {
		return Plan{}, fmt.Errorf("repo path is required")
	}
	absRepo, err := filepath.Abs(repoPath)
	if err != nil {
		return Plan{}, fmt.Errorf("resolve repo path: %w", err)
	}
	goal := strings.TrimSpace(opts.Goal)
	if goal == "" {
		return Plan{}, fmt.Errorf("goal is required")
	}
	if opts.MaxIterations < 0 {
		return Plan{}, fmt.Errorf("max iterations must be non-negative")
	}
	maxIterations := opts.MaxIterations
	if maxIterations == 0 {
		maxIterations = defaultMaxIterations
	}
	if maxIterations > maxIterationsLimit {
		return Plan{}, fmt.Errorf("max iterations must be %d or less", maxIterationsLimit)
	}
	providerID := strings.TrimSpace(opts.Provider)
	if providerID == "" {
		providerID = "provider-neutral"
	} else {
		p, err := provider.Lookup(providerID)
		if err != nil {
			return Plan{}, err
		}
		providerID = p.ID
	}
	verifyCommands := normalizeVerifyCommands(opts.VerifyCommand)
	safetyNotes := []string{
		"loop plan is not executed by open-ralphglasses",
		"review repo state, permission mode, and verification gates before adding a runner",
	}
	if len(verifyCommands) == 0 {
		safetyNotes = append(safetyNotes, "no automated verification command supplied")
	}

	return Plan{
		Goal:            goal,
		RepoPath:        absRepo,
		Provider:        providerID,
		MaxIterations:   maxIterations,
		VerifyCommands:  verifyCommands,
		ReviewRequired:  true,
		ExecutionStatus: "planned_only",
		SafetyNotes:     safetyNotes,
		Steps:           buildSteps(len(verifyCommands) > 0),
		StopConditions: []string{
			"goal is satisfied",
			"verification fails and requires human review",
			"max_iterations is reached",
		},
	}, nil
}

func normalizeVerifyCommands(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, "\n")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func buildSteps(hasVerify bool) []Step {
	verifyDescription := "perform manual verification and capture evidence"
	if hasVerify {
		verifyDescription = "run the configured verification command and capture output"
	}
	return []Step{
		{Order: 1, Name: "inspect", Description: "read repo state, task context, and current diffs"},
		{Order: 2, Name: "plan_slice", Description: "choose one bounded implementation slice"},
		{Order: 3, Name: "implement", Description: "apply the smallest changes needed for the slice"},
		{Order: 4, Name: "verify", Description: verifyDescription},
		{Order: 5, Name: "record", Description: "record result, residual risk, and next action"},
		{Order: 6, Name: "decide", Description: "stop on completion or start the next bounded iteration"},
	}
}
