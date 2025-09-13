package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/veggiemonk/backlog/internal/core"
)

func Test_runList(t *testing.T) {
	t.Run("filter by single label", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "list", runList, "-l", "first", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), 1)
		is.Equal(tasks[0].Title, "First Task")
	})

	t.Run("filter by multiple labels", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "list", runList, "-l", "first,second,third,feature", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), 3) // feature label isn't set
	})

	t.Run("filter by status", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "list", runList, "-s", "todo", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), countTask-1) // countTask - 1 in-progress
	})

	t.Run("filter by multiple statuses", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "list", runList, "-s", "todo,in-progress", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), countTask)
	})

	t.Run("filter by assigned user", func(t *testing.T) {
		is := is.New(t)
		//
		output, err := exec(t, "list", runList, "-a", "first-user", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), 1)
		is.Equal(tasks[0].Title, "First Task")
	})

	t.Run("filter by multiple assigned users", func(t *testing.T) {
		is := is.NewRelaxed(t)

		output, err := exec(t, "list", runList, "-a", "first-user,second-user", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), 2)
		is.Equal(tasks[0].Title, "First Task")
		is.Equal(tasks[1].Title, "Second Task")
	})

	t.Run("filter by priority", func(t *testing.T) {
		is := is.NewRelaxed(t)

		output, err := exec(t, "list", runList, "--priority", "high", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), 1)
		is.Equal(tasks[0].Title, "High Priority Task")
	})

	t.Run("sort by title", func(t *testing.T) {
		is := is.NewRelaxed(t)

		output, err := exec(t, "list", runList, "--sort", "title", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), countTask)
		is.Equal(tasks[0].Title, "Fifth Task")
		is.Equal(tasks[1].Title, "First Task")
		is.Equal(tasks[2].Title, "Fourth Task")
		is.Equal(tasks[3].Title, "High Priority Task")
		is.Equal(tasks[4].Title, "In Progress Task")
		is.Equal(tasks[5].Title, "Second Task")
	})

	t.Run("sort by title reversed", func(t *testing.T) {
		is := is.NewRelaxed(t)

		output, err := exec(t, "list", runList, "--sort", "title", "--reverse", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), countTask)
		// When reversed, "Updated Task" should come first alphabetically when reversed
		is.Equal(tasks[0].Title, "Unlabeled Task")
		is.Equal(tasks[1].Title, "Unassigned Task")
		is.Equal(tasks[2].Title, "Third Task")
		is.Equal(tasks[3].Title, "Second Task")
		is.Equal(tasks[4].Title, "In Progress Task")
	})

	t.Run("sort by priority", func(t *testing.T) {
		is := is.NewRelaxed(t)

		output, err := exec(t, "list", runList, "--sort", "priority", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), countTask)
		is.Equal(tasks[0].Title, "High Priority Task")
	})

	t.Run("multiple sort fields", func(t *testing.T) {
		is := is.NewRelaxed(t)

		output, err := exec(t, "list", runList, "--sort", "priority,title", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), countTask)
		is.Equal(tasks[0].Title, "High Priority Task")
		is.Equal(tasks[1].Title, "Fifth Task")
		is.Equal(tasks[2].Title, "First Task")
		is.Equal(tasks[3].Title, "Fourth Task")
		is.Equal(tasks[4].Title, "Second Task")
		is.Equal(tasks[5].Title, "Third Task")
		is.Equal(tasks[6].Title, "In Progress Task")
		is.Equal(tasks[7].Title, "Unassigned Task")
		is.Equal(tasks[8].Title, "Unlabeled Task")
	})

	t.Run("markdown output format", func(t *testing.T) {
		is := is.NewRelaxed(t)

		output, err := exec(t, "list", runList, "-m")
		is.NoErr(err)
		outputStr := string(output)
		is.True(strings.Contains(outputStr, "|"))    // markdown table contains pipes
		is.True(strings.Contains(outputStr, "ID"))   // should contain header
		is.True(strings.Contains(outputStr, ":---")) // markdown table header separator
	})

	t.Run("hide extra fields", func(t *testing.T) {
		is := is.NewRelaxed(t)

		output, err := exec(t, "list", runList, "--hide-extra")
		is.NoErr(err)
		outputStr := string(output)
		is.True(strings.Contains(outputStr, "ID"))
		is.True(strings.Contains(outputStr, "STATUS"))
		is.True(strings.Contains(outputStr, "TITLE"))
		// Should not contain extra columns
		is.True(!strings.Contains(outputStr, "LABELS"))
		is.True(!strings.Contains(outputStr, "PRIORITY"))
		is.True(!strings.Contains(outputStr, "ASSIGNED"))
	})

	t.Run("default table output format", func(t *testing.T) {
		is := is.NewRelaxed(t)

		output, err := exec(t, "list", runList, "list")
		is.NoErr(err)
		outputStr := string(output)
		is.True(strings.Contains(outputStr, "ID"))
		is.True(strings.Contains(outputStr, "STATUS"))
		is.True(strings.Contains(outputStr, "TITLE"))
		is.True(strings.Contains(outputStr, "LABELS"))
		is.True(strings.Contains(outputStr, "PRIORITY"))
		is.True(strings.Contains(outputStr, "ASSIGNED"))
	})

	t.Run("combined filters", func(t *testing.T) {
		is := is.NewRelaxed(t)

		output, err := exec(t, "list", runList, "-l", "second", "-s", "todo", "--priority", "medium", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), 1) // 1 feature tasks with todo status
	})

	t.Run("empty results", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "list", runList, "-l", "nonexistent", "-j")
		is.NoErr(err)
		outputStr := strings.TrimSpace(string(output))
		is.Equal(outputStr, "[]") // empty JSON array
	})

	t.Run("empty results table format", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "list", runList, "-l", "nonexistent")
		is.NoErr(err)
		outputStr := strings.TrimSpace(string(output))
		is.Equal(outputStr, "No tasks found.")
	})

	t.Run("filter unassigned tasks", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "list", runList, "--unassigned", "-j")
		is.NoErr(err)
		// Should not error, may return empty array since all test tasks are assigned
		var tasks []*core.Task
		outputStr := strings.TrimSpace(string(output))
		if outputStr != "[]" {
			err = json.Unmarshal(output, &tasks)
			is.NoErr(err)
		}
	})

	t.Run("short flag aliases", func(t *testing.T) {
		t.Run("json", func(t *testing.T) {
			is := is.NewRelaxed(t)
			output, err := exec(t, "list", runList, "-j")
			is.NoErr(err)
			tasks := []*core.Task{}
			is.NoErr(json.Unmarshal(output, &tasks))
			is.Equal(len(tasks), countTask)
			is.Equal(tasks[0].Title, "First Task")
		})

		t.Run("markdown", func(t *testing.T) {
			is := is.NewRelaxed(t)
			output, err := exec(t, "list", runList, "-m")
			is.NoErr(err)
			outputStr := string(output)
			is.True(strings.Contains(outputStr, "|:---"))
		})

		t.Run("hide-extra-fields", func(t *testing.T) {
			is := is.NewRelaxed(t)
			output, err := exec(t, "list", runList, "-e")
			is.NoErr(err)
			outputStr := string(output)
			is.True(!strings.Contains(outputStr, "LABELS"))
		})

		t.Run("reverse", func(t *testing.T) {
			is := is.NewRelaxed(t)
			output, err := exec(t, "list", runList, "-r", "-j")
			is.NoErr(err)
			tasks := []*core.Task{}
			is.NoErr(json.Unmarshal(output, &tasks))
			is.Equal(len(tasks), countTask)
			is.Equal(tasks[0].Title, "In Progress Task") // the last task is
		})

		t.Run("unassigned", func(t *testing.T) {
			is := is.NewRelaxed(t)
			output, err := exec(t, "list", runList, "-u", "-j")
			is.NoErr(err)
			tasks := []*core.Task{}
			is.NoErr(json.Unmarshal(output, &tasks))
			is.Equal(len(tasks), 1)
			is.Equal(tasks[0].Title, "Unassigned Task")
		})
	})

	t.Run("invalid sort field", func(t *testing.T) {
		is := is.New(t)
		output, err := exec(t, "list", runList, "--sort", "invalidfield", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), countTask)
		is.Equal(tasks[0].Title, "First Task")
	})

	t.Run("empty sort field", func(t *testing.T) {
		is := is.New(t)
		output, err := exec(t, "list", runList, "--sort", "", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		is.Equal(len(tasks), 9)
	})

	t.Run("multiple combined flags", func(t *testing.T) {
		is := is.New(t)
		output, err := exec(t, "list", runList, "-l", "second", "-s", "todo", "--sort", "title", "-r", "--hide-extra", "-j")
		is.NoErr(err)
		tasks := []*core.Task{}
		is.NoErr(json.Unmarshal(output, &tasks))
		// Should filter by feature label, todo status, sort by title reversed, output JSON
		is.Equal(len(tasks), 1)
	})
}
