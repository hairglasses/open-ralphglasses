// Package session implements a tiny durable session ledger.
//
// Private ralphglasses launches and supervises real provider processes. This
// public subset keeps the safer foundation: validate launch intent, normalize
// provider/repo/prompt metadata, and persist an auditable JSONL record. Downstream
// users can wire real process execution on top without inheriting private local
// automation.
package session

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hairglasses/open-ralphglasses/internal/provider"
)

const stateDirName = ".open-ralph"

// Session is the durable public record of one requested agent run.
type Session struct {
	ID        string    `json:"id"`
	Provider  string    `json:"provider"`
	RepoPath  string    `json:"repo_path"`
	Prompt    string    `json:"prompt"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// StartOptions is intentionally small. Keep provider secrets, browser state,
// tenant ids, and machine-specific paths outside this public contract.
type StartOptions struct {
	Provider string
	RepoPath string
	Prompt   string
	Now      time.Time
}

// New validates launch intent and returns a session record. It does not start a
// child process; callers can make that an explicit, separately reviewed layer.
func New(opts StartOptions) (Session, error) {
	p, err := provider.Lookup(opts.Provider)
	if err != nil {
		return Session{}, err
	}
	repoPath := strings.TrimSpace(opts.RepoPath)
	if repoPath == "" {
		return Session{}, fmt.Errorf("repo path is required")
	}
	absRepo, err := filepath.Abs(repoPath)
	if err != nil {
		return Session{}, fmt.Errorf("resolve repo path: %w", err)
	}
	prompt := strings.TrimSpace(opts.Prompt)
	if prompt == "" {
		return Session{}, fmt.Errorf("prompt is required")
	}
	now := opts.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return Session{
		ID:        newID(now),
		Provider:  p.ID,
		RepoPath:  absRepo,
		Prompt:    prompt,
		Status:    "planned",
		CreatedAt: now.UTC(),
	}, nil
}

// Store appends and reads JSONL records under a workspace-local state
// directory. The file format is deliberately boring so operators can inspect it
// with common shell tools.
type Store struct {
	Root string
}

// Path returns the JSONL path used by this store.
func (s Store) Path() string {
	root := strings.TrimSpace(s.Root)
	if root == "" {
		root = "."
	}
	return filepath.Join(root, stateDirName, "sessions.jsonl")
}

// Append writes one session record with owner-readable permissions.
func (s Store) Append(sess Session) error {
	path := s.Path()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create session state dir: %w", err)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open session ledger: %w", err)
	}
	defer func() { _ = file.Close() }()
	encoded, err := json.Marshal(sess)
	if err != nil {
		return fmt.Errorf("encode session: %w", err)
	}
	if _, err := file.Write(append(encoded, '\n')); err != nil {
		return fmt.Errorf("append session: %w", err)
	}
	return nil
}

// List reads all sessions in append order.
func (s Store) List() ([]Session, error) {
	file, err := os.Open(s.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open session ledger: %w", err)
	}
	defer func() { _ = file.Close() }()
	var sessions []Session
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var sess Session
		if err := json.Unmarshal([]byte(line), &sess); err != nil {
			return nil, fmt.Errorf("decode session ledger: %w", err)
		}
		sessions = append(sessions, sess)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read session ledger: %w", err)
	}
	return sessions, nil
}

func newID(now time.Time) string {
	var buf [4]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return fmt.Sprintf("sess-%s", now.UTC().Format("20060102T150405Z"))
	}
	return fmt.Sprintf("sess-%s-%s", now.UTC().Format("20060102T150405Z"), hex.EncodeToString(buf[:]))
}
