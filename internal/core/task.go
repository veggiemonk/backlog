// Package core defines all the core functionalities to work with tasks.
package core

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/veggiemonk/backlog/internal/logging"
	"go.yaml.in/yaml/v4"
)

var ErrInvalid = errors.New("invalid value")

// NewTask creates a new Task with default values.
func NewTask() *Task {
	return &Task{
		Status: StatusTodo, // Default status
	}
}

// Task represents a single task in the backlog.
// It includes metadata from the front matter and content from the markdown body.
// The task ID represents a hierarchical structure using dot notation (e.g., "1", "1.1", "1.2", "2").
// The task name is derived from the ID and prefixed with "task-" (e.g., "T1", "T1.1").
type Task struct {
	// --- Front Matter Fields ---

	ID           TaskID           `yaml:"id" json:"id"`
	Title        string           `yaml:"title" json:"title"`
	Status       Status           `yaml:"status" json:"status"`
	Parent       TaskID           `yaml:"parent" json:"parent"`
	Assigned     MaybeStringArray `yaml:"assigned,omitempty" json:"assigned,omitempty"`
	Labels       MaybeStringArray `yaml:"labels,omitempty" json:"labels,omitempty"`
	Dependencies MaybeStringArray `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	Priority     Priority         `yaml:"priority,omitempty" json:"priority,omitempty"`
	CreatedAt    time.Time        `yaml:"created_at" json:"created_at"`
	UpdatedAt    time.Time        `yaml:"updated_at,omitempty" json:"updated_at,omitzero"`
	History      []HistoryEntry   `yaml:"history,omitempty" json:"history,omitempty"`

	// --- Markdown Body Fields ---

	Description         string                `json:"description"`
	AcceptanceCriteria  []AcceptanceCriterion `json:"acceptance_criteria"`
	ImplementationPlan  string                `json:"implementation_plan"`
	ImplementationNotes string                `json:"implementation_notes"`
}

var slugRegex = regexp.MustCompile(`[^a-zA-Z0-9_.\(\)\[\]]+`)

const fileFormat = "%s-%s.md" // e.g., T1-implement-feature-x.md

// FileName creates a filename from task ID and a slugified title.
func (t *Task) FileName() string {
	// Create filename slug from title
	slug := strings.ToLower(t.Title)
	slug = slugRegex.ReplaceAllString(slug, "_")
	// Trim underscores from start and end
	slug = strings.Trim(slug, "_")
	// Limit length and ensure it's not empty
	if len(slug) > 50 {
		slug = slug[:50]
	}
	if slug == "" {
		slug = "untitled_task"
	}
	return fmt.Sprintf(fileFormat, t.ID.Name(), slug)
}

func (t *Task) Bytes() []byte {
	frontmatter := &Frontmatter{
		ID:           t.ID.String(),
		Title:        t.Title,
		Status:       string(t.Status),
		Assignee:     t.Assigned,
		Labels:       t.Labels,
		Parent:       t.Parent.String(),
		Priority:     t.Priority.String(),
		Dependencies: t.Dependencies,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		History:      t.History,
	}

	frontMatterBytes, err := yaml.Marshal(frontmatter)
	if err != nil {
		logging.Error("failed to marshal frontmatter", "task_id", t.ID, "error", err)
		return nil
	}

	var body bytes.Buffer
	body.WriteString(fmt.Sprintf("%s\n\n%s\n\n", descHeader, t.Description))
	body.WriteString(fmt.Sprintf("%s\n%s\n\n", acHeader, acStartComment))
	for _, ac := range t.AcceptanceCriteria {
		checked := " "
		if ac.Checked {
			checked = "x"
		}
		body.WriteString(fmt.Sprintf("- [%s] #%d %s\n", checked, ac.Index, ac.Text))
	}
	body.WriteString(fmt.Sprintf("\n%s\n\n", acEndComment))
	body.WriteString(fmt.Sprintf("%s\n\n%s\n\n", planHeader, t.ImplementationPlan))
	body.WriteString(fmt.Sprintf("%s\n\n%s\n", notesHeader, t.ImplementationNotes))

	// Combine front matter and body
	var fullContent bytes.Buffer
	fullContent.WriteString("---\n")
	fullContent.Write(frontMatterBytes)
	fullContent.WriteString("---\n")
	fullContent.Write(body.Bytes())

	return fullContent.Bytes()
}

// AcceptanceCriterion represents a single item in the acceptance criteria list.
type AcceptanceCriterion struct {
	Text    string `json:"text"`
	Checked bool   `json:"checked"`
	Index   int    `json:"index"`
}

// HistoryEntry represents a single entry in the task's history.
type HistoryEntry struct {
	Timestamp time.Time `yaml:"timestamp" json:"timestamp"`
	Change    string    `yaml:"change" json:"change"`
}

// RecordChange adds a history entry to the task for the given change
func RecordChange(task *Task, change string) {
	entry := HistoryEntry{
		Timestamp: time.Now().UTC(),
		Change:    change,
	}
	task.History = append(task.History, entry)
}
