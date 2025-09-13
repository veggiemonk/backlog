// Package core defines all the core functionalities to work with tasks.
package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/agnivade/levenshtein"
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

// Status represents the state of a task.
type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
	StatusArchived   Status = "archived"
	StatusRejected   Status = "rejected"
)

var statuses = []string{
	string(StatusTodo),
	string(StatusInProgress),
	string(StatusDone),
	string(StatusCancelled),
	string(StatusArchived),
	string(StatusRejected),
}

func printValidStatuses() string {
	return strings.Join(statuses, ",")
}

func ParseStatus(s string) (Status, error) {
	if s == "" {
		return StatusTodo, nil
	}
	sc := strings.ReplaceAll(strings.ToLower(s), " ", "")
	for _, validStatus := range statuses {
		distance := levenshtein.ComputeDistance(sc, validStatus)
		if distance < 3 {
			return Status(validStatus), nil
		}
	}
	return "", fmt.Errorf("only valid statuses are %q: %q %w", printValidStatuses(), s, ErrInvalid)
}

type Priority int

const (
	PriorityUnknown Priority = iota
	PriorityLow
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

var priorities = []string{
	"unknown",
	"low",
	"medium",
	"high",
	"critical",
}

func (p Priority) String() string {
	switch p {
	case PriorityUnknown:
		return "unknown"
	case PriorityLow:
		return "low"
	case PriorityMedium:
		return "medium"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return ""
	}
}

func ParsePriority(s string) (Priority, error) {
	if s == "" {
		return PriorityUnknown, nil
	}

	sp := strings.ReplaceAll(strings.ToLower(s), " ", "")
	for _, validPriority := range priorities {
		distance := levenshtein.ComputeDistance(sp, validPriority)
		if distance < 3 {
			switch validPriority {
			case "unknown":
				return PriorityUnknown, nil
			case "low":
				return PriorityLow, nil
			case "medium":
				return PriorityMedium, nil
			case "high":
				return PriorityHigh, nil
			case "critical":
				return PriorityCritical, nil
			}
		}
	}
	return PriorityUnknown, fmt.Errorf("only valid priorities are %q: %q %w", strings.Join(priorities, ","), s, ErrInvalid)
}

// MaybeStringArray is a custom type that can unmarshal from either a single string or an array of strings.
// It also marshals back to the same format, preserving the original structure.
// origin: https://carlosbecker.com/posts/go-custom-marshaling/
type MaybeStringArray []string

var (
	_ yaml.Unmarshaler = &MaybeStringArray{}
	_ yaml.Marshaler   = MaybeStringArray{}
	_ json.Unmarshaler = &MaybeStringArray{}
	_ json.Marshaler   = MaybeStringArray{}
)

func (a *MaybeStringArray) ToSlice() []string {
	if a == nil {
		return nil
	}
	return []string(*a)
}

func (a *MaybeStringArray) UnmarshalYAML(value *yaml.Node) error {
	var slice []string
	if err := value.Decode(&slice); err == nil {
		*a = slice
		return nil
	}

	var single string
	if err := value.Decode(&single); err != nil {
		return err
	}
	*a = []string{single}
	return nil
}

func (a MaybeStringArray) MarshalYAML() (any, error) {
	switch len(a) {
	case 0:
		return nil, nil
	case 1:
		return a[0], nil
	default:
		return []string(a), nil
	}
}

func (a *MaybeStringArray) UnmarshalJSON(data []byte) error {
	var slice []string
	if err := json.Unmarshal(data, &slice); err == nil {
		*a = slice
		return nil
	}

	var single string
	if err := json.Unmarshal(data, &single); err != nil {
		return err
	}
	*a = []string{single}
	return nil
}

func (a MaybeStringArray) MarshalJSON() ([]byte, error) {
	switch len(a) {
	case 0:
		return json.Marshal(nil)
	case 1:
		return json.Marshal(a[0])
	default:
		return json.Marshal([]string(a))
	}
}
