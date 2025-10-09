package cmd

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/veggiemonk/backlog/internal/core"
)

func Test_runSearch(t *testing.T) {
	t.Run("basic search by title", func(t *testing.T) {
		is := is.New(t)
		output, err := exec(t, "list", "--query", "First", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")
	})

	t.Run("search by partial title", func(t *testing.T) {
		is := is.New(t)
		output, err := exec(t, "list", "--query", "Task", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), countTask) // All tasks contain "Task" in title
	})

	t.Run("search by description", func(t *testing.T) {
		is := is.New(t)
		output, err := exec(t, "list", "--query", "First description", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")
	})
}
