package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/core"
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

# Search with pagination
backlog search "api" --limit 5                  # Show first 5 search results
backlog search "bug" --limit 3 --offset 5       # Show 3 results starting from 6th match
backlog search "feature" --status "todo" --limit 10  # Show first 10 "todo" feature results
	`,
	RunE: runSearch,
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
	// pagination
	searchLimitFlag  int
	searchOffsetFlag int
)

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
	// pagination
	cmd.Flags().IntVar(&searchLimitFlag, "limit", 25, "Maximum number of tasks to return (0 means no limit)")
	cmd.Flags().IntVar(&searchOffsetFlag, "offset", 0, "Number of tasks to skip from the beginning")
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]
	sortFieldsSlice := parseSortFields(searchSortFields)

	var limit, offset *int
	// Track if limit was explicitly set by user (not default)
	limitExplicit := cmd.Flags().Changed("limit")
	if searchLimitFlag > 0 {
		limit = &searchLimitFlag
	}
	if searchOffsetFlag > 0 {
		offset = &searchOffsetFlag
	}

	// Apply configuration defaults and limits
	limit, offset = ApplyDefaultPagination(limit, offset)

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
		Limit:         limit,
		Offset:        offset,
	}
	store := cmd.Context().Value(ctxKeyStore).(TaskStore)

	// Get total search count without pagination for metadata
	totalParams := params
	totalParams.Limit = nil
	totalParams.Offset = nil
	allTasks, err := store.Search(query, totalParams)
	if err != nil {
		return fmt.Errorf("failed to search tasks for query %q: %w", query, err)
	}
	totalCount := len(allTasks)

	// Get paginated search results
	tasks, err := store.Search(query, params)
	if err != nil {
		return fmt.Errorf("failed to search tasks for query %q: %w", query, err)
	}

	messagePrefix := ""
	if !searchJSONOutput {
		switch {
		case totalCount == 0:
			messagePrefix = fmt.Sprintf("No tasks found matching '%s'.", query)
		case totalCount > 0:
			if limit != nil || offset != nil {
				// Show pagination info in search prefix
				messagePrefix = fmt.Sprintf("Found %d task(s) matching '%s':", totalCount, query)
			} else {
				messagePrefix = fmt.Sprintf("Found %d task(s) matching '%s':", len(tasks), query)
			}
		}
	}

	// Create pagination info only when pagination was explicitly requested
	var paginationInfo *core.PaginationInfo
	// Only create pagination info if offset was provided or limit was explicitly set by user
	if (offset != nil && *offset > 0) || limitExplicit {
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

	if err := renderTaskResultsWithPagination(cmd.OutOrStdout(), tasks, searchJSONOutput, searchMarkdownOutput, searchHideExtraFields, messagePrefix, paginationInfo); err != nil {
		return fmt.Errorf("failed to render task results for query %q: %w", query, err)
	}
	return nil
}
