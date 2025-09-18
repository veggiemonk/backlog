package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/agnivade/levenshtein"
	"go.yaml.in/yaml/v4"
)

type Priority int

const (
	PriorityUnknown Priority = iota
	PriorityLow
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

var (
	_ yaml.Unmarshaler = (*Priority)(nil)
	_ yaml.Marshaler   = (*Priority)(nil)
	_ json.Unmarshaler = (*Priority)(nil)
	_ json.Marshaler   = (*Priority)(nil)
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

// UnmarshalJSON implements json.Unmarshaler.
func (p *Priority) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var err error
	*p, err = ParsePriority(s)
	if err != nil {
		return err
	}
	return nil
}

func (p Priority) MarshalJSON() ([]byte, error) { return json.Marshal(p.String()) }

// UnmarshalYAML implements yaml.Unmarshaler.
func (p *Priority) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	var err error
	*p, err = ParsePriority(s)
	if err != nil {
		return err
	}
	return nil
}
func (p Priority) MarshalYAML() (any, error) { return p.String(), nil }

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
