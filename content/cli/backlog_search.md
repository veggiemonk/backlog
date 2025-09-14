---
title: "backlog search"
description: "Search for tasks containing the specified query string."
weight: 6
---

## backlog search

Search tasks by content.

### Synopsis

Performs a full-text search for tasks containing the specified query string. The search covers titles, descriptions, comments, and other fields.

```
backlog search <query> [flags]
```

{{< expand "Show Examples" "true" >}}
```bash
# Search for tasks containing "login" in any field
backlog search "login"

# Search for tasks assigned to a specific person
backlog search "@john"

# Search for tasks with specific labels
backlog search "frontend"

# Search with markdown output
backlog search "api" --markdown

# Search with JSON output
backlog search "api" --json
```
{{< /expand >}}

### Options

In addition to the query, you can use the same filtering and sorting flags available in the `list` command.

| Flag | Description |
| --- | --- |
| `-a, --assigned` | Filter tasks by assigned names. |
| `-d, --depended-on` | Filter tasks that are depended on by other tasks. |
| `-c, --has-dependency` | Include tasks that have dependencies. |
| `-h, --help` | Help for search. |
| `-e, --hide-extra` | Hide extra fields (labels, priority, assigned). |
| `-j, --json` | Print JSON output. |
| `-l, --labels` | Filter tasks by labels. |
| `-m, --markdown` | Print markdown table. |
| `-p, --parent` | Filter tasks by parent ID. |
| `-r, --reverse` | Reverse the order of tasks. |
| `--sort` | Sort tasks by comma-separated fields. |
| `-s, --status` | Filter tasks by status. |
| `-u, --unassigned` | List tasks that have no one assigned. |

{{< expand "Inherited Options" >}}
| Flag | Default | Description |
| --- | --- | --- |
| `--auto-commit` | `true` | Auto-committing changes to git repository. |
| `--folder` | `.backlog` | Directory for backlog tasks. |
{{< /expand >}}

### SEE ALSO

- [**`backlog`**]({{< relref "backlog" >}}) - The main backlog command.
- [**`backlog list`**]({{< relref "backlog_list" >}}) - For more advanced filtering without a search query.