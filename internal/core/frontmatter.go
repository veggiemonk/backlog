package core

import (
	"time"

	"go.yaml.in/yaml/v4"
)

type Frontmatter struct {
	ID           string           `yaml:"id"`
	Title        string           `yaml:"title"`
	Status       string           `yaml:"status"`
	Assignee     MaybeStringArray `yaml:"assignee,omitempty"`
	Labels       MaybeStringArray `yaml:"labels,omitempty"`
	Dependencies MaybeStringArray `yaml:"dependencies,omitempty"`
	Parent       string           `yaml:"parent,omitempty"`
	Priority     string           `yaml:"priority,omitempty"`
	CreatedAt    time.Time        `yaml:"created_at"`
	UpdatedAt    time.Time        `yaml:"updated_at,omitempty"`
	History      []HistoryEntry   `yaml:"history,omitempty"`
}

func parseFrontMatter(content []byte) (*Frontmatter, error) {
	var matter Frontmatter
	err := yaml.Unmarshal(content, &matter)
	if err != nil {
		return nil, err
	}
	return &matter, nil
}
