## backlog create

Create a new task

### Synopsis

Creates a new task in the backlog.

```
backlog create <title> [flags]
```

### Examples

```
# Basic Task Creation

# This is the simplest way to create a task, providing only a title.
backlog create "Fix the login button styling"

# Task with Description

# Use the -d or --description flag to add more detailed information about the task.
backlog create "Implement password reset" --description "Users should be able to request a password reset link via their email. This involves creating a new API endpoint and a front-end form."

# Assigning to Single Person

# You can assign a task to one or more team members using the -a or --assigned flag.
backlog create "Design the new dashboard" --assigned "alex"

# Assigning to Multiple People

backlog create "Code review for the payment gateway" --assigned "jordan,casey"

# Adding Labels

# Use the -l or --labels flag to categorize the task with comma-separated labels.
backlog create "Update third-party dependencies" --labels "bug,backend,security"

# Setting High Priority

# Specify the task's priority with the --priority flag. The default is "medium".
backlog create "Hotfix: Production database is down" --priority "high"

# Setting Low Priority

backlog create "Refactor the old user model" --priority "low"

# Defining Acceptance Criteria

# Use the --ac flag multiple times to list the conditions that must be met for the task to be considered complete.
backlog create "Develop user profile page" --ac "Users can view their own profile information.,Users can upload a new profile picture.,The page is responsive on mobile devices."

# Creating a Sub-task

# Link a new task to a parent task using the -p or --parent flag. This is useful for breaking down larger tasks.
backlog create "Add Google OAuth login" --parent "15"

# Setting Single Dependency

# Use the --deps flag to specify that this task depends on other tasks being completed first.
backlog create "Deploy user authentication" --deps "T15"

# Setting Multiple Dependencies

# This means the task cannot be started until tasks T15, T18, and T20 are completed.
backlog create "Integration testing" --deps "T15,T18,T20"

# Task with Implementation Notes

# Use the --notes flag to add implementation notes to help with development.
backlog create "Optimize database queries" --notes "Focus on the user lookup queries in the authentication module. Consider adding indexes on email and username fields."

# Task with Implementation Plan

# Use the --plan flag to add a structured implementation plan.
backlog create "Implement user registration flow" --plan "1. Design registration form UI\n2. Create user validation logic\n3. Set up email verification\n4. Add password strength requirements\n5. Write integration tests"

# Complex Example with Multiple Flags

# Here is a comprehensive example that uses several flags at once to create a very detailed task.
backlog create "Build the new reporting feature" --priority "high" --ac "Report generation logic is accurate.,Users can select a date range for the report.,The exported PDF has the correct branding and layout." --parent "23" --description "Create a new section in the app that allows users to generate and export monthly performance reports in PDF format." --assigned "drew" --labels "feature,frontend,backend"
```

### Options

```
      --ac strings           Acceptance criterion (can be specified multiple times)
  -a, --assigned strings     Assignee for the task (can be specified multiple times)
      --deps strings         Add a dependency (can be used multiple times)
  -d, --description string   Description of the task
  -h, --help                 help for create
  -l, --labels strings       Comma-separated labels for the task
      --notes string         Additional notes for the task
  -p, --parent string        Parent task ID
      --plan string          Implementation plan for the task
      --priority string      Priority of the task (low, medium, high, critical) (default "medium")
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

