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

# Edit tasks using the "backlog edit" command with its different flags.
# Let's assume you have a task with ID "42" that you want to modify.
# Here are some examples of how to use this command effectively:

# 1. Changing the Title
# Use the -t or --title flag to give the task a new title.
backlog edit 42 -t "Fix the main login button styling"

# 2. Updating the Description
# Use the -d or --description flag to replace the existing description with a new one.
backlog edit 42 -d "The login button on the homepage is misaligned on mobile devices. It should be centered."

# 3. Changing the Status
# Update the task's progress by changing its status with the -s or --status flag.
backlog edit 42 -s "in-progress"

# 4. Re-assigning a Task
# You can change the assigned names for a task using the -a or --assignee flag.
# This will replace the current list of assigned names.
# Assign to a single person:
backlog edit 42 -a "jordan"
# Assign to multiple people:
backlog edit 42 -a "jordan" -a "casey"

# 5. Updating Labels
# Use the -l or --labels flag to replace the existing labels.
backlog edit 42 -l "bug,frontend"

# 6. Changing the Priority
# Adjust the task's priority with the --priority flag.
backlog edit 42 --priority "high"

# 7. Managing Acceptance Criteria
# You can add, check, uncheck, or remove acceptance criteria.
# Add a new AC:
backlog edit 42 --ac "The button is centered on screens smaller than 576px."
# Check the first AC (assuming it's at index 1):
backlog edit 42 --check-ac 1
# Uncheck the first AC:
backlog edit 42 --uncheck-ac 1
# Remove the second AC (at index 2):
backlog edit 42 --remove-ac 2

# 8. Changing the Parent Task
# Move a task to be a sub-task of a different parent using the -p or --parent flag.
backlog edit 42 -p "18"
# To remove a parent, pass an empty string:
backlog edit 42 -p ""

# 9. Adding Implementation Notes
# Use the --notes flag to add or update technical notes for implementation.
backlog edit 42 --notes "The issue is in the 'main.css' file, specifically in the '.login-container' class. Need to adjust the media query."

# 10. Complex Example (Combining Multiple Flags)
# You can combine several flags to make multiple changes at once.
backlog edit 42 \
  -s "in-review" \
  -a "alex" \
  --priority "critical" \
  --notes "The fix is ready for review. Please check on both iOS and Android." \
  --check-ac 1 \
  --check-ac 2

# 11. Updating the Implementation Plan
# Use the --plan flag to add or update the implementation plan for the task.
backlog edit 42 --plan "1. Refactor login button\\n2. Test on mobile\\n3. Review with team"

# 12. Adding Dependencies
# Use the --deps flag to add one or more task dependencies.
# This will replace all existing dependencies with the new ones.
backlog edit 42 --deps "T1" --deps "T2"

# 13. Setting a Single Dependency
# If you want to make a task depend on another specific task:
backlog edit 42 --deps "T15"
# This makes task 42 dependent on task T15, meaning T15 must be completed before T42 can be started.

# 14. Setting Multiple Dependencies
# You can make a task depend on multiple other tasks:
backlog edit 42 --deps "T15" --deps "T18" --deps "T20"
# This makes task 42 dependent on tasks T15, T18, and T20.
# 15. Editing the construction plan
backlog edit 42 --plan "1. Dig hole 2. Pour foundation"


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

