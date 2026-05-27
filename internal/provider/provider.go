// Package provider contains the public, provider-neutral launch catalog.
//
// The private ralphglasses tree has many provider, account, and local runtime
// integrations. This public package keeps only the portable contract: a provider
// has a stable id, a human label, an executable name, and a default model. The
// CLI can print or validate that contract without knowing any private account,
// key-pool, or machine-specific details.
package provider

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

// Provider is the public description of one supported agent runtime.
//
// Command intentionally stores only the executable name. Arguments are assembled
// by the caller so applications can choose their own prompting, sandboxing, and
// environment policy without this package needing secrets or local paths.
type Provider struct {
	ID           string
	DisplayName  string
	Command      string
	DefaultModel string
	Notes        string
}

// Catalog returns the built-in public providers. These are examples of the
// provider-normalization pattern from ralphglasses, reduced to data that is safe
// to publish and easy to adapt.
func Catalog() []Provider {
	return []Provider{
		{
			ID:           "codex",
			DisplayName:  "OpenAI Codex CLI",
			Command:      "codex",
			DefaultModel: "gpt-5",
			Notes:        "coding-oriented agent runtime",
		},
		{
			ID:           "claude",
			DisplayName:  "Claude Code",
			Command:      "claude",
			DefaultModel: "sonnet",
			Notes:        "interactive coding agent runtime",
		},
		{
			ID:           "gemini",
			DisplayName:  "Gemini CLI",
			Command:      "gemini",
			DefaultModel: "gemini-pro",
			Notes:        "general-purpose agent runtime",
		},
	}
}

// IDs returns provider ids in catalog order.
func IDs() []string {
	providers := Catalog()
	ids := make([]string, 0, len(providers))
	for _, p := range providers {
		ids = append(ids, p.ID)
	}
	return ids
}

// Lookup resolves a provider id using case-insensitive matching.
func Lookup(id string) (Provider, error) {
	id = strings.ToLower(strings.TrimSpace(id))
	for _, p := range Catalog() {
		if p.ID == id {
			return p, nil
		}
	}
	return Provider{}, fmt.Errorf("unknown provider %q (want one of: %s)", id, strings.Join(IDs(), ", "))
}

// Installed reports whether the provider executable is available on PATH.
func (p Provider) Installed() bool {
	_, err := exec.LookPath(p.Command)
	return err == nil
}

// ValidateID is useful at API boundaries that should reject unsupported values
// before building commands or writing session records.
func ValidateID(id string) error {
	_, err := Lookup(id)
	return err
}

// IsKnown reports whether id appears in the public catalog.
func IsKnown(id string) bool {
	return slices.Contains(IDs(), strings.ToLower(strings.TrimSpace(id)))
}
