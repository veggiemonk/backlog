---
layout: page
title: backlog create
---

## backlog create

Create a new task

### Synopsis

Creates a new task in the backlog.

```
backlog create <title> [flags]
```

### Examples

```

# Create tasks using the "backlog create" command with its different flags.
# Here are some examples of how to use this command effectively:
# 1. Basic Task Creation
# This is the simplest way to create a task, providing only a title.
backlog create "Fix the login button styling"

# 2. Task with a Description. Use the -d or --description flag to add more detailed information about the task.
backlog create "Implement password reset" -d "Users should be able to request a password reset link via their email. This involves creating a new API endpoint and a front-end form."

# 3. Assigning a Task. You can assign a task to one or more team members using the -a or --assigned flag. 
# Assign to a single person: 
backlog create "Design the new dashboard" -a "alex"
# Assign to multiple people: 
backlog create "Code review for the payment gateway" -a "jordan" -a "casey"

# 4. Adding Labels. Use the -l or --labels flag to categorize the task with comma-separated labels.
backlog create "Update third-party dependencies" -l "bug,backend,security"

# 5. Setting a Priority
# Specify the task's priority with the --priority flag. The default is "medium".
backlog create "Hotfix: Production database is down" --priority "high"
backlog create "Refactor the old user model" --priority "low"

# 6. Defining Acceptance Criteria
# Use the --ac flag multiple times to list the conditions that must be met for the task to be considered complete.
backlog create "Develop user profile page" \
  --ac "Users can view their own profile information." \
  --ac "Users can upload a new profile picture." \
  --ac "The page is responsive on mobile devices."

# 7. Creating a Sub-task. Link a new task to a parent task using the -p or --parent flag. This is useful for breaking down larger tasks.
# First, create the parent task
backlog create "Implement User Authentication"
# Now, create a sub-task (assuming the parent task ID is 15)
backlog create "Add Google OAuth login" -p "15"

# 8. Setting Task Dependencies
# Use the --deps flag to specify that this task depends on other tasks being completed first.
# Single dependency:
backlog create "Deploy user authentication" --deps "T15"
# Multiple dependencies:
backlog create "Integration testing" --deps "T15" --deps "T18" --deps "T20"
# This means the task cannot be started until tasks T15, T18, and T20 are completed.

# 9. Complex Example (Combining Multiple Flags). Here is a comprehensive example that uses several flags at once to create a very detailed task.
backlog create "Build the new reporting feature" \
  -d "Create a new section in the app that allows users to generate and export monthly performance reports in PDF format." \
  -a "drew" \
  -l "feature,frontend,backend" \
  --priority "high" \
  --ac "Report generation logic is accurate." \
  --ac "Users can select a date range for the report." \
  --ac "The exported PDF has the correct branding and layout." \
  -p "23"	

```


### Options

#### Environment Variables

```
	(name)				(default)
	BACKLOG_AUTO_COMMIT	false
	BACKLOG_FOLDER		.backlog
	BACKLOG_LOG_FILE	
	BACKLOG_LOG_FORMAT	text
	BACKLOG_LOG_LEVEL	info
```

#### Flags


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
```

### SEE ALSO

* [backlog](backlog.md)	 - Backlog is a git-native, markdown-based task manager

