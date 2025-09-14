---
title: "backlog edit"
description: "Edit an existing task by providing its ID and flags for the fields to update."
weight: 4
---

## backlog edit

Edit an existing task.

### Synopsis

Edit an existing task by providing its ID and flags for the fields to update. This allows for partial updates, so you only change what you need.

```
backlog edit <id> [flags]
```

{{< expand "Show Examples" "true" >}}
{{< tabs "edit-examples" >}}
{{< tab "Basic Edits" >}}
##### 1. Changing the Title
Use `-t` or `--title` to rename the task.
```bash
backlog edit T42 -t "Fix the main login button styling"
```

##### 2. Updating the Description
Use `-d` or `--description` to replace the description.
```bash
backlog edit T42 -d "The login button is misaligned on mobile. It should be centered."
```

##### 3. Changing the Status
Update progress with `-s` or `--status`.
```bash
backlog edit T42 -s "in-progress"
```
{{< /tab >}}
{{< tab "Metadata Edits" >}}
##### 4. Re-assigning a Task
Use `-a` or `--assigned` to replace the current assignees.
```bash
# Assign to a single person
backlog edit T42 -a "jordan"

# Assign to multiple people
backlog edit T42 -a "jordan" -a "casey"
```

##### 5. Updating Labels
Use `-l` or `--labels` to replace the existing labels.
```bash
backlog edit T42 -l "bug,frontend"
```

##### 6. Changing the Priority
Adjust priority with the `--priority` flag.
```bash
backlog edit T42 --priority "high"
```
{{< /tab >}}
{{< tab "Advanced Edits" >}}
##### 7. Managing Acceptance Criteria
Add, check, uncheck, or remove acceptance criteria.
```bash
# Add a new AC
backlog edit T42 --ac "Button is centered on screens < 576px."

# Check the first AC (1-based index)
backlog edit T42 --check-ac 1

# Uncheck the first AC
backlog edit T42 --uncheck-ac 1

# Remove the second AC
backlog edit T42 --remove-ac 2
```

##### 8. Changing the Parent Task
Move a task under a different parent with `-p` or `--parent`.
```bash
backlog edit T42 -p "T18"

# To remove a parent, pass an empty string
backlog edit T42 -p ""
```

##### 9. Updating Dependencies
Use `--dep` to replace existing dependencies.
```bash
backlog edit T42 --dep "T1" --dep "T2"
```
{{< /tab >}}
{{< tab "Complex Edit" >}}
##### 10. Combining Multiple Flags
Make multiple changes at once.
```bash
backlog edit T42 \
  -s "in-review" \
  -a "alex" \
  --priority "critical" \
  --notes "The fix is ready for review. Please check on iOS and Android." \
  --check-ac 1 \
  --check-ac 2
```
{{< /tab >}}
{{< /tabs >}}
{{< /expand >}}

### Options

| Flag | Description |
| --- | --- |
| `--ac` | Add a new acceptance criterion (can be used multiple times). |
| `-a, --assigned` | Add assigned names for the task (can be used multiple times). |
| `--check-ac` | Check an acceptance criterion by its 1-based index. |
| `--dep` | Set dependencies (replaces existing, can be used multiple times). |
| `-d, --description` | New description for the task. |
| `-h, --help` | Help for edit. |
| `-l, --labels` | Add labels for the task (can be used multiple times). |
| `--notes` | New implementation notes for the task. |
| `-p, --parent` | New parent for the task. |
| `--plan` | New implementation plan for the task. |
| `--priority` | New priority for the task. |
| `--remove-ac` | Remove an acceptance criterion by its 1-based index. |
| `-A, --remove-assigned` | Assigned names to remove from the task. |
| `-L, --remove-labels` | Labels to remove from the task. |
| `-s, --status` | New status for the task. |
| `-t, --title` | New title for the task. |
| `--uncheck-ac` | Uncheck an acceptance criterion by its 1-based index. |

{{< expand "Inherited Options" >}}
| Flag | Default | Description |
| --- | --- | --- |
| `--auto-commit` | `true` | Auto-committing changes to git repository. |
| `--folder` | `.backlog` | Directory for backlog tasks. |
{{< /expand >}}

### SEE ALSO

- [**`backlog`**]({{< relref "backlog" >}}) - The main backlog command.