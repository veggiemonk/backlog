package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/urfave/cli/v3"
	"github.com/veggiemonk/backlog/internal/core"
)

const listExamples = `
# List all tasks
backlog list                                    # List all tasks with all columns
backlog list --status "todo"                    # List tasks with status "todo"
backlog list --status "todo,in-progress"        # List tasks with status "todo" or "in-progress"
backlog list --status "done"                    # List tasks with status "done"
backlog list --parent "12345"                   # List tasks that are sub-tasks of the task with ID "12345"
backlog list --status "todo" --parent "12345"   # List "todo" sub-tasks of task "12345"
backlog list --assigned "alice"                 # List tasks assigned to alice
backlog list --unassigned                       # List tasks that have no one assigned
backlog list --labels "bug"                     # List tasks containing the label "bug"
backlog list --labels "bug,feature"             # List tasks containing either "bug" or "feature" labels
backlog list --priority "high"                  # List all high priority tasks

# Search
backlog list --query "refactor"                 # Search for tasks with the word "refactor" in them

# dependency filters
backlog list --has-dependency                   # List tasks that have at least one dependency
backlog list --depended-on                      # List tasks that are depended on by other tasks
backlog list --depended-on --status "todo"      # List all the blocking tasks.

# column visibility
backlog list --hide-extra                       # Hide extra fields (labels, priority, assigned)
backlog list -e                                 # Hide extra fields (labels, priority, assigned)
backlog list --status "todo" --hide-extra       # List "todo" tasks with minimal columns

# sorting
backlog list --sort "priority"                  # Sort tasks by priority
backlog list --sort "updated,priority"          # Sort tasks by updated date, then priority
backlog list --sort "status,created"            # Sort tasks by status, then creation date
backlog list --reverse                          # Reverse the order of tasks
backlog list --sort "priority" --reverse        # Sort by priority in reverse order
backlog list --status "todo" \
    --priority "medium"  \
    --sort "priority"    \
    --reverse                               # Combine all options

# output format
backlog list -m                                 # List tasks in markdown format
backlog list -markdown                          # List tasks in markdown format
backlog list --json                             # List tasks in JSON format
backlog list -j                                 # List tasks in JSON format
backlog list --status "todo" --json             # List "todo" tasks in JSON format

# pagination
backlog list --limit 10                         # List first 10 tasks
backlog list --limit 5 --offset 10              # List 5 tasks starting from 11th task
backlog list --status "todo" --limit 3          # List first 3 "todo" tasks
backlog list --sort "priority" --limit 10       # List top 10 tasks by priority
`

func newListCommand(rt *runtime) *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List all tasks",
		Description: "Lists all tasks in the backlog except archived tasks.\n\nExamples:\n" + listExamples,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "parent", Aliases: []string{"p"}, Usage: "Filter tasks by parent ID"},
			&cli.StringFlag{Name: "priority", Usage: "Filter tasks by priority"},
			&cli.StringSliceFlag{Name: "status", Aliases: []string{"s"}, Usage: "Filter tasks by status"},
			&cli.StringSliceFlag{Name: "assigned", Aliases: []string{"a"}, Usage: "Filter tasks by assigned names"},
			&cli.StringSliceFlag{Name: "labels", Aliases: []string{"l"}, Usage: "Filter tasks by labels"},
			&cli.StringFlag{Name: "query", Aliases: []string{"q"}, Usage: "Search query to filter tasks by"},
			&cli.BoolFlag{Name: "unassigned", Aliases: []string{"u"}, Usage: "Filter tasks that have no one assigned"},
			&cli.BoolFlag{Name: "has-dependency", Aliases: []string{"c"}, Usage: "Filter tasks that have dependencies"},
			&cli.BoolFlag{Name: "depended-on", Aliases: []string{"d"}, Usage: "Filter tasks that are depended on by other tasks"},
			&cli.StringFlag{Name: "sort", Usage: "Sort tasks by comma-separated fields (id, title, status, priority, created, updated)"},
			&cli.BoolFlag{Name: "reverse", Aliases: []string{"r"}, Usage: "Reverse the order of tasks"},
			&cli.BoolFlag{Name: "hide-extra", Aliases: []string{"e"}, Usage: "Hide extra fields (labels, priority, assigned)"},
			&cli.BoolFlag{Name: "markdown", Aliases: []string{"m"}, Usage: "Print markdown table"},
			&cli.BoolFlag{Name: "json", Aliases: []string{"j"}, Usage: "Print JSON output"},
			&cli.IntFlag{Name: "limit", Usage: "Maximum number of tasks to return (0 means no limit)"},
			&cli.IntFlag{Name: "offset", Usage: "Number of tasks to skip from the beginning"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() > 0 {
				return cli.Exit("list does not accept positional arguments", 1)
			}

			store := rt.store
			if store == nil {
				return fmt.Errorf("task store not initialized")
			}

			params := core.ListTasksParams{
				Parent:        cmd.String("parent"),
				Priority:      cmd.String("priority"),
				Status:        cmd.StringSlice("status"),
				Assigned:      cmd.StringSlice("assigned"),
				Labels:        cmd.StringSlice("labels"),
				Query:         cmd.String("query"),
				Unassigned:    cmd.Bool("unassigned"),
				HasDependency: cmd.Bool("has-dependency"),
				DependedOn:    cmd.Bool("depended-on"),
				Sort:          parseSortFields(cmd.String("sort")),
				Reverse:       cmd.Bool("reverse"),
				Limit:         cmd.Int("limit"),
				Offset:        cmd.Int("offset"),
			}

			listResult, err := store.List(params)
			if err != nil {
				return fmt.Errorf("failed to list tasks: %w", err)
			}

			if err := renderTaskResultsWithPagination(cmd.Root().Writer, listResult, cmd.Bool("json"), cmd.Bool("markdown"), cmd.Bool("hide-extra"), ""); err != nil {
				return fmt.Errorf("failed to render task results: %w", err)
			}
			return nil
		},
	}
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
func renderTaskResultsWithPagination(w io.Writer, listResult core.ListResult, jsonOutput, markdownOutput, hideExtraFields bool, messagePrefix string) error {
	// For JSON output with pagination info
	if jsonOutput && listResult.Pagination != nil {
		if err := json.NewEncoder(w).Encode(listResult); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		return nil
	}

	// Add pagination info to message prefix if not JSON output
	if listResult.Pagination != nil && !jsonOutput {
		if messagePrefix == "" {
			messagePrefix = fmt.Sprintf("Showing %d-%d of %d tasks",
				listResult.Pagination.Offset+1,
				listResult.Pagination.Offset+listResult.Pagination.DisplayedResults,
				listResult.Pagination.TotalResults)
			if listResult.Pagination.HasMore {
				messagePrefix += fmt.Sprintf(" (use --offset %d for more)",
					listResult.Pagination.Offset+listResult.Pagination.DisplayedResults)
			}
		} else {
			messagePrefix += fmt.Sprintf(" [%d-%d of %d total]",
				listResult.Pagination.Offset+1,
				listResult.Pagination.Offset+listResult.Pagination.DisplayedResults,
				listResult.Pagination.TotalResults)
		}
	}

	return renderTaskResults(w, listResult.Tasks, jsonOutput, markdownOutput, hideExtraFields, messagePrefix)
}

// renderTaskResults renders a slice of tasks using the specified output format
func renderTaskResults(w io.Writer, tasks []core.Task, jsonOutput, markdownOutput, hideExtraFields bool, messagePrefix string) error {
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

	if jsonOutput {
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		return nil
	}

	if messagePrefix != "" {
		if _, err := fmt.Fprintf(w, "%s\n", messagePrefix); err != nil {
			return fmt.Errorf("writer: %v", err)
		}
	}

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
		Row: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignLeft}},
	}

	opts := []tablewriter.Option{tablewriter.WithConfig(cfg)}
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
	return tablewriter.NewTable(w, opts...)
}
