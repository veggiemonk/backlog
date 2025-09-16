package paths

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func TestResolveTasksDir_AbsolutePath(t *testing.T) {
	fs := afero.NewMemMapFs()
	abs := filepath.Join(string(os.PathSeparator), "tmp", "backlog", ".backlog")

	got, err := ResolveTasksDir(fs, abs)
	if err != nil {
		t.Fatalf("ResolveTasksDir returned error: %v", err)
	}
	if got != abs {
		t.Fatalf("expected absolute path returned as-is, got %q, want %q", got, abs)
	}
}

func TestResolveTasksDir_RelativeExistsInCWD(t *testing.T) {
	fs := afero.NewMemMapFs()
	if err := fs.MkdirAll(".backlog", 0o755); err != nil {
		t.Fatalf("setup mkdir .backlog: %v", err)
	}
	got, err := ResolveTasksDir(fs, ".backlog")
	if err != nil {
		t.Fatalf("ResolveTasksDir returned error: %v", err)
	}
	if got != ".backlog" {
		t.Fatalf("expected '.backlog', got %q", got)
	}
}

func TestResolveTasksDir_SearchUpwards(t *testing.T) {
	fs := afero.NewMemMapFs()
	// Create .backlog at parent of CWD to be discovered by the upward search.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	parent := filepath.Dir(cwd)
	candidate := filepath.Join(parent, ".backlog")
	if err := fs.MkdirAll(candidate, 0o755); err != nil {
		t.Fatalf("setup mkdir %s: %v", candidate, err)
	}

	got, err := ResolveTasksDir(fs, ".backlog")
	if err != nil {
		t.Fatalf("ResolveTasksDir returned error: %v", err)
	}
	if got != candidate {
		t.Fatalf("expected upward-found path %q, got %q", candidate, got)
	}
}

func TestResolveTasksDir_GitRootAnchor(t *testing.T) {
	fs := afero.NewMemMapFs()
	// Ensure neither CWD nor upward search find a directory by NOT creating any.
	// Mock the git root finder to return a fake root.
	fakeRoot := filepath.Join(string(os.PathSeparator), "fake", "gitroot")

	got, err := ResolveTasksDir(fs, ".backlog")
	if err != nil {
		t.Fatalf("ResolveTasksDir returned error: %v", err)
	}
	want := filepath.Join(fakeRoot, ".backlog")
	if got != want {
		t.Fatalf("expected git-root anchored path %q, got %q", want, got)
	}
}

func TestResolveTasksDir_FallbackToCWD(t *testing.T) {
	fs := afero.NewMemMapFs()
	// Neither CWD nor upward search find .backlog; mock git root finder to fail.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	absCWD, err := filepath.Abs(cwd)
	if err != nil {
		t.Fatalf("abs cwd: %v", err)
	}
	want := filepath.Join(absCWD, ".backlog")

	got, err := ResolveTasksDir(fs, ".backlog")
	if err != nil {
		t.Fatalf("ResolveTasksDir returned error: %v", err)
	}
	if got != want {
		t.Fatalf("expected fallback path %q, got %q", want, got)
	}
}

// TODO: write test for env var BACKLOG_FOLDER
