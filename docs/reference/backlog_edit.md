---
layout: page
title: backlog edit
---

# NAME

edit - Edit an existing task

# SYNOPSIS

edit

```
[--ac]=[value]
[--assigned|-a]=[value]
[--check-ac]=[value]
[--deps]=[value]
[--description|-d]=[value]
[--labels|-l]=[value]
[--notes]=[value]
[--parent|-p]=[value]
[--plan]=[value]
[--priority]=[value]
[--remove-ac]=[value]
[--remove-assigned|-A]=[value]
[--remove-labels|-L]=[value]
[--status|-s]=[value]
[--title|-t]=[value]
[--uncheck-ac]=[value]
```

# DESCRIPTION


Edit an existing task by providing its ID and flags for the fields to update.

Examples:
```

  # Change title
  backlog edit T42 -t "Fix the main login button"

  # Update description
  backlog edit T42 -d "The login button is misaligned on mobile"

  # Change status
  backlog edit T42 -s "in-progress"
  backlog edit T42 -s "done"

  # Assign/re-assign users
  backlog edit T42 -a "jordan"                    # Add jordan
  backlog edit T42 -a "jordan" -a "casey"         # Add multiple users
  backlog edit T42 --remove-assigned "alex"       # Remove alex

  # Update labels
  backlog edit T42 -l "bug,frontend"              # Add labels
  backlog edit T42 --remove-labels "backend"      # Remove label

  # Change priority
  backlog edit T42 --priority "high"

  # Manage acceptance criteria
  backlog edit T42 --ac "Button centered on mobile"     # Add new AC
  backlog edit T42 --check-ac 1                         # Check AC #1
  backlog edit T42 --uncheck-ac 1                       # Uncheck AC #1
  backlog edit T42 --remove-ac 2                        # Remove AC #2

  # Change parent
  backlog edit T42 -p "T18"                       # Make it a sub-task of T18
  backlog edit T42 -p ""                          # Remove parent

  # Add implementation notes
  backlog edit T42 --notes "Fixed in main.css, line 234"

  # Update implementation plan
  backlog edit T42 --plan "1. Refactor login\\n2. Test on mobile\\n3. Review"

  # Set dependencies
  backlog edit T42 --deps "T15"                   # Single dependency
  backlog edit T42 --deps "T15" --deps "T18"      # Multiple dependencies

  # Complex example (multiple changes at once)
  backlog edit T42 \
    -s "in-review" \
    -a "alex" \
    --priority "critical" \
    --notes "Ready for review on iOS and Android" \
    --check-ac 1 \
    --check-ac 2

```

**Usage**:

```
backlog edit <id>
```

# GLOBAL OPTIONS

**--ac**="": Add a new acceptance criterion (can be used multiple times) (default: [])

**--assigned, -a**="": Add assigned names for the task (can be used multiple times) (default: [])

**--check-ac**="": Check an acceptance criterion by its index (default: [])

**--deps**="": Set dependencies (can be used multiple times) (default: [])

**--description, -d**="": New description for the task

**--labels, -l**="": Add labels for the task (can be used multiple times) (default: [])

**--notes**="": New implementation notes for the task

**--parent, -p**="": New parent for the task

**--plan**="": New implementation plan for the task

**--priority**="": New priority for the task

**--remove-ac**="": Remove an acceptance criterion by its index (default: [])

**--remove-assigned, -A**="": Assigned names to remove from the task (can be used multiple times) (default: [])

**--remove-labels, -L**="": Labels to remove from the task (can be used multiple times) (default: [])

**--status, -s**="": New status for the task

**--title, -t**="": New title for the task

**--uncheck-ac**="": Uncheck an acceptance criterion by its index (default: [])

