package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `Lists all tasks in the backlog except archived tasks.`,
	Example: `
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
	`,
	Run: runList,
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
}

func runList(cmd *cobra.Command, args []string) {
	sortFieldsSlice := parseSortFields(sortFields)
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
	}

	store := cmd.Context().Value(ctxKeyStore).(TaskStore)
	tasks, err := store.List(params)
	if err != nil {
		logging.Error("failed to list tasks", "error", err)
		os.Exit(1)
	}
	if err := renderTaskResults(cmd.OutOrStdout(), tasks, jsonOutput, markdownOutput, hideExtraFields, ""); err != nil {
		logging.Error("failed to render task results", "error", err)
		os.Exit(1)
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

// renderTaskResults renders a slice of tasks using the specified output format
func renderTaskResults(w io.Writer, tasks []*core.Task, jsonOutput, markdownOutput, hideExtraFields bool, messagePrefix string) error {
	// Handle empty task list
	if len(tasks) == 0 {
		switch {
		case jsonOutput:
			fmt.Fprintln(w, "[]")
		case markdownOutput:
			fmt.Fprintln(w, "| No tasks found. |")
		default:
			fmt.Fprintln(w, "No tasks found.")
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
		fmt.Fprintf(w, "%s\n", messagePrefix)
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
