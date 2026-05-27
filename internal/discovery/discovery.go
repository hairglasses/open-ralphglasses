// Package discovery scans a workspace for repos that have opted into
// open-ralphglasses.
//
// The private system scans several local state formats. This public package uses
// only explicit, example-friendly markers: a Git checkout plus an optional
// `.open-ralphrc` file. That keeps discovery useful without depending on private
// workspace layout or hidden operator state.
package discovery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const ConfigFile = ".open-ralphrc"

// Repo is the public read model for one discovered repository.
type Repo struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	ConfigPath string `json:"config_path,omitempty"`
	Enabled    bool   `json:"enabled"`
}

// Scan walks root up to maxDepth directory levels and returns Git repositories.
// If a repo contains `.open-ralphrc`, Enabled is true. Plain Git repos are still
// returned so users can see what might be onboarded next.
func Scan(ctx context.Context, root string, maxDepth int) ([]Repo, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		root = "."
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve scan root: %w", err)
	}
	if maxDepth <= 0 {
		maxDepth = 3
	}

	var repos []Repo
	err = filepath.WalkDir(absRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if !entry.IsDir() {
			return nil
		}
		if shouldSkipDir(entry.Name()) && path != absRoot {
			return filepath.SkipDir
		}
		depth, err := relativeDepth(absRoot, path)
		if err != nil {
			return err
		}
		if depth > maxDepth {
			return filepath.SkipDir
		}
		gitDir := filepath.Join(path, ".git")
		if _, err := os.Stat(gitDir); err != nil {
			return nil
		}
		configPath := filepath.Join(path, ConfigFile)
		enabled := false
		if _, err := os.Stat(configPath); err == nil {
			enabled = true
		} else {
			configPath = ""
		}
		repos = append(repos, Repo{
			Name:       filepath.Base(path),
			Path:       path,
			ConfigPath: configPath,
			Enabled:    enabled,
		})
		return filepath.SkipDir
	})
	if err != nil {
		return nil, err
	}
	slices.SortFunc(repos, func(a, b Repo) int {
		return strings.Compare(a.Path, b.Path)
	})
	return repos, nil
}

func shouldSkipDir(name string) bool {
	switch name {
	case ".git", "node_modules", "vendor", ".cache", "dist", "bin":
		return true
	default:
		return false
	}
}

func relativeDepth(root, path string) (int, error) {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return 0, err
	}
	if rel == "." {
		return 0, nil
	}
	return len(strings.Split(rel, string(os.PathSeparator))), nil
}
