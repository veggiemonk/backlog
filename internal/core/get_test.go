package core_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

func TestGetTask(t *testing.T) {
	is := is.New(t)
	store := core.NewFileTaskStore(afero.NewMemMapFs(), ".backlog", core.NewMockLocker())
	_, _ = store.Create(core.CreateTaskParams{Title: "View Me"})

	task, err := store.Get("1")
	is.NoErr(err)
	is.Equal("View Me", task.Title)
	is.Equal("T01", task.ID.Name())
}
