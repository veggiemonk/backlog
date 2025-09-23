package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/veggiemonk/backlog/internal/core"
)

func Test_runSearch(t *testing.T) {
	t.Run("basic search by title", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "First", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")
	})

	t.Run("search by partial title", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), countTask) // All tasks contain "Task" in title
	})

	t.Run("search by description", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "First description", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")
	})

	t.Run("search by label", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "first", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")
	})

	t.Run("search by assigned user", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "first-user", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")
	})

	t.Run("search by priority", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "high", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "High Priority Task")
	})

	t.Run("search by acceptance criteria", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "First AC", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")
	})

	t.Run("search by implementation notes", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "implementation notes", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 5) // First-Fifth tasks have implementation notes
	})

	t.Run("search by implementation plan", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Second implementation plan", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "Second Task")
	})

	t.Run("case insensitive search", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "FIRST", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")
	})

	t.Run("search with no results", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "nonexistent", "-j")
		is.NoErr(err)
		outputStr := strings.TrimSpace(string(output))
		is.Equal(outputStr, "{\"tasks\":[],\"pagination\":{\"total_results\":0,\"displayed_results\":0,\"offset\":0,\"limit\":25,\"has_more\":false}}") // empty ListResult JSON
	})

	t.Run("search with status filter", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "-s", "todo", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), countTask-1) // All tasks except the in-progress one
	})

	t.Run("search with multiple status filter", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "-s", "todo,in-progress", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), countTask) // All tasks
	})

	t.Run("search with assigned filter", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "-a", "first-user", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")
	})

	t.Run("search with multiple assigned filter", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "-a", "first-user,second-user", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 2)
		is.Equal(listResult.Tasks[0].Title, "First Task")
		is.Equal(listResult.Tasks[1].Title, "Second Task")
	})

	t.Run("search with labels filter", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "-l", "first", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")
	})

	t.Run("search with multiple labels filter", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "-l", "first,second", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 2)
		is.Equal(listResult.Tasks[0].Title, "First Task")
		is.Equal(listResult.Tasks[1].Title, "Second Task")
	})

	t.Run("search with unassigned filter", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "-u", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "Unassigned Task")
	})

	t.Run("search with sorting by title", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "--sort", "title", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), countTask)
		is.Equal(listResult.Tasks[0].Title, "Fifth Task") // Alphabetically first
	})

	t.Run("search with sorting by priority", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "--sort", "priority", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), countTask)
		is.Equal(listResult.Tasks[0].Title, "High Priority Task") // Highest priority first
	})

	t.Run("search with reverse order", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "-r", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), countTask)
		is.Equal(listResult.Tasks[0].Title, "In Progress Task") // Last task when reversed
	})

	t.Run("search with combined filters", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "-l", "second", "-s", "todo", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "Second Task")
	})

	t.Run("search markdown output format", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "First Task", "-m")
		is.NoErr(err)
		outputStr := string(output)
		is.True(strings.Contains(outputStr, "|"))    // markdown table contains pipes
		is.True(strings.Contains(outputStr, "ID"))   // should contain header
		is.True(strings.Contains(outputStr, ":---")) // markdown table header separator
		is.True(strings.Contains(outputStr, "Found") && strings.Contains(outputStr, "matching"))
	})

	t.Run("search hide extra fields", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "First Task", "--hide-extra")
		is.NoErr(err)
		outputStr := string(output)
		is.True(strings.Contains(outputStr, "ID"))
		is.True(strings.Contains(outputStr, "STATUS"))
		is.True(strings.Contains(outputStr, "TITLE"))
		// Should not contain extra columns
		is.True(!strings.Contains(outputStr, "LABELS"))
		is.True(!strings.Contains(outputStr, "PRIORITY"))
		is.True(!strings.Contains(outputStr, "ASSIGNED"))
		is.True(strings.Contains(outputStr, "Found") && strings.Contains(outputStr, "matching"))
	})

	t.Run("search default table output format", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "First Task")
		is.NoErr(err)
		outputStr := string(output)
		is.True(strings.Contains(outputStr, "ID"))
		is.True(strings.Contains(outputStr, "STATUS"))
		is.True(strings.Contains(outputStr, "TITLE"))
		is.True(strings.Contains(outputStr, "LABELS"))
		is.True(strings.Contains(outputStr, "PRIORITY"))
		is.True(strings.Contains(outputStr, "ASSIGNED"))
		is.True(strings.Contains(outputStr, "Found") && strings.Contains(outputStr, "matching"))
	})

	t.Run("search JSON output format", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "First Task", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")
		// Should not contain message prefix in JSON output
		outputStr := string(output)
		is.True(!strings.Contains(outputStr, "Found"))
	})

	t.Run("search empty results table format", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "nonexistent")
		is.NoErr(err)
		outputStr := string(output)
		is.True(strings.Contains(outputStr, "No tasks found"))
	})

	t.Run("search empty results JSON format", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "nonexistent", "-j")
		is.NoErr(err)
		outputStr := strings.TrimSpace(string(output))
		is.Equal(outputStr, "{\"tasks\":[],\"pagination\":{\"total_results\":0,\"displayed_results\":0,\"offset\":0,\"limit\":25,\"has_more\":false}}") // empty ListResult JSON
	})

	t.Run("short flag aliases", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// Test short flag for JSON output
		output, err := exec(t, "search", runSearch, "First Task", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")

		// Test short flag for markdown output
		output, err = exec(t, "search", runSearch, "First Task", "-m")
		is.NoErr(err)
		outputStr := string(output)
		is.True(strings.Contains(outputStr, "|:---"))

		// Test short flag for hide extra fields
		output, err = exec(t, "search", runSearch, "First Task", "-e")
		is.NoErr(err)
		outputStr = string(output)
		is.True(!strings.Contains(outputStr, "LABELS"))

		// Test short flag for reverse order
		output, err = exec(t, "search", runSearch, "Task", "-r", "-j")
		is.NoErr(err)
		listResult = &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), countTask)
		is.Equal(listResult.Tasks[0].Title, "In Progress Task")

		// Test short flag for status
		output, err = exec(t, "search", runSearch, "Task", "-s", "todo", "-j")
		is.NoErr(err)
		listResult = &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), countTask-1)

		// Test short flag for assigned
		output, err = exec(t, "search", runSearch, "Task", "-a", "first-user", "-j")
		is.NoErr(err)
		listResult = &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")

		// Test short flag for labels
		output, err = exec(t, "search", runSearch, "Task", "-l", "first", "-j")
		is.NoErr(err)
		listResult = &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "First Task")

		// Test short flag for unassigned
		output, err = exec(t, "search", runSearch, "Task", "-u", "-j")
		is.NoErr(err)
		listResult = &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "Unassigned Task")
	})

	t.Run("multiple combined flags", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "-l", "second", "-s", "todo", "--sort", "title", "-r", "--hide-extra", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		// Should filter by second label, todo status, sort by title reversed, output JSON
		is.Equal(len(listResult.Tasks), 1)
		is.Equal(listResult.Tasks[0].Title, "Second Task")
	})

	t.Run("search with parent filter", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// Since no parent tasks are created in test data, this should return empty
		output, err := exec(t, "search", runSearch, "Task", "-p", "1", "-j")
		is.NoErr(err)
		outputStr := strings.TrimSpace(string(output))
		is.Equal(outputStr, "{\"tasks\":[],\"pagination\":{\"total_results\":0,\"displayed_results\":0,\"offset\":0,\"limit\":25,\"has_more\":false}}") // empty ListResult JSON
	})

	t.Run("search with has-dependency filter", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// Since no dependencies are created in test data, this should return empty
		output, err := exec(t, "search", runSearch, "Task", "-c", "-j")
		is.NoErr(err)
		outputStr := strings.TrimSpace(string(output))
		is.Equal(outputStr, "{\"tasks\":[],\"pagination\":{\"total_results\":0,\"displayed_results\":0,\"offset\":0,\"limit\":25,\"has_more\":false}}") // empty ListResult JSON
	})

	t.Run("search with depended-on filter", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// Since no dependencies are created in test data, this should return empty
		output, err := exec(t, "search", runSearch, "Task", "-d", "-j")
		is.NoErr(err)
		outputStr := strings.TrimSpace(string(output))
		is.Equal(outputStr, "{\"tasks\":[],\"pagination\":{\"total_results\":0,\"displayed_results\":0,\"offset\":0,\"limit\":25,\"has_more\":false}}") // empty ListResult JSON
	})

	t.Run("search with invalid sort field", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "--sort", "invalidfield", "-j")
		// Command should still execute but may not sort properly
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), countTask)
	})

	t.Run("search with empty sort field", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "search", runSearch, "Task", "--sort", "", "-j")
		is.NoErr(err)
		listResult := &core.ListResult{}
		is.NoErr(json.Unmarshal(output, listResult))
		is.Equal(len(listResult.Tasks), countTask)
	})
}
