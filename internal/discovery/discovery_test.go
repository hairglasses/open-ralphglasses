package discovery

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestScanFindsGitReposAndOptInConfig(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "example")
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir git dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, ConfigFile), []byte("provider=codex\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	repos, err := Scan(context.Background(), root, 2)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("repos=%+v", repos)
	}
	if repos[0].Name != "example" || !repos[0].Enabled || repos[0].ConfigPath == "" {
		t.Fatalf("unexpected repo: %+v", repos[0])
	}
}
