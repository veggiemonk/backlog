---
title: "backlog list"
description: "Lists all tasks in the backlog except archived tasks."
weight: 3
---

## backlog list

List all tasks.

### Synopsis

Lists all tasks in the backlog except for archived tasks. Provides powerful filtering and sorting options to help you find exactly what you're looking for.

```
backlog list [flags]
```

{{< expand "Show Examples" "true" >}}
{{< tabs "list-examples" >}}
{{< tab "Filtering" >}}
##### Filtering by Metadata
```bash
# List tasks with status "todo"
backlog list --status "todo"

# List tasks with status "todo" or "in-progress"
backlog list --status "todo,in-progress"

# List sub-tasks of task "T12"
backlog list --parent "T12"

# List tasks assigned to alice
backlog list --assigned "alice"

# List tasks that have no one assigned
backlog list --unassigned

# List tasks with the "bug" label
backlog list --labels "bug"

# List all high priority tasks
backlog list --priority "high"
```
{{< /tab >}}
{{< tab "Dependencies" >}}
##### Filtering by Dependencies
```bash
# List tasks that have at least one dependency
backlog list --has-dependency

# List tasks that are depended on by other tasks
backlog list --depended-on

# List all blocking tasks (are depended on and not done)
backlog list --depended-on --status "todo,in-progress"
```
{{< /tab >}}
{{< tab "Sorting & Output" >}}
##### Sorting
```bash
# Sort tasks by priority
backlog list --sort "priority"

# Sort by updated date, then priority
backlog list --sort "updated,priority"

# Sort by priority in reverse order
backlog list --sort "priority" --reverse
```

##### Output Formats
```bash
# List tasks in a compact table
backlog list -e

# List tasks in markdown format
backlog list -m

# List tasks in JSON format
backlog list -j
```
{{< /tab >}}
{{< /tabs >}}
{{< /expand >}}

### Options

| Flag | Description |
| --- | --- |
| `-a, --assigned` | Filter tasks by assigned names. |
| `-d, --depended-on` | Filter tasks that are depended on by other tasks. |
| `-c, --has-dependency` | Filter tasks that have dependencies. |
| `-h, --help` | Help for list. |
| `-e, --hide-extra` | Hide extra fields (labels, priority, assigned). |
| `-j, --json` | Print JSON output. |
| `-l, --labels` | Filter tasks by labels. |
| `-m, --markdown` | Print markdown table. |
| `-p, --parent` | Filter tasks by parent ID. |
| `--priority` | Filter tasks by priority. |
| `-r, --reverse` | Reverse the order of tasks. |
| `--sort` | Sort tasks by comma-separated fields (id, title, status, priority, created, updated). |
| `-s, --status` | Filter tasks by status. |
| `-u, --unassigned` | Filter tasks that have no one assigned. |

{{< expand "Inherited Options" >}}
| Flag | Default | Description |
| --- | --- | --- |
| `--auto-commit` | `true` | Auto-committing changes to git repository. |
| `--folder` | `.backlog` | Directory for backlog tasks. |
{{< /expand >}}

### SEE ALSO

- [**`backlog`**]({{< relref "backlog" >}}) - The main backlog command.