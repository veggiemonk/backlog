---
layout: page
title: backlog doctor
---

# NAME

doctor - Diagnose and fix task ID conflicts

# SYNOPSIS

doctor

```
[--dry-run]
[--fix]
[--json|-j]
[--strategy]=[value]
```

# DESCRIPTION

Diagnose and fix task ID conflicts that can occur when creating tasks
in separate Git branches. Conflicts arise when multiple branches generate the same task IDs.

This command provides conflict detection and resolution capabilities to maintain
task ID uniqueness and data integrity.

Conflict types detected:
- Duplicate IDs (same ID in multiple files)
- Orphaned children (tasks with non-existent parents)
- Invalid hierarchy (parent-child ID mismatch)


Examples:
  backlog doctor                    # Detect conflicts in text format
  backlog doctor --json             # Detect conflicts in JSON format
  backlog doctor --fix              # Detect and automatically fix conflicts
  backlog doctor --fix --dry-run    # Show what would be fixed without making changes
  backlog doctor --fix --strategy=auto    # Use auto-renumbering strategy


**Usage**:

```
doctor [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--dry-run**: Show what would be changed without making changes (use with --fix)

**--fix**: Automatically fix detected conflicts

**--json, -j**: Output in JSON format

**--strategy**="": Resolution strategy when using --fix (chronological|auto|manual) (default: chronological)

