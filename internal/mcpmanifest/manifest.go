// Package mcpmanifest exposes a tiny JSON manifest for MCP-style adapters.
//
// It is not a full MCP server. The goal is to document the stable public command
// surface in a machine-readable shape that future stdio/HTTP adapters can reuse.
package mcpmanifest

// Tool describes one public command/action pair.
type Tool struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ReadOnly    bool     `json:"read_only"`
	Inputs      []string `json:"inputs,omitempty"`
}

// Manifest returns the initial public tool set.
func Manifest() []Tool {
	return []Tool{
		{
			Name:        "open_ralph_doctor",
			Description: "Check provider executables and public state directory readiness",
			ReadOnly:    true,
		},
		{
			Name:        "open_ralph_provider_list",
			Description: "List configured public providers",
			ReadOnly:    true,
		},
		{
			Name:        "open_ralph_budget_estimate",
			Description: "Estimate provider token cost and evaluate optional budget headroom",
			ReadOnly:    true,
			Inputs:      []string{"provider", "input_tokens", "output_tokens", "budget", "spent"},
		},
		{
			Name:        "open_ralph_launch_plan",
			Description: "Build a review-only provider CLI command plan without executing it",
			ReadOnly:    true,
			Inputs:      []string{"provider", "repo", "prompt", "model", "permission_mode", "budget", "max_turns"},
		},
		{
			Name:        "open_ralph_session_plan",
			Description: "Validate and record a planned provider session",
			ReadOnly:    false,
			Inputs:      []string{"provider", "repo", "prompt"},
		},
		{
			Name:        "open_ralph_session_list",
			Description: "List planned sessions from the local JSONL ledger",
			ReadOnly:    true,
		},
		{
			Name:        "open_ralph_repo_scan",
			Description: "Scan a workspace for Git repos and explicit public opt-in markers",
			ReadOnly:    true,
			Inputs:      []string{"root", "depth"},
		},
		{
			Name:        "open_ralph_worktree_path",
			Description: "Compute a deterministic managed worktree path",
			ReadOnly:    true,
			Inputs:      []string{"root", "repo", "label"},
		},
	}
}
