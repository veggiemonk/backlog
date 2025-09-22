package core

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var errWrongIDFormat = errors.New("wrong id format")

func ParseTask(content []byte) (*Task, error) {
	matter, err := parseFrontMatter(content)
	if err != nil {
		return nil, fmt.Errorf("could not parse frontmatter: %w", err)
	}
	id, err := parseTaskID(matter.ID)
	if err != nil {
		return nil, fmt.Errorf("task ID %q: %w", matter.ID, err)
	}

	// Convert parent ID string to TaskID if present
	var pid TaskID
	if matter.Parent != "" {
		pid, err = parseTaskID(matter.Parent)
		if err != nil && matter.Parent != "" {
			return nil, fmt.Errorf("parent task ID %q: %w", matter.Parent, err)
		}
	}
	var status Status
	if matter.Status != "" {
		status, err = ParseStatus(matter.Status)
		if err != nil {
			return nil, fmt.Errorf("status %q: %w", matter.Status, err)
		}
	}

	var priority Priority
	if matter.Priority != "" {
		priority, err = ParsePriority(matter.Priority)
		if err != nil {
			return nil, fmt.Errorf("priority %q: %w", matter.Priority, err)
		}
	}

	task := &Task{
		ID:           id,
		Title:        matter.Title,
		Status:       status,
		Assigned:     matter.Assignee,
		Labels:       matter.Labels,
		Parent:       pid,
		Priority:     priority,
		Dependencies: matter.Dependencies,
		CreatedAt:    matter.CreatedAt,
		UpdatedAt:    matter.UpdatedAt,
		History:      matter.History,
	}

	parseMarkdownBody(task, content)
	return task, nil
}

func parseTaskIDfromFileName(fileName string) (TaskID, error) {
	// Find the part before the first '-' (not N=1, but N=2)
	parts := strings.SplitN(fileName, fieldSeparator, 2)
	if len(parts) < 2 {
		return TaskID{}, errWrongIDFormat
	}
	id, err := parseTaskID(parts[0])
	if err != nil {
		return TaskID{}, err
	}
	return id, nil
}

func parseMarkdownBody(task *Task, content []byte) {
	// This is a simplified parser. A more robust solution might use a proper markdown AST parser.
	sections := splitByHeaders(string(content))

	task.Description = getSectionContent(sections, descHeader)
	task.ImplementationPlan = getSectionContent(sections, planHeader)
	task.ImplementationNotes = getSectionContent(sections, notesHeader)

	acContent := getSectionContent(sections, acHeader)
	task.AcceptanceCriteria = parseAcceptanceCriteria(acContent)
}

func splitByHeaders(content string) map[string]string {
	headers := []string{descHeader, acHeader, planHeader, notesHeader}
	sections := make(map[string]string)

	for i, header := range headers {
		start := strings.Index(content, header)
		if start == -1 {
			continue
		}

		end := len(content)
		for j := i + 1; j < len(headers); j++ {
			nextHeaderPos := strings.Index(content, headers[j])
			if nextHeaderPos != -1 {
				end = nextHeaderPos
				break
			}
		}

		sectionContent := content[start+len(header) : end]
		sections[header] = strings.TrimSpace(sectionContent)
	}

	return sections
}

func getSectionContent(sections map[string]string, header string) string {
	if content, ok := sections[header]; ok {
		return content
	}
	return ""
}

func parseAcceptanceCriteria(content string) []AcceptanceCriterion {
	var criteria []AcceptanceCriterion
	scanner := bufio.NewScanner(strings.NewReader(content))
	re := regexp.MustCompile(`- \[( |x)\] #(\d+) (.*)`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 4 {
			index, _ := strconv.Atoi(matches[2])
			criteria = append(criteria, AcceptanceCriterion{
				Checked: matches[1] == "x",
				Index:   index,
				Text:    matches[3],
			})
		}
	}
	// Ensure ACs are sorted by index, as they might not be in order in the file.
	sort.Slice(criteria, func(i, j int) bool {
		return criteria[i].Index < criteria[j].Index
	})
	return criteria
}
