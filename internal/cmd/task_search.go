package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search tasks by content",
	Long:  `Search for tasks containing the specified query string.`,
	Args:  cobra.ExactArgs(1),
	Example: `
# Search for tasks containing "login" in any field
backlog search "login"

# Search for tasks containing "bug" 
backlog search "bug"

# Search for tasks assigned to a specific person
backlog search "@john"

# Search for tasks with specific labels
backlog search "frontend"

# Search in acceptance criteria
backlog search "validation"

# Search with markdown output
backlog search "api" --markdown

# Search with JSON output
backlog search "api" --json

# Search with additional columns displayed
backlog search "user" --labels --priority --assigned
	`,
	Run: runSearch,
}
var (
	searchParent        string
	searchStatus        []string
	searchAssigned      []string
	searchLabels        []string
	searchUnassigned    bool
	searchHasDependency bool
	searchDependedon    bool
	// sorting
	searchSortFields   string
	searchReverseOrder bool
	// column visibility
	searchHideExtraFields bool
	// output format
	searchMarkdownOutput bool
	searchJSONOutput     bool
)

// markdownOutput bool
// jsonOutput     bool
// hideExtraFields bool // defined in task_list.go

func init() {
	rootCmd.AddCommand(searchCmd)
	// filtering
	searchCmd.Flags().StringVarP(&searchParent, "parent", "p", "", "Filter tasks by parent ID")
	searchCmd.Flags().StringSliceVarP(&searchStatus, "status", "s", nil, "Filter tasks by status")
	searchCmd.Flags().StringSliceVarP(&searchAssigned, "assigned", "a", nil, "Filter tasks by assigned names")
	searchCmd.Flags().StringSliceVarP(&searchLabels, "labels", "l", nil, "Filter tasks by labels")
	searchCmd.Flags().BoolVarP(&searchUnassigned, "unassigned", "u", false, "List tasks that have no one assigned")
	searchCmd.Flags().BoolVarP(&searchHasDependency, "has-dependency", "c", false, "Include tasks that have dependencies")
	searchCmd.Flags().BoolVarP(&searchDependedon, "depended-on", "d", false, "Filter tasks that are depended on by other tasks")
	// sorting
	searchCmd.Flags().StringVar(&searchSortFields, "sort", "", "Sort tasks by comma-separated fields (id, title, status, priority, created, updated)")
	searchCmd.Flags().BoolVarP(&searchReverseOrder, "reverse", "r", false, "Reverse the order of tasks")
	// column visibility
	searchCmd.Flags().BoolVarP(&searchHideExtraFields, "hide-extra", "e", false, "Hide extra fields (labels, priority, assigned)")
	// output format
	searchCmd.Flags().BoolVarP(&searchMarkdownOutput, "markdown", "m", false, "Print markdown table")
	searchCmd.Flags().BoolVarP(&searchJSONOutput, "json", "j", false, "Print JSON output")
}

func runSearch(cmd *cobra.Command, args []string) {
	query := args[0]
	// Parse sort fields if provided
	var sortFieldsSlice []string
	if sortFields != "" {
		sortFieldsSlice = strings.Split(sortFields, ",")
		for i, field := range sortFieldsSlice {
			sortFieldsSlice[i] = strings.TrimSpace(field)
		}
	}
	params := core.ListTasksParams{
		Parent:        &searchParent,
		Status:        searchStatus,
		Assigned:      searchAssigned,
		Labels:        searchLabels,
		Unassigned:    searchUnassigned,
		HasDependency: searchHasDependency,
		DependedOn:    searchDependedon,
		Sort:          sortFieldsSlice,
		Reverse:       searchReverseOrder,
	}
	store := cmd.Context().Value(ctxKeyStore).(TaskStore)
	tasks, err := store.Search(query, params)
	if err != nil {
		logging.Error("failed to search tasks", "query", query, "error", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		if searchJSONOutput {
			fmt.Println("[]")
		} else {
			fmt.Printf("No tasks found matching '%s'.\n", query)
		}
		return
	}

	// Handle JSON output
	if searchJSONOutput {
		if err := json.NewEncoder(os.Stdout).Encode(tasks); err != nil {
			logging.Error("failed to encode JSON", "query", query, "error", err)
			os.Exit(1)
		}
		return
	}

	// Set table header based on hidden columns
	header := []string{"ID", "Status", "Title", "Dependencies"}
	if !searchHideExtraFields {
		header = append(header, "Labels", "Priority", "Assigned")
	}

	table := tableWriter(searchMarkdownOutput)
	table.Header(header)

	fmt.Printf("Found %d task(s) matching '%s':\n\n", len(tasks), query)

	for _, t := range tasks {
		row := []string{
			t.ID.String(),
			string(t.Status),
			t.Title,
			strings.Join(t.Dependencies, ", ")}
		if !searchHideExtraFields {
			row = append(row,
				strings.Join(t.Labels, ", "),
				t.Priority.String(),
				strings.Join(t.Assigned, ", "),
			)
		}
		if err := table.Append(row); err != nil {
			logging.Error("failed to append table row", "task_id", t.ID, "error", err)
			os.Exit(1)
		}
	}

	if err := table.Render(); err != nil {
		logging.Error("failed to render table", "query", query, "error", err)
		os.Exit(1)
	}
}
