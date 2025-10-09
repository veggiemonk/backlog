---
layout: page
title: backlog create
---

# NAME

create - Create a new task

# SYNOPSIS

create

```
[--ac]=[value]
[--assigned|-a]=[value]
[--deps]=[value]
[--description|-d]=[value]
[--labels|-l]=[value]
[--notes]=[value]
[--parent|-p]=[value]
[--plan]=[value]
[--priority]=[value]
```

# DESCRIPTION


Creates a new task in the backlog.

Examples:
```

  # Basic task
  backlog create "Fix the login button styling"

  # Task with description
  backlog create "Implement password reset" \
    -d "Users should be able to request a password reset link via email"

  # Assign to team members
  backlog create "Design the new dashboard" -a "alex"
  backlog create "Code review" -a "jordan" -a "casey"    # Multiple assignees

  # Add labels
  backlog create "Update dependencies" -l "bug,backend,security"

  # Set priority (low, medium, high, critical)
  backlog create "Hotfix: Database down" --priority "high"
  backlog create "Refactor old code" --priority "low"

  # Define acceptance criteria
  backlog create "Develop user profile page" \
    --ac "Users can view their profile" \
    --ac "Users can upload a profile picture" \
    --ac "Page is responsive on mobile"

  # Create sub-task with parent
  backlog create "Implement User Authentication"           # Creates T01
  backlog create "Add Google OAuth" -p "T01"               # Creates T01.01

  # Set dependencies
  backlog create "Deploy to production" --deps "T15"       # Single dependency
  backlog create "Integration testing" \
    --deps "T15" --deps "T18" --deps "T20"                 # Multiple dependencies

  # Complex example with multiple flags
  backlog create "Build reporting feature" \
    -d "Monthly performance reports in PDF format" \
    -a "drew" \
    -l "feature,frontend,backend" \
    --priority "high" \
    --ac "Report generation logic is accurate" \
    --ac "Users can select date range" \
    --ac "PDF has correct branding" \
    -p "23"

```

**Usage**:

```
backlog create <title>
```

# GLOBAL OPTIONS

**--ac**="": Acceptance criterion (can be specified multiple times) (default: [])

**--assigned, -a**="": Assignee for the task (can be specified multiple times) (default: [])

**--deps**="": Add a dependency (can be used multiple times) (default: [])

**--description, -d**="": Description of the task

**--labels, -l**="": Comma-separated labels for the task (default: [])

**--notes**="": Additional notes for the task

**--parent, -p**="": Parent task ID

**--plan**="": Implementation plan for the task

**--priority**="": Priority of the task (low, medium, high, critical) (default: medium)

