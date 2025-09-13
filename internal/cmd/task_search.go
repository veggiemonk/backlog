package cmd

import (
	"fmt"
	"os"
	"slices"

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
	setSearchFlags(searchCmd)
}

func setSearchFlags(cmd *cobra.Command) {
	// filtering
	cmd.Flags().StringVarP(&searchParent, "parent", "p", "", "Filter tasks by parent ID")
	cmd.Flags().StringSliceVarP(&searchStatus, "status", "s", nil, "Filter tasks by status")
	cmd.Flags().StringSliceVarP(&searchAssigned, "assigned", "a", nil, "Filter tasks by assigned names")
	cmd.Flags().StringSliceVarP(&searchLabels, "labels", "l", nil, "Filter tasks by labels")
	cmd.Flags().BoolVarP(&searchUnassigned, "unassigned", "u", false, "List tasks that have no one assigned")
	cmd.Flags().BoolVarP(&searchHasDependency, "has-dependency", "c", false, "Include tasks that have dependencies")
	cmd.Flags().BoolVarP(&searchDependedon, "depended-on", "d", false, "Filter tasks that are depended on by other tasks")
	// sorting
	cmd.Flags().StringVar(&searchSortFields, "sort", "", "Sort tasks by comma-separated fields (id, title, status, priority, created, updated)")
	cmd.Flags().BoolVarP(&searchReverseOrder, "reverse", "r", false, "Reverse the order of tasks")
	// column visibility
	cmd.Flags().BoolVarP(&searchHideExtraFields, "hide-extra", "e", false, "Hide extra fields (labels, priority, assigned)")
	// output format
	cmd.Flags().BoolVarP(&searchMarkdownOutput, "markdown", "m", false, "Print markdown table")
	cmd.Flags().BoolVarP(&searchJSONOutput, "json", "j", false, "Print JSON output")
}

func runSearch(cmd *cobra.Command, args []string) {
	query := args[0]
	sortFieldsSlice := parseSortFields(searchSortFields)
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
	if reverseOrder {
		slices.Reverse(tasks)
	}
	messagePrefix := ""
	if len(tasks) > 0 && !searchJSONOutput {
		messagePrefix = fmt.Sprintf("Found %d task(s) matching '%s':", len(tasks), query)
	} else if len(tasks) == 0 && !searchJSONOutput {
		messagePrefix = fmt.Sprintf("No tasks found matching '%s'.", query)
	}

	if err := renderTaskResults(cmd.OutOrStdout(), tasks, searchJSONOutput, searchMarkdownOutput, searchHideExtraFields, messagePrefix); err != nil {
		logging.Error("failed to render task results", "query", query, "error", err)
		os.Exit(1)
	}
}
