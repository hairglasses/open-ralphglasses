// Package worktree contains public-safe helpers for managed worktree paths.
//
// Private ralphglasses has provider launchers and local workspace policies. This
// package keeps only deterministic path planning. It never creates or deletes a
// worktree by itself.
package worktree

import (
	"path/filepath"
	"regexp"
	"strings"
)

var unsafeName = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

// SanitizeName converts a branch/task label into a filesystem-safe leaf.
func SanitizeName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "work"
	}
	value = unsafeName.ReplaceAllString(value, "-")
	value = strings.Trim(value, ".-_")
	if value == "" {
		return "work"
	}
	return value
}

// ManagedPath returns where a managed worktree would live under root.
func ManagedPath(root, repoName, label string) string {
	root = strings.TrimSpace(root)
	if root == "" {
		root = ".open-ralph/worktrees"
	}
	return filepath.Join(root, SanitizeName(repoName), SanitizeName(label))
}
