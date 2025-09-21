package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/core"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all tasks",
	Long:    `Lists all tasks in the backlog except archived tasks.`,
	Example: generateExampleText(ListExamples),
	RunE:    runList,
}

var (
	filterParent     string
	filterPriority   string
	filterStatus     []string
	filterAssigned   []string
	filterLabels     []string
	filterUnassigned bool
	hasDependency    bool
	dependedon       bool
	// sorting
	sortFields   string
	reverseOrder bool
	// column visibility
	hideExtraFields bool
	// output format
	markdownOutput bool
	jsonOutput     bool
	// pagination
	limitFlag  int
	offsetFlag int
)

func init() {
	rootCmd.AddCommand(listCmd)
	setListFlags(listCmd)
}

func setListFlags(cmd *cobra.Command) {
	// filtering
	cmd.Flags().StringVarP(&filterParent, "parent", "p", "", "Filter tasks by parent ID")
	cmd.Flags().StringVar(&filterPriority, "priority", "", "Filter tasks by priority")
	cmd.Flags().StringSliceVarP(&filterStatus, "status", "s", nil, "Filter tasks by status")
	cmd.Flags().StringSliceVarP(&filterAssigned, "assigned", "a", nil, "Filter tasks by assigned names")
	cmd.Flags().StringSliceVarP(&filterLabels, "labels", "l", nil, "Filter tasks by labels")
	cmd.Flags().BoolVarP(&filterUnassigned, "unassigned", "u", false, "Filter tasks that have no one assigned")
	cmd.Flags().BoolVarP(&hasDependency, "has-dependency", "c", false, "Filter tasks that have dependencies")
	cmd.Flags().BoolVarP(&dependedon, "depended-on", "d", false, "Filter tasks that are depended on by other tasks")
	// sorting
	cmd.Flags().StringVar(&sortFields, "sort", "", "Sort tasks by comma-separated fields (id, title, status, priority, created, updated)")
	cmd.Flags().BoolVarP(&reverseOrder, "reverse", "r", false, "Reverse the order of tasks")
	// column visibility
	cmd.Flags().BoolVarP(&hideExtraFields, "hide-extra", "e", false, "Hide extra fields (labels, priority, assigned)")
	// output format
	cmd.Flags().BoolVarP(&markdownOutput, "markdown", "m", false, "print markdown table")
	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Print JSON output")
	// pagination
	cmd.Flags().IntVar(&limitFlag, "limit", 0, "Maximum number of tasks to return (0 means no limit)")
	cmd.Flags().IntVar(&offsetFlag, "offset", 0, "Number of tasks to skip from the beginning")
}

func runList(cmd *cobra.Command, args []string) error {
	sortFieldsSlice := parseSortFields(sortFields)

	var limit, offset *int
	if limitFlag > 0 {
		limit = &limitFlag
	}
	if offsetFlag > 0 {
		offset = &offsetFlag
	}

	// Apply configuration defaults and limits
	limit, offset = ApplyDefaultPagination(limit, offset)

	params := core.ListTasksParams{
		Parent:        &filterParent,
		Priority:      &filterPriority,
		Status:        filterStatus,
		Assigned:      filterAssigned,
		Labels:        filterLabels,
		Unassigned:    filterUnassigned,
		HasDependency: hasDependency,
		DependedOn:    dependedon,
		Sort:          sortFieldsSlice,
		Reverse:       reverseOrder,
		Limit:         limit,
		Offset:        offset,
	}

	store := cmd.Context().Value(ctxKeyStore).(TaskStore)

	// Get total count without pagination for metadata
	totalParams := params
	totalParams.Limit = nil
	totalParams.Offset = nil
	allTasks, err := store.List(totalParams)
	if err != nil {
		return fmt.Errorf("list tasks: %v", err)
	}
	totalCount := len(allTasks)

	// Get paginated results
	tasks, err := store.List(params)
	if err != nil {
		return fmt.Errorf("list tasks: %v", err)
	}

	// Create pagination info
	var paginationInfo *core.PaginationInfo
	if limit != nil || offset != nil {
		offsetVal := 0
		if offset != nil {
			offsetVal = *offset
		}
		limitVal := 0
		if limit != nil {
			limitVal = *limit
		}
		hasMore := (offsetVal + len(tasks)) < totalCount

		paginationInfo = &core.PaginationInfo{
			TotalResults:     totalCount,
			DisplayedResults: len(tasks),
			Offset:           offsetVal,
			Limit:            limitVal,
			HasMore:          hasMore,
		}
	}

	if err := renderTaskResultsWithPagination(cmd.OutOrStdout(), tasks, jsonOutput, markdownOutput, hideExtraFields, "", paginationInfo); err != nil {
		return fmt.Errorf("render task results: %v", err)
	}
	return nil
}

// parseSortFields parses a comma-separated string of sort fields
func parseSortFields(sortFields string) []string {
	if sortFields == "" {
		return nil
	}
	sortFieldsSlice := strings.Split(sortFields, ",")
	for i, field := range sortFieldsSlice {
		sortFieldsSlice[i] = strings.TrimSpace(field)
	}
	return sortFieldsSlice
}

// renderTaskResultsWithPagination renders a slice of tasks with pagination info
func renderTaskResultsWithPagination(w io.Writer, tasks []*core.Task, jsonOutput, markdownOutput, hideExtraFields bool, messagePrefix string, paginationInfo *core.PaginationInfo) error {
	// For JSON output with pagination info
	if jsonOutput && paginationInfo != nil {
		result := core.ListResult{
			Tasks:      tasks,
			Pagination: paginationInfo,
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		return nil
	}

	// Add pagination info to message prefix if not JSON output
	if paginationInfo != nil && !jsonOutput {
		if messagePrefix == "" {
			messagePrefix = fmt.Sprintf("Showing %d-%d of %d tasks",
				paginationInfo.Offset+1,
				paginationInfo.Offset+paginationInfo.DisplayedResults,
				paginationInfo.TotalResults)
			if paginationInfo.HasMore {
				messagePrefix += fmt.Sprintf(" (use --offset %d for more)",
					paginationInfo.Offset+paginationInfo.DisplayedResults)
			}
		} else {
			messagePrefix += fmt.Sprintf(" [%d-%d of %d total]",
				paginationInfo.Offset+1,
				paginationInfo.Offset+paginationInfo.DisplayedResults,
				paginationInfo.TotalResults)
		}
	}

	return renderTaskResults(w, tasks, jsonOutput, markdownOutput, hideExtraFields, messagePrefix)
}

// renderTaskResults renders a slice of tasks using the specified output format
func renderTaskResults(w io.Writer, tasks []*core.Task, jsonOutput, markdownOutput, hideExtraFields bool, messagePrefix string) error {
	// Handle empty task list
	if len(tasks) == 0 {
		switch {
		case jsonOutput:
			if _, err := fmt.Fprintln(w, "[]"); err != nil {
				return fmt.Errorf("writer: %v", err)
			}
		case markdownOutput:
			if _, err := fmt.Fprintln(w, "| No tasks found. |"); err != nil {
				return fmt.Errorf("writer: %v", err)
			}
		default:
			if _, err := fmt.Fprintln(w, "No tasks found."); err != nil {
				return fmt.Errorf("writer: %v", err)
			}
		}
		return nil
	}

	// Handle JSON output
	if jsonOutput {
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		return nil
	}

	// Print message prefix if provided
	if messagePrefix != "" {
		if _, err := fmt.Fprintf(w, "%s\n", messagePrefix); err != nil {
			return fmt.Errorf("writer: %v", err)
		}
	}

	// Set table header based on hidden columns
	header := []string{"ID", "Status", "Title", "Dependencies"}
	if !hideExtraFields {
		header = append(header, "Labels", "Priority", "Assigned")
	}

	table := tableWriter(w, markdownOutput)
	table.Header(header)

	for _, t := range tasks {
		row := []string{
			t.ID.String(),
			string(t.Status),
			t.Title,
			strings.Join(t.Dependencies, ", "),
		}
		if !hideExtraFields {
			row = append(row,
				strings.Join(t.Labels, ", "),
				t.Priority.String(),
				strings.Join(t.Assigned, ", "),
			)
		}
		if err := table.Append(row); err != nil {
			return fmt.Errorf("failed to append table row for task %s: %w", t.ID, err)
		}
	}

	if err := table.Render(); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	return nil
}

func tableWriter(w io.Writer, md bool) *tablewriter.Table {
	cfg := tablewriter.Config{
		Header: tw.CellConfig{
			Formatting: tw.CellFormatting{
				AutoFormat: tw.On,
				AutoWrap:   int(tw.Off),
			},
		},
		// MaxWidth: 150,
		Row: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignLeft}},
	}

	opts := []tablewriter.Option{}
	opts = append(opts, tablewriter.WithConfig(cfg))
	if md {
		opts = append(opts,
			tablewriter.WithConfig(tablewriter.Config{
				Header: tw.CellConfig{
					Alignment: tw.CellAlignment{Global: tw.AlignLeft},
				},
			}),
			tablewriter.WithRenderer(renderer.NewMarkdown()),
			tablewriter.WithRowAutoWrap(tw.WrapNone),
		)
	}
	table := tablewriter.NewTable(w, opts...)
	return table
}
