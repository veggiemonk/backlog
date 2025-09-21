## backlog edit

Edit an existing task

### Synopsis

Edit an existing task by providing its ID and flags for the fields to update.

```
backlog edit <id> [flags]
```

### Examples

```
# Change Title

# Use the -t or --title flag to give the task a new title.
backlog backlog edit "42" --title "Fix the main login button styling"

# Update Description

# Use the -d or --description flag to replace the existing description with a new one.
backlog backlog edit "42" --description "The login button on the homepage is misaligned on mobile devices. It should be centered."

# Change Status

# Update the task's progress by changing its status with the -s or --status flag.
backlog backlog edit "42" --status "in-progress"

# Re-assign to Single Person

# You can change the assigned names for a task using the -a or --assignee flag.
backlog backlog edit "42" --assigned "jordan"

# Re-assign to Multiple People

backlog backlog edit "42" --assigned "jordan,casey"

# Update Labels

# Use the -l or --labels flag to replace the existing labels.
backlog backlog edit "42" --labels "bug,frontend"

# Change Priority

# Adjust the task's priority with the --priority flag.
backlog backlog edit "42" --priority "high"

# Add Acceptance Criteria

# Add a new AC
backlog backlog edit "42" --ac "The button is centered on screens smaller than 576px."

# Check Acceptance Criteria

# Check the first AC (assuming it's at index 1)
backlog backlog edit "42" --check-ac "1"

# Uncheck Acceptance Criteria

# Uncheck the first AC
backlog backlog edit "42" --uncheck-ac "1"

# Remove Acceptance Criteria

# Remove the second AC (at index 2)
backlog backlog edit "42" --remove-ac "2"

# Change Parent Task

# Move a task to be a sub-task of a different parent using the -p or --parent flag.
backlog backlog edit "42" --parent "18"

# Remove Parent

# To remove a parent, pass an empty string
backlog backlog edit "42" --parent

# Add Implementation Notes

# Use the --notes flag to add or update technical notes for implementation.
backlog backlog edit "42" --notes "The issue is in the 'main.css' file, specifically in the '.login-container' class. Need to adjust the media query."

# Update Implementation Plan

# Use the --plan flag to add or update the implementation plan for the task.
backlog backlog edit "42" --plan "1. Refactor login button\n2. Test on mobile\n3. Review with team"

# Set Single Dependency

# If you want to make a task depend on another specific task
backlog backlog edit "42" --dep "T15"

# Set Multiple Dependencies

# You can make a task depend on multiple other tasks
backlog backlog edit "42" --dep "T15,T18,T20"

# Complex Example with Multiple Changes

# You can combine several flags to make multiple changes at once.
backlog backlog edit "42" --assigned "alex" --priority "critical" --notes "The fix is ready for review. Please check on both iOS and Android." --check-ac "1,2" --status "in-review"
```

### Options

```
      --ac strings                Add a new acceptance criterion (can be used multiple times)
  -a, --assigned strings          Add assigned names for the task (can be used multiple times)
      --check-ac ints             Check an acceptance criterion by its index
      --dep strings               Set dependencies (can be used multiple times)
  -d, --description string        New description for the task
  -h, --help                      help for edit
  -l, --labels strings            Add labels for the task (can be used multiple times)
      --notes string              New implementation notes for the task
  -p, --parent string             New parent for the task
      --plan string               New implementation plan for the task
      --priority string           New priority for the task
      --remove-ac ints            Remove an acceptance criterion by its index
  -A, --remove-assigned strings   Assigned names to remove from the task (can be used multiple times)
  -L, --remove-labels strings     Labels to remove from the task (can be used multiple times)
  -s, --status string             New status for the task
  -t, --title string              New title for the task
      --uncheck-ac ints           Uncheck an acceptance criterion by its index
```

### Options inherited from parent commands

```
      --auto-commit         Auto-committing changes to git repository
      --folder string       Directory for backlog tasks (default ".backlog")
      --log-file string     Log file path (defaults to stderr)
      --log-format string   Log format (json, text) (default "text")
      --log-level string    Log level (debug, info, warn, error) (default "info")
      --max-limit int       Maximum limit for pagination (default 1000)
      --page-size int       Default page size for pagination (default 25)
```

### SEE ALSO

* [backlog](backlog.md)	 - backlog is a git-native, markdown-based task manager

