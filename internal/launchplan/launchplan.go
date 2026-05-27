// Package launchplan builds public, reviewable provider command plans.
//
// Private ralphglasses can launch, supervise, meter, and recover real provider
// processes. This public package deliberately stops one step earlier: it turns a
// provider-neutral request into a command shape that humans and tests can
// inspect before any process execution layer exists. That keeps the useful
// normalization pattern while avoiding secrets, local launcher wrappers,
// account policy, browser state, and machine-specific process management.
package launchplan

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hairglasses/open-ralphglasses/internal/provider"
)

// Options describes the portable launch intent shared by provider CLIs.
//
// Empty optional fields are ignored. Prompt and RepoPath are required because a
// useful command plan should always have both a target checkout and a task.
type Options struct {
	Provider       string
	RepoPath       string
	Prompt         string
	Model          string
	PermissionMode string
	MaxTurns       int
	BudgetUSD      float64
	Agent          string
}

// Plan is the public command plan. It is not executed by this package.
type Plan struct {
	Provider        string   `json:"provider"`
	Command         string   `json:"command"`
	Args            []string `json:"args"`
	RepoPath        string   `json:"repo_path"`
	Model           string   `json:"model,omitempty"`
	PermissionMode  string   `json:"permission_mode,omitempty"`
	EnvPolicy       string   `json:"env_policy"`
	ReviewRequired  bool     `json:"review_required"`
	SafetyNotes     []string `json:"safety_notes"`
	Unsupported     []string `json:"unsupported,omitempty"`
	PromptInArgs    bool     `json:"prompt_in_args"`
	ExecutionStatus string   `json:"execution_status"`
}

// Build returns a provider-specific command plan. The plan is intentionally
// conservative: high-risk permission modes are rejected, environment handling
// is represented as a policy string instead of real variables, and unsupported
// provider features are reported rather than silently dropped.
func Build(opts Options) (Plan, error) {
	p, err := provider.Lookup(opts.Provider)
	if err != nil {
		return Plan{}, err
	}
	repoPath := strings.TrimSpace(opts.RepoPath)
	if repoPath == "" {
		return Plan{}, fmt.Errorf("repo path is required")
	}
	absRepo, err := filepath.Abs(repoPath)
	if err != nil {
		return Plan{}, fmt.Errorf("resolve repo path: %w", err)
	}
	prompt := strings.TrimSpace(opts.Prompt)
	if prompt == "" {
		return Plan{}, fmt.Errorf("prompt is required")
	}
	permission, err := normalizePermissionMode(opts.PermissionMode)
	if err != nil {
		return Plan{}, err
	}
	if opts.MaxTurns < 0 || opts.BudgetUSD < 0 {
		return Plan{}, fmt.Errorf("max turns and budget must be non-negative")
	}
	model := strings.TrimSpace(opts.Model)
	if model == "" {
		model = p.DefaultModel
	}

	plan := Plan{
		Provider:        p.ID,
		Command:         p.Command,
		RepoPath:        absRepo,
		Model:           model,
		PermissionMode:  permission,
		EnvPolicy:       "inherit shell environment; caller is responsible for credential injection",
		ReviewRequired:  true,
		ExecutionStatus: "planned_only",
		SafetyNotes: []string{
			"command is not executed by open-ralphglasses",
			"review args and environment policy before wiring a process runner",
		},
	}

	switch p.ID {
	case "claude":
		plan.Args, plan.PromptInArgs, plan.Unsupported = claudeArgs(opts, model, permission, prompt)
	case "codex":
		plan.Args, plan.PromptInArgs, plan.Unsupported = codexArgs(opts, model, permission, prompt)
	case "gemini":
		plan.Args, plan.PromptInArgs, plan.Unsupported = geminiArgs(opts, model, permission, prompt)
	default:
		return Plan{}, fmt.Errorf("provider %q cannot be planned", p.ID)
	}
	return plan, nil
}

func claudeArgs(opts Options, model, permission, prompt string) ([]string, bool, []string) {
	args := []string{"-p", "--output-format", "stream-json", "--model", model}
	var unsupported []string
	if opts.BudgetUSD > 0 {
		args = append(args, "--max-budget-usd", formatFloat(opts.BudgetUSD))
	}
	if opts.MaxTurns > 0 {
		args = append(args, "--max-turns", strconv.Itoa(opts.MaxTurns))
	}
	if opts.Agent != "" {
		args = append(args, "--agent", strings.TrimSpace(opts.Agent))
	}
	if permission != "" {
		args = append(args, "--permission-mode", permission)
	}
	args = append(args, prompt)
	return args, true, unsupported
}

func codexArgs(opts Options, model, permission, prompt string) ([]string, bool, []string) {
	args := []string{"exec", "--model", model, "--json"}
	var unsupported []string
	if opts.BudgetUSD > 0 {
		unsupported = append(unsupported, "budget_usd")
	}
	if opts.MaxTurns > 0 {
		unsupported = append(unsupported, "max_turns")
	}
	if opts.Agent != "" {
		prompt = fmt.Sprintf("Use the project agent %q if available.\n\nTask:\n%s", strings.TrimSpace(opts.Agent), prompt)
	}
	if sandboxMode := codexSandboxMode(permission); sandboxMode != "" {
		args = append(args, "--sandbox", sandboxMode)
	}
	args = append(args, prompt)
	return args, true, unsupported
}

func geminiArgs(opts Options, model, permission, prompt string) ([]string, bool, []string) {
	args := []string{"--output-format", "stream-json", "--model", model}
	var unsupported []string
	if opts.BudgetUSD > 0 {
		unsupported = append(unsupported, "budget_usd")
	}
	if opts.MaxTurns > 0 {
		unsupported = append(unsupported, "max_turns")
	}
	if opts.Agent != "" {
		args = append(args, strings.TrimSpace(opts.Agent))
	}
	args = append(args, "--approval-mode", geminiApprovalMode(permission), "-p", prompt)
	return args, true, unsupported
}

func normalizePermissionMode(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "", "default", "auto":
		return "", nil
	case "plan", "read", "read-only", "readonly":
		return "read-only", nil
	case "edit", "workspace", "workspace-write", "workspace_write":
		return "workspace-write", nil
	}
	if strings.Contains(value, "bypass") || strings.Contains(value, "danger") || strings.Contains(value, "full") {
		return "", fmt.Errorf("permission mode %q is intentionally not supported by the public planner", value)
	}
	return "", fmt.Errorf("unsupported permission mode %q", value)
}

func codexSandboxMode(permission string) string {
	switch permission {
	case "read-only":
		return "read-only"
	case "workspace-write":
		return "workspace-write"
	default:
		return ""
	}
}

func geminiApprovalMode(permission string) string {
	switch permission {
	case "read-only":
		return "default"
	case "workspace-write":
		return "auto_edit"
	default:
		return "default"
	}
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
