package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v4"
)

const (
	// TaskIDPrefix is the prefix used for task filenames.
	TaskIDPrefix       = "T"
	fieldSeparator     = "-"
	fieldTaskSeperator = "."
)

var _ yaml.Unmarshaler = &TaskID{}
var _ yaml.Marshaler = TaskID{}
var _ json.Marshaler = TaskID{}
var _ json.Unmarshaler = &TaskID{}

type TaskID struct {
	seg []int
}

// parseTaskID parses a task ID string (e.g., "T1.2.3" or "1.2.3") into a TaskID struct.
func parseTaskID(id string) (TaskID, error) {
	if strings.HasPrefix(id, TaskIDPrefix) {
		id = id[1:]
	}
	parts := strings.Split(id, fieldTaskSeperator)
	var segments []int
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return TaskID{}, fmt.Errorf("invalid segment in task ID %q: %w", id, err)
		}
		segments = append(segments, num)
	}
	return TaskID{seg: segments}, nil
}

func (t TaskID) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *TaskID) UnmarshalText(text []byte) error {
	parsed, err := parseTaskID(string(text))
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

func (t TaskID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
}

func (t *TaskID) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	parsed, err := parseTaskID(str)
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

func (t TaskID) MarshalYAML() (any, error) {
	return t.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (t *TaskID) UnmarshalYAML(value *yaml.Node) error {
	var str string
	if err := value.Decode(&str); err != nil {
		return err
	}
	parsed, err := parseTaskID(str)
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

// String returns the string representation of the TaskID (e.g., "1.02.03").
func (t TaskID) String() string {
	buf := bytes.Buffer{}
	for i, s := range t.seg {
		if i > 0 {
			buf.WriteString(".")
		}
		buf.WriteString(fmt.Sprintf("%02d", s))
	}
	return buf.String()
}

// Name returns the filename prefix for the task, e.g., "T1.02.03".
func (t TaskID) Name() string {
	return TaskIDPrefix + t.String()
}

// HasSubTasks returns true if the task ID has subtask segments (i.e., more than one segment).
func (t TaskID) HasSubTasks() bool {
	return len(t.seg) > 1
}

func (t TaskID) Less(other TaskID) bool {
	// Compare segment by segment
	minLen := len(t.seg)
	if len(other.seg) < minLen {
		minLen = len(other.seg)
	}
	for i := 0; i < minLen; i++ {
		if t.seg[i] < other.seg[i] {
			return true
		} else if t.seg[i] > other.seg[i] {
			return false
		}
	}
	// If all compared segments are equal, the shorter ID is "less"
	return len(t.seg) < len(other.seg)
}

func (t TaskID) Equals(other TaskID) bool {
	if len(t.seg) != len(other.seg) {
		return false
	}
	for i := range t.seg {
		if t.seg[i] != other.seg[i] {
			return false
		}
	}
	return true
}

func (t TaskID) Parent() *TaskID {
	if t.HasSubTasks() {
		return &TaskID{seg: t.seg[:len(t.seg)-1]}
	}
	return nil
}
func (t TaskID) NextSubTaskID() TaskID {
	newSeg := make([]int, len(t.seg))
	copy(newSeg, t.seg)
	newSeg = append(newSeg, 1)
	return TaskID{seg: newSeg}
}

func (t TaskID) NextSiblingID() TaskID {
	if len(t.seg) == 0 {
		return TaskID{seg: []int{1}}
	}
	newSeg := make([]int, len(t.seg))
	copy(newSeg, t.seg)
	newSeg[len(newSeg)-1]++
	return TaskID{seg: newSeg}
}
