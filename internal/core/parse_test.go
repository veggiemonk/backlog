package core

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestGetNextTaskID(t *testing.T) {
	filenames := []string{
		"T24.1-cli-kanban-board-milestone-view.md",
		"T200-Add-Claude-Code-integration-with-workflow-commands-during-init.md",
		"T208-Add-paste-as-markdown-support-in-Web-UI.md",
		"T217-Create-web-UI-for-sequences-with-drag-and-drop.md",
		"T217.02-Sequences-web-UI-list-sequences.md",
		"T217.03-Sequences-web-UI-move-tasks-and-update-dependencies.md",
		"T217.04-Sequences-web-UI-tests.md",
		"T218-Update-documentation-and-tests-for-sequences.md",
		"T222-Improve-task-and-subtask-visualization-in-web-UI.md",
	}

	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")

	is := is.New(t)
	for _, name := range filenames {
		_, err := fs.Create(".backlog/" + name)
		is.NoErr(err)
	}

	t.Run("next top-level ID", func(t *testing.T) {
		is := is.New(t)
		nextID, err := store.getNextTaskID()
		is.NoErr(err)
		is.Equal("T223", nextID.Name())
	})

	t.Run("next subtask ID", func(t *testing.T) {
		is := is.New(t)
		nextID, err := store.getNextTaskID(217)
		is.NoErr(err)
		is.Equal("T217.05", nextID.Name())
	})

	t.Run("next subtask ID deeper level", func(t *testing.T) {
		is := is.New(t)
		nextID, err := store.getNextTaskID(217, 3)
		is.NoErr(err)
		is.Equal("T217.03.01", nextID.Name())
	})
}

func TestWriteTask(t *testing.T) {
	is := is.New(t)
	var err error

	task := NewTask()
	task.ID = TaskID{seg: []int{1, 2, 3}}
	task.Parent = TaskID{seg: []int{1, 2}}
	task.Title = "Test Task"
	task.Description = "This is a test task."
	task.Status = "todo"
	task.Assigned = []string{"John Doe"}
	task.Labels = []string{"test", "unit"}

	task.Priority, err = ParsePriority("High")
	is.NoErr(err)

	task.AcceptanceCriteria = []AcceptanceCriterion{
		{Text: "Acceptance Criteria 1", Checked: true, Index: 1},
		{Text: "Acceptance Criteria 2", Checked: false, Index: 2},
	}
	task.ImplementationPlan = "This is the implementation plan."
	task.ImplementationNotes = "This is the implementation notes."

	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")
	err = store.write(task)
	is.NoErr(err)

	// Check if the file was actually created in the memory fs
	// The filename should be generated from ID and title
	expectedFilePath := ".backlog/T01.02.03-test_task.md"
	exists, err := afero.Exists(fs, expectedFilePath)
	is.NoErr(err)
	is.True(exists)

	// Check the content of the created file
	contentBytes, err := afero.ReadFile(fs, expectedFilePath)
	is.NoErr(err)
	content := string(contentBytes)

	is.True(strings.Contains(content, "title: Test Task"))
	is.True(strings.Contains(content, "parent: \"01.02\""))
	is.True(strings.Contains(content, "## Description\n\nThis is a test task."))
	is.True(strings.Contains(content, "## Acceptance Criteria"))
	is.True(strings.Contains(content, "<!-- AC:BEGIN -->"))
	is.True(strings.Contains(content, "- [x] #1 Acceptance Criteria 1"))
	is.True(strings.Contains(content, "- [ ] #2 Acceptance Criteria 2"))
	is.True(strings.Contains(content, "<!-- AC:END -->"))
	is.True(strings.Contains(content, "## Implementation Plan\n\nThis is the implementation plan."))
	is.True(strings.Contains(content, "## Implementation Notes\n\nThis is the implementation notes."))
	is.True(strings.Contains(content, "status: todo"))
	is.True(strings.Contains(content, "assignee"))
	is.True(strings.Contains(content, "John Doe"))
	is.True(strings.Contains(content, "labels:"))
	is.True(strings.Contains(content, "- test"))
	is.True(strings.Contains(content, "- unit"))
	is.True(strings.Contains(content, "priority: high"))
}

func TestParseJSONArray(t *testing.T) {
	is := is.New(t)

	output := `[
  {
    "id": "01",
    "title": "Initial Task",
    "status": "todo",
    "parent": "",
    "assigned": "initial-user",
    "labels": "bug",
    "priority": "medium",
    "created_at": "2025-09-13T09:00:59.522213Z",
    "description": "Initial description.",
    "acceptance_criteria": [
      {
        "text": "Initial AC.",
        "checked": false,
        "index": 1
      }
    ],
    "implementation_plan": "Initial implementation plan.",
    "implementation_notes": "Initial implementation notes."
  }
]`
	var tasks []Task
	err := json.Unmarshal([]byte(output), &tasks)
	is.NoErr(err)

	is.Equal(len(tasks), 1)
}
