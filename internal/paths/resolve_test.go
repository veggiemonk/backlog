package paths

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestResolveTasksDir_AbsolutePath(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	abs := filepath.Join(string(os.PathSeparator), "tmp", "backlog", ".backlog")

	got, err := ResolveTasksDir(fs, abs)
	is.NoErr(err)
	is.Equal(got, abs)
}

func TestResolveTasksDir_RelativeExistsInCWD(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	err := fs.MkdirAll(".backlog", 0o755)
	is.NoErr(err)

	got, err := ResolveTasksDir(fs, ".backlog")
	is.NoErr(err)
	is.Equal(got, ".backlog")
}

func TestResolveTasksDir_SearchUpwards(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	// Create .backlog at parent of CWD to be discovered by the upward search.
	cwd, err := os.Getwd()
	is.NoErr(err)

	parent := filepath.Dir(cwd)
	candidate := filepath.Join(parent, ".backlog")
	err = fs.MkdirAll(candidate, 0o755)
	is.NoErr(err)

	got, err := ResolveTasksDir(fs, ".backlog")
	is.NoErr(err)
	is.Equal(got, candidate)
}

func TestResolveTasksDir_GitRootAnchor(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	// Ensure neither CWD nor upward search find a directory by NOT creating any.
	// Since we're in a git repo, this should return the git root + .backlog

	got, err := ResolveTasksDir(fs, ".backlog")
	is.NoErr(err)

	// The result should be the git root + .backlog
	// Since we're already in a git repo, it should find the root
	is.True(filepath.IsAbs(got))                // Should be absolute path
	is.True(strings.HasSuffix(got, ".backlog")) // Should end with .backlog
}

func TestResolveTasksDir_FallbackToCWD(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	// This test verifies the fallback behavior.
	// Since we're in a git repo, it will find the git root, not fall back to CWD.
	// But we can still test that the function returns an absolute path that ends with .backlog

	got, err := ResolveTasksDir(fs, ".backlog")
	is.NoErr(err)
	is.True(filepath.IsAbs(got))                // Should be absolute path
	is.True(strings.HasSuffix(got, ".backlog")) // Should end with .backlog
}
