---
title: "backlog view"
description: "View a task by providing its ID in markdown or JSON format."
weight: 5
---

## backlog view

View a task by providing its ID.

### Synopsis

Displays the full details of a single task, including its title, description, metadata, and acceptance criteria. You can choose to output in standard markdown or JSON format.

```
backlog view <id> [flags]
```

### Examples

```bash
# View task T01 in the default markdown format
backlog view T01

# View task T01 in JSON format
backlog view T01 --json

# View task T01 in JSON format (short flag)
backlog view T01 -j
```

### Options

| Flag | Description |
| --- | --- |
| `-h, --help` | Help for view. |
| `-j, --json` | Print JSON output. |

{{< expand "Inherited Options" >}}
| Flag | Default | Description |
| --- | --- | --- |
| `--auto-commit` | `true` | Auto-committing changes to git repository. |
| `--folder` | `.backlog` | Directory for backlog tasks. |
{{< /expand >}}

### SEE ALSO

- [**`backlog`**]({{< relref "backlog" >}}) - The main backlog command.