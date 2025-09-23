package core

import (
	"fmt"
	"strings"

	"github.com/agnivade/levenshtein"
)

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
var allStatuses = strings.Join(statuses, ",")

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
	return "", fmt.Errorf("only valid statuses are %q: %q %w", allStatuses, s, ErrInvalid)
}
