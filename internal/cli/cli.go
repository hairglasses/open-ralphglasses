// Package cli wires the small public control-plane primitives into a CLI.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/hairglasses/open-ralphglasses/internal/budget"
	"github.com/hairglasses/open-ralphglasses/internal/discovery"
	"github.com/hairglasses/open-ralphglasses/internal/mcpmanifest"
	"github.com/hairglasses/open-ralphglasses/internal/provider"
	"github.com/hairglasses/open-ralphglasses/internal/session"
	"github.com/hairglasses/open-ralphglasses/internal/worktree"
)

const version = "0.1.0"

// Run executes the CLI and returns a process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printHelp(stdout)
		return 0
	}
	switch args[0] {
	case "help", "-h", "--help":
		printHelp(stdout)
		return 0
	case "version":
		fmt.Fprintf(stdout, "open-ralphglasses %s\n", version)
		return 0
	case "budget":
		return runBudget(args[1:], stdout, stderr)
	case "doctor":
		return runDoctor(stdout)
	case "providers":
		return runProviders(stdout)
	case "session":
		return runSession(args[1:], stdout, stderr)
	case "mcp":
		return runMCP(args[1:], stdout, stderr)
	case "repos":
		return runRepos(args[1:], stdout, stderr)
	case "worktree":
		return runWorktree(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		printHelp(stderr)
		return 2
	}
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, `open-ralphglasses - public multi-provider agent control-plane seed

Usage:
  open-ralphglasses doctor
  open-ralphglasses providers
  open-ralphglasses budget estimate --provider codex --input-tokens 1000 --output-tokens 500
  open-ralphglasses session start --provider codex --repo . --prompt "Summarize this repo"
  open-ralphglasses session list
  open-ralphglasses mcp manifest
  open-ralphglasses repos scan --root . --depth 3
  open-ralphglasses worktree path --repo example --label refactor-api`)
}

func runDoctor(stdout io.Writer) int {
	fmt.Fprintln(stdout, "provider\tinstalled\tcommand")
	for _, p := range provider.Catalog() {
		fmt.Fprintf(stdout, "%s\t%t\t%s\n", p.ID, p.Installed(), p.Command)
	}
	return 0
}

func runProviders(stdout io.Writer) int {
	tw := tabwriter.NewWriter(stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tDISPLAY\tCOMMAND\tDEFAULT MODEL\tNOTES")
	for _, p := range provider.Catalog() {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", p.ID, p.DisplayName, p.Command, p.DefaultModel, p.Notes)
	}
	_ = tw.Flush()
	return 0
}

func runBudget(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "estimate" {
		fmt.Fprintln(stderr, "usage: open-ralphglasses budget estimate --provider codex --input-tokens 1000 --output-tokens 500 [--budget 1.00 --spent 0.25]")
		return 2
	}
	flags := parseFlags(args[1:])
	inputTokens, ok := parseOptionalIntFlag(flags, "input-tokens", stderr)
	if !ok {
		return 2
	}
	outputTokens, ok := parseOptionalIntFlag(flags, "output-tokens", stderr)
	if !ok {
		return 2
	}
	estimate, err := budget.Estimate(budget.EstimateInput{
		Provider:     flags["provider"],
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
	})
	if err != nil {
		fmt.Fprintf(stderr, "estimate budget: %v\n", err)
		return 2
	}
	result := map[string]any{"estimate": estimate}
	if flags["budget"] != "" || flags["spent"] != "" {
		budgetUSD, ok := parseOptionalFloatFlag(flags, "budget", stderr)
		if !ok {
			return 2
		}
		spentUSD, ok := parseOptionalFloatFlag(flags, "spent", stderr)
		if !ok {
			return 2
		}
		headroomPct, ok := parseOptionalFloatFlag(flags, "headroom-pct", stderr)
		if !ok {
			return 2
		}
		status, err := budget.Status(budget.StatusInput{
			SpentUSD:    spentUSD,
			BudgetUSD:   budgetUSD,
			HeadroomPct: headroomPct,
		})
		if err != nil {
			fmt.Fprintf(stderr, "estimate budget: %v\n", err)
			return 2
		}
		result["status"] = status
	}
	encoded, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintln(stdout, string(encoded))
	return 0
}

func runSession(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "session subcommand is required: start or list")
		return 2
	}
	switch args[0] {
	case "start":
		return runSessionStart(args[1:], stdout, stderr)
	case "list":
		store := session.Store{Root: "."}
		sessions, err := store.List()
		if err != nil {
			fmt.Fprintf(stderr, "list sessions: %v\n", err)
			return 1
		}
		encoded, _ := json.MarshalIndent(sessions, "", "  ")
		fmt.Fprintln(stdout, string(encoded))
		return 0
	default:
		fmt.Fprintf(stderr, "unknown session subcommand %q\n", args[0])
		return 2
	}
}

func runSessionStart(args []string, stdout, stderr io.Writer) int {
	flags := parseFlags(args)
	sess, err := session.New(session.StartOptions{
		Provider: flags["provider"],
		RepoPath: flags["repo"],
		Prompt:   flags["prompt"],
	})
	if err != nil {
		fmt.Fprintf(stderr, "plan session: %v\n", err)
		return 2
	}
	if err := (session.Store{Root: "."}).Append(sess); err != nil {
		fmt.Fprintf(stderr, "write session: %v\n", err)
		return 1
	}
	encoded, _ := json.MarshalIndent(sess, "", "  ")
	fmt.Fprintln(stdout, string(encoded))
	return 0
}

func runMCP(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 || args[0] != "manifest" {
		fmt.Fprintln(stderr, "usage: open-ralphglasses mcp manifest")
		return 2
	}
	encoded, _ := json.MarshalIndent(mcpmanifest.Manifest(), "", "  ")
	fmt.Fprintln(stdout, string(encoded))
	return 0
}

func runRepos(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "scan" {
		fmt.Fprintln(stderr, "usage: open-ralphglasses repos scan --root . --depth 3")
		return 2
	}
	flags := parseFlags(args[1:])
	root := flags["root"]
	if root == "" {
		root = "."
	}
	depth := 3
	if flags["depth"] != "" {
		parsed, err := strconv.Atoi(flags["depth"])
		if err != nil || parsed < 1 {
			fmt.Fprintf(stderr, "invalid depth %q\n", flags["depth"])
			return 2
		}
		depth = parsed
	}
	repos, err := discovery.Scan(context.Background(), root, depth)
	if err != nil {
		fmt.Fprintf(stderr, "scan repos: %v\n", err)
		return 1
	}
	encoded, _ := json.MarshalIndent(repos, "", "  ")
	fmt.Fprintln(stdout, string(encoded))
	return 0
}

func runWorktree(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "path" {
		fmt.Fprintln(stderr, "usage: open-ralphglasses worktree path --repo example --label task")
		return 2
	}
	flags := parseFlags(args[1:])
	root := flags["root"]
	if root == "" {
		root = filepath.Join(".", ".open-ralph", "worktrees")
	}
	repoName := flags["repo"]
	if repoName == "" {
		repoName = "repo"
	}
	label := flags["label"]
	if label == "" {
		label = "work"
	}
	fmt.Fprintln(stdout, worktree.ManagedPath(root, repoName, label))
	return 0
}

func parseFlags(args []string) map[string]string {
	flags := make(map[string]string)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		key := strings.TrimPrefix(arg, "--")
		if key == "" {
			continue
		}
		if strings.Contains(key, "=") {
			parts := strings.SplitN(key, "=", 2)
			flags[parts[0]] = parts[1]
			continue
		}
		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
			flags[key] = args[i+1]
			i++
			continue
		}
		flags[key] = "true"
	}
	if flags["repo"] == "" {
		if cwd, err := os.Getwd(); err == nil {
			flags["repo"] = cwd
		}
	}
	return flags
}

func parseOptionalIntFlag(flags map[string]string, key string, stderr io.Writer) (int, bool) {
	raw := strings.TrimSpace(flags[key])
	if raw == "" {
		return 0, true
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		fmt.Fprintf(stderr, "invalid %s %q\n", key, raw)
		return 0, false
	}
	return value, true
}

func parseOptionalFloatFlag(flags map[string]string, key string, stderr io.Writer) (float64, bool) {
	raw := strings.TrimSpace(flags[key])
	if raw == "" {
		return 0, true
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil || value < 0 {
		fmt.Fprintf(stderr, "invalid %s %q\n", key, raw)
		return 0, false
	}
	return value, true
}
