// Package hookgate evaluates hook decisions for proposed tool use.
package hookgate

import (
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

// Event is a provider hook lifecycle event.
type Event string

const (
	EventPreToolUse        Event = "PreToolUse"
	EventPermissionRequest Event = "PermissionRequest"
	EventPostToolUse       Event = "PostToolUse"
	EventNotification      Event = "Notification"
	EventSessionStart      Event = "SessionStart"
	EventUserPromptSubmit  Event = "UserPromptSubmit"
	EventStop              Event = "Stop"
	EventSubagentStart     Event = "SubagentStart"
	EventSubagentStop      Event = "SubagentStop"
	EventPreCompact        Event = "PreCompact"
)

// Verdict is the coarse gate result.
type Verdict string

const (
	VerdictAllow Verdict = "allow"
	VerdictWarn  Verdict = "warn"
	VerdictBlock Verdict = "block"
)

// CheckInput describes one proposed provider action.
type CheckInput struct {
	Event    string
	Tool     string
	Command  string
	Path     string
	RepoPath string
}

// Decision is the JSON-friendly public gate result.
type Decision struct {
	Event           string      `json:"event"`
	Tool            string      `json:"tool"`
	Verdict         Verdict     `json:"verdict"`
	ReviewRequired  bool        `json:"review_required"`
	Reasons         []string    `json:"reasons,omitempty"`
	Recommendations []string    `json:"recommendations,omitempty"`
	MatchedRules    []RuleMatch `json:"matched_rules,omitempty"`
}

// RuleMatch identifies one rule that contributed to a verdict.
type RuleMatch struct {
	ID       string  `json:"id"`
	Verdict  Verdict `json:"verdict"`
	Fragment string  `json:"fragment,omitempty"`
}

var pipeToShellPattern = regexp.MustCompile(`\b(curl|wget)\b.*\|\s*(sh|bash)\b`)
var removeRootPattern = regexp.MustCompile(`\brm\s+-[a-z]*r[a-z]*f[a-z]*\s+/(?:\s|$)`)

// AllEvents returns the public event set.
func AllEvents() []Event {
	return []Event{
		EventPreToolUse,
		EventPermissionRequest,
		EventPostToolUse,
		EventNotification,
		EventSessionStart,
		EventUserPromptSubmit,
		EventStop,
		EventSubagentStart,
		EventSubagentStop,
		EventPreCompact,
	}
}

// NormalizeEvent maps provider names and short aliases to canonical events.
func NormalizeEvent(value string) (Event, error) {
	switch strings.TrimSpace(value) {
	case string(EventPreToolUse), "pre-tool", "before-tool", "BeforeTool":
		return EventPreToolUse, nil
	case string(EventPermissionRequest), "permission", "permission-request", "approval-request":
		return EventPermissionRequest, nil
	case string(EventPostToolUse), "post-tool", "after-tool", "AfterTool":
		return EventPostToolUse, nil
	case string(EventNotification), "notification":
		return EventNotification, nil
	case string(EventSessionStart), "session-start":
		return EventSessionStart, nil
	case string(EventUserPromptSubmit), "prompt-submit", "user-prompt-submit":
		return EventUserPromptSubmit, nil
	case string(EventStop), "stop":
		return EventStop, nil
	case string(EventSubagentStart), "subagent-start":
		return EventSubagentStart, nil
	case string(EventSubagentStop), "subagent-stop":
		return EventSubagentStop, nil
	case string(EventPreCompact), "pre-compact":
		return EventPreCompact, nil
	default:
		return "", fmt.Errorf("unsupported hook event %q", value)
	}
}

// Check evaluates a proposed action with conservative public rules.
func Check(in CheckInput) (Decision, error) {
	event, err := NormalizeEvent(in.Event)
	if err != nil {
		return Decision{}, err
	}
	tool := strings.TrimSpace(in.Tool)
	if tool == "" {
		tool = "unknown"
	}
	decision := Decision{
		Event:   string(event),
		Tool:    tool,
		Verdict: VerdictAllow,
		Recommendations: []string{
			"review the decision before wiring it into an automatic hook runner",
		},
	}

	if !isGateEvent(event) {
		decision.Reasons = append(decision.Reasons, "event is informational in the public checker")
		return decision, nil
	}

	command := strings.TrimSpace(in.Command)
	if command != "" && isShellTool(tool) {
		evaluateCommand(command, &decision)
	}
	if reason := outsideRepoReason(in.Path, in.RepoPath); reason != "" {
		addMatch(&decision, VerdictWarn, "path_outside_repo", strings.TrimSpace(in.Path), reason)
	}
	if decision.Verdict == VerdictAllow {
		decision.Reasons = append(decision.Reasons, "no public hook-gate rules matched")
	}
	decision.ReviewRequired = decision.Verdict != VerdictAllow
	return decision, nil
}

func isGateEvent(event Event) bool {
	return event == EventPreToolUse ||
		event == EventPermissionRequest ||
		event == EventUserPromptSubmit
}

func isShellTool(tool string) bool {
	switch strings.ToLower(strings.TrimSpace(tool)) {
	case "bash", "shell", "sh", "zsh", "exec", "terminal":
		return true
	default:
		return false
	}
}

func evaluateCommand(command string, decision *Decision) {
	lower := strings.ToLower(command)
	if pipeToShellPattern.MatchString(lower) {
		addMatch(decision, VerdictBlock, "pipe_remote_script_to_shell", command, "remote download piped to a shell is blocked")
	}
	if removeRootPattern.MatchString(lower) {
		addMatch(decision, VerdictBlock, "remove_root", "rm -rf /", "recursive deletion of root is blocked")
	}
	for _, rule := range blockRules() {
		if strings.Contains(lower, rule.fragment) {
			addMatch(decision, VerdictBlock, rule.id, rule.fragment, rule.reason)
		}
	}
	for _, rule := range warnRules() {
		if strings.Contains(lower, rule.fragment) {
			addMatch(decision, VerdictWarn, rule.id, rule.fragment, rule.reason)
		}
	}
}

type fragmentRule struct {
	id       string
	fragment string
	reason   string
}

func blockRules() []fragmentRule {
	return []fragmentRule{
		{id: "format_filesystem", fragment: "mkfs", reason: "filesystem formatting is blocked"},
		{id: "raw_disk_write", fragment: "dd if=", reason: "raw disk writes are blocked"},
		{id: "fork_bomb", fragment: ":(){", reason: "shell fork-bomb pattern is blocked"},
		{id: "terraform_destroy", fragment: "terraform destroy", reason: "infrastructure destroy commands are blocked"},
		{id: "sudo_destructive_delete", fragment: "sudo rm", reason: "privileged destructive deletion is blocked"},
		{id: "sudo_shutdown", fragment: "sudo shutdown", reason: "privileged shutdown is blocked"},
		{id: "sudo_reboot", fragment: "sudo reboot", reason: "privileged reboot is blocked"},
	}
}

func warnRules() []fragmentRule {
	return []fragmentRule{
		{id: "recursive_delete", fragment: "rm -rf", reason: "recursive deletion requires review"},
		{id: "recursive_delete_alt", fragment: "rm -fr", reason: "recursive deletion requires review"},
		{id: "git_push", fragment: "git push", reason: "remote git mutation requires review"},
		{id: "git_commit", fragment: "git commit", reason: "git history mutation requires review"},
		{id: "git_clean", fragment: "git clean", reason: "workspace cleanup can delete untracked files"},
		{id: "docker_prune", fragment: "docker system prune", reason: "docker prune can delete local state"},
		{id: "terraform_apply", fragment: "terraform apply", reason: "infrastructure apply requires review"},
		{id: "kubectl_delete", fragment: "kubectl delete", reason: "cluster delete requires review"},
		{id: "chmod_open", fragment: "chmod -r 777", reason: "broad permission changes require review"},
	}
}

func outsideRepoReason(pathValue, repoPath string) string {
	pathValue = strings.TrimSpace(pathValue)
	repoPath = strings.TrimSpace(repoPath)
	if pathValue == "" || repoPath == "" || !filepath.IsAbs(pathValue) {
		return ""
	}
	absRepo, err := filepath.Abs(repoPath)
	if err != nil {
		return ""
	}
	absPath, err := filepath.Abs(pathValue)
	if err != nil {
		return ""
	}
	rel, err := filepath.Rel(absRepo, absPath)
	if err != nil {
		return ""
	}
	if rel == "." || rel == "" {
		return ""
	}
	if strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return "absolute path is outside the target repository"
	}
	return ""
}

func addMatch(decision *Decision, verdict Verdict, id, fragment, reason string) {
	if slices.ContainsFunc(decision.MatchedRules, func(match RuleMatch) bool {
		return match.ID == id
	}) {
		return
	}
	decision.MatchedRules = append(decision.MatchedRules, RuleMatch{
		ID:       id,
		Verdict:  verdict,
		Fragment: fragment,
	})
	decision.Reasons = append(decision.Reasons, reason)
	if verdictRank(verdict) > verdictRank(decision.Verdict) {
		decision.Verdict = verdict
	}
}

func verdictRank(verdict Verdict) int {
	switch verdict {
	case VerdictBlock:
		return 2
	case VerdictWarn:
		return 1
	default:
		return 0
	}
}
